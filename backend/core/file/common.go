package file

import (
	"encoding/json"
)

func HandleSpecialValues(value map[string]any) map[string]any {
	result := handleSpecialValues(value)
	if v, ok := result.(map[string]any); ok || v == nil {
		return result.(map[string]any)
	}
	var rsp map[string]any
	bindata, _ := json.Marshal(result)
	json.Unmarshal(bindata, &rsp)
	return rsp
}
func handleSpecialValues(value any) any {
	if value == nil {
		return value
	}
	if _, ok := value.(map[string]any); ok {
		res := map[string]any{}
		for k, v := range value.(map[string]any) {
			res[k] = handleSpecialValues(v)
		}
		return res
	}
	if _, ok := value.([]any); ok {
		res_list := []any{}
		for _, item := range value.([]any) {
			res_list = append(res_list, handleSpecialValues(item))
		}
		return res_list
	}
	if _, ok := value.(*File); ok {
		return value.(*File).ToDict()
	}
	return value
}

func GetFileVarFromValue(value any /* Union[dict, list]*/) map[string]any {
	/*
		Get file var from value
		:param value: variable value
		:return
	*/
	if value == nil {
		return nil
	}
	switch data := value.(type) {
	case map[string]any:
		if _, ok := data["gofy_model_identity"]; ok {
			if _, ok := data["gofy_model_identity"].(string); ok {
				if data["gofy_model_identity"].(string) == FILE_MODEL_IDENTITY {
					return data
				}
			}
		}
	case *File:
		return data.ToDict()
	}

	return nil
}
func FetchFilesFromVariableValue(value any /* Union[dict, list]*/) []map[string]any {
	/*
		Fetch files from variable value
		:param value: variable value
		:return
	*/
	if value == nil {
		return nil
	}
	files := []map[string]any{}
	switch data := value.(type) {
	case []any:
		for _, item := range data {
			f := GetFileVarFromValue(item)
			if f != nil {
				files = append(files, f)
			}
		}
	case map[string]any:
		f := GetFileVarFromValue(data)
		if f != nil {
			files = append(files, f)
		}
	}

	return files
}

func FetchFilesFromNodeOutputs(outputs_dict map[string]any) []map[string]any {
	/*
		Fetch files from node outputs
		:param outputs_dict: node outputs dict
		:return
	*/
	if outputs_dict == nil {
		return nil
	}
	files := []map[string]any{}
	for _, output_value := range outputs_dict {
		f := FetchFilesFromVariableValue(output_value)
		if f != nil {
			for _, v := range f {
				if v != nil {
					files = append(files, v)
				}
			}
		}
	}

	return files
}
