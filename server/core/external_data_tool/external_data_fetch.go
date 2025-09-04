package externaldatatool

import (
	"sync"

	appconfigentities "mlib.com/gofy/server/entities/app/config"
	"mlib.com/mlog"
)

type ExternalDataFetch struct{}

// Fetch 并发拉取外部数据工具结果并回填到 inputs
func (e *ExternalDataFetch) Fetch(
	tenantID string,
	appID string,
	tools []*appconfigentities.ExternalDataVariableEntity,
	inputs map[string]any,
	query string,
) map[string]any {
	// 并发安全的结果收集
	var (
		mu     sync.Mutex
		result = make(map[string]string, len(tools))
		wg     sync.WaitGroup
	)

	for _, tool := range tools {
		wg.Add(1)
		go func(t *appconfigentities.ExternalDataVariableEntity) {
			defer wg.Done()

			varName, val := e.queryExternalDataTool(
				tenantID,
				appID,
				t,
				inputs,
				query,
			)
			if varName != "" && val != "" {
				mu.Lock()
				result[varName] = val
				mu.Unlock()
			}
		}(tool)
	}

	wg.Wait()

	// 回填到原始 map
	for k, v := range result {
		inputs[k] = v
	}
	return inputs
}

// queryExternalDataTool 单工具查询（等价 Python _query_external_data_tool）
func (e *ExternalDataFetch) queryExternalDataTool(
	tenantID string,
	appID string,
	tool *appconfigentities.ExternalDataVariableEntity,
	inputs map[string]any,
	query string,
) (string, string) {
	// 用工厂创建工具实例
	factory, err := NewExternalDataToolFactory(
		tool.Type,
		tenantID,
		appID,
		tool.Variable,
		tool.Config,
	)
	if err != nil {
		mlog.Errorf("create ExternalDataToolFactory failed:%v", err)
		return "", ""
	}

	// 假设 Query 返回 (value, error)
	val := factory.Query(inputs, query)
	return tool.Variable, val
}
