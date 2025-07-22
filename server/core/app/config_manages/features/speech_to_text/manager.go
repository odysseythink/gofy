package speechtotext

type SpeechToTextConfigManager struct{}

func (mgr *SpeechToTextConfigManager) Convert(config map[string]any) bool {
	/*
	   Convert model config to model config

	   :param config: model config args
	*/
	speech_to_text := false
	var speech_to_text_dict map[string]any
	if _, ok := config["speech_to_text"]; ok {
		if _, ok := config["speech_to_text"].(map[string]any); ok {
			speech_to_text_dict = config["speech_to_text"].(map[string]any)
		}
	}

	if len(speech_to_text_dict) > 0 {
		if _, ok := speech_to_text_dict["enabled"]; ok {
			if _, ok := speech_to_text_dict["enabled"].(bool); ok {
				speech_to_text = speech_to_text_dict["enabled"].(bool)
			}
		}
	}
	return speech_to_text
}

func (mgr *SpeechToTextConfigManager) ValidateAndSetDefaults(config map[string]any) (map[string]any, []string) {
	/*
		Validate and set defaults for speech to text feature

		:param config: app model config args
	*/

	if _, ok := config["speech_to_text"]; ok {
		if _, ok := config["speech_to_text"].(map[string]any); !ok {
			config["speech_to_text"] = map[string]any{"enabled": false}
		}
	} else {
		config["speech_to_text"] = map[string]any{"enabled": false}
	}

	if _, ok := config["speech_to_text"].(map[string]any)["enabled"]; ok {
		if _, ok := config["speech_to_text"].(map[string]any)["enabled"].(bool); !ok {
			config["speech_to_text"].(map[string]any)["enabled"] = false
		}
	} else {
		config["speech_to_text"].(map[string]any)["enabled"] = false
	}

	return config, []string{"speech_to_text"}
}
