package file

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"mlib.com/confy"
	"mlib.com/gofy/server/core/exceptions"
	modelruntimeentities "mlib.com/gofy/server/entities/model_runtime"
)

// GetAttr retrieves the specified attribute of a file.
func GetAttr(file *File, attr FileAttribute) any {
	switch attr {
	case FileAttribute_TYPE:
		return file.Type
	case FileAttribute_SIZE:
		return file.Size
	case FileAttribute_NAME:
		return file.Filename
	case FileAttribute_MIME_TYPE:
		return file.MimeType
	case FileAttribute_TRANSFER_METHOD:
		return file.TransferMethod
	case FileAttribute_URL:
		return file.RemoteURL
	case FileAttribute_EXTENSION:
		return file.Extension
	default:
		return nil
	}
}

// ToPromptMessageContent converts a file to a prompt message content.
func ToPromptMessageContent(f *File, imageDetailConfig modelruntimeentities.ImagePromptMessageContentDETAIL) (modelruntimeentities.MultiModalPromptMessageContent, error) {
	if f.Extension == "" {
		return nil, exceptions.NewValueError("missing file extension")
	}
	if f.MimeType == "" {
		return nil, exceptions.NewValueError("missing file mime_type")
	}

	params := map[string]any{
		"base64_data": "",
		"url":         "",
		"format":      strings.TrimPrefix(f.Extension, "."),
		"mime_type":   f.MimeType,
	}
	if confy.Get[string]("MULTIMODAL_SEND_FORMAT") == "base64" {
		params["base64_data"] = getEncodedString(f)
	}
	if confy.Get[string]("MULTIMODAL_SEND_FORMAT") == "url" {
		params["url"], _ = toURL(f)
	}

	if f.Type == File_IMAGE {
		if imageDetailConfig == "" {
			params["detail"] = modelruntimeentities.ImagePromptMessageContentDETAIL_LOW
		} else {
			params["detail"] = string(imageDetailConfig)
		}
	}
	switch f.Type {
	case File_IMAGE:
		return modelruntimeentities.NewImagePromptMessageContent(params), nil
	case File_AUDIO:
		return modelruntimeentities.NewAudioPromptMessageContent(params), nil
	case File_VIDEO:
		return modelruntimeentities.NewVideoPromptMessageContent(params), nil
	case File_DOCUMENT:
		return modelruntimeentities.NewDocumentPromptMessageContent(params), nil
	default:
		return nil, exceptions.NewValueError(fmt.Sprintf("file type %s is not supported", f.Type))
	}
}

// Download downloads the file content based on the transfer method.
func Download(f *File) ([]byte, error) {
	switch f.TransferMethod {
	case FileTransferMethod_REMOTE_URL:
		resp, err := http.Get(f.RemoteURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to download file: %s", resp.Status)
		}
		return io.ReadAll(resp.Body)
	case FileTransferMethod_LOCAL_FILE:
		return downloadFileContent(f.StorageKey)
	case FileTransferMethod_TOOL_FILE:
		return downloadFileContent(f.StorageKey)
	default:
		return nil, fmt.Errorf("unsupported transfer method: %d", f.TransferMethod)
	}
}

// downloadFileContent downloads the file content from storage.
func downloadFileContent(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

// getEncodedString returns the base64 encoded string of the file content.
func getEncodedString(f *File) string {
	data, err := Download(f)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}

// toURL returns the URL of the file based on the transfer method.
func toURL(f *File) (string, error) {
	switch f.TransferMethod {
	case FileTransferMethod_REMOTE_URL:
		if f.RemoteURL == "" {
			return "", exceptions.NewValueError("missing file remote_url")
		}
		return f.RemoteURL, nil
	case FileTransferMethod_LOCAL_FILE:
		if f.RelatedID == "" {
			return "", exceptions.NewValueError("missing file related_id")
		}
		return getSignedFileURL(f.RelatedID), nil
	case FileTransferMethod_TOOL_FILE:
		if f.RelatedID == "" || f.Extension == "" {
			return "", exceptions.NewValueError("missing file related_id or extension")
		}
		return SignFile(f.RelatedID, f.Extension), nil
	default:
		return "", exceptions.NewValueError(fmt.Sprintf("unsupported transfer method: %v", f.TransferMethod))
	}
}

// // getSignedFileURL generates a signed URL for the file.
// func getSignedFileURL(uploadFileID string) string {
// 	// Implement the logic to generate a signed URL
// 	return "signed-url"
// }

// // signToolFile generates a signed URL for the tool file.
// func signToolFile(toolFileID, extension string) string {
// 	// Implement the logic to sign a tool file
// 	return "signed-tool-file-url"
// }

// func main() {
// 	// Example usage
// 	file := &File{
// 		Type:          FILE_TYPE_IMAGE,
// 		Size:          1024,
// 		Filename:      "example.jpg",
// 		MimeType:      "image/jpeg",
// 		TransferMethod: FILE_TRANSFER_METHOD_REMOTE_URL,
// 		RemoteURL:     "https://example.com/file.jpg",
// 		StorageKey:    "path/to/storage",
// 		RelatedID:     uuid.NewV4().String(),
// 		Extension:     ".jpg",
// 	}

// 	// Get attribute
// 	fileType := GetAttr(file, FILE_ATTRIBUTE_TYPE)
// 	fmt.Println("File Type:", fileType)

// 	// Convert to prompt message content
// 	promptContent, err := ToPromptMessageContent(file, nil)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("Prompt Content:", promptContent)

// 	// Download file
// 	data, err := Download(file)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("Downloaded data size:", len(data))

// 	// Get URL
// 	url, err := toURL(file)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("File URL:", url)
// }

// def get_attr(*, file: File, attr: FileAttribute):
//     match attr:
//         case FileAttribute.TYPE:
//             return file.type.value
//         case FileAttribute.SIZE:
//             return file.size
//         case FileAttribute.NAME:
//             return file.filename
//         case FileAttribute.MIME_TYPE:
//             return file.mime_type
//         case FileAttribute.TRANSFER_METHOD:
//             return file.transfer_method.value
//         case FileAttribute.URL:
//             return file.remote_url
//         case FileAttribute.EXTENSION:
//             return file.extension

// def to_prompt_message_content(
//     f: File,
//     /,
//     *,
//     image_detail_config: ImagePromptMessageContent.DETAIL | None = None,
// ) -> MultiModalPromptMessageContent:
//     if f.extension is None:
//         raise ValueError("Missing file extension")
//     if f.mime_type is None:
//         raise ValueError("Missing file mime_type")

//     params = {
//         "base64_data": _get_encoded_string(f) if dify_config.MULTIMODAL_SEND_FORMAT == "base64" else "",
//         "url": _to_url(f) if dify_config.MULTIMODAL_SEND_FORMAT == "url" else "",
//         "format": f.extension.removeprefix("."),
//         "mime_type": f.mime_type,
//     }
//     if f.type == FileType.IMAGE:
//         params["detail"] = image_detail_config or ImagePromptMessageContent.DETAIL.LOW

//     prompt_class_map: Mapping[FileType, type[MultiModalPromptMessageContent]] = {
//         FileType.IMAGE: ImagePromptMessageContent,
//         FileType.AUDIO: AudioPromptMessageContent,
//         FileType.VIDEO: VideoPromptMessageContent,
//         FileType.DOCUMENT: DocumentPromptMessageContent,
//     }

//     try:
//         return prompt_class_map[f.type].model_validate(params)
//     except KeyError:
//         raise ValueError(f"file type {f.type} is not supported")

// def download(f: File, /):
//     if f.transfer_method in (FileTransferMethod.TOOL_FILE, FileTransferMethod.LOCAL_FILE):
//         return _download_file_content(f._storage_key)
//     elif f.transfer_method == FileTransferMethod.REMOTE_URL:
//         response = ssrf_proxy.get(f.remote_url, follow_redirects=True)
//         response.raise_for_status()
//         return response.content
//     raise ValueError(f"unsupported transfer method: {f.transfer_method}")

// def _download_file_content(path: str, /):
//     """
//     Download and return the contents of a file as bytes.

//     This function loads the file from storage and ensures it's in bytes format.

//     Args:
//         path (str): The path to the file in storage.

//     Returns:
//         bytes: The contents of the file as a bytes object.

//     Raises:
//         ValueError: If the loaded file is not a bytes object.
//     """
//     data = storage.load(path, stream=False)
//     if not isinstance(data, bytes):
//         raise ValueError(f"file {path} is not a bytes object")
//     return data

// def _get_encoded_string(f: File, /):
//     match f.transfer_method:
//         case FileTransferMethod.REMOTE_URL:
//             response = ssrf_proxy.get(f.remote_url, follow_redirects=True)
//             response.raise_for_status()
//             data = response.content
//         case FileTransferMethod.LOCAL_FILE:
//             data = _download_file_content(f._storage_key)
//         case FileTransferMethod.TOOL_FILE:
//             data = _download_file_content(f._storage_key)

//     encoded_string = base64.b64encode(data).decode("utf-8")
//     return encoded_string

// def _to_url(f: File, /):
//     if f.transfer_method == FileTransferMethod.REMOTE_URL:
//         if f.remote_url is None:
//             raise ValueError("Missing file remote_url")
//         return f.remote_url
//     elif f.transfer_method == FileTransferMethod.LOCAL_FILE:
//         if f.related_id is None:
//             raise ValueError("Missing file related_id")
//         return f.remote_url or helpers.get_signed_file_url(upload_file_id=f.related_id)
//     elif f.transfer_method == FileTransferMethod.TOOL_FILE:
//         # add sign url
//         if f.related_id is None or f.extension is None:
//             raise ValueError("Missing file related_id or extension")
//         return ToolFileParser.get_tool_file_manager().sign_file(tool_file_id=f.related_id, extension=f.extension)
//     else:
//         raise ValueError(f"Unsupported transfer method: {f.transfer_method}")
