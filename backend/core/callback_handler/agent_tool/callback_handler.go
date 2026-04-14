package agenttool

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/odysseythink/confy"
)

var (
	_TEXT_COLOR_MAPPING = map[string]string{
		"blue":   "36;1",
		"yellow": "33;1",
		"pink":   "38;5;200",
		"green":  "32;1",
		"red":    "31;1",
	}
)

func get_colored_text(text string, color string) string {
	/*Get colored text.*/
	color_str := _TEXT_COLOR_MAPPING[color]
	return fmt.Sprintf("\u001b[%sm\033[1;3m%s\u001b[0m", color_str, text)
}

func print_text(text string, color string) {
	/*Print text with highlighting and no end characters.*/
	text_to_print := text
	if color != "" {
		text_to_print = get_colored_text(text, color)
	}
	fmt.Println(text_to_print)
}

type AgentCallbackHandler struct {
	/*Callback Handler that prints to std out.*/

	color        string
	current_loop int
}

func New(color string) *AgentCallbackHandler {
	// use a specific color is not specified
	if color == "" {
		color = "green"
	}
	return &AgentCallbackHandler{
		color:        color,
		current_loop: 1,
	}
}
func (handler *AgentCallbackHandler) OnToolStart(

	tool_name string,
	tool_inputs map[string]any,
) {
	/*Do nothing.*/
	if confy.GetWithDefault[bool]("DEBUG", false) {
		bindata, _ := json.Marshal(tool_inputs)
		print_text("\n[OnToolStart] ToolCall:"+tool_name+"\n"+string(bindata)+"\n", handler.color)
	}
}
func (handler *AgentCallbackHandler) OnToolEnd(

	tool_name string,
	tool_inputs map[string]any,
	tool_outputs any, //Iterable[ToolInvokeMessage] | str,
	message_id string,
	// timer: Optional[Any] = None,
	// trace_manager: Optional[TraceQueueManager] = None,
) {
	/*If not the final action, print out observation.*/
	if confy.GetWithDefault[bool]("DEBUG", false) {
		bindata, _ := json.Marshal(tool_inputs)
		print_text("\n[OnToolEnd]\n", handler.color)
		print_text("Tool: "+tool_name+"\n", handler.color)
		print_text("Inputs: "+string(bindata)+"\n", handler.color)
		print_text("Outputs: "+fmt.Sprintf("%v", tool_outputs)+"\n", handler.color)
		print_text("\n", "")
	}
	// if trace_manager:
	//     trace_manager.add_trace_task(
	//         TraceTask(
	//             TraceTaskName.TOOL_TRACE,
	//             message_id=message_id,
	//             tool_name=tool_name,
	//             tool_inputs=tool_inputs,
	//             tool_outputs=tool_outputs,
	//             timer=timer,
	//         )
	//     )
}
func (handler *AgentCallbackHandler) OnToolError(err error) {
	/*Do nothing.*/
	if confy.GetWithDefault[bool]("DEBUG", false) {
		print_text("\n[OnToolError] Error: "+err.Error()+"\n", "red")
	}
}
func (handler *AgentCallbackHandler) OnAgentStart(thought string) {
	/*Run on agent start.*/
	if confy.GetWithDefault[bool]("DEBUG", false) {
		if thought != "" {
			print_text(
				"\n[OnAgentStart] \nCurrent Loop: "+strconv.Itoa(handler.current_loop)+"\nThought: "+thought+"\n",
				handler.color,
			)
		}
	} else {
		print_text("\n[OnAgentStart] \nCurrent Loop: "+strconv.Itoa(handler.current_loop)+"\n", handler.color)
	}
}
func (handler *AgentCallbackHandler) OnAgentFinish(color string) {
	/*Run on agent end.*/
	if confy.GetWithDefault[bool]("DEBUG", false) {
		print_text("\n[OnAgentFinish]\n Loop: "+strconv.Itoa(handler.current_loop)+"\n", handler.color)
	}
	handler.current_loop += 1

}
func (handler *AgentCallbackHandler) ignore_agent() bool {
	/*Whether to ignore agent callbacks.*/
	return !confy.GetWithDefault[bool]("DEBUG", false)

}
func (handler *AgentCallbackHandler) ignore_chat_model() bool {
	/*Whether to ignore chat model callbacks.*/
	return !confy.GetWithDefault[bool]("DEBUG", false)
}
