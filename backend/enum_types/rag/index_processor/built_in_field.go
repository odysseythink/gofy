package indexprocessor

type BuiltInFieldType string

const (
	BuiltInField_document_name    BuiltInFieldType = "document_name"
	BuiltInField_uploader         BuiltInFieldType = "uploader"
	BuiltInField_upload_date      BuiltInFieldType = "upload_date"
	BuiltInField_last_update_date BuiltInFieldType = "last_update_date"
	BuiltInField_source           BuiltInFieldType = "source"
)

type MetadataDataSourceType string

const (
	MetadataDataSource_upload_file   MetadataDataSourceType = "file_upload"
	MetadataDataSource_website_crawl MetadataDataSourceType = "website"
	MetadataDataSource_notion_import MetadataDataSourceType = "notion"
)
