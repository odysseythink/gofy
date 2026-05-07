package extractprocessor

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	baseextractor "github.com/odysseythink/gofy/backend/core/rag/extractors/base_extractor"
	csvextractor "github.com/odysseythink/gofy/backend/core/rag/extractors/csv_extractor"
	excelextractor "github.com/odysseythink/gofy/backend/core/rag/extractors/excel_extractor"
	firecrawlextractor "github.com/odysseythink/gofy/backend/core/rag/extractors/firecrawl_extractor"
	htmlextractor "github.com/odysseythink/gofy/backend/core/rag/extractors/html_extractor"
	jinareaderextractor "github.com/odysseythink/gofy/backend/core/rag/extractors/jina_reader_extractor"
	markdownextractor "github.com/odysseythink/gofy/backend/core/rag/extractors/markdown_extractor"
	notionextractor "github.com/odysseythink/gofy/backend/core/rag/extractors/notion_extractor"
	pdfextractor "github.com/odysseythink/gofy/backend/core/rag/extractors/pdf_extractor"
	textextractor "github.com/odysseythink/gofy/backend/core/rag/extractors/text_extractor"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	extractorentities "github.com/odysseythink/gofy/backend/entities/rag/extractor"
	extractorenumtypes "github.com/odysseythink/gofy/backend/enum_types/rag/extractor"
	"github.com/odysseythink/gofy/backend/utils"
)

type ExtractProcessor struct {
}

func (processor *ExtractProcessor) Extract(
	extract_setting *extractorentities.ExtractSetting, is_automatic bool, file_path string,
) []*ragentities.Document {
	var extractor baseextractor.Extractor
	if extract_setting.DatasourceType == extractorenumtypes.Datasource_FILE {
		temp_dir := os.TempDir()
		if file_path == "" {
			if extract_setting.UploadFile == nil {
				panic(exceptions.NewValueError("upload_file is required"))
			}

			upload_file := extract_setting.UploadFile
			suffix := filepath.Ext(upload_file.Key)

			// FIXME mypy: Cannot determine type of 'tempfile._get_candidate_names' better not use it here
			file_path := filepath.Join(temp_dir, strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(time.Now().Format("2006-01-02 15:04:05.000000"), "-", ""), " ", ""), ":", ""), ".", "")+suffix) // type: ignore
			utils.DownloadFileFromRemoteURL(upload_file.Key, file_path)
		}
		file_extension := strings.ToLower(filepath.Ext(file_path))
		// etl_type := confy.GetWithDefault[string]("ETL_TYPE", "gofy")

		// unstructured_api_url := confy.Get[string]("unstructured_api_url")
		// unstructured_api_key := confy.Get[string]("unstructured_api_key")

		if slices.Contains([]string{".xlsx", ".xls"}, file_extension) {
			extractor = excelextractor.New(file_path, "", false)
		} else if file_extension == ".pdf" {
			extractor = pdfextractor.New(file_path)
		} else if slices.Contains([]string{".md", ".markdown", ".mdx"}, file_extension) {
			extractor = markdownextractor.New(file_path, false, false, "", true)
		} else if slices.Contains([]string{".htm", ".html"}, file_extension) {
			extractor = htmlextractor.New(file_path)
		} else if file_extension == ".csv" {
			extractor = csvextractor.New(file_path, "", true, "", nil)
		} else {
			// txt
			extractor = textextractor.New(file_path, "", true)
		}

		return extractor.Extract()
	} else if extract_setting.DatasourceType == extractorenumtypes.Datasource_NOTION {
		if extract_setting.NotionInfo == nil {
			panic(exceptions.NewValueError("notion_info is required"))
		}
		extractor = notionextractor.New(
			extract_setting.NotionInfo.NotionWorkspaceID,
			extract_setting.NotionInfo.NotionObjID,
			extract_setting.NotionInfo.NotionPageType,
			extract_setting.NotionInfo.TenantID,
			extract_setting.NotionInfo.Document,
			"",
		)
		return extractor.Extract()
	} else if extract_setting.DatasourceType == extractorenumtypes.Datasource_WEBSITE {
		if extract_setting.WebsiteInfo == nil {
			panic(exceptions.NewValueError("website_info is required"))
		}
		if extract_setting.WebsiteInfo.Provider == "firecrawl" {
			extractor = firecrawlextractor.New(
				extract_setting.WebsiteInfo.URL,
				extract_setting.WebsiteInfo.JobID,
				extract_setting.WebsiteInfo.TenantID,
				extract_setting.WebsiteInfo.Mode,
				extract_setting.WebsiteInfo.OnlyMainContent,
			)
			return extractor.Extract()
		} else if extract_setting.WebsiteInfo.Provider == "jinareader" {
			extractor = jinareaderextractor.New(
				extract_setting.WebsiteInfo.URL,
				extract_setting.WebsiteInfo.JobID,
				extract_setting.WebsiteInfo.TenantID,
				extract_setting.WebsiteInfo.Mode,
				extract_setting.WebsiteInfo.OnlyMainContent,
			)
			return extractor.Extract()
		} else {
			panic(exceptions.NewValueError("Unsupported website provider: " + extract_setting.WebsiteInfo.Provider))
		}
	} else {
		panic(exceptions.NewValueError("Unsupported datasource type: " + string(extract_setting.DatasourceType)))
	}
}
