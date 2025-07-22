package app

type QueueEvent string

const (
	/*
	   QueueEvent enum
	*/

	QueueEvent_LLM_CHUNK                     QueueEvent = "llm_chunk"
	QueueEvent_TEXT_CHUNK                    QueueEvent = "text_chunk"
	QueueEvent_AGENT_MESSAGE                 QueueEvent = "agent_message"
	QueueEvent_MESSAGE_REPLACE               QueueEvent = "message_replace"
	QueueEvent_MESSAGE_END                   QueueEvent = "message_end"
	QueueEvent_ADVANCED_CHAT_MESSAGE_END     QueueEvent = "advanced_chat_message_end"
	QueueEvent_WORKFLOW_STARTED              QueueEvent = "workflow_started"
	QueueEvent_WORKFLOW_SUCCEEDED            QueueEvent = "workflow_succeeded"
	QueueEvent_WORKFLOW_FAILED               QueueEvent = "workflow_failed"
	QueueEvent_WORKFLOW_PARTIAL_SUCCEEDED    QueueEvent = "workflow_partial_succeeded"
	QueueEvent_ITERATION_START               QueueEvent = "iteration_start"
	QueueEvent_ITERATION_NEXT                QueueEvent = "iteration_next"
	QueueEvent_ITERATION_COMPLETED           QueueEvent = "iteration_completed"
	QueueEvent_NODE_STARTED                  QueueEvent = "node_started"
	QueueEvent_NODE_SUCCEEDED                QueueEvent = "node_succeeded"
	QueueEvent_NODE_FAILED                   QueueEvent = "node_failed"
	QueueEvent_NODE_EXCEPTION                QueueEvent = "node_exception"
	QueueEvent_RETRIEVER_RESOURCES           QueueEvent = "retriever_resources"
	QueueEvent_ANNOTATION_REPLY              QueueEvent = "annotation_reply"
	QueueEvent_AGENT_THOUGHT                 QueueEvent = "agent_thought"
	QueueEvent_MESSAGE_FILE                  QueueEvent = "message_file"
	QueueEvent_PARALLEL_BRANCH_RUN_STARTED   QueueEvent = "parallel_branch_run_started"
	QueueEvent_PARALLEL_BRANCH_RUN_SUCCEEDED QueueEvent = "parallel_branch_run_succeeded"
	QueueEvent_PARALLEL_BRANCH_RUN_FAILED    QueueEvent = "parallel_branch_run_failed"
	QueueEvent_ERROR                         QueueEvent = "error"
	QueueEvent_PING                          QueueEvent = "ping"
	QueueEvent_STOP                          QueueEvent = "stop"
	QueueEvent_RETRY                         QueueEvent = "retry"
)

// StopBy represents the reason for stopping.
type QueueStopEvent_StopBy string

const (
	QueueStopEvent_StopBy_USER_MANUAL       QueueStopEvent_StopBy = "user-manual"
	QueueStopEvent_StopBy_ANNOTATION_REPLY  QueueStopEvent_StopBy = "annotation-reply"
	QueueStopEvent_StopBy_OUTPUT_MODERATION QueueStopEvent_StopBy = "output-moderation"
	QueueStopEvent_StopBy_INPUT_MODERATION  QueueStopEvent_StopBy = "input-moderation"
)
