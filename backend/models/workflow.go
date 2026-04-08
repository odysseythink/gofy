package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"mlib.com/confy/cast"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/core/variables"
	dbengine "mlib.com/gofy/server/db_engine"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	variablefactory "mlib.com/gofy/server/factories/variable_factory"
	"mlib.com/gofy/server/libs/helper"
	"mlib.com/gofy/server/utils/validate"
	"mlib.com/mlog"
)

type WorkflowAppLogCreatedFrom string

const (
	/*
	   Workflow App Log Created From Enum
	*/

	WorkflowAppLogCreatedFrom_SERVICE_API   = "service-api"
	WorkflowAppLogCreatedFrom_WEB_APP       = "web-app"
	WorkflowAppLogCreatedFrom_INSTALLED_APP = "installed-app"

	// @classmethod
	// def value_of(cls, value: str) -> "WorkflowAppLogCreatedFrom":
	//     """
	//     Get value of given mode.

	//     :param value: mode value
	//     :return: mode
	//     """
	//     for mode in cls:
	//         if mode.value == value:
	//             return mode
	//     raise ValueError(f"invalid workflow app log created from value {value}")
)

type WorkflowNodeExecutionTriggeredFrom string

const (
	/*
	   Workflow Node Execution Triggered From Enum
	*/

	WorkflowNodeExecutionTriggeredFrom_SINGLE_STEP  WorkflowNodeExecutionTriggeredFrom = "single-step"
	WorkflowNodeExecutionTriggeredFrom_WORKFLOW_RUN WorkflowNodeExecutionTriggeredFrom = "workflow-run"

	// @classmethod
	// def value_of(cls, value: str) -> "WorkflowNodeExecutionTriggeredFrom":
	//     """
	//     Get value of given mode.

	//     :param value: mode value
	//     :return: mode
	//     """
	//     for mode in cls:
	//         if mode.value == value:
	//             return mode
	//     raise ValueError(f"invalid workflow node execution triggered from value {value}")
)

type CreatedByRole string

const (
	CreatedByRole_ACCOUNT  CreatedByRole = "account"
	CreatedByRole_END_USER CreatedByRole = "end_user"
)

func (e CreatedByRole) Valid() bool {
	return e == CreatedByRole_ACCOUNT ||
		e == CreatedByRole_END_USER
}

type UserFrom string

const (
	UserFrom_ACCOUNT  UserFrom = "account"
	UserFrom_END_USER UserFrom = "end-user"
)

type WorkflowRunTriggeredFrom string

const (
	WorkflowRunTriggeredFrom_DEBUGGING WorkflowRunTriggeredFrom = "debugging"
	WorkflowRunTriggeredFrom_APP_RUN   WorkflowRunTriggeredFrom = "app-run"
)

type WorkflowType string

const (
	/*
	   Workflow Type Enum
	*/

	Workflow_WORKFLOW WorkflowType = "workflow"
	Workflow_CHAT     WorkflowType = "chat"
)

func (wft WorkflowType) from_app_mode(app_mode AppMode) WorkflowType {
	// """
	// Get workflow type from app mode.

	// :param app_mode: app mode
	// :return: workflow type
	// """
	if app_mode == AppMode_WORKFLOW {
		return Workflow_WORKFLOW
	} else {
		return Workflow_CHAT
	}
}

type WorkflowNodeExecutionStatus string

const (
	/*
	   Workflow Node Execution Status Enum
	*/

	WorkflowNodeExecutionStatus_RUNNING   WorkflowNodeExecutionStatus = "running"
	WorkflowNodeExecutionStatus_SUCCEEDED WorkflowNodeExecutionStatus = "succeeded"
	WorkflowNodeExecutionStatus_FAILED    WorkflowNodeExecutionStatus = "failed"
	WorkflowNodeExecutionStatus_EXCEPTION WorkflowNodeExecutionStatus = "exception"
	WorkflowNodeExecutionStatus_RETRY     WorkflowNodeExecutionStatus = "retry"

	// @classmethod
	// def value_of(cls, value: str) -> "WorkflowNodeExecutionStatus":
	//     """
	//     Get value of given mode.

	//     :param value: mode value
	//     :return: mode
	//     """
	//     for mode in cls:
	//         if mode.value == value:
	//             return mode
	//     raise ValueError(f"invalid workflow node execution status value {value}")
)

// Workflow [...]
type Workflow struct {
	ID                       string       `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID                 string       `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Tenant                   *Tenant      `json:"tenant" form:"tenant" gorm:"foreignKey:TenantID;references:ID;"`
	AppID                    string       `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App                      *App         `json:"app" form:"app" gorm:"foreignKey:AppID;references:ID;"`
	Type                     WorkflowType `gorm:"column:type;type:varchar(255);not null" json:"type"`
	Version                  string       `gorm:"column:version;type:varchar(255);not null" json:"version"`
	Graph                    string       `gorm:"column:graph;type:text" json:"graph"`
	FeaturesStr              string       `gorm:"column:features;type:text" json:"features"`
	CreatedBy                string       `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedByAccount         *Account     `json:"created_by_account" form:"created_by_account" gorm:"foreignKey:CreatedBy;references:ID;"`
	CreatedAt                *time.Time   `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedBy                string       `gorm:"column:updated_by;type:varchar(36)" json:"updated_by"`
	UpdatedByAccount         *Account     `json:"updated_by_account" form:"updated_by_account" gorm:"foreignKey:UpdatedBy;references:ID;"`
	UpdatedAt                *time.Time   `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
	EnvironmentVariablesStr  string       `gorm:"column:environment_variables;type:varchar(255);default:{}" json:"environment_variables"`
	ConversationVariablesStr string       `gorm:"column:conversation_variables;type:varchar(255);default:{}" json:"conversation_variables"`
}

// TableName get sql table name.获取数据库表名
func (Workflow) TableName() string {
	return "workflows"
}

func (wf *Workflow) GraphDict() map[string]any {
	if wf.Graph == "" {
		return nil
	} else {
		var tmp map[string]any
		err := json.Unmarshal([]byte(wf.Graph), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", wf.Graph, err)
			return nil
		}
		return tmp
	}
}

func (wf *Workflow) Features() string {
	if wf.FeaturesStr != "" {
		return wf.FeaturesStr
	} else {
		var tmp map[string]any
		err := json.Unmarshal([]byte(wf.FeaturesStr), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", wf.Graph, err)
			return ""
		}
		if _, ok := tmp["file_upload"]; ok {
			if file_upload, ok := tmp["file_upload"].(map[string]any); ok && file_upload != nil {
				if _, ok := file_upload["image"]; ok {
					if image, ok := file_upload["image"].(map[string]any); ok && image != nil {
						if _, ok := image["enabled"]; ok {
							if enabled, ok := file_upload["enabled"].(bool); ok && enabled {
								image_enabled := true
								image_number_limits := 1
								if _, ok := image["number_limits"]; ok {
									val, err := cast.ToE[int](image["number_limits"])
									if err == nil {
										image_number_limits = val
									}
								}
								image_transfer_methods := []string{"remote_url", "local_file"}
								if _, ok := image["transfer_methods"]; ok {
									bindata, err := json.Marshal(image["transfer_methods"])
									if err == nil {
										tmplist := []string{}
										err := json.Unmarshal(bindata, &tmplist)
										if err == nil {
											image_transfer_methods = tmplist
										}
									}
								}

								file_upload["enabled"] = image_enabled
								file_upload["number_limits"] = image_number_limits
								file_upload["allowed_file_upload_methods"] = image_transfer_methods
								file_upload["allowed_file_types"] = []string{"image"}
								file_upload["allowed_file_extensions"] = []any{}
								delete(file_upload, "image")
								tmp["file_upload"] = file_upload
							}
						}
					}
				}
			}
		}
		bindata, _ := json.Marshal(tmp)
		return string(bindata)
	}
}

func (wf *Workflow) FeaturesDict() map[string]any {
	if wf.FeaturesStr == "" {
		return nil
	} else {
		var tmp map[string]any
		err := json.Unmarshal([]byte(wf.FeaturesStr), &tmp)
		if err != nil {
			mlog.Errorf("json unmarshal(%s) failed:%v", wf.FeaturesStr, err)
			return nil
		}
		return tmp
	}
}

func (wf *Workflow) UserInputForm(to_old_structure bool) []any {
	// # get start node from graph
	if wf.Graph == "" {
		return nil
	}

	graph_dict := wf.GraphDict()
	if _, ok := graph_dict["nodes"]; !ok {
		return nil
	}
	if _, ok := graph_dict["nodes"].([]any); !ok {
		return nil
	}
	var start_node map[string]any
	for _, node := range graph_dict["nodes"].([]any) {
		if _, ok := node.(map[string]any); !ok {
			mlog.Errorf("node is not a map")
			continue
		}
		err := validate.StringMapTypeVerify(node.(map[string]any), validate.Rules{
			"data": {validate.RuleTypeOfField(reflect.Map), validate.NotEmpty()},
		})
		if err != nil {
			mlog.Errorf("node validate failed:%v", err)
			continue
		}
		err = validate.StringMapTypeVerify(node.(map[string]any)["data"].(map[string]any), validate.Rules{
			"type": {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
		})
		if err != nil {
			mlog.Errorf("node data validate failed:%v", err)
			continue
		}
		if node.(map[string]any)["data"].(map[string]any)["type"].(string) == "start" {
			start_node = node.(map[string]any)
			break
		}
	}

	if len(start_node) == 0 {
		mlog.Error("no start node")
		return nil
	}

	if _, ok := start_node["data"]; !ok || start_node["data"] == nil {
		mlog.Errorf("start_node don't have data field")
		return nil
	}
	if _, ok := start_node["data"].(map[string]any); !ok {
		mlog.Errorf("start_node data field is not map")
		return nil
	}
	data_map := start_node["data"].(map[string]any)
	if _, ok := data_map["variables"]; !ok {
		mlog.Errorf("start_node don't have variables field")
		return nil
	}
	if _, ok := data_map["variables"].([]any); !ok {
		mlog.Errorf("start_node variables field is not slice")
		return nil
	}
	data_variables := data_map["variables"].([]any)

	if to_old_structure {
		old_structure_variables := []any{}
		for _, val := range data_variables {
			if _, ok := val.(map[string]any); !ok {
				mlog.Errorf("start_node variables field is not map")
				return nil
			}
			if _, ok := val.(map[string]any)["type"]; ok {
				if _, ok := val.(map[string]any)["type"].(string); ok {
					old_structure_variables = append(old_structure_variables, map[string]any{val.(map[string]any)["type"].(string): val})
				}
			}
		}
		return old_structure_variables
	}
	return data_variables
}
func (wf *Workflow) UniqueHash() string {
	// """
	// Get hash of workflow.

	// :return: hash
	// """
	entity := map[string]any{"graph": wf.GraphDict(), "features": wf.FeaturesDict()}
	bindata, _ := json.Marshal(entity)

	return helper.GenerateTextHash(string(bindata))
}

func (wf *Workflow) GetEnvironmentVariables() []variables.Variabler {
	// # TODO: find some way to init `wf._environment_variables` when instance created.
	if wf.EnvironmentVariablesStr == "" {
		wf.EnvironmentVariablesStr = "{}"
	}
	// tenant_id = contexts.tenant_id.get()
	var environment_variables_dict map[string]map[string]any
	json.Unmarshal([]byte(wf.EnvironmentVariablesStr), &environment_variables_dict)
	results := []variables.Variabler{}
	for _, v := range environment_variables_dict {
		v1 := variablefactory.BuildEnvironmentVariableFromMapping(v)
		results = append(results, v1)
	}
	return results
}

func (wf *Workflow) SetEnvironmentVariables(value []variables.Variabler) {
	if len(value) == 0 {
		wf.EnvironmentVariablesStr = "{}"
		return
	}
	variables_dictionary := map[string]map[string]any{}
	for _, v := range value {
		if v.GetID() == "" {
			panic(exceptions.NewValueError("environment variable require a unique id"))
		}
		variables_dictionary[v.GetName()] = v.ToDict(v)
	}
	bindata, err := json.Marshal(variables_dictionary)
	if err != nil {
		mlog.Errorf("json marshal variables_dictionary=%#v failed:%v", variables_dictionary, err)
		panic(exceptions.NewValueError(fmt.Sprintf("json marshal environment_variables_dictionary=%#v failed:%v", variables_dictionary, err)))
	}

	wf.EnvironmentVariablesStr = string(bindata)
}
func (wf *Workflow) ToDict(include_secret bool) map[string]any {
	environment_variables := wf.GetEnvironmentVariables()
	// environment_variables = [
	//     v if not isinstance(v, SecretVariable) or include_secret else v.model_copy(update={"value": ""})
	//     for v in environment_variables
	// ]

	return map[string]any{
		"graph":                  wf.GraphDict(),
		"features":               wf.FeaturesDict(),
		"environment_variables":  environment_variables,
		"conversation_variables": wf.GetConversationVariables(),
	}
}
func (wf *Workflow) GetConversationVariables() []variables.Variabler {
	// TODO: find some way to init `wf._conversation_variables` when instance created.
	if wf.ConversationVariablesStr == "" {
		wf.ConversationVariablesStr = "{}"
	}
	var variables_dict map[string]map[string]any
	json.Unmarshal([]byte(wf.ConversationVariablesStr), &variables_dict)
	results := []variables.Variabler{}
	// variables_dict:
	// 	dict[str, Any] = json.loads(wf._conversation_variables)
	// results := []*variables.Variable{}
	for _, v := range variables_dict {
		mlog.Debugf("------conversation variable=%#v", v)
		v1 := variablefactory.BuildConversationVariableFromMapping(v)

		results = append(results, v1)
	}
	return results
}

func (wf *Workflow) SetConversationVariables(value []variables.Variabler) {
	if len(value) == 0 {
		wf.ConversationVariablesStr = "{}"
		return
	}
	variables_dictionary := map[string]map[string]any{}
	for _, v := range value {
		if v.GetID() == "" {
			panic(exceptions.NewValueError("conversation variable require a unique id"))
		}
		variables_dictionary[v.GetName()] = v.ToDict(v)
	}
	bindata, err := json.Marshal(variables_dictionary)
	if err != nil {
		mlog.Errorf("json marshal variables_dictionary=%#v failed:%v", variables_dictionary, err)
		panic(exceptions.NewValueError(fmt.Sprintf("json marshal conversation_variables_dictionary=%#v failed:%v", variables_dictionary, err)))
	}

	wf.ConversationVariablesStr = string(bindata)
}

func (wf *Workflow) ToolPublished() bool {
	var count int64
	err := dbengine.Instance().DB.Model(&WorkflowToolProvider{}).Where("tenant_id = ? and app_id = ?", wf.TenantID, wf.AppID).Count(&count).Error
	if err != nil {
		mlog.Errorf("get WorkflowToolProvider by tenant=%#v, app_id= %#v failed:%v", wf.TenantID, wf.AppID, err)
		return false
	}
	return count > 0
}

type WorkflowRunStatus string

// """
// Workflow Run Status Enum
// """
const (
	WorkflowRunStatus_RUNNING           = "running"
	WorkflowRunStatus_SUCCEEDED         = "succeeded"
	WorkflowRunStatus_FAILED            = "failed"
	WorkflowRunStatus_STOPPED           = "stopped"
	WorkflowRunStatus_PARTIAL_SUCCESSED = "partial-succeeded"

// @classmethod
// def value_of(cls, value: str) -> "WorkflowRunStatus":
//     """
//     Get value of given mode.

// :param value: mode value
// :return: mode
// """
// for mode in cls:
//
//	if mode.value == value:
//	    return mode
)

func (s WorkflowRunStatus) Validate() {
	if s != WorkflowRunStatus_RUNNING ||
		s != WorkflowRunStatus_RUNNING ||
		s != WorkflowRunStatus_RUNNING ||
		s != WorkflowRunStatus_RUNNING ||
		s != WorkflowRunStatus_RUNNING {
		panic(exceptions.NewValueError("invalid workflow run status value " + string(s)))
	}
}

// WorkflowRun [...]
type WorkflowRun struct {
	ID              string  `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID        string  `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	AppID           string  `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	SequenceNumber  int     `gorm:"column:sequence_number;type:int;not null" json:"sequence_number"`
	WorkflowID      string  `gorm:"column:workflow_id;type:varchar(36);not null" json:"workflow_id"`
	Type            string  `gorm:"column:type;type:varchar(255);not null" json:"type"`
	TriggeredFrom   string  `gorm:"column:triggered_from;type:varchar(255);not null" json:"triggered_from"`
	Version         string  `gorm:"column:version;type:varchar(255);not null" json:"version"`
	Graph           string  `gorm:"column:graph;type:text" json:"graph_str"`
	Inputs          string  `gorm:"column:inputs;type:text" json:"inputs_str"`
	Status          string  `gorm:"column:status;type:varchar(255);not null" json:"status"`
	Outputs         string  `gorm:"column:outputs;type:text" json:"outputs_str"`
	Error           string  `gorm:"column:error;type:text" json:"error"`
	ElapsedTime     float64 `gorm:"column:elapsed_time;type:double;not null;default:0" json:"elapsed_time"`
	TotalTokens     int     `gorm:"column:total_tokens;type:int;not null;default:0" json:"total_tokens"`
	TotalSteps      int     `gorm:"column:total_steps;type:int;default:0" json:"total_steps"`
	CreatedByRole   `gorm:"column:created_by_role;type:varchar(255);not null" json:"created_by_role"`
	CreatedBy       string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	FinishedAt      *time.Time `gorm:"column:finished_at;type:timestamp" json:"finished_at"`
	ExceptionsCount int        `gorm:"column:exceptions_count;" json:"exceptions_count"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowRun) TableName() string {
	return "workflow_runs"
}

func (wfr *WorkflowRun) GraphDict() map[string]any {
	if wfr.Graph == "" {
		return nil
	}
	var ret map[string]any
	json.Unmarshal([]byte(wfr.Graph), &ret)
	return ret
}

func (wfr *WorkflowRun) InputsDict() map[string]any {
	if wfr.Inputs == "" {
		return nil
	}
	var ret map[string]any
	json.Unmarshal([]byte(wfr.Inputs), &ret)
	return ret
}
func (wfr *WorkflowRun) OutputsDict() map[string]any {
	if wfr.Outputs == "" {
		return nil
	}
	var ret map[string]any
	json.Unmarshal([]byte(wfr.Outputs), &ret)
	return ret
}

func (wfr *WorkflowRun) ToDict() map[string]any {
	return map[string]any{
		"id":               wfr.ID,
		"tenant_id":        wfr.TenantID,
		"app_id":           wfr.AppID,
		"sequence_number":  wfr.SequenceNumber,
		"workflow_id":      wfr.WorkflowID,
		"type":             wfr.Type,
		"triggered_from":   wfr.TriggeredFrom,
		"version":          wfr.Version,
		"graph":            wfr.GraphDict(),
		"inputs":           wfr.InputsDict(),
		"status":           wfr.Status,
		"outputs":          wfr.Outputs,
		"error":            wfr.Error,
		"elapsed_time":     wfr.ElapsedTime,
		"total_tokens":     wfr.TotalTokens,
		"total_steps":      wfr.TotalSteps,
		"created_by_role":  wfr.CreatedByRole,
		"created_by":       wfr.CreatedBy,
		"created_at":       wfr.CreatedAt,
		"finished_at":      wfr.FinishedAt,
		"exceptions_count": wfr.ExceptionsCount,
	}
}
func (wfr *WorkflowRun) FromDict(data map[string]any) *WorkflowRun {
	if data == nil {
		mlog.Warningf("missing data")
		return nil
	}
	res := new(WorkflowRun)
	bindata, _ := json.Marshal(data)
	err := json.Unmarshal(bindata, &res)
	if err != nil {
		mlog.Warningf("json unmarshal WorkflowRun faild:%v", err)
		return nil
	}
	if _, ok := data["graph"]; !ok {
		mlog.Warningf("missing graph")
		return nil
	} else {
		bindata, _ := json.Marshal(data["graph"])
		res.Graph = string(bindata)
	}
	if _, ok := data["inputs"]; !ok {
		mlog.Warningf("missing inputs")
		return nil
	} else {
		bindata, _ := json.Marshal(data["inputs"])
		res.Inputs = string(bindata)
	}
	if _, ok := data["outputs"]; !ok {
		mlog.Warningf("missing outputs")
		return nil
	} else {
		bindata, _ := json.Marshal(data["outputs"])
		res.Outputs = string(bindata)
	}
	return res
}

func (wfr *WorkflowRun) Message() *Message {
	msg := new(Message)
	err := dbengine.Instance().DB.Model(&Message{}).Where("app_id = ? and workflow_run_id = ?", wfr.AppID, wfr.ID).First(msg).Error
	if err != nil {
		mlog.Errorf("get message failed:%v", err)
		return nil
	}
	return msg
}

func (wfr *WorkflowRun) CreatedByAccount() *Account {
	if wfr.CreatedByRole == CreatedByRole_ACCOUNT {
		acc := new(Account)
		err := dbengine.Instance().DB.Model(&Account{}).Where("id = ?", wfr.CreatedBy).First(acc).Error
		if err != nil {
			mlog.Errorf("get Account failed:%v", err)
			return nil
		}
		return acc
	} else {
		return nil
	}
}
func (wfr *WorkflowRun) CreatedByEndUser() *EndUser {
	if wfr.CreatedByRole == CreatedByRole_END_USER {
		user := new(EndUser)
		err := dbengine.Instance().DB.Model(&EndUser{}).Where("id = ?", wfr.CreatedBy).First(user).Error
		if err != nil {
			mlog.Errorf("get EndUser failed:%v", err)
			return nil
		}
		return user
	} else {
		return nil
	}
}
func (wfr *WorkflowRun) App() *App {
	if wfr.AppID != "" {
		ret := new(App)
		err := dbengine.Instance().DB.Model(&App{}).Where("id = ?", wfr.AppID).Scan(ret).Error
		if err != nil {
			mlog.Errorf("get App from mysql failed:%v", err)
			return nil
		}
		return ret
	}
	return nil
}
func (wfr *WorkflowRun) Workflow() *Workflow {
	if wfr.WorkflowID != "" {
		ret := new(Workflow)
		err := dbengine.Instance().DB.Model(&Workflow{}).Where("id = ?", wfr.WorkflowID).Scan(ret).Error
		if err != nil {
			mlog.Errorf("get Workflow from mysql failed:%v", err)
			return nil
		}
		return ret
	}
	return nil
}

func (wfr *WorkflowRun) Tenant() *Tenant {
	if wfr.TenantID != "" {
		ret := new(Tenant)
		err := dbengine.Instance().DB.Model(&Tenant{}).Where("id = ?", wfr.TenantID).Scan(ret).Error
		if err != nil {
			mlog.Errorf("get Tenant from mysql failed:%v", err)
			return nil
		}
		return ret
	}
	return nil
}

type WorkflowWithMessage struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	*WorkflowRun
}

func NewWorkflowWithMessage(workflow_run *WorkflowRun) *WorkflowWithMessage {
	return &WorkflowWithMessage{
		WorkflowRun: workflow_run,
	}
}

// WorkflowNodeExecution [...]
type WorkflowNodeExecution struct {
	ID                string                  `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID          string                  `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Tenant            *Tenant                 `json:"tenant" form:"tenant" gorm:"foreignKey:TenantID;references:ID;"`
	AppID             string                  `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App               *App                    `json:"app" form:"app" gorm:"foreignKey:AppID;references:ID;"`
	WorkflowID        string                  `gorm:"column:workflow_id;type:varchar(36);not null" json:"workflow_id"`
	Workflow          *Workflow               `json:"workflow" form:"workflow" gorm:"foreignKey:WorkflowID;references:ID;"`
	TriggeredFrom     string                  `gorm:"column:triggered_from;type:varchar(255);not null" json:"triggered_from"`
	WorkflowRunID     string                  `gorm:"column:workflow_run_id;type:varchar(36)" json:"workflow_run_id"`
	WorkflowRun       *WorkflowRun            `json:"workflow_run" form:"workflow_run" gorm:"foreignKey:WorkflowRunID;references:ID;"`
	Index             int                     `gorm:"column:index;type:int;not null" json:"index"`
	PredecessorNodeID string                  `gorm:"column:predecessor_node_id;type:varchar(255)" json:"predecessor_node_id"`
	NodeExecutionID   string                  `gorm:"column:node_execution_id;type:varchar(255)" json:"node_execution_id"`
	NodeID            string                  `gorm:"column:node_id;type:varchar(255);not null" json:"node_id"`
	NodeType          nodesenumtypes.NodeType `gorm:"column:node_type;type:varchar(255);not null" json:"node_type"`
	Title             string                  `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Inputs            string                  `gorm:"column:inputs;type:text" json:"inputs"`
	ProcessData       string                  `gorm:"column:process_data;type:text" json:"process_data"`
	Outputs           string                  `gorm:"column:outputs;type:text" json:"outputs"`
	Status            string                  `gorm:"column:status;type:varchar(255);not null" json:"status"`
	Error             string                  `gorm:"column:error;type:text" json:"error"`
	ElapsedTime       float64                 `gorm:"column:elapsed_time;type:double;not null;default:0" json:"elapsed_time"`
	ExecutionMetadata string                  `gorm:"column:execution_metadata;type:text" json:"execution_metadata"`
	CreatedAt         *time.Time              `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedByRole     `gorm:"column:created_by_role;type:varchar(255);not null" json:"created_by_role"`
	CreatedBy         string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	FinishedAt        *time.Time `gorm:"column:finished_at;type:timestamp" json:"finished_at"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowNodeExecution) TableName() string {
	return "workflow_node_executions"
}

func (wfne *WorkflowNodeExecution) InputsDict() map[string]any {
	if wfne.Inputs == "" {
		return nil
	}
	var ret map[string]any
	json.Unmarshal([]byte(wfne.Inputs), &ret)
	return ret
}
func (wfne *WorkflowNodeExecution) OutputsDict() map[string]any {
	if wfne.Outputs == "" {
		return nil
	}
	var ret map[string]any
	json.Unmarshal([]byte(wfne.Outputs), &ret)
	return ret
}
func (wfne *WorkflowNodeExecution) ProcessDataDict() map[string]any {
	if wfne.ProcessData == "" {
		return nil
	}
	var ret map[string]any
	json.Unmarshal([]byte(wfne.ProcessData), &ret)
	return ret
}
func (wfne *WorkflowNodeExecution) ExecutionMetadataDict() map[string]any {
	if wfne.ExecutionMetadata == "" {
		return nil
	}
	var ret map[string]any
	json.Unmarshal([]byte(wfne.ExecutionMetadata), &ret)
	return ret
}

func (wfne *WorkflowNodeExecution) Extras() map[string]any {
	// from core.tools.tool_manager import ToolManager

	// extras = {}
	// if self.execution_metadata_dict:
	//     from core.workflow.nodes import NodeType

	//     if self.node_type == NodeType.TOOL.value and "tool_info" in self.execution_metadata_dict:
	//         tool_info = self.execution_metadata_dict["tool_info"]
	//         extras["icon"] = ToolManager.get_tool_icon(
	//             tenant_id=self.tenant_id,
	//             provider_type=tool_info["provider_type"],
	//             provider_id=tool_info["provider_id"],
	//         )

	// return extras
	return map[string]any{}
}

func (wfne *WorkflowNodeExecution) CreatedByAccount() *Account {
	if wfne.CreatedByRole == CreatedByRole_ACCOUNT {
		acc := new(Account)
		err := dbengine.Instance().DB.Model(&Account{}).Where("id = ?", wfne.CreatedBy).First(acc).Error
		if err != nil {
			mlog.Errorf("get Account failed:%v", err)
			return nil
		}
		return acc
	} else {
		return nil
	}
}
func (wfne *WorkflowNodeExecution) CreatedByEndUser() *EndUser {
	if wfne.CreatedByRole == CreatedByRole_END_USER {
		user := new(EndUser)
		err := dbengine.Instance().DB.Model(&EndUser{}).Where("id = ?", wfne.CreatedBy).First(user).Error
		if err != nil {
			mlog.Errorf("get EndUser failed:%v", err)
			return nil
		}
		return user
	} else {
		return nil
	}
}

// WorkflowAppLog [...]
type WorkflowAppLog struct {
	ID               string       `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID         string       `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	Tenant           *Tenant      `json:"tenant" form:"tenant" gorm:"foreignKey:TenantID;references:ID;"`
	AppID            string       `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	App              *App         `json:"app" form:"app" gorm:"foreignKey:AppID;references:ID;"`
	WorkflowID       string       `gorm:"column:workflow_id;type:varchar(36);not null" json:"workflow_id"`
	Workflow         *Workflow    `json:"workflow" form:"workflow" gorm:"foreignKey:WorkflowID;references:ID;"`
	WorkflowRunID    string       `gorm:"column:workflow_run_id;type:varchar(36);not null" json:"workflow_run_id"`
	WorkflowRun      *WorkflowRun `json:"workflow_run" form:"workflow_run" gorm:"foreignKey:WorkflowRunID;references:ID;"`
	CreatedFrom      string       `gorm:"column:created_from;type:varchar(255);not null" json:"created_from"`
	CreatedByRole    `gorm:"column:created_by_role;type:varchar(255);not null" json:"created_by_role"`
	CreatedBy        string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	CreatedByAccount *Account   `json:"created_by_account" form:"created_by_account" gorm:"foreignKey:CreatedBy;references:ID;"`
	CreatedAt        *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowAppLog) TableName() string {
	return "workflow_app_logs"
}

// ConversationVariable [...]
type ConversationVariable struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	ConversationID string     `gorm:"primaryKey;column:conversation_id;type:varchar(36);not null" json:"conversation_id"`
	AppID          string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	Data           string     `gorm:"column:data;type:text;not null" json:"data"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (ConversationVariable) TableName() string {
	return "workflow_conversation_variables"
}

func (cv *ConversationVariable) FromVariable(app_id, conversation_id string, val variables.Variabler) *ConversationVariable {
	cls := &ConversationVariable{
		ID:             val.GetID(),
		AppID:          app_id,
		ConversationID: conversation_id,
	}
	bindata, _ := json.Marshal(val.ToDict(val))
	cls.Data = string(bindata)
	return cls
}
func (cv *ConversationVariable) ToVariable() variables.Variabler {
	var mapping map[string]any
	err := json.Unmarshal([]byte(cv.Data), &mapping)
	if err != nil {
		mlog.Errorf("json unmarshal(%s) failed:%v", cv.Data, err)
		return nil
	}
	res := variablefactory.BuildConversationVariableFromMapping(mapping)

	return res
}

// WorkflowConversationVariable [...]
type WorkflowConversationVariable struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	ConversationID string     `gorm:"column:conversation_id;type:varchar(36);not null" json:"conversation_id"`
	AppID          string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	Data           string     `gorm:"column:data;type:text;not null" json:"data"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowConversationVariable) TableName() string {
	return "workflow_conversation_variable"
}

// WorkflowNodeExecutionOffload [...]
type WorkflowNodeExecutionOffload struct {
	ID              string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	TenantID        string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	AppID           string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	NodeExecutionID *string    `gorm:"column:node_execution_id;type:varchar(36)" json:"node_execution_id"`
	Type            string     `gorm:"column:type;type:varchar(50);not null" json:"type"`
	FileID          string     `gorm:"column:file_id;type:varchar(36);not null" json:"file_id"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowNodeExecutionOffload) TableName() string {
	return "workflow_node_execution_offload"
}

// WorkflowArchiveLog [...]
type WorkflowArchiveLog struct {
	ID                 string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	TenantID           string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	AppID              string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	WorkflowID         string     `gorm:"column:workflow_id;type:varchar(36);not null" json:"workflow_id"`
	WorkflowRunID      string     `gorm:"column:workflow_run_id;type:varchar(36);not null" json:"workflow_run_id"`
	CreatedByRole      string     `gorm:"column:created_by_role;type:varchar(255);not null" json:"created_by_role"`
	CreatedBy          string     `gorm:"column:created_by;type:varchar(36);not null" json:"created_by"`
	LogID              *string    `gorm:"column:log_id;type:varchar(36)" json:"log_id"`
	LogCreatedAt       *time.Time `gorm:"column:log_created_at;type:timestamp" json:"log_created_at"`
	LogCreatedFrom     *string    `gorm:"column:log_created_from;type:varchar(255)" json:"log_created_from"`
	RunVersion         string     `gorm:"column:run_version;type:varchar(255);not null" json:"run_version"`
	RunStatus          string     `gorm:"column:run_status;type:varchar(255);not null" json:"run_status"`
	RunTriggeredFrom   string     `gorm:"column:run_triggered_from;type:varchar(255);not null" json:"run_triggered_from"`
	RunError           *string    `gorm:"column:run_error;type:text" json:"run_error"`
	RunElapsedTime     float64    `gorm:"column:run_elapsed_time;type:float;not null;default:0" json:"run_elapsed_time"`
	RunTotalTokens     int64      `gorm:"column:run_total_tokens;type:bigint;default:0" json:"run_total_tokens"`
	RunTotalSteps      *int       `gorm:"column:run_total_steps;type:int;default:0" json:"run_total_steps"`
	RunCreatedAt       time.Time  `gorm:"column:run_created_at;type:timestamp;not null" json:"run_created_at"`
	RunFinishedAt      *time.Time `gorm:"column:run_finished_at;type:timestamp" json:"run_finished_at"`
	RunExceptionsCount *int       `gorm:"column:run_exceptions_count;type:int;default:0" json:"run_exceptions_count"`
	TriggerMetadata    *string    `gorm:"column:trigger_metadata;type:text" json:"trigger_metadata"`
	ArchivedAt         *time.Time `gorm:"column:archived_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"archived_at"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowArchiveLog) TableName() string {
	return "workflow_archive_logs"
}

// WorkflowDraftVariable [...]
type WorkflowDraftVariable struct {
	ID              string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	AppID           string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	UserID          *string    `gorm:"column:user_id;type:varchar(36)" json:"user_id"`
	LastEditedAt    *time.Time `gorm:"column:last_edited_at;type:timestamp" json:"last_edited_at"`
	NodeID          string     `gorm:"column:node_id;type:varchar(255);not null" json:"node_id"`
	Name            string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description     string     `gorm:"column:description;type:varchar(255);not null;default:''" json:"description"`
	Selector        string     `gorm:"column:selector;type:varchar(255);not null" json:"selector"`
	ValueType       string     `gorm:"column:value_type;type:varchar(20);not null" json:"value_type"`
	Value           string     `gorm:"column:value;type:text;not null" json:"value"`
	Visible         bool       `gorm:"column:visible;type:bool;not null;default:true" json:"visible"`
	Editable        bool       `gorm:"column:editable;type:bool;not null;default:false" json:"editable"`
	NodeExecutionID *string    `gorm:"column:node_execution_id;type:varchar(36)" json:"node_execution_id"`
	FileID          *string    `gorm:"column:file_id;type:varchar(36)" json:"file_id"`
	IsDefaultValue  bool       `gorm:"column:is_default_value;type:bool;not null;default:false" json:"is_default_value"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowDraftVariable) TableName() string {
	return "workflow_draft_variables"
}

// WorkflowDraftVariableFile [...]
type WorkflowDraftVariableFile struct {
	ID           string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	CreatedAt    *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	TenantID     string     `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	AppID        string     `gorm:"column:app_id;type:varchar(36);not null" json:"app_id"`
	UserID       string     `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	UploadFileID string     `gorm:"column:upload_file_id;type:varchar(36);not null" json:"upload_file_id"`
	Size         *int64     `gorm:"column:size;type:bigint;not null" json:"size"`
	Length       *int       `gorm:"column:length;type:int" json:"length"`
	ValueType    string     `gorm:"column:value_type;type:varchar(20);not null" json:"value_type"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowDraftVariableFile) TableName() string {
	return "workflow_draft_variable_files"
}

// WorkflowPause [...]
type WorkflowPause struct {
	ID             string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	CreatedAt      *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	WorkflowID     string     `gorm:"column:workflow_id;type:varchar(36);not null" json:"workflow_id"`
	WorkflowRunID  string     `gorm:"column:workflow_run_id;type:varchar(36);not null;uniqueIndex" json:"workflow_run_id"`
	ResumedAt      *time.Time `gorm:"column:resumed_at;type:timestamp" json:"resumed_at"`
	StateObjectKey string     `gorm:"column:state_object_key;type:varchar(255);not null" json:"state_object_key"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowPause) TableName() string {
	return "workflow_pauses"
}

// WorkflowPauseReason [...]
type WorkflowPauseReason struct {
	ID        string     `gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	PauseID   string     `gorm:"column:pause_id;type:varchar(36);not null;index" json:"pause_id"`
	Type      string     `gorm:"column:type;type:varchar(50);not null" json:"type"`
	FormID    string     `gorm:"column:form_id;type:varchar(36);not null;default:''" json:"form_id"`
	Message   string     `gorm:"column:message;type:varchar(255);not null;default:''" json:"message"`
	NodeID    string     `gorm:"column:node_id;type:varchar(255);not null;default:''" json:"node_id"`
}

// TableName get sql table name.获取数据库表名
func (WorkflowPauseReason) TableName() string {
	return "workflow_pause_reasons"
}
