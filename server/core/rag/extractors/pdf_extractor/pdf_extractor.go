package pdfextractor

import (
	ragentities "mlib.com/gofy/server/entities/rag"
)

type PdfExtractor struct {
}

func New(file_path string) *PdfExtractor {
	return &PdfExtractor{}
}

func (extractor *PdfExtractor) Extract() []*ragentities.Document {
	return nil
}
