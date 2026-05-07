package baseextractor

import ragentities "github.com/odysseythink/gofy/backend/entities/rag"

type Extractor interface {
	Extract() []*ragentities.Document
}
