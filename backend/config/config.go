package config

type DsnProvider interface {
	Dsn() string
}

// Embeded 结构体可以压平到上一层，从而保持 config 文件的结构和原来一样
// 见 playground: https://go.dev/play/p/KIcuhqEoxmY

// GeneralDB 也被 Pgsql 和 Mysql 原样使用
type GeneralDB struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`                               // 服务器地址:端口
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                               //:端口
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                         // 高级配置
	Dbname       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`                      // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                   // 数据库用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                   // 数据库密码
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         //全局表前缀，单独定义TableName则不生效
	Singular     bool   `mapstructure:"singular" json:"singular" yaml:"singular"`                   //是否开启全局禁用复数，true表示开启
	Engine       string `mapstructure:"engine" json:"engine" yaml:"engine" default:"InnoDB"`        //数据库引擎，默认InnoDB
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`                   // 是否开启Gorm全局日志
	LogZap       bool   `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`                      // 是否通过zap写入日志文件
}

// type Config struct {
// 	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
// 	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
// 	Email   Email   `mapstructure:"email" json:"email" yaml:"email"`
// 	System  System  `mapstructure:"system" json:"system" yaml:"system"`
// 	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
// 	// auto
// 	AutoCode Autocode `mapstructure:"autocode" json:"autocode" yaml:"autocode"`
// 	// gorm
// 	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`

// 	// oss
// 	Local Local `mapstructure:"local" json:"local" yaml:"local"`

// 	Excel Excel `mapstructure:"excel" json:"excel" yaml:"excel"`
// 	Timer Timer `mapstructure:"timer" json:"timer" yaml:"timer"`

// 	// 跨域配置
// 	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`

// 	IntentTemplates map[string][]string `mapstructure:"intent_templates" json:"intent_templates" yaml:"intent_templates"`

// 	Edition                  string `mapstructure:"EDITION" json:"EDITION" yaml:"EDITION"`
// 	SecretKey                string `mapstructure:"SECRET_KEY" json:"SECRET_KEY" yaml:"SECRET_KEY"`
// 	EnterpriseEnabled        bool   `mapstructure:"ENTERPRISE_ENABLED" json:"ENTERPRISE_ENABLED" yaml:"ENTERPRISE_ENABLED"`
// 	CheckUpdateURL           string `mapstructure:"CHECK_UPDATE_URL" json:"CHECK_UPDATE_URL" yaml:"CHECK_UPDATE_URL"`
// 	CurrentVersion           string `mapstructure:"CURRENT_VERSION" json:"CURRENT_VERSION" yaml:"CURRENT_VERSION"`
// 	CanReplaceLogo           bool   `mapstructure:"CAN_REPLACE_LOGO" json:"CAN_REPLACE_LOGO" yaml:"CAN_REPLACE_LOGO"`
// 	ModelLBEnabled           bool   `mapstructure:"MODEL_LB_ENABLED" json:"MODEL_LB_ENABLED" yaml:"MODEL_LB_ENABLED"`
// 	VectorStore              string `mapstructure:"vector_store" json:"vector_store" yaml:"vector_store"`
// 	LoginLockoutDuration     int    `mapstructure:"login_lockout_duration" json:"login_lockout_duration" yaml:"login_lockout_duration"`
// 	EnableEmailCodeLogin     bool   `mapstructure:"enable_email_code_login" json:"enable_email_code_login" yaml:"enable_email_code_login"`
// 	EnableEmailPasswordLogin bool   `mapstructure:"enable_email_password_login" json:"enable_email_password_login" yaml:"enable_email_password_login"`
// 	EnableSocialOauthLogin   bool   `mapstructure:"enable_social_oauth_login" json:"enable_social_oauth_login" yaml:"enable_social_oauth_login"`
// 	AllowRegister            bool   `mapstructure:"allow_register" json:"allow_register" yaml:"allow_register"`
// 	AllowCreateWorkspace     bool   `mapstructure:"allow_create_workspace" json:"allow_create_workspace" yaml:"allow_create_workspace"`
// 	MailType                 string `mapstructure:"mail_type" json:"mail_type" yaml:"mail_type"`
// 	RefreshTokenExpireDays   int    `mapstructure:"refresh_token_expire_days" json:"refresh_token_expire_days" yaml:"refresh_token_expire_days"`
// 	DatasetOperatorEnabled   bool   `mapstructure:"dataset_operator_enabled" json:"dataset_operator_enabled" yaml:"dataset_operator_enabled"`
// 	BillingEnabled           bool   `mapstructure:"billing_enabled" json:"billing_enabled" yaml:"billing_enabled"`
// 	*Upload                  `mapstructure:"upload" json:"upload" yaml:"upload"`
// 	FilesURL                 string `mapstructure:"FILES_URL" json:"FILES_URL" yaml:"FILES_URL"`
// 	ServiceApiURL            string `mapstructure:"SERVICE_API_URL" json:"SERVICE_API_URL" yaml:"SERVICE_API_URL"`
// 	ConsoleApiURL            string `mapstructure:"CONSOLE_API_URL" json:"CONSOLE_API_URL" yaml:"CONSOLE_API_URL"`
// 	ConsoleWebURL            string `mapstructure:"CONSOLE_WEB_URL" json:"CONSOLE_WEB_URL" yaml:"CONSOLE_WEB_URL"`
// 	AppWebURL                string `mapstructure:"APP_WEB_URL" json:"APP_WEB_URL" yaml:"APP_WEB_URL"`

// 	FilesAccessTimeout        int    `mapstructure:"FILES_ACCESS_TIMEOUT" json:"FILES_ACCESS_TIMEOUT" yaml:"FILES_ACCESS_TIMEOUT"`
// 	MultimodalSendFormat      string `mapstructure:"MULTIMODAL_SEND_FORMAT" json:"MULTIMODAL_SEND_FORMAT" yaml:"MULTIMODAL_SEND_FORMAT"`
// 	PromptGenerationMaxTokens int    `mapstructure:"PROMPT_GENERATION_MAX_TOKENS" json:"PROMPT_GENERATION_MAX_TOKENS" yaml:"PROMPT_GENERATION_MAX_TOKENS"`
// 	CodeGenerationMaxTokens   int    `mapstructure:"CODE_GENERATION_MAX_TOKENS" json:"CODE_GENERATION_MAX_TOKENS" yaml:"CODE_GENERATION_MAX_TOKENS"`

// 	// Workflow runtime configuration
// 	Workflow struct {
// 		MaxExecutionSteps  int `mapstructure:"max_execution_steps" json:"max_execution_steps" yaml:"max_execution_steps"`
// 		MaxExecutionTime   int `mapstructure:"max_execution_time" json:"max_execution_time" yaml:"max_execution_time"`
// 		CallMaxDepth       int `mapstructure:"call_max_depth" json:"call_max_depth" yaml:"call_max_depth"`
// 		ParallelDepthLimit int `mapstructure:"parallel_depth_limit" json:"parallel_depth_limit" yaml:"parallel_depth_limit"`
// 		MaxVariableSize    int `mapstructure:"max_variable_size" json:"max_variable_size" yaml:"max_variable_size"`
// 	} `mapstructure:"workflow" json:"workflow" yaml:"workflow"`

// 	HttpNodeConfig struct {
// 		MaxConnectTimeout int `mapstructure:"max_connect_timeout" json:"max_connect_timeout" yaml:"max_connect_timeout"`
// 		MaxReadTimeout    int `mapstructure:"max_read_timeout" json:"max_read_timeout" yaml:"max_read_timeout"`
// 		MaxWriteTimeout   int `mapstructure:"max_write_timeout" json:"max_write_timeout" yaml:"max_write_timeout"`
// 		MaxBinarySize     int `mapstructure:"max_binary_size" json:"max_binary_size" yaml:"max_binary_size"`
// 		MaxTextSize       int `mapstructure:"max_text_size" json:"max_text_size" yaml:"max_text_size"`
// 	} `mapstructure:"http_node_config" json:"http_node_config" yaml:"http_node_config"`

// 	AppConfig struct {
// 		MaxExecutionTime  int `mapstructure:"max_execution_time" json:"max_execution_time" yaml:"max_execution_time"`
// 		MaxActiveRequests int `mapstructure:"max_active_requests" json:"max_active_requests" yaml:"max_active_requests"`
// 	} `mapstructure:"app_config" json:"app_config" yaml:"app_config"`

// 	Debug bool `mapstructure:"debug" json:"debug" yaml:"debug"`

// 	Moderation struct {
// 		BufferSize int `mapstructure:"buffer_size" json:"buffer_size" yaml:"buffer_size"` //"Size of the buffer for content moderation processing"
// 	} `mapstructure:"moderation" json:"moderation" yaml:"moderation"`
// }

// type Upload struct {
// 	FileSizeLimit      int `mapstructure:"file_size_limit" json:"file_size_limit" yaml:"file_size_limit"`
// 	FileBatchLimit     int `mapstructure:"file_batch_limit" json:"file_batch_limit" yaml:"file_batch_limit"`
// 	ImageFileSizeLimit int `mapstructure:"image_file_size_limit" json:"image_file_size_limit" yaml:"image_file_size_limit"`
// 	VideoFileSizeLimit int `mapstructure:"video_file_size_limit" json:"video_file_size_limit" yaml:"video_file_size_limit"`
// 	AudioFileSizeLimit int `mapstructure:"audio_file_size_limit" json:"audio_file_size_limit" yaml:"audio_file_size_limit"`
// 	WorkflowFileLimit  int `mapstructure:"workflow_file_limit" json:"workflow_file_limit" yaml:"workflow_file_limit"`
// }
