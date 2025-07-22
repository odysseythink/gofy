package services

import (
	"reflect"

	dbengine "mlib.com/gofy/server/db_engine"
	nodesenumtypes "mlib.com/gofy/server/enum_types/nodes"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/validate"
	"mlib.com/mlog"
)

type WorkflowNodeExecutionService struct {
}

func (s *WorkflowNodeExecutionService) Get(id string) (wfr *models.WorkflowNodeExecution, err error) {
	wfr = &models.WorkflowNodeExecution{}
	err = dbengine.Instance().DB.Model(&models.WorkflowNodeExecution{}).Where("id = ?", id).Preload("Tenant").Preload("App").Preload("Workflow").Preload("WorkflowRun").First(wfr).Error
	if err != nil {
		wfr = nil
	}
	return
}

func (s *WorkflowNodeExecutionService) Extras(wfne *models.WorkflowNodeExecution) map[string]any {
	// from core.tools.tool_manager import ToolManager
	execution_metadata_dict := wfne.ExecutionMetadataDict()
	extras := map[string]any{}
	if execution_metadata_dict != nil {
		// from core.workflow.nodes import NodeType

		if wfne.NodeType == nodesenumtypes.Node_TOOL {
			if _, ok := execution_metadata_dict["tool_info"]; ok {
				// tool_info := execution_metadata_dict["tool_info"]
				if tool_info, ok := execution_metadata_dict["tool_info"].(map[string]any); ok && tool_info != nil {
					err := validate.StringMapTypeVerify(tool_info, validate.Rules{
						"provider_type": {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
						"provider_id":   {validate.RuleTypeOfField(reflect.String), validate.NotEmpty()},
					})
					if err != nil {
						mlog.Errorf("node type validate failed:%v", err)
					} else {
						extras["icon"], _ = ServiceGroupApp.Tool.GetToolIcon(wfne.TenantID, tool_info["provider_type"].(string), tool_info["provider_id"].(string))
					}
				}
			}
		}
	}
	return extras
}
