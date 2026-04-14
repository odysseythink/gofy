package datasource

import (
	"fmt"
	"io"
	"sync"

	"github.com/odysseythink/mlog"
)

// DatasourceType defines the type of data source.
type DatasourceType string

const (
	DatasourceTypeLocalFile      DatasourceType = "local_file"
	DatasourceTypeWebsiteCrawl   DatasourceType = "website_crawl"
	DatasourceTypeOnlineDocument DatasourceType = "online_document"
	DatasourceTypeOnlineDrive    DatasourceType = "online_drive"
)

// DatasourceDocument represents an extracted document from a datasource.
type DatasourceDocument struct {
	Title    string
	Content  string
	Metadata map[string]any
	MimeType string
	Size     int64
}

// DatasourceProvider is the interface for data source connectors.
type DatasourceProvider interface {
	Type() DatasourceType
	ValidateCredentials(credentials map[string]any) error
	FetchDocuments(config map[string]any, credentials map[string]any) ([]*DatasourceDocument, error)
	FetchDocumentContent(config map[string]any, credentials map[string]any, documentID string) (io.ReadCloser, error)
}

// DatasourceManager manages data source providers.
type DatasourceManager struct {
	mu        sync.RWMutex
	providers map[DatasourceType]DatasourceProvider
}

var globalDatasourceManager *DatasourceManager

// InitDatasourceManager initializes the global manager.
func InitDatasourceManager() {
	globalDatasourceManager = &DatasourceManager{
		providers: make(map[DatasourceType]DatasourceProvider),
	}
	// Register built-in providers
	globalDatasourceManager.Register(NewLocalFileProvider())
	globalDatasourceManager.Register(NewWebsiteCrawlProvider())
	mlog.Info("datasource manager initialized")
}

// GetDatasourceManager returns the global manager.
func GetDatasourceManager() *DatasourceManager {
	return globalDatasourceManager
}

// Register adds a provider.
func (dm *DatasourceManager) Register(provider DatasourceProvider) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.providers[provider.Type()] = provider
}

// GetProvider returns a provider by type.
func (dm *DatasourceManager) GetProvider(dsType DatasourceType) (DatasourceProvider, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	p, ok := dm.providers[dsType]
	if !ok {
		return nil, fmt.Errorf("datasource provider not found: %s", dsType)
	}
	return p, nil
}
