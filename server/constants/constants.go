package constants

const (
	HIDDEN_VALUE = "[__HIDDEN__]"
	UUID_NIL     = "00000000-0000-0000-0000-000000000000"
)

var (
	SUPPORT_URL_CONTENT_TYPES = []string{"application/pdf", "text/plain", "application/json"}
	ALLOW_CREATE_APP_MODES    = []string{"chat", "agent-chat", "advanced-chat", "workflow", "completion"}
	IMAGE_EXTENSIONS          = []string{"jpg", "jpeg", "png", "webp", "gif", "svg"}

	VIDEO_EXTENSIONS = []string{"mp4", "mov", "mpeg", "mpga"}

	AUDIO_EXTENSIONS = []string{"mp3", "m4a", "wav", "webm", "amr"}

	DOCUMENT_EXTENSIONS = []string{"txt", "markdown", "md", "mdx", "pdf", "html", "htm", "xlsx", "xls", "docx", "csv", "eml", "msg", "pptx", "xml", "epub"}
)
