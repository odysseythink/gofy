package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"mlib.com/confy"
	"mlib.com/mlog"
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
	if !confy.InConfig("redis") {
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
	// 	confy.Get[bool]("redis.is_cluster"),
	// 	)
	if confy.Get[bool]("redis.is_cluster") {
		addrs := strings.ReplaceAll(confy.Get[string]("redis.addr"), " ", "")
		rdb := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(addrs, ","),
			Password: confy.Get[string]("redis.password"),
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
			Addr:     confy.Get[string]("redis.addr"),
			Password: confy.Get[string]("redis.password"),
			DB:       confy.Get[int]("redis.db"),
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
