package textextractor

import (
	"os"

	"github.com/anaskhan96/soup"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	"github.com/odysseythink/mlog"
)

type TextExtractor struct {
	_file_path           string
	_encoding            string
	_autodetect_encoding bool
}

func New(
	file_path string,
	encoding string,
	autodetect_encoding bool,
) *TextExtractor {
	return &TextExtractor{
		_file_path:           file_path,
		_encoding:            encoding,
		_autodetect_encoding: autodetect_encoding,
	}
}

func (extractor *TextExtractor) Extract() []*ragentities.Document {
	bindata, err := os.ReadFile(extractor._file_path)
	if err != nil {
		mlog.Errorf("read file=%s failed:%v", extractor._file_path, err)
		panic(err)
	}

	doc := soup.HTMLParse(string(bindata))
	return []*ragentities.Document{
		{
			PageContent: doc.FullText(),
			Provider:    "gofy",
			Metadata:    map[string]any{"source": extractor._file_path},
		},
	}
}
