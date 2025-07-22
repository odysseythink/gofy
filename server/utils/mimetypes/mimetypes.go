package mimetypes

import (
	"mime"
	"net/url"
	"path"
	"strings"
)

func GuessType(url_str string, strict bool) string {
	parsed_url, err := url.Parse(url_str)
	if err != nil {
		return ""
	}

	var scheme string
	var filePath string

	if parsed_url.Scheme != "" && len(parsed_url.Scheme) > 1 {
		scheme = parsed_url.Scheme
		filePath = parsed_url.Path
	} else {
		scheme = ""
		filePath = path.Clean(url_str)
	}

	if scheme == "data" {
		commaIndex := strings.Index(filePath, ",")
		if commaIndex < 0 {
			return ""
		}
		semiIndex := strings.Index(filePath[:commaIndex], ";")
		if semiIndex >= 0 {
			mimeType := filePath[:semiIndex]
			if strings.Contains(mimeType, "=") || !strings.Contains(mimeType, "/") {
				return "text/plain"
			}
			return mimeType
		} else {
			mimeType := filePath[:commaIndex]
			if strings.Contains(mimeType, "=") || !strings.Contains(mimeType, "/") {
				return "text/plain"
			}
			return mimeType
		}
	}

	_, ext := path.Split(filePath)
	extLower := strings.ToLower(ext)
	mime_type := mime.TypeByExtension(extLower)

	return mime_type
}
