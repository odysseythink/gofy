package eventhandlers

import (
	"encoding/json"
	"time"

	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	toolnodesentities "github.com/odysseythink/gofy/backend/entities/nodes/tool"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/utils"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

func create_installed_app_when_app_created(app *models.App, acc *models.Account) {
	now := time.Now()
	tenant := app.Tenant()

	installed_app := &models.InstalledApp{
		ID:         uuid.NewV4().String(),
		AppID:      app.ID,
		LastUsedAt: &now,
		CreatedAt:  &now,
	}
	if tenant != nil {
		installed_app.TenantID = tenant.ID
		installed_app.AppOwnerTenantID = tenant.ID
	}
	dbengine.Instance().DB.Create(installed_app)
}

func create_site_record_when_app_created(app *models.App, acc *models.Account) {
	if acc != nil {
		now := time.Now()
		site := &models.Site{
			ID:                     uuid.NewV4().String(),
			AppID:                  app.ID,
			Title:                  app.Name,
			IconType:               app.IconType,
			Icon:                   app.Icon,
			IconBackground:         app.IconBackground,
			DefaultLanguage:        acc.InterfaceLanguage,
			CustomizeTokenStrategy: "not_allow",
			CreatedAt:              &now,
			UpdatedAt:              &now,
			Code:                   (models.Site{}).GenerateCode(16),

			CreatedBy: app.CreatedBy,
			UpdatedBy: app.UpdatedBy,
		}
		dbengine.Instance().DB.Create(site)
	}
}

func UpdateAppDatasetJoinWhenAppPublishedWorkflowUpdated(app *models.App, wf *models.Workflow) {
	// app = sender
	// published_workflow = kwargs.get("published_workflow")
	// published_workflow = cast(Workflow, published_workflow)

	// dataset_ids = get_dataset_ids_from_workflow(published_workflow)
	// app_dataset_joins = db.session.query(AppDatasetJoin).filter(AppDatasetJoin.app_id == app.id).all()

	// removed_dataset_ids: set[str] = set()
	// if not app_dataset_joins:
	//     added_dataset_ids = dataset_ids
	// else:
	//     old_dataset_ids: set[str] = set()
	//     old_dataset_ids.update(app_dataset_join.dataset_id for app_dataset_join in app_dataset_joins)

	//     added_dataset_ids = dataset_ids - old_dataset_ids
	//     removed_dataset_ids = old_dataset_ids - dataset_ids

	// if removed_dataset_ids:
	//     for dataset_id in removed_dataset_ids:
	//         db.session.query(AppDatasetJoin).filter(
	//             AppDatasetJoin.app_id == app.id, AppDatasetJoin.dataset_id == dataset_id
	//         ).delete()

	// if added_dataset_ids:
	//     for dataset_id in added_dataset_ids:
	//         app_dataset_join = AppDatasetJoin(app_id=app.id, dataset_id=dataset_id)
	//         db.session.add(app_dataset_join)

	// db.session.commit()
}

func DeleteToolParametersCacheWhenSyncDraftWorkflow(app *models.App, synced_draft_workflow *models.Workflow) {
	if synced_draft_workflow == nil {
		return
	}
	nodes := []any{}
	graph_dict := synced_draft_workflow.GraphDict()
	if _, ok := graph_dict["nodes"]; ok {
		if _, ok := graph_dict["nodes"].([]any); ok {
			nodes = graph_dict["nodes"].([]any)
		} else if _, ok := graph_dict["nodes"].([]map[string]any); ok {
			for _, item := range graph_dict["nodes"].([]map[string]any) {
				nodes = append(nodes, item)
			}
		}
	}
	for _, node_data := range nodes {
		if _, ok := node_data.(map[string]any); !ok {
			continue
		}
		node_data_dict := node_data.(map[string]any)
		data := map[string]any{}
		if _, ok := node_data_dict["data"]; ok {
			if _, ok := node_data_dict["data"].(map[string]any); ok {
				data = node_data_dict["data"].(map[string]any)
			}
		}
		node_type := ""
		if _, ok := data["type"]; ok {
			if _, ok := data["type"].(string); ok {
				node_type = data["type"].(string)
			}
		}

		if node_type == string(nodesenumtypes.Node_TOOL) {
			func() {
				defer func() {
					if r := recover(); r != nil {
						mlog.Errorf("panic recover:%v\n%s", r, utils.GetCurrentGoroutineStack())
					}
				}()
				tool_entity := new(toolnodesentities.ToolEntity)
				bindata, _ := json.Marshal(data)
				err := json.Unmarshal(bindata, tool_entity)
				if err != nil {
					mlog.Errorf("json unmarshal(%#v) to ToolEntity failed:%v", string(bindata), err)
				} else {

				}
				// tool_entity = ToolEntity(**node_data["data"])
				// tool_runtime = ToolManager.get_tool_runtime(
				//     provider_type=tool_entity.provider_type,
				//     provider_id=tool_entity.provider_id,
				//     tool_name=tool_entity.tool_name,
				//     tenant_id=app.tenant_id,
				// )
				// manager = ToolParameterConfigurationManager(
				//     tenant_id=app.tenant_id,
				//     tool_runtime=tool_runtime,
				//     provider_name=tool_entity.provider_name,
				//     provider_type=tool_entity.provider_type,
				//     identity_id=f"WORKFLOW.{app.id}.{node_data.get('id')}",
				// )
				// manager.delete_tool_parameters_cache()
			}()
		}
	}
}
