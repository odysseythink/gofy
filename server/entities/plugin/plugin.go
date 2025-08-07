package plugin

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"mlib.com/gofy/server/core/exceptions"
	agententities "mlib.com/gofy/server/entities/agent"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
	toolsentities "mlib.com/gofy/server/entities/tools"
	commontypes "mlib.com/gofy/server/types/common"
	"mlib.com/mlog"
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
)

type PluginDependencyType string

const (
	PluginDependency_Github      PluginDependencyType = PluginDependencyType(PluginInstallationSource_Github)
	PluginDependency_Marketplace PluginDependencyType = PluginDependencyType(PluginInstallationSource_Marketplace)
	PluginDependency_Package     PluginDependencyType = PluginDependencyType(PluginInstallationSource_Package)
)

type PluginResourceRequirements struct {
	Memory     int `json:"memory"`
	Permission *struct {
		Tool *struct {
			Enabled bool `json:"enabled"`
		} `json:"tool"`
		Model *struct {
			Enabled       bool `json:"enabled"`
			LLM           bool `json:"llm"`
			TextEmbedding bool `json:"text_embedding"`
			Rerank        bool `json:"rerank"`
			TTS           bool `json:"tts"`
			Speech2Text   bool `json:"speech2text"`
			Moderation    bool `json:"moderation"`
		} `json:"model"`
		Node *struct {
			Enabled bool `json:"enabled"`
		} `json:"node"`
		Endpoint *struct {
			Enabled bool `json:"enabled"`
		} `json:"endpoint"`
		Storage *struct {
			Enabled bool `json:"enabled"`
			Size    int  `json:"size"` //(ge=1024, le=1073741824, default=1048576)
		} `json:"storage"`
	} `json:"permission"`
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
		Tools     []string `json:"tools"`
		Models    []string `json:"models"`
		Endpoints []string `json:"endpoints"`
	} `json:"plugins"`
	Tags          []string                                   `json:"tags"`
	Repo          string                                     `json:"repo"`
	Verified      bool                                       `json:"verified"`
	Tool          *toolsentities.ToolProviderEntity          `json:"tool"`
	Model         *modelruntimeentities.ProviderEntity       `json:"model"`
	Endpoint      *EndpointProviderDeclaration               `json:"endpoint"`
	AgentStrategy *agententities.AgentStrategyProviderEntity `json:"agent_strategy"`
	Meta          struct {
		MinimumDifyVersion string `json:"minimum_dify_version"` // pattern=r"^\d{1,4}(\.\d{1,4}){1,3}(-\w{1,16})?$")
		Version            string `json:"version"`
	} `json:"meta"`

	// @model_validator(mode="before")
	// @classmethod
	// def validate_category(cls, values: dict) -> dict:
	//     # auto detect category
	//     if values.get("tool"):
	//         values["category"] = PluginCategory.Tool
	//     elif values.get("model"):
	//         values["category"] = PluginCategory.Model
	//     elif values.get("agent_strategy"):
	//         values["category"] = PluginCategory.AgentStrategy
	//     else:
	//         values["category"] = PluginCategory.Extension
	//     return values
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
		// check if matches [a-z0-9_-]+, if yes, append with langgenius/$value/$value
		matched, err = regexp.Match("^[a-z0-9_-]+$", []byte(value))
		if err != nil {
			mlog.Errorf("Invalid plugin id %v, reason:%v", value, err)
			panic(exceptions.NewValueError("Invalid plugin id " + value))
		}
		if matched {
			value = fmt.Sprintf("langgenius/%s/%s", value, value)
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
func (id *GenericProviderID) IsLanggenius() bool {
	return id.Organization == "langgenius"
}
func (id *GenericProviderID) PluginID() string {
	return fmt.Sprintf("%s/%s", id.Organization, id.PluginName)
}

type ModelProviderID struct {
	*GenericProviderID
}

func NewModelProviderID(value string, is_hardcoded bool) *ModelProviderID {
	mp := &ModelProviderID{
		GenericProviderID: NewGenericProviderID(value, is_hardcoded),
	}
	if mp.Organization == "langgenius" && mp.ProviderName == "google" {
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
	if tp.Organization == "langgenius" && slices.Contains([]string{"jina", "siliconflow", "stepfun", "gitee_ai"}, tp.ProviderName) {
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
