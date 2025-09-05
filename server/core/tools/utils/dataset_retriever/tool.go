package datasetretriever

import "mlib.com/gofy/server/core/tools/base"

type DatasetRetrieverToolor interface {
	Run(query string) string
}
type DatasetRetrieverBaseTool struct {
    name: str = "dataset"
    description: str = "use this to retrieve a dataset. "
    tenant_id: str
    top_k: int = 2
    score_threshold: Optional[float] = None
    hit_callbacks: list[DatasetIndexToolCallbackHandler] = []
    return_resource: bool
    retriever_from: str
    model_config = ConfigDict(arbitrary_types_allowed=True)
}
    @abstractmethod
    def _run(self, query: str) -> str:
        """Use the tool.

        Add run_manager: Optional[CallbackManagerForToolRun] = None
        to child implementations to enable tracing,
        """

type DatasetRetrieverTool struct {
	base.Tool
}
