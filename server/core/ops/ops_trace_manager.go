package ops

import (
	"encoding/json"
	"fmt"

	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	opsentities "mlib.com/gofy/server/entities/ops"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

var (
	provider_config_map = map[opsentities.TracingProviderType]map[string]any{
		opsentities.TracingProvider_LANGFUSE: map[string]any{
			"config_class": &opsentities.LangfuseConfig{},
			"secret_keys":  []string{"public_key", "secret_key"},
			"other_keys":   []string{"host", "project_key"},
			// "trace_instance": &opsentities.LangFuseDataTrace{},
		},
		opsentities.TracingProvider_LANGSMITH: map[string]any{
			"config_class": &opsentities.LangSmithConfig{},
			"secret_keys":  []string{"api_key"},
			"other_keys":   []string{"project", "endpoint"},
			// "trace_instance": &opsentities.LangSmithDataTrace{},
		},
		opsentities.TracingProvider_OPIK: map[string]any{
			"config_class": &opsentities.OpikConfig{},
			"secret_keys":  []string{"api_key"},
			"other_keys":   []string{"project", "url", "workspace"},
			// "trace_instance": &opsentities.OpikDataTrace{},
		},
	}
)

type OpsTraceManager struct{}

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
