package jinareaderextractor

import (
	ragentities "mlib.com/gofy/server/entities/rag"
	"mlib.com/gofy/server/services"
)

type JinaReaderWebExtractor struct {
	_url              string
	job_id            string
	tenant_id         string
	mode              string
	only_main_content bool
}

func New(
	url string, job_id string, tenant_id string, mode string, only_main_content bool,
) *JinaReaderWebExtractor {
	if mode == "" {
		mode = "crawl"
	}
	return &JinaReaderWebExtractor{
		_url:              url,
		job_id:            job_id,
		tenant_id:         tenant_id,
		mode:              mode,
		only_main_content: only_main_content,
	}
}

func (extractor *JinaReaderWebExtractor) Extract() []*ragentities.Document {
	documents := []*ragentities.Document{}
	if extractor.mode == "crawl" {
		crawl_data := (&services.WebsiteService{}).GetCrawlURLData(extractor.job_id, "jinareader", extractor._url, extractor.tenant_id)
		if len(crawl_data) == 0 {
			return nil
		}
		content := ""
		if _, ok := crawl_data["content"]; ok {
			if _, ok := crawl_data["content"].(string); ok {
				content = crawl_data["content"].(string)
			}
		}
		document := &ragentities.Document{
			PageContent: content,
			Metadata: map[string]any{
				"source_url":  crawl_data["url"],
				"description": crawl_data["description"],
				"title":       crawl_data["title"],
			},
			Provider: "gofy",
		}
		documents = append(documents, document)
	}
	return documents
}
