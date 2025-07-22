package baseextractor

import ragentities "mlib.com/gofy/server/entities/rag"

type Extractor interface {
	Extract() []*ragentities.Document
}
