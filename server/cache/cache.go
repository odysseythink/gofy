package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"mlib.com/mlog"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Cache struct {
	// rdsClusterCli *redis.ClusterClient
	// rdsCli        *redis.Client
	rdsCli redis.UniversalClient
}

func (m *Cache) validatorConfig() bool {
	return true
}

func (m *Cache) Init(args ...any) error {
	if viper.Get("redis") == nil {
		mlog.Errorf("don't exist redis config, return")
		return errors.New("don't exist redis config, return")
	}

	if !m.validatorConfig() {
		mlog.Errorf("redis config is invalid")
		return fmt.Errorf("redis config is invalid")
	}
	//  = NewRedisClient(
	//
	// 	,
	// 	viper.GetBool("redis.is_cluster"),
	// 	)
	if viper.GetBool("redis.is_cluster") {
		addrs := strings.ReplaceAll(viper.GetString("redis.addr"), " ", "")
		rdb := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(addrs, ","),
			Password: viper.GetString("redis.password"),
		})
		err := rdb.ForEachShard(context.Background(), func(ctx context.Context, shared *redis.Client) error {
			return shared.Ping(ctx).Err()
		})
		if err != nil {
			mlog.Errorf("redis ring cluster ping failed:%v", err)
			return fmt.Errorf("redis ring cluster ping failed:%v", err)
		}
		// m.rdsClusterCli = rdb
		m.rdsCli = rdb
	} else {
		rdb := redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		})
		err := rdb.Ping(context.Background()).Err()
		if err != nil {
			mlog.Errorf("redis ping failed:%v", err)
			return fmt.Errorf("redis ping failed:%v", err)
		}
		// m.rdsCli = rdb
		m.rdsCli = rdb
	}

	if m.rdsCli == nil {
		mlog.Errorf("new redis client failed")
		return fmt.Errorf("new redis client failed")
	}
	return nil

}

func (m *Cache) RunOnce(ctx context.Context) error {
	return nil
}

func (m *Cache) Destroy() {
}
func (m *Cache) UserData() any {
	return m
}

var (
	gOnce     sync.Once
	gInstance *Cache
)

func Instance() *Cache {
	gOnce.Do(func() {
		gInstance = &Cache{}
	})
	return gInstance
}
