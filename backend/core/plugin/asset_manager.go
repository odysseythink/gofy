package plugin

import (
	"fmt"
	"path/filepath"
	"strings"

	"mlib.com/gofy/server/storage"
	"mlib.com/mlog"
)

// AssetManager handles plugin static assets (icons, images).
type AssetManager struct{}

// NewAssetManager creates an asset manager.
func NewAssetManager() *AssetManager {
	return &AssetManager{}
}

// SaveAsset stores a plugin asset.
func (am *AssetManager) SaveAsset(uniqueIdentifier, assetName string, data []byte) (string, error) {
	key := am.assetKey(uniqueIdentifier, assetName)
	if err := storage.Save(key, data); err != nil {
		return "", fmt.Errorf("failed to save asset: %w", err)
	}
	return key, nil
}

// GetAsset retrieves a plugin asset.
func (am *AssetManager) GetAsset(uniqueIdentifier, assetName string) ([]byte, error) {
	key := am.assetKey(uniqueIdentifier, assetName)
	return storage.LoadOnce(key)
}

// GetIcon retrieves the plugin's icon.
func (am *AssetManager) GetIcon(uniqueIdentifier string) ([]byte, string, error) {
	// Try common icon formats
	for _, ext := range []string{"svg", "png", "jpg", "ico"} {
		key := am.assetKey(uniqueIdentifier, "icon."+ext)
		data, err := storage.LoadOnce(key)
		if err == nil {
			return data, mimeTypeForExt(ext), nil
		}
	}
	return nil, "", fmt.Errorf("icon not found for plugin %s", uniqueIdentifier)
}

// DeleteAssets removes all assets for a plugin.
func (am *AssetManager) DeleteAssets(uniqueIdentifier string) error {
	mlog.Infof("cleaning assets for plugin %s", uniqueIdentifier)
	// TODO: Implement directory-level deletion in storage
	return nil
}

// ListAssets lists all assets for a plugin.
func (am *AssetManager) ListAssets(uniqueIdentifier string) []string {
	// TODO: Implement listing in storage backend
	return nil
}

func (am *AssetManager) assetKey(uniqueIdentifier, assetName string) string {
	// Sanitize to prevent path traversal
	assetName = filepath.Base(assetName)
	return fmt.Sprintf("plugins/assets/%s/%s", uniqueIdentifier, assetName)
}

func mimeTypeForExt(ext string) string {
	switch strings.ToLower(ext) {
	case "svg":
		return "image/svg+xml"
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "ico":
		return "image/x-icon"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
