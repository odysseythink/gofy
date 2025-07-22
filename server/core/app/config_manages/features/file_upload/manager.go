package fileupload

import "mlib.com/gofy/server/core/file"

type FileUploadConfigManager struct{}

func (mgr *FileUploadConfigManager) Convert(config map[string]any, is_vision bool /* = True*/) *file.FileUploadConfig {
	/*
	   Convert model config to model config

	   :param config: model config args
	   :param is_vision: if True, the feature is vision feature
	*/
	var file_upload_dict map[string]any
	if _, ok := config["file_upload"]; ok {
		if _, ok := config["file_upload"].(map[string]any); ok {
			file_upload_dict = config["file_upload"].(map[string]any)
		}
	}
	if len(file_upload_dict) > 0 {
		var enabled bool
		if _, ok := file_upload_dict["enabled"]; ok {
			if _, ok := file_upload_dict["enabled"].(bool); ok {
				enabled = file_upload_dict["enabled"].(bool)
			}
		}
		if enabled {
			var transform_methods []any
			if _, ok := file_upload_dict["allowed_file_upload_methods"]; ok {
				if _, ok := file_upload_dict["allowed_file_upload_methods"].([]any); ok {
					transform_methods = file_upload_dict["allowed_file_upload_methods"].([]any)
				}
			}

			data := map[string]any{
				"image_config": map[string]any{
					"number_limits":    file_upload_dict["number_limits"],
					"transfer_methods": transform_methods,
				},
			}

			if is_vision {
				var image map[string]any
				var detail string = "low"
				if _, ok := file_upload_dict["image"]; ok {
					if _, ok := file_upload_dict["image"].(map[string]any); ok {
						image = file_upload_dict["image"].(map[string]any)
						if _, ok := image["detail"]; ok {
							if _, ok := image["detail"].(string); ok {
								detail = image["detail"].(string)
							}
						}
					}
				}
				data["image_config"].(map[string]any)["detail"] = detail
			}

			return file.NewFileUploadConfigFromDict(data)
		}
	}
	return nil
}
func (mgr *FileUploadConfigManager) ValidateAndSetDefaults(config map[string]any) (map[string]any, []string) {
	/*
	   Validate and set defaults for file upload feature

	   :param config: app model config args
	*/
	var file_upload map[string]any
	if _, ok := config["file_upload"]; ok {
		if _, ok := config["file_upload"].(map[string]any); ok {
			file_upload = config["file_upload"].(map[string]any)
		}
	}
	if len(file_upload) == 0 {
		config["file_upload"] = map[string]any{}
	} else {
		file.NewFileUploadConfigFromDict(file_upload)
		// FileUploadConfig.model_validate(config["file_upload"])
	}
	return config, []string{"file_upload"}
}
