package ops

import (
	"fmt"
	"strings"
)

type TracingProviderType string

const (
	TracingProvider_LANGFUSE  TracingProviderType = "langfuse"
	TracingProvider_LANGSMITH TracingProviderType = "langsmith"
	TracingProvider_OPIK      TracingProviderType = "opik"
)

type BaseTracingConfig struct {
	// 基础追踪配置字段可以在这里定义
}

type LangfuseConfig struct {
	BaseTracingConfig
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key"`
	Host      string `json:"host"`
}

func (l *LangfuseConfig) SetValue(host string) error {
	if host == "" {
		host = "https://api.langfuse.com"
	}
	if !(strings.HasPrefix(host, "https://") || strings.HasPrefix(host, "http://")) {
		return fmt.Errorf("host must start with https:// or http://")
	}
	l.Host = host
	return nil
}

type LangSmithConfig struct {
	BaseTracingConfig
	APIKey   string `json:"api_key"`
	Project  string `json:"project"`
	Endpoint string `json:"endpoint"`
}

func (l *LangSmithConfig) SetValue(endpoint string) error {
	if endpoint == "" {
		endpoint = "https://api.smith.langchain.com"
	}
	if !strings.HasPrefix(endpoint, "https://") {
		return fmt.Errorf("endpoint must start with https://")
	}
	l.Endpoint = endpoint
	return nil
}

type OpikConfig struct {
	BaseTracingConfig
	APIKey    *string `json:"api_key,omitempty"`
	Project   *string `json:"project,omitempty"`
	Workspace *string `json:"workspace,omitempty"`
	URL       string  `json:"url"`
}

func (o *OpikConfig) ProjectValidator(project string) {
	if project == "" {
		project = "Default Project"
	}
	o.Project = &project
}

func (o *OpikConfig) URLValidator(url string) error {
	if url == "" {
		url = "https://www.comet.com/opik/api/"
	}
	if !(strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://")) {
		return fmt.Errorf("url must start with https:// or http://")
	}
	if !strings.HasSuffix(url, "/api/") {
		return fmt.Errorf("url should ends with /api/")
	}
	o.URL = url
	return nil
}

const (
	OPS_FILE_PATH        = "ops_trace/"
	OPS_TRACE_FAILED_KEY = "FAILED_OPS_TRACE"
)
