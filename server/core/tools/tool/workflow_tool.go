package tool

import (
	"encoding/json"
	"fmt"
	"maps"

	wfgenerator "mlib.com/gofy/server/core/app/generatores/workflow"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/file"
	dbengine "mlib.com/gofy/server/db_engine"
	toolsentities "mlib.com/gofy/server/entities/tools"
	toolsenumtypes "mlib.com/gofy/server/enum_types/tools"
	filefactory "mlib.com/gofy/server/factories/file_factory"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type WorkflowTool struct {
	*Tool
	WorkflowAppID     string         `json:"workflow_app_id"`
	Version           string         `json:"version"`
	WorkflowEntities  map[string]any `json:"workflow_entities"`
	WorkflowCallDepth int            `json:"workflow_call_depth"`
	ThreadPoolID      string         `json:"thread_pool_id"`
	Label             string         `json:"label"`
}

func (t *WorkflowTool) _get_app(app_id string) *models.App {
	app := new(models.App)
	err := dbengine.Instance().DB.Model(&models.App{}).Where("id = ?", app_id).First(app).Error
	if err != nil {
		mlog.Errorf("get App from mysql failed:%v", err)
		app = nil
	}

	if app == nil {
		panic(exceptions.NewValueError("app not found"))
	}
	return app
}
func (t *WorkflowTool) _get_user(user_id string) (*models.EndUser, *models.Account) {
	end_user := new(models.EndUser)
	err := dbengine.Instance().DB.Model(&models.EndUser{}).Where("id = ?", user_id).First(end_user).Error
	if err != nil {
		mlog.Errorf("get EndUser from mysql failed:%v", err)
		end_user = nil
	}
	if end_user == nil {
		acc_user := new(models.Account)
		err := dbengine.Instance().DB.Model(&models.Account{}).Where("id = ?", user_id).First(acc_user).Error
		if err != nil {
			mlog.Errorf("get Account from mysql failed:%v", err)
			acc_user = nil
			panic(exceptions.NewValueError("user not found"))
		} else {
			return nil, acc_user
		}
	} else {
		return end_user, nil
	}
}
func (t *WorkflowTool) _get_workflow(app_id string, version string) *models.Workflow {
	var workflow *models.Workflow
	if version == "" {
		workflow = new(models.Workflow)
		err := dbengine.Instance().DB.Model(&models.Workflow{}).Where("app_id = ? and version != ?", app_id, "draft").Order("created_at DESC").First(workflow).Error
		if err != nil {
			mlog.Errorf("get Workflow from mysql failed:%v", err)
			workflow = nil
		}
	} else {
		workflow = new(models.Workflow)
		err := dbengine.Instance().DB.Model(&models.Workflow{}).Where("app_id = ? and version = ?", app_id, version).First(workflow).Error
		if err != nil {
			mlog.Errorf("get Workflow from mysql failed:%v", err)
			workflow = nil
		}
	}
	if workflow == nil {
		panic(exceptions.NewValueError("workflow not found or not published"))
	}
	return workflow
}
func (t *WorkflowTool) _extract_files(outputs map[string]any) (map[string]any, []*file.File) {
	files := []*file.File{}
	result := map[string]any{}
	for key, value := range outputs {
		if _, ok := value.([]any); ok {
			for _, item := range value.([]any) {
				if item_dict, ok := item.(map[string]any); ok {
					if _, ok := item_dict["dify_model_identity"]; ok {
						if _, ok := item_dict["dify_model_identity"].(string); ok {
							if item_dict["dify_model_identity"].(string) == file.FILE_MODEL_IDENTITY {
								item_dict["tool_file_id"] = item_dict["related_id"]
								file := filefactory.BuildFromMapping(item_dict, t.Runtime.TenantID, nil)
								files = append(files, file)
							}
						}
					}

				}
			}
		} else if _, ok := value.(map[string]any); ok {
			if _, ok := value.(map[string]any)["dify_model_identity"]; ok {
				if _, ok := value.(map[string]any)["dify_model_identity"].(string); ok {
					if value.(map[string]any)["dify_model_identity"].(string) == file.FILE_MODEL_IDENTITY {
						value.(map[string]any)["tool_file_id"] = value.(map[string]any)["related_id"]
						file := filefactory.BuildFromMapping(value.(map[string]any), t.Runtime.TenantID, nil)
						files = append(files, file)
					}
				}
			}
		}
		result[key] = value
	}
	return result, files
}

func (t *WorkflowTool) _transform_args(tool_parameters map[string]any) (map[string]any, []map[string]any) {
	parameter_rules := t.GetAllRuntimeParameters()
	parameters_result := map[string]any{}
	files := []map[string]any{}
	for _, parameter := range parameter_rules {
		if parameter.Type == toolsenumtypes.ToolParameter_SYSTEM_FILES {
			if _, ok := tool_parameters[parameter.Name]; ok {
				if _, ok := tool_parameters[parameter.Name].([]*file.File); ok {
					for _, param_file := range tool_parameters[parameter.Name].([]*file.File) {
						file_dict := map[string]any{
							"transfer_method": string(param_file.TransferMethod),
							"type":            string(param_file.Type),
						}
						if param_file.TransferMethod == file.FileTransferMethod_TOOL_FILE {
							file_dict["tool_file_id"] = param_file.RelatedID
						} else if param_file.TransferMethod == file.FileTransferMethod_LOCAL_FILE {
							file_dict["upload_file_id"] = param_file.RelatedID
						} else if param_file.TransferMethod == file.FileTransferMethod_REMOTE_URL {
							file_dict["url"] = param_file.GenerateURL()
						}
						files = append(files, file_dict)
					}
				} else {
					mlog.Errorf("tool_parameters[%s]=%#v is not file list", parameter.Name, tool_parameters[parameter.Name])
				}
			} else {
				mlog.Errorf("parameter.Name=%s not exist in tool_parameters", parameter.Name)
			}
		} else {
			parameters_result[parameter.Name] = tool_parameters[parameter.Name]
		}
	}
	return parameters_result, files
}

func (t *WorkflowTool) ToolProviderType() toolsenumtypes.ToolProviderType {
	return toolsenumtypes.ToolProvider_WORKFLOW
}
func (t *WorkflowTool) RealInvoke(user_id string, tool_parameters map[string]any) (*toolsentities.ToolInvokeMessage, []*toolsentities.ToolInvokeMessage) {
	if t.Runtime == nil {
		panic(exceptions.NewValueError("workflow tool runtime is nil"))
	}
	t.Runtime.InvokeFrom.Validate()

	app := t._get_app(t.WorkflowAppID)
	wf := t._get_workflow(t.WorkflowAppID, t.Version)

	// transform the tool parameters
	tool_parameters, files := t._transform_args(tool_parameters)
	end_user, acc_user := t._get_user(user_id)
	var result map[string]any
	if end_user != nil {
		generator := wfgenerator.New[*models.EndUser]()
		result, _ = generator.Generate(
			app,
			wf,
			end_user,
			map[string]any{"inputs": tool_parameters, "files": files},
			t.Runtime.InvokeFrom,
			false,
			t.WorkflowCallDepth+1,
		)
	} else if acc_user != nil {
		generator := wfgenerator.New[*models.Account]()
		result, _ = generator.Generate(
			app,
			wf,
			acc_user,
			map[string]any{"inputs": tool_parameters, "files": files},
			t.Runtime.InvokeFrom,
			false,
			t.WorkflowCallDepth+1,
		)
	}

	data := map[string]any{}
	if _, ok := result["data"]; ok {
		if _, ok := result["data"].(map[string]any); ok {
			data = result["data"].(map[string]any)
		}
	}
	if _, ok := data["error"]; ok {
		panic(exceptions.NewValueError(fmt.Sprintf("%v", data["error"])))
	}

	r := []*toolsentities.ToolInvokeMessage{}
	outputs := map[string]any{}
	if _, ok := data["outputs"]; ok {
		if _, ok := data["outputs"].(map[string]any); ok {
			outputs = data["outputs"].(map[string]any)
		}
	}
	if len(outputs) > 0 {
		var extracted_files []*file.File
		outputs, extracted_files = t._extract_files(outputs)
		for _, f := range extracted_files {
			r = append(r, t.CreateFileMessage(f))
		}
	}
	bindata, _ := json.Marshal(outputs)
	r = append(r, t.CreateTextMessage(string(bindata), ""))
	r = append(r, t.CreateJsonMessage(outputs))

	return nil, r
}

func (t *WorkflowTool) fork_tool_runtime(runtime map[string]any) *WorkflowTool {
	return &WorkflowTool{
		Tool:              t.ForkToolRuntime(runtime),
		WorkflowAppID:     t.WorkflowAppID,
		Version:           t.Version,
		WorkflowEntities:  maps.Clone(t.WorkflowEntities),
		WorkflowCallDepth: t.WorkflowCallDepth,
		ThreadPoolID:      t.ThreadPoolID,
		Label:             t.Label,
	}
}
