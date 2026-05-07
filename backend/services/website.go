package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	datamanageres "github.com/odysseythink/gofy/backend/data_manageres"
	"github.com/odysseythink/gofy/backend/utils"
	firecrawlutils "github.com/odysseythink/gofy/backend/utils/firecrawl"
	"github.com/odysseythink/mlog"
)

type WebsiteService struct {
}

func (service *WebsiteService) GetScrapeURLData(provider string, url string, tenant_id string, only_main_content bool) map[string]any {
	credentials := datamanageres.ManagerGroupApp.ApiKeyAuth.GetAuthCredentials(tenant_id, "website", provider)
	// decrypt api_key
	api_key := ""
	if _, ok := credentials["config"]; ok {
		if _, ok := credentials["config"].(map[string]any); ok {
			if _, ok := credentials["config"].(map[string]any)["api_key"]; ok {
				if _, ok := credentials["config"].(map[string]any)["api_key"].(string); ok {
					api_key = credentials["config"].(map[string]any)["api_key"].(string)
				}
			}
		}
	}
	if provider == "firecrawl" {
		base_url := ""
		if _, ok := credentials["config"]; ok {
			if _, ok := credentials["config"].(map[string]any); ok {
				if _, ok := credentials["config"].(map[string]any)["base_url"]; ok {
					if _, ok := credentials["config"].(map[string]any)["base_url"].(string); ok {
						base_url = credentials["config"].(map[string]any)["base_url"].(string)
					}
				}
			}
		}
		firecrawl_app := firecrawlutils.New(api_key, base_url)
		params := map[string]any{"onlyMainContent": only_main_content}
		return firecrawl_app.ScrapeURL(url, params)
	} else {
		mlog.Errorf("invalid provider=%s", provider)
		panic(exceptions.NewValueError("Invalid provider"))
	}
}

func (service *WebsiteService) GetCrawlURLData(job_id string, provider string, url string, tenant_id string) map[string]any {
	credentials := datamanageres.ManagerGroupApp.ApiKeyAuth.GetAuthCredentials(tenant_id, "website", provider)
	// decrypt api_key
	api_key := ""
	if _, ok := credentials["config"]; ok {
		if _, ok := credentials["config"].(map[string]any); ok {
			if _, ok := credentials["config"].(map[string]any)["api_key"]; ok {
				if _, ok := credentials["config"].(map[string]any)["api_key"].(string); ok {
					api_key = credentials["config"].(map[string]any)["api_key"].(string)
				}
			}
		}
	}

	// FIXME data is redefine too many times here, use Any to ease the type checking, fix it later
	// data: Any
	if provider == "firecrawl" {
		data_list := []map[string]any{}
		file_key := "website_files/" + job_id + ".txt"
		if utils.FileExist(file_key) {
			data, err := os.ReadFile(file_key)
			if err != nil {
				mlog.Errorf("read file=%s failed:%v", file_key, err)
				panic(exceptions.NewValueError("read Crawl job file failed"))
			}
			err = json.Unmarshal(data, &data_list)
			if err != nil {
				mlog.Errorf("json Unmarshal data=%s to []any failed:%v", string(data), err)
				panic(exceptions.NewValueError("json Unmarshal data to []any failed"))
			}
		} else {
			base_url := ""
			if _, ok := credentials["config"]; ok {
				if _, ok := credentials["config"].(map[string]any); ok {
					if _, ok := credentials["config"].(map[string]any)["base_url"]; ok {
						if _, ok := credentials["config"].(map[string]any)["base_url"].(string); ok {
							base_url = credentials["config"].(map[string]any)["base_url"].(string)
						}
					}
				}
			}
			firecrawl_app := firecrawlutils.New(api_key, base_url)
			result := firecrawl_app.CheckCrawlStatus(job_id)
			status := ""
			if _, ok := result["status"]; ok {
				if _, ok := result["status"].(string); ok {
					status = result["status"].(string)
				}
			}
			if status != "completed" {
				panic(exceptions.NewValueError("Crawl job is not completed"))
			}
			if _, ok := result["data"]; ok {
				if _, ok := result["data"].([]any); ok {
					for _, v := range result["data"].([]any) {
						if _, ok := v.(map[string]any); ok {
							data_list = append(data_list, v.(map[string]any))
						} else {
							mlog.Errorf("result[\"data\"]=%#v must be []map[string]any", result["data"])
							panic(exceptions.NewValueError("result[\"data\"] must be []map[string]any"))
						}
					}
				} else if _, ok := result["data"].([]map[string]any); ok {
					data_list = result["data"].([]map[string]any)
				}
			}
		}
		if len(data_list) > 0 {
			for _, item := range data_list {
				if _, ok := item["source_url"]; ok {
					if _, ok := item["source_url"].(string); ok {
						if item["source_url"].(string) == url {
							return item
						}
					}
				}
			}
		}
		return nil
	} else if provider == "jinareader" {
		if job_id == "" {
			req, err := http.NewRequest("GET", fmt.Sprintf("https://r.jina.ai/%s", url), nil)
			if err != nil {
				mlog.Errorf("creating http request failed:%v", err)
				panic(exceptions.NewValueError("creating http request failed"))
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key))
			// 发送请求
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				mlog.Errorf("sending request failed:%v", err)
				panic(exceptions.NewValueError("sending request failed"))
			}
			defer resp.Body.Close()
			bodydata, err := io.ReadAll(resp.Body)
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
			code := 0
			if _, ok := bodydict["code"]; ok {
				if _, ok := bodydict["code"].(int); ok {
					code = bodydict["code"].(int)
				}
			}
			if code != 200 {
				mlog.Errorf("Failed to crawl")
				panic(exceptions.NewValueError("Failed to crawl"))
			}
			data := map[string]any{}
			if _, ok := bodydict["data"]; ok {
				if _, ok := bodydict["data"].(map[string]any); ok {
					data = bodydict["data"].(map[string]any)
				}
			}
			return data
		} else {
			jsonPayload, _ := json.Marshal(map[string]any{"taskId": job_id})

			// 创建 POST 请求
			req, err := http.NewRequest("POST", "https://adaptivecrawlstatus-kir3wx7b3a-uc.a.run.app", bytes.NewBuffer(jsonPayload))
			if err != nil {
				mlog.Errorf("creating http request failed:%v", err)
				panic(exceptions.NewValueError("creating http request failed"))
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key))

			// 发送请求
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				mlog.Warningf("sending request failed:%v", err)
				panic(exceptions.NewValueError("sending request failed"))
			}
			defer resp.Body.Close()
			bodydata, err := io.ReadAll(resp.Body)
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
			status := ""
			if _, ok := data["status"]; ok {
				if _, ok := data["status"].(string); ok {
					status = data["status"].(string)
				}
			}
			if status != "completed" {
				mlog.Errorf("Crawl job is not completed")
				panic(exceptions.NewValueError("Crawl job is not completed"))
			}
			processed := map[string]any{}
			if _, ok := data["processed"]; ok {
				if _, ok := data["processed"].(map[string]any); ok {
					processed = data["processed"].(map[string]any)
				}
			}
			jsonPayload, _ = json.Marshal(map[string]any{"taskId": job_id, "urls": maps.Keys(processed)})
			req, _ = http.NewRequest("POST", "https://adaptivecrawlstatus-kir3wx7b3a-uc.a.run.app", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key))
			resp1, err := http.DefaultClient.Do(req)
			if err != nil {
				mlog.Warningf("sending request failed:%v", err)
				panic(exceptions.NewValueError("sending request failed"))
			}
			defer resp1.Body.Close()
			bodydata1, err := io.ReadAll(resp1.Body)
			if err != nil {
				mlog.Errorf("read response body failed:%v", err)
				panic(err)
			}
			bodydict1 := map[string]any{}
			err = json.Unmarshal(bodydata1, &bodydict1)
			if err != nil {
				mlog.Errorf("json unmarshal response body=%s to dict failed:%v", string(bodydata), err)
				panic(err)
			}
			data = map[string]any{}
			if _, ok := bodydict1["data"]; ok {
				if _, ok := bodydict1["data"].(map[string]any); ok {
					data = bodydict1["data"].(map[string]any)
				}
			}
			processed = map[string]any{}
			if _, ok := data["processed"]; ok {
				if _, ok := data["processed"].(map[string]any); ok {
					processed = data["processed"].(map[string]any)
				}
			}
			for _, item := range processed {
				if _, ok := item.(map[string]any); ok {
					if item_data, ok := item.(map[string]any)["data"]; ok {
						if item_data_dict, ok := item_data.(map[string]any); ok {
							if _, ok := item_data_dict["url"]; ok {
								if _, ok := item_data_dict["url"].(string); ok && item_data_dict["url"].(string) == url {
									return item_data_dict
								}
							}
						}
					}
				}
			}
		}
		return nil
	} else {
		mlog.Errorf("invalid provider=%s", provider)
		panic(exceptions.NewValueError("Invalid provider"))
	}
}
