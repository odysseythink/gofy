package filefactory

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/odysseythink/gofy/backend/constants"
	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/gofy/backend/core/file"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
)

func get_file_type_by_mimetype(mime_type string) file.FileType {
	file_type := file.FileType("")
	if strings.Contains(mime_type, "image") {
		file_type = file.File_IMAGE
	} else if strings.Contains(mime_type, "video") {
		file_type = file.File_VIDEO
	} else if strings.Contains(mime_type, "audio") {
		file_type = file.File_AUDIO
	} else if strings.Contains(mime_type, "text") && strings.Contains(mime_type, "pdf") {
		file_type = file.File_DOCUMENT
	} else {
		file_type = file.File_CUSTOM
	}
	return file_type
}

func get_file_type_by_extension(extension string) file.FileType {
	extension = strings.ToLower(strings.TrimPrefix(extension, "."))
	if slices.Contains(constants.IMAGE_EXTENSIONS, extension) {
		return file.File_IMAGE
	} else if slices.Contains(constants.VIDEO_EXTENSIONS, extension) {
		return file.File_VIDEO
	} else if slices.Contains(constants.AUDIO_EXTENSIONS, extension) {
		return file.File_AUDIO
	} else if slices.Contains(constants.DOCUMENT_EXTENSIONS, extension) {
		return file.File_DOCUMENT
	}
	return file.FileType("")
}
func standardize_file_type(file_type file.FileType, extension string, mime_type string) file.FileType {
	/*
	   If custom type, try to guess the file type by extension and mime_type.
	*/
	if file_type != file.File_CUSTOM {
		return file.FileType(file_type)
	}
	var guessed_type file.FileType
	if extension != "" {
		guessed_type = get_file_type_by_extension(extension)
	}
	if guessed_type == file.FileType("") && mime_type != "" {
		guessed_type = get_file_type_by_mimetype(mime_type)
	}
	if guessed_type != file.FileType("") {
		return guessed_type
	}
	return file.File_CUSTOM
}

func get_remote_file_info(urlStr string) (mime string, filename string, size int64) {
	// 解析 URL
	url, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// 提取文件名
	filename = path.Base(url.Path)
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}
	if filename == "" {
		filename = "unknown_file"
	}

	// 猜测 MIME 类型
	mime = http.DetectContentType([]byte(filename))

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 创建 HEAD 请求
	req, err := http.NewRequest("HEAD", urlStr, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return
	}

	// 获取文件名
	if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "" {
		// 解析 Content-Disposition 头
		disposition := strings.Split(contentDisposition, ";")
		for _, part := range disposition {
			if strings.Contains(part, "filename=") {
				filename = strings.TrimSpace(strings.TrimPrefix(part, "filename="))
				if strings.HasPrefix(filename, `"`) && strings.HasSuffix(filename, `"`) {
					filename = filename[1 : len(filename)-1]
				}
				break
			}
		}
	}

	// 获取文件大小
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		if size, err = strconv.ParseInt(contentLength, 10, 64); err != nil {
			fmt.Println("Error parsing Content-Length:", err)
			size = -1
		}
	}

	// 获取 MIME 类型
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		mime = contentType
	}

	return
}

func build_from_tool_file(
	mapping map[string]any,
	tenant_id string,
	transfer_method file.FileTransferMethod,
) *file.File {
	var tool_file_id string
	if _, ok := mapping["tool_file_id"]; ok {
		if _, ok := mapping["tool_file_id"].(string); ok {
			tool_file_id = mapping["tool_file_id"].(string)
		}
	}
	tool_file := new(models.ToolFile)
	err := dbengine.Instance().DB.Model(&models.ToolFile{}).Where("id = ? and tenant_id = ?", tool_file_id, tenant_id).First(tool_file).Error
	if err != nil {
		panic(exceptions.NewValueError(fmt.Sprintf("ToolFile %s not found", tool_file_id)))
	}
	extension := ".bin"
	if strings.Contains(tool_file.FileKey, ".") {
		tmplist := strings.Split(tool_file.FileKey, ".")
		extension = "." + tmplist[len(tmplist)-1]
	}
	file_type := file.FileType("custom")
	if _, ok := mapping["type"]; ok {
		if _, ok := mapping["type"].(string); ok {
			file_type = file.FileType(mapping["type"].(string))
		}
	}

	file_type = standardize_file_type(file_type, extension, tool_file.Mimetype)
	var id string
	if _, ok := mapping["id"]; ok {
		if _, ok := mapping["id"].(string); ok {
			id = mapping["id"].(string)
		}
	}
	return &file.File{
		ID:             id,
		TenantID:       tenant_id,
		Type:           file.FileType(file_type),
		TransferMethod: transfer_method,
		RemoteURL:      tool_file.OriginalURL,
		RelatedID:      tool_file.ID,
		Filename:       tool_file.Name,
		Extension:      extension,
		MimeType:       tool_file.Mimetype,
		Size:           tool_file.Size,
		StorageKey:     tool_file.FileKey,
	}
}

func build_from_remote_url(
	mapping map[string]any,
	tenant_id string,
	transfer_method file.FileTransferMethod,
) *file.File {
	var url string
	if _, ok := mapping["url"]; ok {
		if _, ok := mapping["url"].(string); ok {
			url = mapping["url"].(string)
		}
	}
	if url == "" {
		if _, ok := mapping["remote_url"]; ok {
			if _, ok := mapping["remote_url"].(string); ok {
				url = mapping["remote_url"].(string)
			}
		}
	}

	if url == "" {
		panic(exceptions.NewValueError("Invalid file url"))
	}
	mime_type, filename, file_size := get_remote_file_info(url)
	extension := ""
	if mime_type != "" {
		extension = http.DetectContentType([]byte(filename))
	} else if strings.Contains(filename, ".") {
		extension = "." + strings.Split(filename, ".")[len(strings.Split(filename, "."))-1]
	} else {
		extension = ".bin"
	}
	file_type := file.FileType("custom")
	if _, ok := mapping["type"]; ok {
		if _, ok := mapping["type"].(string); ok {
			file_type = file.FileType(mapping["type"].(string))
		}
	}

	file_type = standardize_file_type(file_type, extension, mime_type)
	var id string
	if _, ok := mapping["id"]; ok {
		if _, ok := mapping["id"].(string); ok {
			id = mapping["id"].(string)
		}
	}
	return &file.File{
		ID:             id,
		TenantID:       tenant_id,
		Type:           file_type,
		TransferMethod: transfer_method,
		RemoteURL:      url,
		Filename:       filename,
		Extension:      extension,
		MimeType:       mime_type,
		Size:           int(file_size),
		StorageKey:     "",
	}
}

func build_from_local_file(
	mapping map[string]any,
	tenant_id string,
	transfer_method file.FileTransferMethod,
) *file.File {
	var upload_file_id string
	if _, ok := mapping["upload_file_id"]; ok {
		if _, ok := mapping["upload_file_id"].(string); ok {
			upload_file_id = mapping["upload_file_id"].(string)
		}
	}
	if upload_file_id == "" {
		panic(exceptions.NewValueError("Invalid upload file id"))
	}
	// check if upload_file_id is a valid uuid

	_, err := uuid.FromString(upload_file_id)
	if err != nil {
		mlog.Errorf("Invalid upload file id format")
		panic(exceptions.NewValueError("Invalid upload file id format"))
	}
	stmt := new(models.UploadFile)
	err = dbengine.Instance().DB.Model(&models.UploadFile{}).Where("id = ? and tenant_id = ?", upload_file_id, tenant_id).First(stmt).Error
	if err != nil {
		panic(exceptions.NewValueError("Invalid upload file"))
	}
	file_type := file.FileType("custom")
	if _, ok := mapping["type"]; ok {
		if _, ok := mapping["type"].(string); ok {
			file_type = file.FileType(mapping["type"].(string))
		}
	}
	file_type = standardize_file_type(file_type, "."+stmt.Extension, stmt.MimeType)
	var id string
	if _, ok := mapping["id"]; ok {
		if _, ok := mapping["id"].(string); ok {
			id = mapping["id"].(string)
		}
	}
	return &file.File{
		ID:             id,
		TenantID:       tenant_id,
		Type:           file_type,
		TransferMethod: transfer_method,
		RemoteURL:      stmt.SourceURL,
		Filename:       stmt.Name,
		Extension:      "." + stmt.Extension,
		MimeType:       stmt.MimeType,
		RelatedID:      upload_file_id,
		Size:           stmt.Size,
		StorageKey:     stmt.Key,
	}
}

func is_file_valid_with_config(
	input_file_type file.FileType,
	file_extension string,
	file_transfer_method file.FileTransferMethod,
	config *file.FileUploadConfig,
) bool {
	if len(config.AllowedFileTypes) > 0 &&
		!slices.Contains(config.AllowedFileTypes, input_file_type) &&
		input_file_type != file.File_CUSTOM {
		return false
	}
	if input_file_type == file.File_CUSTOM &&
		len(config.AllowedFileExtensions) > 0 &&
		!slices.Contains(config.AllowedFileExtensions, file_extension) {
		return false
	}
	if input_file_type == file.File_IMAGE && config.ImageConfig != nil {
		if len(config.ImageConfig.TransferMethods) > 0 && !slices.Contains(config.ImageConfig.TransferMethods, file_transfer_method) {
			return false
		}
	}
	return true
}

func BuildFromMapping(
	mapping map[string]any,
	tenant_id string,
	config *file.FileUploadConfig,
) *file.File {
	var transfer_method_str string
	if _, ok := mapping["transfer_method"]; ok {
		if _, ok := mapping["transfer_method"].(string); ok {
			transfer_method_str = mapping["transfer_method"].(string)
		}
	}
	transfer_method := file.FileTransferMethod(transfer_method_str)
	build_functions := map[file.FileTransferMethod]func(map[string]any, string, file.FileTransferMethod) *file.File{
		file.FileTransferMethod_LOCAL_FILE: build_from_local_file,
		file.FileTransferMethod_REMOTE_URL: build_from_remote_url,
		file.FileTransferMethod_TOOL_FILE:  build_from_tool_file,
	}
	if _, ok := build_functions[transfer_method]; !ok {
		panic(exceptions.NewValueError(fmt.Sprintf("Invalid file transfer method: %v", transfer_method)))
	}

	f := build_functions[transfer_method](
		mapping,
		tenant_id,
		transfer_method,
	)

	input_file_type := file.File_CUSTOM
	if _, ok := mapping["type"]; ok {
		if _, ok := mapping["type"].(string); ok {
			input_file_type = file.FileType(mapping["type"].(string))
		}
	}
	if config != nil && !is_file_valid_with_config(
		input_file_type,
		f.Extension,
		f.TransferMethod,
		config,
	) {
		panic(exceptions.NewValueError(fmt.Sprintf("File validation failed for file: %s", f.Filename)))
	}
	return f
}

func BuildFromMappings(
	mappings []map[string]any,
	tenant_id string,
	config *file.FileUploadConfig,
) []*file.File {
	var files []*file.File
	var image_file_num int
	for _, mapping := range mappings {
		f := BuildFromMapping(
			mapping,
			tenant_id,
			config,
		)
		if files == nil {
			files = make([]*file.File, 0)
		}
		files = append(files, f)
		if f.Type == file.File_IMAGE {
			image_file_num += 1
		}
	}

	if config != nil && // If image config is set.
		config.ImageConfig != nil && image_file_num > config.ImageConfig.NumberLimits { // And the number of image files exceeds the maximum limit
		panic(exceptions.NewValueError(fmt.Sprintf("Number of image files exceeds the maximum limit %d", config.ImageConfig.NumberLimits)))
	}
	if config != nil && config.NumberLimits > 0 && len(files) > config.NumberLimits {
		panic(exceptions.NewValueError(fmt.Sprintf("Number of files exceeds the maximum limit %d", config.NumberLimits)))
	}
	return files
}
