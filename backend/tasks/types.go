package tasks

// Task type constants for all background jobs.
const (
	// Document indexing
	TaskDocumentIndexing       = "document:indexing"
	TaskDocumentIndexingUpdate = "document:indexing:update"
	TaskDocumentIndexingSync   = "document:indexing:sync"
	TaskDocumentClean          = "document:clean"
	TaskDatasetClean           = "dataset:clean"

	// Segment operations
	TaskSegmentCreate      = "segment:create"
	TaskSegmentBatchCreate = "segment:batch_create"
	TaskSegmentDelete      = "segment:delete"
	TaskSegmentEnable      = "segment:enable"
	TaskSegmentDisable     = "segment:disable"

	// Workflow
	TaskWorkflowExecution     = "workflow:execution"
	TaskWorkflowNodeExecution = "workflow:node_execution"
	TaskWorkflowSchedule      = "workflow:schedule"

	// Ops
	TaskOpsTrace = "ops:trace"

	// Cleanup (scheduled)
	TaskCleanMessages       = "clean:messages"
	TaskCleanWorkflowRuns   = "clean:workflow_runs"
	TaskCleanEmbeddingCache = "clean:embedding_cache"
	TaskCleanUnusedDatasets = "clean:unused_datasets"
)
