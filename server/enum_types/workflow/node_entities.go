package workflow

type NodeRunMetadataKey string

const (
	NodeRunMetadataKey_TOTAL_TOKENS                  NodeRunMetadataKey = "total_tokens"
	NodeRunMetadataKey_TOTAL_PRICE                   NodeRunMetadataKey = "total_price"
	NodeRunMetadataKey_CURRENCY                      NodeRunMetadataKey = "currency"
	NodeRunMetadataKey_TOOL_INFO                     NodeRunMetadataKey = "tool_info"
	NodeRunMetadataKey_ITERATION_ID                  NodeRunMetadataKey = "iteration_id"
	NodeRunMetadataKey_ITERATION_INDEX               NodeRunMetadataKey = "iteration_index"
	NodeRunMetadataKey_PARALLEL_ID                   NodeRunMetadataKey = "parallel_id"
	NodeRunMetadataKey_PARALLEL_START_NODE_ID        NodeRunMetadataKey = "parallel_start_node_id"
	NodeRunMetadataKey_PARENT_PARALLEL_ID            NodeRunMetadataKey = "parent_parallel_id"
	NodeRunMetadataKey_PARENT_PARALLEL_START_NODE_ID NodeRunMetadataKey = "parent_parallel_start_node_id"
	NodeRunMetadataKey_PARALLEL_MODE_RUN_ID          NodeRunMetadataKey = "parallel_mode_run_id"
	NodeRunMetadataKey_ITERATION_DURATION_MAP        NodeRunMetadataKey = "iteration_duration_map" // single iteration duration if iteration node runs
	NodeRunMetadataKey_ERROR_STRATEGY                NodeRunMetadataKey = "error_strategy"         // node in continue on error mode return the field
)
