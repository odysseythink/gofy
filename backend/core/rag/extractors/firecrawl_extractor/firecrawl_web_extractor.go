package firecrawlextractor

import (
	ragentities "mlib.com/gofy/server/entities/rag"
	"mlib.com/gofy/server/services"
)

type FirecrawlWebExtractor struct {
	_url              string
	job_id            string
	tenant_id         string
	mode              string
	only_main_content bool
}

func New(
	url string,
	job_id string,
	tenant_id string,
	mode string,
	only_main_content bool, /* = true*/
) *FirecrawlWebExtractor {
	if mode == "" {
		mode = "crawl"
	}
	return &FirecrawlWebExtractor{
		_url:              url,
		job_id:            job_id,
		tenant_id:         tenant_id,
		mode:              mode,
		only_main_content: only_main_content,
	}
}

func (extractor *FirecrawlWebExtractor) Extract() []*ragentities.Document {
	documents := []*ragentities.Document{}
	if extractor.mode == "crawl" {
		crawl_data := (&services.WebsiteService{}).GetCrawlURLData(extractor.job_id, "firecrawl", extractor._url, extractor.tenant_id)
		if len(crawl_data) == 0 {
			return nil
		}
		markdown := ""
		if _, ok := crawl_data["markdown"]; ok {
			if _, ok := crawl_data["markdown"].(string); ok {
				markdown = crawl_data["markdown"].(string)
			}
		}
		document := &ragentities.Document{
			PageContent: markdown,
			Metadata: map[string]any{
				"source_url":  crawl_data["source_url"],
				"description": crawl_data["description"],
				"title":       crawl_data["title"],
			},
			Provider: "dify",
		}
		documents = append(documents, document)
	} else if extractor.mode == "scrape" {
		scrape_data := (&services.WebsiteService{}).GetScrapeURLData("firecrawl", extractor._url, extractor.tenant_id, extractor.only_main_content)
		markdown := ""
		if _, ok := scrape_data["markdown"]; ok {
			if _, ok := scrape_data["markdown"].(string); ok {
				markdown = scrape_data["markdown"].(string)
			}
		}
		document := &ragentities.Document{
			PageContent: markdown,
			Metadata: map[string]any{
				"source_url":  scrape_data["source_url"],
				"description": scrape_data["description"],
				"title":       scrape_data["title"],
			},
		}
		documents = append(documents, document)
	}
	return documents
}
