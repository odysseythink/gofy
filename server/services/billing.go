package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"

	"mlib.com/gofy/server/core/exceptions"
)

type BillingService struct {
	secret_key string
	base_url   string
}

func (s *BillingService) _get_base_url() string {
	if s.base_url == "" {
		s.base_url = os.Getenv("BILLING_API_URL")
	}
	return s.base_url
}
func (s *BillingService) _get_secret_key() string {
	if s.secret_key == "" {
		s.secret_key = os.Getenv("BILLING_API_SECRET_KEY")
	}
	return s.secret_key
}

func (s *BillingService) GetInfo(tenant_id string) map[string]any {

	return map[string]any{}
}

func (s *BillingService) _send_request(method string, endpoint string, jsonData any, params url.Values) map[string]any {
	if !slices.Contains([]string{"GET", "POST", "DELETE"}, method) {
		panic(exceptions.NewValueError(fmt.Sprintf("unsupported method=%s", method)))
	}
	// 创建请求 URL
	u, err := url.Parse(s._get_base_url() + endpoint)
	if err != nil {
		panic(exceptions.NewValueError(fmt.Sprintf("failed to parse URL: %w", err)))
	}

	// 添加查询参数
	if params != nil {
		u.RawQuery = params.Encode()
	}

	// 设置请求头
	headers := http.Header{
		"Content-Type":           {"application/json"},
		"Billing-Api-Secret-Key": {s._get_secret_key()},
	}

	// 创建请求体
	var body io.Reader
	if jsonData != nil {
		jsonDataBytes, err := json.Marshal(jsonData)
		if err != nil {
			panic(exceptions.NewValueError(fmt.Sprintf("failed to marshal JSON: %w", err)))
		}
		body = bytes.NewBuffer(jsonDataBytes)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		panic(exceptions.NewValueError(fmt.Sprintf("failed to create request: %w", err)))
	}
	req.Header = headers

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(exceptions.NewValueError(fmt.Sprintf("failed to send request: %w", err)))
	}
	defer resp.Body.Close()

	// 检查响应状态码（仅针对 GET 请求）
	if method == "GET" && resp.StatusCode != http.StatusOK {
		panic(exceptions.NewValueError("unable to retrieve billing information. Please try again later or contact support"))
	}

	// 解析响应内容
	var responseData map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		panic(exceptions.NewValueError(fmt.Sprintf("failed to decode response: %w", err)))
	}

	return responseData
}

func (s *BillingService) IsEmailInFreeze(email string) (status bool) {
	params := map[string][]string{"email": []string{email}}
	return func() (status bool) {
		defer func() {
			if r := recover(); r != nil {
				status = false
			}
		}()
		response := s._send_request("GET", "/account/in-freeze", nil, params)
		status = false
		if _, ok := response["data"]; ok {
			if _, ok := response["data"].(bool); ok {
				status = response["data"].(bool)
			}
		}
		return
	}()
}
