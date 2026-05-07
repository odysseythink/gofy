package excelextractor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	ragentities "github.com/odysseythink/gofy/backend/entities/rag"
	"github.com/odysseythink/mlog"
	"github.com/xuri/excelize/v2"
)

type ExcelExtractor struct {
	_file_path           string
	_encoding            string
	_autodetect_encoding bool
}

func New(
	file_path string,
	encoding string,
	autodetect_encoding bool,
) *ExcelExtractor {
	return &ExcelExtractor{
		_file_path:           file_path,
		_encoding:            encoding,
		_autodetect_encoding: autodetect_encoding,
	}
}
func (extractor *ExcelExtractor) _get_cell_col_name_by_idx(idx int) string {
	if idx < 0 {
		mlog.Errorf("invalid col idx=%d", idx)
		panic(exceptions.NewValueError("invalid col idx"))
	}
	tmp := idx % 26
	tmp1 := idx / 26
	if tmp1 > 25 {
		mlog.Errorf("col idx=%d is beyond max range", idx)
		panic(exceptions.NewValueError("col idx is beyond max range"))
	}
	if tmp1 == 0 {
		return string('A' + tmp)
	} else {
		return string('A'+tmp1) + string('A'+tmp)
	}
}
func (extractor *ExcelExtractor) Extract() []*ragentities.Document {

	documents := []*ragentities.Document{}
	f, err := excelize.OpenFile(extractor._file_path)
	if err != nil {
		mlog.Errorf("open excel file=%s failed:%v", extractor._file_path, err)
		panic(err)
	}
	defer f.Close()
	for _, sheet := range f.GetSheetList() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			mlog.Errorf("get excel file=%s sheet=%s rows failed:%v", extractor._file_path, sheet, err)
			panic(err)
		}
		if len(rows) < 1 {
			continue
		}
		columns := map[int]string{}
		for idx, v := range rows[0] {
			columns[idx] = v
		}
		if len(rows) > 1 {
			for index, row := range rows[1:] {
				page_content := []string{}

				for col_index, val := range row {
					cell_name := extractor._get_cell_col_name_by_idx(col_index)
					cell_name = cell_name + strconv.Itoa(index+2)
					ok, link_target, err := f.GetCellHyperLink(sheet, cell_name)
					if err != nil {
						mlog.Errorf("get excel file=%s sheet=%s cell=%s hyper link failed:%v", extractor._file_path, sheet, cell_name, err)
						panic(err)
					}

					if ok {
						value := fmt.Sprintf("[%s](%s)", val, link_target)
						page_content = append(page_content, fmt.Sprintf(`"%s":"%s"`, columns[col_index], value))
					} else {
						page_content = append(page_content, fmt.Sprintf(`"%s":"%s"`, columns[col_index], val))
					}
				}
				documents = append(documents, &ragentities.Document{
					PageContent: strings.Join(page_content, ";"),
					Metadata:    map[string]any{"source": extractor._file_path},
					Provider:    "gofy",
				})
			}
		}
	}
	return documents
}
