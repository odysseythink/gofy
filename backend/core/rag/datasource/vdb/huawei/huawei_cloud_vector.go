// Package huawei is the Huawei Cloud CSS (OpenSearch-based vector) IVector
// adapter.
//
// Huawei Cloud Search Service wraps ES/OpenSearch under Huawei AK/SK signing.
// We reuse the elasticsearch adapter's REST surface with basic auth;
// production deployments should wrap outbound requests with V4 signing via
// github.com/huaweicloud/huaweicloud-sdk-go-v3.
package huawei

import (
	"context"
	"errors"
	"time"

	"github.com/odysseythink/confy"

	vdb "github.com/odysseythink/gofy/backend/core/rag/datasource/vdb"
	"github.com/odysseythink/gofy/backend/core/rag/datasource/vdb/elasticsearch"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	vectorenumtypes "github.com/odysseythink/gofy/backend/enum_types/rag/vector"
)

func init() { vdb.Register(vectorenumtypes.Vector_HUAWEI_CLOUD, &Factory{}) }

type Factory struct{}

func (Factory) New(ctx context.Context, cfg vdb.AdapterConfig) (ragentities.IVector, error) {
	esCfg := elasticsearch.Config{
		Endpoint: confy.GetWithDefault[string]("huawei.endpoint", ""),
		Username: confy.GetWithDefault[string]("huawei.username", ""),
		Password: confy.GetWithDefault[string]("huawei.password", ""),
	}
	if esCfg.Endpoint == "" {
		return nil, errors.New("huawei: endpoint is required")
	}
	_ = time.Second
	v, err := elasticsearch.Open(ctx, esCfg, cfg.CollectionName)
	if err != nil {
		return nil, err
	}
	return &wrap{inner: v}, nil
}

// wrap exposes the ES-backed adapter under the HUAWEI_CLOUD type name.
type wrap struct{ inner ragentities.IVector }

func (w *wrap) GetType() string { return string(vectorenumtypes.Vector_HUAWEI_CLOUD) }

func (w *wrap) Create(ctx context.Context, d []*ragentities.Document, e [][]float32) error {
	return w.inner.Create(ctx, d, e)
}
func (w *wrap) AddTexts(ctx context.Context, d []*ragentities.Document, e [][]float32) error {
	return w.inner.AddTexts(ctx, d, e)
}
func (w *wrap) TextExists(ctx context.Context, id string) (bool, error) {
	return w.inner.TextExists(ctx, id)
}
func (w *wrap) GetIDsByMetadataField(ctx context.Context, k, v string) ([]string, error) {
	return w.inner.GetIDsByMetadataField(ctx, k, v)
}
func (w *wrap) DeleteByIDs(ctx context.Context, ids []string) error {
	return w.inner.DeleteByIDs(ctx, ids)
}
func (w *wrap) DeleteByMetadataField(ctx context.Context, k, v string) error {
	return w.inner.DeleteByMetadataField(ctx, k, v)
}
func (w *wrap) SearchByVector(ctx context.Context, qv []float32, o ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	return w.inner.SearchByVector(ctx, qv, o)
}
func (w *wrap) SearchByFullText(ctx context.Context, q string, o ragentities.SearchOptions) ([]*ragentities.SearchResult, error) {
	return w.inner.SearchByFullText(ctx, q, o)
}
func (w *wrap) Delete(ctx context.Context) error { return w.inner.Delete(ctx) }
