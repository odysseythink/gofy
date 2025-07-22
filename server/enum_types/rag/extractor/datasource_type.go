package extractor

type DatasourceType string

const (
	Datasource_FILE    DatasourceType = "upload_file"
	Datasource_NOTION  DatasourceType = "notion_import"
	Datasource_WEBSITE DatasourceType = "website_crawl"
)
