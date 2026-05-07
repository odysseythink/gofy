package services

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/odysseythink/gofy/backend/cache"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/gofy/backend/core/variables"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/events"
	variablefactory "github.com/odysseythink/gofy/backend/factories/variable_factory"
	"github.com/odysseythink/gofy/backend/models"
	pbentities "github.com/odysseythink/gofy/backend/proto/entities"
	"github.com/odysseythink/gofy/backend/utils"
	versionutils "github.com/odysseythink/gofy/backend/utils/version"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v3"
)

const (
	IMPORT_INFO_REDIS_KEY_PREFIX = "app_import_info:"
	IMPORT_INFO_REDIS_EXPIRY     = 180 // 3 minutes
	CURRENT_DSL_VERSION          = "0.1.5"
)

type ImportModeType string

const (
	ImportMode_YAML_CONTENT ImportModeType = "yaml-content"
	ImportMode_YAML_URL     ImportModeType = "yaml-url"
)

type ImportStatusType string

const (
	ImportStatus_COMPLETED               ImportStatusType = "completed"
	ImportStatus_COMPLETED_WITH_WARNINGS ImportStatusType = "completed-with-warnings"
	ImportStatus_PENDING                 ImportStatusType = "pending"
	ImportStatus_FAILED                  ImportStatusType = "failed"
)

func NewPendingData(content string) *pbentities.PendingData {
	data := new(pbentities.PendingData)
	err := json.Unmarshal([]byte(content), data)
	if err != nil {
		mlog.Errorf("json unmarshal(%s) failed:%v", content, err)
		panic(exceptions.NewValueError(""))
	}
	return data
}

func NewImport() *pbentities.Import {
	return &pbentities.Import{
		CurrentDslVersion: CURRENT_DSL_VERSION,
	}
}
func _check_version_compatibility(imported_version string) ImportStatusType {

	current_ver, err := versionutils.Parse(CURRENT_DSL_VERSION)
	if err != nil {
		mlog.Errorf("parse version failed:%v", err)
		return ImportStatus_FAILED
	}
	imported_ver, err := versionutils.Parse(imported_version)
	if err != nil {
		mlog.Errorf("parse version failed:%v", err)
		return ImportStatus_FAILED
	}

	// Compare major version and minor version
	if current_ver.Major() != imported_ver.Major() || current_ver.Minor() != imported_ver.Minor() {
		return ImportStatus_PENDING
	}
	if current_ver.Patch() != imported_ver.Patch() {
		return ImportStatus_COMPLETED_WITH_WARNINGS
	}
	return ImportStatus_COMPLETED
}

type AppDSLService struct {
}

func (service *AppDSLService) ImportApp(
	account *models.Account,
	import_mode ImportModeType,
	yaml_content string,
	yaml_url string,
	name string,
	description string,
	icon_type string,
	icon string,
	icon_background string,
	app_id string,
) *pbentities.Import {
	/*Import an app from YAML content or URL.*/
	import_id := uuid.NewV4().String()

	// Validate import mode
	if import_mode != ImportMode_YAML_CONTENT {
		panic(exceptions.NewValueError("Invalid import_mode: " + string(import_mode)))
	}

	// Get YAML content
	if yaml_content == "" {
		imp := NewImport()
		imp.Id = import_id
		imp.Status = string(ImportStatus_FAILED)
		imp.Error = "yaml_content is required when import_mode is yaml-content"
		return imp
	}
	content := yaml_content

	// Process YAML content
	return func() (imp *pbentities.Import) {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				imp = NewImport()
				imp.Id = import_id
				imp.Status = string(ImportStatus_FAILED)
				imp.Error = fmt.Sprintf("%v", r)
			}
		}()
		// Parse YAML to validate format
		var data map[string]any
		err := yaml.Unmarshal([]byte(content), &data)
		if err != nil {
			mlog.Errorf("yaml unmarshal=%s failed:%v", content, err)
			imp = NewImport()
			imp.Id = import_id
			imp.Status = string(ImportStatus_FAILED)
			imp.Error = "Invalid YAML format: content must be a mapping"
			return imp
		}
		// Validate and fix DSL version
		if _, ok := data["version"]; !ok {
			data["version"] = "0.1.0"
		}

		data["kind"] = "app"
		imported_version := "0.1.0"
		if _, ok := data["version"]; ok {
			if _, ok := data["version"].(string); ok {
				imported_version = data["version"].(string)
			} else {
				panic(exceptions.NewValueError(fmt.Sprintf("Invalid version type, expected str, got %#v", data["version"])))
			}
		}

		status := _check_version_compatibility(imported_version)

		// Extract app data
		if _, ok := data["app"]; !ok {
			mlog.Errorf("Missing app data in YAML content")
			imp = NewImport()
			imp.Id = import_id
			imp.Status = string(ImportStatus_FAILED)
			imp.Error = "Missing app data in YAML content"
			return imp
		}

		// If app_id is provided, check if it exists
		var app *models.App
		if app_id != "" {
			app := new(models.App)
			err := dbengine.Instance().DB.Model(&models.App{}).Where("id =? and tenant_id = ?", app_id, account.CurrentTenantID()).Find(app).Error
			if err != nil {
				mlog.Errorf("count App failed:%v", err)
			}

			if app == nil {
				mlog.Errorf("App not found")
				imp = NewImport()
				imp.Id = import_id
				imp.Status = string(ImportStatus_FAILED)
				imp.Error = "App not found"
				return imp
			}
			if !slices.Contains([]models.AppMode{models.AppMode_WORKFLOW, models.AppMode_ADVANCED_CHAT}, app.Mode) {
				mlog.Errorf("Only workflow or advanced chat apps can be overwritten")
				imp = NewImport()
				imp.Id = import_id
				imp.Status = string(ImportStatus_FAILED)
				imp.Error = "Only workflow or advanced chat apps can be overwritten"
				return imp
			}
		}
		// If major version mismatch, store import info in Redis
		if status == ImportStatus_PENDING {
			panding_data := pbentities.PendingData{
				ImportMode:     string(import_mode),
				YamlContent:    content,
				Name:           name,
				Description:    description,
				IconType:       icon_type,
				Icon:           icon,
				IconBackground: icon_background,
				AppId:          app_id,
			}
			bindata, _ := json.Marshal(panding_data)
			cache.Instance().SetEx(IMPORT_INFO_REDIS_KEY_PREFIX+import_id, string(bindata), IMPORT_INFO_REDIS_EXPIRY*time.Second)
			imp = NewImport()
			imp.Id = import_id
			imp.Status = string(status)
			imp.AppId = app_id
			imp.ImportedDslVersion = imported_version
			return imp
		}

		// Create or update app
		app = service._create_or_update_app(
			app,
			data,
			account,
			name,
			description,
			icon_type,
			icon,
			icon_background,
		)
		imp = NewImport()
		imp.Id = import_id
		imp.Status = string(status)
		imp.AppId = app.ID
		imp.ImportedDslVersion = imported_version
		return imp
	}()
}

func (service *AppDSLService) ConfirmImport(import_id string, account *models.Account) *pbentities.Import {
	/*
	   Confirm an import that requires confirmation
	*/
	redis_key := IMPORT_INFO_REDIS_KEY_PREFIX + import_id
	pending_data_str := cache.Instance().GetString(redis_key)

	if pending_data_str == "" {
		mlog.Errorf("Import information expired or does not exist")
		imp := NewImport()
		imp.Id = import_id
		imp.Status = string(ImportStatus_FAILED)
		imp.Error = "Import information expired or does not exist"
		return imp
	}
	return func() (imp *pbentities.Import) {
		defer func() {
			if r := recover(); r != nil {
				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
				imp = NewImport()
				imp.Id = import_id
				imp.Status = string(ImportStatus_FAILED)
				imp.Error = fmt.Sprintf("%v", r)
			}
		}()

		pending_data := NewPendingData(pending_data_str)
		var data map[string]any
		err := yaml.Unmarshal([]byte(pending_data.YamlContent), &data)
		if err != nil {
			mlog.Errorf("yaml unmarshal(%s) failed:%v", pending_data.YamlContent, err)
			panic(exceptions.NewValueError("pending data yaml content is not a dict"))
		}

		var app *models.App
		if pending_data.AppId != "" {
			app = new(models.App)
			err := dbengine.Instance().DB.Model(&models.App{}).Where("id =? and tenant_id=?", pending_data.AppId, account.CurrentTenantID()).Find(app).Error
			if err != nil {
				mlog.Errorf("count models.App failed:%v", err)
				app = nil
			}
		}
		// Create or update app
		app = service._create_or_update_app(
			app,
			data,
			account,
			pending_data.Name,
			pending_data.Description,
			pending_data.IconType,
			pending_data.Icon,
			pending_data.IconBackground,
		)

		// Delete import info from Redis
		version := "0.1.0"
		if _, ok := data["version"]; ok {
			if _, ok := data["version"].(string); ok {
				version = data["version"].(string)
			}
		}
		cache.Instance().DelKey(redis_key)
		imp = NewImport()
		imp.Id = import_id
		imp.Status = string(ImportStatus_COMPLETED)
		imp.AppId = app.ID
		imp.CurrentDslVersion = CURRENT_DSL_VERSION
		imp.ImportedDslVersion = version
		return
	}()

}
func (service *AppDSLService) _create_or_update_app(
	app *models.App,
	data map[string]any,
	account *models.Account,
	name string,
	description string,
	icon_type string,
	icon string,
	icon_background string,
) *models.App {
	/*Create a new app or update an existing one.*/
	app_data := map[string]any{}
	if _, ok := data["app"]; ok {
		if _, ok := data["app"].(map[string]any); ok {
			app_data = data["app"].(map[string]any)
		}
	}
	app_mode_str := ""
	if _, ok := app_data["mode"]; ok {
		if _, ok := app_data["mode"].(string); ok {
			app_mode_str = app_data["mode"].(string)
		}
	}

	if app_mode_str == "" {
		panic(exceptions.NewValueError("loss app mode"))
	}
	app_mode := models.AppMode(app_mode_str)

	// Set icon type
	icon_type_value := icon_type
	if icon_type_value == "" {
		if _, ok := app_data["icon_type"]; ok {
			if _, ok := app_data["icon_type"].(string); ok {
				icon_type_value = app_data["icon_type"].(string)
			}
		}
	}

	if slices.Contains([]string{"emoji", "link"}, icon_type_value) {
		icon_type = icon_type_value
	} else {
		icon_type = "emoji"
	}
	if icon == "" {
		if _, ok := app_data["icon"]; ok {
			if _, ok := app_data["icon"].(string); ok {
				icon = app_data["icon"].(string)
			}
		}
	}

	if app != nil {
		// Update existing app
		app.Name = name
		if app.Name == "" {
			if _, ok := app_data["name"]; ok {
				if _, ok := app_data["name"].(string); ok {
					app.Name = app_data["name"].(string)
				}
			}
		}

		app.Description = description
		if app.Description == "" {
			if _, ok := app_data["description"]; ok {
				if _, ok := app_data["description"].(string); ok {
					app.Description = app_data["description"].(string)
				}
			}
		}

		app.IconType = icon_type
		app.Icon = icon
		app.IconBackground = icon_background
		if app.IconBackground == "" {
			if _, ok := app_data["icon_background"]; ok {
				if _, ok := app_data["icon_background"].(string); ok {
					app.IconBackground = app_data["icon_background"].(string)
				}
			}
		}
		app.UpdatedBy = account.ID
	} else {
		if account.CurrentTenant == nil {
			panic(exceptions.NewValueError("Current tenant is not set"))
		}
		now := time.Now()
		// Create new app
		app = &models.App{
			ID:                  uuid.NewV4().String(),
			TenantID:            account.CurrentTenantID(),
			Name:                name,
			Mode:                app_mode,
			Icon:                icon,
			IconBackground:      icon_background,
			Status:              "normal",
			EnableSite:          true,
			EnableAPI:           true,
			CreatedAt:           &now,
			UpdatedAt:           &now,
			Description:         description,
			IconType:            icon_type,
			CreatedBy:           account.ID,
			UpdatedBy:           account.ID,
			UseIconAsAnswerIcon: false,
		}
		if app.Name == "" {
			if _, ok := app_data["name"]; ok {
				if _, ok := app_data["name"].(string); ok {
					app.Name = app_data["name"].(string)
				}
			}
		}
		if app.Description == "" {
			if _, ok := app_data["description"]; ok {
				if _, ok := app_data["description"].(string); ok {
					app.Description = app_data["description"].(string)
				}
			}
		}
		if app.IconBackground == "" {
			if _, ok := app_data["icon_background"]; ok {
				if _, ok := app_data["icon_background"].(string); ok {
					app.IconBackground = app_data["icon_background"].(string)
				}
			}
		}

		if app.IconBackground == "" {
			app.IconBackground = "//FFFFFF"
		}
		if _, ok := app_data["use_icon_as_answer_icon"]; ok {
			if _, ok := app_data["use_icon_as_answer_icon"].(bool); ok {
				app.UseIconAsAnswerIcon = app_data["use_icon_as_answer_icon"].(bool)
			}
		}
		dbengine.Instance().DB.Create(app)
		events.Instance.AppWasCreatedSig.Emit(app, account)
	}

	// Initialize app based on mode
	if slices.Contains([]models.AppMode{models.AppMode_ADVANCED_CHAT, models.AppMode_WORKFLOW}, app_mode) {
		workflow_data := map[string]any{}
		if _, ok := data["workflow"]; ok {
			if _, ok := data["workflow"].(map[string]any); ok {
				workflow_data = data["workflow"].(map[string]any)
			}
		}

		if len(workflow_data) == 0 {
			panic(exceptions.NewValueError("Missing workflow data for workflow/advanced chat app"))
		}
		environment_variables_list := []map[string]any{}
		if _, ok := workflow_data["environment_variables"]; ok {
			if _, ok := workflow_data["environment_variables"].([]any); ok {
				for _, v := range workflow_data["environment_variables"].([]any) {
					if _, ok := v.(map[string]any); ok {
						environment_variables_list = append(environment_variables_list, v.(map[string]any))
					} else {
						mlog.Errorf("invalid workflow environment_variables(%#v) data for workflow/advanced chat app", workflow_data["environment_variables"])
						panic(exceptions.NewValueError("invalid workflow environment_variables data for workflow/advanced chat app"))
					}
				}
			} else if _, ok := workflow_data["environment_variables"].([]map[string]any); ok {
				environment_variables_list = workflow_data["environment_variables"].([]map[string]any)
			}
		}
		environment_variables := []variables.Variabler{}
		for _, obj := range environment_variables_list {
			environment_variables = append(environment_variables, variablefactory.BuildEnvironmentVariableFromMapping(obj))
		}
		conversation_variables_list := []map[string]any{}
		if _, ok := workflow_data["conversation_variables"]; ok {
			if _, ok := workflow_data["conversation_variables"].([]any); ok {
				for _, v := range workflow_data["conversation_variables"].([]any) {
					if _, ok := v.(map[string]any); ok {
						conversation_variables_list = append(conversation_variables_list, v.(map[string]any))
					} else {
						mlog.Errorf("invalid workflow conversation_variables(%#v) data for workflow/advanced chat app", workflow_data["conversation_variables"])
						panic(exceptions.NewValueError("invalid workflow conversation_variables data for workflow/advanced chat app"))
					}
				}
			} else if _, ok := workflow_data["conversation_variables"].([]map[string]any); ok {
				conversation_variables_list = workflow_data["conversation_variables"].([]map[string]any)
			}
		}
		conversation_variables := []variables.Variabler{}
		for _, obj := range conversation_variables_list {
			conversation_variables = append(conversation_variables, variablefactory.BuildConversationVariableFromMapping(obj))
		}

		current_draft_workflow := ServiceGroupApp.Workflow.GetDraftWorkflow(app)
		unique_hash := ""
		if current_draft_workflow != nil {
			unique_hash = current_draft_workflow.UniqueHash()
		}
		graph_dict := map[string]any{}
		if _, ok := workflow_data["graph"]; ok {
			if _, ok := workflow_data["graph"].(map[string]any); ok {
				graph_dict = workflow_data["graph"].(map[string]any)
			} else {
				mlog.Errorf("invalid workflow graph(%#v) data for workflow/advanced chat app", workflow_data["graph"])
				panic(exceptions.NewValueError("invalid workflow graph data for workflow/advanced chat app"))
			}
		}
		features_dict := map[string]any{}
		if _, ok := workflow_data["features"]; ok {
			if _, ok := workflow_data["features"].(map[string]any); ok {
				features_dict = workflow_data["features"].(map[string]any)
			} else {
				mlog.Errorf("invalid workflow features(%#v) data for workflow/advanced chat app", workflow_data["features"])
				panic(exceptions.NewValueError("invalid workflow features data for workflow/advanced chat app"))
			}
		}
		ServiceGroupApp.Workflow.SyncDraftWorkflow(
			app,
			graph_dict,
			features_dict,
			unique_hash,
			account,
			environment_variables,
			conversation_variables,
		)
	} else if slices.Contains([]models.AppMode{models.AppMode_CHAT, models.AppMode_AGENT_CHAT, models.AppMode_COMPLETION}, app_mode) {
		// Initialize model config
		model_config := map[string]any{}
		if _, ok := data["model_config"]; ok {
			if _, ok := data["model_config"].(map[string]any); ok {
				model_config = data["model_config"].(map[string]any)
			}
		}
		if len(model_config) == 0 {
			panic(exceptions.NewValueError("Missing model_config for chat/agent-chat/completion app"))
		}
		// Initialize or update model config
		app_model_config := app.AppModelConfig()
		if app_model_config == nil {
			app_model_config = new(models.AppModelConfig)
			app_model_config.FromModelConfigDict(model_config)
			app_model_config.ID = uuid.NewV4().String()
			app_model_config.AppID = app.ID
			app_model_config.CreatedBy = account.ID
			app_model_config.UpdatedBy = account.ID
			dbengine.Instance().DB.Create(app_model_config)

			app.AppModelConfigID = app_model_config.ID

			events.Instance.AppModelConfigWasUpdatedSig.Emit(app, app_model_config)
		}
	} else {
		panic(exceptions.NewValueError("Invalid app mode"))
	}
	return app

}

func (service *AppDSLService) ExportDSL(app_model *models.App, include_secret bool) string {
	/*
	   Export app
	   :param app_model: *models.App instance
	   :return:
	*/
	app_mode := app_model.Mode

	export_data := map[string]any{
		"version": CURRENT_DSL_VERSION,
		"kind":    "app",
		"app": map[string]any{
			"name":                    app_model.Name,
			"mode":                    app_model.Mode,
			"description":             app_model.Description,
			"use_icon_as_answer_icon": app_model.UseIconAsAnswerIcon,
		},
	}
	if app_model.IconType == "image" {
		export_data["app"].(map[string]any)["icon"] = "🤖"
		export_data["app"].(map[string]any)["icon_background"] = "#FFEAD5"
	} else {
		export_data["app"].(map[string]any)["icon"] = app_model.Icon
		export_data["app"].(map[string]any)["icon"] = app_model.IconBackground
	}

	if slices.Contains([]models.AppMode{models.AppMode_ADVANCED_CHAT, models.AppMode_WORKFLOW}, app_mode) {
		service._append_workflow_export_data(export_data, app_model, include_secret)
	} else {
		service._append_model_config_export_data(export_data, app_model)
	}
	bindata, _ := yaml.Marshal(export_data)
	return string(bindata) // type: ignore

}
func (service *AppDSLService) _append_workflow_export_data(export_data map[string]any, app_model *models.App, include_secret bool) {
	/*
	   Append workflow export data
	   :param export_data: export data
	   :param app_model: *models.App instance
	*/
	wf := ServiceGroupApp.Workflow.GetDraftWorkflow(app_model)
	if wf == nil {
		panic(exceptions.NewValueError("Missing draft workflow configuration, please check."))
	}
	export_data["workflow"] = wf.ToDict(include_secret)
}
func (service *AppDSLService) _append_model_config_export_data(export_data map[string]any, app_model *models.App) {
	/*
	   Append model config export data
	   :param export_data: export data
	   :param app_model *models.App instance
	*/
	app_model_config := app_model.AppModelConfig()
	if app_model_config == nil {
		panic(exceptions.NewValueError("Missing app configuration, please check."))
	}
	export_data["model_config"] = app_model_config.ToDict()
}
