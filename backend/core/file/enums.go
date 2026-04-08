package file

type FileType string

const (
	File_IMAGE    FileType = "image"
	File_DOCUMENT FileType = "document"
	File_AUDIO    FileType = "audio"
	File_VIDEO    FileType = "video"
	File_CUSTOM   FileType = "custom"

	// @staticmethod
	// def value_of(value):
	//     for member in FileType:
	//         if member.value == value:
	//             return member
	//     raise ValueError(f"No matching enum found for value '{value}'")

)

type FileTransferMethod string

const (
	FileTransferMethod_REMOTE_URL FileTransferMethod = "remote_url"
	FileTransferMethod_LOCAL_FILE FileTransferMethod = "local_file"
	FileTransferMethod_TOOL_FILE  FileTransferMethod = "tool_file"

	// @staticmethod
	// def value_of(value):
	//     for member in FileTransferMethod:
	//         if member.value == value:
	//             return member
	//     raise ValueError(f"No matching enum found for value '{value}'")

)

type FileBelongsTo string

const (
	FileBelongsTo_USER      FileBelongsTo = "user"
	FileBelongsTo_ASSISTANT FileBelongsTo = "assistant"

	// @staticmethod
	// def value_of(value):
	//     for member in FileBelongsTo:
	//         if member.value == value:
	//             return member
	//     raise ValueError(f"No matching enum found for value '{value}'")

)

type FileAttribute string

const (
	FileAttribute_TYPE            FileAttribute = "type"
	FileAttribute_SIZE            FileAttribute = "size"
	FileAttribute_NAME            FileAttribute = "name"
	FileAttribute_MIME_TYPE       FileAttribute = "mime_type"
	FileAttribute_TRANSFER_METHOD FileAttribute = "transfer_method"
	FileAttribute_URL             FileAttribute = "url"
	FileAttribute_EXTENSION       FileAttribute = "extension"
)

func FileAttributeContains(attr string) bool {
	return FileAttribute(attr) == FileAttribute_TYPE ||
		FileAttribute(attr) == FileAttribute_SIZE ||
		FileAttribute(attr) == FileAttribute_NAME ||
		FileAttribute(attr) == FileAttribute_MIME_TYPE ||
		FileAttribute(attr) == FileAttribute_TRANSFER_METHOD ||
		FileAttribute(attr) == FileAttribute_URL ||
		FileAttribute(attr) == FileAttribute_EXTENSION
}

type ArrayFileAttribute string

const (
	ArrayFileAttribute_LENGTH ArrayFileAttribute = "length"
)
