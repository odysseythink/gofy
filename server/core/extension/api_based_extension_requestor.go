package extension

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"mlib.com/confy"
	"mlib.com/gofy/server/core/exceptions"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

var (
	timeout = [2]int{5, 60}
)

type APIBasedExtensionRequestor struct {
	api_endpoint string
	api_key      string
	httpClient   *http.Client
}

func NewAPIBasedExtensionRequestor(apiEndpoint, apiKey string) *APIBasedExtensionRequestor {
	return &APIBasedExtensionRequestor{
		api_endpoint: apiEndpoint,
		api_key:      apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // 对应 Python (5,60) 的总超时
		},
	}
}

// Request 等价于 Python 的 request 方法
func (r *APIBasedExtensionRequestor) Request(
	point models.APIBasedExtensionPointType,
	params map[string]any,
) map[string]any {
	body := map[string]any{
		"point":  point,
		"params": params,
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, r.api_endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		mlog.Errorf("create request: %v", err)
		panic(exceptions.NewValueError(fmt.Sprintf("create request: %v", err)))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.api_key)

	// 代理支持（对应 Python 的 SSRF_PROXY_*）
	if httpProxy := confy.GetWithDefault[string]("ssrf.proxy.http_url", ""); httpProxy != "" {
		proxyURL, _ := url.Parse(httpProxy)
		r.httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	if httpsProxy := confy.GetWithDefault[string]("ssrf.proxy.https_url", ""); httpsProxy != "" {
		proxyURL, _ := url.Parse(httpsProxy)
		r.httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		// 简单区分超时/连接错误
		if e, ok := err.(interface{ Timeout() bool }); ok && e.Timeout() {
			mlog.Errorf("request timeout")
			panic(exceptions.NewValueError("request timeout"))
		}
		mlog.Errorf("request connection error: %w", err)
		panic(exceptions.NewValueError(fmt.Sprintf("request connection error: %w", err)))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := make([]byte, 100)
		n, _ := resp.Body.Read(buf)
		mlog.Errorf("request error, status_code: %d, content: %s", resp.StatusCode, string(buf[:n]))
		panic(exceptions.NewValueError(fmt.Sprintf("request error, status_code: %d, content: %s", resp.StatusCode, string(buf[:n]))))
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		mlog.Errorf("decode response failed: %w", err)
		panic(exceptions.NewValueError(fmt.Sprintf("decode response failed: %w", err)))
	}
	return result
}
