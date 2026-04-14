// Package httpvec provides a small HTTP client helper shared by REST-based
// vector store adapters (qdrant, milvus, weaviate, elasticsearch, opensearch,
// chroma, upstash, vikingdb, baidu, huawei, …).
package httpvec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Endpoint string
	Headers  map[string]string
	HTTP     *http.Client
}

func New(endpoint string, timeout time.Duration, headers map[string]string) *Client {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		Endpoint: strings.TrimRight(endpoint, "/"),
		Headers:  headers,
		HTTP:     &http.Client{Timeout: timeout},
	}
}

// Do executes method+path with optional JSON body, decodes response JSON into
// out (if non-nil). 4xx/5xx produces an error containing status and body.
func (c *Client) Do(ctx context.Context, method, path string, body, out any) error {
	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(buf)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.Endpoint+path, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("http %s %s -> %d: %s", method, path, resp.StatusCode, string(buf))
	}
	if out != nil && len(buf) > 0 {
		return json.Unmarshal(buf, out)
	}
	return nil
}

// Is404 reports whether err looks like a 404 response from Do.
func Is404(err error) bool {
	return err != nil && strings.Contains(err.Error(), "-> 404")
}
