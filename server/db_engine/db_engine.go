package dbengine

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mlib.com/gofy/server/config"
	"mlib.com/mlog"
)

type DBEngine struct {
	*gorm.DB
}

func (m *DBEngine) Init() error {
	tmp := viper.Get("mysql")
	bindata, err := json.Marshal(tmp)
	if err != nil {
		mlog.Warningf("json marshal=%#v to string failed:%v", tmp, err)
		return err
	}
	cfg := config.Mysql{}
	err = json.Unmarshal(bindata, &cfg)
	if err != nil {
		mlog.Warningf("json Unmarshal=%#v to Mysql config failed:%v", string(bindata), err)
		return err
	}

	if cfg.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       cfg.Dsn(), // DSN data source name
		DefaultStringSize:         191,       // string 类型字段的默认长度
		SkipInitializeWithVersion: false,     // 根据版本自动配置

	}

	if m.DB, err = gorm.Open(mysql.New(mysqlConfig), NewGormConfig(cfg.GetLogMode(), cfg.Prefix, cfg.Singular)); err != nil {
		mlog.Warningf("gorm.Open(%s) failed:%v", cfg.Dsn(), err)
		return fmt.Errorf("gorm.Open(%s) failed:%v", cfg.Dsn(), err)
	}

	m.DB.InstanceSet("gorm:table_options", "ENGINE="+cfg.Engine)
	sqlDB, err := m.DB.DB()
	if err != nil {
		mlog.Warningf("can't get sql.DB failed:%v", err)
		return fmt.Errorf("can't get sql.DB failed:%v", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		mlog.Warningf("sql.DB ping failed:%v", err)
		return fmt.Errorf("sql.DB ping failed:%v", err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

func (m *DBEngine) RunOnce(ctx context.Context) {

}

func (m *DBEngine) Destroy() {

}

var (
	gOnce     sync.Once
	gInstance *DBEngine
)

func Instance() *DBEngine {
	gOnce.Do(func() {
		gInstance = &DBEngine{}
	})
	return gInstance
}
