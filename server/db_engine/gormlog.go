package dbengine

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"mlib.com/mlog"
)

type GormLogger struct {
	Log_all       bool
	SlowThreshold time.Duration
}

func NewGormLogger() GormLogger {
	sql_log := viper.GetBool("system.sql_log")
	slow_log := viper.GetInt("system.slow_log")
	if slow_log == 0 {
		slow_log = 1000
	}
	return GormLogger{
		Log_all:       sql_log,
		SlowThreshold: time.Duration(slow_log) * time.Millisecond,
	}
}

// LogMode 实现 gormlogger.Interface 的 LogMode 方法
func (l GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return GormLogger{
		SlowThreshold: l.SlowThreshold,
	}
}

// Info 实现 gormlogger.Interface 的 Info 方法
func (l GormLogger) Info(ctx context.Context, str string, args ...any) {
	mlog.Debugf(str, args...)
}

// Warn 实现 gormlogger.Interface 的 Warn 方法
func (l GormLogger) Warn(ctx context.Context, str string, args ...any) {
	mlog.Warningf(str, args...)
}

// Error 实现 gormlogger.Interface 的 Error 方法
func (l GormLogger) Error(ctx context.Context, str string, args ...any) {
	mlog.Errorf(str, args...)
}

// Trace 实现 gormlogger.Interface 的 Trace 方法
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	} else {
		file = "???"
		line = 0
	}
	// 获取运行时间
	elapsed := time.Since(begin)
	// 获取 SQL 请求和返回条数
	sql, _ := fc()

	// Gorm 错误
	if err != nil {
		// 记录未找到的错误使用 warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			mlog.Warningf("[ %s:%d ] DB Warn: %s  %s ", file, line, elapsed.String(), sql)
		} else {
			mlog.Errorf("[ %s:%d ] DB Error: %s  %s %+v", file, line, elapsed.String(), sql, err)
		}
		// 如果有错误，就 return 掉，防止后面重复输出
		return
	}

	// 慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		mlog.Warningf("[ %s:%d ] DB Slow: %s  %s ", file, line, elapsed.String(), sql)
		// 如果是慢查询，就 return 掉，防止后面重复输出
		return
	}
	if l.Log_all {
		// 记录所有 SQL 请求
		mlog.Debugf("[ %s:%d ] DB Query: %s  %s ", file, line, elapsed.String(), sql)
	}
}
