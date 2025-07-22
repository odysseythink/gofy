package response

type SimpleEndUserResponse struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	IsAnonymous bool   `json:"is_anonymous"`
	SessionID   string `json:"session_id"`
}
