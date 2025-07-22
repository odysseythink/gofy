SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for account_integrates
-- ----------------------------
DROP TABLE IF EXISTS `account_integrates`;
CREATE TABLE `account_integrates` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `account_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider` varchar(16) COLLATE utf8mb4_general_ci NOT NULL,
  `open_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `encrypted_token` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of account_integrates
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for accounts
-- ----------------------------
DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `email` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `password_salt` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `avatar` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `interface_language` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `interface_theme` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `timezone` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `last_login_at` timestamp NULL DEFAULT NULL,
  `last_login_ip` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `status` varchar(16) COLLATE utf8mb4_general_ci DEFAULT 'active',
  `initialized_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_active_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of accounts
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for alembic_version
-- ----------------------------
DROP TABLE IF EXISTS `alembic_version`;
CREATE TABLE `alembic_version` (
  `version_num` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`version_num`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of alembic_version
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for api_based_extensions
-- ----------------------------
DROP TABLE IF EXISTS `api_based_extensions`;
CREATE TABLE `api_based_extensions` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `api_endpoint` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `api_key` text COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of api_based_extensions
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for api_requests
-- ----------------------------
DROP TABLE IF EXISTS `api_requests`;
CREATE TABLE `api_requests` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `api_token_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `path` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `request` text COLLATE utf8mb4_general_ci,
  `response` text COLLATE utf8mb4_general_ci,
  `ip` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of api_requests
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for api_tokens
-- ----------------------------
DROP TABLE IF EXISTS `api_tokens`;
CREATE TABLE `api_tokens` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `type` varchar(16) COLLATE utf8mb4_general_ci NOT NULL,
  `token` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `last_used_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of api_tokens
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for app_annotation_hit_histories
-- ----------------------------
DROP TABLE IF EXISTS `app_annotation_hit_histories`;
CREATE TABLE `app_annotation_hit_histories` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `annotation_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `source` text COLLATE utf8mb4_general_ci NOT NULL,
  `question` text COLLATE utf8mb4_general_ci NOT NULL,
  `account_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `score` double NOT NULL DEFAULT '0',
  `message_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `annotation_question` text COLLATE utf8mb4_general_ci NOT NULL,
  `annotation_content` text COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of app_annotation_hit_histories
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for app_annotation_settings
-- ----------------------------
DROP TABLE IF EXISTS `app_annotation_settings`;
CREATE TABLE `app_annotation_settings` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `score_threshold` double NOT NULL DEFAULT '0',
  `collection_binding_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_user_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_user_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of app_annotation_settings
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for app_dataset_joins
-- ----------------------------
DROP TABLE IF EXISTS `app_dataset_joins`;
CREATE TABLE `app_dataset_joins` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of app_dataset_joins
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for app_model_configs
-- ----------------------------
DROP TABLE IF EXISTS `app_model_configs`;
CREATE TABLE `app_model_configs` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `model_id` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `configs` json DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `opening_statement` text COLLATE utf8mb4_general_ci,
  `suggested_questions` text COLLATE utf8mb4_general_ci,
  `suggested_questions_after_answer` text COLLATE utf8mb4_general_ci,
  `more_like_this` text COLLATE utf8mb4_general_ci,
  `model` text COLLATE utf8mb4_general_ci,
  `user_input_form` text COLLATE utf8mb4_general_ci,
  `pre_prompt` text COLLATE utf8mb4_general_ci,
  `agent_mode` text COLLATE utf8mb4_general_ci,
  `speech_to_text` text COLLATE utf8mb4_general_ci,
  `sensitive_word_avoidance` text COLLATE utf8mb4_general_ci,
  `retriever_resource` text COLLATE utf8mb4_general_ci,
  `dataset_query_variable` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `prompt_type` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'simple',
  `chat_prompt_config` text COLLATE utf8mb4_general_ci,
  `completion_prompt_config` text COLLATE utf8mb4_general_ci,
  `dataset_configs` text COLLATE utf8mb4_general_ci,
  `external_data_tools` text COLLATE utf8mb4_general_ci,
  `file_upload` text COLLATE utf8mb4_general_ci,
  `text_to_speech` text COLLATE utf8mb4_general_ci,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `updated_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of app_model_configs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for apps
-- ----------------------------
DROP TABLE IF EXISTS `apps`;
CREATE TABLE `apps` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `mode` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `icon` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `icon_background` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `app_model_config_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `status` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'normal',
  `enable_site` tinyint(1) NOT NULL,
  `enable_api` tinyint(1) NOT NULL,
  `api_rpm` int NOT NULL DEFAULT '0',
  `api_rph` int NOT NULL DEFAULT '0',
  `is_demo` tinyint(1) NOT NULL DEFAULT '0',
  `is_public` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_universal` tinyint(1) NOT NULL DEFAULT '0',
  `workflow_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  `tracing` text COLLATE utf8mb4_general_ci,
  `max_active_requests` int DEFAULT NULL,
  `icon_type` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `updated_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `use_icon_as_answer_icon` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of apps
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for celery_taskmeta
-- ----------------------------
DROP TABLE IF EXISTS `celery_taskmeta`;
CREATE TABLE `celery_taskmeta` (
  `id` int NOT NULL AUTO_INCREMENT,
  `task_id` varchar(155) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `status` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `result` blob,
  `date_done` timestamp NULL DEFAULT NULL,
  `traceback` text COLLATE utf8mb4_general_ci,
  `name` varchar(155) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `args` blob,
  `kwargs` blob,
  `worker` varchar(155) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `retries` int DEFAULT NULL,
  `queue` varchar(155) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of celery_taskmeta
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for celery_tasksetmeta
-- ----------------------------
DROP TABLE IF EXISTS `celery_tasksetmeta`;
CREATE TABLE `celery_tasksetmeta` (
  `id` int NOT NULL AUTO_INCREMENT,
  `taskset_id` varchar(155) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `result` blob,
  `date_done` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of celery_tasksetmeta
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for child_chunks
-- ----------------------------
DROP TABLE IF EXISTS `child_chunks`;
CREATE TABLE `child_chunks` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `document_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `segment_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `position` int NOT NULL,
  `content` text COLLATE utf8mb4_general_ci NOT NULL,
  `word_count` int NOT NULL,
  `index_node_id` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `index_node_hash` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `type` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'automatic',
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `indexing_at` timestamp NULL DEFAULT NULL,
  `completed_at` timestamp NULL DEFAULT NULL,
  `error` text COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of child_chunks
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for conversations
-- ----------------------------
DROP TABLE IF EXISTS `conversations`;
CREATE TABLE `conversations` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_model_config_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `model_provider` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `override_model_configs` text COLLATE utf8mb4_general_ci,
  `model_id` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `mode` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `summary` text COLLATE utf8mb4_general_ci,
  `inputs` text COLLATE utf8mb4_general_ci NOT NULL,
  `introduction` text COLLATE utf8mb4_general_ci,
  `system_instruction` text COLLATE utf8mb4_general_ci,
  `system_instruction_tokens` int NOT NULL DEFAULT '0',
  `status` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `from_source` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `from_end_user_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `from_account_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `read_at` timestamp NULL DEFAULT NULL,
  `read_account_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `invoke_from` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `dialogue_count` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of conversations
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for data_source_api_key_auth_bindings
-- ----------------------------
DROP TABLE IF EXISTS `data_source_api_key_auth_bindings`;
CREATE TABLE `data_source_api_key_auth_bindings` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `category` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `provider` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `credentials` text COLLATE utf8mb4_general_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `disabled` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of data_source_api_key_auth_bindings
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for data_source_oauth_bindings
-- ----------------------------
DROP TABLE IF EXISTS `data_source_oauth_bindings`;
CREATE TABLE `data_source_oauth_bindings` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `access_token` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `provider` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `source_info` json NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `disabled` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of data_source_oauth_bindings
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for dataset_auto_disable_logs
-- ----------------------------
DROP TABLE IF EXISTS `dataset_auto_disable_logs`;
CREATE TABLE `dataset_auto_disable_logs` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `document_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `notified` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of dataset_auto_disable_logs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for dataset_collection_bindings
-- ----------------------------
DROP TABLE IF EXISTS `dataset_collection_bindings`;
CREATE TABLE `dataset_collection_bindings` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider_name` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `model_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `collection_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `type` varchar(40) COLLATE utf8mb4_general_ci DEFAULT 'dataset',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of dataset_collection_bindings
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for dataset_keyword_tables
-- ----------------------------
DROP TABLE IF EXISTS `dataset_keyword_tables`;
CREATE TABLE `dataset_keyword_tables` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `keyword_table` text COLLATE utf8mb4_general_ci NOT NULL,
  `data_source_type` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'database',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of dataset_keyword_tables
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for dataset_permissions
-- ----------------------------
DROP TABLE IF EXISTS `dataset_permissions`;
CREATE TABLE `dataset_permissions` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `account_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `has_permission` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of dataset_permissions
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for dataset_process_rules
-- ----------------------------
DROP TABLE IF EXISTS `dataset_process_rules`;
CREATE TABLE `dataset_process_rules` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `mode` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'automatic',
  `rules` text COLLATE utf8mb4_general_ci,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of dataset_process_rules
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for dataset_queries
-- ----------------------------
DROP TABLE IF EXISTS `dataset_queries`;
CREATE TABLE `dataset_queries` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `content` text COLLATE utf8mb4_general_ci NOT NULL,
  `source` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `source_app_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_by_role` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of dataset_queries
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for dataset_retriever_resources
-- ----------------------------
DROP TABLE IF EXISTS `dataset_retriever_resources`;
CREATE TABLE `dataset_retriever_resources` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `message_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `position` int NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_name` text COLLATE utf8mb4_general_ci NOT NULL,
  `document_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `document_name` text COLLATE utf8mb4_general_ci NOT NULL,
  `data_source_type` text COLLATE utf8mb4_general_ci,
  `segment_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `score` double DEFAULT NULL,
  `content` text COLLATE utf8mb4_general_ci NOT NULL,
  `hit_count` int DEFAULT NULL,
  `word_count` int DEFAULT NULL,
  `segment_position` int DEFAULT NULL,
  `index_node_hash` text COLLATE utf8mb4_general_ci,
  `retriever_from` text COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of dataset_retriever_resources
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for datasets
-- ----------------------------
DROP TABLE IF EXISTS `datasets`;
CREATE TABLE `datasets` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `description` text COLLATE utf8mb4_general_ci,
  `provider` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'vendor',
  `permission` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'only_me',
  `data_source_type` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `indexing_technique` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `index_struct` text COLLATE utf8mb4_general_ci,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `embedding_model` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'text-embedding-ada-002',
  `embedding_model_provider` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'openai',
  `collection_binding_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `retrieval_model` json DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of datasets
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for dify_setups
-- ----------------------------
DROP TABLE IF EXISTS `dify_setups`;
CREATE TABLE `dify_setups` (
  `version` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `setup_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of dify_setups
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for document_segments
-- ----------------------------
DROP TABLE IF EXISTS `document_segments`;
CREATE TABLE `document_segments` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `document_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `position` int NOT NULL,
  `content` text COLLATE utf8mb4_general_ci NOT NULL,
  `word_count` int NOT NULL,
  `tokens` int NOT NULL,
  `keywords` json DEFAULT NULL,
  `index_node_id` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `index_node_hash` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `hit_count` int NOT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `disabled_at` timestamp NULL DEFAULT NULL,
  `disabled_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `status` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'waiting',
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `indexing_at` timestamp NULL DEFAULT NULL,
  `completed_at` timestamp NULL DEFAULT NULL,
  `error` text COLLATE utf8mb4_general_ci,
  `stopped_at` timestamp NULL DEFAULT NULL,
  `answer` text COLLATE utf8mb4_general_ci,
  `updated_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of document_segments
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for documents
-- ----------------------------
DROP TABLE IF EXISTS `documents`;
CREATE TABLE `documents` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `position` int NOT NULL,
  `data_source_type` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `data_source_info` text COLLATE utf8mb4_general_ci,
  `dataset_process_rule_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `batch` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_from` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_api_request_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `processing_started_at` timestamp NULL DEFAULT NULL,
  `file_id` text COLLATE utf8mb4_general_ci,
  `word_count` int DEFAULT NULL,
  `parsing_completed_at` timestamp NULL DEFAULT NULL,
  `cleaning_completed_at` timestamp NULL DEFAULT NULL,
  `splitting_completed_at` timestamp NULL DEFAULT NULL,
  `tokens` int DEFAULT NULL,
  `indexing_latency` double DEFAULT NULL,
  `completed_at` timestamp NULL DEFAULT NULL,
  `is_paused` tinyint(1) DEFAULT '0',
  `paused_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `paused_at` timestamp NULL DEFAULT NULL,
  `error` text COLLATE utf8mb4_general_ci,
  `stopped_at` timestamp NULL DEFAULT NULL,
  `indexing_status` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'waiting',
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `disabled_at` timestamp NULL DEFAULT NULL,
  `disabled_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `archived` tinyint(1) NOT NULL DEFAULT '0',
  `archived_reason` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `archived_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `archived_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `doc_type` varchar(40) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `doc_metadata` json DEFAULT NULL,
  `doc_form` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'text_model',
  `doc_language` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of documents
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for embeddings
-- ----------------------------
DROP TABLE IF EXISTS `embeddings`;
CREATE TABLE `embeddings` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `hash` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `embedding` blob NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `model_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'text-embedding-ada-002',
  `provider_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of embeddings
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for end_users
-- ----------------------------
DROP TABLE IF EXISTS `end_users`;
CREATE TABLE `end_users` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `type` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `external_user_id` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `is_anonymous` tinyint(1) NOT NULL DEFAULT '1',
  `session_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of end_users
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for external_knowledge_apis
-- ----------------------------
DROP TABLE IF EXISTS `external_knowledge_apis`;
CREATE TABLE `external_knowledge_apis` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `description` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `settings` text COLLATE utf8mb4_general_ci,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of external_knowledge_apis
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for external_knowledge_bindings
-- ----------------------------
DROP TABLE IF EXISTS `external_knowledge_bindings`;
CREATE TABLE `external_knowledge_bindings` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `external_knowledge_api_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `dataset_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `external_knowledge_id` text COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of external_knowledge_bindings
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for installed_apps
-- ----------------------------
DROP TABLE IF EXISTS `installed_apps`;
CREATE TABLE `installed_apps` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_owner_tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `position` int NOT NULL,
  `is_pinned` tinyint(1) NOT NULL DEFAULT '0',
  `last_used_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of installed_apps
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for invitation_codes
-- ----------------------------
DROP TABLE IF EXISTS `invitation_codes`;
CREATE TABLE `invitation_codes` (
  `id` int NOT NULL,
  `batch` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `code` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  `status` varchar(16) COLLATE utf8mb4_general_ci DEFAULT 'unused',
  `used_at` timestamp NULL DEFAULT NULL,
  `used_by_tenant_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `used_by_account_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `deprecated_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of invitation_codes
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for load_balancing_model_configs
-- ----------------------------
DROP TABLE IF EXISTS `load_balancing_model_configs`;
CREATE TABLE `load_balancing_model_configs` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `model_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `model_type` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `encrypted_config` text COLLATE utf8mb4_general_ci,
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of load_balancing_model_configs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for message_agent_thoughts
-- ----------------------------
DROP TABLE IF EXISTS `message_agent_thoughts`;
CREATE TABLE `message_agent_thoughts` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `message_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `message_chain_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `position` int NOT NULL,
  `thought` text COLLATE utf8mb4_general_ci,
  `tool` text COLLATE utf8mb4_general_ci,
  `tool_input` text COLLATE utf8mb4_general_ci,
  `observation` text COLLATE utf8mb4_general_ci,
  `tool_process_data` text COLLATE utf8mb4_general_ci,
  `message` text COLLATE utf8mb4_general_ci,
  `message_token` int DEFAULT NULL,
  `message_unit_price` decimal(10,0) DEFAULT NULL,
  `answer` text COLLATE utf8mb4_general_ci,
  `answer_token` int DEFAULT NULL,
  `answer_unit_price` decimal(10,0) DEFAULT NULL,
  `tokens` int DEFAULT NULL,
  `total_price` decimal(10,0) DEFAULT NULL,
  `currency` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `latency` double DEFAULT NULL,
  `created_by_role` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `message_price_unit` decimal(10,7) NOT NULL DEFAULT '0.0010000',
  `answer_price_unit` decimal(10,7) NOT NULL DEFAULT '0.0010000',
  `message_files` text COLLATE utf8mb4_general_ci,
  `tool_labels_str` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '{}',
  `tool_meta_str` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '{}',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of message_agent_thoughts
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for message_annotations
-- ----------------------------
DROP TABLE IF EXISTS `message_annotations`;
CREATE TABLE `message_annotations` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `conversation_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `message_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `content` text COLLATE utf8mb4_general_ci NOT NULL,
  `account_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `question` text COLLATE utf8mb4_general_ci,
  `hit_count` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of message_annotations
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for message_chains
-- ----------------------------
DROP TABLE IF EXISTS `message_chains`;
CREATE TABLE `message_chains` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `message_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `type` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `input` text COLLATE utf8mb4_general_ci,
  `output` text COLLATE utf8mb4_general_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of message_chains
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for message_feedbacks
-- ----------------------------
DROP TABLE IF EXISTS `message_feedbacks`;
CREATE TABLE `message_feedbacks` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `conversation_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `message_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `rating` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `content` text COLLATE utf8mb4_general_ci,
  `from_source` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `from_end_user_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `from_account_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of message_feedbacks
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for message_files
-- ----------------------------
DROP TABLE IF EXISTS `message_files`;
CREATE TABLE `message_files` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `message_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `type` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `transfer_method` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `url` text COLLATE utf8mb4_general_ci,
  `upload_file_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_by_role` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `belongs_to` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of message_files
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `model_provider` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `model_id` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `override_model_configs` text COLLATE utf8mb4_general_ci,
  `conversation_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `inputs` text COLLATE utf8mb4_general_ci NOT NULL,
  `query` text COLLATE utf8mb4_general_ci NOT NULL,
  `message` text COLLATE utf8mb4_general_ci NOT NULL,
  `message_tokens` int NOT NULL DEFAULT '0',
  `message_unit_price` decimal(10,4) NOT NULL,
  `answer` text COLLATE utf8mb4_general_ci NOT NULL,
  `answer_tokens` int NOT NULL DEFAULT '0',
  `answer_unit_price` decimal(10,4) NOT NULL,
  `provider_response_latency` double NOT NULL DEFAULT '0',
  `total_price` decimal(10,7) DEFAULT NULL,
  `currency` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `from_source` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `from_end_user_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `from_account_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `agent_based` tinyint(1) NOT NULL DEFAULT '0',
  `message_price_unit` decimal(10,7) NOT NULL DEFAULT '0.0010000',
  `answer_price_unit` decimal(10,7) NOT NULL DEFAULT '0.0010000',
  `workflow_run_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `status` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'normal',
  `error` text COLLATE utf8mb4_general_ci,
  `message_metadata` text COLLATE utf8mb4_general_ci,
  `invoke_from` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `parent_message_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of messages
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for operation_logs
-- ----------------------------
DROP TABLE IF EXISTS `operation_logs`;
CREATE TABLE `operation_logs` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `account_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `action` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `content` json DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_ip` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of operation_logs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for pinned_conversations
-- ----------------------------
DROP TABLE IF EXISTS `pinned_conversations`;
CREATE TABLE `pinned_conversations` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `conversation_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_role` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'end_user',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of pinned_conversations
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for provider_model_settings
-- ----------------------------
DROP TABLE IF EXISTS `provider_model_settings`;
CREATE TABLE `provider_model_settings` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `model_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `model_type` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `load_balancing_enabled` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of provider_model_settings
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for provider_models
-- ----------------------------
DROP TABLE IF EXISTS `provider_models`;
CREATE TABLE `provider_models` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `model_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `model_type` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `encrypted_config` text COLLATE utf8mb4_general_ci,
  `is_valid` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of provider_models
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for provider_orders
-- ----------------------------
DROP TABLE IF EXISTS `provider_orders`;
CREATE TABLE `provider_orders` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `account_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `payment_product_id` varchar(191) COLLATE utf8mb4_general_ci NOT NULL,
  `payment_id` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `transaction_id` varchar(191) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `quantity` int NOT NULL DEFAULT '1',
  `currency` varchar(40) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `total_amount` int DEFAULT NULL,
  `payment_status` varchar(40) COLLATE utf8mb4_general_ci DEFAULT 'wait_pay',
  `paid_at` timestamp NULL DEFAULT NULL,
  `pay_failed_at` timestamp NULL DEFAULT NULL,
  `refunded_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of provider_orders
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for providers
-- ----------------------------
DROP TABLE IF EXISTS `providers`;
CREATE TABLE `providers` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `provider_type` varchar(40) COLLATE utf8mb4_general_ci DEFAULT 'custom',
  `encrypted_config` text COLLATE utf8mb4_general_ci,
  `is_valid` tinyint(1) NOT NULL DEFAULT '0',
  `last_used` timestamp NULL DEFAULT NULL,
  `quota_type` varchar(40) COLLATE utf8mb4_general_ci DEFAULT '',
  `quota_limit` bigint DEFAULT NULL,
  `quota_used` bigint DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of providers
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for recommended_apps
-- ----------------------------
DROP TABLE IF EXISTS `recommended_apps`;
CREATE TABLE `recommended_apps` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `description` json NOT NULL,
  `copyright` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `privacy_policy` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `category` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `position` int NOT NULL,
  `is_listed` tinyint(1) NOT NULL,
  `install_count` int NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `language` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'en-US',
  `custom_disclaimer` text COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of recommended_apps
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for saved_messages
-- ----------------------------
DROP TABLE IF EXISTS `saved_messages`;
CREATE TABLE `saved_messages` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `message_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_role` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'end_user',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of saved_messages
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for sites
-- ----------------------------
DROP TABLE IF EXISTS `sites`;
CREATE TABLE `sites` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `title` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `icon` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `icon_background` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `description` text COLLATE utf8mb4_general_ci,
  `default_language` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `copyright` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `privacy_policy` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `customize_domain` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `customize_token_strategy` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `prompt_public` tinyint(1) NOT NULL DEFAULT '0',
  `status` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'normal',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `code` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `custom_disclaimer` text COLLATE utf8mb4_general_ci NOT NULL,
  `show_workflow_steps` tinyint(1) NOT NULL DEFAULT '1',
  `chat_color_theme` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `chat_color_theme_inverted` tinyint(1) NOT NULL DEFAULT '0',
  `icon_type` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `updated_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `use_icon_as_answer_icon` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of sites
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tag_bindings
-- ----------------------------
DROP TABLE IF EXISTS `tag_bindings`;
CREATE TABLE `tag_bindings` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `tag_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `target_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tag_bindings
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tags
-- ----------------------------
DROP TABLE IF EXISTS `tags`;
CREATE TABLE `tags` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `type` varchar(16) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tags
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tenant_account_joins
-- ----------------------------
DROP TABLE IF EXISTS `tenant_account_joins`;
CREATE TABLE `tenant_account_joins` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `account_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `role` varchar(16) COLLATE utf8mb4_general_ci DEFAULT 'normal',
  `invited_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `current` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tenant_account_joins
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tenant_default_models
-- ----------------------------
DROP TABLE IF EXISTS `tenant_default_models`;
CREATE TABLE `tenant_default_models` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `model_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `model_type` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tenant_default_models
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tenant_preferred_model_providers
-- ----------------------------
DROP TABLE IF EXISTS `tenant_preferred_model_providers`;
CREATE TABLE `tenant_preferred_model_providers` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `preferred_provider_type` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tenant_preferred_model_providers
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tenants
-- ----------------------------
DROP TABLE IF EXISTS `tenants`;
CREATE TABLE `tenants` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `encrypt_public_key` text COLLATE utf8mb4_general_ci,
  `plan` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'basic',
  `status` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'normal',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `custom_config` text COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tenants
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tidb_auth_bindings
-- ----------------------------
DROP TABLE IF EXISTS `tidb_auth_bindings`;
CREATE TABLE `tidb_auth_bindings` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `cluster_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `cluster_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '0',
  `status` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'CREATING',
  `account` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tidb_auth_bindings
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tool_api_providers
-- ----------------------------
DROP TABLE IF EXISTS `tool_api_providers`;
CREATE TABLE `tool_api_providers` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `schema` text COLLATE utf8mb4_general_ci NOT NULL,
  `schema_type_str` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `user_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tools_str` text COLLATE utf8mb4_general_ci NOT NULL,
  `icon` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `credentials_str` text COLLATE utf8mb4_general_ci NOT NULL,
  `description` text COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `privacy_policy` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `custom_disclaimer` text COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tool_api_providers
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tool_builtin_providers
-- ----------------------------
DROP TABLE IF EXISTS `tool_builtin_providers`;
CREATE TABLE `tool_builtin_providers` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `user_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `encrypted_credentials` text COLLATE utf8mb4_general_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tool_builtin_providers
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tool_conversation_variables
-- ----------------------------
DROP TABLE IF EXISTS `tool_conversation_variables`;
CREATE TABLE `tool_conversation_variables` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `user_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `conversation_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `variables_str` text COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tool_conversation_variables
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tool_files
-- ----------------------------
DROP TABLE IF EXISTS `tool_files`;
CREATE TABLE `tool_files` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `user_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `conversation_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `file_key` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `mimetype` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `original_url` text COLLATE utf8mb4_general_ci,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `size` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tool_files
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tool_label_bindings
-- ----------------------------
DROP TABLE IF EXISTS `tool_label_bindings`;
CREATE TABLE `tool_label_bindings` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tool_id` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `tool_type` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `label_name` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tool_label_bindings
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tool_model_invokes
-- ----------------------------
DROP TABLE IF EXISTS `tool_model_invokes`;
CREATE TABLE `tool_model_invokes` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `user_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `provider` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `tool_type` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `tool_name` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `model_parameters` text COLLATE utf8mb4_general_ci NOT NULL,
  `prompt_messages` text COLLATE utf8mb4_general_ci NOT NULL,
  `model_response` text COLLATE utf8mb4_general_ci NOT NULL,
  `prompt_tokens` int NOT NULL DEFAULT '0',
  `answer_tokens` int NOT NULL DEFAULT '0',
  `answer_unit_price` decimal(10,4) NOT NULL,
  `answer_price_unit` decimal(10,7) NOT NULL DEFAULT '0.0010000',
  `provider_response_latency` double NOT NULL DEFAULT '0',
  `total_price` decimal(10,7) DEFAULT NULL,
  `currency` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tool_model_invokes
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tool_published_apps
-- ----------------------------
DROP TABLE IF EXISTS `tool_published_apps`;
CREATE TABLE `tool_published_apps` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `user_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `description` text COLLATE utf8mb4_general_ci NOT NULL,
  `llm_description` text COLLATE utf8mb4_general_ci NOT NULL,
  `query_description` text COLLATE utf8mb4_general_ci NOT NULL,
  `query_name` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `tool_name` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `author` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tool_published_apps
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for tool_workflow_providers
-- ----------------------------
DROP TABLE IF EXISTS `tool_workflow_providers`;
CREATE TABLE `tool_workflow_providers` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(40) COLLATE utf8mb4_general_ci NOT NULL,
  `icon` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `user_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `description` text COLLATE utf8mb4_general_ci NOT NULL,
  `parameter_configuration` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '[]',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `privacy_policy` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  `version` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  `label` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of tool_workflow_providers
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for trace_app_config
-- ----------------------------
DROP TABLE IF EXISTS `trace_app_config`;
CREATE TABLE `trace_app_config` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tracing_provider` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `tracing_config` json DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of trace_app_config
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for upload_files
-- ----------------------------
DROP TABLE IF EXISTS `upload_files`;
CREATE TABLE `upload_files` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `storage_type` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `key` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `size` int NOT NULL,
  `extension` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `mime_type` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `used` tinyint(1) NOT NULL DEFAULT '0',
  `used_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `used_at` timestamp NULL DEFAULT NULL,
  `hash` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_by_role` varchar(255) COLLATE utf8mb4_general_ci DEFAULT 'account',
  `source_url` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of upload_files
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for whitelists
-- ----------------------------
DROP TABLE IF EXISTS `whitelists`;
CREATE TABLE `whitelists` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `category` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of whitelists
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for workflow_app_logs
-- ----------------------------
DROP TABLE IF EXISTS `workflow_app_logs`;
CREATE TABLE `workflow_app_logs` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `workflow_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `workflow_run_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_from` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by_role` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of workflow_app_logs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for workflow_conversation_variables
-- ----------------------------
DROP TABLE IF EXISTS `workflow_conversation_variables`;
CREATE TABLE `workflow_conversation_variables` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `conversation_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `data` text COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`,`conversation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of workflow_conversation_variables
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for workflow_node_executions
-- ----------------------------
DROP TABLE IF EXISTS `workflow_node_executions`;
CREATE TABLE `workflow_node_executions` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `workflow_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `triggered_from` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `workflow_run_id` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `index` int NOT NULL,
  `predecessor_node_id` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `node_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `node_type` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `title` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `inputs` text COLLATE utf8mb4_general_ci,
  `process_data` text COLLATE utf8mb4_general_ci,
  `outputs` text COLLATE utf8mb4_general_ci,
  `status` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `error` text COLLATE utf8mb4_general_ci,
  `elapsed_time` double NOT NULL DEFAULT '0',
  `execution_metadata` text COLLATE utf8mb4_general_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by_role` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `finished_at` timestamp NULL DEFAULT NULL,
  `node_execution_id` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of workflow_node_executions
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for workflow_runs
-- ----------------------------
DROP TABLE IF EXISTS `workflow_runs`;
CREATE TABLE `workflow_runs` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `sequence_number` int NOT NULL,
  `workflow_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `type` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `triggered_from` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `version` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `graph` text COLLATE utf8mb4_general_ci,
  `inputs` text COLLATE utf8mb4_general_ci,
  `status` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `outputs` text COLLATE utf8mb4_general_ci,
  `error` text COLLATE utf8mb4_general_ci,
  `elapsed_time` double NOT NULL DEFAULT '0',
  `total_tokens` bigint NOT NULL DEFAULT '0',
  `total_steps` int DEFAULT '0',
  `created_by_role` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `finished_at` timestamp NULL DEFAULT NULL,
  `exceptions_count` int DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of workflow_runs
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for workflows
-- ----------------------------
DROP TABLE IF EXISTS `workflows`;
CREATE TABLE `workflows` (
  `id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `tenant_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `app_id` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `type` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `version` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `graph` text COLLATE utf8mb4_general_ci NOT NULL,
  `features` text COLLATE utf8mb4_general_ci NOT NULL,
  `created_by` varchar(36) COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` varchar(36) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `updated_at` timestamp NOT NULL,
  `environment_variables` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `conversation_variables` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of workflows
-- ----------------------------
BEGIN;
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
