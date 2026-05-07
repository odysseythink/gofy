package utils

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/mlog"
)

func ParseJsonMarkdown(json_string string) map[string]any {
	mlog.Debugf("------json_string=%#v", json_string)
	// Get json from the backticks/braces
	json_string = strings.TrimSpace(json_string)
	starts := []string{"```json", "```", "``", "`", "{"}
	ends := []string{"```", "``", "`", "}"}
	end_index := -1
	start_index := 0
	parsed := map[string]any{}
	for _, s := range starts {
		start_index = strings.Index(json_string, s)
		if start_index != -1 {
			if json_string[start_index] != '{' {
				start_index += len(s)
			}
			break
		}
	}
	if start_index != -1 {
		for _, e := range ends {
			end_index = strings.LastIndex(json_string[start_index:], e)
			if end_index != -1 {
				end_index += start_index
				if json_string[end_index] == '}' {
					end_index += 1
				}
				break
			}
		}
	}
	if start_index != -1 && end_index != -1 && start_index < end_index {
		extracted_content := json_string[start_index:end_index]
		extracted_content = strings.TrimPrefix(extracted_content, "\n")
		if strings.HasPrefix(extracted_content, "[") {
			list_parsed := []map[string]any{}
			err := json.Unmarshal([]byte(extracted_content), &list_parsed)
			if err != nil || len(list_parsed) != 1 {
				mlog.Errorf("json marshal json block(%s) failed:%v", extracted_content, err)
				panic(exceptions.NewValueError("json unmarshal json block failed."))
			}
			parsed = maps.Clone(list_parsed[0])
		} else {
			err := json.Unmarshal([]byte(extracted_content), &parsed)
			if err != nil {
				mlog.Errorf("json marshal json block(%s) failed:%v", extracted_content, err)
				panic(exceptions.NewValueError("json unmarshal json block failed."))
			}
		}
	} else {
		mlog.Error("could not find json block in the output.")
		panic(exceptions.NewValueError("could not find json block in the output."))
	}
	return parsed
}

func ParseAndCheckJsonMarkdown(text string, expected_keys []string) map[string]any {
	json_obj := func(text string) map[string]any {
		if r := recover(); r != nil {
			if exp, ok := r.(error); ok {
				mlog.Errorf("accur error:%v", exp)
				panic(exceptions.NewOutputParserError(fmt.Sprintf("got invalid json object. error: %v", exp)))
			} else {
				panic(r)
			}
		}
		return ParseJsonMarkdown(text)
	}(text)

	for _, key := range expected_keys {
		if _, ok := json_obj[key]; !ok {
			panic(exceptions.NewOutputParserError(fmt.Sprintf("got invalid return object. expected key `%s` to be present, but got %#v", key, json_obj)))
		}
	}
	return json_obj
}
