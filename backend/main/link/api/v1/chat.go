package v1

type ChatApi struct {
}

// func (api *ChatApi) Chat(c *gin.Context) {
// 	rawapp, _ := c.Get("app_model")
// 	app := rawapp.(*models.App)

// 	if !slices.Contains([]models.AppMode{models.AppMode_ADVANCED_CHAT, models.AppMode_CHAT, models.AppMode_AGENT_CHAT}, app.Mode) {
// 		panic(httpexceptions.NewNotChatAppError())
// 	}
// 	type Req struct {
// 		Inputs           map[string]any `json:"inputs"`
// 		Query            string         `json:"query"`
// 		ResponseMode     string         `json:"response_mode"`
// 		ConversationID   string         `json:"conversation_id"`
// 		RetrieverFrom    string         `json:"retriever_from"`
// 		AutoGenerateName bool           `json:"auto_generate_name"`
// 		User             string         `json:"user"`
// 	}
// 	req := Req{
// 		RetrieverFrom:    "dev",
// 		AutoGenerateName: true,
// 	}
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
// 		return
// 	}
// 	if req.Query == "" {
// 		mlog.Error("query user")
// 		response.InvalidArgError(c)
// 		return
// 	}
// 	if req.User == "" {
// 		mlog.Error("missing user")
// 		response.InvalidArgError(c)
// 		return
// 	}
// 	if !slices.Contains([]string{"blocking", "streaming"}, req.ResponseMode) {
// 		mlog.Errorf("invalid response_mode=%v", req.ResponseMode)
// 		response.InvalidArgError(c)
// 		return
// 	}
// 	end_user := models.CreateOrGetEndUserByUserID(app, req.User)

// 	streaming := req.ResponseMode == "streaming"
// 	defer func() {
// 		if r := recover(); r != nil {
// 			if _, ok := r.(*exceptions.ConversationNotExistsError); ok {
// 				panic(httpexceptions.NewNotFound("Conversation Not Exists."))
// 			} else if exp, ok := r.(*httpexceptions.ConversationCompletedError); ok {
// 				panic(exp)
// 			} else if _, ok := r.(*exceptions.AppModelConfigBrokenError); ok {
// 				mlog.Errorf("App model config broken.")
// 				panic(httpexceptions.NewAppUnavailableError())
// 			} else if exp, ok := r.(*exceptions.ProviderTokenNotInitError); ok {
// 				panic(httpexceptions.NewProviderNotInitializeError(exp.Error()))
// 			} else if _, ok := r.(*exceptions.QuotaExceededError); ok {
// 				panic(httpexceptions.NewProviderQuotaExceededError())
// 			} else if _, ok := r.(*exceptions.ModelCurrentlyNotSupportError); ok {
// 				panic(httpexceptions.NewProviderModelCurrentlyNotSupportError())
// 			} else if exp, ok := r.(*modelruntimeexceptions.InvokeError); ok {
// 				panic(httpexceptions.NewCompletionRequestError(exp.Error()))
// 			} else if exp, ok := r.(*exceptions.ValueError); ok {
// 				panic(exp)
// 			} else {
// 				mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
// 				panic(httpexceptions.NewInternalServerError(""))
// 			}
// 		}
// 	}()
// 	bindata, _ := json.Marshal(req)
// 	args := map[string]any{}
// 	json.Unmarshal(bindata, &args)
// 	rsp, rspiter := services.ServiceGroupApp.AppGenerate.Generate(app, end_user, args, appenumtypes.InvokeFrom_SERVICE_API, streaming)
// 	mlog.Debugf("***********realrsp=%#v", rsp)
// 	if rsp != nil {
// 		c.JSON(http.StatusOK, rsp)
// 		return
// 	}

// 	for item := range rspiter {
// 		mlog.Debugf("-----send:%s", string(item))
// 		c.Data(http.StatusOK, "text/event-stream", []byte(item))
// 	}
// }
