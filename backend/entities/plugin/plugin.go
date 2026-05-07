package plugin

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	agententities "github.com/odysseythink/gofy/backend/entities/agent"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	toolsentities "github.com/odysseythink/gofy/backend/entities/tools"
	commontypes "github.com/odysseythink/gofy/backend/types/common"
	"github.com/odysseythink/mlog"
)

type PluginInstallationSourceType string

const (
	PluginInstallationSource_Github      PluginInstallationSourceType = "github"
	PluginInstallationSource_Marketplace PluginInstallationSourceType = "marketplace"
	PluginInstallationSource_Package     PluginInstallationSourceType = "package"
	PluginInstallationSource_Remote      PluginInstallationSourceType = "remote"
)

type PluginCategoryType string

const (
	PluginCategory_Tool          PluginCategoryType = "tool"
	PluginCategory_Model         PluginCategoryType = "model"
	PluginCategory_Extension     PluginCategoryType = "extension"
	PluginCategory_AgentStrategy PluginCategoryType = "agent-strategy"
	PluginCategory_Datasource    PluginCategoryType = "datasource"
	PluginCategory_Trigger       PluginCategoryType = "trigger"
)

type PluginDependencyType string

const (
	PluginDependency_Github      PluginDependencyType = PluginDependencyType(PluginInstallationSource_Github)
	PluginDependency_Marketplace PluginDependencyType = PluginDependencyType(PluginInstallationSource_Marketplace)
	PluginDependency_Package     PluginDependencyType = PluginDependencyType(PluginInstallationSource_Package)
)

type PluginPermissionToolRequirement struct {
	Enabled bool `json:"enabled"`
}

type PluginPermissionModelRequirement struct {
	Enabled       bool `json:"enabled"`
	LLM           bool `json:"llm"`
	TextEmbedding bool `json:"text_embedding"`
	Rerank        bool `json:"rerank"`
	TTS           bool `json:"tts"`
	Speech2Text   bool `json:"speech2text"`
	Moderation    bool `json:"moderation"`
}

type PluginPermissionNodeRequirement struct {
	Enabled bool `json:"enabled"`
}

type PluginPermissionEndpointRequirement struct {
	Enabled bool `json:"enabled"`
}

type PluginPermissionAppRequirement struct {
	Enabled bool `json:"enabled"`
}

type PluginPermissionStorageRequirement struct {
	Enabled bool `json:"enabled"`
	Size    int  `json:"size"` //(ge=1024, le=1073741824, default=1048576)
}

type PluginPermissionRequirement struct {
	Tool     *PluginPermissionToolRequirement     `json:"tool,omitempty"`
	Model    *PluginPermissionModelRequirement    `json:"model,omitempty"`
	Node     *PluginPermissionNodeRequirement     `json:"node,omitempty"`
	Endpoint *PluginPermissionEndpointRequirement `json:"endpoint,omitempty"`
	App      *PluginPermissionAppRequirement      `json:"app,omitempty"`
	Storage  *PluginPermissionStorageRequirement  `json:"storage,omitempty"`
}

func (p *PluginPermissionRequirement) AllowInvokeTool() bool {
	return p != nil && p.Tool != nil && p.Tool.Enabled
}

func (p *PluginPermissionRequirement) AllowInvokeModel() bool {
	return p != nil && p.Model != nil && p.Model.Enabled
}

func (p *PluginPermissionRequirement) AllowInvokeLLM() bool {
	return p != nil && p.Model != nil && p.Model.Enabled && p.Model.LLM
}

func (p *PluginPermissionRequirement) AllowInvokeTextEmbedding() bool {
	return p != nil && p.Model != nil && p.Model.Enabled && p.Model.TextEmbedding
}

func (p *PluginPermissionRequirement) AllowInvokeRerank() bool {
	return p != nil && p.Model != nil && p.Model.Enabled && p.Model.Rerank
}

func (p *PluginPermissionRequirement) AllowInvokeTTS() bool {
	return p != nil && p.Model != nil && p.Model.Enabled && p.Model.TTS
}

func (p *PluginPermissionRequirement) AllowInvokeSpeech2Text() bool {
	return p != nil && p.Model != nil && p.Model.Enabled && p.Model.Speech2Text
}

func (p *PluginPermissionRequirement) AllowInvokeModeration() bool {
	return p != nil && p.Model != nil && p.Model.Enabled && p.Model.Moderation
}

func (p *PluginPermissionRequirement) AllowInvokeNode() bool {
	return p != nil && p.Node != nil && p.Node.Enabled
}

func (p *PluginPermissionRequirement) AllowInvokeApp() bool {
	return p != nil && p.App != nil && p.App.Enabled
}

func (p *PluginPermissionRequirement) AllowRegisterEndpoint() bool {
	return p != nil && p.Endpoint != nil && p.Endpoint.Enabled
}

func (p *PluginPermissionRequirement) AllowInvokeStorage() bool {
	return p != nil && p.Storage != nil && p.Storage.Enabled
}

type PluginResourceRequirements struct {
	Memory     int                          `json:"memory"`
	Permission *PluginPermissionRequirement `json:"permission,omitempty"`
}

type PluginDeclaration struct {
	Version     string                     `json:"version"` //pattern=r"^\d{1,4}(\.\d{1,4}){1,3}(-\w{1,16})?$")
	Author      string                     `json:"author"`  //pattern=r"^[a-zA-Z0-9_-]{1,64}$")
	Name        string                     `json:"name"`    //pattern=r"^[a-z0-9_-]{1,128}$")
	Description commontypes.I18nObject     `json:"description"`
	Icon        string                     `json:"icon"`
	IconDark    string                     `json:"icon_dark"`
	Label       commontypes.I18nObject     `json:"label"`
	Category    PluginCategoryType         `json:"category"`
	CreatedAt   time.Time                  `json:"created_at"`
	Resource    PluginResourceRequirements `json:"resource"`
	Plugins     struct {
		Tools           []string `json:"tools"`
		Models          []string `json:"models"`
		Endpoints       []string `json:"endpoints"`
		AgentStrategies []string `json:"agent_strategies"`
		Datasources     []string `json:"datasources"`
		Triggers        []string `json:"triggers"`
	} `json:"plugins"`
	Tags     []string `json:"tags"`
	Repo     string   `json:"repo"`
	Privacy  string   `json:"privacy,omitempty"`
	Verified bool     `json:"verified"`

	Tool          *toolsentities.ToolProviderEntity          `json:"tool,omitempty"`
	Model         *modelruntimeentities.ProviderEntity       `json:"model,omitempty"`
	Endpoint      *EndpointProviderDeclaration               `json:"endpoint,omitempty"`
	AgentStrategy *agententities.AgentStrategyProviderEntity `json:"agent_strategy,omitempty"`
	// Datasource and Trigger are placeholders for future provider declarations.
	Datasource any `json:"datasource,omitempty"`
	Trigger    any `json:"trigger,omitempty"`

	Meta struct {
		MinimumGofyVersion string `json:"minimum_gofy_version"` // pattern=r"^\d{1,4}(\.\d{1,4}){1,3}(-\w{1,16})?$")
		Version            string `json:"version"`
	} `json:"meta"`
}

// DetectCategory auto-detects the plugin category based on provider declarations.
func (d *PluginDeclaration) DetectCategory() PluginCategoryType {
	if d.Tool != nil || len(d.Plugins.Tools) > 0 {
		return PluginCategory_Tool
	}
	if d.Model != nil || len(d.Plugins.Models) > 0 {
		return PluginCategory_Model
	}
	if d.Datasource != nil || len(d.Plugins.Datasources) > 0 {
		return PluginCategory_Datasource
	}
	if d.AgentStrategy != nil || len(d.Plugins.AgentStrategies) > 0 {
		return PluginCategory_AgentStrategy
	}
	if d.Trigger != nil || len(d.Plugins.Triggers) > 0 {
		return PluginCategory_Trigger
	}
	return PluginCategory_Extension
}

type PluginInstallation struct {
	BasePluginEntity
	TenantID               string                       `json:"tenant_id"`
	EndpointsSetups        int                          `json:"endpoints_setups"`
	EndpointsActive        int                          `json:"endpoints_active"`
	RuntimeType            string                       `json:"runtime_type"`
	Source                 PluginInstallationSourceType `json:"source"`
	Meta                   map[string]any               `json:"meta"`
	PluginID               string                       `json:"plugin_id"`
	PluginUniqueIdentifier string                       `json:"plugin_unique_identifier"`
	Version                string                       `json:"version"`
	Checksum               string                       `json:"checksum"`
	Declaration            PluginDeclaration            `json:"declaration"`
}

type PluginEntity struct {
	*PluginInstallation
	Name           string `json:"name"`
	InstallationID string `json:"installation_id"`
	Version        string `json:"version"`

	// @model_validator(mode="after")
	// def set_plugin_id(self):
	//     if self.declaration.tool:
	//         self.declaration.tool.plugin_id = self.plugin_id
	//     return self
}

type GenericProviderID struct {
	Organization string `json:"organization"`
	PluginName   string `json:"plugin_name"`
	ProviderName string `json:"provider_name"`
	IsHardcoded  bool   `json:"is_hardcoded"`
}

func (id *GenericProviderID) ToString() string {
	return fmt.Sprintf("%s/%s/%s", id.Organization, id.PluginName, id.ProviderName)
}

func NewGenericProviderID(value string, is_hardcoded bool) *GenericProviderID {
	if value == "" {
		panic(exceptions.NewValueError("plugin not found, please add plugin"))
	}
	// check if the value is a valid plugin id with format: $organization/$plugin_name/$provider_name
	matched, err := regexp.Match("^[a-z0-9_-]+\\/[a-z0-9_-]+\\/[a-z0-9_-]+$", []byte(value))
	if err != nil {
		mlog.Errorf("Invalid plugin id %v, reason:%v", value, err)
		panic(exceptions.NewValueError("Invalid plugin id " + value))
	}
	if !matched {
		// check if matches [a-z0-9_-]+, if yes, append with odysseythink/$value/$value
		matched, err = regexp.Match("^[a-z0-9_-]+$", []byte(value))
		if err != nil {
			mlog.Errorf("Invalid plugin id %v, reason:%v", value, err)
			panic(exceptions.NewValueError("Invalid plugin id " + value))
		}
		if matched {
			value = fmt.Sprintf("odysseythink/%s/%s", value, value)
		} else {
			mlog.Errorf("Invalid plugin id %v, reason:%v", value, err)
			panic(exceptions.NewValueError("Invalid plugin id " + value))
		}
	}
	tmps := strings.Split(value, "/")
	return &GenericProviderID{
		Organization: tmps[0],
		PluginName:   tmps[1],
		ProviderName: tmps[2],
		IsHardcoded:  is_hardcoded,
	}
}
func (id *GenericProviderID) IsOdysseythink() bool {
	return id.Organization == "odysseythink"
}
func (id *GenericProviderID) PluginID() string {
	return fmt.Sprintf("%s/%s", id.Organization, id.PluginName)
}

func (id *GenericProviderID) String() string {
	return fmt.Sprintf("%s/%s/%s", id.Organization, id.PluginName, id.PluginName)
}

type ModelProviderID struct {
	*GenericProviderID
}

func NewModelProviderID(value string, is_hardcoded bool) *ModelProviderID {
	mp := &ModelProviderID{
		GenericProviderID: NewGenericProviderID(value, is_hardcoded),
	}
	if mp.Organization == "odysseythink" && mp.ProviderName == "google" {
		mp.PluginName = "gemini"
	}
	return mp
}

type ToolProviderID struct {
	*GenericProviderID
}

func NewToolProviderID(value string, is_hardcoded bool) *ToolProviderID {
	tp := &ToolProviderID{
		GenericProviderID: NewGenericProviderID(value, is_hardcoded),
	}
	if tp.Organization == "odysseythink" && slices.Contains([]string{"jina", "siliconflow", "stepfun", "gitee_ai"}, tp.ProviderName) {
		tp.PluginName = fmt.Sprintf("%s_tool", tp.ProviderName)
	}
	return tp
}

type PluginDependencyGithub struct {
	Repo                         string  `json:"repo"`
	Version                      string  `json:"version"`
	Package                      string  `json:"package"`
	GithubPluginUniqueIdentifier string  `json:"github_plugin_unique_identifier"`
	CurrentIdentifier            *string `json:"current_identifier"`
}

func (dep *PluginDependencyGithub) PluginUniqueIdentifier() string {
	return dep.GithubPluginUniqueIdentifier
}

type PluginDependencyMarketplace struct {
	MarketplacePluginUniqueIdentifier string  `json:"marketplace_plugin_unique_identifier"`
	CurrentIdentifier                 *string `json:"current_identifier"`
}

func (dep *PluginDependencyMarketplace) PluginUniqueIdentifier() string {
	return dep.MarketplacePluginUniqueIdentifier
}

type PluginDependencyPackage struct {
	PluginUniqueIdentifier string  `json:"plugin_unique_identifier"`
	CurrentIdentifier      *string `json:"current_identifier"`
}
type MissingPluginDependency struct {
	PluginUniqueIdentifier string  `json:"plugin_unique_identifier"`
	CurrentIdentifier      *string `json:"current_identifier"`
}
