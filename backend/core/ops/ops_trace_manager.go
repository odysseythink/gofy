package ops

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	opsentities "github.com/odysseythink/gofy/backend/entities/ops"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

var (
	provider_config_map = map[opsentities.TracingProviderType]map[string]any{
		opsentities.TracingProvider_LANGFUSE: {
			"config_class": &opsentities.LangfuseConfig{},
			"secret_keys":  []string{"public_key", "secret_key"},
			"other_keys":   []string{"host", "project_key"},
		},
		opsentities.TracingProvider_LANGSMITH: {
			"config_class": &opsentities.LangSmithConfig{},
			"secret_keys":  []string{"api_key"},
			"other_keys":   []string{"project", "endpoint"},
		},
		opsentities.TracingProvider_OPIK: {
			"config_class": &opsentities.OpikConfig{},
			"secret_keys":  []string{"api_key"},
			"other_keys":   []string{"project", "url", "workspace"},
		},
	}
)

// OpsTraceManager manages tracing configuration and delegates to trace backends.
type OpsTraceManager struct {
	backend TraceInstance
}

// InitTraceBackend initialises the active trace backend based on the provider
// configuration stored for the given app. If tracing is not enabled or no
// provider is configured, the manager operates with a nil backend (no-op).
func (mgr *OpsTraceManager) InitTraceBackend(appID string) error {
	cfg := mgr.GetAppTracingConfig(appID)

	enabled, _ := cfg["enabled"].(bool)
	if !enabled {
		mgr.backend = nil
		return nil
	}

	provider, _ := cfg["tracing_provider"].(string)
	switch opsentities.TracingProviderType(provider) {
	case opsentities.TracingProvider_LANGFUSE:
		publicKey, _ := cfg["public_key"].(string)
		secretKey, _ := cfg["secret_key"].(string)
		host, _ := cfg["host"].(string)
		mgr.backend = NewLangfuseInstance(LangfuseTraceConfig{
			PublicKey: publicKey,
			SecretKey: secretKey,
			Host:      host,
		})
	case opsentities.TracingProvider_LANGSMITH:
		apiKey, _ := cfg["api_key"].(string)
		project, _ := cfg["project"].(string)
		endpoint, _ := cfg["endpoint"].(string)
		mgr.backend = NewLangSmithInstance(LangSmithTraceConfig{
			APIKey:   apiKey,
			Project:  project,
			Endpoint: endpoint,
		})
	default:
		mlog.Warningf("ops: unsupported tracing provider %q, tracing disabled", provider)
		mgr.backend = nil
	}
	return nil
}

// RecordTrace delegates a span to the active trace backend. If no backend is
// configured the call is a no-op.
func (mgr *OpsTraceManager) RecordTrace(ctx context.Context, span *TraceSpan) error {
	if mgr.backend == nil {
		return nil
	}
	return mgr.backend.Trace(ctx, span)
}

// FlushTraces flushes any buffered spans in the active backend.
func (mgr *OpsTraceManager) FlushTraces(ctx context.Context) error {
	if mgr.backend == nil {
		return nil
	}
	return mgr.backend.Flush(ctx)
}

// CloseTraceBackend shuts down the active trace backend.
func (mgr *OpsTraceManager) CloseTraceBackend() error {
	if mgr.backend == nil {
		return nil
	}
	return mgr.backend.Close()
}

// ActiveProvider returns the provider type of the active backend, or "" if none.
func (mgr *OpsTraceManager) ActiveProvider() TraceProvider {
	if mgr.backend == nil {
		return ""
	}
	return mgr.backend.Provider()
}

func (mgr *OpsTraceManager) GetAppTracingConfig(app_id string) map[string]any {
	app := new(models.App)
	err := dbengine.Instance().DB.Model(&models.App{}).Where("id = ?", app_id).First(app).Error
	if err != nil {
		mlog.Errorf("get app failed:%v", err)
		app = nil
	}
	if app == nil {
		panic(exceptions.NewValueError("App not found"))
	}
	if app.Tracing == "" {
		return map[string]any{"enabled": false, "tracing_provider": ""}
	}

	app_trace_config := map[string]any{}
	err = json.Unmarshal([]byte(app.Tracing), &app_trace_config)
	if err != nil {
		mlog.Errorf("json unmarshal app trace(%s) failed:%v", app.Tracing, err)
		panic(exceptions.NewValueError("App tracing is invalid"))
	}
	return app_trace_config
}

func (mgr *OpsTraceManager) UpdateAppTracingConfig(app_id string, enabled bool, tracing_provider string) {
	// auth check
	if tracing_provider != "" {
		if _, ok := provider_config_map[opsentities.TracingProviderType(tracing_provider)]; !ok {
			panic(exceptions.NewValueError(fmt.Sprintf("Invalid tracing provider: %s", tracing_provider)))
		}
	}
	app := new(models.App)
	err := dbengine.Instance().DB.Model(&models.App{}).Where("id = ?", app_id).First(app).Error
	if err != nil {
		mlog.Errorf("get app failed:%v", err)
		app = nil
	}
	if app == nil {
		panic(exceptions.NewValueError("App not found"))
	}
	bindata, _ := json.Marshal(map[string]any{
		"enabled":          enabled,
		"tracing_provider": tracing_provider,
	})
	dbengine.Instance().DB.Updates(&models.App{ID: app_id, Tracing: string(bindata)})
}
