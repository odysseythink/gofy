package services

import "fmt"

type AudioService struct{}

func (s *AudioService) TranscriptASR(appID string, fileData []byte, filename string) (string, error) {
	// TODO: Integrate with speech-to-text model via model_runtime
	return "", fmt.Errorf("ASR not yet implemented")
}

func (s *AudioService) TranscriptTTS(appID string, text string, voice string) ([]byte, error) {
	// TODO: Integrate with TTS model via model_runtime
	return nil, fmt.Errorf("TTS not yet implemented")
}

func (s *AudioService) GetTTSVoices(tenantID string, language string) []map[string]any {
	return []map[string]any{}
}
