package services

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"mlib.com/confy"
	"mlib.com/gofy/server/core/exceptions"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	"mlib.com/gofy/server/core/manageres"
	modelmanager "mlib.com/gofy/server/core/manageres/model_manager"
	dbengine "mlib.com/gofy/server/db_engine"
	agententities "mlib.com/gofy/server/entities/agent"
	modelruntimeenumtypes "mlib.com/gofy/server/enum_types/model_runtime"
	"mlib.com/gofy/server/events"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/request"
	"mlib.com/gofy/server/models/response"
	"mlib.com/gofy/server/proto/pbapi"
	"mlib.com/gofy/server/utils"
	"mlib.com/mlog"
)

type AppService struct {
}

// Create 创建App记录
func (s *AppService) Create(app *models.App) (err error) {
	err = dbengine.Instance().DB.Create(app).Error
	return err
}

// Delete 删除App记录
func (s *AppService) Delete(app models.App) (err error) {
	err = dbengine.Instance().DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.App{}).Where("id = ?", app.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&app).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByIds 批量删除App记录
func (s *AppService) DeleteByIds(ids []string) (err error) {
	if err := dbengine.Instance().DB.Model(&models.App{}).Delete("id in ?", ids).Error; err != nil {
		return err
	}

	return nil
}

// Update 更新App记录
func (s *AppService) Update(app *models.App) (err error) {
	err = dbengine.Instance().DB.Save(app).Error
	return err
}

// Get 根据id获取App记录
func (s *AppService) Get(id string) (app models.App, err error) {
	err = dbengine.Instance().DB.Where("id = ?", id).First(&app).Error
	return
}

func (s *AppService) GetByIDAndTenantID(id, tenant_id string) (*models.App, error) {
	app := new(models.App)
	err := dbengine.Instance().DB.Where("id = ? and tenant_id = ?", id, tenant_id).First(app).Error
	if err != nil {
		mlog.Errorf("get app failed:%v", err)
		return nil, err
	}
	return app, nil
}

// GetList 分页获取App记录
func (s *AppService) GetList(info request.PageInfo) (list []models.App, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := dbengine.Instance().DB.Model(&models.App{})
	var endpoints []models.App

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&endpoints).Error
	return endpoints, total, err
}

func (s *AppService) ModeCompatibleWitAgent(a *models.App) models.AppMode {
	if a.Mode == models.AppMode_CHAT && a.IsAgent() {
		return models.AppMode_AGENT_CHAT
	}
	return a.Mode
}

func (s *AppService) apiBaseUrl() string {
	return confy.Get[string]("SERVICE_API_URL") + "/v1"
}

func (s *AppService) GetPaginateApps(user_id, tenant_id string, args *pbapi.ListAppsRequest) *response.AppPaginationResponse {
	// """
	// Get app list with pagination
	// :param tenant_id: tenant id
	// :param args: request args
	// :return:
	// """
	db := dbengine.Instance().DB.Model(&models.App{}).Where("tenant_id = ? and is_universal = ?", tenant_id, false)

	switch args.Mode {
	case "workflow":
		db = db.Where("`mode` IN ?", []models.AppMode{models.AppMode_WORKFLOW, models.AppMode_COMPLETION})
	case "chat":
		db = db.Where("`mode` IN ?", []models.AppMode{models.AppMode_CHAT, models.AppMode_ADVANCED_CHAT})
	case "agent-chat":
		db = db.Where("`mode` = ?", models.AppMode_AGENT_CHAT)
	case "channel":
		db = db.Where("`mode` = ?", models.AppMode_CHANNEL)
	}

	if args.IsCreatedByMe {
		db = db.Where("created_by == ?", user_id)
	}

	if args.Name != "" {
		db = db.Where("`name` LIKE ?", "%"+args.Name+"%")
	}

	if len(args.TagIds) > 0 {
		targetIDs := (&TagService{}).GetTargetIDsByIDs("app", tenant_id, args.TagIds)
		if len(targetIDs) != 0 {
			// return nil,  // 或者返回错误，取决于你的业务逻辑
			db = db.Where("id in ?", targetIDs)
		}
	}
	rsp := &response.AppPaginationResponse{
		Page:  args.Page,
		Limit: args.Limit,
	}
	{
		db.Count(&rsp.Total)
	}
	if args.Page < 1 {
		args.Page = 1
	}
	// 使用Gorm的分页功能
	var apps []*models.App
	err := db.Order("created_at desc").Limit(int(args.Limit)).Offset(int((args.Page - 1) * args.Limit)).Find(&apps).Error
	if err != nil {
		mlog.Errorf("get app failed:%v", err)
		return nil
	}
	rsp.HasMore = rsp.Total > int64(len(apps))
	for _, app := range apps {
		if rsp.Data == nil {
			rsp.Data = make([]*response.AppPartialResponse, 0)
		}
		val := &response.AppPartialResponse{
			ID:                  app.ID,
			Name:                app.Name,
			MaxActiveRequests:   app.MaxActiveRequests,
			Description:         app.Description,
			Mode:                app.Mode,
			IconType:            app.IconType,
			Icon:                app.Icon,
			IconBackground:      app.IconBackground,
			IconURL:             app.Icon,
			Tracing:             app.Tracing,
			UseIconAsAnswerIcon: app.UseIconAsAnswerIcon,
			CreatedBy:           app.CreatedBy,
			CreatedAt:           app.CreatedAt.Unix(),
			UpdatedBy:           app.UpdatedBy,
			UpdatedAt:           app.UpdatedAt.Unix(),
			Tags:                make([]*response.TagFields, 0),
		}
		app_model_config := app.AppModelConfig()
		if app_model_config != nil {
			val.ModelConfig = &response.ModelConfigFields{
				OpeningStatement:              app_model_config.OpeningStatement,
				SuggestedQuestions:            app_model_config.SuggestedQuestionsList(),
				SuggestedQuestionsAfterAnswer: app_model_config.SuggestedQuestionsAfterAnswerDict(),
				SpeechToText:                  app_model_config.SpeechToTextDict(),
				TextToSpeech:                  app_model_config.TextToSpeechDict(),
				RetrieverResource:             app_model_config.RetrieverResourceDict(),
				AnnotationReply:               app_model_config.AnnotationReplyDict(),
				MoreLikeThis:                  app_model_config.MoreLikeThisDict(),
				SensitiveWordAvoidance:        app_model_config.SensitiveWordAvoidanceDict(),
				ExternalDataTools:             app_model_config.ExternalDataToolsList(),
				Model:                         app_model_config.ModelDict(),
				UserInputForm:                 app_model_config.UserInputFormList(),
				DatasetQueryVariable:          app_model_config.DatasetQueryVariable,
				PrePrompt:                     app_model_config.PrePrompt,
				AgentMode:                     app_model_config.AgentModeDict(),
				PromptType:                    app_model_config.PromptType,
				ChatPromptConfig:              app_model_config.ChatPromptConfigDict(),
				CompletionPromptConfig:        app_model_config.CompletionPromptConfigDict(),
				DatasetConfigs:                app_model_config.DatasetConfigsDict(),
				FileUpload:                    app_model_config.FileUploadDict(),
				CreatedBy:                     app_model_config.CreatedBy,
				CreatedAt:                     app_model_config.CreatedAt.Unix(),
				UpdatedBy:                     app_model_config.UpdatedBy,
				UpdatedAt:                     app_model_config.UpdatedAt.Unix(),
			}
		}
		wf := app.Workflow()
		if wf != nil {
			val.Workflow = &response.WorkflowPartialFields{
				ID:        wf.ID,
				CreatedBy: wf.CreatedBy,
				CreatedAt: wf.CreatedAt.Unix(),
				UpdatedBy: wf.UpdatedBy,
				UpdatedAt: wf.UpdatedAt.Unix(),
			}
		}

		for _, v := range app.Tags() {
			if val.Tags == nil {
				val.Tags = make([]*response.TagFields, 0)
			}
			val.Tags = append(val.Tags, &response.TagFields{
				ID:   v.ID,
				Name: v.Name,
				Type: v.Type,
			})
		}

		rsp.Data = append(rsp.Data, val)
	}
	return rsp
}

func (s *AppService) GetAppModel(app_id string, current_user *models.Account, limit_models []models.AppMode) *models.App {
	if app_id == "" {
		mlog.Error("invalid app_id arg")
		panic(exceptions.NewValueError("missing app_id in path parameters"))
	}
	if current_user == nil {
		mlog.Error("invalid current_user arg")
		panic(exceptions.NewValueError("missing current_user"))
	}
	app := new(models.App)
	err := dbengine.Instance().DB.Model(&models.App{}).Where("id = ? and tenant_id = ? and status = ?", app_id, current_user.CurrentTenantID(), "normal").First(app).Error
	if err != nil || app.ID == "" {
		mlog.Error("get app failed:", err)
		panic(httpexceptions.NewAppNotFoundError(""))
	}
	if app.Mode == models.AppMode_CHANNEL {
		mlog.Error("app is channel")
		panic(httpexceptions.NewAppNotFoundError(""))
	}
	if len(limit_models) > 0 {
		if !slices.Contains(limit_models, app.Mode) {
			panic(httpexceptions.NewAppNotFoundError(fmt.Sprintf("App mode is not in the supported list: %v", limit_models)))
		}
	}
	return app
}

func (s *AppService) GetApp(id string, app *models.App) *models.App {
	// get original app model config
	if app.Mode == models.AppMode_AGENT_CHAT || app.IsAgent() {
		model_config := app.AppModelConfig()
		agent_mode := model_config.AgentModeDict()
		var tools []map[string]any
		if _, ok := agent_mode["tools"]; ok {
			if _, ok := agent_mode["tools"].([]map[string]any); ok {
				tools = agent_mode["tools"].([]map[string]any)
			}
		}
		// decrypt agent tool parameters if it's secret-input
		for _, tool := range tools {
			if len(tool) <= 3 {
				continue
			}
			bindata, _ := json.Marshal(tool)
			agent_tool_entity := new(agententities.AgentToolEntity)
			err := json.Unmarshal(bindata, agent_tool_entity)
			if err != nil {
				mlog.Errorf("json.Unmarshal(%#v) to AgentToolEntity failed:%v", tool, err)
				continue
			}
			// get tool
			// try:
			// 	tool_runtime = ToolManager.get_agent_tool_runtime(
			// 		tenant_id=current_user.current_tenant_id,
			// 		app_id=app.id,
			// 		agent_tool=agent_tool_entity,
			// 	)
			// 	manager = ToolParameterConfigurationManager(
			// 		tenant_id=current_user.current_tenant_id,
			// 		tool_runtime=tool_runtime,
			// 		provider_name=agent_tool_entity.provider_id,
			// 		provider_type=agent_tool_entity.provider_type,
			// 		identity_id=f"AGENT.{app.id}",
			// 	)

			// 	// get decrypted parameters
			// 	if agent_tool_entity.tool_parameters{
			// 		parameters = manager.decrypt_tool_parameters(agent_tool_entity.tool_parameters or {})
			// 		masked_parameter = manager.mask_tool_parameters(parameters or {})
			// 		}else{
			// 		masked_parameter = {}
			// 	}
			// 	// override tool parameters
			// 	tool["tool_parameters"] = masked_parameter
			// except Exception as e:
			// 	pass
		}
		// override agent mode
		bindata, _ := json.Marshal(agent_mode)
		model_config.AgentMode = string(bindata)

		// class ModifiedApp(App):
		// 	"""
		// 	Modified App class
		// 	"""

		// 	def __init__(self, app):
		// 		self.__dict__.update(app.__dict__)

		// 	@property
		// 	def app_model_config(self):
		// 		return model_config

		// app = ModifiedApp(app)
	}
	return app
}

func (s *AppService) CreateApp(tenant_id string, in *pbapi.CreateAppRequest, account *models.Account) *models.App {
	app_mode := models.AppMode(in.Mode)
	app_template := models.DEFAULT_APP_TEMPLATES[app_mode]

	// get model config
	default_model_config := app_template["model_config"]

	if len(default_model_config) > 0 {
		if _, ok := default_model_config["model"]; ok {
			// get model provider
			var model_instance *modelmanager.ModelInstance

			// get default model instance
			model_instance = func(tenant_id string) *modelmanager.ModelInstance {
				defer func() {
					if r := recover(); r != nil {
						mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
						if _, ok := r.(*exceptions.ProviderTokenNotInitError); ok {
							model_instance = nil
						} else if _, ok := r.(*exceptions.LLMBadRequestError); ok {
							model_instance = nil
						} else if _, ok := r.(error); ok {
							mlog.Errorf("Get default model instance failed, tenant_id: %s", tenant_id)
							model_instance = nil
						} else {
							panic(r)
						}
					}

				}()
				return manageres.Instance.Model.GetDefaultModelInstance(
					tenant_id, modelruntimeenumtypes.Model_LLM,
				)
			}(account.CurrentTenantID())
			default_model_dict := map[string]any{}

			if _, ok := default_model_config["model"]; ok {
				if _, ok := default_model_config["model"].(map[string]any); ok {
					default_model_dict = default_model_config["model"].(map[string]any)
				}
			}

			if model_instance != nil {
				default_model_name := ""
				default_model_provider := ""

				if _, ok := default_model_dict["name"]; ok {
					if _, ok := default_model_dict["name"].(string); ok {
						default_model_name = default_model_dict["name"].(string)
					}
				}
				if _, ok := default_model_dict["provider"]; ok {
					if _, ok := default_model_dict["provider"].(string); ok {
						default_model_provider = default_model_dict["provider"].(string)
					}
				}

				if model_instance.Model == default_model_name &&
					model_instance.Provider == default_model_provider {
					default_model_dict = default_model_dict
				} else {
					// llm_model := model_instance.ModelTypeInstance.(modelruntimeentities.LargeLanguageModeler)
					model_schema := model_instance.ModelTypeInstance.GetModelSchema(model_instance.ModelTypeInstance, model_instance.Model, model_instance.Credentials)
					if model_schema == nil {
						panic(exceptions.NewValueError(fmt.Sprintf("model schema not found for model %s", model_instance.Model)))
					}
					default_model_dict = map[string]any{
						"provider":          model_instance.Provider,
						"name":              model_instance.Model,
						"mode":              model_schema.ModelProperties[modelruntimeenumtypes.ModelPropertyKey_MODE],
						"completion_params": map[string]any{},
					}
				}
			} else {
				provider, model := manageres.Instance.Model.GetDefaultProviderModelName(
					account.CurrentTenantID(), modelruntimeenumtypes.Model_LLM,
				)
				if default_model_dict == nil {
					default_model_dict = make(map[string]any)
				}
				default_model_dict["provider"] = provider
				default_model_dict["name"] = model
			}
			bindata, _ := json.Marshal(default_model_dict)
			default_model_config["model"] = string(bindata)
		}
	}
	app := models.NewApp(app_template["app"])
	app.ID = uuid.NewV4().String()
	app.Name = in.Name
	app.Description = in.Description
	app.Mode = models.AppMode(in.Mode)
	app.IconType = in.IconType
	if app.IconType == "" {
		app.IconType = "emoji"
	}
	app.Icon = in.Icon
	app.IconBackground = in.IconBackground
	app.APIRph = int(in.ApiRph)
	app.IconBackground = in.IconBackground
	app.TenantID = tenant_id
	app.APIRpm = int(in.ApiRpm)
	app.CreatedBy = account.ID
	app.UpdatedBy = account.ID
	now := time.Now()
	app.CreatedAt = &now
	app.UpdatedAt = &now

	dbengine.Instance().DB.Create(app)

	if len(default_model_config) > 0 {
		app_model_config := models.NewAppModelConfig(default_model_config)
		app_model_config.ID = uuid.NewV4().String()
		app_model_config.AppID = app.ID
		app_model_config.CreatedBy = account.ID
		app_model_config.UpdatedBy = account.ID
		dbengine.Instance().DB.Create(app_model_config)

		app.AppModelConfigID = app_model_config.ID
		dbengine.Instance().DB.Updates(&models.App{ID: app.ID, AppModelConfigID: app_model_config.ID})
	}

	events.Instance.AppWasCreatedSig.Emit(app, account)

	return app
}

func (s *AppService) UpdateAppName(app *models.App, name string, acc *models.Account) *models.App {
	now := time.Now()
	dbengine.Instance().DB.Updates(&models.App{ID: app.ID, Name: name, UpdatedBy: acc.ID, UpdatedAt: &now})
	app.Name = name
	app.UpdatedBy = acc.ID
	app.UpdatedAt = &now
	return app
}
func (s *AppService) UpdateAppIcon(app *models.App, icon, icon_background string, acc *models.Account) *models.App {
	now := time.Now()
	dbengine.Instance().DB.Updates(&models.App{ID: app.ID, Icon: icon, IconBackground: icon_background, UpdatedBy: acc.ID, UpdatedAt: &now})
	app.Icon = icon
	app.IconBackground = icon_background
	app.UpdatedBy = acc.ID
	app.UpdatedAt = &now
	return app
}

func (s *AppService) UpdateSiteStatus(app *models.App, enable_site bool, acc *models.Account) *models.App {
	if app.EnableSite == enable_site {
		return app
	}
	now := time.Now()
	dbengine.Instance().DB.Updates(&models.App{ID: app.ID, EnableSite: enable_site, UpdatedBy: acc.ID, UpdatedAt: &now})
	app.EnableSite = enable_site
	app.UpdatedBy = acc.ID
	app.UpdatedAt = &now
	return app
}
func (s *AppService) UpdateApiStatus(app *models.App, enable_api bool, acc *models.Account) *models.App {
	if app.EnableAPI == enable_api {
		return app
	}
	now := time.Now()
	dbengine.Instance().DB.Updates(&models.App{ID: app.ID, EnableAPI: enable_api, UpdatedBy: acc.ID, UpdatedAt: &now})
	app.EnableAPI = enable_api
	app.UpdatedBy = acc.ID
	app.UpdatedAt = &now
	return app
}
