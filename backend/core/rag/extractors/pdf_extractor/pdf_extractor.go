package pdfextractor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/odysseythink/mlog"
	ragentities "mlib.com/gofy/server/entities/rag"
)

type PdfExtractor struct {
	_file_path string
}

func New(file_path string) *PdfExtractor {
	return &PdfExtractor{_file_path: file_path}
}

// Extract extracts text from a PDF file.
// Uses pdftotext command-line tool if available, otherwise falls back to basic text extraction.
func (extractor *PdfExtractor) Extract() []*ragentities.Document {
	text, err := extractor.extractWithPdftotext()
	if err != nil {
		mlog.Errorf("PDF extraction with pdftotext failed for %s: %v", extractor._file_path, err)
		// Try fallback: read raw bytes and do basic text extraction
		text, err = extractor.extractBasic()
		if err != nil {
			mlog.Errorf("PDF basic extraction also failed for %s: %v", extractor._file_path, err)
			return nil
		}
	}

	if strings.TrimSpace(text) == "" {
		return nil
	}

	return []*ragentities.Document{
		{
			PageContent: text,
			Metadata: map[string]any{
				"source": extractor._file_path,
				"type":   "pdf",
			},
			Provider: "gofy",
		},
	}
}

// extractWithPdftotext uses the pdftotext command-line tool (from poppler-utils).
func (extractor *PdfExtractor) extractWithPdftotext() (string, error) {
	// Check if pdftotext is available
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return "", fmt.Errorf("pdftotext not found: %w", err)
	}

	cmd := exec.Command("pdftotext", "-layout", "-enc", "UTF-8", extractor._file_path, "-")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("pdftotext failed: %w", err)
	}

	return string(output), nil
}

// extractBasic does basic text extraction by reading readable ASCII/UTF-8 sequences from the PDF binary.
func (extractor *PdfExtractor) extractBasic() (string, error) {
	data, err := os.ReadFile(extractor._file_path)
	if err != nil {
		return "", fmt.Errorf("read file failed: %w", err)
	}

	// Very basic: extract printable text sequences from the PDF binary
	var builder strings.Builder
	var current strings.Builder
	inText := false

	for i := 0; i < len(data); i++ {
		b := data[i]
		if b >= 32 && b < 127 || b == '\n' || b == '\r' || b == '\t' {
			current.WriteByte(b)
			if current.Len() > 4 {
				inText = true
			}
		} else {
			if inText && current.Len() > 10 {
				builder.WriteString(current.String())
				builder.WriteByte('\n')
			}
			current.Reset()
			inText = false
		}
	}

	result := builder.String()
	if len(result) < 10 {
		return "", fmt.Errorf("no readable text found in PDF")
	}
	return result, nil
}
