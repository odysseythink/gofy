package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odysseythink/gofy/backend/cluster"
	servicesentities "github.com/odysseythink/gofy/backend/entities/services"
	modelruntimeenumtypes "github.com/odysseythink/gofy/backend/enum_types/model_runtime"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/models/response"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/mlog"
)

var (
	currentProviders = `[       {
            "provider": "openai",
            "label": {
                "zh_Hans": "OpenAI",
                "en_US": "OpenAI"
            },
            "description": {
                "zh_Hans": "OpenAI \u63d0\u4f9b\u7684\u6a21\u578b\uff0c\u4f8b\u5982 GPT-3.5-Turbo \u548c GPT-4\u3002",
                "en_US": "Models provided by OpenAI, such as GPT-3.5-Turbo and GPT-4."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openai/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openai/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openai/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openai/icon_large/en_US"
            },
            "background": "#E5E7EB",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece OpenAI \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from OpenAI"
                },
                "url": {
                    "zh_Hans": "https://platform.openai.com/account/api-keys",
                    "en_US": "https://platform.openai.com/account/api-keys"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding",
                "speech2text",
                "moderation",
                "tts"
            ],
            "configurate_methods": [
                "predefined-model",
                "customizable-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "openai_api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "openai_organization",
                        "label": {
                            "zh_Hans": "\u7ec4\u7ec7 ID",
                            "en_US": "Organization"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u7ec4\u7ec7 ID",
                            "en_US": "Enter your Organization ID"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "openai_api_base",
                        "label": {
                            "zh_Hans": "API Base",
                            "en_US": "API Base"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Base, \u5982\uff1ahttps://api.openai.com",
                            "en_US": "Enter your API Base, e.g. https://api.openai.com"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "openai_api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "openai_organization",
                        "label": {
                            "zh_Hans": "\u7ec4\u7ec7 ID",
                            "en_US": "Organization"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u7ec4\u7ec7 ID",
                            "en_US": "Enter your Organization ID"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "openai_api_base",
                        "label": {
                            "zh_Hans": "API Base",
                            "en_US": "API Base"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Base",
                            "en_US": "Enter your API Base"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "anthropic",
            "label": {
                "zh_Hans": "Anthropic",
                "en_US": "Anthropic"
            },
            "description": {
                "zh_Hans": "Anthropic \u7684\u5f3a\u5927\u6a21\u578b\uff0c\u4f8b\u5982 Claude 3\u3002",
                "en_US": "Anthropic\u2019s powerful models, such as Claude 3."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/anthropic/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/anthropic/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/anthropic/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/anthropic/icon_large/en_US"
            },
            "background": "#F0F0EB",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece Anthropic \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from Anthropic"
                },
                "url": {
                    "zh_Hans": "https://console.anthropic.com/account/keys",
                    "en_US": "https://console.anthropic.com/account/keys"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "anthropic_api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "anthropic_api_url",
                        "label": {
                            "zh_Hans": "API URL",
                            "en_US": "API URL"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API URL",
                            "en_US": "Enter your API URL"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "azure_openai",
            "label": {
                "zh_Hans": "Azure OpenAI Service Model",
                "en_US": "Azure OpenAI Service Model"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/azure_openai/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/azure_openai/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/azure_openai/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/azure_openai/icon_large/en_US"
            },
            "background": "#E3F0FF",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece Azure \u83b7\u53d6 API Key",
                    "en_US": "Get your API key from Azure"
                },
                "url": {
                    "zh_Hans": "https://azure.microsoft.com/en-us/products/ai-services/openai-service",
                    "en_US": "https://azure.microsoft.com/en-us/products/ai-services/openai-service"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding",
                "speech2text",
                "tts"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u90e8\u7f72\u540d\u79f0",
                        "en_US": "Deployment Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u90e8\u7f72\u540d\u79f0\uff0c\u4e0e Azure \u90e8\u7f72\u540d\u79f0\u5339\u914d\u3002",
                        "en_US": "Enter your Deployment Name here, matching the Azure deployment name."
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "openai_api_base",
                        "label": {
                            "zh_Hans": "API \u57df\u540d",
                            "en_US": "API Endpoint URL"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API \u57df\u540d\uff0c\u5982\uff1ahttps://example.com/xxx",
                            "en_US": "Enter your API Endpoint, eg: https://example.com/xxx"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "openai_api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API key here"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "openai_api_version",
                        "label": {
                            "zh_Hans": "API \u7248\u672c",
                            "en_US": "API Version"
                        },
                        "type": "select",
                        "required": true,
                        "default": null,
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "2024-05-01-preview",
                                    "en_US": "2024-05-01-preview"
                                },
                                "value": "2024-05-01-preview",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "2024-04-01-preview",
                                    "en_US": "2024-04-01-preview"
                                },
                                "value": "2024-04-01-preview",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "2024-03-01-preview",
                                    "en_US": "2024-03-01-preview"
                                },
                                "value": "2024-03-01-preview",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "2024-02-15-preview",
                                    "en_US": "2024-02-15-preview"
                                },
                                "value": "2024-02-15-preview",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "2023-12-01-preview",
                                    "en_US": "2023-12-01-preview"
                                },
                                "value": "2023-12-01-preview",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "2024-02-01",
                                    "en_US": "2024-02-01"
                                },
                                "value": "2024-02-01",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "2024-06-01",
                                    "en_US": "2024-06-01"
                                },
                                "value": "2024-06-01",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u9009\u62e9\u60a8\u7684 API \u7248\u672c",
                            "en_US": "Select your API Version here"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "base_model_name",
                        "label": {
                            "zh_Hans": "\u57fa\u7840\u6a21\u578b",
                            "en_US": "Base Model"
                        },
                        "type": "select",
                        "required": true,
                        "default": null,
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "gpt-35-turbo",
                                    "en_US": "gpt-35-turbo"
                                },
                                "value": "gpt-35-turbo",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-35-turbo-0125",
                                    "en_US": "gpt-35-turbo-0125"
                                },
                                "value": "gpt-35-turbo-0125",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-35-turbo-16k",
                                    "en_US": "gpt-35-turbo-16k"
                                },
                                "value": "gpt-35-turbo-16k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4",
                                    "en_US": "gpt-4"
                                },
                                "value": "gpt-4",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4-32k",
                                    "en_US": "gpt-4-32k"
                                },
                                "value": "gpt-4-32k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4o-mini",
                                    "en_US": "gpt-4o-mini"
                                },
                                "value": "gpt-4o-mini",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4o-mini-2024-07-18",
                                    "en_US": "gpt-4o-mini-2024-07-18"
                                },
                                "value": "gpt-4o-mini-2024-07-18",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4o",
                                    "en_US": "gpt-4o"
                                },
                                "value": "gpt-4o",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4o-2024-05-13",
                                    "en_US": "gpt-4o-2024-05-13"
                                },
                                "value": "gpt-4o-2024-05-13",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4-turbo",
                                    "en_US": "gpt-4-turbo"
                                },
                                "value": "gpt-4-turbo",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4-turbo-2024-04-09",
                                    "en_US": "gpt-4-turbo-2024-04-09"
                                },
                                "value": "gpt-4-turbo-2024-04-09",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4-0125-preview",
                                    "en_US": "gpt-4-0125-preview"
                                },
                                "value": "gpt-4-0125-preview",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4-1106-preview",
                                    "en_US": "gpt-4-1106-preview"
                                },
                                "value": "gpt-4-1106-preview",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-4-vision-preview",
                                    "en_US": "gpt-4-vision-preview"
                                },
                                "value": "gpt-4-vision-preview",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "gpt-35-turbo-instruct",
                                    "en_US": "gpt-35-turbo-instruct"
                                },
                                "value": "gpt-35-turbo-instruct",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "text-embedding-ada-002",
                                    "en_US": "text-embedding-ada-002"
                                },
                                "value": "text-embedding-ada-002",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "text-embedding"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "text-embedding-3-small",
                                    "en_US": "text-embedding-3-small"
                                },
                                "value": "text-embedding-3-small",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "text-embedding"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "text-embedding-3-large",
                                    "en_US": "text-embedding-3-large"
                                },
                                "value": "text-embedding-3-large",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "text-embedding"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "whisper-1",
                                    "en_US": "whisper-1"
                                },
                                "value": "whisper-1",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "speech2text"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "tts-1",
                                    "en_US": "tts-1"
                                },
                                "value": "tts-1",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "tts"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "tts-1-hd",
                                    "en_US": "tts-1-hd"
                                },
                                "value": "tts-1-hd",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "tts"
                                    }
                                ]
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u6a21\u578b\u7248\u672c",
                            "en_US": "Enter your model version"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "google",
            "label": {
                "zh_Hans": "Google",
                "en_US": "Google"
            },
            "description": {
                "zh_Hans": "\u8c37\u6b4c\u63d0\u4f9b\u7684 Gemini \u6a21\u578b.",
                "en_US": "Google's Gemini model."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/google/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/google/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/google/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/google/icon_large/en_US"
            },
            "background": "#FCFDFF",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece Google \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from Google"
                },
                "url": {
                    "zh_Hans": "https://ai.google.dev/",
                    "en_US": "https://ai.google.dev/"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "google_api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "vertex_ai",
            "label": {
                "zh_Hans": "Vertex AI | Google Cloud Platform",
                "en_US": "Vertex AI | Google Cloud Platform"
            },
            "description": {
                "zh_Hans": "Vertex AI in Google Cloud Platform.",
                "en_US": "Vertex AI in Google Cloud Platform."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/vertex_ai/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/vertex_ai/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/vertex_ai/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/vertex_ai/icon_large/en_US"
            },
            "background": "#FCFDFF",
            "help": {
                "title": {
                    "zh_Hans": "Get your Access Details from Google",
                    "en_US": "Get your Access Details from Google"
                },
                "url": {
                    "zh_Hans": "https://cloud.google.com/vertex-ai/",
                    "en_US": "https://cloud.google.com/vertex-ai/"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "vertex_project_id",
                        "label": {
                            "zh_Hans": "Project ID",
                            "en_US": "Project ID"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Enter your Google Cloud Project ID",
                            "en_US": "Enter your Google Cloud Project ID"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "vertex_location",
                        "label": {
                            "zh_Hans": "Location",
                            "en_US": "Location"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Enter your Google Cloud Location",
                            "en_US": "Enter your Google Cloud Location"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "vertex_service_account_key",
                        "label": {
                            "zh_Hans": "Service Account Key (Leave blank if you use Application Default Credentials)",
                            "en_US": "Service Account Key (Leave blank if you use Application Default Credentials)"
                        },
                        "type": "secret-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Enter your Google Cloud Service Account Key in base64 format",
                            "en_US": "Enter your Google Cloud Service Account Key in base64 format"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "nvidia",
            "label": {
                "zh_Hans": "API Catalog",
                "en_US": "API Catalog"
            },
            "description": {
                "zh_Hans": "API Catalog",
                "en_US": "API Catalog"
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/nvidia/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/nvidia/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/nvidia/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/nvidia/icon_large/en_US"
            },
            "background": "#FFFFFF",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece NVIDIA \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from NVIDIA"
                },
                "url": {
                    "zh_Hans": "https://build.nvidia.com/explore/discover",
                    "en_US": "https://build.nvidia.com/explore/discover"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding",
                "rerank"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "nvidia_nim",
            "label": {
                "zh_Hans": "NVIDIA NIM",
                "en_US": "NVIDIA NIM"
            },
            "description": {
                "zh_Hans": "NVIDIA NIM\uff0c\u4e00\u7ec4\u6613\u4e8e\u4f7f\u7528\u7684\u6a21\u578b\u63a8\u7406\u5fae\u670d\u52a1\u3002",
                "en_US": "NVIDIA NIM, a set of easy-to-use inference microservices."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/nvidia_nim/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/nvidia_nim/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/nvidia_nim/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/nvidia_nim/icon_large/en_US"
            },
            "background": "#EFFDFD",
            "help": {
                "title": {
                    "zh_Hans": "\u4e86\u89e3 NVIDIA NIM \u66f4\u591a\u4fe1\u606f",
                    "en_US": "Learn more about NVIDIA NIM"
                },
                "url": {
                    "zh_Hans": "https://www.nvidia.com/en-us/ai/",
                    "en_US": "https://www.nvidia.com/en-us/ai/"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u5168\u79f0",
                        "en_US": "Enter full model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "endpoint_url",
                        "label": {
                            "zh_Hans": "API endpoint URL",
                            "en_US": "API endpoint URL"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Base URL, e.g. http://192.168.1.100:8000/v1",
                            "en_US": "Base URL, e.g. http://192.168.1.100:8000/v1"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "mode",
                        "label": {
                            "zh_Hans": "Completion mode",
                            "en_US": "Completion mode"
                        },
                        "type": "select",
                        "required": false,
                        "default": "chat",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u8865\u5168",
                                    "en_US": "Completion"
                                },
                                "value": "completion",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5bf9\u8bdd",
                                    "en_US": "Chat"
                                },
                                "value": "chat",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u9009\u62e9\u5bf9\u8bdd\u7c7b\u578b",
                            "en_US": "Select completion mode"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Model context size"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Enter your Model context size"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "max_tokens_to_sample",
                        "label": {
                            "zh_Hans": "\u6700\u5927 token \u4e0a\u9650",
                            "en_US": "Upper bound for max tokens"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "cohere",
            "label": {
                "zh_Hans": "Cohere",
                "en_US": "Cohere"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/cohere/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/cohere/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/cohere/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/cohere/icon_large/en_US"
            },
            "background": "#ECE9E3",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece cohere \u83b7\u53d6 API Key",
                    "en_US": "Get your API key from cohere"
                },
                "url": {
                    "zh_Hans": "https://dashboard.cohere.com/api-keys",
                    "en_US": "https://dashboard.cohere.com/api-keys"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding",
                "rerank"
            ],
            "configurate_methods": [
                "predefined-model",
                "customizable-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "base_url",
                        "label": {
                            "zh_Hans": "API Base",
                            "en_US": "API Base"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Base\uff0c\u5982 https://api.cohere.ai/v1",
                            "en_US": "Enter your API Base, e.g. https://api.cohere.ai/v1"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "mode",
                        "label": {
                            "zh_Hans": "Completion mode",
                            "en_US": "Completion mode"
                        },
                        "type": "select",
                        "required": false,
                        "default": "chat",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u8865\u5168",
                                    "en_US": "Completion"
                                },
                                "value": "completion",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5bf9\u8bdd",
                                    "en_US": "Chat"
                                },
                                "value": "chat",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u9009\u62e9\u5bf9\u8bdd\u7c7b\u578b",
                            "en_US": "Select completion mode"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "base_url",
                        "label": {
                            "zh_Hans": "API Base",
                            "en_US": "API Base"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Base\uff0c\u5982 https://api.cohere.ai/v1",
                            "en_US": "Enter your API Base, e.g. https://api.cohere.ai/v1"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "upstage",
            "label": {
                "zh_Hans": "Upstage",
                "en_US": "Upstage"
            },
            "description": {
                "zh_Hans": "Upstage \u63d0\u4f9b\u7684\u6a21\u578b\uff0c\u4f8b\u5982 Solar-1-mini-chat.",
                "en_US": "Models provided by Upstage, such as Solar-1-mini-chat."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/upstage/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/upstage/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/upstage/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/upstage/icon_large/en_US"
            },
            "background": "#FFFFF",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece Upstage \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from Upstage"
                },
                "url": {
                    "zh_Hans": "https://console.upstage.ai/api-keys",
                    "en_US": "https://console.upstage.ai/api-keys"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "upstage_api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "upstage_api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "bedrock",
            "label": {
                "zh_Hans": "AWS",
                "en_US": "AWS"
            },
            "description": {
                "zh_Hans": "AWS Bedrock's models.",
                "en_US": "AWS Bedrock's models."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/bedrock/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/bedrock/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/bedrock/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/bedrock/icon_large/en_US"
            },
            "background": "#FCFDFF",
            "help": {
                "title": {
                    "zh_Hans": "Get your Access Key and Secret Access Key from AWS Console",
                    "en_US": "Get your Access Key and Secret Access Key from AWS Console"
                },
                "url": {
                    "zh_Hans": "https://console.aws.amazon.com/",
                    "en_US": "https://console.aws.amazon.com/"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "aws_access_key_id",
                        "label": {
                            "zh_Hans": "Access Key",
                            "en_US": "Access Key (If not provided, credentials are obtained from the running environment.)"
                        },
                        "type": "secret-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Access Key",
                            "en_US": "Enter your Access Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "aws_secret_access_key",
                        "label": {
                            "zh_Hans": "Secret Access Key",
                            "en_US": "Secret Access Key"
                        },
                        "type": "secret-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Secret Access Key",
                            "en_US": "Enter your Secret Access Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "aws_region",
                        "label": {
                            "zh_Hans": "AWS \u5730\u533a",
                            "en_US": "AWS Region"
                        },
                        "type": "select",
                        "required": true,
                        "default": "us-east-1",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u7f8e\u56fd\u4e1c\u90e8 (\u5f17\u5409\u5c3c\u4e9a\u5317\u90e8)",
                                    "en_US": "US East (N. Virginia)"
                                },
                                "value": "us-east-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u7f8e\u56fd\u897f\u90e8 (\u4fc4\u52d2\u5188\u5dde)",
                                    "en_US": "US West (Oregon)"
                                },
                                "value": "us-west-2",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e9a\u592a\u5730\u533a (\u65b0\u52a0\u5761)",
                                    "en_US": "Asia Pacific (Singapore)"
                                },
                                "value": "ap-southeast-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e9a\u592a\u5730\u533a (\u4e1c\u4eac)",
                                    "en_US": "Asia Pacific (Tokyo)"
                                },
                                "value": "ap-northeast-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u6b27\u6d32 (\u6cd5\u5170\u514b\u798f)",
                                    "en_US": "Europe (Frankfurt)"
                                },
                                "value": "eu-central-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u6b27\u6d32\u897f\u90e8 (\u4f26\u6566)",
                                    "en_US": "Eu west London (London)"
                                },
                                "value": "eu-west-2",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "AWS GovCloud (US-West)",
                                    "en_US": "AWS GovCloud (US-West)"
                                },
                                "value": "us-gov-west-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e9a\u592a\u5730\u533a (\u6089\u5c3c)",
                                    "en_US": "Asia Pacific (Sydney)"
                                },
                                "value": "ap-southeast-2",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "model_for_validation",
                        "label": {
                            "zh_Hans": "\u53ef\u7528\u6a21\u578b\u540d\u79f0",
                            "en_US": "Available Model Name"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u4e3a\u4e86\u8fdb\u884c\u9a8c\u8bc1\uff0c\u8bf7\u8f93\u5165\u4e00\u4e2a\u60a8\u53ef\u7528\u7684\u6a21\u578b\u540d\u79f0 (\u4f8b\u5982\uff1aamazon.titan-text-lite-v1)",
                            "en_US": "A model you have access to (e.g. amazon.titan-text-lite-v1) for validation."
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "togetherai",
            "label": {
                "zh_Hans": "together.ai",
                "en_US": "together.ai"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/togetherai/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/togetherai/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/togetherai/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/togetherai/icon_large/en_US"
            },
            "background": "#F1EFED",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece together.ai \u83b7\u53d6 API Key",
                    "en_US": "Get your API key from together.ai"
                },
                "url": {
                    "zh_Hans": "https://api.together.xyz/",
                    "en_US": "https://api.together.xyz/"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u5168\u79f0",
                        "en_US": "Enter full model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "mode",
                        "label": {
                            "zh_Hans": "Completion mode",
                            "en_US": "Completion mode"
                        },
                        "type": "select",
                        "required": false,
                        "default": "chat",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u8865\u5168",
                                    "en_US": "Completion"
                                },
                                "value": "completion",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5bf9\u8bdd",
                                    "en_US": "Chat"
                                },
                                "value": "chat",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u9009\u62e9\u5bf9\u8bdd\u7c7b\u578b",
                            "en_US": "Select completion mode"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Model context size"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Enter your Model context size"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "max_tokens_to_sample",
                        "label": {
                            "zh_Hans": "\u6700\u5927 token \u4e0a\u9650",
                            "en_US": "Upper bound for max tokens"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "openrouter",
            "label": {
                "zh_Hans": "OpenRouter",
                "en_US": "OpenRouter"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openrouter/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openrouter/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openrouter/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openrouter/icon_large/en_US"
            },
            "background": "#F1EFED",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece openrouter.ai \u83b7\u53d6 API Key",
                    "en_US": "Get your API key from openrouter.ai"
                },
                "url": {
                    "zh_Hans": "https://openrouter.ai/keys",
                    "en_US": "https://openrouter.ai/keys"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model",
                "customizable-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u5168\u79f0",
                        "en_US": "Enter full model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "mode",
                        "label": {
                            "zh_Hans": "Completion mode",
                            "en_US": "Completion mode"
                        },
                        "type": "select",
                        "required": false,
                        "default": "chat",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u8865\u5168",
                                    "en_US": "Completion"
                                },
                                "value": "completion",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5bf9\u8bdd",
                                    "en_US": "Chat"
                                },
                                "value": "chat",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u9009\u62e9\u5bf9\u8bdd\u7c7b\u578b",
                            "en_US": "Select completion mode"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Model context size"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Enter your Model context size"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "max_tokens_to_sample",
                        "label": {
                            "zh_Hans": "\u6700\u5927 token \u4e0a\u9650",
                            "en_US": "Upper bound for max tokens"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "vision_support",
                        "label": {
                            "zh_Hans": "\u662f\u5426\u652f\u6301 Vision",
                            "en_US": "Vision Support"
                        },
                        "type": "radio",
                        "required": false,
                        "default": "no_support",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u662f",
                                    "en_US": "Yes"
                                },
                                "value": "support",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5426",
                                    "en_US": "No"
                                },
                                "value": "no_support",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "ollama",
            "label": {
                "zh_Hans": "Ollama",
                "en_US": "Ollama"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/ollama/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/ollama/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/ollama/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/ollama/icon_large/en_US"
            },
            "background": "#F9FAFB",
            "help": {
                "title": {
                    "zh_Hans": "\u5982\u4f55\u96c6\u6210 Ollama",
                    "en_US": "How to integrate with Ollama"
                },
                "url": {
                    "zh_Hans": "https://docs.gofy.ai/tutorials/model-configuration/ollama",
                    "en_US": "https://docs.gofy.ai/tutorials/model-configuration/ollama"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "base_url",
                        "label": {
                            "zh_Hans": "\u57fa\u7840 URL",
                            "en_US": "Base URL"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Ollama server \u7684\u57fa\u7840 URL\uff0c\u4f8b\u5982 http://192.168.1.100:11434",
                            "en_US": "Base url of Ollama server, e.g. http://192.168.1.100:11434"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "mode",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u7c7b\u578b",
                            "en_US": "Completion mode"
                        },
                        "type": "select",
                        "required": true,
                        "default": "chat",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u8865\u5168",
                                    "en_US": "Completion"
                                },
                                "value": "completion",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5bf9\u8bdd",
                                    "en_US": "Chat"
                                },
                                "value": "chat",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u9009\u62e9\u5bf9\u8bdd\u7c7b\u578b",
                            "en_US": "Select completion mode"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Model context size"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Enter your Model context size"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "max_tokens",
                        "label": {
                            "zh_Hans": "\u6700\u5927 token \u4e0a\u9650",
                            "en_US": "Upper bound for max tokens"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "vision_support",
                        "label": {
                            "zh_Hans": "\u662f\u5426\u652f\u6301 Vision",
                            "en_US": "Vision support"
                        },
                        "type": "radio",
                        "required": false,
                        "default": "false",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u662f",
                                    "en_US": "Yes"
                                },
                                "value": "true",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5426",
                                    "en_US": "No"
                                },
                                "value": "false",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "mistralai",
            "label": {
                "zh_Hans": "MistralAI",
                "en_US": "MistralAI"
            },
            "description": {
                "zh_Hans": "MistralAI \u63d0\u4f9b\u7684\u6a21\u578b\uff0c\u4f8b\u5982 open-mistral-7b \u548c mistral-large-latest\u3002",
                "en_US": "Models provided by MistralAI, such as open-mistral-7b and mistral-large-latest."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/mistralai/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/mistralai/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/mistralai/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/mistralai/icon_large/en_US"
            },
            "background": "#FFFFFF",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece MistralAI \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from MistralAI"
                },
                "url": {
                    "zh_Hans": "https://console.mistral.ai/api-keys/",
                    "en_US": "https://console.mistral.ai/api-keys/"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "groq",
            "label": {
                "zh_Hans": "GroqCloud",
                "en_US": "GroqCloud"
            },
            "description": {
                "zh_Hans": "GroqCloud \u63d0\u4f9b\u5bf9 Groq Cloud API \u7684\u8bbf\u95ee\uff0c\u5176\u4e2d\u6258\u7ba1\u4e86 LLama2 \u548c Mixtral \u7b49\u6a21\u578b\u3002",
                "en_US": "GroqCloud provides access to the Groq Cloud API, which hosts models like LLama2 and Mixtral."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/groq/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/groq/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/groq/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/groq/icon_large/en_US"
            },
            "background": "#F5F5F4",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece GroqCloud \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from GroqCloud"
                },
                "url": {
                    "zh_Hans": "https://console.groq.com/",
                    "en_US": "https://console.groq.com/"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "replicate",
            "label": {
                "zh_Hans": "Replicate",
                "en_US": "Replicate"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/replicate/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/replicate/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/replicate/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/replicate/icon_large/en_US"
            },
            "background": "#E5E7EB",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece Replicate \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from Replicate"
                },
                "url": {
                    "zh_Hans": "https://replicate.com/account/api-tokens",
                    "en_US": "https://replicate.com/account/api-tokens"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": null
                },
                "credential_form_schemas": [
                    {
                        "variable": "replicate_api_token",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Replicate API Key",
                            "en_US": "Enter your Replicate API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "model_version",
                        "label": {
                            "zh_Hans": "Model Version",
                            "en_US": "Model Version"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u6a21\u578b\u7248\u672c\uff0c\u9ed8\u8ba4\u4e3a\u6700\u65b0\u7248\u672c",
                            "en_US": "Enter your model version, default to the latest version"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "huggingface_hub",
            "label": {
                "zh_Hans": "Hugging Face Model",
                "en_US": "Hugging Face Model"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/huggingface_hub/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/huggingface_hub/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/huggingface_hub/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/huggingface_hub/icon_large/en_US"
            },
            "background": "#FFF8DC",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece Hugging Face Hub \u83b7\u53d6 API Key",
                    "en_US": "Get your API key from Hugging Face Hub"
                },
                "url": {
                    "zh_Hans": "https://huggingface.co/settings/tokens",
                    "en_US": "https://huggingface.co/settings/tokens"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": null
                },
                "credential_form_schemas": [
                    {
                        "variable": "huggingfacehub_api_type",
                        "label": {
                            "zh_Hans": "\u7aef\u70b9\u7c7b\u578b",
                            "en_US": "Endpoint Type"
                        },
                        "type": "radio",
                        "required": true,
                        "default": "hosted_inference_api",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "Hosted Inference API",
                                    "en_US": "Hosted Inference API"
                                },
                                "value": "hosted_inference_api",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "Inference Endpoints",
                                    "en_US": "Inference Endpoints"
                                },
                                "value": "inference_endpoints",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "huggingfacehub_api_token",
                        "label": {
                            "zh_Hans": "API Token",
                            "en_US": "API Token"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Hugging Face Hub API Token",
                            "en_US": "Enter your Hugging Face Hub API Token here"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "huggingface_namespace",
                        "label": {
                            "zh_Hans": "\u7528\u6237\u540d / \u7ec4\u7ec7\u540d\u79f0",
                            "en_US": "User Name / Organization Name"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u7528\u6237\u540d / \u7ec4\u7ec7\u540d\u79f0",
                            "en_US": "Enter your User Name / Organization Name here"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "text-embedding"
                            },
                            {
                                "variable": "huggingfacehub_api_type",
                                "value": "inference_endpoints"
                            }
                        ]
                    },
                    {
                        "variable": "huggingfacehub_endpoint_url",
                        "label": {
                            "zh_Hans": "\u7aef\u70b9 URL",
                            "en_US": "Endpoint URL"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u7aef\u70b9 URL",
                            "en_US": "Enter your Endpoint URL here"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "huggingfacehub_api_type",
                                "value": "inference_endpoints"
                            }
                        ]
                    },
                    {
                        "variable": "task_type",
                        "label": {
                            "zh_Hans": "Task",
                            "en_US": "Task"
                        },
                        "type": "select",
                        "required": true,
                        "default": null,
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "Text-to-Text Generation",
                                    "en_US": "Text-to-Text Generation"
                                },
                                "value": "text2text-generation",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u6587\u672c\u751f\u6210",
                                    "en_US": "Text Generation"
                                },
                                "value": "text-generation",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Feature Extraction",
                                    "en_US": "Feature Extraction"
                                },
                                "value": "feature-extraction",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "text-embedding"
                                    }
                                ]
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "huggingfacehub_api_type",
                                "value": "inference_endpoints"
                            }
                        ]
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "xinference",
            "label": {
                "zh_Hans": "Xorbits Inference",
                "en_US": "Xorbits Inference"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/xinference/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/xinference/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/xinference/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/xinference/icon_large/en_US"
            },
            "background": "#FAF5FF",
            "help": {
                "title": {
                    "zh_Hans": "\u5982\u4f55\u90e8\u7f72 Xinference",
                    "en_US": "How to deploy Xinference"
                },
                "url": {
                    "zh_Hans": "https://github.com/xorbitsai/inference",
                    "en_US": "https://github.com/xorbitsai/inference"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding",
                "rerank",
                "speech2text",
                "tts"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "server_url",
                        "label": {
                            "zh_Hans": "\u670d\u52a1\u5668URL",
                            "en_US": "Server url"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165Xinference\u7684\u670d\u52a1\u5668\u5730\u5740\uff0c\u5982 http://192.168.1.100:9997",
                            "en_US": "Enter the url of your Xinference, e.g. http://192.168.1.100:9997"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "model_uid",
                        "label": {
                            "zh_Hans": "\u6a21\u578bUID",
                            "en_US": "Model uid"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684Model UID",
                            "en_US": "Enter the model uid"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API\u5bc6\u94a5",
                            "en_US": "API key"
                        },
                        "type": "secret-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684API\u5bc6\u94a5",
                            "en_US": "Enter the api key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "triton_inference_server",
            "label": {
                "zh_Hans": "Triton Inference Server",
                "en_US": "Triton Inference Server"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/triton_inference_server/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/triton_inference_server/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/triton_inference_server/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/triton_inference_server/icon_large/en_US"
            },
            "background": "#EFFDFD",
            "help": {
                "title": {
                    "zh_Hans": "\u5982\u4f55\u90e8\u7f72 Triton Inference Server",
                    "en_US": "How to deploy Triton Inference Server"
                },
                "url": {
                    "zh_Hans": "https://github.com/triton-inference-server/server",
                    "en_US": "https://github.com/triton-inference-server/server"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "server_url",
                        "label": {
                            "zh_Hans": "\u670d\u52a1\u5668URL",
                            "en_US": "Server url"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165 Triton Inference Server \u7684\u670d\u52a1\u5668\u5730\u5740\uff0c\u5982 http://192.168.1.100:8000",
                            "en_US": "Enter the url of your Triton Inference Server, e.g. http://192.168.1.100:8000"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u4e0a\u4e0b\u6587\u5927\u5c0f",
                            "en_US": "Context size"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "2048",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u4e0a\u4e0b\u6587\u5927\u5c0f",
                            "en_US": "Enter the context size"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "completion_type",
                        "label": {
                            "zh_Hans": "\u8865\u5168\u7c7b\u578b",
                            "en_US": "Model type"
                        },
                        "type": "select",
                        "required": true,
                        "default": "chat",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u8865\u5168\u6a21\u578b",
                                    "en_US": "Completion model"
                                },
                                "value": "completion",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5bf9\u8bdd\u6a21\u578b",
                                    "en_US": "Chat model"
                                },
                                "value": "chat",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u8865\u5168\u7c7b\u578b",
                            "en_US": "Enter the completion type"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "stream",
                        "label": {
                            "zh_Hans": "\u6d41\u5f0f\u8f93\u51fa",
                            "en_US": "Stream output"
                        },
                        "type": "select",
                        "required": true,
                        "default": "true",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u662f",
                                    "en_US": "Yes"
                                },
                                "value": "true",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5426",
                                    "en_US": "No"
                                },
                                "value": "false",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u662f\u5426\u652f\u6301\u6d41\u5f0f\u8f93\u51fa",
                            "en_US": "Whether to support stream output"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "zhipuai",
            "label": {
                "zh_Hans": "\u667a\u8c31 AI",
                "en_US": "ZHIPU AI"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/zhipuai/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/zhipuai/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/zhipuai/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/zhipuai/icon_large/en_US"
            },
            "background": "#EFF1FE",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece\u667a\u8c31 AI \u83b7\u53d6 API Key",
                    "en_US": "Get your API key from ZHIPU AI"
                },
                "url": {
                    "zh_Hans": "https://open.bigmodel.cn/usercenter/apikeys",
                    "en_US": "https://open.bigmodel.cn/usercenter/apikeys"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "APIKey",
                            "en_US": "APIKey"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 APIKey",
                            "en_US": "Enter your APIKey"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "baichuan",
            "label": {
                "zh_Hans": "Baichuan",
                "en_US": "Baichuan"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/baichuan/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/baichuan/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/baichuan/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/baichuan/icon_large/en_US"
            },
            "background": "#FFF6F2",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece\u767e\u5ddd\u667a\u80fd\u83b7\u53d6\u60a8\u7684 API Key",
                    "en_US": "Get your API Key from BAICHUAN AI"
                },
                "url": {
                    "zh_Hans": "https://www.baichuan-ai.com",
                    "en_US": "https://www.baichuan-ai.com"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "secret_key",
                        "label": {
                            "zh_Hans": "Secret Key",
                            "en_US": "Secret Key"
                        },
                        "type": "secret-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Secret Key",
                            "en_US": "Enter your Secret Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "spark",
            "label": {
                "zh_Hans": "\u8baf\u98de\u661f\u706b",
                "en_US": "iFLYTEK SPARK"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/spark/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/spark/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/spark/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/spark/icon_large/en_US"
            },
            "background": "#EBF8FF",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece\u8baf\u98de\u661f\u706b\u83b7\u53d6 API Keys",
                    "en_US": "Get your API key from iFLYTEK SPARK"
                },
                "url": {
                    "zh_Hans": "https://www.xfyun.cn/solutions/xinghuoAPI",
                    "en_US": "https://www.xfyun.cn/solutions/xinghuoAPI"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "app_id",
                        "label": {
                            "zh_Hans": "APPID",
                            "en_US": "APPID"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 APPID",
                            "en_US": "Enter your APPID"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "api_secret",
                        "label": {
                            "zh_Hans": "APISecret",
                            "en_US": "APISecret"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 APISecret",
                            "en_US": "Enter your APISecret"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "APIKey",
                            "en_US": "APIKey"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 APIKey",
                            "en_US": "Enter your APIKey"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "minimax",
            "label": {
                "zh_Hans": "Minimax",
                "en_US": "Minimax"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/minimax/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/minimax/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/minimax/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/minimax/icon_large/en_US"
            },
            "background": "#FFEFEF",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece Minimax \u83b7\u53d6\u60a8\u7684 API Key",
                    "en_US": "Get your API Key from Minimax"
                },
                "url": {
                    "zh_Hans": "https://api.minimax.chat/user-center/basic-information/interface-key",
                    "en_US": "https://api.minimax.chat/user-center/basic-information/interface-key"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "minimax_api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "minimax_group_id",
                        "label": {
                            "zh_Hans": "Group ID",
                            "en_US": "Group ID"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Group ID",
                            "en_US": "Enter your group ID"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "tongyi",
            "label": {
                "zh_Hans": "\u901a\u4e49\u5343\u95ee",
                "en_US": "TONGYI"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/tongyi/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/tongyi/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/tongyi/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/tongyi/icon_large/en_US"
            },
            "background": "#EFF1FE",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece\u963f\u91cc\u4e91\u83b7\u53d6 API Key",
                    "en_US": "Get your API key from AliCloud"
                },
                "url": {
                    "zh_Hans": "https://dashscope.console.aliyun.com/api-key_management",
                    "en_US": "https://dashscope.console.aliyun.com/api-key_management"
                }
            },
            "supported_model_types": [
                "llm",
                "tts",
                "text-embedding"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "dashscope_api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "wenxin",
            "label": {
                "zh_Hans": "\u6587\u5fc3\u4e00\u8a00",
                "en_US": "WenXin"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/wenxin/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/wenxin/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/wenxin/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/wenxin/icon_large/en_US"
            },
            "background": "#E8F5FE",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece\u6587\u5fc3\u4e00\u8a00\u83b7\u53d6\u60a8\u7684 API Key",
                    "en_US": "Get your API Key from WenXin"
                },
                "url": {
                    "zh_Hans": "https://cloud.baidu.com/wenxin.html",
                    "en_US": "https://cloud.baidu.com/wenxin.html"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "secret_key",
                        "label": {
                            "zh_Hans": "Secret Key",
                            "en_US": "Secret Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Secret Key",
                            "en_US": "Enter your Secret Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "moonshot",
            "label": {
                "zh_Hans": "\u6708\u4e4b\u6697\u9762",
                "en_US": "Moonshot"
            },
            "description": {
                "zh_Hans": "Moonshot \u63d0\u4f9b\u7684\u6a21\u578b\uff0c\u4f8b\u5982 moonshot-v1-8k\u3001moonshot-v1-32k \u548c moonshot-v1-128k\u3002",
                "en_US": "Models provided by Moonshot, such as moonshot-v1-8k, moonshot-v1-32k, and moonshot-v1-128k."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/moonshot/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/moonshot/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/moonshot/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/moonshot/icon_large/en_US"
            },
            "background": "#FFFFFF",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece Moonshot \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from Moonshot"
                },
                "url": {
                    "zh_Hans": "https://platform.moonshot.cn/console/api-keys",
                    "en_US": "https://platform.moonshot.cn/console/api-keys"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model",
                "customizable-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "endpoint_url",
                        "label": {
                            "zh_Hans": "API Base",
                            "en_US": "API Base"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Base URL, \u5982\uff1ahttps://api.moonshot.cn/v1",
                            "en_US": "Base URL, e.g. https://api.moonshot.cn/v1"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Model context size"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Enter your Model context size"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "max_tokens",
                        "label": {
                            "zh_Hans": "\u6700\u5927 token \u4e0a\u9650",
                            "en_US": "Upper bound for max tokens"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "function_calling_type",
                        "label": {
                            "zh_Hans": "Function calling",
                            "en_US": "Function calling"
                        },
                        "type": "select",
                        "required": false,
                        "default": "no_call",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u4e0d\u652f\u6301",
                                    "en_US": "Not supported"
                                },
                                "value": "no_call",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "Tool Call",
                                    "en_US": "Tool Call"
                                },
                                "value": "tool_call",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "tencent",
            "label": {
                "zh_Hans": "\u817e\u8baf\u4e91",
                "en_US": "Tencent"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/tencent/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/tencent/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/tencent/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/tencent/icon_large/en_US"
            },
            "background": "#E5E7EB",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece\u817e\u8baf\u4e91\u83b7\u53d6 API Key",
                    "en_US": "Get your API key from Tencent AI"
                },
                "url": {
                    "zh_Hans": "https://cloud.tencent.com/product/asr",
                    "en_US": "https://cloud.tencent.com/product/asr"
                }
            },
            "supported_model_types": [
                "speech2text"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "app_id",
                        "label": {
                            "zh_Hans": "APPID",
                            "en_US": "APPID"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u817e\u8baf\u8bed\u97f3\u8bc6\u522b\u670d\u52a1\u7684 APPID",
                            "en_US": "Enter the APPID of your Tencent Cloud ASR service"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "secret_id",
                        "label": {
                            "zh_Hans": "SecretId",
                            "en_US": "SecretId"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u817e\u8baf\u8bed\u97f3\u8bc6\u522b\u670d\u52a1\u7684 SecretId",
                            "en_US": "Enter the SecretId of your Tencent Cloud ASR service"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "secret_key",
                        "label": {
                            "zh_Hans": "SecretKey",
                            "en_US": "SecretKey"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u817e\u8baf\u8bed\u97f3\u8bc6\u522b\u670d\u52a1\u7684 SecretKey",
                            "en_US": "Enter the SecretKey of your Tencent Cloud ASR service"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "jina",
            "label": {
                "zh_Hans": "Jina",
                "en_US": "Jina"
            },
            "description": {
                "zh_Hans": "Embedding and Rerank Model Supported",
                "en_US": "Embedding and Rerank Model Supported"
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/jina/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/jina/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/jina/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/jina/icon_large/en_US"
            },
            "background": "#EFFDFD",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece Jina \u83b7\u53d6 API Key",
                    "en_US": "Get your API key from Jina AI"
                },
                "url": {
                    "zh_Hans": "https://jina.ai/",
                    "en_US": "https://jina.ai/"
                }
            },
            "supported_model_types": [
                "text-embedding",
                "rerank"
            ],
            "configurate_methods": [
                "predefined-model",
                "customizable-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "base_url",
                        "label": {
                            "zh_Hans": "\u670d\u52a1\u5668 URL",
                            "en_US": "Base URL"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "https://api.jina.ai/v1",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Base URL, e.g. https://api.jina.ai/v1",
                            "en_US": "Base URL, e.g. https://api.jina.ai/v1"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u4e0a\u4e0b\u6587\u5927\u5c0f",
                            "en_US": "Context size"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": "8192",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8f93\u5165\u4e0a\u4e0b\u6587\u5927\u5c0f",
                            "en_US": "Enter context size"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "chatglm",
            "label": {
                "zh_Hans": "ChatGLM",
                "en_US": "ChatGLM"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/chatglm/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/chatglm/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/chatglm/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/chatglm/icon_large/en_US"
            },
            "background": "#F4F7FF",
            "help": {
                "title": {
                    "zh_Hans": "\u90e8\u7f72\u60a8\u7684\u672c\u5730 ChatGLM",
                    "en_US": "Deploy ChatGLM to your local"
                },
                "url": {
                    "zh_Hans": "https://github.com/THUDM/ChatGLM3",
                    "en_US": "https://github.com/THUDM/ChatGLM3"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_base",
                        "label": {
                            "zh_Hans": "API URL",
                            "en_US": "API URL"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API URL",
                            "en_US": "Enter your API URL"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "yi",
            "label": {
                "zh_Hans": "\u96f6\u4e00\u4e07\u7269",
                "en_US": "01.AI"
            },
            "description": {
                "zh_Hans": "\u96f6\u4e00\u4e07\u7269\u63d0\u4f9b\u7684\u6a21\u578b\uff0c\u4f8b\u5982 yi-34b-chat \u548c yi-vl-plus\u3002",
                "en_US": "Models provided by 01.AI, such as yi-34b-chat and yi-vl-plus."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/yi/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/yi/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/yi/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/yi/icon_large/en_US"
            },
            "background": "#E9F1EC",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece\u96f6\u4e00\u4e07\u7269\u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from 01.ai"
                },
                "url": {
                    "zh_Hans": "https://platform.lingyiwanwu.com/apikeys",
                    "en_US": "https://platform.lingyiwanwu.com/apikeys"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "endpoint_url",
                        "label": {
                            "zh_Hans": "\u81ea\u5b9a\u4e49 API endpoint \u5730\u5740",
                            "en_US": "Custom API endpoint URL"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Base URL, e.g. https://api.lingyiwanwu.com/v1",
                            "en_US": "Base URL, e.g. https://api.lingyiwanwu.com/v1"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "openllm",
            "label": {
                "zh_Hans": "OpenLLM",
                "en_US": "OpenLLM"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openllm/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openllm/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openllm/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/openllm/icon_large/en_US"
            },
            "background": "#F9FAFB",
            "help": {
                "title": {
                    "zh_Hans": "\u5982\u4f55\u90e8\u7f72 OpenLLM",
                    "en_US": "How to deploy OpenLLM"
                },
                "url": {
                    "zh_Hans": "https://github.com/bentoml/OpenLLM",
                    "en_US": "https://github.com/bentoml/OpenLLM"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "server_url",
                        "label": {
                            "zh_Hans": "\u670d\u52a1\u5668URL",
                            "en_US": "Server url"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165OpenLLM\u7684\u670d\u52a1\u5668\u5730\u5740\uff0c\u5982 http://192.168.1.100:3000",
                            "en_US": "Enter the url of your OpenLLM, e.g. http://192.168.1.100:3000"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "localai",
            "label": {
                "zh_Hans": "LocalAI",
                "en_US": "LocalAI"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/localai/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/localai/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/localai/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/localai/icon_large/en_US"
            },
            "background": "#F3F4F6",
            "help": {
                "title": {
                    "zh_Hans": "\u5982\u4f55\u90e8\u7f72 LocalAI",
                    "en_US": "How to deploy LocalAI"
                },
                "url": {
                    "zh_Hans": "https://github.com/go-skynet/LocalAI",
                    "en_US": "https://github.com/go-skynet/LocalAI"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding",
                "rerank",
                "speech2text"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "completion_type",
                        "label": {
                            "zh_Hans": "Completion type",
                            "en_US": "Completion type"
                        },
                        "type": "select",
                        "required": false,
                        "default": "chat_completion",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u8865\u5168",
                                    "en_US": "Completion"
                                },
                                "value": "completion",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5bf9\u8bdd",
                                    "en_US": "ChatCompletion"
                                },
                                "value": "chat_completion",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u9009\u62e9\u5bf9\u8bdd\u7c7b\u578b",
                            "en_US": "Select completion type"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "server_url",
                        "label": {
                            "zh_Hans": "\u670d\u52a1\u5668URL",
                            "en_US": "Server url"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165LocalAI\u7684\u670d\u52a1\u5668\u5730\u5740\uff0c\u5982 http://192.168.1.100:8080",
                            "en_US": "Enter the url of your LocalAI, e.g. http://192.168.1.100:8080"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u4e0a\u4e0b\u6587\u5927\u5c0f",
                            "en_US": "Context size"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8f93\u5165\u4e0a\u4e0b\u6587\u5927\u5c0f",
                            "en_US": "Enter context size"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "volcengine_maas",
            "label": {
                "zh_Hans": "Volcengine",
                "en_US": "Volcengine"
            },
            "description": {
                "zh_Hans": "\u706b\u5c71\u65b9\u821f\u63d0\u4f9b\u7684\u6a21\u578b\uff0c\u4f8b\u5982 Doubao-pro-4k\u3001Doubao-pro-32k \u548c Doubao-pro-128k\u3002",
                "en_US": "Volcengine Ark models."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/volcengine_maas/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/volcengine_maas/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/volcengine_maas/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/volcengine_maas/icon_large/en_US"
            },
            "background": "#F9FAFB",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece\u706b\u5c71\u5f15\u64ce\u63a7\u5236\u53f0\u83b7\u53d6\u60a8\u7684 Access Key \u548c Secret Access Key",
                    "en_US": "Get your Access Key and Secret Access Key from Volcengine Console"
                },
                "url": {
                    "zh_Hans": "https://console.volcengine.com/iam/keymanage/",
                    "en_US": "https://console.volcengine.com/iam/keymanage/"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your Model Name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "volc_access_key_id",
                        "label": {
                            "zh_Hans": "Access Key",
                            "en_US": "Access Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8f93\u5165\u60a8\u7684 Access Key",
                            "en_US": "Enter your Access Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "volc_secret_access_key",
                        "label": {
                            "zh_Hans": "Secret Access Key",
                            "en_US": "Secret Access Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8f93\u5165\u60a8\u7684 Secret Access Key",
                            "en_US": "Enter your Secret Access Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "volc_region",
                        "label": {
                            "zh_Hans": "\u706b\u5c71\u5f15\u64ce\u5730\u57df",
                            "en_US": "Volcengine Region"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "cn-beijing",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8f93\u5165\u706b\u5c71\u5f15\u64ce\u5730\u57df",
                            "en_US": "Enter Volcengine Region"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "api_endpoint_host",
                        "label": {
                            "zh_Hans": "API Endpoint Host",
                            "en_US": "API Endpoint Host"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "maas-api.ml-platform-cn-beijing.volces.com",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8f93\u5165 API Endpoint Host",
                            "en_US": "Enter your API Endpoint Host"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "endpoint_id",
                        "label": {
                            "zh_Hans": "Endpoint ID",
                            "en_US": "Endpoint ID"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8f93\u5165\u60a8\u7684 Endpoint ID",
                            "en_US": "Enter your Endpoint ID"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "base_model_name",
                        "label": {
                            "zh_Hans": "\u57fa\u7840\u6a21\u578b",
                            "en_US": "Base Model"
                        },
                        "type": "select",
                        "required": true,
                        "default": null,
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "Doubao-pro-4k",
                                    "en_US": "Doubao-pro-4k"
                                },
                                "value": "Doubao-pro-4k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Doubao-lite-4k",
                                    "en_US": "Doubao-lite-4k"
                                },
                                "value": "Doubao-lite-4k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Doubao-pro-32k",
                                    "en_US": "Doubao-pro-32k"
                                },
                                "value": "Doubao-pro-32k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Doubao-lite-32k",
                                    "en_US": "Doubao-lite-32k"
                                },
                                "value": "Doubao-lite-32k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Doubao-pro-128k",
                                    "en_US": "Doubao-pro-128k"
                                },
                                "value": "Doubao-pro-128k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Doubao-lite-128k",
                                    "en_US": "Doubao-lite-128k"
                                },
                                "value": "Doubao-lite-128k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Llama3-8B",
                                    "en_US": "Llama3-8B"
                                },
                                "value": "Llama3-8B",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Llama3-70B",
                                    "en_US": "Llama3-70B"
                                },
                                "value": "Llama3-70B",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Moonshot-v1-8k",
                                    "en_US": "Moonshot-v1-8k"
                                },
                                "value": "Moonshot-v1-8k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Moonshot-v1-32k",
                                    "en_US": "Moonshot-v1-32k"
                                },
                                "value": "Moonshot-v1-32k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Moonshot-v1-128k",
                                    "en_US": "Moonshot-v1-128k"
                                },
                                "value": "Moonshot-v1-128k",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "GLM3-130B",
                                    "en_US": "GLM3-130B"
                                },
                                "value": "GLM3-130B",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "GLM3-130B-Fin",
                                    "en_US": "GLM3-130B-Fin"
                                },
                                "value": "GLM3-130B-Fin",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Mistral-7B",
                                    "en_US": "Mistral-7B"
                                },
                                "value": "Mistral-7B",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "llm"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "Doubao-embedding",
                                    "en_US": "Doubao-embedding"
                                },
                                "value": "Doubao-embedding",
                                "show_on": [
                                    {
                                        "variable": "__model_type",
                                        "value": "text-embedding"
                                    }
                                ]
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u81ea\u5b9a\u4e49",
                                    "en_US": "Custom"
                                },
                                "value": "Custom",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "mode",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u7c7b\u578b",
                            "en_US": "Completion Mode"
                        },
                        "type": "select",
                        "required": true,
                        "default": "chat",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u8865\u5168",
                                    "en_US": "Completion"
                                },
                                "value": "completion",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5bf9\u8bdd",
                                    "en_US": "Chat"
                                },
                                "value": "chat",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u9009\u62e9\u5bf9\u8bdd\u7c7b\u578b",
                            "en_US": "Select Completion Mode"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            },
                            {
                                "variable": "base_model_name",
                                "value": "Custom"
                            }
                        ]
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Model Context Size"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8f93\u5165\u60a8\u7684\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Enter your Model Context Size"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "base_model_name",
                                "value": "Custom"
                            }
                        ]
                    },
                    {
                        "variable": "max_tokens",
                        "label": {
                            "zh_Hans": "\u6700\u5927 token \u4e0a\u9650",
                            "en_US": "Upper Bound for Max Tokens"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8f93\u5165\u60a8\u7684\u6a21\u578b\u6700\u5927 token \u4e0a\u9650",
                            "en_US": "Enter your model Upper Bound for Max Tokens"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            },
                            {
                                "variable": "base_model_name",
                                "value": "Custom"
                            }
                        ]
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "openai_api_compatible",
            "label": {
                "zh_Hans": "OpenAI-API-compatible",
                "en_US": "OpenAI-API-compatible"
            },
            "description": {
                "zh_Hans": "\u517c\u5bb9 OpenAI API \u7684\u6a21\u578b\u4f9b\u5e94\u5546\uff0c\u4f8b\u5982 LM Studio \u3002",
                "en_US": "Model providers compatible with OpenAI's API standard, such as LM Studio."
            },
            "icon_small": null,
            "icon_large": null,
            "background": null,
            "help": null,
            "supported_model_types": [
                "llm",
                "text-embedding",
                "speech2text"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u5168\u79f0",
                        "en_US": "Enter full model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "endpoint_url",
                        "label": {
                            "zh_Hans": "API endpoint URL",
                            "en_US": "API endpoint URL"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Base URL, e.g. https://api.openai.com/v1",
                            "en_US": "Base URL, e.g. https://api.openai.com/v1"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "mode",
                        "label": {
                            "zh_Hans": "Completion mode",
                            "en_US": "Completion mode"
                        },
                        "type": "select",
                        "required": false,
                        "default": "chat",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u8865\u5168",
                                    "en_US": "Completion"
                                },
                                "value": "completion",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5bf9\u8bdd",
                                    "en_US": "Chat"
                                },
                                "value": "chat",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u9009\u62e9\u5bf9\u8bdd\u7c7b\u578b",
                            "en_US": "Select completion mode"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Model context size"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Enter your Model context size"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Model context size"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Enter your Model context size"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "text-embedding"
                            }
                        ]
                    },
                    {
                        "variable": "max_tokens_to_sample",
                        "label": {
                            "zh_Hans": "\u6700\u5927 token \u4e0a\u9650",
                            "en_US": "Upper bound for max tokens"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "4096",
                        "options": null,
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "function_calling_type",
                        "label": {
                            "zh_Hans": "Function calling",
                            "en_US": "Function calling"
                        },
                        "type": "select",
                        "required": false,
                        "default": "no_call",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "Function Call",
                                    "en_US": "Function Call"
                                },
                                "value": "function_call",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "Tool Call",
                                    "en_US": "Tool Call"
                                },
                                "value": "tool_call",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e0d\u652f\u6301",
                                    "en_US": "Not Support"
                                },
                                "value": "no_call",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "stream_function_calling",
                        "label": {
                            "zh_Hans": "Stream function calling",
                            "en_US": "Stream function calling"
                        },
                        "type": "select",
                        "required": false,
                        "default": "not_supported",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u652f\u6301",
                                    "en_US": "Support"
                                },
                                "value": "supported",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e0d\u652f\u6301",
                                    "en_US": "Not Support"
                                },
                                "value": "not_supported",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "vision_support",
                        "label": {
                            "zh_Hans": "Vision \u652f\u6301",
                            "en_US": "Vision Support"
                        },
                        "type": "select",
                        "required": false,
                        "default": "no_support",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u652f\u6301",
                                    "en_US": "Support"
                                },
                                "value": "support",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e0d\u652f\u6301",
                                    "en_US": "Not Support"
                                },
                                "value": "no_support",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "stream_mode_delimiter",
                        "label": {
                            "zh_Hans": "\u6d41\u6a21\u5f0f\u8fd4\u56de\u7ed3\u679c\u7684\u5206\u9694\u7b26",
                            "en_US": "Delimiter for streaming results"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "\\n\\n",
                        "options": null,
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "deepseek",
            "label": {
                "zh_Hans": "\u6df1\u5ea6\u6c42\u7d22",
                "en_US": "deepseek"
            },
            "description": {
                "zh_Hans": "\u6df1\u5ea6\u6c42\u7d22\u63d0\u4f9b\u7684\u6a21\u578b\uff0c\u4f8b\u5982 deepseek-chat\u3001deepseek-coder \u3002",
                "en_US": "Models provided by deepseek, such as deepseek-chat\u3001deepseek-coder."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/deepseek/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/deepseek/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/deepseek/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/deepseek/icon_large/en_US"
            },
            "background": "#c0cdff",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece\u6df1\u5ea6\u6c42\u7d22\u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from deepseek"
                },
                "url": {
                    "zh_Hans": "https://platform.deepseek.com/api_keys",
                    "en_US": "https://platform.deepseek.com/api_keys"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "endpoint_url",
                        "label": {
                            "zh_Hans": "\u81ea\u5b9a\u4e49 API endpoint \u5730\u5740",
                            "en_US": "Custom API endpoint URL"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Base URL, e.g. https://api.deepseek.com/v1 or https://api.deepseek.com",
                            "en_US": "Base URL, e.g. https://api.deepseek.com/v1 or https://api.deepseek.com"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "hunyuan",
            "label": {
                "zh_Hans": "\u817e\u8baf\u6df7\u5143",
                "en_US": "Hunyuan"
            },
            "description": {
                "zh_Hans": "\u817e\u8baf\u6df7\u5143\u63d0\u4f9b\u7684\u6a21\u578b\uff0c\u4f8b\u5982 hunyuan-standard\u3001 hunyuan-standard-256k, hunyuan-pro \u548c hunyuan-lite\u3002",
                "en_US": "Models provided by Tencent Hunyuan, such as hunyuan-standard, hunyuan-standard-256k, hunyuan-pro and hunyuan-lite."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/hunyuan/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/hunyuan/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/hunyuan/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/hunyuan/icon_large/en_US"
            },
            "background": "#F6F7F7",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece\u817e\u8baf\u6df7\u5143\u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from Tencent Hunyuan"
                },
                "url": {
                    "zh_Hans": "https://console.cloud.tencent.com/cam/capi",
                    "en_US": "https://console.cloud.tencent.com/cam/capi"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "secret_id",
                        "label": {
                            "zh_Hans": "Secret ID",
                            "en_US": "Secret ID"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Secret ID",
                            "en_US": "Enter your Secret ID"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "secret_key",
                        "label": {
                            "zh_Hans": "Secret Key",
                            "en_US": "Secret Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Secret Key",
                            "en_US": "Enter your Secret Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "siliconflow",
            "label": {
                "zh_Hans": "\u7845\u57fa\u6d41\u52a8",
                "en_US": "SiliconFlow"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/siliconflow/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/siliconflow/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/siliconflow/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/siliconflow/icon_large/en_US"
            },
            "background": "#ffecff",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece SiliconFlow \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from SiliconFlow"
                },
                "url": {
                    "zh_Hans": "https://cloud.siliconflow.cn/account/ak",
                    "en_US": "https://cloud.siliconflow.cn/account/ak"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding",
                "rerank",
                "speech2text"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "perfxcloud",
            "label": {
                "zh_Hans": "PerfXCloud",
                "en_US": "PerfXCloud"
            },
            "description": {
                "zh_Hans": "PerfXCloud\uff08\u6f8e\u5cf0\u79d1\u6280\uff09\u4e3a\u5f00\u53d1\u8005\u548c\u4f01\u4e1a\u91cf\u8eab\u6253\u9020\u7684AI\u5f00\u53d1\u548c\u90e8\u7f72\u5e73\u53f0\uff0c\u63d0\u4f9b\u591a\u79cd\u6a21\u578b\u7684\u7684\u63a8\u7406\u80fd\u529b\u3002",
                "en_US": "PerfXCloud (Pengfeng Technology) is an AI development and deployment platform tailored for developers and enterprises, providing reasoning capabilities for multiple models."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/perfxcloud/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/perfxcloud/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/perfxcloud/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/perfxcloud/icon_large/en_US"
            },
            "background": "#e3f0ff",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece PerfXCloud \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from PerfXCloud"
                },
                "url": {
                    "zh_Hans": "https://cloud.perfxlab.cn/panel/token",
                    "en_US": "https://cloud.perfxlab.cn/panel/token"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "endpoint_url",
                        "label": {
                            "zh_Hans": "\u81ea\u5b9a\u4e49 API endpoint \u5730\u5740",
                            "en_US": "Custom API endpoint URL"
                        },
                        "type": "text-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "Base URL, e.g. https://cloud.perfxlab.cn/v1",
                            "en_US": "Base URL, e.g. https://cloud.perfxlab.cn/v1"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "zhinao",
            "label": {
                "zh_Hans": "360 \u667a\u8111",
                "en_US": "360 AI"
            },
            "description": {
                "zh_Hans": "360 \u667a\u8111\u63d0\u4f9b\u7684\u6a21\u578b\u3002",
                "en_US": "Models provided by 360 AI."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/zhinao/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/zhinao/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/zhinao/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/zhinao/icon_large/en_US"
            },
            "background": "#e3f0ff",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece360 \u667a\u8111\u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from 360 AI."
                },
                "url": {
                    "zh_Hans": "https://ai.360.com/platform/keys",
                    "en_US": "https://ai.360.com/platform/keys"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "sagemaker",
            "label": {
                "zh_Hans": "Sagemaker",
                "en_US": "Sagemaker"
            },
            "description": {
                "zh_Hans": "Sagemaker\u4e0a\u7684\u79c1\u6709\u5316\u90e8\u7f72\u7684\u6a21\u578b",
                "en_US": "Customized model on Sagemaker"
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/sagemaker/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/sagemaker/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/sagemaker/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/sagemaker/icon_large/en_US"
            },
            "background": "#ECE9E3",
            "help": {
                "title": {
                    "zh_Hans": "\u5982\u4f55\u5728Sagemaker\u4e0a\u7684\u79c1\u6709\u5316\u90e8\u7f72\u7684\u6a21\u578b",
                    "en_US": "How to deploy customized model on Sagemaker"
                },
                "url": {
                    "zh_Hans": "https://github.com/aws-samples/gofy-aws-tool/blob/main/README_ZH.md#%E5%A6%82%E4%BD%95%E9%83%A8%E7%BD%B2sagemaker%E6%8E%A8%E7%90%86%E7%AB%AF%E7%82%B9",
                    "en_US": "https://github.com/aws-samples/gofy-aws-tool/blob/main/README.md#how-to-deploy-sagemaker-endpoint"
                }
            },
            "supported_model_types": [
                "llm",
                "text-embedding",
                "rerank"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "mode",
                        "label": {
                            "zh_Hans": "Completion mode",
                            "en_US": "Completion mode"
                        },
                        "type": "select",
                        "required": false,
                        "default": "chat",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u8865\u5168",
                                    "en_US": "Completion"
                                },
                                "value": "completion",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u5bf9\u8bdd",
                                    "en_US": "Chat"
                                },
                                "value": "chat",
                                "show_on": []
                            }
                        ],
                        "placeholder": {
                            "zh_Hans": "\u9009\u62e9\u5bf9\u8bdd\u7c7b\u578b",
                            "en_US": "Select completion mode"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "sagemaker_endpoint",
                        "label": {
                            "zh_Hans": "sagemaker endpoint",
                            "en_US": "sagemaker endpoint"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8bf7\u8f93\u51fa\u4f60\u7684Sagemaker\u63a8\u7406\u7aef\u70b9",
                            "en_US": "Enter your Sagemaker Inference endpoint"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "aws_access_key_id",
                        "label": {
                            "zh_Hans": "Access Key (\u5982\u679c\u672a\u63d0\u4f9b\uff0c\u51ed\u8bc1\u5c06\u4ece\u8fd0\u884c\u73af\u5883\u4e2d\u83b7\u53d6\u3002)",
                            "en_US": "Access Key (If not provided, credentials are obtained from the running environment.)"
                        },
                        "type": "secret-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Access Key",
                            "en_US": "Enter your Access Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "aws_secret_access_key",
                        "label": {
                            "zh_Hans": "Secret Access Key",
                            "en_US": "Secret Access Key"
                        },
                        "type": "secret-input",
                        "required": false,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Secret Access Key",
                            "en_US": "Enter your Secret Access Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "aws_region",
                        "label": {
                            "zh_Hans": "AWS \u5730\u533a",
                            "en_US": "AWS Region"
                        },
                        "type": "select",
                        "required": false,
                        "default": "us-east-1",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u7f8e\u56fd\u4e1c\u90e8 (\u5f17\u5409\u5c3c\u4e9a\u5317\u90e8)",
                                    "en_US": "US East (N. Virginia)"
                                },
                                "value": "us-east-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u7f8e\u56fd\u897f\u90e8 (\u4fc4\u52d2\u5188\u5dde)",
                                    "en_US": "US West (Oregon)"
                                },
                                "value": "us-west-2",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e9a\u592a\u5730\u533a (\u65b0\u52a0\u5761)",
                                    "en_US": "Asia Pacific (Singapore)"
                                },
                                "value": "ap-southeast-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e9a\u592a\u5730\u533a (\u4e1c\u4eac)",
                                    "en_US": "Asia Pacific (Tokyo)"
                                },
                                "value": "ap-northeast-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u6b27\u6d32 (\u6cd5\u5170\u514b\u798f)",
                                    "en_US": "Europe (Frankfurt)"
                                },
                                "value": "eu-central-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "AWS GovCloud (US-West)",
                                    "en_US": "AWS GovCloud (US-West)"
                                },
                                "value": "us-gov-west-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e9a\u592a\u5730\u533a (\u6089\u5c3c)",
                                    "en_US": "Asia Pacific (Sydney)"
                                },
                                "value": "ap-southeast-2",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e2d\u56fd\u5317\u4eac (cn-north-1)",
                                    "en_US": "AWS Beijing (cn-north-1)"
                                },
                                "value": "cn-north-1",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "\u4e2d\u56fd\u5b81\u590f (cn-northwest-1)",
                                    "en_US": "AWS Ningxia (cn-northwest-1)"
                                },
                                "value": "cn-northwest-1",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "novita",
            "label": {
                "zh_Hans": "novita.ai",
                "en_US": "novita.ai"
            },
            "description": {
                "zh_Hans": "\u9002\u914d\u591a\u79cd\u6d77\u5916\u5e94\u7528\u573a\u666f\u7684\u9ad8\u6027\u4ef7\u6bd4 LLM API",
                "en_US": "An LLM API that matches various application scenarios with high cost-effectiveness."
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/novita/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/novita/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/novita/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/novita/icon_large/en_US"
            },
            "background": "#eadeff",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece novita.ai \u83b7\u53d6 API Key",
                    "en_US": "Get your API key from novita.ai"
                },
                "url": {
                    "zh_Hans": "https://novita.ai/settings#key-management?utm_source=gofy&utm_medium=ch&utm_campaign=api",
                    "en_US": "https://novita.ai/settings#key-management?utm_source=gofy&utm_medium=ch&utm_campaign=api"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "leptonai",
            "label": {
                "zh_Hans": "Lepton AI",
                "en_US": "Lepton AI"
            },
            "description": null,
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/leptonai/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/leptonai/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/leptonai/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/leptonai/icon_large/en_US"
            },
            "background": "#F5F5F4",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece Lepton AI \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from Lepton AI"
                },
                "url": {
                    "zh_Hans": "https://dashboard.lepton.ai",
                    "en_US": "https://dashboard.lepton.ai"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": null,
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "huggingface_tei",
            "label": {
                "zh_Hans": "Text Embedding Inference",
                "en_US": "Text Embedding Inference"
            },
            "description": {
                "zh_Hans": "\u7528\u4e8e\u6587\u672c\u5d4c\u5165\u6a21\u578b\u7684\u8d85\u5feb\u901f\u63a8\u7406\u89e3\u51b3\u65b9\u6848\u3002",
                "en_US": "A blazing fast inference solution for text embeddings models."
            },
            "icon_small": null,
            "icon_large": null,
            "background": "#FFF8DC",
            "help": {
                "title": {
                    "zh_Hans": "\u5982\u4f55\u90e8\u7f72 Text Embedding Inference",
                    "en_US": "How to deploy Text Embedding Inference"
                },
                "url": {
                    "zh_Hans": "https://github.com/huggingface/text-embeddings-inference",
                    "en_US": "https://github.com/huggingface/text-embeddings-inference"
                }
            },
            "supported_model_types": [
                "text-embedding",
                "rerank"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "server_url",
                        "label": {
                            "zh_Hans": "\u670d\u52a1\u5668URL",
                            "en_US": "Server url"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165Text Embedding Inference\u7684\u670d\u52a1\u5668\u5730\u5740\uff0c\u5982 http://192.168.1.100:8080",
                            "en_US": "Enter the url of your Text Embedding Inference, e.g. http://192.168.1.100:8080"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "stepfun",
            "label": {
                "zh_Hans": "\u9636\u8dc3\u661f\u8fb0",
                "en_US": "Stepfun"
            },
            "description": {
                "zh_Hans": "\u9636\u8dc3\u661f\u8fb0\u63d0\u4f9b\u7684\u6a21\u578b\uff0c\u4f8b\u5982 step-1-8k\u3001step-1-32k\u3001step-1v-8k\u3001step-1v-32k\u3001step-1-128k \u548c step-1-256k\u3002",
                "en_US": "Models provided by stepfun, such as step-1-8k, step-1-32k\u3001step-1v-8k\u3001step-1v-32k, step-1-128k and step-1-256k"
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/stepfun/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/stepfun/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/stepfun/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/stepfun/icon_large/en_US"
            },
            "background": "#FFFFFF",
            "help": {
                "title": {
                    "zh_Hans": "\u4ece stepfun \u83b7\u53d6 API Key",
                    "en_US": "Get your API Key from stepfun"
                },
                "url": {
                    "zh_Hans": "https://platform.stepfun.com/interface-key",
                    "en_US": "https://platform.stepfun.com/interface-key"
                }
            },
            "supported_model_types": [
                "llm"
            ],
            "configurate_methods": [
                "predefined-model",
                "customizable-model"
            ],
            "provider_credential_schema": {
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 API Key",
                            "en_US": "Enter your API Key"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "context_size",
                        "label": {
                            "zh_Hans": "\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Model context size"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "8192",
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684\u6a21\u578b\u4e0a\u4e0b\u6587\u957f\u5ea6",
                            "en_US": "Enter your Model context size"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "max_tokens",
                        "label": {
                            "zh_Hans": "\u6700\u5927 token \u4e0a\u9650",
                            "en_US": "Upper bound for max tokens"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": "8192",
                        "options": null,
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "function_calling_type",
                        "label": {
                            "zh_Hans": "Function calling",
                            "en_US": "Function calling"
                        },
                        "type": "select",
                        "required": false,
                        "default": "no_call",
                        "options": [
                            {
                                "label": {
                                    "zh_Hans": "\u4e0d\u652f\u6301",
                                    "en_US": "Not supported"
                                },
                                "value": "no_call",
                                "show_on": []
                            },
                            {
                                "label": {
                                    "zh_Hans": "Tool Call",
                                    "en_US": "Tool Call"
                                },
                                "value": "tool_call",
                                "show_on": []
                            }
                        ],
                        "placeholder": null,
                        "max_length": 0,
                        "show_on": []
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        },
        {
            "provider": "azure_ai_studio",
            "label": {
                "zh_Hans": "Azure AI Studio",
                "en_US": "Azure AI Studio"
            },
            "description": {
                "zh_Hans": "Azure AI Studio",
                "en_US": "Azure AI Studio"
            },
            "icon_small": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/azure_ai_studio/icon_small/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/azure_ai_studio/icon_small/en_US"
            },
            "icon_large": {
                "zh_Hans": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/azure_ai_studio/icon_large/zh_Hans",
                "en_US": "http://172.26.228.36:5002/console/api/workspaces/current/model-providers/azure_ai_studio/icon_large/en_US"
            },
            "background": "#93c5fd",
            "help": {
                "title": {
                    "zh_Hans": "\u5982\u4f55\u5728Azure AI Studio\u4e0a\u7684\u79c1\u6709\u5316\u90e8\u7f72\u7684\u6a21\u578b",
                    "en_US": "How to deploy customized model on Azure AI Studio"
                },
                "url": {
                    "zh_Hans": "https://learn.microsoft.com/zh-cn/azure/ai-studio/how-to/deploy-models",
                    "en_US": "https://learn.microsoft.com/en-us/azure/ai-studio/how-to/deploy-models"
                }
            },
            "supported_model_types": [
                "llm",
                "rerank"
            ],
            "configurate_methods": [
                "customizable-model"
            ],
            "provider_credential_schema": null,
            "model_credential_schema": {
                "model": {
                    "label": {
                        "zh_Hans": "\u6a21\u578b\u540d\u79f0",
                        "en_US": "Model Name"
                    },
                    "placeholder": {
                        "zh_Hans": "\u8f93\u5165\u6a21\u578b\u540d\u79f0",
                        "en_US": "Enter your model name"
                    }
                },
                "credential_form_schemas": [
                    {
                        "variable": "endpoint",
                        "label": {
                            "zh_Hans": "Azure AI Studio Endpoint",
                            "en_US": "Azure AI Studio Endpoint"
                        },
                        "type": "text-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u8bf7\u8f93\u5165\u4f60\u7684Azure AI Studio\u63a8\u7406\u7aef\u70b9",
                            "en_US": "Enter your API Endpoint, eg: https://example.com"
                        },
                        "max_length": 0,
                        "show_on": []
                    },
                    {
                        "variable": "api_key",
                        "label": {
                            "zh_Hans": "API Key",
                            "en_US": "API Key"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Azure AI Studio API Key",
                            "en_US": "Enter your Azure AI Studio API Key"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "llm"
                            }
                        ]
                    },
                    {
                        "variable": "jwt_token",
                        "label": {
                            "zh_Hans": "JWT\u4ee4\u724c",
                            "en_US": "JWT Token"
                        },
                        "type": "secret-input",
                        "required": true,
                        "default": null,
                        "options": null,
                        "placeholder": {
                            "zh_Hans": "\u5728\u6b64\u8f93\u5165\u60a8\u7684 Azure AI Studio \u63a8\u7406 API Key",
                            "en_US": "Enter your Azure AI Studio JWT Token"
                        },
                        "max_length": 0,
                        "show_on": [
                            {
                                "variable": "__model_type",
                                "value": "rerank"
                            }
                        ]
                    }
                ]
            },
            "preferred_provider_type": "custom",
            "custom_configuration": {
                "status": "no-configure"
            },
            "system_configuration": {
                "enabled": false,
                "current_quota_type": null,
                "quota_configurations": []
            }
        }]`
)

type ModelProvideApi struct {
}

func (api *ModelProvideApi) ProviderList(c *gin.Context) {
	// data := make([]any, 0)
	// err := json.Unmarshal([]byte(currentProviders), &data)
	// if err != nil {
	// 	mlog.Errorf("json unmarshl failed:%v", err)
	// 	c.JSON(http.StatusNoContent, gin.H{"data": nil, "result": "failed", "code": 7})
	// } else {
	// 	c.JSON(http.StatusOK, gin.H{"data": data, "result": "success", "code": 0})
	// }
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	model_type := c.Query("model_type")
	if model_type != "" {
		if !modelruntimeenumtypes.ModelType(model_type).Valid() {
			mlog.Errorf("model_type=%s is invalid", model_type)
			response.InvalidArgErrorWithDetail(c, fmt.Sprintf("model_type=%s is invalid", model_type))
			return
		}
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetModelProviderList(c, &pbapi.GetModelProviderListRequest{TenantId: acc.CurrentTenantID(), ModelType: model_type})
		if err != nil {
			mlog.Errorf("remote call GetModelProviderList failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetModelProviderList failed",
			})
			return
		} else {
			mlog.Infof("remote call GetModelProviderList return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.ProviderListStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					var tmps []*servicesentities.ProviderResponse
					if err1 := json.Unmarshal([]byte(pbrsp.ProviderListStr), &tmps); err1 != nil {
						mlog.Errorf("json unmarshal provider_list_str=%s to ProviderResponse list failed:%v", pbrsp.ProviderListStr, err1)
						c.JSON(http.StatusBadRequest, gin.H{"result": "no providerr list", "code": 7})
					} else {
						c.JSON(http.StatusOK, gin.H{"data": tmps, "result": "success", "code": 0})
					}
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *ModelProvideApi) GetCredentials(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tenant_id := acc.CurrentTenantID()
	if tenant_id == "" {
		mlog.Errorf("missing tenant_id")
		response.Forbidden(c)
		return
	}
	provider := c.Param("provider")
	if provider == "" {
		mlog.Error("missing provider")
		response.InvalidArgErrorWithDetail(c, "missing provider")
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).GetModelProviderCredentials(c, &pbapi.GetModelProviderCredentialsRequest{TenantId: tenant_id, Provider: provider})
		if err != nil {
			mlog.Errorf("remote call GetModelProviderCredentials failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call GetModelProviderCredentials failed",
			})
			return
		} else {
			mlog.Infof("remote call GetModelProviderCredentials return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				if pbrsp.CredentialsDictStr == "" {
					mlog.Error("查询失败!")
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": "fail",
						"data":   "internal server error",
					})
				} else {
					var tmps map[string]any
					if err1 := json.Unmarshal([]byte(pbrsp.CredentialsDictStr), &tmps); err1 != nil {
						mlog.Errorf("json unmarshal provider_list_str=%s to dict failed:%v", pbrsp.CredentialsDictStr, err1)
						c.JSON(http.StatusBadRequest, gin.H{"result": "no credentials", "code": 7})
					} else {
						c.JSON(http.StatusOK, gin.H{"credentials": tmps})
					}
				}
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}

func (api *ModelProvideApi) Update(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tenant_id := acc.CurrentTenantID()
	if tenant_id == "" {
		mlog.Errorf("missing tenant_id")
		response.Forbidden(c)
		return
	}

	provider := c.Param("provider")
	if provider == "" {
		mlog.Error("missing provider")
		response.InvalidArgErrorWithDetail(c, "missing provider")
		return
	}
	type Req struct {
		Credentials map[string]any `json:"credentials"`
	}
	var req Req
	err := c.ShouldBindJSON(&req)
	if err != nil {
		mlog.Error("bind failed:", err)
		c.JSON(http.StatusBadRequest, map[string]any{"result": "invalid arg"})
		return
	}
	bindata, _ := json.Marshal(req.Credentials)
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).UpdateModelProvider(c, &pbapi.UpdateModelProviderRequest{TenantId: acc.CurrentTenantID(), Provider: provider, CredentialsDictStr: string(bindata)})
		if err != nil {
			mlog.Errorf("remote call UpdateModelProvider failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call UpdateModelProvider failed",
			})
			return
		} else {
			mlog.Infof("remote call UpdateModelProvider return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusCreated, gin.H{"result": "success"})
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
func (api *ModelProvideApi) Delete(c *gin.Context) {
	rawuser, _ := c.Get("current_user")
	acc := rawuser.(*models.Account)
	tenant_id := acc.CurrentTenantID()
	if tenant_id == "" {
		mlog.Errorf("missing tenant_id")
		response.Forbidden(c)
		return
	}

	provider := c.Param("provider")
	if provider == "" {
		mlog.Error("missing provider")
		response.InvalidArgErrorWithDetail(c, "missing provider")
		return
	}
	conn := cluster.Instance().GetRpcClientByModule("admin")
	if conn != nil {
		pbrsp, err := pbapi.NewAdminClient(conn).DelModelProvider(c, &pbapi.DelModelProviderRequest{TenantId: acc.CurrentTenantID(), Provider: provider})
		if err != nil {
			mlog.Errorf("remote call DelModelProvider failed:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": "fail",
				"data":   "remote call DelModelProvider failed",
			})
			return
		} else {
			mlog.Infof("remote call DelModelProvider return:%#v", pbrsp)
			if pbrsp.Exp != nil {
				response.PbHttpException(c, pbrsp.Exp)
			} else {
				c.JSON(http.StatusNoContent, gin.H{"result": "success"})
			}
			return
		}
	} else {
		mlog.Errorf("get rpc client failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": "fail",
			"data":   "internal server error",
		})
		return
	}
}
