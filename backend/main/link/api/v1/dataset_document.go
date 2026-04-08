package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mlib.com/gofy/server/services"
)

type DatasetDocumentApi struct{}

// GetDocuments returns paginated documents in a dataset.
func (a *DatasetDocumentApi) GetDocuments(c *gin.Context) {
	datasetID := c.Param("dataset_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("keyword")
	status := c.Query("indexing_status")

	docs, total := services.ServiceGroupApp.Document.GetPaginatedDocuments(datasetID, page, limit, search, status)
	c.JSON(http.StatusOK, gin.H{"data": docs, "total": total, "page": page, "limit": limit})
}

// DeleteDocument deletes a document.
func (a *DatasetDocumentApi) DeleteDocument(c *gin.Context) {
	documentID := c.Param("document_id")
	doc := services.ServiceGroupApp.Document.GetDocumentByID(documentID)
	if doc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	if err := services.ServiceGroupApp.Document.DeleteDocument(doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// RenameDocument renames a document.
func (a *DatasetDocumentApi) RenameDocument(c *gin.Context) {
	datasetID := c.Param("dataset_id")
	documentID := c.Param("document_id")
	var body struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	doc, err := services.ServiceGroupApp.Document.RenameDocument(datasetID, documentID, body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doc)
}

// PauseDocument pauses document indexing.
func (a *DatasetDocumentApi) PauseDocument(c *gin.Context) {
	documentID := c.Param("document_id")
	doc := services.ServiceGroupApp.Document.GetDocumentByID(documentID)
	if doc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	if err := services.ServiceGroupApp.Document.PauseDocument(doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// RecoverDocument resumes paused document indexing.
func (a *DatasetDocumentApi) RecoverDocument(c *gin.Context) {
	documentID := c.Param("document_id")
	doc := services.ServiceGroupApp.Document.GetDocumentByID(documentID)
	if doc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	if err := services.ServiceGroupApp.Document.RecoverDocument(doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// BatchUpdateStatus enables/disables/archives documents in batch.
func (a *DatasetDocumentApi) BatchUpdateStatus(c *gin.Context) {
	datasetID := c.Param("dataset_id")
	var body struct {
		DocumentIDs []string `json:"document_ids"`
		Action      string   `json:"action"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.ServiceGroupApp.Document.BatchUpdateDocumentStatus(datasetID, body.DocumentIDs, body.Action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// GetSegments returns paginated segments for a document.
func (a *DatasetDocumentApi) GetSegments(c *gin.Context) {
	documentID := c.Param("document_id")
	tenantID := c.GetString("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	var statusList []string
	if status != "" {
		statusList = []string{status}
	}

	segments, total := services.ServiceGroupApp.Segment.GetSegments(documentID, tenantID, statusList, keyword, page, limit)
	c.JSON(http.StatusOK, gin.H{"data": segments, "total": total})
}

// CreateSegment creates a new segment in a document.
func (a *DatasetDocumentApi) CreateSegment(c *gin.Context) {
	datasetID := c.Param("dataset_id")
	documentID := c.Param("document_id")

	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc := services.ServiceGroupApp.Document.GetDocumentByID(documentID)
	if doc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}

	dataset, _ := services.ServiceGroupApp.Dataset.GetByIDAndTenantID(datasetID, "")
	if dataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dataset not found"})
		return
	}

	segment, err := services.ServiceGroupApp.Segment.CreateSegment(body, doc, dataset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, segment)
}

// DeleteSegment deletes a segment.
func (a *DatasetDocumentApi) DeleteSegment(c *gin.Context) {
	datasetID := c.Param("dataset_id")
	documentID := c.Param("document_id")
	segmentID := c.Param("segment_id")

	if err := services.ServiceGroupApp.Segment.DeleteSegment(segmentID, documentID, datasetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// HitTesting performs retrieval testing on a dataset.
func (a *DatasetDocumentApi) HitTesting(c *gin.Context) {
	// TODO: Implement hit testing using RAG retrieval
	c.JSON(http.StatusOK, gin.H{"records": []any{}})
}
