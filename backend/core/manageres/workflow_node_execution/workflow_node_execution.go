package workflownodeexecution

import (
	toolmanager "github.com/odysseythink/gofy/backend/core/manageres/tool_manager"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type WorkflowNodeExecutionManager struct {
}

func (mgr *WorkflowNodeExecutionManager) CreatedByAccount(wfne *models.WorkflowNodeExecution) *models.Account {
	if wfne.CreatedByRole == models.CreatedByRole_ACCOUNT {
		acc := new(models.Account)
		err := dbengine.Instance().DB.Model(&models.Account{}).Where("id = ?", wfne.CreatedBy).First(acc).Error
		if err != nil {
			mlog.Errorf("get account(%s) failed:%v", wfne.CreatedBy, err)
			return nil
		}
		return acc
	} else {
		return nil
	}

}
func (mgr *WorkflowNodeExecutionManager) CreatedByEndUser(wfne *models.WorkflowNodeExecution) *models.EndUser {
	if wfne.CreatedByRole == models.CreatedByRole_END_USER {
		user := new(models.EndUser)
		err := dbengine.Instance().DB.Model(&models.EndUser{}).Where("id = ?", wfne.CreatedBy).First(user).Error
		if err != nil {
			mlog.Errorf("get EndUser(%s) failed:%v", wfne.CreatedBy, err)
			return nil
		}
		return user
	} else {
		return nil
	}
}

func (mgr *WorkflowNodeExecutionManager) Extras(wfne *models.WorkflowNodeExecution) map[string]any {

	extras := map[string]any{}
	metadata_dict := wfne.ExecutionMetadataDict()
	if len(metadata_dict) > 0 {
		if _, ok := metadata_dict["tool_info"]; ok && wfne.NodeType == nodesenumtypes.Node_TOOL {
			if _, ok := metadata_dict["tool_info"].(map[string]any); ok {
				tool_info := metadata_dict["tool_info"].(map[string]any)
				provider_type := ""
				if _, ok := tool_info["provider_type"]; ok {
					if _, ok := tool_info["provider_type"].(string); ok {
						provider_type = tool_info["provider_type"].(string)
					}
				}
				provider_id := ""
				if _, ok := tool_info["provider_id"]; ok {
					if _, ok := tool_info["provider_id"].(string); ok {
						provider_id = tool_info["provider_id"].(string)
					}
				}
				if provider_type != "" && provider_id != "" {
					extras["icon"] = (&toolmanager.ToolManager{}).GetToolIcon(
						wfne.TenantID,
						provider_type,
						provider_id,
					)
				}
			}
		}
	}
	return extras
}
