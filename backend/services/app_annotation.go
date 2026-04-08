package services

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/cache"
	httpexceptions "mlib.com/gofy/server/core/exceptions/http"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type AppAnnotationService struct {
}

func (service *AppAnnotationService) UpInsertAppAnnotationFromMessage(
	args map[string]any,
	app_id string,
	current_user *models.Account,
) *models.MessageAnnotation {
	// get app info
	app := new(models.App)
	err := dbengine.Instance().DB.Model(&models.App{}).Where("id=? and tenant_id=? and status = ?", app_id, current_user.CurrentTenantID(), "normal").First(app).Error
	if err != nil {
		mlog.Errorf("get app failed:%v", err)
		app = nil
	}

	if app == nil {
		panic(httpexceptions.NewNotFound("App not found"))
	}
	message_id := ""
	if _, ok := args["message_id"]; ok {
		if _, ok := args["message_id"].(string); ok {
			message_id = args["message_id"].(string)
		}
	}
	var annotation *models.MessageAnnotation
	if message_id != "" {
		// get message info
		message := new(models.Message)
		err := dbengine.Instance().DB.Model(&models.Message{}).Where("id=? and app_id=?", message_id, app_id).First(message).Error
		if err != nil {
			mlog.Errorf("get Message failed:%v", err)
			message = nil
		}

		if message == nil {
			panic(httpexceptions.NewNotFound("Message Not Exists."))
		}

		annotation = message.Annotation()
		// save the message annotation
		if annotation != nil {
			if _, ok := args["answer"]; ok {
				if _, ok := args["answer"].(string); ok {
					annotation.Content = args["answer"].(string)
				}
			}
			if _, ok := args["question"]; ok {
				if _, ok := args["question"].(string); ok {
					annotation.Question = args["question"].(string)
				}
			}

		} else {
			now := time.Now()
			annotation = &models.MessageAnnotation{
				ID:             uuid.NewV4().String(),
				AppID:          app_id,
				ConversationID: message.ConversationID,
				MessageID:      message.ID,
				// Content       :,
				AccountID: current_user.ID,
				CreatedAt: &now,
				UpdatedAt: &now,
			}
			if _, ok := args["answer"]; ok {
				if _, ok := args["answer"].(string); ok {
					annotation.Content = args["answer"].(string)
				}
			}
			if _, ok := args["question"]; ok {
				if _, ok := args["question"].(string); ok {
					annotation.Question = args["question"].(string)
				}
			}
		}
	} else {
		now := time.Now()
		annotation = &models.MessageAnnotation{
			ID:        uuid.NewV4().String(),
			AppID:     app_id,
			AccountID: current_user.ID,
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		if _, ok := args["answer"]; ok {
			if _, ok := args["answer"].(string); ok {
				annotation.Content = args["answer"].(string)
			}
		}
		if _, ok := args["question"]; ok {
			if _, ok := args["question"].(string); ok {
				annotation.Question = args["question"].(string)
			}
		}
	}
	dbengine.Instance().DB.Save(annotation)
	// if annotation reply is enabled , add annotation to index
	annotation_setting := new(models.AppAnnotationSetting)
	err = dbengine.Instance().DB.Model(&models.AppAnnotationSetting{}).Where("app_id=?", app_id).First(annotation_setting).Error
	if err != nil {
		mlog.Errorf("get AppAnnotationSetting failed:%v", err)
		annotation_setting = nil
	}
	// if annotation_setting != nil{
	//     add_annotation_to_index_task.delay(
	//         annotation.id,
	//         args["question"],
	//         current_user.current_tenant_id,
	//         app_id,
	//         annotation_setting.collection_binding_id,
	//     )
	// }
	return annotation

}
func (service *AppAnnotationService) EnableAppAnnotation(args map[string]any, app_id string) map[string]any {
	enable_app_annotation_key := fmt.Sprintf("enable_app_annotation_%s", app_id)
	cache_result := cache.Instance().GetString(enable_app_annotation_key)
	if cache_result != "" {
		return map[string]any{"job_id": cache_result, "job_status": "processing"}
	}
	// async job
	job_id := uuid.NewV4().String()
	enable_app_annotation_job_key := fmt.Sprintf("enable_app_annotation_job_%s", job_id)
	// send batch add segments task
	cache.Instance().SetNX(enable_app_annotation_job_key, "waiting", -1)
	// enable_annotation_reply_task.delay(
	//     str(job_id),
	//     app_id,
	//     current_user.id,
	//     current_user.current_tenant_id,
	//     args["score_threshold"],
	//     args["embedding_provider_name"],
	//     args["embedding_model_name"],
	// )
	return map[string]any{"job_id": job_id, "job_status": "waiting"}

}

// func(service *AppAnnotationService) DisableAppAnnotation(app_id string) map[string]any{
//         disable_app_annotation_key = "disable_app_annotation_{}".format(str(app_id))
//         cache_result = redis_client.get(disable_app_annotation_key)
//         if cache_result is not None{
//             return {"job_id": cache_result, "job_status": "processing"}
// 		}
//         // async job
//         job_id = str(uuid.uuid4())
//         disable_app_annotation_job_key = "disable_app_annotation_job_{}".format(str(job_id))
//         // send batch add segments task
//         redis_client.setnx(disable_app_annotation_job_key, "waiting")
//         disable_annotation_reply_task.delay(str(job_id), app_id, current_user.current_tenant_id)
//         return {"job_id": job_id, "job_status": "waiting"}

//     }
// func(service *AppAnnotationService) get_annotation_list_by_app_id(app_id string, page: int, limit: int, keyword string){
//         // get app info
//         err := dbengine.Instance().DB.Model(&models.App{}).Where("id=? and tenant_id=? and status = ?", app_id,current_user.CurrentTenant.ID, "normal").First(app).Error
//         if err != nil {
// 			mlog.Errorf("get app failed:%v", err)
// 			app = nil
// 		}

//         if app == nil{
//             panic(httpexceptions.NewNotFound("App not found"))
// 		}
//         if keyword{
//             annotations = (
//                 MessageAnnotation.query.filter(MessageAnnotation.app_id == app_id)
//                 .filter(
//                     or_(
//                         MessageAnnotation.question.ilike("%{}%".format(keyword)),
//                         MessageAnnotation.content.ilike("%{}%".format(keyword)),
//                     )
//                 )
//                 .order_by(MessageAnnotation.created_at.desc(), MessageAnnotation.id.desc())
//                 .paginate(page=page, per_page=limit, max_per_page=100, error_out=False)
//             )
//         }else{
//             annotations = (
//                 MessageAnnotation.query.filter(MessageAnnotation.app_id == app_id)
//                 .order_by(MessageAnnotation.created_at.desc(), MessageAnnotation.id.desc())
//                 .paginate(page=page, per_page=limit, max_per_page=100, error_out=False)
//             )
//         return annotations.items, annotations.total

//     }
// func(service *AppAnnotationService) export_annotation_list_by_app_id(app_id string){
//         // get app info
//         err := dbengine.Instance().DB.Model(&models.App{}).Where("id=? and tenant_id=? and status = ?", app_id,current_user.CurrentTenant.ID, "normal").First(app).Error
//         if err != nil {
// 			mlog.Errorf("get app failed:%v", err)
// 			app = nil
// 		}

//         if app == nil{
//             panic(httpexceptions.NewNotFound("App not found"))
// 		}
//         annotations = (
//             db.session.query(MessageAnnotation)
//             .filter(MessageAnnotation.app_id == app_id)
//             .order_by(MessageAnnotation.created_at.desc())
//             .all()
//         )
//         return annotations

//     }
// func(service *AppAnnotationService) insert_app_annotation_directly(args map[string]any, app_id string) *models.MessageAnnotation{
//         // get app info
//         err := dbengine.Instance().DB.Model(&models.App{}).Where("id=? and tenant_id=? and status = ?", app_id,current_user.CurrentTenant.ID, "normal").First(app).Error
//         if err != nil {
// 			mlog.Errorf("get app failed:%v", err)
// 			app = nil
// 		}

//         if app == nil{
//             panic(httpexceptions.NewNotFound("App not found"))
// 		}

//         annotation = MessageAnnotation(
//             app_id=app.id, content=args["answer"], question=args["question"], account_id=current_user.id
//         )
//         db.session.add(annotation)
//         db.session.commit()
//         // if annotation reply is enabled , add annotation to index
//         annotation_setting = (
//             db.session.query(AppAnnotationSetting).filter(AppAnnotationSetting.app_id == app_id).first()
//         )
//         if annotation_setting{
//             add_annotation_to_index_task.delay(
//                 annotation.id,
//                 args["question"],
//                 current_user.current_tenant_id,
//                 app_id,
//                 annotation_setting.collection_binding_id,
//             )
//         return annotation

//     }
// func(service *AppAnnotationService) update_app_annotation_directly(args map[string]any, app_id string, annotation_id string){
//         // get app info
//         err := dbengine.Instance().DB.Model(&models.App{}).Where("id=? and tenant_id=? and status = ?", app_id,current_user.CurrentTenant.ID, "normal").First(app).Error
//         if err != nil {
// 			mlog.Errorf("get app failed:%v", err)
// 			app = nil
// 		}

//         if app == nil{
//             panic(httpexceptions.NewNotFound("App not found"))
// 		}

//         annotation = db.session.query(MessageAnnotation).filter(MessageAnnotation.id == annotation_id).first()

//         if not annotation{
//             panic(httpexceptions.NewNotFound("Annotation not found")

//         annotation.content = args["answer"]
//         annotation.question = args["question"]

//         db.session.commit()
//         // if annotation reply is enabled , add annotation to index
//         app_annotation_setting = (
//             db.session.query(AppAnnotationSetting).filter(AppAnnotationSetting.app_id == app_id).first()
//         )

//         if app_annotation_setting{
//             update_annotation_to_index_task.delay(
//                 annotation.id,
//                 annotation.question,
//                 current_user.current_tenant_id,
//                 app_id,
//                 app_annotation_setting.collection_binding_id,
//             )

//         return annotation

//     }
// func(service *AppAnnotationService) delete_app_annotation(app_id string, annotation_id string){
//         // get app info
//         err := dbengine.Instance().DB.Model(&models.App{}).Where("id=? and tenant_id=? and status = ?", app_id,current_user.CurrentTenant.ID, "normal").First(app).Error
//         if err != nil {
// 			mlog.Errorf("get app failed:%v", err)
// 			app = nil
// 		}

//         if app == nil{
//             panic(httpexceptions.NewNotFound("App not found"))
// 		}

//         annotation = db.session.query(MessageAnnotation).filter(MessageAnnotation.id == annotation_id).first()

//         if not annotation{
//             panic(httpexceptions.NewNotFound("Annotation not found")

//         db.session.delete(annotation)

//         annotation_hit_histories = (
//             db.session.query(AppAnnotationHitHistory)
//             .filter(AppAnnotationHitHistory.annotation_id == annotation_id)
//             .all()
//         )
//         if annotation_hit_histories{
//             for annotation_hit_history in annotation_hit_histories{
//                 db.session.delete(annotation_hit_history)

//         db.session.commit()
//         // if annotation reply is enabled , delete annotation index
//         app_annotation_setting = (
//             db.session.query(AppAnnotationSetting).filter(AppAnnotationSetting.app_id == app_id).first()
//         )

//         if app_annotation_setting{
//             delete_annotation_index_task.delay(
//                 annotation.id, app_id, current_user.current_tenant_id, app_annotation_setting.collection_binding_id
//             )

//     }
// func(service *AppAnnotationService) batch_import_app_annotations(app_id, file: FileStorage) map[string]any{
//         // get app info
//         err := dbengine.Instance().DB.Model(&models.App{}).Where("id=? and tenant_id=? and status = ?", app_id,current_user.CurrentTenant.ID, "normal").First(app).Error
//         if err != nil {
// 			mlog.Errorf("get app failed:%v", err)
// 			app = nil
// 		}

//         if app == nil{
//             panic(httpexceptions.NewNotFound("App not found"))
// 		}

//         try{
//             // Skip the first row
//             df = pd.read_csv(file)
//             result = []
//             for index, row in df.iterrows(){
//                 content = {"question": row.iloc[0], "answer": row.iloc[1]}
//                 result.append(content)
//             if len(result) == 0{
//                 raise ValueError("The CSV file is empty.")
//             // check annotation limit
//             features = FeatureService.get_features(current_user.current_tenant_id)
//             if features.billing.enabled{
//                 annotation_quota_limit = features.annotation_quota_limit
//                 if annotation_quota_limit.limit < len(result) + annotation_quota_limit.size{
//                     raise ValueError("The number of annotations exceeds the limit of your subscription.")
//             // async job
//             job_id = str(uuid.uuid4())
//             indexing_cache_key = "app_annotation_batch_import_{}".format(str(job_id))
//             // send batch add segments task
//             redis_client.setnx(indexing_cache_key, "waiting")
//             batch_import_annotations_task.delay(
//                 str(job_id), result, app_id, current_user.current_tenant_id, current_user.id
//             )
//         except Exception as e{
//             return {"error_msg" string(e)}
//         return {"job_id": job_id, "job_status": "waiting"}

//     }
// func(service *AppAnnotationService) get_annotation_hit_histories(app_id string, annotation_id string, page, limit){
//         // get app info
//         err := dbengine.Instance().DB.Model(&models.App{}).Where("id=? and tenant_id=? and status = ?", app_id,current_user.CurrentTenant.ID, "normal").First(app).Error
//         if err != nil {
// 			mlog.Errorf("get app failed:%v", err)
// 			app = nil
// 		}

//         if app == nil{
//             panic(httpexceptions.NewNotFound("App not found"))
// 		}

//         annotation = db.session.query(MessageAnnotation).filter(MessageAnnotation.id == annotation_id).first()

//         if not annotation{
//             panic(httpexceptions.NewNotFound("Annotation not found")

//         annotation_hit_histories = (
//             AppAnnotationHitHistory.query.filter(
//                 AppAnnotationHitHistory.app_id == app_id,
//                 AppAnnotationHitHistory.annotation_id == annotation_id,
//             )
//             .order_by(AppAnnotationHitHistory.created_at.desc())
//             .paginate(page=page, per_page=limit, max_per_page=100, error_out=False)
//         )
//         return annotation_hit_histories.items, annotation_hit_histories.total

//     }
// func(service *AppAnnotationService) get_annotation_by_id(annotation_id string) -> MessageAnnotation | None{
//         annotation = db.session.query(MessageAnnotation).filter(MessageAnnotation.id == annotation_id).first()

//         if not annotation{
//             return None
//         return annotation

//     }
// func(service *AppAnnotationService) add_annotation_history(
//         cls,
//         annotation_id string,
//         app_id string,
//         annotation_question string,
//         annotation_content string,
//         query string,
//         user_id string,
//         message_id string,
//         from_source string,
//         score: float,
//     ){
//         // add hit count to annotation
//         db.session.query(MessageAnnotation).filter(MessageAnnotation.id == annotation_id).update(
//             {MessageAnnotation.hit_count: MessageAnnotation.hit_count + 1}, synchronize_session=False
//         )

//         annotation_hit_history = AppAnnotationHitHistory(
//             annotation_id=annotation_id,
//             app_id=app_id,
//             account_id=user_id,
//             question=query,
//             source=from_source,
//             score=score,
//             message_id=message_id,
//             annotation_question=annotation_question,
//             annotation_content=annotation_content,
//         )
//         db.session.add(annotation_hit_history)
//         db.session.commit()

//     }
// func(service *AppAnnotationService) get_app_annotation_setting_by_app_id(app_id string){
//         // get app info
//         err := dbengine.Instance().DB.Model(&models.App{}).Where("id=? and tenant_id=? and status = ?", app_id,current_user.CurrentTenant.ID, "normal").First(app).Error
//         if err != nil {
// 			mlog.Errorf("get app failed:%v", err)
// 			app = nil
// 		}

//         if app == nil{
//             panic(httpexceptions.NewNotFound("App not found"))
// 		}

//         annotation_setting = (
//             db.session.query(AppAnnotationSetting).filter(AppAnnotationSetting.app_id == app_id).first()
//         )
//         if annotation_setting{
//             collection_binding_detail = annotation_setting.collection_binding_detail
//             return {
//                 "id": annotation_setting.id,
//                 "enabled": True,
//                 "score_threshold": annotation_setting.score_threshold,
//                 "embedding_model": {
//                     "embedding_provider_name": collection_binding_detail.provider_name,
//                     "embedding_model_name": collection_binding_detail.model_name,
//                 },
//             }
//         return {"enabled": False}

//     }
// func(service *AppAnnotationService) update_app_annotation_setting(app_id string, annotation_setting_id string, args map[string]any){
//         // get app info
//         err := dbengine.Instance().DB.Model(&models.App{}).Where("id=? and tenant_id=? and status = ?", app_id,current_user.CurrentTenant.ID, "normal").First(app).Error
//         if err != nil {
// 			mlog.Errorf("get app failed:%v", err)
// 			app = nil
// 		}

//         if app == nil{
//             panic(httpexceptions.NewNotFound("App not found"))
// 		}

//         annotation_setting = (
//             db.session.query(AppAnnotationSetting)
//             .filter(
//                 AppAnnotationSetting.app_id == app_id,
//                 AppAnnotationSetting.id == annotation_setting_id,
//             )
//             .first()
//         )
//         if not annotation_setting{
//             panic(httpexceptions.NewNotFound("App annotation not found")
// 		}
//         annotation_setting.score_threshold = args["score_threshold"]
//         annotation_setting.updated_user_id = current_user.id
//         annotation_setting.updated_at = datetime.datetime.now(datetime.UTC).replace(tzinfo=None)
//         db.session.add(annotation_setting)
//         db.session.commit()

//         collection_binding_detail = annotation_setting.collection_binding_detail

//         return {
//             "id": annotation_setting.id,
//             "enabled": True,
//             "score_threshold": annotation_setting.score_threshold,
//             "embedding_model": {
//                 "embedding_provider_name": collection_binding_detail.provider_name,
//                 "embedding_model_name": collection_binding_detail.model_name,
//             },
//         }
//  }
