package plugin

// from typing import Optional

// from pydantic import BaseModel, Field, model_validator

// from core.model_runtime.entities.provider_entities import ProviderEntity
// from core.plugin.entities.endpoint import EndpointProviderDeclaration
// from core.plugin.entities.plugin import PluginResourceRequirements
// from core.tools.entities.common_entities import I18nObject
// from core.tools.entities.tool_entities import ToolProviderEntity
import (
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	toolsentities "mlib.com/gofy/server/entities/tools"
	parameterenumtypes "mlib.com/gofy/server/enum_types/parameter"
	commontypes "mlib.com/gofy/server/types/common"
)

type MarketplacePluginDeclaration[T1 parameterenumtypes.AppSelectorScopeType | parameterenumtypes.ModelSelectorScopeType | parameterenumtypes.ToolSelectorScopeType, T2 int | string] struct {
	Name                    string                                    `json:"name"`                      //description="Unique identifier for the plugin within the marketplace")
	Org                     string                                    `json:"org"`                       //description="Organization or developer responsible for creating and maintaining the plugin")
	PluginID                string                                    `json:"plugin_id"`                 //description="Globally unique identifier for the plugin across all marketplaces")
	Icon                    string                                    `json:"icon"`                      //description="URL or path to the plugin's visual representation")
	Label                   commontypes.I18nObject                    `json:"label"`                     //description="Localized display name for the plugin in different languages")
	Brief                   commontypes.I18nObject                    `json:"brief"`                     //description="Short, localized description of the plugin's functionality")
	Resource                PluginResourceRequirements                `json:"resource"`                  //escription="Specification of computational resources needed to run the plugin"
	Endpoint                *EndpointProviderDeclaration[T1, T2]      `json:"endpoint"`                  //description="Configuration for the plugin's API endpoint, if applicable"
	Model                   *modelruntimeentities.ProviderEntity      `json:"model"`                     //description="Details of the AI model used by the plugin, if any")
	Tool                    *toolsentities.ToolProviderEntity[T1, T2] `json:"tool"`                      //description="Information about the tool functionality provided by the plugin, if any"
	LatestVersion           string                                    `json:"latest_version"`            //"Most recent version number of the plugin available in the marketplace"
	LatestPackageIdentifier string                                    `json:"latest_package_identifier"` //"Unique identifier for the latest package release of the plugin"
	Status                  string                                    `json:"status"`                    //description="Indicate the status of marketplace plugin, enum from `active` `deleted`")
	DeprecatedReason        string                                    `json:"deprecated_reason"`         //"Not empty when status='deleted', indicates the reason why this plugin is deleted(deprecated)"
	AlternativePluginID     string                                    `json:"alternative_plugin_id"`     //"Optional, indicates the alternative plugin for user to switch to"
}

// @model_validator(mode="before")
// @classmethod
// def transform_declaration(cls, data: dict):
//     if "endpoint" in data and not data["endpoint"]:
//         del data["endpoint"]
//     if "model" in data and not data["model"]:
//         del data["model"]
//     if "tool" in data and not data["tool"]:
//         del data["tool"]
//     return data
