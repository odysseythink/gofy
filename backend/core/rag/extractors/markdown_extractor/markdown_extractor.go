package markdownextractor

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	ragentities "mlib.com/gofy/server/entities/rag"
)

type MarkdownExtractor struct {
	_file_path           string
	_remove_hyperlinks   bool
	_remove_images       bool
	_encoding            string
	_autodetect_encoding bool
}

func New(
	file_path string,
	remove_hyperlinks bool,
	remove_images bool,
	encoding string,
	autodetect_encoding bool, /* = True*/
) *MarkdownExtractor {
	return &MarkdownExtractor{
		_file_path:           file_path,
		_remove_hyperlinks:   remove_hyperlinks,
		_remove_images:       remove_images,
		_encoding:            encoding,
		_autodetect_encoding: autodetect_encoding,
	}
}

func (extractor *MarkdownExtractor) Extract() []*ragentities.Document {
	tups := extractor.parse_tups(extractor._file_path)
	documents := []*ragentities.Document{}
	for _, tup := range tups {
		header := tup[0]
		value := tup[1]
		value = strings.TrimSpace(value)
		if header == "" {
			documents = append(documents, &ragentities.Document{PageContent: value, Provider: "gofy"})
		} else {
			documents = append(documents, &ragentities.Document{PageContent: fmt.Sprintf("\n\n%s\n%s", header, value), Provider: "gofy"})
		}
	}
	return documents
	// return []*ragentities.Document{
	// 	{
	// 		PageContent: doc.FullText(),
	// 		Provider:    "gofy",
	// 		Metadata:    map[string]any{"source": extractor._file_path},
	// 	},
	// }
}

func (extractor *MarkdownExtractor) markdown_to_tups(markdown_text string) [][2]string {
	/*Convert a markdown file to a dictionary.

	  The keys are the headers and the values are the text under each header.

	*/
	markdown_tups := [][2]string{}
	lines := strings.Split(markdown_text, "\n")

	current_header := ""
	current_text := ""
	code_block_flag := false
	header_regx, err := regexp.Compile(`^#+\s`)
	if err != nil {
		mlog.Error("Error with regex:", err)
		panic(exceptions.NewValueError("Error with regex MatchString:" + err.Error()))
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			code_block_flag = !code_block_flag
			current_text += line + "\n"
			continue
		}
		if code_block_flag {
			current_text += line + "\n"
			continue
		}

		header_match := header_regx.Match([]byte(line))
		if header_match {
			if current_header != "" {
				markdown_tups = append(markdown_tups, [2]string{current_header, current_text})
			}
			current_header = line
			current_text = ""
		} else {
			current_text += line + "\n"
		}
	}
	markdown_tups = append(markdown_tups, [2]string{current_header, current_text})

	if current_header != "" {
		// pass linting, assert keys are defined
		for idx, tuple := range markdown_tups {
			key := tuple[0]
			value := tuple[1]

			// 替换 '#' 并去首尾空白
			key = strings.Trim(regexp.MustCompile(`#`).ReplaceAllString(key, ""), " ")

			// 替换 HTML 标签
			value = regexp.MustCompile(`<.*?>`).ReplaceAllString(value, "")
			markdown_tups[idx] = [2]string{key, value}
		}
	} else {
		for idx, tuple := range markdown_tups {
			key := tuple[0]
			value := tuple[1]
			value = regexp.MustCompile(`\n`).ReplaceAllString(value, "")
			markdown_tups[idx] = [2]string{key, value}
		}
	}
	return markdown_tups
}
func (extractor *MarkdownExtractor) remove_images(content string) string {
	/*Get a dictionary of a markdown file from its path.*/
	pattern := `!{1}\[\[(.*)\]\]`
	re := regexp.MustCompile(pattern)
	content = re.ReplaceAllString(content, "")
	return content
}
func (extractor *MarkdownExtractor) remove_hyperlinks(content string) string {
	/*Get a dictionary of a markdown file from its path.*/
	pattern := `\[(.*?)\]\((.*?)\)`
	re := regexp.MustCompile(pattern)
	content = re.ReplaceAllString(content, `\1`)
	return content
}
func (extractor *MarkdownExtractor) parse_tups(file_path string) [][2]string {
	/*Parse file into tuples.*/
	// try{
	bin_content, err := os.ReadFile(file_path)
	if err != nil {
		mlog.Errorf("loading %s failed:%v", file_path, err)
		panic(exceptions.NewRuntimeError(fmt.Sprintf("loading %s failed:%v", file_path, err)))
	}
	content := string(bin_content)
	if extractor._remove_hyperlinks {
		content = extractor.remove_hyperlinks(content)
	}
	if extractor._remove_images {
		content = extractor.remove_images(content)
	}
	return extractor.markdown_to_tups(content)
}
