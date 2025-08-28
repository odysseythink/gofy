package services

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	toolsentities "mlib.com/gofy/server/entities/tools"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type ToolsService struct {
}

func (s *ToolsService) GetToolIcon(tenant_id, provider_type, provider_id string) (any, error) {
	// """
	// get the tool icon

	// :param tenant_id: the id of the tenant
	// :param provider_type: the type of the provider
	// :param provider_id: the id of the provider
	// :return:
	// """
	// provider_type = provider_type
	// provider_id = provider_id
	// provider: Optional[Union[BuiltinToolProvider, ApiToolProvider, WorkflowToolProvider]] = None
	if provider_type == "builtin" {
		return viper.GetString("CONSOLE_API_URL") + "/console/api/workspaces/current/tool-provider/builtin/" + provider_id + "/icon", nil
	} else if provider_type == "api" {
		provider := new(models.ApiToolProvider)
		err := dbengine.Instance().DB.Model(&models.ApiToolProvider{}).Where("tenant_id = ? and id = ?", tenant_id, provider_id).First(provider).Error
		if err != nil {
			mlog.Errorf("get ApiToolProvider by tenant_id = %s and id = '%s' failed:%v", tenant_id, provider_id, err)
			return map[string]any{"background": "#252525", "content": `\ud83d\ude01`}, exceptions.NewToolProviderNotFoundError(fmt.Sprintf("api provider %s not found", provider_id))
		}
		icon_dic := map[string]any{}
		err = json.Unmarshal([]byte(provider.Icon), &icon_dic)
		if err != nil {
			return provider.Icon, nil
		}

		// return icon_dic, nil
		// 	icon = json.loads(provider.icon)
		// 	if isinstance(icon, (str, dict)):
		// 		return icon
		// 	return {"background": "#252525", "content": "\ud83d\ude01"}
		// except:
		// 	return {"background": "#252525", "content": "\ud83d\ude01"}
	} else if provider_type == "workflow" {
		provider := new(models.WorkflowToolProvider)
		err := dbengine.Instance().DB.Model(&models.WorkflowToolProvider{}).Where("tenant_id = ? and id = ?", tenant_id, provider_id).First(provider).Error
		if err != nil {
			mlog.Errorf("get WorkflowToolProvider by tenant_id = %s and id = '%s' failed:%v", tenant_id, provider_id, err)
			return map[string]any{"background": "#252525", "content": `\ud83d\ude01`}, exceptions.NewToolProviderNotFoundError(fmt.Sprintf("workflow provider %s not found", provider_id))
		}
		icon_dic := map[string]any{}
		err = json.Unmarshal([]byte(provider.Icon), &icon_dic)
		if err != nil {
			return provider.Icon, nil
		}
		// provider = (
		// 	db.session.query(WorkflowToolProvider)
		// 	.filter(WorkflowToolProvider.tenant_id == tenant_id, WorkflowToolProvider.id == provider_id)
		// 	.first()
		// )
		// if provider is None:
		// 	raise ToolProviderNotFoundError(f"workflow provider {provider_id} not found")

		// try:
		// 	icon = json.loads(provider.icon)
		// 	if isinstance(icon, (str, dict)):
		// 		return icon
		// 	return {"background": "#252525", "content": "\ud83d\ude01"}
		// except:
		// 	return {"background": "#252525", "content": "\ud83d\ude01"}
	}
	return nil, exceptions.NewValueError(fmt.Sprintf("provider type %s not found", provider_type))

}

func (service *ToolsService) ListToolLabels() []*toolsentities.ToolLabel {
	tmplist := []*toolsentities.ToolLabel{}
	for _, v := range toolsentities.DefaultToolLabelDict {
		tmplist = append(tmplist, &v)
	}
	return tmplist
}

func (service *ToolsService) ListToolProviders(user_id string, tenant_id string, typ string) []map[string]any {
	// providers = ToolManager.user_list_providers(user_id, tenant_id, typ)

	// // add icon
	// for _,provider := range providers{
	//     ToolTransformService.repack_provider(provider)
	// }
	// result = [provider.to_dict() for provider in providers]

	return nil
}
func (service *ToolsService) RetrieveMCPTools(tenant_id string, for_list bool) []*toolsentities.ToolProviderApiEntity {
	var datas []*models.MCPToolProvider
	err := dbengine.Instance().DB.Model(&models.MCPToolProvider{}).Where("tenant_id = ?", tenant_id).Find(&datas).Error
	if err != nil {
		mlog.Errorf("get MCPToolProvider by tenant_id = %s  failed:%v", tenant_id, err)
		return nil
	}
	var entities []*toolsentities.ToolProviderApiEntity
	for _, mcp_provider := range datas {
		data := ServiceGroupApp.ToolsTransform.MCPProviderToUserProvider(mcp_provider, for_list)
		if data != nil {
			if entities == nil {
				entities = make([]*toolsentities.ToolProviderApiEntity, 0)
			}
			entities = append(entities, data)
		}
	}
	return entities
}
