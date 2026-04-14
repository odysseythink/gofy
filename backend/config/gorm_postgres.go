package config

import "fmt"

type Postgres struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`                               // 服务器地址
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                               // 端口
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                         // 额外连接参数（如 sslmode=disable）
	Dbname       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`                      // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                   // 用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                   // 密码
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 全局表前缀
	Singular     bool   `mapstructure:"singular" json:"singular" yaml:"singular"`                   // 单数表名
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 最大打开连接数
	LogMode      string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`                   // Gorm 日志级别
	LogZap       bool   `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`                      // 是否写 zap
}

func (p *Postgres) Dsn() string {
	extra := p.Config
	if extra == "" {
		extra = "sslmode=disable TimeZone=Local"
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
		p.Path, p.Port, p.Username, p.Password, p.Dbname, extra)
}

func (p *Postgres) GetLogMode() string {
	return p.LogMode
}
