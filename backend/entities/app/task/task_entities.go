package task

import (
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
)

// TaskState represents the state of a task.
type TaskState struct {
	Metadata map[string]any `json:"metadata"`
}

// EasyUITaskState represents the state of an EasyUI task.
type EasyUITaskState struct {
	*TaskState
	LLMResult *modelruntimeentities.LLMResult `json:"llm_result"`
}

// WorkflowTaskState represents the state of a workflow task.
type WorkflowTaskState struct {
	*TaskState
	Answer string `json:"answer"`
}

func NewWorkflowTaskState() *WorkflowTaskState {
	return &WorkflowTaskState{
		TaskState: &TaskState{
			Metadata: make(map[string]any),
		},
	}
}

// Helper function to convert a struct to a map.
// func structToMap(obj any) map[string]any {
// 	result := make(map[string]any)
// 	value := reflect.Indirect(reflect.ValueOf(obj))
// 	typeOf := value.Type()

// 	for i := 0; i < value.NumField(); i++ {
// 		field := value.Field(i)
// 		tag := typeOf.Field(i).Tag.Get("json")
// 		if tag == "-" {
// 			continue
// 		}
// 		if tag == "" {
// 			continue
// 		}
// 		if idx := strings.Index(tag, ","); idx != -1 {
// 			tag = tag[:idx]
// 		}
// 		result[tag] = field.Interface()
// 	}

// 	return result
// }
