package httprequest

import (
	"fmt"
	"iter"
	"strings"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/gofy/backend/core/workflow/nodes/base"
	variabletemplateparser "github.com/odysseythink/gofy/backend/core/workflow/utils/variable_template_parser"
	httprequestnodesentities "github.com/odysseythink/gofy/backend/entities/nodes/http_request"
	workflowentities "github.com/odysseythink/gofy/backend/entities/workflow"
	nodesenumtypes "github.com/odysseythink/gofy/backend/enum_types/nodes"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
)

type HttpRequestNode struct {
	*base.BaseNode[*httprequestnodesentities.HttpRequestNodeData]
}

func (n *HttpRequestNode) Type() nodesenumtypes.NodeType {
	return nodesenumtypes.Node_HTTP_REQUEST
}

func (n *HttpRequestNode) Run() (run_result *workflowentities.NodeRunResult, stream_run_result iter.Seq[any]) {
	process_data := map[string]any{}
	defer func() {
		if r := recover(); r != nil {
			if exp, ok := r.(*exceptions.HttpRequestNodeError); ok {
				mlog.Warningf("http request node %s failed to run: %v", n.GetNodeID(), exp)
				run_result = &workflowentities.NodeRunResult{
					Status:      models.WorkflowNodeExecutionStatus_FAILED,
					Error:       exp.Error(),
					ProcessData: process_data,
					ErrorType:   "HttpRequestNodeError",
				}
			} else {
				panic(r)
			}
		}
	}()
	http_executor := NewExecutor(
		n.NodeData,
		n.get_request_timeout(n.NodeData),
		n.GetGraphRuntimeState().VariablePool,
		0,
	)

	process_data["request"] = http_executor.to_log()

	response := http_executor.invoke()

	if n.ShouldContinueOnError() || n.ShouldRetry() {
		return &workflowentities.NodeRunResult{
			Status: models.WorkflowNodeExecutionStatus_FAILED,
			Outputs: map[string]any{
				"status_code": response.StatusCode,
				"body":        response.Text(),
				"headers":     response.Header,
				"files":       nil,
			},
			ProcessData: map[string]any{
				"request": process_data["request"],
			},
			Error:     fmt.Sprintf("Request failed with status code %d", response.StatusCode),
			ErrorType: "HTTPResponseCodeError",
		}, nil
	}
	headerdict := map[string]string{}
	defer response.Body.Close()
	for k, v := range map[string][]string(response.Header) {
		headerdict[k] = strings.Join(v, "; ")
	}
	return &workflowentities.NodeRunResult{
		Status: models.WorkflowNodeExecutionStatus_SUCCEEDED,
		Outputs: map[string]any{
			"status_code": response.StatusCode,
			"body":        response.Text(),
			"headers":     headerdict,
			"files":       nil,
		},
		ProcessData: map[string]any{
			"request": process_data["request"],
		},
	}, nil
	// except HttpRequestNodeError as e:
	// 	logger.warning(f"http request node {n.node_id} failed to run: {e}")
	// 	return NodeRunResult(
	// 		status=WorkflowNodeExecutionStatus.FAILED,
	// 		error=str(e),
	// 		process_data=process_data,
	// 		error_type=type(e).__name__,
	// 	)
}

func (n *HttpRequestNode) GetDefaultConfig(filters map[string]any) map[string]any {
	return map[string]any{
		"type": "http-request",
		"config": map[string]any{
			"method": "get",
			"authorization": map[string]any{
				"type": "no-auth",
			},
			"body": map[string]any{"type": "none"},
			"timeout": map[string]any{
				"connect":             confy.GetWithDefault[int]("http_node_config.max_connect_timeout", 300),
				"read":                confy.GetWithDefault[int]("http_node_config.max_read_timeout", 600),
				"write":               confy.GetWithDefault[int]("http_node_config.max_write_timeout", 600),
				"max_connect_timeout": confy.GetWithDefault[int]("http_node_config.max_connect_timeout", 300),
				"max_read_timeout":    confy.GetWithDefault[int]("http_node_config.max_read_timeout", 600),
				"max_write_timeout":   confy.GetWithDefault[int]("http_node_config.max_write_timeout", 600),
			},
		},
		"retry_config": map[string]any{
			"max_retries":    confy.GetWithDefault[int]("ssrf.default_max_retries", 3),
			"retry_interval": 0.5 * (2 * 2),
			"retry_enabled":  true,
		},
	}
}

func (n *HttpRequestNode) ExtractVariableSelectorToVariableMapping(graph_config map[string]any, node_id string, node_data *httprequestnodesentities.HttpRequestNodeData) map[string][]string {
	selectors := []*workflowentities.VariableSelector{}
	selectors = append(selectors, variabletemplateparser.ExtractSelectorsFromTemplate(node_data.URL)...)
	selectors = append(selectors, variabletemplateparser.ExtractSelectorsFromTemplate(node_data.Headers)...)
	selectors = append(selectors, variabletemplateparser.ExtractSelectorsFromTemplate(node_data.Params)...)
	if node_data.Body != nil {
		body_type := node_data.Body.Type
		data := node_data.Body.Data
		switch body_type {
		case "binary":
			if len(data) != 1 {
				panic(exceptions.NewRequestBodyError("invalid body data, should have only one item"))
			}
			selector := data[0].File
			selectors = append(selectors, &workflowentities.VariableSelector{
				Variable:      "#" + strings.Join(selector, ".") + "#",
				ValueSelector: selector,
			})
		case "json", "raw-text":
			if len(data) != 1 {
				panic(exceptions.NewRequestBodyError("invalid body data, should have only one item"))
			}
			selectors = append(selectors, variabletemplateparser.ExtractSelectorsFromTemplate(data[0].Key)...)
			selectors = append(selectors, variabletemplateparser.ExtractSelectorsFromTemplate(data[0].Value)...)
		case "x-www-form-urlencoded":
			for _, item := range data {
				selectors = append(selectors, variabletemplateparser.ExtractSelectorsFromTemplate(item.Key)...)
				selectors = append(selectors, variabletemplateparser.ExtractSelectorsFromTemplate(item.Value)...)
			}
		case "form-data":
			for _, item := range data {
				selectors = append(selectors, variabletemplateparser.ExtractSelectorsFromTemplate(item.Key)...)
				if item.Type == "text" {
					selectors = append(selectors, variabletemplateparser.ExtractSelectorsFromTemplate(item.Value)...)
				} else if item.Type == "file" {
					selectors = append(selectors, &workflowentities.VariableSelector{
						Variable:      "#" + strings.Join(item.File, ".") + "#",
						ValueSelector: item.File,
					})
				}
			}

		}
	}
	mapping := map[string][]string{}
	for _, selector_iter := range selectors {
		mapping[node_id+"."+selector_iter.Variable] = selector_iter.ValueSelector
	}
	return mapping
}

func (n *HttpRequestNode) get_request_timeout(node_data *httprequestnodesentities.HttpRequestNodeData) *httprequestnodesentities.HttpRequestNodeTimeout {
	timeout := node_data.Timeout
	if timeout == nil {
		return &httprequestnodesentities.HttpRequestNodeTimeout{
			Connect: confy.GetWithDefault[int]("http_node_config.max_connect_timeout", 300),
			Read:    confy.GetWithDefault[int]("http_node_config.max_read_timeout", 600),
			Write:   confy.GetWithDefault[int]("http_node_config.max_write_timeout", 600),
		}
	}
	if timeout.Connect <= 0 {
		timeout.Connect = confy.GetWithDefault[int]("http_node_config.max_connect_timeout", 300)
	}
	if timeout.Read <= 0 {
		timeout.Read = confy.GetWithDefault[int]("http_node_config.max_read_timeout", 600)
	}
	if timeout.Write <= 0 {
		timeout.Write = confy.GetWithDefault[int]("http_node_config.max_write_timeout", 600)
	}
	return timeout
}
func New() *HttpRequestNode {
	return &HttpRequestNode{
		BaseNode: &base.BaseNode[*httprequestnodesentities.HttpRequestNodeData]{},
	}
}
