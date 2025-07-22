package character

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GbkToUtf8(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return ""
	}
	return string(d)
}

func Gb2312ToUtf8(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.HZGB2312.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return ""
	}
	return string(d)
}

func IoReaderGb2312ToUtf8(r io.Reader) string {
	reader := transform.NewReader(r, simplifiedchinese.HZGB2312.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return ""
	}
	return string(d)
}

func BytesGb2312ToUtf8(r []byte) string {
	reader := transform.NewReader(bytes.NewReader(r), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return ""
	}
	return string(d)
}

// func Utf8ToGbk(s string) string {
// 	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewEncoder())
// 	d, e := io.ReadAll(reader)
// 	if e != nil {
// 		return ""
// 	}
// 	return string(d)
// }
