package app

type QueueEventType string

const (
	/*
	   QueueEvent enum
	*/
	QueueEvent_LLM_CHUNK                     QueueEventType = "llm_chunk"
	QueueEvent_TEXT_CHUNK                    QueueEventType = "text_chunk"
	QueueEvent_AGENT_MESSAGE                 QueueEventType = "agent_message"
	QueueEvent_MESSAGE_REPLACE               QueueEventType = "message_replace"
	QueueEvent_MESSAGE_END                   QueueEventType = "message_end"
	QueueEvent_ADVANCED_CHAT_MESSAGE_END     QueueEventType = "advanced_chat_message_end"
	QueueEvent_WORKFLOW_STARTED              QueueEventType = "workflow_started"
	QueueEvent_WORKFLOW_SUCCEEDED            QueueEventType = "workflow_succeeded"
	QueueEvent_WORKFLOW_FAILED               QueueEventType = "workflow_failed"
	QueueEvent_WORKFLOW_PARTIAL_SUCCEEDED    QueueEventType = "workflow_partial_succeeded"
	QueueEvent_ITERATION_START               QueueEventType = "iteration_start"
	QueueEvent_ITERATION_NEXT                QueueEventType = "iteration_next"
	QueueEvent_ITERATION_COMPLETED           QueueEventType = "iteration_completed"
	QueueEvent_LOOP_START                    QueueEventType = "loop_start"
	QueueEvent_LOOP_NEXT                     QueueEventType = "loop_next"
	QueueEvent_LOOP_COMPLETED                QueueEventType = "loop_completed"
	QueueEvent_NODE_STARTED                  QueueEventType = "node_started"
	QueueEvent_NODE_SUCCEEDED                QueueEventType = "node_succeeded"
	QueueEvent_NODE_FAILED                   QueueEventType = "node_failed"
	QueueEvent_NODE_EXCEPTION                QueueEventType = "node_exception"
	QueueEvent_RETRIEVER_RESOURCES           QueueEventType = "retriever_resources"
	QueueEvent_ANNOTATION_REPLY              QueueEventType = "annotation_reply"
	QueueEvent_AGENT_THOUGHT                 QueueEventType = "agent_thought"
	QueueEvent_MESSAGE_FILE                  QueueEventType = "message_file"
	QueueEvent_PARALLEL_BRANCH_RUN_STARTED   QueueEventType = "parallel_branch_run_started"
	QueueEvent_PARALLEL_BRANCH_RUN_SUCCEEDED QueueEventType = "parallel_branch_run_succeeded"
	QueueEvent_PARALLEL_BRANCH_RUN_FAILED    QueueEventType = "parallel_branch_run_failed"
	QueueEvent_AGENT_LOG                     QueueEventType = "agent_log"
	QueueEvent_ERROR                         QueueEventType = "error"
	QueueEvent_PING                          QueueEventType = "ping"
	QueueEvent_STOP                          QueueEventType = "stop"
	QueueEvent_RETRY                         QueueEventType = "retry"
)

// StopBy represents the reason for stopping.
type QueueStopEvent_StopBy string

const (
	QueueStopEvent_StopBy_USER_MANUAL       QueueStopEvent_StopBy = "user-manual"
	QueueStopEvent_StopBy_ANNOTATION_REPLY  QueueStopEvent_StopBy = "annotation-reply"
	QueueStopEvent_StopBy_OUTPUT_MODERATION QueueStopEvent_StopBy = "output-moderation"
	QueueStopEvent_StopBy_INPUT_MODERATION  QueueStopEvent_StopBy = "input-moderation"
)
