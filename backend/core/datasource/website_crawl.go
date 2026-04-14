package datasource

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/odysseythink/mlog"
)

type WebsiteCrawlProvider struct {
	client *http.Client
}

func NewWebsiteCrawlProvider() *WebsiteCrawlProvider {
	return &WebsiteCrawlProvider{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *WebsiteCrawlProvider) Type() DatasourceType {
	return DatasourceTypeWebsiteCrawl
}

func (p *WebsiteCrawlProvider) ValidateCredentials(credentials map[string]any) error {
	return nil
}

func (p *WebsiteCrawlProvider) FetchDocuments(config map[string]any, credentials map[string]any) ([]*DatasourceDocument, error) {
	url, ok := config["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url is required")
	}

	// Check for external crawl service (Firecrawl/Jina)
	crawlMode, _ := config["mode"].(string)

	switch crawlMode {
	case "firecrawl":
		return p.crawlWithFirecrawl(url, credentials)
	case "jina":
		return p.crawlWithJina(url, credentials)
	default:
		return p.directCrawl(url)
	}
}

func (p *WebsiteCrawlProvider) FetchDocumentContent(config map[string]any, credentials map[string]any, documentID string) (io.ReadCloser, error) {
	url, ok := config["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url is required")
	}
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (p *WebsiteCrawlProvider) directCrawl(url string) ([]*DatasourceDocument, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc := &DatasourceDocument{
		Title:    url,
		Content:  string(body),
		MimeType: resp.Header.Get("Content-Type"),
		Size:     int64(len(body)),
		Metadata: map[string]any{
			"url":         url,
			"status_code": resp.StatusCode,
		},
	}

	return []*DatasourceDocument{doc}, nil
}

func (p *WebsiteCrawlProvider) crawlWithFirecrawl(url string, credentials map[string]any) ([]*DatasourceDocument, error) {
	mlog.Infof("firecrawl mode for %s (not yet implemented)", url)
	// TODO: Implement Firecrawl API integration
	return p.directCrawl(url)
}

func (p *WebsiteCrawlProvider) crawlWithJina(url string, credentials map[string]any) ([]*DatasourceDocument, error) {
	// Jina Reader API: prepend https://r.jina.ai/ to the URL
	jinaURL := "https://r.jina.ai/" + url
	resp, err := p.client.Get(jinaURL)
	if err != nil {
		return nil, fmt.Errorf("jina reader failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc := &DatasourceDocument{
		Title:    url,
		Content:  string(body),
		MimeType: "text/markdown",
		Size:     int64(len(body)),
		Metadata: map[string]any{"url": url, "source": "jina_reader"},
	}

	return []*DatasourceDocument{doc}, nil
}
