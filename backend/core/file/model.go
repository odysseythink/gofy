package file

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	modelruntimeentities "github.com/odysseythink/gofy/backend/entities/model_runtime"
	"github.com/odysseythink/gofy/backend/utils"
	"github.com/odysseythink/mlog"
)

// ImageConfig represents the configuration for image uploads.
type ImageConfig struct {
	NumberLimits    int                                                  `json:"number_limits"`
	TransferMethods []FileTransferMethod                                 `json:"transfer_methods"`
	Detail          modelruntimeentities.ImagePromptMessageContentDETAIL `json:"detail,omitempty"`
}

// FileUploadConfig represents the configuration for file uploads.
type FileUploadConfig struct {
	ImageConfig              *ImageConfig         `json:"image_config,omitempty"`
	AllowedFileTypes         []FileType           `json:"allowed_file_types"`
	AllowedFileExtensions    []string             `json:"allowed_file_extensions"`
	AllowedFileUploadMethods []FileTransferMethod `json:"allowed_file_upload_methods"`
	NumberLimits             int                  `json:"number_limits"`
}

func NewFileUploadConfigFromDict(config map[string]any) *FileUploadConfig {
	cfg := new(FileUploadConfig)
	tmpbindata, _ := json.Marshal(config)
	err := json.Unmarshal(tmpbindata, cfg)
	if err != nil {
		mlog.Errorf("json unmarshal failed:%v", err)
		return nil
	}
	return cfg
}

type File struct {
	GofyModelIdentity string `json:"gofy_model_identity"`

	ID             string/*# message file id */ `json:"id"`
	TenantID       string             `json:"tenant_id"`
	Type           FileType           `json:"type"`
	TransferMethod FileTransferMethod `json:"transfer_method"`
	RemoteURL      string             //# remote url `json:"remote_url"`
	RelatedID      string             `json:"related_id"`
	Filename       string             `json:"filename"`
	Extension      string/*Field(default=None, description="File extension, should contains dot")*/ `json:"extension"`
	MimeType       string `json:"mime_type"`
	Size           int    `json:"size"`

	//# Those properties are private, should not be exposed to the outside.
	StorageKey string `json:"-"`
}

func NewFile(id string,
	tenantID string,
	fileType FileType,
	transferMethod FileTransferMethod,
	remoteURL string,
	relatedID string,
	filename string,
	extension string,
	mimeType string,
	size int,
	storageKey string) *File {
	return &File{
		GofyModelIdentity: FILE_MODEL_IDENTITY,
		ID:                id,
		TenantID:          tenantID,
		Type:              fileType,
		TransferMethod:    transferMethod,
		RemoteURL:         remoteURL,
		RelatedID:         relatedID,
		Filename:          filename,
		Extension:         extension,
		MimeType:          mimeType,
		StorageKey:        storageKey,
		Size:              size,
	}
}

// ToDict converts the File to a map.
func (f *File) ToDict() map[string]any {
	bindata, _ := json.Marshal(f)
	data := map[string]any{}
	json.Unmarshal(bindata, &data)
	data["url"] = f.GenerateURL()
	return data
}

// Markdown generates the Markdown representation of the File.
func (f *File) Markdown() string {
	url := f.GenerateURL()
	if f.Type == File_IMAGE {
		if f.Filename == "" {
			return fmt.Sprintf("![](%s)", url)
		} else {
			return fmt.Sprintf("![%s](%s)", f.Filename, url)
		}
	} else {
		if f.Filename == "" {
			return fmt.Sprintf("![%s](%s)", url, url)
		} else {
			return fmt.Sprintf("![%s](%s)", f.Filename, url)
		}
	}
}

// GenerateURL generates the URL for the File.
func (f *File) GenerateURL() string {
	if f.Type == File_IMAGE {
		switch f.TransferMethod {
		case FileTransferMethod_REMOTE_URL:
			return f.RemoteURL
		case FileTransferMethod_LOCAL_FILE:
			if f.RelatedID == "" {
				panic(exceptions.NewValueError("Missing file related_id"))
			}
			return GetSignedFileURL(f.RelatedID)
		case FileTransferMethod_TOOL_FILE:
			if f.RelatedID == "" || f.Extension == "" {
				panic(exceptions.NewValueError("Missing file related_id or extension"))
			}
			return SignFile(f.RelatedID, f.Extension)
		default:
			return ""
		}
	} else {
		switch f.TransferMethod {
		case FileTransferMethod_REMOTE_URL:
			return f.RemoteURL
		case FileTransferMethod_LOCAL_FILE:
			if f.RelatedID == "" {
				panic(exceptions.NewValueError("Missing file related_id"))
			}
			return GetSignedFileURL(f.RelatedID)
		case FileTransferMethod_TOOL_FILE:
			if f.RelatedID == "" || f.Extension == "" {
				panic("Missing file related_id or extension")
			}
			return SignFile(f.RelatedID, f.Extension)
		default:
			return ""
		}
	}
}

// FilenameOrDefault returns the filename or a default value if it is nil.
func (f *File) FilenameOrDefault() string {
	if f.Filename != "" {
		return f.Filename
	}
	return ""
}

// Validate validates the File after initialization.
func (f *File) Validate() error {
	switch f.TransferMethod {
	case FileTransferMethod_REMOTE_URL:
		if f.RemoteURL == "" {
			return exceptions.NewValueError("missing file url")
		}
		if _, err := url.Parse(f.RemoteURL); err != nil || !strings.HasPrefix(f.RemoteURL, "http") {
			return exceptions.NewValueError("invalid file url")
		}
	case FileTransferMethod_LOCAL_FILE:
		if f.RelatedID == "" {
			return exceptions.NewValueError("missing file related_id")
		}
	case FileTransferMethod_TOOL_FILE:
		if f.RelatedID == "" {
			return exceptions.NewValueError("missing file related_id")
		}
	default:
		return exceptions.NewValueError("unknown transfer method")
	}
	return nil
}

// Helper functions (need to be implemented)
func getSignedFileURL(uploadFileID string) string {
	// Implement the logic to generate a signed URL for the file
	return "signed-url"
}

func SignFile(tool_file_id string, extension string) string {
	base_url := confy.Get[string]("FILES_URL")
	file_preview_url := fmt.Sprintf("%s/files/tools/%s%s", base_url, tool_file_id, extension)

	timestamp := strconv.Itoa(int(time.Now().Unix()))
	nonce := utils.GenerateRandomHex(16)
	data_to_sign := fmt.Sprintf("file-preview|%s|%s|%s", tool_file_id, timestamp, nonce)
	secretKey := []byte(confy.Get[string]("SECRET_KEY")) // 替换为实际的密钥

	// 创建 HMAC-SHA256 签名
	sign := hmac.New(sha256.New, secretKey)
	sign.Write([]byte(data_to_sign))
	sign_bytes := sign.Sum(nil)
	encoded_sign := base64.URLEncoding.EncodeToString(sign_bytes)

	return fmt.Sprintf("%s?timestamp=%s&nonce=%s&sign=%s", file_preview_url, url.QueryEscape(timestamp), url.QueryEscape(nonce), url.QueryEscape(encoded_sign))
}
