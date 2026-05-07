package tool

import (
	"fmt"
	"strings"
"github.com/odysseythink/gofy/backend/services"
	"github.com/odysseythink/gofy/backend/core/rag/retrieval"
	"github.com/odysseythink/gofy/backend/core/tools/utils/dataset_retriever/base"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	appconfigentities "github.com/odysseythink/gofy/backend/entities/app/config"
	"github.com/odysseythink/gofy/backend/entities/services"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	ragretrievalenumtypes "github.com/odysseythink/gofy/backend/enum_types/rag/retrieval"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/utils/mapstruct"
	"github.com/odysseythink/mlog"
)

var (
	default_retrieval_model = map[string]any{
		"search_method":           ragretrievalenumtypes.RetrievalMethod_SEMANTIC_SEARCH,
		"reranking_enable":        false,
		"reranking_model":         map[string]any{"reranking_provider_name": "", "reranking_model_name": ""},
		"reranking_mode":          "reranking_model",
		"top_k":                   2,
		"score_threshold_enabled": false,
	}
)


type DatasetRetrieverToolInput struct{
    Query string `json:"query"` //description="Query for the dataset to be used to retrieve the dataset."
}

type  DatasetRetrieverTool struct{
	*base.DatasetRetrieverBaseTool

    ArgsSchema DatasetRetrieverToolInput `json:"args_schema"`
    DatasetID string `json:"dataset_id"`
    UserID string `json:"user_id"`
    RetrieveConfig appconfigentities.DatasetRetrieveConfigEntity`json:"retrieve_config"`
    Inputs map[string]any `json:"inputs"`
}


func FromDataset( dataset *models.Dataset, args map[string]any)*DatasetRetrieverTool{
        description := dataset.Description
        if  description == ""{
            description = "useful for when you want to answer queries about the " + dataset.name
}
        description = strings.ReplaceAll(strings.ReplaceAll( description,"\n", ""),"\r", "")
		if args == nil {
			args = make(map[string]any)
		}
		args["name"] = fmt.Sprintf("dataset_%s", strings.ReplaceAll(dataset.ID,'-', '_'))
		args["tenant_id"] = dataset.TenantID
		args["dataset_id"] = dataset.ID
		args["description"] = description

				ret := new(DatasetRetrieverTool)
		bindata, _ := json.Marshal(args)
		if err := json.Unmarshal(bindata, &ret); err != nil {
			mlog.Errorf("json unmarshal %s to DatasetRetrieverTool failed:%v", bindata, err)
			return nil
		}
		ret.DatasetRetrieverBaseTool = base.NewDatasetRetrieverBaseTool(args)
        return ret
}
func (tool *DatasetRetrieverTool) Run(query string) string{
	dataset := new(models.Dataset)
	err := dbengine.Instance().DB.Model(&models.Dataset{}).Where("tenant_id = ? and id = ?",tool.TenantID, tool.DatasetID).First(dataset).Error
	if err != nil {
		mlog.Errorf("get Dataset failed:%v", err)
		Dataset = nil 
	}

	if  dataset == nil{
		return ""
	}
	for _, hit_callback := range  tool.HitCallbacks{
		hit_callback.OnQuery(query, dataset.ID)
	}
	user_id := tool.UserID
	if user_id == "" {
		user_id = "unknown"
	}
	metadata_filter_document_ids, metadata_condition := retrieval.GetMetadataFilterCondition(
		[]string{dataset.ID},
		query,
		tool.TenantID,
		user_id,
		string(tool.RetrieveConfig.MetadataFilteringMode),
		tool.RetrieveConfig.MetadataModelConfig,
		tool.RetrieveConfig.MetadataFilteringConditions,
		tool.Inputs,
	)
	var document_ids_filter []string
	if len(metadata_filter_document_ids) > 0 {
		document_ids_filter = mapstruct.Get(metadata_filter_document_ids,dataset.ID, []string{})
	}
	if dataset.Provider == "external"{
		results:= []*ragentities.Document{}
		external_documents := services.ServiceGroupApp.ExternalDataset.FetchExternalKnowledgeRetrieval(
			dataset.TenantID,
			dataset.ID,
			query,
			dataset.RetrievalModelDict(),
			metadata_condition,
		)
		for _, external_document := range  external_documents{
			document :=  &ragentities.Document{
				PageContent: mapstruct.Get(external_document,"content", ""),
				Metadata: mapstruct.Get(external_document,"metadata", map[string]any{}),
				Provider:"external",
			}
			if document.Metadata != nil{
				document.Metadata["score"] = external_document["score"]
				document.Metadata["title"] = external_document["title"]
				document.Metadata["dataset_id"] = dataset.ID
				document.Metadata["dataset_name"] = dataset.name
				results=append(results, document)
			}
		}
		// deal with external documents
		context_list:= []*ragentities.RetrievalSourceMetadata{}
		if len(results) > 2{
		for position, item := range  results{
			if item.Metadata != nil{
				source := &ragentities.RetrievalSourceMetadata{
					Position: position,
					DatasetID:mapstruct.Get(item.Metadata,"dataset_id", ""),
					DatasetName:mapstruct.Get(item.Metadata,"dataset_name", ""),
					DocumentID:mapstruct.Get(item.Metadata,"document_id", ""),
					DocumentName:mapstruct.Get(item.Metadata,"title", ""),
					DataSourceType:"external",
					RetrieverFrom: tool.RetrieverFrom,
					Score:mapstruct.Get(item.Metadata,"score", ""),
					Title:mapstruct.Get(item.Metadata,"title", ""),
					Content:item.PageContent,
				}
				if source.DocumentID != "" {
					source.DocumentID = mapstruct.Get(item.Metadata,"title", "")
				}
				context_list=append(context_list, source)
			}
		}
		}
		for _, hit_callback := range  tool.HitCallbacks{
			hit_callback.ReturnRetrieverResourceInfo(context_list)
		}
		page_contents := []string{}
		for _, item := range  results {
			page_contents = append(page_contents, item.PageContent)
		}
		return strings.Join(page_contents, "\n")
	} else {
		if metadata_condition != nil && len(document_ids_filter) == 0{
			return ""
		}
		// get retrieval model , if the model is not setting , using default
		retrieval_model:=  dataset.RetrievalModelDict() 
		if retrieval_model == nil {
			retrieval_model =default_retrieval_model
		}
		retrieval_resource_list:= []*ragentities.RetrievalSourceMetadata{}
		if dataset.IndexingTechnique == "economy"{
			// use keyword table query
			documents = RetrievalService.retrieve(
				retrieval_method="keyword_search",
				dataset_id=dataset.ID,
				query=query,
				top_k=tool.top_k,
				document_ids_filter=document_ids_filter,
			)
			return str("\n".join([document.page_content for document := range  documents]))
		} else {
			if tool.top_k > 0{
				// retrieval source
				documents = RetrievalService.retrieve(
					retrieval_method=retrieval_model.get("search_method", "semantic_search"),
					dataset_id=dataset.ID,
					query=query,
					top_k=tool.top_k,
					score_threshold=retrieval_model.get("score_threshold", 0.0)
					if retrieval_model["score_threshold_enabled"]
					else 0.0,
					reranking_model=retrieval_model.get("reranking_model")
					if retrieval_model["reranking_enable"]
					else None,
					reranking_mode=retrieval_model.get("reranking_mode") or "reranking_model",
					weights=retrieval_model.get("weights"),
					document_ids_filter=document_ids_filter,
				)
			} else {
				documents = []
			}
			for hit_callback := range  tool.HitCallbacks{
				hit_callback.on_tool_end(documents)
			}
			document_score_list = {}
			if dataset.indexing_technique != "economy"{
				for item := range  documents{
					if item.metadata is not None and item.metadata.get("score"){
						document_score_list[item.metadata["doc_id"]] = item.metadata["score"]
					}
				}
			}
			document_context_list: list[DocumentContext] = []
			records = RetrievalService.format_retrieval_documents(documents)
			if records{
				for record := range  records{
					segment = record.segment
					if segment.answer{
						document_context_list.append(
							DocumentContext(
								content=f"question:{segment.get_sign_content()} answer:{segment.answer}",
								score=record.score,
							)
						)
					} else {
						document_context_list.append(
							DocumentContext(
								content=segment.get_sign_content(),
								score=record.score,
							)
						)
					}
				}
				if tool.return_resource{
					for record := range  records{
						segment = record.segment
						dataset = db.session.query(Dataset).filter_by(id=segment.DatasetID).first()
						document = (
							db.session.query(DatasetDocument)  // type: ignore
							.where(
								DatasetDocument.ID == segment.document_id,
								DatasetDocument.enabled == True,
								DatasetDocument.archived == False,
							)
							.first()
						)
						if dataset and document{
							source = RetrievalSourceMetadata(
								dataset_id=dataset.ID,
								dataset_name=dataset.name,
								document_id=document.ID,  // type: ignore
								document_name=document.name,  // type: ignore
								data_source_type=document.data_source_type,  // type: ignore
								segment_id=segment.ID,
								retriever_from=tool.retriever_from,
								score=record.score or 0.0,
								doc_metadata=document.doc_metadata,  // type: ignore
							)

							if tool.retriever_from == "dev"{
								source.hit_count = segment.hit_count
								source.word_count = segment.word_count
								source.segment_position = segment.position
								source.index_node_hash = segment.index_node_hash
							}
							if segment.answer{
								source.content = f"question:{segment.content} \nanswer:{segment.answer}"
							} else {
								source.content = segment.content
							}
							retrieval_resource_list.append(source)
						}
					}
				}
			}
		}
		if tool.return_resource and retrieval_resource_list{
			retrieval_resource_list = sorted(
				retrieval_resource_list,
				key=lambda x: x.score or 0.0,
				reverse=True,
			)
			for position, item := range  enumerate(retrieval_resource_list, start=1){  // type: ignore
				item.position = position  // type: ignore
			}
			for hit_callback := range  tool.HitCallbacks{
				hit_callback.return_retriever_resource_info(retrieval_resource_list)
			}
		}
		if document_context_list{
			document_context_list = sorted(document_context_list, key=lambda x: x.score or 0.0, reverse=True)
			return str("\n".join([document_context.content for document_context := range  document_context_list]))
		}
		return ""
	}
}
