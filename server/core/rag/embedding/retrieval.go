package embedding


type RetrievalChildChunk struct{
    ID  string `json:"id"`
    Content  string `json:"content"`
    Score float64 `json:"score"`
    Position int `json:"position"`
}

type RetrievalSegments struct{
    model_config = {"arbitrary_types_allowed": True}`json:"id"`
    segment: DocumentSegment`json:"id"`
    child_chunks: Optional[list[RetrievalChildChunk]] = None`json:"id"`
    Score float64 `json:"score"`
}

func NewRetrievalSegments() *RetrievalSegments{
    return &RetrievalSegments{

    }
}