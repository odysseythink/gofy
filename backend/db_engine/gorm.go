package dbengine

import (
	// "fmt"
	// "log"
	// "os"
	// "time"

	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	// "github.com/odysseythink/mlog"
)

// Gorm 初始化数据库并产生数据库全局变量
// Author SliverHorn
// func Gorm() *gorm.DB {
// 	return GormMysql()
// }

// Config gorm 自定义配置
func NewGormConfig(logmode, prefix string, singular bool) *gorm.Config {
	// lm_log := NewGormLogger()
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		// Set logger
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	// _default := logger.New(NewGormWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
	// 	SlowThreshold: 200 * time.Millisecond,
	// 	LogLevel:      logger.Warn,
	// 	Colorful:      true,
	// })

	// switch logmode {
	// case "silent", "Silent":
	// 	config.Logger = _default.LogMode(logger.Silent)
	// case "error", "Error":
	// 	config.Logger = _default.LogMode(logger.Error)
	// case "warn", "Warn":
	// 	config.Logger = _default.LogMode(logger.Warn)
	// case "info", "Info":
	// 	config.Logger = _default.LogMode(logger.Info)
	// default:
	// 	config.Logger = _default.LogMode(logger.Info)
	// }
	return config
}

// type gormwriter struct {
// 	logger.Writer
// }

// // NewWriter writer 构造函数
// func NewGormWriter(w logger.Writer) *gormwriter {
// 	return &gormwriter{Writer: w}
// }

// // Printf 格式化打印日志
// func (w *gormwriter) Printf(message string, data ...any) {

// 	mlog.Info(fmt.Sprintf(message+"\n", data...))
// }
