package plugin

// from core.plugin.entities.plugin import PluginDeclaration, PluginInstallationSource

type PluginBundleDependencyGithub struct {
	RepoAddress string `json:"repo_address"`
	Repo        string `json:"repo"`
	Release     string `json:"release"`
	Packages    string `json:"packages"`
}
type PluginBundleDependencyMarketplace struct {
	Organization string `json:"organization"`
	Plugin       string `json:"plugin"`
	Version      string `json:"version"`
}
type PluginBundleDependencyPackage struct {
	UniqueIdentifier string            `json:"unique_identifier"`
	Manifest         PluginDeclaration `json:"manifest"`
}
