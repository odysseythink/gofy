package router

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"mlib.com/confy"
	v1 "mlib.com/gofy/server/main/link/api/v1"
	"mlib.com/gofy/server/main/link/middleware"
	"mlib.com/mlog"
)

// type RouterGroup struct {
// 	System  system.RouterGroup
// 	Example example.RouterGroup
// 	Bot     bot.RouterGroup
// }

// var RouterGroupApp = new(RouterGroup)

type RouterInitFunc func(*gin.RouterGroup)

func getAllFile(path string, filter string, s []string) ([]string, error) {
	if s == nil {
		s = make([]string, 0)
	}
	rd, err := os.ReadDir(path)
	if err != nil {
		mlog.Errorf("read dir(%s) failed:%v", path, err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fulldir := filepath.Join(path, fi.Name())
			s, err := getAllFile(fulldir, filter, s)
			if err != nil {
				mlog.Errorf("read dir(%s) failed:%v", path, err)
				return s, err
			}
		} else {
			var extlist []string
			if filter != "" {
				extlist = strings.Split(filter, "|")
			}
			if len(extlist) > 0 {
				for _, ext := range extlist {
					if strings.HasSuffix(fi.Name(), ext) {
						s = append(s, filepath.Join(path, fi.Name()))
						break
					}
				}
			} else {
				s = append(s, filepath.Join(path, fi.Name()))
			}
		}
	}
	return s, nil
}

func InitRouters() *gin.Engine {
	r := gin.Default()

	// 如果想要不使用nginx代理前端网页，可以修改 web/.env.production 下的
	// VUE_APP_BASE_API = /
	// VUE_APP_BASE_PATH = http://localhost
	// 然后执行打包命令 npm run build。在打开下面4行注释
	// r.LoadHTMLGlob("./dist/*.html") // npm打包成dist的路径
	// r.Static("/favicon.ico", "./dist/favicon.ico")
	// r.Static("/static", "./dist/assets")   // dist里面的静态资源
	// r.Static("/assets", "./dist/assets")   // dist里面的静态资源
	// r.StaticFile("/", "./dist/index.html") // 前端网页入口页面
	rd, err := os.ReadDir(filepath.Join("model_runtime", "model_provides"))
	if err != nil {
		mlog.Errorf("read dir(model_provides) failed:%v", err)
	} else {
		for _, fi := range rd {
			if fi.IsDir() && fi.Name() != "base" {
				srd, err := os.ReadDir(filepath.Join("model_runtime", "model_provides", fi.Name(), "_assets"))
				if err == nil {
					for _, sfi := range srd {
						if !sfi.IsDir() && (strings.HasSuffix(strings.ToLower(sfi.Name()), ".png") || strings.HasSuffix(strings.ToLower(sfi.Name()), ".jpg")) {
							r.StaticFile("/model_runtime/model_provides/"+fi.Name()+"/assets/"+sfi.Name(), filepath.Join("model_runtime", "model_provides", fi.Name(), "_assets", sfi.Name()))
						}
					}
				}
			}
		}
	}

	// r.StaticFS(confy.Get[string]("local.path") , http.Dir(confy.Get[string]("local.store-path"))) // 为用户头像和文件提供静态地址
	// Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")
	// 跨域，如需跨域可以打开下面的注释
	r.Use(middleware.Cors()) // 直接放行全部跨域请求
	// r.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	//mlog.Info("use middleware cors")
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	mlog.Info("register swagger handler")
	// 方便统一添加路由组前缀 多服务器上线使用
	r.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")                                // 允许跨域
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // 允许的方法
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")     // 允许的头信息
			c.Header("Allow", "GET, POST, PUT, DELETE, OPTIONS")                        // 允许的HTTP方法
			c.AbortWithStatus(http.StatusOK)                                            // 直接返回204状态码
			return
		}
	})

	consoleRouter := r.Group("console")
	apiRouter := consoleRouter.Group("api")
	apiRouter.Use(middleware.ErrorMiddleware())
	{
		apiRouter.GET("setup", v1.ApiGroupApp.SetupApi.GetSetupStatus)
		apiRouter.POST("setup", v1.ApiGroupApp.SetupApi.Setup)
		apiRouter.POST("login", v1.ApiGroupApp.LoginApi.Login)
		apiRouter.GET("system-features", v1.ApiGroupApp.FeatureApi.ListSystem)
		apiRouter.GET("version", v1.ApiGroupApp.VersionApi.GetVersion)
	}
	{
		authapiRouter := apiRouter
		authapiRouter.Use(middleware.JWTAuth())
		authapiRouter.GET("logout", v1.ApiGroupApp.LoginApi.Logout)
		authapiRouter.GET("files/upload", v1.ApiGroupApp.FileApi.GetUploadConfig)

		authapiRouter.GET("apps/:app_id", v1.ApiGroupApp.AppApi.Find)
		authapiRouter.PUT("apps/:app_id", v1.ApiGroupApp.AppApi.Update)
		authapiRouter.DELETE("apps/:app_id", v1.ApiGroupApp.AppApi.Delete)
		authapiRouter.GET("apps", v1.ApiGroupApp.AppApi.GetList)
		authapiRouter.POST("apps", v1.ApiGroupApp.AppApi.Create)
		authapiRouter.POST("apps/imports", v1.ApiGroupApp.AppApi.Import)
		authapiRouter.GET("apps/:app_id/workflow-app-logs", v1.ApiGroupApp.AppApi.GetWorkflowAppLogList)
		authapiRouter.POST("apps/:app_id/name", v1.ApiGroupApp.AppApi.SetName)
		authapiRouter.POST("apps/:app_id/icon", v1.ApiGroupApp.AppApi.SetIcon)
		authapiRouter.POST("apps/:app_id/copy", v1.ApiGroupApp.AppApi.Copy)
		authapiRouter.GET("apps/:app_id/export", v1.ApiGroupApp.AppApi.Export)
		authapiRouter.POST("apps/:app_id/site-enable", v1.ApiGroupApp.AppApi.UpdateSiteStatus)
		authapiRouter.POST("apps/:app_id/api-enable", v1.ApiGroupApp.AppApi.UpdateApiStatus)
		authapiRouter.GET("apps/:app_id/trace", v1.ApiGroupApp.AppApi.GetTrace)
		authapiRouter.POST("apps/:app_id/trace", v1.ApiGroupApp.AppApi.SetTrace)
		authapiRouter.GET("apps/:app_id/trace-config", v1.ApiGroupApp.OpsTraceApi.GetTraceAppConfig)

		authapiRouter.GET("apps/:app_id/workflows/draft", v1.ApiGroupApp.DraftWorkflowApi.Find)
		authapiRouter.POST("apps/:app_id/workflows/draft", v1.ApiGroupApp.DraftWorkflowApi.Sync)
		authapiRouter.POST("apps/:app_id/workflows/draft/run", v1.ApiGroupApp.DraftWorkflowApi.Run)
		authapiRouter.POST("apps/:app_id/advanced-chat/workflows/draft/run", v1.ApiGroupApp.DraftWorkflowApi.Run)
		authapiRouter.POST("apps/:app_id/workflows/draft/nodes/:node_id/run", v1.ApiGroupApp.DraftWorkflowApi.NodeRun)
		authapiRouter.GET("apps/:app_id/workflows/draft/config", v1.ApiGroupApp.DraftWorkflowApi.Config)
		authapiRouter.GET("apps/:app_id/workflows/default-workflow-block-configs", v1.ApiGroupApp.DraftWorkflowApi.DefaultBlockConfigs)
		authapiRouter.GET("apps/:app_id/workflows/default-workflow-block-configs/:block_type", v1.ApiGroupApp.DraftWorkflowApi.DefaultBlockConfig)
		authapiRouter.GET("apps/:app_id/workflows/publish", v1.ApiGroupApp.DraftWorkflowApi.GetPublished)
		authapiRouter.POST("apps/:app_id/workflows/publish", v1.ApiGroupApp.DraftWorkflowApi.SetPublished)

		authapiRouter.GET("apps/:app_id/api-keys", v1.ApiGroupApp.ApiKeyApi.GetApiKeyListResource)
		authapiRouter.POST("apps/:app_id/api-keys", v1.ApiGroupApp.ApiKeyApi.SetApiKeyListResource)
		authapiRouter.GET("datasets/:app_id/api-keys", v1.ApiGroupApp.ApiKeyApi.GetApiKeyListResource)
		authapiRouter.POST("datasets/:app_id/api-keys", v1.ApiGroupApp.ApiKeyApi.SetApiKeyListResource)
		authapiRouter.DELETE("apps/:app_id/api-keys/:api_key_id", v1.ApiGroupApp.ApiKeyApi.DelApiKeyResource)
		authapiRouter.DELETE("datasets/:app_id/api-keys/:api_key_id", v1.ApiGroupApp.ApiKeyApi.DelApiKeyResource)

		authapiRouter.GET("account/profile", v1.ApiGroupApp.AccountApi.Profile)
		authapiRouter.POST("/account/name", v1.ApiGroupApp.AccountApi.Update)
		authapiRouter.POST("/account/avatar", v1.ApiGroupApp.AccountApi.Update)
		authapiRouter.POST("/account/interface-language", v1.ApiGroupApp.AccountApi.Update)
		authapiRouter.POST("/account/interface-theme", v1.ApiGroupApp.AccountApi.Update)
		authapiRouter.POST("/account/timezone", v1.ApiGroupApp.AccountApi.Update)

		authapiRouter.GET("apps/:app_id/workflow/statistics/daily-conversations", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/workflow/statistics/daily-terminals", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/workflow/statistics/token-costs", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/workflow/statistics/average-app-interactions", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/statistics/daily-messages", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/statistics/daily-conversations", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/statistics/daily-end-users", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/statistics/token-costs", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/statistics/average-session-interactions", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/statistics/user-satisfaction-rate", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/statistics/average-response-time", v1.ApiGroupApp.StatisticApi.Statistic)
		authapiRouter.GET("apps/:app_id/statistics/tokens-per-second", v1.ApiGroupApp.StatisticApi.Statistic)

		authapiRouter.GET("apps/:app_id/chat-messages/:message_id/suggested-questions", v1.ApiGroupApp.MessageApi.GetSuggestedQuestion)
		authapiRouter.GET("apps/:app_id/chat-messages", v1.ApiGroupApp.MessageApi.ChatMessageList)
		authapiRouter.POST("apps/:app_id/feedbacks", v1.ApiGroupApp.MessageApi.MessageFeedback)
		authapiRouter.POST("apps/:app_id/annotations", v1.ApiGroupApp.MessageApi.SetMessageAnnotation)
		authapiRouter.GET("apps/:app_id/annotations/count", v1.ApiGroupApp.MessageApi.MessageAnnotationCount)
		authapiRouter.GET("apps/:app_id/messages/:message_id", v1.ApiGroupApp.MessageApi.GetMessage)

		authapiRouter.GET("apps/:app_id/chat-conversations", v1.ApiGroupApp.ConversationApi.GetChatConversationPagination)
		authapiRouter.GET("apps/:app_id/chat-conversations/:conversation_id", v1.ApiGroupApp.ConversationApi.ChatConversationDetail)
		authapiRouter.DELETE("apps/:app_id/chat-conversations/:conversation_id", v1.ApiGroupApp.ConversationApi.DelChatConversation)

		authapiRouter.GET("apps/:app_id/advanced-chat/workflow-runs", v1.ApiGroupApp.WorkflowRunApi.AppList)
		authapiRouter.GET("apps/:app_id/workflow-runs", v1.ApiGroupApp.WorkflowRunApi.AppList)
		authapiRouter.GET("apps/:app_id/workflow-runs/:workflow_run_id", v1.ApiGroupApp.WorkflowRunApi.Detail)
		authapiRouter.GET("apps/:app_id/workflow-runs/:workflow_run_id/node-executions", v1.ApiGroupApp.WorkflowRunApi.Detail)

		authapiRouter.POST("workspaces/current/default-model", v1.ApiGroupApp.ModelsApi.SetDefaultModel)
		authapiRouter.GET("workspaces/current/default-model", v1.ApiGroupApp.ModelsApi.DefaultModel)
		authapiRouter.GET("workspaces/current/models/model-types/:model_type", v1.ApiGroupApp.ModelsApi.Available)
		authapiRouter.GET("workspaces/current/model-providers/:provider/models", v1.ApiGroupApp.ModelsApi.GetModel)
		authapiRouter.POST("workspaces/current/model-providers/:provider/models", v1.ApiGroupApp.ModelsApi.SetModel)
		authapiRouter.PATCH("workspaces/current/model-providers/:provider/models/enable", v1.ApiGroupApp.ModelsApi.EnableModel)
		authapiRouter.PATCH("workspaces/current/model-providers/:provider/models/disable", v1.ApiGroupApp.ModelsApi.EnableModel)
		authapiRouter.GET("workspaces/current/model-providers/:provider/models/parameter-rules", v1.ApiGroupApp.ModelsApi.GetParameterRules)

		authapiRouter.GET("workspaces/current/model-providers", v1.ApiGroupApp.ModelProvideApi.ProviderList)
		authapiRouter.POST("workspaces/current/model-providers/:provider", v1.ApiGroupApp.ModelProvideApi.Update)
		authapiRouter.DELETE("workspaces/current/model-providers/:provider", v1.ApiGroupApp.ModelProvideApi.Delete)
		authapiRouter.GET("workspaces/current/model-providers/:provider/credentials", v1.ApiGroupApp.ModelProvideApi.GetCredentials)

		authapiRouter.GET("all-workspaces", v1.ApiGroupApp.WorkspaceApi.List)
		authapiRouter.GET("workspaces", v1.ApiGroupApp.WorkspaceApi.All)
		authapiRouter.GET("workspaces/current", v1.ApiGroupApp.WorkspaceApi.GetCurrent)

		authapiRouter.GET("features", v1.ApiGroupApp.FeatureApi.List)

		authapiRouter.GET("workspaces/current/members", v1.ApiGroupApp.MemberApi.List)

		authapiRouter.GET("datasets/retrieval-setting", v1.ApiGroupApp.DatasetApi.RetrievalSetting)
		authapiRouter.GET("datasets", v1.ApiGroupApp.DatasetApi.DatasetList)
		authapiRouter.GET("datasets/external-knowledge-api", v1.ApiGroupApp.DatasetApi.ExternalKnowledgeApiList)

		// Dataset document & segment routes
		docApi := v1.ApiGroupApp.DatasetDocumentApi
		authapiRouter.GET("datasets/:dataset_id/documents", docApi.GetDocuments)
		authapiRouter.DELETE("datasets/:dataset_id/documents/:document_id", docApi.DeleteDocument)
		authapiRouter.POST("datasets/:dataset_id/documents/:document_id/rename", docApi.RenameDocument)
		authapiRouter.PATCH("datasets/:dataset_id/documents/:document_id/pause", docApi.PauseDocument)
		authapiRouter.PATCH("datasets/:dataset_id/documents/:document_id/recover", docApi.RecoverDocument)
		authapiRouter.POST("datasets/:dataset_id/documents/status", docApi.BatchUpdateStatus)
		authapiRouter.GET("datasets/:dataset_id/documents/:document_id/segments", docApi.GetSegments)
		authapiRouter.POST("datasets/:dataset_id/documents/:document_id/segments", docApi.CreateSegment)
		authapiRouter.DELETE("datasets/:dataset_id/documents/:document_id/segments/:segment_id", docApi.DeleteSegment)
		authapiRouter.POST("datasets/:dataset_id/hit-testing", docApi.HitTesting)

		authapiRouter.GET("tags", v1.ApiGroupApp.TagApi.List)
		authapiRouter.POST("tags", v1.ApiGroupApp.TagApi.Add)
		authapiRouter.PATCH("tags/:tag_id", v1.ApiGroupApp.TagApi.Update)
		authapiRouter.DELETE("tags/:tag_id", v1.ApiGroupApp.TagApi.Del)
		authapiRouter.POST("tag-bindings/create", v1.ApiGroupApp.TagApi.AddTagBinding)
		authapiRouter.POST("tag-bindings/remove", v1.ApiGroupApp.TagApi.DelTagBinding)

		authapiRouter.GET("code-based-extension", v1.ApiGroupApp.ExtensionApi.GetCodeBasedExtension)

		authapiRouter.GET("workspaces/current/tool-labels", v1.ApiGroupApp.ToolsApi.ListToolLabels)
		authapiRouter.GET("workspaces/current/tool-providers", v1.ApiGroupApp.ToolsApi.ListToolProvider)
		authapiRouter.GET("workspaces/current/tools/:tool_type", v1.ApiGroupApp.ToolsApi.GetToolList)
		authapiRouter.GET("workspaces/current/tools/mcp", v1.ApiGroupApp.ToolsApi.GetMCPToolList)

		authapiRouter.GET("workspaces/current/plugin/preferences/fetch", v1.ApiGroupApp.PluginApi.PluginFetchPreferences)
		authapiRouter.GET("workspaces/current/plugin/tasks", v1.ApiGroupApp.PluginApi.PluginFetchInstallTasks)

		authapiRouter.POST("rule-generate", v1.ApiGroupApp.RuleGenerateApi.RuleGenerate)
	}
	{
		noauthapiRouter := r.Group("v1")
		noauthapiRouter.Use(middleware.AppAuth())
		noauthapiRouter.POST("workflows/run", v1.ApiGroupApp.WorkflowRunApi.AppRun)
		noauthapiRouter.GET("workflows/run/:workflow_run_id", v1.ApiGroupApp.WorkflowRunApi.GetDetailOnlyByApp)
		// authapiRouter.POST("/workflows/tasks/<string:task_id>/stop", v1.ApiGroupApp.WorkflowRunApi.WorkflowTaskStopApi)
		// authapiRouter.POST("/workflows/logs", v1.ApiGroupApp.WorkflowRunApi.WorkflowAppLogApi)
		noauthapiRouter.POST("chat-messages", v1.ApiGroupApp.WorkflowRunApi.AppRun)

		// Web/end-user API routes
		webApi := v1.ApiGroupApp.WebApi
		noauthapiRouter.GET("conversations", webApi.GetConversations)
		noauthapiRouter.DELETE("conversations/:conversation_id", webApi.DeleteConversation)
		noauthapiRouter.POST("conversations/:conversation_id/name", webApi.RenameConversation)
		noauthapiRouter.GET("messages", webApi.GetMessages)
		noauthapiRouter.POST("messages/:message_id/feedbacks", webApi.MessageFeedback)
		noauthapiRouter.GET("messages/:message_id/suggested", webApi.GetSuggestedQuestions)
		noauthapiRouter.POST("files/upload", webApi.UploadFile)
		noauthapiRouter.GET("parameters", webApi.GetAppParameters)
		noauthapiRouter.GET("meta", webApi.GetAppMeta)
	}

	// GitHub API proxy (replaces frontend Next.js API route)
	r.GET("/api/repos/:owner/:repo/releases", func(c *gin.Context) {
		owner := c.Param("owner")
		repo := c.Param("repo")

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token := confy.GetWithDefault[string]("github.access_token", "")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		req.Header.Set("Accept", "application/vnd.github+json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	mlog.Info("router register success")
	return r
}
