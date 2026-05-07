package csvextractor

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	"github.com/odysseythink/mlog"
)

type CSVExtractor struct {
	_file_path           string
	_encoding            string
	_autodetect_encoding bool
	source_column        string
	csv_args             map[string]any
}

func New(
	file_path string,
	encoding string,
	autodetect_encoding bool,
	source_column string,
	csv_args map[string]any,
) *CSVExtractor {
	return &CSVExtractor{
		_file_path:           file_path,
		_encoding:            encoding,
		_autodetect_encoding: autodetect_encoding,
		source_column:        source_column,
		csv_args:             csv_args,
	}
}

func (extractor *CSVExtractor) Extract() []*ragentities.Document {
	f, err := os.Open(extractor._file_path)
	if err != nil {
		mlog.Errorf("open file=%s failed:%v", extractor._file_path, err)
		panic(exceptions.NewRuntimeError("error loading " + extractor._file_path))
	}
	defer f.Close()
	return extractor._read_from_file(f)
}

func (extractor *CSVExtractor) _read_from_file(csvfile *os.File) []*ragentities.Document {
	docs := []*ragentities.Document{}

	csv_reader := csv.NewReader(csvfile)
	data, err := csv_reader.ReadAll()
	if err != nil {
		mlog.Errorf("read csv file=%s failed:%v", extractor._file_path, err)
		panic(err)
	}
	if len(data) < 1 {
		return nil
	}
	columns := map[string]int{}
	for idx, v := range data[0] {
		columns[v] = idx
	}
	// check source column exists
	if extractor.source_column != "" && !slices.Contains(data[0], extractor.source_column) {
		panic(exceptions.NewValueError(fmt.Sprintf("Source column '%s' not found in CSV file.", extractor.source_column)))
	}

	//     // create document objects
	if len(data) > 1 {
		for idx, row := range data[1:] {
			tmplist := []string{}
			for col_idx, col := range data[0] {
				tmplist = append(tmplist, fmt.Sprintf("%s: %s", strings.TrimSpace(col), strings.TrimSpace(row[col_idx])))
			}
			content := strings.Join(tmplist, ";")
			source := ""
			if extractor.source_column != "" {
				source = row[columns[extractor.source_column]]
			}
			metadata := map[string]any{"source": source, "row": idx}
			doc := &ragentities.Document{
				PageContent: content,
				Metadata:    metadata,
				Provider:    "gofy",
			}
			docs = append(docs, doc)
		}
	}

	return docs
}
