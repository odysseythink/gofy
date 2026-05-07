package htmlextractor

import (
	"os"

	"github.com/anaskhan96/soup"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	"github.com/odysseythink/mlog"
)

type HtmlExtractor struct {
	_file_path string
}

func New(
	file_path string,
) *HtmlExtractor {
	return &HtmlExtractor{
		_file_path: file_path,
	}
}

func (extractor *HtmlExtractor) Extract() []*ragentities.Document {
	bindata, err := os.ReadFile(extractor._file_path)
	if err != nil {
		mlog.Errorf("read file=%s failed:%v", extractor._file_path, err)
		panic(err)
	}

	doc := soup.HTMLParse(string(bindata))
	return []*ragentities.Document{
		&ragentities.Document{
			PageContent: doc.FullText(),
			Provider:    "gofy",
		},
	}
}
