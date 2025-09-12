// from __future__ import annotations

// from abc import ABC, abstractmethod
// from typing import Any

// from core.rag.models.document import Document
// from models.dataset import Dataset
package keyword

import (
	ragentities "mlib.com/gofy/server/entities/rag"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils/mapstruct"
)

type Keywordor interface {
	Create(texts []*ragentities.Document, args ...any) Keywordor
	AddTexts(texts []*ragentities.Document, args ...any)
	TextExists(id string) bool
	DeleteByIDs(ids []string)
	Delete()
	Search(query string, args ...any) []*ragentities.Document
}

type BaseKeyword struct {
	Dataset *models.Dataset
}

func NewBaseKeyword(dataset *models.Dataset) *BaseKeyword {
	return &BaseKeyword{
		Dataset: dataset,
	}
}

func (key *BaseKeyword) FilterDuplicateTexts(intf Keywordor, texts []*ragentities.Document) []*ragentities.Document {
	new_texts := []*ragentities.Document{}
	for _, text := range texts {
		if len(text.Metadata) == 0 {
			continue
		}
		doc_id := mapstruct.Get(text.Metadata, "doc_id", "")
		exists_duplicate_node := intf.TextExists(doc_id)
		if !exists_duplicate_node {
			new_texts = append(new_texts, text)
		}
	}
	return new_texts
}
func (key *BaseKeyword) GetUUIDs(texts []*ragentities.Document) []string {
	ids := []string{}
	for _, text := range texts {
		if len(text.Metadata) == 0 {
			continue
		}
		ids = append(ids, mapstruct.Get(text.Metadata, "doc_id", ""))
	}
	return ids
}
