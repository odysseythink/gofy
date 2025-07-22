package app

type StreamEventType string

const (
	/*
	   Stream event
	*/

	StreamEvent_PING                     StreamEventType = "ping"
	StreamEvent_ERROR                    StreamEventType = "error"
	StreamEvent_MESSAGE                  StreamEventType = "message"
	StreamEvent_MESSAGE_END              StreamEventType = "message_end"
	StreamEvent_TTS_MESSAGE              StreamEventType = "tts_message"
	StreamEvent_TTS_MESSAGE_END          StreamEventType = "tts_message_end"
	StreamEvent_MESSAGE_FILE             StreamEventType = "message_file"
	StreamEvent_MESSAGE_REPLACE          StreamEventType = "message_replace"
	StreamEvent_AGENT_THOUGHT            StreamEventType = "agent_thought"
	StreamEvent_AGENT_MESSAGE            StreamEventType = "agent_message"
	StreamEvent_WORKFLOW_STARTED         StreamEventType = "workflow_started"
	StreamEvent_WORKFLOW_FINISHED        StreamEventType = "workflow_finished"
	StreamEvent_NODE_STARTED             StreamEventType = "node_started"
	StreamEvent_NODE_FINISHED            StreamEventType = "node_finished"
	StreamEvent_NODE_RETRY               StreamEventType = "node_retry"
	StreamEvent_PARALLEL_BRANCH_STARTED  StreamEventType = "parallel_branch_started"
	StreamEvent_PARALLEL_BRANCH_FINISHED StreamEventType = "parallel_branch_finished"
	StreamEvent_ITERATION_STARTED        StreamEventType = "iteration_started"
	StreamEvent_ITERATION_NEXT           StreamEventType = "iteration_next"
	StreamEvent_ITERATION_COMPLETED      StreamEventType = "iteration_completed"
	StreamEvent_TEXT_CHUNK               StreamEventType = "text_chunk"
	StreamEvent_TEXT_REPLACE             StreamEventType = "text_replace"
)
