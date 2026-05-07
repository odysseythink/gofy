package firecrawl

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"maps"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	"github.com/odysseythink/gofy/backend/utils"
	"github.com/odysseythink/mlog"
	"golang.org/x/exp/errors/fmt"
)

type FirecrawlApp struct {
	api_key  string
	base_url string
}

func New(
	api_key string, base_url string,
) *FirecrawlApp {
	if base_url == "" {
		base_url = "https://api.firecrawl.dev"
	}
	if api_key == "" && base_url == "https://api.firecrawl.dev" {
		panic(exceptions.NewValueError("No API key provided"))
	}
	return &FirecrawlApp{
		api_key:  api_key,
		base_url: base_url,
	}
}
func (app *FirecrawlApp) _prepare_headers() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + app.api_key,
	}
}
func (app *FirecrawlApp) _post_request(url string, data map[string]any, headers map[string]string, retries int, backoff_factor float64) *http.Response {
	if retries <= 0 {
		retries = 3
	}
	if backoff_factor <= 0 {
		backoff_factor = 0.5
	}

	for attempt := range retries {
		mlog.Debugf("------attempt %d times", attempt)
		// 将请求体转换为 JSON 格式
		jsonPayload, err := json.Marshal(data)
		if err != nil {
			mlog.Errorf("marshaling JSON=%#v failed:%v", data, err)
			panic(exceptions.NewValueError("marshaling JSON failed"))
		}

		// 创建 POST 请求
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			mlog.Errorf("creating http request failed:%v", err)
			panic(exceptions.NewValueError("creating http request failed"))
		}

		// 设置请求头
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		// 发送请求
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			mlog.Warningf("sending request failed:%v", err)
			panic(exceptions.NewValueError("sending request failed"))
		}
		if res.StatusCode == 502 {
			tmp := 2 << attempt
			ftmp := backoff_factor * float64(tmp)
			time.Sleep(time.Duration(ftmp) * time.Second)
		} else {
			return res
		}
	}
	panic(exceptions.NewValueError("http request timeout"))
}
func (app *FirecrawlApp) _format_crawl_status_response(
	status string, crawl_status_response map[string]any, url_data_list []map[string]any,
) map[string]any {
	return map[string]any{
		"status":  status,
		"total":   crawl_status_response["total"],
		"current": crawl_status_response["completed"],
		"data":    url_data_list,
	}
}
func (app *FirecrawlApp) _extract_common_fields(item map[string]any) map[string]any {
	metadata := map[string]any{}
	if _, ok := item["metadata"]; ok {
		if _, ok := item["metadata"].(map[string]any); ok {
			metadata = item["metadata"].(map[string]any)
		}
	}
	return map[string]any{
		"title":       metadata["title"],
		"description": metadata["description"],
		"source_url":  metadata["sourceURL"],
		"markdown":    item["markdown"],
	}
}

func (app *FirecrawlApp) _get_request(url string, headers map[string]string, retries int, backoff_factor float64) *http.Response {
	if retries <= 0 {
		retries = 3
	}
	if backoff_factor <= 0 {
		backoff_factor = 0.5
	}

	for attempt := range retries {
		mlog.Debugf("------attempt %d times", attempt)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			mlog.Errorf("creating http request failed:%v", err)
			panic(exceptions.NewValueError("creating http request failed"))
		}

		// 设置请求头
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		// 发送请求
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			mlog.Warningf("sending request failed:%v", err)
			panic(exceptions.NewValueError("sending request failed"))
		}
		if res.StatusCode == 502 {
			tmp := 2 << attempt
			ftmp := backoff_factor * float64(tmp)
			time.Sleep(time.Duration(ftmp) * time.Second)
		} else {
			return res
		}
	}
	panic(exceptions.NewValueError("http request timeout"))
}
func (app *FirecrawlApp) _handle_error(response *http.Response, action string) {
	bodydata, err := io.ReadAll(response.Body)
	if err != nil {
		mlog.Errorf("read response body failed:%v", err)
		panic(err)
	}
	bodydict := map[string]any{}
	err = json.Unmarshal(bodydata, &bodydict)
	if err != nil {
		mlog.Errorf("json unmarshal response body=%s to dict failed:%v", string(bodydata), err)
		panic(err)
	}
	errormsg := "Unknown error occurred"
	if _, ok := bodydict["error"]; ok {
		if _, ok := bodydict["error"].(string); ok {
			errormsg = bodydict["error"].(string)
		}
	}
	panic(fmt.Errorf("Failed to %s. Status code: %d. Error: %s", action, response.StatusCode, errormsg))
}

func (app *FirecrawlApp) ScrapeURL(url string, params map[string]any) map[string]any {
	// Documentation: https://docs.firecrawl.dev/api-reference/endpoint/scrape
	headers := app._prepare_headers()
	json_data := map[string]any{
		"url":             url,
		"formats":         []string{"markdown"},
		"onlyMainContent": true,
		"timeout":         30000,
	}
	maps.Copy(json_data, params)
	response := app._post_request(fmt.Sprintf("%s/v1/scrape", app.base_url), json_data, headers, 0, 0.0)
	defer response.Body.Close()
	if response.StatusCode == 200 {
		bodydata, err := io.ReadAll(response.Body)
		if err != nil {
			mlog.Errorf("read response body failed:%v", err)
			panic(err)
		}
		bodydict := map[string]any{}
		err = json.Unmarshal(bodydata, &bodydict)
		if err != nil {
			mlog.Errorf("json unmarshal response body=%s to dict failed:%v", string(bodydata), err)
			panic(err)
		}
		data := map[string]any{}
		if _, ok := bodydict["data"]; ok {
			if _, ok := bodydict["data"].(map[string]any); ok {
				data = bodydict["data"].(map[string]any)
			}
		}

		return app._extract_common_fields(data)
	} else if slices.Contains([]int{402, 409, 500, 429, 408}, response.StatusCode) {
		app._handle_error(response, "scrape URL")
		return nil // Avoid additional exception after handling error
	} else {
		panic(fmt.Errorf("Failed to scrape URL. Status code: %d", response.StatusCode))
	}
}

func (app *FirecrawlApp) CrawlURL(url string, params map[string]any) string {
	// Documentation: https://docs.firecrawl.dev/api-reference/endpoint/crawl-post
	headers := app._prepare_headers()
	json_data := map[string]any{"url": url}
	maps.Copy(json_data, params)
	response := app._post_request(fmt.Sprintf("%s/v1/crawl", app.base_url), json_data, headers, 0, 0.)
	defer response.Body.Close()
	if response.StatusCode == 200 {
		// There's also another two fields in the response: "success" (bool) and "url" (str)
		bodydata, err := io.ReadAll(response.Body)
		if err != nil {
			mlog.Errorf("read response body failed:%v", err)
			panic(err)
		}
		bodydict := map[string]any{}
		err = json.Unmarshal(bodydata, &bodydict)
		if err != nil {
			mlog.Errorf("json unmarshal response body=%s to dict failed:%v", string(bodydata), err)
			panic(err)
		}
		job_id := ""
		if _, ok := bodydict["id"]; ok {
			if _, ok := bodydict["id"].(string); ok {
				job_id = bodydict["id"].(string)
			}
		}
		return job_id
	} else {
		app._handle_error(response, "start crawl job")
		// FIXME: unreachable code for mypy
		return "" // unreachable
	}
}

func (app *FirecrawlApp) CheckCrawlStatus(job_id string) map[string]any {
	headers := app._prepare_headers()
	response := app._get_request(fmt.Sprintf("%s/v1/crawl/{job_id}", app.base_url), headers, 0, 0.)
	defer response.Body.Close()
	if response.StatusCode == 200 {
		bodydata, err := io.ReadAll(response.Body)
		if err != nil {
			mlog.Errorf("read response body failed:%v", err)
			panic(err)
		}
		crawl_status_response := map[string]any{}
		err = json.Unmarshal(bodydata, &crawl_status_response)
		if err != nil {
			mlog.Errorf("json unmarshal response body=%s to dict failed:%v", string(bodydata), err)
			panic(err)
		}
		status := ""
		if _, ok := crawl_status_response["status"]; ok {
			if _, ok := crawl_status_response["status"].(string); ok {
				status = crawl_status_response["status"].(string)
			}
		}
		if status == "completed" {
			total := 0
			if _, ok := crawl_status_response["total"]; ok {
				if _, ok := crawl_status_response["total"].(int); ok {
					total = crawl_status_response["total"].(int)
				}
			}
			if total == 0 {
				panic(errors.New("failed to check crawl status. Error: No page found"))
			}
			data := []any{}
			if _, ok := crawl_status_response["data"]; ok {
				if _, ok := crawl_status_response["data"].([]any); ok {
					data = crawl_status_response["data"].([]any)
				}
			}

			url_data_list := []map[string]any{}
			for _, item := range data {
				if _, ok := item.(map[string]any); ok {
					if _, ok := item.(map[string]any)["metadata"]; ok {
						if _, ok := item.(map[string]any)["markdown"]; ok {
							url_data := app._extract_common_fields(item.(map[string]any))
							url_data_list = append(url_data_list, url_data)
						}
					}
				}
			}
			if len(url_data_list) > 0 {
				file_key := "website_files/" + job_id + ".txt"
				// try{
				if utils.FileExist(file_key) {
					os.Remove(file_key)
				}
				bindata, _ := json.Marshal(url_data_list)
				os.WriteFile(file_key, bindata, 0644)
			}
			return app._format_crawl_status_response("completed", crawl_status_response, url_data_list)
		} else {
			return app._format_crawl_status_response(status, crawl_status_response, nil)
		}
	} else {
		app._handle_error(response, "check crawl status")
		// FIXME: unreachable code for mypy
		return nil // unreachable
	}
}
