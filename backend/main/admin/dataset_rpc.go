package main

import (
	"context"
	"net/http"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
	"google.golang.org/grpc/peer"
	enumtypes "mlib.com/gofy/server/enum_types/rag/vector"
	pbexceptions "mlib.com/gofy/server/proto/exceptions"
	"mlib.com/gofy/server/proto/pbapi"
)

var (
	RetrievalMethod_SEMANTIC_SEARCH  = "semantic_search"
	RetrievalMethod_FULL_TEXT_SEARCH = "full_text_search"
	RetrievalMethod_HYBRID_SEARCH    = "hybrid_search"
)

func (s *AdminService) GetDatasetRetrievalSetting(ctx context.Context, in *pbapi.GetDatasetRetrievalSettingRequest) (out *pbapi.GetDatasetRetrievalSettingReply, err error) {
	p, _ := peer.FromContext(ctx)
	mlog.Infof("remote[%s] admin.GetDatasetRetrievalSetting call:%#v", p.Addr.String(), in)

	out = &pbapi.GetDatasetRetrievalSettingReply{}
	switch enumtypes.VectorType(confy.Get[string]("vector_store")) {
	case enumtypes.Vector_MILVUS,
		enumtypes.Vector_RELYT,
		enumtypes.Vector_PGVECTOR,
		enumtypes.Vector_TIDB_VECTOR,
		enumtypes.Vector_CHROMA,
		enumtypes.Vector_TENCENT:
		out.RetrievalMethod = []string{RetrievalMethod_SEMANTIC_SEARCH}
	case enumtypes.Vector_QDRANT,
		enumtypes.Vector_WEAVIATE,
		enumtypes.Vector_OPENSEARCH,
		enumtypes.Vector_ANALYTICDB,
		enumtypes.Vector_MYSCALE,
		enumtypes.Vector_ORACLE,
		enumtypes.Vector_ELASTICSEARCH:
		out.RetrievalMethod = []string{RetrievalMethod_SEMANTIC_SEARCH, RetrievalMethod_FULL_TEXT_SEARCH, RetrievalMethod_HYBRID_SEARCH}
	default:
		mlog.Errorf("Unsupported vector db type %s.", confy.Get[string]("vector_store"))
		out.Exp = &pbexceptions.HTTPException{
			Status:  http.StatusNoContent,
			Message: "Unsupported vector db type",
		}
		return
	}
	return
}
