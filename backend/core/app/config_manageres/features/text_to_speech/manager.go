package texttospeech

import (
	appconfigentities "github.com/odysseythink/gofy/backend/entities/app/config"
)

type TextToSpeechConfigManager struct{}

func (mgr *TextToSpeechConfigManager) Convert(config map[string]any) *appconfigentities.TextToSpeechEntity {
	/*
	   Convert model config to model config

	   :param config: model config args
	*/
	var text_to_speech *appconfigentities.TextToSpeechEntity
	var text_to_speech_dict map[string]any
	if _, ok := config["text_to_speech"]; ok {
		if _, ok := config["text_to_speech"].(map[string]any); ok {
			text_to_speech_dict = config["text_to_speech"].(map[string]any)
		}
	}
	if len(text_to_speech_dict) > 0 {
		var enabled bool
		if _, ok := text_to_speech_dict["enabled"]; ok {
			if _, ok := text_to_speech_dict["enabled"].(bool); ok {
				enabled = text_to_speech_dict["enabled"].(bool)
			}
		}
		if enabled {
			text_to_speech = &appconfigentities.TextToSpeechEntity{
				Enabled: enabled,
				// voice=text_to_speech_dict.get("voice"),
				// language=text_to_speech_dict.get("language"),
			}
			if _, ok := text_to_speech_dict["voice"]; ok {
				if _, ok := text_to_speech_dict["voice"].(string); ok {
					text_to_speech.Voice = text_to_speech_dict["voice"].(string)
				}
			}
			if _, ok := text_to_speech_dict["language"]; ok {
				if _, ok := text_to_speech_dict["language"].(string); ok {
					text_to_speech.Language = text_to_speech_dict["language"].(string)
				}
			}
		}
	}
	return text_to_speech
}

func (mgr *TextToSpeechConfigManager) ValidateAndSetDefaults(config map[string]any) (map[string]any, []string) {
	/*
	   Validate and set defaults for text to speech feature

	   :param config: app model config args
	*/
	var text_to_speech map[string]any
	if _, ok := config["text_to_speech"]; ok {
		if _, ok := config["text_to_speech"].(map[string]any); ok {
			text_to_speech = config["text_to_speech"].(map[string]any)
		}
	}
	if len(text_to_speech) == 0 {
		config["text_to_speech"] = map[string]any{"enabled": false, "voice": "", "language": ""}
	}
	var enabled bool
	if _, ok := text_to_speech["enabled"]; !ok {
		if _, ok := text_to_speech["enabled"].(bool); ok {
			enabled = text_to_speech["enabled"].(bool)
		}
	}

	if !enabled {
		config["text_to_speech"].(map[string]any)["enabled"] = false
		config["text_to_speech"].(map[string]any)["voice"] = ""
		config["text_to_speech"].(map[string]any)["language"] = ""
	}

	return config, []string{"text_to_speech"}
}
