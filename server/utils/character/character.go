package character

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Utf8ToGbk(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return ""
	}
	return string(d)
}

func NlpTrimStr(src string) string {
	src = strings.ToLower(src)

	for _, v := range []string{"\r", "\r\n", "\n"} {
		src = strings.ReplaceAll(src, v, "")
	}
	for _, v := range []string{" ", "	", "。", "?", "？", "，", "~", "、", ".", "：", "_", "-", ",", "，", "【", "】"} {
		for strings.HasPrefix(src, v) {
			src = strings.TrimPrefix(src, v)
		}
		for strings.HasSuffix(src, v) {
			src = strings.TrimSuffix(src, v)
		}
	}

	return src
}
