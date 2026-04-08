package plugin

import (
	"encoding/json"
	"sync"
	"time"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
)

// DeclarationCache caches plugin declarations in memory.
type DeclarationCache struct {
	mu      sync.RWMutex
	entries map[string]*cachedDeclaration
	ttl     time.Duration
}

type cachedDeclaration struct {
	declaration *models.PluginDeclaration
	parsed      map[string]any
	loadedAt    time.Time
}

// NewDeclarationCache creates a declaration cache.
func NewDeclarationCache(ttl time.Duration) *DeclarationCache {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &DeclarationCache{
		entries: make(map[string]*cachedDeclaration),
		ttl:     ttl,
	}
}

// Get returns a cached declaration, loading from DB if needed.
func (dc *DeclarationCache) Get(uniqueIdentifier string) (*models.PluginDeclaration, map[string]any) {
	dc.mu.RLock()
	if entry, ok := dc.entries[uniqueIdentifier]; ok {
		if time.Since(entry.loadedAt) < dc.ttl {
			dc.mu.RUnlock()
			return entry.declaration, entry.parsed
		}
	}
	dc.mu.RUnlock()

	// Load from DB
	return dc.load(uniqueIdentifier)
}

// Invalidate removes a cached entry.
func (dc *DeclarationCache) Invalidate(uniqueIdentifier string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	delete(dc.entries, uniqueIdentifier)
}

// InvalidateAll clears the entire cache.
func (dc *DeclarationCache) InvalidateAll() {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.entries = make(map[string]*cachedDeclaration)
}

// Set manually sets a cache entry.
func (dc *DeclarationCache) Set(uniqueIdentifier string, decl *models.PluginDeclaration) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	parsed := make(map[string]any)
	if decl != nil {
		json.Unmarshal([]byte(decl.Declaration), &parsed)
	}
	dc.entries[uniqueIdentifier] = &cachedDeclaration{
		declaration: decl,
		parsed:      parsed,
		loadedAt:    time.Now(),
	}
}

// GetToolProviders extracts all tool providers from cached declarations.
func (dc *DeclarationCache) GetToolProviders() []map[string]any {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	var providers []map[string]any
	for id, entry := range dc.entries {
		if tools, ok := entry.parsed["tools"].([]any); ok && len(tools) > 0 {
			providers = append(providers, map[string]any{
				"unique_identifier": id,
				"tools":             tools,
			})
		}
	}
	return providers
}

// GetModelProviders extracts all model providers from cached declarations.
func (dc *DeclarationCache) GetModelProviders() []map[string]any {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	var providers []map[string]any
	for id, entry := range dc.entries {
		if models, ok := entry.parsed["models"].([]any); ok && len(models) > 0 {
			providers = append(providers, map[string]any{
				"unique_identifier": id,
				"models":            models,
			})
		}
	}
	return providers
}

// Size returns the number of cached entries.
func (dc *DeclarationCache) Size() int {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return len(dc.entries)
}

func (dc *DeclarationCache) load(uniqueIdentifier string) (*models.PluginDeclaration, map[string]any) {
	var decl models.PluginDeclaration
	if err := dbengine.Instance().DB.Where("plugin_unique_identifier = ?", uniqueIdentifier).First(&decl).Error; err != nil {
		return nil, nil
	}

	parsed := make(map[string]any)
	json.Unmarshal([]byte(decl.Declaration), &parsed)

	dc.mu.Lock()
	dc.entries[uniqueIdentifier] = &cachedDeclaration{
		declaration: &decl,
		parsed:      parsed,
		loadedAt:    time.Now(),
	}
	dc.mu.Unlock()

	return &decl, parsed
}

