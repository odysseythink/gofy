package notionextractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
	"gorm.io/datatypes"
	"gorm.io/gen"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	ragentities "mlib.com/gofy/server/entities/rag"
	"mlib.com/gofy/server/models"
)

const (
	BLOCK_CHILD_URL_TMPL = "https://api.notion.com/v1/blocks/{block_id}/children"
	DATABASE_URL_TMPL    = "https://api.notion.com/v1/databases/{database_id}/query"
	SEARCH_URL           = "https://api.notion.com/v1/search"

	RETRIEVE_PAGE_URL_TMPL     = "https://api.notion.com/v1/pages/{page_id}"
	RETRIEVE_DATABASE_URL_TMPL = "https://api.notion.com/v1/databases/{database_id}"
)

var (
	// if user want split by headings, use the corresponding splitter
	HEADING_SPLITTER = map[string]string{
		"heading_1": "# ",
		"heading_2": "## ",
		"heading_3": "### ",
	}
)

type NotionExtractor struct {
	_notion_access_token string
	_notion_obj_id       string
	_notion_page_type    string
	_notion_workspace_id string
	_document_model      *models.Document
}

func New(
	notion_workspace_id string,
	notion_obj_id string,
	notion_page_type string,
	tenant_id string,
	document_model *models.Document,
	notion_access_token string,
) *NotionExtractor {
	ext := &NotionExtractor{
		_notion_access_token: notion_access_token,
		_notion_obj_id:       notion_obj_id,
		_notion_page_type:    notion_page_type,
		_notion_workspace_id: notion_workspace_id,
		_document_model:      document_model,
	}
	if ext._notion_access_token == "" {
		ext._notion_access_token = ext._get_access_token(tenant_id, ext._notion_workspace_id)
		if ext._notion_access_token == "" {
			integration_token := confy.Get[string]("notion.integration_token")
			if integration_token == "" {
				panic(exceptions.NewValueError("Must specify `integration_token` or set environment variable `NOTION_INTEGRATION_TOKEN`."))
			}
			ext._notion_access_token = integration_token
		}
	}

	return ext

}

func (extractor *NotionExtractor) _read_table_rows(block_id string) string {
	/*Read table rows.*/
	if extractor._notion_access_token == "" {
		panic(exceptions.NewValueError("Notion access token is required"))
	}
	done := false
	result_lines_arr := []string{}
	var start_cursor any
	block_url := strings.ReplaceAll(BLOCK_CHILD_URL_TMPL, "{block_id}", block_id)
	for !done {
		params := url.Values{}
		if start_cursor != nil {
			params.Add("start_cursor", fmt.Sprintf("%v", start_cursor))
		}
		full_url := ""
		if len(params) > 0 {
			full_url = fmt.Sprintf("%s?%s", block_url, params.Encode())
		} else {
			full_url = block_url
		}

		req, err := http.NewRequest("GET", full_url, nil)
		if err != nil {
			mlog.Errorf("creating http request failed:%v", err)
			panic(exceptions.NewValueError("creating http request failed"))
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", extractor._notion_access_token))
		req.Header.Set("Notion-Version", confy.GetWithDefault[string]("notion.version", "2022-06-28"))
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
		// get table headers text
		table_header_cell_texts := []string{}
		if _, ok := bodydict["results"].([]any); ok && len(bodydict["results"].([]any)) > 0 {
			if _, ok := bodydict["results"].([]any)[0].(map[string]any); ok {
				if _, ok := bodydict["results"].([]any)[0].(map[string]any)["table_row"]; ok {
					if table_row, ok := bodydict["results"].([]any)[0].(map[string]any)["table_row"].(map[string]any); ok {

						if _, ok := table_row["cells"]; ok {
							if table_header_cells, ok := table_row["cells"].([]any); ok {
								for _, table_header_cell := range table_header_cells {
									if table_header_cell != nil {
										if _, ok := table_header_cell.([]any); ok {
											for _, table_header_cell_text := range table_header_cell.([]any) {
												if _, ok := table_header_cell_text.(map[string]any); ok {
													if _, ok := table_header_cell_text.(map[string]any)["text"]; ok {
														if _, ok := table_header_cell_text.(map[string]any)["text"].(map[string]any); ok {
															if _, ok := table_header_cell_text.(map[string]any)["text"].(map[string]any)["content"]; ok {
																if _, ok := table_header_cell_text.(map[string]any)["text"].(map[string]any)["content"].(string); ok {
																	text := table_header_cell_text.(map[string]any)["text"].(map[string]any)["content"].(string)
																	table_header_cell_texts = append(table_header_cell_texts, text)
																}
															}
														}
													}
												}
											}
										}
									} else {
										table_header_cell_texts = append(table_header_cell_texts, "")
									}
								}
							}
						}

					}
				}
			}
		}

		// Initialize Markdown table with headers
		markdown_table := "| " + strings.Join(table_header_cell_texts, " | ") + " |\n"
		markdown_table += "| " + strings.Join(slices.Repeat([]string{"---"}, len(table_header_cell_texts)), " | ") + " |\n"

		// Process data to format each row in Markdown table format
		if _, ok := bodydict["results"].([]any); ok && len(bodydict["results"].([]any)) > 1 {
			for _, result := range bodydict["results"].([]any)[1:] {
				column_texts := []string{}

				if _, ok := result.(map[string]any); ok {
					if _, ok := result.(map[string]any)["table_row"]; ok {
						if table_row, ok := result.(map[string]any)["table_row"].(map[string]any); ok {
							if _, ok := table_row["cells"]; ok {
								if table_column_cells, ok := table_row["cells"].([]any); ok {
									for _, table_column_cell := range table_column_cells {
										if table_column_cell != nil {
											if _, ok := table_column_cell.([]any); ok {
												for _, table_column_cell_text := range table_column_cell.([]any) {
													if _, ok := table_column_cell_text.(map[string]any); ok {
														if _, ok := table_column_cell_text.(map[string]any)["text"]; ok {
															if _, ok := table_column_cell_text.(map[string]any)["text"].(map[string]any); ok {
																if _, ok := table_column_cell_text.(map[string]any)["text"].(map[string]any)["content"]; ok {
																	if _, ok := table_column_cell_text.(map[string]any)["text"].(map[string]any)["content"].(string); ok {
																		text := table_column_cell_text.(map[string]any)["text"].(map[string]any)["content"].(string)
																		table_header_cell_texts = append(table_header_cell_texts, text)
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}

						}
					}
				}
				markdown_table += "| " + strings.Join(column_texts, " | ") + " |\n"
			}
		}

		result_lines_arr = append(result_lines_arr, markdown_table)
		if _, ok := bodydict["next_cursor"]; !ok || bodydict["next_cursor"] == nil {
			done = true
			break
		} else {
			start_cursor = bodydict["next_cursor"]
		}
	}
	result_lines := strings.Join(result_lines_arr, "\n")
	return result_lines
}

func (extractor *NotionExtractor) _read_block(block_id string, num_tabs int /* = 0*/) string {
	/*Read a block.*/
	if extractor._notion_access_token == "" {
		panic(exceptions.NewValueError("Notion access token is required"))
	}
	result_lines_arr := []string{}
	var start_cursor any
	block_url := strings.ReplaceAll(BLOCK_CHILD_URL_TMPL, "{block_id}", block_id)
	for {
		params := url.Values{}
		if start_cursor != nil {
			params.Add("start_cursor", fmt.Sprintf("%v", start_cursor))
		}
		full_url := ""
		if len(params) > 0 {
			full_url = fmt.Sprintf("%s?%s", block_url, params.Encode())
		} else {
			full_url = block_url
		}

		req, err := http.NewRequest("GET", full_url, nil)
		if err != nil {
			mlog.Errorf("creating http request failed:%v", err)
			panic(exceptions.NewValueError("creating http request failed"))
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", extractor._notion_access_token))
		req.Header.Set("Notion-Version", confy.GetWithDefault[string]("notion.version", "2022-06-28"))
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
		if _, ok := bodydict["results"]; !ok || bodydict["results"] == nil {
			break
		}
		if _, ok := bodydict["results"].([]any); !ok {
			mlog.Warningf("results=%#v must be type of []any", bodydict["results"])
			break
		}

		for _, result := range bodydict["results"].([]any) {
			if _, ok := result.(map[string]any); !ok {
				mlog.Warningf("results=%#v must be type of []map[string]any", bodydict["results"])
				break
			}
			result_type := ""
			if _, ok := result.(map[string]any)["type"]; ok {
				if _, ok := result.(map[string]any)["type"].(string); ok {
					result_type = result.(map[string]any)["type"].(string)
				}
			}
			result_block_id := ""
			if _, ok := result.(map[string]any)["id"]; ok {
				if _, ok := result.(map[string]any)["id"].(string); ok {
					result_block_id = result.(map[string]any)["id"].(string)
				}
			}
			cur_result_text_arr := []string{}
			if result_type == "table" {
				text := extractor._read_table_rows(result_block_id)
				result_lines_arr = append(result_lines_arr, text)
			} else {
				result_obj := map[string]any{}
				if _, ok := result.(map[string]any)[result_type]; ok {
					if _, ok := result.(map[string]any)[result_type].(map[string]any); ok {
						result_obj = result.(map[string]any)[result_type].(map[string]any)
					}
				}

				if _, ok := result_obj["rich_text"]; ok {
					if _, ok := result_obj["rich_text"].([]any); ok {
						for _, rich_text := range result_obj["rich_text"].([]any) {
							if _, ok := rich_text.(map[string]any); ok {
								// skip if doesn't have text object
								if _, ok := rich_text.(map[string]any)["text"]; ok {
									if _, ok := rich_text.(map[string]any)["text"].(map[string]any); ok {
										if _, ok := rich_text.(map[string]any)["text"].(map[string]any)["content"]; ok {
											if _, ok := rich_text.(map[string]any)["text"].(map[string]any)["content"].(string); ok {
												text := rich_text.(map[string]any)["text"].(map[string]any)["content"].(string)
												cur_result_text_arr = append(cur_result_text_arr, strings.Repeat("\t", num_tabs)+text)
											}
										}
									}
								}
							}
						}
					}
				}
				has_children := false
				if _, ok := result.(map[string]any)["has_children"]; ok {
					if _, ok := result.(map[string]any)["has_children"].(bool); ok {
						has_children = result.(map[string]any)["has_children"].(bool)
					}
				}

				block_type := result_type
				if has_children && block_type != "child_page" {
					children_text := extractor._read_block(result_block_id, num_tabs+1)
					cur_result_text_arr = append(cur_result_text_arr, children_text)
				}

				cur_result_text := strings.Join(cur_result_text_arr, "\n")
				if _, ok := HEADING_SPLITTER[result_type]; ok {
					result_lines_arr = append(result_lines_arr, HEADING_SPLITTER[result_type]+cur_result_text)
				} else {
					result_lines_arr = append(result_lines_arr, cur_result_text+"\n\n")
				}
			}
		}
		if _, ok := bodydict["next_cursor"]; !ok || bodydict["next_cursor"] == nil {
			break
		} else {
			start_cursor = bodydict["next_cursor"]
		}
	}

	return strings.Join(result_lines_arr, "\n")
}
func (extractor *NotionExtractor) get_notion_last_edited_time() string {
	if extractor._notion_access_token == "" {
		panic(exceptions.NewValueError("Notion access token is required"))
	}
	obj_id := extractor._notion_obj_id
	page_type := extractor._notion_page_type
	retrieve_page_url := ""
	if page_type == "database" {
		retrieve_page_url = strings.ReplaceAll(RETRIEVE_DATABASE_URL_TMPL, "{database_id}", obj_id)
	} else {
		retrieve_page_url = strings.ReplaceAll(RETRIEVE_PAGE_URL_TMPL, "{page_id}", obj_id)
	}
	req, err := http.NewRequest("GET", retrieve_page_url, nil)
	if err != nil {
		mlog.Errorf("creating http request failed:%v", err)
		panic(exceptions.NewValueError("creating http request failed"))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", extractor._notion_access_token))
	req.Header.Set("Notion-Version", confy.GetWithDefault[string]("notion.version", "2022-06-28"))
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
	last_edited_time := ""
	if _, ok := bodydict["last_edited_time"]; ok {
		if _, ok := bodydict["last_edited_time"].(string); ok {
			last_edited_time = bodydict["last_edited_time"].(string)
		}
	}

	return last_edited_time

}
func (extractor *NotionExtractor) _get_access_token(tenant_id string, notion_workspace_id string) string {
	data_source_binding := new(models.DataSourceOauthBinding)
	err := dbengine.Instance().DB.Model(&models.DataSourceOauthBinding{}).Where("tenant_id =? and provider = ? and disabled = ?", tenant_id, "notion", false).Where(gen.Cond(datatypes.JSONQuery("source_info").Equals(notion_workspace_id, "workspace_id"))).First(data_source_binding).Error
	if err != nil {
		mlog.Errorf("get DataSourceOauthBinding failed:%v", err)
		data_source_binding = nil
	}

	if data_source_binding == nil {
		panic(fmt.Errorf("No notion data source binding found for tenant %s and notion workspace %s", tenant_id, notion_workspace_id))
	}
	return data_source_binding.AccessToken
}

func (extractor *NotionExtractor) Extract() []*ragentities.Document {
	extractor.update_last_edited_time(extractor._document_model)

	text_docs := extractor._load_data_as_documents(extractor._notion_obj_id, extractor._notion_page_type)

	return text_docs
}
func (extractor *NotionExtractor) _get_notion_database_data(database_id string, query_dict map[string]any) []*ragentities.Document {
	/*Get all the pages from a Notion database.*/
	if extractor._notion_access_token == "" {
		panic(exceptions.NewValueError("Notion access token is required"))
	}
	json_payload, _ := json.Marshal(query_dict)
	req, err := http.NewRequest("POST", strings.ReplaceAll(DATABASE_URL_TMPL, "{database_id}", database_id), bytes.NewBuffer(json_payload))
	if err != nil {
		mlog.Errorf("creating http request failed:%v", err)
		panic(exceptions.NewValueError("creating http request failed"))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", extractor._notion_access_token))
	req.Header.Set("Notion-Version", confy.GetWithDefault[string]("notion.version", "2022-06-28"))
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
	if _, ok := bodydict["results"]; !ok || bodydict["results"] == nil {
		return nil
	}
	if _, ok := bodydict["results"].([]any); !ok {
		mlog.Warningf("results=%#v must be type of []any", bodydict["results"])
		return nil
	}
	database_content := []string{}
	for _, result := range bodydict["results"].([]any) {
		if _, ok := result.([]map[string]any); ok || result.([]map[string]any) == nil {
			mlog.Warningf("results=%#v must be type of []map[string]any", bodydict["results"])
			continue
		}
		result_dict := result.(map[string]any)
		if _, ok := result_dict["properties"]; !ok || result_dict["properties"] == nil {
			mlog.Warningf("result=%#v don't contain properties field", result_dict)
			continue
		}
		if _, ok := result_dict["properties"].(map[string]any); !ok {
			mlog.Warningf("result=%#v  properties field must be type of map[string]any", result_dict)
			continue
		}
		properties := result_dict["properties"].(map[string]any)
		data := map[string]any{}
		for property_name, property_value := range properties {
			if _, ok := property_value.(map[string]any); !ok {
				mlog.Warningf("result=%#v  properties value must be type of map[string]any", result_dict)
				continue
			}
			property_value_dict := property_value.(map[string]any)
			typ := ""
			if _, ok := property_value_dict["type"]; ok {
				if _, ok := property_value_dict["type"].(string); ok {
					typ = property_value_dict["type"].(string)
				}
			}

			if typ == "multi_select" {
				value := []any{}
				multi_select_list := []map[string]any{}
				if _, ok := property_value_dict[typ]; ok {
					if _, ok := property_value_dict[typ].([]any); ok {
						for _, v := range property_value_dict[typ].([]any) {
							if _, ok := v.(map[string]any); !ok {
								mlog.Warningf("property_value_dict[%s]=%#v  must be type of []map[string]any", typ, property_value_dict[typ])
								break
							} else {
								multi_select_list = append(multi_select_list, v.(map[string]any))
							}
						}
					} else if _, ok := property_value_dict[typ].([]map[string]any); ok {
						multi_select_list = property_value_dict[typ].([]map[string]any)
					}
				}
				for _, multi_select := range multi_select_list {
					value = append(value, multi_select["name"])
				}
				data[property_name] = value
			} else if slices.Contains([]string{"rich_text", "title"}, typ) {
				value := ""
				if _, ok := property_value_dict[typ]; ok {
					if _, ok := property_value_dict[typ].([]any); ok && len(property_value_dict[typ].([]any)) > 0 {
						if _, ok := property_value_dict[typ].([]any)[0].(map[string]any); !ok {
							mlog.Warningf("property_value_dict[%s]=%#v  must be type of []map[string]any", typ, property_value_dict[typ])
						} else {
							if _, ok := property_value_dict[typ].([]any)[0].(map[string]any)["plain_text"]; ok {
								if _, ok := property_value_dict[typ].([]any)[0].(map[string]any)["plain_text"].(string); ok {
									value = property_value_dict[typ].([]any)[0].(map[string]any)["plain_text"].(string)
								}
							}
						}
					} else if _, ok := property_value_dict[typ].([]map[string]any); ok && len(property_value_dict[typ].([]map[string]any)) > 0 {
						if _, ok := property_value_dict[typ].([]map[string]any)[0]["plain_text"]; ok {
							if _, ok := property_value_dict[typ].([]map[string]any)[0]["plain_text"].(string); ok {
								value = property_value_dict[typ].([]map[string]any)[0]["plain_text"].(string)
							}
						}
					}
				}
				data[property_name] = value
			} else if slices.Contains([]string{"select", "status"}, typ) {
				value := ""
				if _, ok := property_value_dict[typ]; ok {
					if _, ok := property_value_dict[typ].(map[string]any); ok {
						if _, ok := property_value_dict[typ].(map[string]any)["name"]; ok {
							if _, ok := property_value_dict[typ].(map[string]any)["plain_text"].(string); ok {
								value = property_value_dict[typ].(map[string]any)["plain_text"].(string)
							}
						}
					}
				}
				data[property_name] = value
			} else {
				data[property_name] = property_value_dict[typ]
			}
		}

		row_dict := map[string]any{}
		for k, v := range data {
			if v != nil {
				row_dict[k] = v
			}
		}
		row_content := ""
		for key, value := range row_dict {
			if _, ok := value.(map[string]any); ok {
				value_dict := map[string]any{}
				for k, v := range value.(map[string]any) {
					if v != nil {
						value_dict[k] = v
					}
				}
				value_list := []string{}
				for k, v := range value_dict {
					value_list = append(value_list, fmt.Sprintf("%s:%s ", k, v))
				}
				value_content := strings.Join(value_list, "")
				row_content = row_content + fmt.Sprintf("%s:%s\n", key, value_content)
			} else {
				row_content = row_content + fmt.Sprintf("%s:%v\n", key, value)
			}
		}
		database_content = append(database_content, row_content)
	}
	return []*ragentities.Document{&ragentities.Document{PageContent: strings.Join(database_content, "\n"), Provider: "gofy"}}
}

func (extractor *NotionExtractor) _get_notion_block_data(page_id string) []string {
	if extractor._notion_access_token == "" {
		panic(exceptions.NewValueError("Notion access token is required"))
	}
	result_lines_arr := []string{}
	var start_cursor any
	block_url := strings.ReplaceAll(BLOCK_CHILD_URL_TMPL, "block_id", page_id)
	for {
		params := url.Values{}
		if start_cursor != nil {
			params.Add("start_cursor", fmt.Sprintf("%v", start_cursor))
		}
		full_url := ""
		if len(params) > 0 {
			full_url = fmt.Sprintf("%s?%s", block_url, params.Encode())
		} else {
			full_url = block_url
		}

		req, err := http.NewRequest("GET", full_url, nil)
		if err != nil {
			mlog.Errorf("creating http request failed:%v", err)
			panic(exceptions.NewValueError("creating http request failed"))
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", extractor._notion_access_token))
		req.Header.Set("Notion-Version", confy.GetWithDefault[string]("notion.version", "2022-06-28"))
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
		if resp.StatusCode != 200 {
			mlog.Errorf("Error fetching Notion block data: %s", string(bodydata))
			panic(exceptions.NewValueError(fmt.Sprintf("Error fetching Notion block data: %s", string(bodydata))))
		}

		bodydict := map[string]any{}
		err = json.Unmarshal(bodydata, &bodydict)
		if err != nil {
			mlog.Errorf("json unmarshal response body=%s to dict failed:%v", string(bodydata), err)
			panic(err)
		}
		if _, ok := bodydict["results"]; !ok {
			panic(exceptions.NewValueError("Error fetching Notion block data"))
		} else {
			if _, ok := bodydict["results"].([]any); !ok {
				panic(exceptions.NewValueError("Error fetching Notion block data"))
			}
		}

		for _, result := range bodydict["results"].([]any) {
			if _, ok := result.(map[string]any); !ok {
				mlog.Errorf("results=%#v must be type of []map[string]any", bodydict["results"])
				panic(exceptions.NewValueError("results must be type of []map[string]any"))
			}
			result_type := ""
			if _, ok := result.(map[string]any)["type"]; ok {
				if _, ok := result.(map[string]any)["type"].(string); ok {
					result_type = result.(map[string]any)["type"].(string)
				}
			}

			cur_result_text_arr := []string{}
			if result_type == "table" {
				result_block_id := ""
				if _, ok := result.(map[string]any)["id"]; ok {
					if _, ok := result.(map[string]any)["id"].(string); ok {
						result_block_id = result.(map[string]any)["id"].(string)
					}
				}
				text := extractor._read_table_rows(result_block_id)
				text += "\n\n"
				result_lines_arr = append(result_lines_arr, text)
			} else {
				result_obj := map[string]any{}
				if _, ok := result.(map[string]any)[result_type]; ok {
					if _, ok := result.(map[string]any)[result_type].(map[string]any); ok {
						result_obj = result.(map[string]any)[result_type].(map[string]any)
					}
				}

				if _, ok := result_obj["rich_text"]; ok {
					if _, ok := result_obj["rich_text"].([]any); ok {
						for _, rich_text := range result_obj["rich_text"].([]any) {
							if _, ok := rich_text.(map[string]any); ok {
								// skip if doesn't have text object
								if _, ok := rich_text.(map[string]any)["text"]; ok {
									if _, ok := rich_text.(map[string]any)["text"].(map[string]any); ok {
										if _, ok := rich_text.(map[string]any)["text"].(map[string]any)["content"]; ok {
											if _, ok := rich_text.(map[string]any)["text"].(map[string]any)["content"].(string); ok {
												text := rich_text.(map[string]any)["text"].(map[string]any)["content"].(string)
												cur_result_text_arr = append(cur_result_text_arr, text)
											}
										}
									}
								}
							}
						}
					}
				}
				result_block_id := ""
				if _, ok := result.(map[string]any)["id"]; ok {
					if _, ok := result.(map[string]any)["id"].(string); ok {
						result_block_id = result.(map[string]any)["id"].(string)
					}
				}
				has_children := false
				if _, ok := result.(map[string]any)["has_children"]; ok {
					if _, ok := result.(map[string]any)["has_children"].(bool); ok {
						has_children = result.(map[string]any)["has_children"].(bool)
					}
				}

				block_type := result_type
				if has_children && block_type != "child_page" {
					children_text := extractor._read_block(result_block_id, 1)
					cur_result_text_arr = append(cur_result_text_arr, children_text)
				}
				cur_result_text := strings.Join(cur_result_text_arr, "\n")
				if _, ok := HEADING_SPLITTER[result_type]; ok {
					result_lines_arr = append(result_lines_arr, HEADING_SPLITTER[result_type]+cur_result_text)
				} else {
					result_lines_arr = append(result_lines_arr, cur_result_text+"\n\n")
				}
			}
		}
		if _, ok := bodydict["next_cursor"]; !ok || bodydict["next_cursor"] == nil {
			break
		} else {
			start_cursor = bodydict["next_cursor"]
		}
	}
	return result_lines_arr
}
func (extractor *NotionExtractor) _load_data_as_documents(notion_obj_id string, notion_page_type string) []*ragentities.Document {
	docs := []*ragentities.Document{}
	if notion_page_type == "database" {
		// get all the pages in the database
		page_text_documents := extractor._get_notion_database_data(extractor._notion_obj_id, nil)
		docs = append(docs, page_text_documents...)
	} else if notion_page_type == "page" {
		page_text_list := extractor._get_notion_block_data(notion_obj_id)
		docs = append(docs, &ragentities.Document{PageContent: strings.Join(page_text_list, "\n"), Provider: "gofy"})
	} else {
		panic(exceptions.NewValueError("notion page type not supported"))
	}
	return docs
}

func (extractor *NotionExtractor) update_last_edited_time(document_model *models.Document) {
	if document_model == nil {
		return
	}
	last_edited_time := extractor.get_notion_last_edited_time()
	data_source_info := document_model.DataSourceInfoDict()
	data_source_info["last_edited_time"] = last_edited_time
	bindata, _ := json.Marshal(data_source_info)
	dbengine.Instance().DB.Updates(&models.Document{ID: document_model.ID, DataSourceInfo: string(bindata)})
}
