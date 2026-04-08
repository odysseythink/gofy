package callbacks

import "fmt"

var (
	_TEXT_COLOR_MAPPING = map[string]string{
		"blue":   "36;1",
		"yellow": "33;1",
		"pink":   "38;5;200",
		"green":  "32;1",
		"red":    "31;1",
	}
)

type BaseCallback struct {
	RaiseError bool
}

func (cb *BaseCallback) PrintText(text string, color string, endarg ...string /*= "\n"*/) {
	/*Print text with highlighting and no end characters.*/
	end := ""
	if len(endarg) > 0 {
		end = endarg[0]
	}
	if end == "" {
		end = "\n"
	}
	text_to_print := ""
	if color != "" {
		text_to_print = cb.getColoredText(text, color)
	} else {
		text_to_print = text
	}
	fmt.Println(text_to_print, end)
}

func (cb *BaseCallback) getColoredText(text string, color string) string {
	/*Get colored text.*/
	color_str := _TEXT_COLOR_MAPPING[color]
	return fmt.Sprintf("\u001b[%sm\033[1;3m%s\u001b[0m", color_str, text)
}
