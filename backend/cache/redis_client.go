package cache

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/odysseythink/mlog"
	"github.com/redis/go-redis/v9"
)

func (c *Cache) Unlink(keys ...string) error {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return fmt.Errorf("redis don't have client")
		// return
	}
	if len(keys) == 0 {
		mlog.Warningf("invalid arg")
		return fmt.Errorf("invalid arg")
		// return
	}
	// return c.rdb.Del(context.Background(), key).Err()
	// if c.rdsCli != nil {
	cmd := c.rdsCli.Unlink(context.Background(), keys...)
	if cmd.Err() != nil {
		mlog.Errorf("redis Unlink(%v) failed:%v", keys, cmd.Err())
		return fmt.Errorf("redis Unlink(%v) failed:%v", keys, cmd.Err())
		// return
	}
	if int(cmd.Val()) != len(keys) {
		mlog.Errorf("not all key(%v) Unlink", keys)
		return fmt.Errorf("not all key(%v) Unlink", keys)
	}
	// } else {
	// 	var cnt int32 = 0
	// 	for _, v := range keys {
	// 		cmd := c.rdsClusterCli.Unlink(context.Background(), v)
	// 		if cmd.Err() != nil {
	// 			mlog.Errorf("redis cluster Unlink(%v) failed:%v", v, cmd.Err())
	// 			return fmt.Errorf("redis cluster Unlink(%v) failed:%v", v, cmd.Err())
	// 			// return
	// 		}
	// 		atomic.AddInt32(&cnt, int32(cmd.Val()))
	// 	}
	// 	if int(cnt) != len(keys) {
	// 		mlog.Errorf("not all key(%v) Unlink", keys)
	// 		return fmt.Errorf("not all key(%v) Unlink", keys)
	// 	}
	// }
	return nil
}

func (c *Cache) DelKey(key string) error {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return fmt.Errorf("redis don't have client")
		// return
	}
	if key == "" {
		mlog.Errorf("invalid arg")
		return fmt.Errorf("invalid arg")
		// return
	}
	// return c.rdb.Del(context.Background(), key).Err()
	// if c.rdsCli != nil {
	cmd := c.rdsCli.Del(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Del(%s) failed:%v", key, cmd.Err())
		return fmt.Errorf("redis Del(%s) failed:%v", key, cmd.Err())
		// return
	}
	if int(cmd.Val()) == 0 {
		mlog.Errorf("no redis key(%v) exist", key)
		return fmt.Errorf("no redis key(%v) exist", key)
	}
	// } else {
	// 	cmd := c.rdsClusterCli.Del(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis cluster Del(%s) failed:%v", key, cmd.Err())
	// 		return fmt.Errorf("redis cluster Del(%s) failed:%v", key, cmd.Err())
	// 		// return
	// 	}
	// 	if int(cmd.Val()) == 0 {
	// 		mlog.Errorf("no redis key(%v) exist", key)
	// 		return fmt.Errorf("no redis key(%v) exist", key)
	// 	}
	// }
	return nil
}

func (c *Cache) ExistsKey(key string) bool {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return false, fmt.Errorf("redis don't have cluster client")
		return false
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return false, fmt.Errorf("invalid arg")
		return false
	}
	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Exists(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Exists(%s) failed:%v", key, cmd.Err())
		// return false, fmt.Errorf("redis Exists(%s) failed:%v", key, cmd.Err())
		return false
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Exists(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Exists(%s) failed:%v", key, cmd.Err())
	// 		// return false, fmt.Errorf("redis Exists(%s) failed:%v", key, cmd.Err())
	// 		return false
	// 	}
	// }
	return cmd.Val() == 1
}

func (c *Cache) TypeOfKey(key string) string {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return false, fmt.Errorf("redis don't have cluster client")
		return ""
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return false, fmt.Errorf("invalid arg")
		return ""
	}
	var cmd *redis.StatusCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Type(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Type(%s) failed:%v", key, cmd.Err())
		// return false, fmt.Errorf("redis Exists(%s) failed:%v", key, cmd.Err())
		return ""
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Type(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Type(%s) failed:%v", key, cmd.Err())
	// 		// return false, fmt.Errorf("redis Exists(%s) failed:%v", key, cmd.Err())
	// 		return ""
	// 	}
	// }
	return cmd.Val()
}

func (c *Cache) IncrKey(key string) int64 {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return 0
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return 0
	}
	// return c.rdb.Incr(context.Background(), key).Err()
	// if c.rdsCli != nil {
	cmd := c.rdsCli.Incr(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Incr(%s) failed:%v", key, cmd.Err())
		return 0
	}
	return cmd.Val()
	// } else {
	// 	cmd := c.rdsClusterCli.Incr(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Incr(%s) failed:%v", key, cmd.Err())
	// 		return 0
	// 	}
	// 	return cmd.Val()
	// }
}

func (c *Cache) ExpireKey(key string, iSeconds int) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	// return c.rdb.Expire(context.Background(), key, time.Duration(iSeconds)*time.Second).Err()
	// if c.rdsCli != nil {
	if iSeconds < 0 {
		cmd := c.rdsCli.Persist(context.Background(), key)
		if cmd.Err() != nil {
			mlog.Errorf("redis Expire(%s) failed:%v", key, cmd.Err())
			return
		}
	} else {
		cmd := c.rdsCli.Expire(context.Background(), key, time.Duration(iSeconds)*time.Second)
		if cmd.Err() != nil {
			mlog.Errorf("redis Expire(%s) failed:%v", key, cmd.Err())
			return
		}
	}
	// } else {
	// 	if iSeconds < 0 {
	// 		cmd := c.rdsClusterCli.Persist(context.Background(), key)
	// 		if cmd.Err() != nil {
	// 			mlog.Errorf("redis Expire(%s) failed:%v", key, cmd.Err())
	// 			return
	// 		}
	// 	} else {
	// 		cmd := c.rdsClusterCli.Expire(context.Background(), key, time.Duration(iSeconds)*time.Second)
	// 		if cmd.Err() != nil {
	// 			mlog.Errorf("redis Expire(%s) failed:%v", key, cmd.Err())
	// 			return
	// 		}
	// 	}
	// }
}

func (c *Cache) TTL(key string) int {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return 0
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return 0
	}
	// return c.rdb.Expire(context.Background(), key, time.Duration(iSeconds)*time.Second).Err()
	// if c.rdsCli != nil {
	cmd := c.rdsCli.TTL(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis TTL(%s) failed:%v", key, cmd.Err())
		return 0
	}
	return int(cmd.Val().Seconds())
	// } else {
	// 	cmd := c.rdsClusterCli.TTL(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis TTL(%s) failed:%v", key, cmd.Err())
	// 		return 0
	// 	}
	// 	return int(cmd.Val().Seconds())
	// }
}

func (c *Cache) SetExKey(key string, val interface{}, expiration time.Duration) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	// return c.rdb.Set(context.Background(), key, val, time.Duration(iSeconds)*time.Second).Err()
	// if c.rdsCli != nil {
	cmd := c.rdsCli.Set(context.Background(), key, val, expiration)
	if cmd.Err() != nil {
		mlog.Errorf("redis Set(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.Set(context.Background(), key, val, expiration)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Set(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) PTTL(key string) int64 {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return 0
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return 0
	}
	// return c.rdb.Set(context.Background(), key, val, time.Duration(iSeconds)*time.Second).Err()
	var cmd *redis.DurationCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.PTTL(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Set(%s) failed:%v", key, cmd.Err())
		return 0
	}
	// } else {
	// 	cmd = c.rdsClusterCli.PTTL(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Set(%s) failed:%v", key, cmd.Err())
	// 		return 0
	// 	}
	// }
	return int64(cmd.Val().Seconds())
}

func (c *Cache) Get(key string, val interface{}) error {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return fmt.Errorf("redis don't have cluster client")
	}
	if key == "" || reflect.ValueOf(val).Kind() != reflect.Ptr {
		mlog.Errorf("invalid arg")
		// return fmt.Errorf("invalid arg")
		return fmt.Errorf("invalid arg")
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Get(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
		return fmt.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Get(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	// 		return fmt.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	// 	}
	// }
	if cmd.Scan(val) != nil {
		mlog.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		return fmt.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
	}
	return nil
}

func (c *Cache) GetString(key string) string {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return ""
	}
	if key == "" {
		mlog.Errorf("invalid arg")
		return ""
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Get(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
		return ""
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Get(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	// 		return ""
	// 	}
	// }
	return cmd.Val()
}

func (c *Cache) GetBool(key string) bool {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return false
	}
	if key == "" {
		mlog.Errorf("invalid arg")
		return false
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Get(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
		return false
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Get(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	// 		return false
	// 	}
	// }
	res, err := cmd.Bool()
	if err != nil {
		mlog.Errorf("redis val to bool failed:%v", err)
		return false
	}
	return res
}

func (c *Cache) GetInt(key string) int {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return 0
	}
	if key == "" {
		mlog.Errorf("invalid arg")
		return 0
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Get(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
		return 0
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Get(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	// 		return 0
	// 	}
	// }
	res, err := cmd.Int()
	if err != nil {
		mlog.Errorf("redis val to int failed:%v", err)
		return 0
	}
	return res
}

func (c *Cache) GetInt64(key string) int64 {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return 0
	}
	if key == "" {
		mlog.Errorf("invalid arg")
		return 0
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Get(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
		return 0
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Get(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	// 		return 0
	// 	}
	// }
	res, err := cmd.Int64()
	if err != nil {
		mlog.Errorf("redis val to int failed:%v", err)
		return 0
	}
	return res
}

func (c *Cache) GetUint64(key string) uint64 {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return 0
	}
	if key == "" {
		mlog.Errorf("invalid arg")
		return 0
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Get(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
		return 0
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Get(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	// 		return 0
	// 	}
	// }
	res, err := cmd.Uint64()
	if err != nil {
		mlog.Errorf("redis val to int failed:%v", err)
		return 0
	}
	return res
}

func (c *Cache) GetFloat32(key string) float32 {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return 0
	}
	if key == "" {
		mlog.Errorf("invalid arg")
		return 0
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Get(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
		return 0
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Get(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	// 		return 0
	// 	}
	// }
	res, err := cmd.Float32()
	if err != nil {
		mlog.Errorf("redis val to int failed:%v", err)
		return 0
	}
	return res
}

func (c *Cache) GetFloat64(key string) float64 {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return 0
	}
	if key == "" {
		mlog.Errorf("invalid arg")
		return 0
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Get(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
		return 0
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Get(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Get(%s) failed:%v", key, cmd.Err())
	// 		return 0
	// 	}
	// }
	res, err := cmd.Float64()
	if err != nil {
		mlog.Errorf("redis val to int failed:%v", err)
		return 0
	}
	return res
}

func (c *Cache) Set(key string, val interface{}) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}

	// return c.rdb.Set(context.Background(), key, val, -1*time.Second).Err()
	// if c.rdsCli != nil {
	cmd := c.rdsCli.Set(context.Background(), key, val, -1*time.Second)
	if cmd.Err() != nil {
		mlog.Errorf("redis Set(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.Set(context.Background(), key, val, -1*time.Second)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Set(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) SetEx(key string, val interface{}, tm time.Duration) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		return fmt.Errorf("redis don't have cluster client")
		// return
	}
	if key == "" {
		mlog.Error("invalid arg")
		return errors.New("invalid arg")
		// return
	}
	// return c.rdb.SetEx(context.Background(), key, val, -1*time.Second).Err()
	// if c.rdsCli != nil {
	cmd := c.rdsCli.SetEx(context.Background(), key, val, tm)
	if cmd.Err() != nil {
		mlog.Errorf("redis SetEx(%s) failed:%v", key, cmd.Err())
		return fmt.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		// return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.SetEx(context.Background(), key, val, tm)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis SetEx(%s) failed:%v", key, cmd.Err())
	// 		return fmt.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
	// 		// return
	// 	}
	// }
	return nil
}

func (c *Cache) SetNX(key string, val interface{}, tm time.Duration) (error, bool) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		return fmt.Errorf("redis don't have cluster client"), false
		// return
	}
	if key == "" {
		mlog.Error("invalid arg")
		return errors.New("invalid arg"), false
		// return
	}
	cmd := c.rdsCli.SetNX(context.Background(), key, val, tm)
	if cmd.Err() != nil {
		mlog.Errorf("redis SetNX(%s) failed:%v", key, cmd.Err())
		return fmt.Errorf("redis SetNX(%s) failed:%v", key, cmd.Err()), false
		// return
	}

	return nil, cmd.Val()
}

func (c *Cache) LIndex(key string, index int64) string {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return 0, fmt.Errorf("redis don't have cluster client")
		return ""
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return 0, fmt.Errorf("invalid arg")
		return ""
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.LIndex(context.Background(), key, index)
	if cmd.Err() != nil {
		mlog.Errorf("redis LIndex(%s, %d) failed:%v", key, index, cmd.Err())
		// return 0, fmt.Errorf("redis LLen(%s) failed:%v", key, cmd.Err())
		return ""
	}
	// } else {
	// 	cmd = c.rdsClusterCli.LIndex(context.Background(), key, index)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis LIndex(%s, %d) failed:%v", key, index, cmd.Err())
	// 		// return 0, fmt.Errorf("redis LLen(%s) failed:%v", key, cmd.Err())
	// 		return ""
	// 	}
	// }
	return cmd.Val()
}

func (c *Cache) LInsert(key string, op string, pivot interface{}, value interface{}) int {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return 0, fmt.Errorf("redis don't have cluster client")
		return 0
	}
	if key == "" || (op != "BEFORE" && op != "AFTER") || pivot == nil || value == nil {
		mlog.Error("invalid arg")
		// return 0, fmt.Errorf("invalid arg")
		return 0
	}
	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.LInsert(context.Background(), key, op, pivot, value)
	if cmd.Err() != nil {
		mlog.Errorf("redis LInsert(%s, %s, %#v, %#v) failed:%v", key, op, pivot, value, cmd.Err())
		return 0
	}
	// } else {
	// 	cmd = c.rdsClusterCli.LInsert(context.Background(), key, op, pivot, value)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis LInsert(%s, %s, %#v, %#v) failed:%v", key, op, pivot, value, cmd.Err())
	// 		return 0
	// 	}
	// }
	return int(cmd.Val())
}

func (c *Cache) LInsertBefore(key string, pivot interface{}, value interface{}) int {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return 0, fmt.Errorf("redis don't have cluster client")
		return 0
	}
	if key == "" || pivot == nil || value == nil {
		mlog.Error("invalid arg")
		// return 0, fmt.Errorf("invalid arg")
		return 0
	}
	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.LInsertBefore(context.Background(), key, pivot, value)
	if cmd.Err() != nil {
		mlog.Errorf("redis LInsertBefore(%s, %#v, %#v) failed:%v", key, pivot, value, cmd.Err())
		return 0
	}
	// } else {
	// 	cmd = c.rdsClusterCli.LInsertBefore(context.Background(), key, pivot, value)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis LInsertBefore(%s, %#v, %#v) failed:%v", key, pivot, value, cmd.Err())
	// 		return 0
	// 	}
	// }
	return int(cmd.Val())
}

func (c *Cache) LInsertAfter(key string, pivot interface{}, value interface{}) int {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return 0, fmt.Errorf("redis don't have cluster client")
		return 0
	}
	if key == "" || pivot == nil || value == nil {
		mlog.Error("invalid arg")
		// return 0, fmt.Errorf("invalid arg")
		return 0
	}
	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.LInsertAfter(context.Background(), key, pivot, value)
	if cmd.Err() != nil {
		mlog.Errorf("redis LInsertAfter(%s, %#v, %#v) failed:%v", key, pivot, value, cmd.Err())
		return 0
	}
	// } else {
	// 	cmd = c.rdsClusterCli.LInsertAfter(context.Background(), key, pivot, value)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis LInsertAfter(%s, %#v, %#v) failed:%v", key, pivot, value, cmd.Err())
	// 		return 0
	// 	}
	// }
	return int(cmd.Val())
}

func (c *Cache) LLen(key string) int64 {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return 0, fmt.Errorf("redis don't have cluster client")
		return 0
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return 0, fmt.Errorf("invalid arg")
		return 0
	}
	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.LLen(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis LLen(%s) failed:%v", key, cmd.Err())
		// return 0, fmt.Errorf("redis LLen(%s) failed:%v", key, cmd.Err())
		return 0
	}
	// } else {
	// 	cmd = c.rdsClusterCli.LLen(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis LLen(%s) failed:%v", key, cmd.Err())
	// 		// return 0, fmt.Errorf("redis LLen(%s) failed:%v", key, cmd.Err())
	// 		return 0
	// 	}
	// }
	return cmd.Val()
}

func (c *Cache) LPop(key string, val interface{}) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" || reflect.ValueOf(val).Kind() != reflect.Ptr {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.LPop(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis LPop(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd = c.rdsClusterCli.LPop(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis LPop(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
	err := cmd.Scan(val)
	if err != nil {
		mlog.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		return
	}
}

func (c *Cache) RPop(key string, val interface{}) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" || reflect.ValueOf(val).Kind() != reflect.Ptr {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.RPop(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis RPop(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd = c.rdsClusterCli.RPop(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis RPop(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
	err := cmd.Scan(val)
	if err != nil {
		mlog.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		return
	}
}

func (c *Cache) LPush(key string, vals ...interface{}) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		return
	}
	if key == "" || len(vals) < 1 {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.LPush(context.Background(), key, vals...)
	if cmd.Err() != nil {
		mlog.Errorf("redis LPush(%s, %#v) failed:%v", key, vals, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.LPush(context.Background(), key, vals...)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis LPush(%s, %#v) failed:%v", key, vals, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) RPush(key string, vals ...interface{}) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		return
	}
	if key == "" || len(vals) < 1 {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.RPush(context.Background(), key, vals...)
	if cmd.Err() != nil {
		mlog.Errorf("redis RPush(%s, %#v) failed:%v", key, vals, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.RPush(context.Background(), key, vals...)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis RPush(%s, %#v) failed:%v", key, vals, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) LRange(key string, start, stop int64) []string {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		return nil
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return nil
	}
	var cmd *redis.StringSliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.LRange(context.Background(), key, start, stop)
	if cmd.Err() != nil {
		mlog.Errorf("redis LRange(%s, %v, %v) failed:%v", key, start, stop, cmd.Err())
		return nil
	}
	// } else {
	// 	cmd = c.rdsClusterCli.LRange(context.Background(), key, start, stop)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis LRange(%s, %v, %v) failed:%v", key, start, stop, cmd.Err())
	// 		return nil
	// 	}
	// }
	return cmd.Val()
}

func (c *Cache) LRem(key string, val interface{}) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" || val == nil {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.LRem(context.Background(), key, 0, val)
	if cmd.Err() != nil {
		mlog.Errorf("redis LRem(%s, %v) failed:%v", key, val, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.LRem(context.Background(), key, 0, val)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis LRem(%s, %v) failed:%v", key, val, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) LSet(key string, idx int64, val interface{}) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.LSet(context.Background(), key, idx, val)
	if cmd.Err() != nil {
		mlog.Errorf("redis LRem(%s,%v, %v) failed:%v", key, idx, val, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.LSet(context.Background(), key, idx, val)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis LRem(%s,%v, %v) failed:%v", key, idx, val, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) HMGetData(key string, attrNames []string) ([]interface{}, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return nil, errors.New("redis don't have cluster client")
	}
	if key == "" || attrNames == nil {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return nil, errors.New("redis HMGetData: invalid arg")
	}
	var cmd *redis.SliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.HMGet(context.Background(), key, attrNames...)
	if cmd.Err() != nil {
		mlog.Errorf("redis HMGet(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis HMGet(%s) failed:%v", key, cmd.Err())
		return nil, fmt.Errorf("redis HMGet(%s) failed:%v", key, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.HMGet(context.Background(), key, attrNames...)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HMGet(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis HMGet(%s) failed:%v", key, cmd.Err())
	// 		return nil, fmt.Errorf("redis HMGet(%s) failed:%v", key, cmd.Err())
	// 	}
	// }
	return cmd.Result()
}

func (c *Cache) HGetAllData(key string) map[string]string {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return nil
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return nil
	}
	var cmd *redis.MapStringStringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.HGetAll(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis HMGet(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis HMGet(%s) failed:%v", key, cmd.Err())
		return nil
	}
	// } else {
	// 	cmd = c.rdsClusterCli.HGetAll(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HMGet(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis HMGet(%s) failed:%v", key, cmd.Err())
	// 		return nil
	// 	}
	// }
	attr, err := cmd.Result()
	// err := cmd.Scan(attr)
	if err != nil {
		mlog.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		return nil
	}
	return attr
}

func (c *Cache) HGetData(key string, field string, val interface{}) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" || field == "" || reflect.ValueOf(val).Kind() != reflect.Ptr {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.HGet(context.Background(), key, field)
	if cmd.Err() != nil {
		mlog.Errorf("redis HGet(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis HGet(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd = c.rdsClusterCli.HGet(context.Background(), key, field)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HGet(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis HGet(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
	err := cmd.Scan(val)
	if err != nil {
		mlog.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis Scan(%s) failed:%v", key, cmd.Err())
		return
	}
}

func (c *Cache) HGet(key string, field string) string {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		// return fmt.Errorf("redis don't have client")
		return ""
	}
	if key == "" || field == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return ""
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.HGet(context.Background(), key, field)
	if cmd.Err() != nil {
		mlog.Errorf("redis HGet(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis HGet(%s) failed:%v", key, cmd.Err())
		return ""
	}
	// } else {
	// 	cmd = c.rdsClusterCli.HGet(context.Background(), key, field)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HGet(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis HGet(%s) failed:%v", key, cmd.Err())
	// 		return ""
	// 	}
	// }
	return cmd.Val()
}

func (c *Cache) HSetData(key string, field string, val interface{}) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return fmt.Errorf("redis don't have cluster client")
	}
	if key == "" || field == "" {
		mlog.Error("invalid arg")
		return fmt.Errorf("invalid arg")
		// return
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.HSet(context.Background(), key, field, val)
	if cmd.Err() != nil {
		mlog.Errorf("redis HSet(%s %s %#v) failed:%v", key, field, val, cmd.Err())
		return fmt.Errorf("redis HSet(%s %s %#v) failed:%v", key, field, val, cmd.Err())
		// return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.HSet(context.Background(), key, field, val)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HSet(%s %s %#v) failed:%v", key, field, val, cmd.Err())
	// 		return fmt.Errorf("redis HSet(%s %s %#v) failed:%v", key, field, val, cmd.Err())
	// 		// return
	// 	}
	// }
	return nil
}

func (c *Cache) HMSetData(key string, attrs map[string]interface{}) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" || attrs == nil || len(attrs) == 0 {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	vals := make([]interface{}, 0)
	for k, v := range attrs {
		vals = append(vals, k)
		vals = append(vals, v)
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.HMSet(context.Background(), key, vals...)
	if cmd.Err() != nil {
		mlog.Errorf("redis HMSet(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis HMSet(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.HMSet(context.Background(), key, vals...)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HMSet(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis HMSet(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) HMSetStrData(key string, attrs map[string]string) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" || attrs == nil || len(attrs) == 0 {
		mlog.Errorf("invalid arg(key=%#v  attrs=%#v)", key, attrs)
		// return fmt.Errorf("invalid arg")
		return
	}
	vals := make([]interface{}, 0)
	for k, v := range attrs {
		vals = append(vals, k)
		vals = append(vals, v)
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.HMSet(context.Background(), key, vals...)
	if cmd.Err() != nil {
		mlog.Errorf("redis HMSet(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis HMSet(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.HMSet(context.Background(), key, vals...)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HMSet(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis HMSet(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) HExistsData(key string, field string) bool {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return false
	}
	if key == "" || field == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return false
	}
	var cmd *redis.BoolCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.HExists(context.Background(), key, field)
	if cmd.Err() != nil {
		mlog.Errorf("redis HExists(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis HExists(%s) failed:%v", key, cmd.Err())
		return false
	}
	// } else {
	// 	cmd = c.rdsClusterCli.HExists(context.Background(), key, field)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HExists(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis HExists(%s) failed:%v", key, cmd.Err())
	// 		return false
	// 	}
	// }
	return cmd.Val()
}

func (c *Cache) HDelData(key string, field string) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" || field == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.HDel(context.Background(), key, field)
	if cmd.Err() != nil {
		mlog.Errorf("redis HDel(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis HDel(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.HDel(context.Background(), key, field)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HDel(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis HDel(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) HIncrBy(key string, field string, inc int64) int64 {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return 0
	}
	if key == "" || field == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return 0
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.HIncrBy(context.Background(), key, field, inc)
	if cmd.Err() != nil {
		mlog.Errorf("redis HIncrBy(%s %s) failed:%v", key, field, cmd.Err())
		// return fmt.Errorf("redis HIncrBy(%s) failed:%v", key, cmd.Err())
		return 0
	}
	return cmd.Val()
	// } else {
	// 	cmd := c.rdsClusterCli.HIncrBy(context.Background(), key, field, inc)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HIncrBy(%s %s) failed:%v", key, field, cmd.Err())
	// 		// return fmt.Errorf("redis HIncrBy(%s) failed:%v", key, cmd.Err())
	// 		return 0
	// 	}
	// 	return cmd.Val()
	// }
}

func (c *Cache) HIncrByFloat(key string, field string, inc float64) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" || field == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.HIncrByFloat(context.Background(), key, field, inc)
	if cmd.Err() != nil {
		mlog.Errorf("redis HIncrBy(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis HIncrBy(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.HIncrByFloat(context.Background(), key, field, inc)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis HIncrBy(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis HIncrBy(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) ZAdd(key string, member interface{}, score float64) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.ZAdd(context.Background(), key, redis.Z{Score: score, Member: member})
	if cmd.Err() != nil {
		mlog.Errorf("redis ZAdd(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis ZAdd(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.ZAdd(context.Background(), key, redis.Z{Score: score, Member: member})
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis ZAdd(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis ZAdd(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) SAdd(key string, members ...interface{}) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.SAdd(context.Background(), key, members)
	if cmd.Err() != nil {
		mlog.Errorf("redis SAdd(%s %v) failed:%v", key, members, cmd.Err())
		// return fmt.Errorf("redis SAdd(%s) failed:%v", key, cmd.Err())
		return
	}
	// } else {
	// 	cmd := c.rdsClusterCli.SAdd(context.Background(), key, members)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis SAdd(%s %v) failed:%v", key, members, cmd.Err())
	// 		// return fmt.Errorf("redis SAdd(%s) failed:%v", key, cmd.Err())
	// 		return
	// 	}
	// }
}

func (c *Cache) SRem(key string, members ...interface{}) error {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have cluster client")
		return fmt.Errorf("redis don't have cluster client")
		// return
	}
	if key == "" {
		mlog.Errorf("invalid arg")
		return fmt.Errorf("invalid arg")
		// return
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.SRem(context.Background(), key, members)
	if cmd.Err() != nil {
		mlog.Errorf("redis SRem(%s) failed:%v", key, cmd.Err())
		return fmt.Errorf("redis SRem(%s) failed:%v", key, cmd.Err())
		// return
	}
	if int(cmd.Val()) != len(members) {
		mlog.Errorf("not all memebers(%v) remove from redis key(%v) ", members, key)
		return fmt.Errorf("not all memebers(%v) remove from redis key(%v) ", members, key)
	}
	// } else {
	// 	cmd := c.rdsClusterCli.SRem(context.Background(), key, members)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis SRem(%s) failed:%v", key, cmd.Err())
	// 		return fmt.Errorf("redis SRem(%s) failed:%v", key, cmd.Err())
	// 		// return
	// 	}
	// 	if int(cmd.Val()) != len(members) {
	// 		mlog.Errorf("not all memebers(%v) remove from redis key(%v) ", members, key)
	// 		return fmt.Errorf("not all memebers(%v) remove from redis key(%v) ", members, key)
	// 	}
	// }
	return nil
}

func (c *Cache) SMembers(key string) []string {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return nil
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return nil
	}
	var cmd *redis.StringSliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.SMembers(context.Background(), key)
	if cmd.Err() != nil {
		mlog.Errorf("redis SMembers(%s) failed:%v", key, cmd.Err())
		// return fmt.Errorf("redis SMembers(%s) failed:%v", key, cmd.Err())
		return nil
	}
	// } else {
	// 	cmd = c.rdsClusterCli.SMembers(context.Background(), key)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis SMembers(%s) failed:%v", key, cmd.Err())
	// 		// return fmt.Errorf("redis SMembers(%s) failed:%v", key, cmd.Err())
	// 		return nil
	// 	}
	// }

	return cmd.Val()
}

func (c *Cache) SIsMembers(key string, member interface{}) bool {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return false, fmt.Errorf("redis don't have cluster client")
		return false
	}
	if key == "" {
		mlog.Error("invalid arg")
		// return false, fmt.Errorf("invalid arg")
		return false
	}
	var cmd *redis.BoolCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.SIsMember(context.Background(), key, member)
	if cmd.Err() != nil {
		mlog.Errorf("redis SMembers(%s) failed:%v", key, cmd.Err())
		// return false, fmt.Errorf("redis SMembers(%s) failed:%v", key, cmd.Err())
		return false
	}
	// } else {
	// 	cmd = c.rdsClusterCli.SIsMember(context.Background(), key, member)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis SMembers(%s) failed:%v", key, cmd.Err())
	// 		// return false, fmt.Errorf("redis SMembers(%s) failed:%v", key, cmd.Err())
	// 		return false
	// 	}
	// }

	return cmd.Val()
}

func (c *Cache) Keys(patten string) []string {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return nil, fmt.Errorf("redis don't have cluster client")
		return nil
	}
	var cmd *redis.StringSliceCmd
	if cli, ok := c.rdsCli.(*redis.Client); ok {
		cmd = cli.Keys(context.Background(), patten)
		if cmd.Err() != nil {
			mlog.Errorf("redis Keys(%s) failed:%v", patten, cmd.Err())
			return nil
		}
		return cmd.Val()
	} else if clustercli, ok := c.rdsCli.(*redis.ClusterClient); ok {
		keys := make([]string, 0)
		var keyMap sync.Map
		err := clustercli.ForEachShard(context.Background(), func(ctx context.Context, client *redis.Client) error {
			cmd := client.Keys(ctx, patten)
			if cmd.Err() != nil {
				return cmd.Err()
			}
			for _, v := range cmd.Val() {
				if _, ok := keyMap.Load(v); !ok {
					keyMap.Store(v, 0)
				}
			}
			// keys = append(keys, cmd.Val()...)
			return nil
		})
		if err != nil {
			mlog.Errorf("redis ForEachMaster failed:%v", err)
			// return nil, fmt.Errorf("redis get keys in ForEachMaster failed:%v", err)
			return nil
		}
		keyMap.Range(func(key, val interface{}) bool {
			keys = append(keys, key.(string))
			return true
		})
		return keys
	}
	return nil
}
func (c *Cache) redisClientScan(cli *redis.Client, pattern string) map[string]int {
	if cli == nil || pattern == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("redis don't have cluster client")
		return nil
	}
	var cursor uint64 = 0
	var cmd *redis.ScanCmd
	var result map[string]int

	for {
		cmd = cli.Scan(context.Background(), cursor, pattern, 1000)
		if cmd.Err() != nil {
			mlog.Errorf("redis Scan(%s) failed:%v", pattern, cmd.Err())
			// return fmt.Errorf("redis SMembers(%s) failed:%v", key, cmd.Err())
			return nil
		}

		var keys []string
		keys, cursor = cmd.Val()
		for _, v := range keys {
			if result == nil {
				result = make(map[string]int)
			}
			if _, ok := result[v]; !ok {
				result[v] = 0
			}
			result[v]++
		}
		if cursor == 0 {
			return result
		}
	}
}

func (c *Cache) Scan(pattern string) []string {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return fmt.Errorf("redis don't have cluster client")
		return nil
	}
	if pattern == "" {
		mlog.Error("invalid arg")
		// return fmt.Errorf("invalid arg")
		return nil
	}

	var result map[string]int

	if cli, ok := c.rdsCli.(*redis.Client); ok {
		result = c.redisClientScan(cli, pattern)
	} else if clustercli, ok := c.rdsCli.(*redis.ClusterClient); ok {
		clustercli.ForEachMaster(context.Background(), func(ctx context.Context, client *redis.Client) error {
			tmp := c.redisClientScan(client, pattern)
			if result == nil {
				result = make(map[string]int)
			}
			for k, v := range tmp {
				if _, ok := result[k]; !ok {
					result[k] = 0
				}
				result[k] += v
			}
			return nil
		})
	}

	var ret []string
	for k := range result {
		if ret == nil {
			ret = make([]string, 0)
		}
		ret = append(ret, k)
	}
	return ret
}
func (c *Cache) Close() {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		// return nil, fmt.Errorf("redis don't have cluster client")
		return
	}
	// if c.rdsCli != nil {
	c.rdsCli.Close()
	// } else {
	// 	c.rdsClusterCli.Close()
	// }
}

func (c *Cache) Publish(channel string, message interface{}) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}
	// if c.rdsCli != nil {
	cmd := c.rdsCli.Publish(context.Background(), channel, message)
	if cmd.Err() != nil {
		mlog.Errorf("redis Publish(%s) failed:%v", channel, cmd.Err())
		return fmt.Errorf("redis Publish(%s) failed:%v", channel, cmd.Err())
	}
	// } else {
	// 	cmd := c.rdsClusterCli.Publish(context.Background(), channel, message)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Publish(%s) failed:%v", channel, cmd.Err())
	// 		return fmt.Errorf("redis Publish(%s) failed:%v", channel, cmd.Err())
	// 	}
	// }
	return nil
}

func (c *Cache) Subscribe(topic string) (*redis.PubSub, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}
	var sub *redis.PubSub
	// if c.rdsCli != nil {
	sub = c.rdsCli.Subscribe(context.Background(), topic)
	if _, err := sub.Receive(context.Background()); err != nil {
		mlog.Errorf("redis Subscribe(%s) failed:%v", topic, err)
		return nil, fmt.Errorf("redis Subscribe(%s) failed:%v", topic, err)
	}
	// } else {
	// 	sub = c.rdsClusterCli.Subscribe(context.Background(), topic)
	// 	if _, err := sub.Receive(context.Background()); err != nil {
	// 		mlog.Errorf("redis Subscribe(%s) failed:%v", topic, err)
	// 		return nil, fmt.Errorf("redis Subscribe(%s) failed:%v", topic, err)
	// 	}
	// }
	return sub, nil
}

func (c *Cache) StreamAdd(streamName string, data map[string]interface{}) (string, error) {

	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		return "", fmt.Errorf("redis don't have client")
	}
	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XAdd(context.Background(), &redis.XAddArgs{
		Stream: streamName,
		Values: data,
	})
	// } else {
	// 	cmd = c.rdsClusterCli.XAdd(context.Background(), &redis.XAddArgs{
	// 		Stream: streamName,
	// 		Values: data,
	// 	})
	// }
	return cmd.Result()
}

func (c *Cache) XAck(stream, group string, ids ...string) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}
	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XAck(context.Background(), stream, group, ids...)
	if cmd.Err() != nil {
		mlog.Errorf("redis XAck(%s %s %#v) failed:%v", stream, group, ids, cmd.Err())
		return fmt.Errorf("redis XAck(%s %s %#v) failed:%v", stream, group, ids, cmd.Err())
		// return
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XAck(context.Background(), stream, group, ids...)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XAck(%s %s %#v) failed:%v", stream, group, ids, cmd.Err())
	// 		return fmt.Errorf("redis XAck(%s %s %#v) failed:%v", stream, group, ids, cmd.Err())
	// 		// return
	// 	}
	// }
	return nil
}

func (c *Cache) XAdd(a *redis.XAddArgs) (string, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return "", fmt.Errorf("redis don't have client")
	}

	var cmd *redis.StringCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XAdd(context.Background(), a)
	if cmd.Err() != nil {
		mlog.Errorf("redis XAdd(%#v) failed:%v", a, cmd.Err())
		return "", fmt.Errorf("redis XAdd(%#v) failed:%v", a, cmd.Err())
		// return
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XAdd(context.Background(), a)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XAdd(%#v) failed:%v", a, cmd.Err())
	// 		return "", fmt.Errorf("redis XAdd(%#v) failed:%v", a, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}

func (c *Cache) XClaim(a *redis.XClaimArgs) (string, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return "", fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XMessageSliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XClaim(context.Background(), a)
	if cmd.Err() != nil {
		mlog.Errorf("redis XClaim(%#v) failed:%v", a, cmd.Err())
		return "", fmt.Errorf("redis XClaim(%#v) failed:%v", a, cmd.Err())
		// return
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XClaim(context.Background(), a)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XClaim(%#v) failed:%v", a, cmd.Err())
	// 		return "", fmt.Errorf("redis XClaim(%#v) failed:%v", a, cmd.Err())
	// 	}
	// }
	return cmd.String(), nil
}

func (c *Cache) XDel(stream string, ids ...string) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XDel(context.Background(), stream, ids...)
	if cmd.Err() != nil {
		mlog.Errorf("redis XDel(%s %#v) failed:%v", stream, ids, cmd.Err())
		return fmt.Errorf("redis XDel(%s %#v) failed:%v", stream, ids, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XDel(context.Background(), stream, ids...)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XDel(%s %#v) failed:%v", stream, ids, cmd.Err())
	// 		return fmt.Errorf("redis XDel(%s %#v) failed:%v", stream, ids, cmd.Err())
	// 	}
	// }
	return nil
}

func (c *Cache) XGroupCreate(stream, group, start string) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.StatusCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XGroupCreate(context.Background(), stream, group, start)
	if cmd.Err() != nil && cmd.Err().Error() != "BUSYGROUP Consumer Group name already exists" {
		mlog.Errorf("redis XGroupCreate(%s %s %s) failed:%v", stream, group, start, cmd.Err())
		return fmt.Errorf("redis XGroupCreate(%s %s %s) failed:%v", stream, group, start, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XGroupCreate(context.Background(), stream, group, start)
	// 	if cmd.Err() != nil && cmd.Err().Error() != "BUSYGROUP Consumer Group name already exists" {
	// 		mlog.Errorf("redis XGroupCreate(%s %s %s) failed:%v", stream, group, start, cmd.Err())
	// 		return fmt.Errorf("redis XGroupCreate(%s %s %s) failed:%v", stream, group, start, cmd.Err())
	// 	}
	// }
	// if cmd.Val() == "OK" {
	// 	return nil
	// } else if cmd.Val() == "BUSYGROUP Consumer Group name already exists" {
	// 	return nil
	// } else {
	// 	return fmt.Errorf("redis XGroupCreate(%s %s %s) failed:%v", stream, group, start, cmd.Val())
	// }
	return nil
}

func (c *Cache) XGroupCreateMkStream(stream, group, start string) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.StatusCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XGroupCreateMkStream(context.Background(), stream, group, start)
	if cmd.Err() != nil && cmd.Err().Error() != "BUSYGROUP Consumer Group name already exists" {
		mlog.Errorf("redis XGroupCreateMkStream(%s %s %s) failed:%v", stream, group, start, cmd.Err())
		return fmt.Errorf("redis XGroupCreateMkStream(%s %s %s) failed:%v", stream, group, start, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XGroupCreateMkStream(context.Background(), stream, group, start)
	// 	if cmd.Err() != nil && cmd.Err().Error() != "BUSYGROUP Consumer Group name already exists" {
	// 		mlog.Errorf("redis XGroupCreateMkStream(%s %s %s) failed:%v", stream, group, start, cmd.Err())
	// 		return fmt.Errorf("redis XGroupCreateMkStream(%s %s %s) failed:%v", stream, group, start, cmd.Err())
	// 	}
	// }
	// if cmd.Val() == "OK" {
	// 	return nil
	// } else if cmd.Val() == "BUSYGROUP Consumer Group name already exists" {
	// 	return nil
	// } else {
	// 	return fmt.Errorf("redis XGroupCreateMkStream(%s %s %s) failed:%v", stream, group, start, cmd.Val())
	// }
	return nil
}

func (c *Cache) XGroupCreateConsumer(stream, group, consumer string) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XGroupCreateConsumer(context.Background(), stream, group, consumer)
	if cmd.Err() != nil {
		mlog.Errorf("redis XGroupCreateConsumer(%s %s %s) failed:%v", stream, group, consumer, cmd.Err())
		return fmt.Errorf("redis XGroupCreateConsumer(%s %s %s) failed:%v", stream, group, consumer, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XGroupCreateConsumer(context.Background(), stream, group, consumer)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XGroupCreateConsumer(%s %s %s) failed:%v", stream, group, consumer, cmd.Err())
	// 		return fmt.Errorf("redis XGroupCreateConsumer(%s %s %s) failed:%v", stream, group, consumer, cmd.Err())
	// 	}
	// }
	return nil
}

func (c *Cache) XGroupDelConsumer(stream, group, consumer string) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XGroupDelConsumer(context.Background(), stream, group, consumer)
	if cmd.Err() != nil {
		mlog.Errorf("redis XGroupDelConsumer(%s %s %s) failed:%v", stream, group, consumer, cmd.Err())
		return fmt.Errorf("redis XGroupDelConsumer(%s %s %s) failed:%v", stream, group, consumer, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XGroupDelConsumer(context.Background(), stream, group, consumer)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XGroupDelConsumer(%s %s %s) failed:%v", stream, group, consumer, cmd.Err())
	// 		return fmt.Errorf("redis XGroupDelConsumer(%s %s %s) failed:%v", stream, group, consumer, cmd.Err())
	// 	}
	// }
	return nil
}
func (c *Cache) XGroupDestroy(stream, group string) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XGroupDestroy(context.Background(), stream, group)
	if cmd.Err() != nil {
		mlog.Errorf("redis XGroupDestroy(%s %s) failed:%v", stream, group, cmd.Err())
		return fmt.Errorf("redis XGroupDestroy(%s %s) failed:%v", stream, group, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XGroupDestroy(context.Background(), stream, group)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XGroupDestroy(%s %s) failed:%v", stream, group, cmd.Err())
	// 		return fmt.Errorf("redis XGroupDestroy(%s %s) failed:%v", stream, group, cmd.Err())
	// 	}
	// }
	return nil
}
func (c *Cache) XGroupSetID(stream, group, start string) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.StatusCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XGroupSetID(context.Background(), stream, group, start)
	if cmd.Err() != nil {
		mlog.Errorf("redis XGroupSetID(%s %s %s) failed:%v", stream, group, start, cmd.Err())
		return fmt.Errorf("redis XGroupSetID(%s %s %s) failed:%v", stream, group, start, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XGroupSetID(context.Background(), stream, group, start)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XGroupSetID(%s %s %s) failed:%v", stream, group, start, cmd.Err())
	// 		return fmt.Errorf("redis XGroupSetID(%s %s %s) failed:%v", stream, group, start, cmd.Err())
	// 	}
	// }
	if cmd.Val() == "OK" {
		return nil
	} else {
		return fmt.Errorf("redis XGroupSetID(%s %s %s) failed:%v", stream, group, start, cmd.Val())
	}
}

func (c *Cache) XInfoConsumers(stream, group string) ([]redis.XInfoConsumer, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XInfoConsumersCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XInfoConsumers(context.Background(), stream, group)
	if cmd.Err() != nil {
		mlog.Errorf("redis XInfoConsumers(%s %s) failed:%v", stream, group, cmd.Err())
		return nil, fmt.Errorf("redis XInfoConsumers(%s %s) failed:%v", stream, group, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XInfoConsumers(context.Background(), stream, group)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XInfoConsumers(%s %s) failed:%v", stream, group, cmd.Err())
	// 		return nil, fmt.Errorf("redis XInfoConsumers(%s %s) failed:%v", stream, group, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}

func (c *Cache) XInfoGroups(stream string) ([]redis.XInfoGroup, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XInfoGroupsCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XInfoGroups(context.Background(), stream)
	if cmd.Err() != nil {
		mlog.Errorf("redis XInfoGroups(%s) failed:%v", stream, cmd.Err())
		return nil, fmt.Errorf("redis XInfoGroups(%s) failed:%v", stream, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XInfoGroups(context.Background(), stream)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XInfoGroups(%s) failed:%v", stream, cmd.Err())
	// 		return nil, fmt.Errorf("redis XInfoGroups(%s) failed:%v", stream, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}

func (c *Cache) XInfoStream(stream string) (*redis.XInfoStream, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XInfoStreamCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XInfoStream(context.Background(), stream)
	if cmd.Err() != nil {
		mlog.Errorf("redis XInfoStream(%s) failed:%v", stream, cmd.Err())
		return nil, fmt.Errorf("redis XInfoStream(%s) failed:%v", stream, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XInfoStream(context.Background(), stream)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XInfoStream(%s) failed:%v", stream, cmd.Err())
	// 		return nil, fmt.Errorf("redis XInfoStream(%s) failed:%v", stream, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}

func (c *Cache) XLen(stream string) (int, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return 0, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XLen(context.Background(), stream)
	if cmd.Err() != nil {
		mlog.Errorf("redis XLen(%s) failed:%v", stream, cmd.Err())
		return 0, fmt.Errorf("redis XLen(%s) failed:%v", stream, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XLen(context.Background(), stream)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XLen(%s) failed:%v", stream, cmd.Err())
	// 		return 0, fmt.Errorf("redis XLen(%s) failed:%v", stream, cmd.Err())
	// 	}
	// }
	return int(cmd.Val()), nil
}

func (c *Cache) XPending(stream, group string) (*redis.XPending, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XPendingCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XPending(context.Background(), stream, group)
	if cmd.Err() != nil {
		mlog.Errorf("redis XPending(%s %s) failed:%v", stream, group, cmd.Err())
		return nil, fmt.Errorf("redis XPending(%s %s) failed:%v", stream, group, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XPending(context.Background(), stream, group)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XPending(%s %s) failed:%v", stream, group, cmd.Err())
	// 		return nil, fmt.Errorf("redis XPending(%s %s) failed:%v", stream, group, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}

func (c *Cache) XPendingExt(a *redis.XPendingExtArgs) ([]redis.XPendingExt, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XPendingExtCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XPendingExt(context.Background(), a)
	if cmd.Err() != nil {
		mlog.Errorf("redis XPendingExt(%#v) failed:%v", a, cmd.Err())
		return nil, fmt.Errorf("redis XPendingExt(%#v) failed:%v", a, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XPendingExt(context.Background(), a)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XPendingExt(%#v) failed:%v", a, cmd.Err())
	// 		return nil, fmt.Errorf("redis XPendingExt(%#v) failed:%v", a, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}
func (c *Cache) XRange(stream, start, stop string) ([]redis.XMessage, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XMessageSliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XRange(context.Background(), stream, start, stop)
	if cmd.Err() != nil {
		mlog.Errorf("redis XRange(%s %s %s) failed:%v", stream, start, stop, cmd.Err())
		return nil, fmt.Errorf("redis XRange(%s %s %s) failed:%v", stream, start, stop, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XRange(context.Background(), stream, start, stop)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XRange(%s %s %s) failed:%v", stream, start, stop, cmd.Err())
	// 		return nil, fmt.Errorf("redis XRange(%s %s %s) failed:%v", stream, start, stop, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}
func (c *Cache) XRangeN(stream, start, stop string, count int64) ([]redis.XMessage, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XMessageSliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XRangeN(context.Background(), stream, start, stop, count)
	if cmd.Err() != nil {
		mlog.Errorf("redis XRangeN(%s %s %s %d) failed:%v", stream, start, stop, count, cmd.Err())
		return nil, fmt.Errorf("redis XRangeN(%s %s %s %d) failed:%v", stream, start, stop, count, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XRangeN(context.Background(), stream, start, stop, count)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XRangeN(%s %s %s %d) failed:%v", stream, start, stop, count, cmd.Err())
	// 		return nil, fmt.Errorf("redis XRangeN(%s %s %s %d) failed:%v", stream, start, stop, count, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}
func (c *Cache) XRead(a *redis.XReadArgs) ([]redis.XStream, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XStreamSliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XRead(context.Background(), a)
	if cmd.Err() != nil && cmd.Err().Error() != "redis: nil" {
		mlog.Errorf("redis XRead(%#v) failed:%v", a, cmd.Err())
		return nil, fmt.Errorf("redis XRead(%#v) failed:%v", a, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XRead(context.Background(), a)
	// 	if cmd.Err() != nil && cmd.Err().Error() != "redis: nil" {
	// 		mlog.Errorf("redis XRead(%#v) failed:%v", a, cmd.Err())
	// 		return nil, fmt.Errorf("redis XRead(%#v) failed:%v", a, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}
func (c *Cache) XReadGroup(a *redis.XReadGroupArgs) ([]redis.XStream, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XStreamSliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XReadGroup(context.Background(), a)
	if cmd.Err() != nil && cmd.Err().Error() != "redis: nil" {
		mlog.Errorf("redis XReadGroup(%#v) failed:%v", a, cmd.Err())
		return nil, fmt.Errorf("redis XReadGroup(%#v) failed:%v", a, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XReadGroup(context.Background(), a)
	// 	if cmd.Err() != nil && cmd.Err().Error() != "redis: nil" {
	// 		mlog.Errorf("redis XReadGroup(%#v) failed:%v", a, cmd.Err())
	// 		return nil, fmt.Errorf("redis XReadGroup(%#v) failed:%v", a, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}
func (c *Cache) XReadStreams(streams ...string) ([]redis.XStream, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XStreamSliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XReadStreams(context.Background(), streams...)
	if cmd.Err() != nil {
		mlog.Errorf("redis XReadStreams(%#v) failed:%v", streams, cmd.Err())
		return nil, fmt.Errorf("redis XReadStreams(%#v) failed:%v", streams, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XReadStreams(context.Background(), streams...)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XReadStreams(%#v) failed:%v", streams, cmd.Err())
	// 		return nil, fmt.Errorf("redis XReadStreams(%#v) failed:%v", streams, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}
func (c *Cache) XRevRange(stream, start, stop string) ([]redis.XMessage, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XMessageSliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XRevRange(context.Background(), stream, start, stop)
	if cmd.Err() != nil {
		mlog.Errorf("redis XRevRange(%#v) failed:%v", []string{stream, start, stop}, cmd.Err())
		return nil, fmt.Errorf("redis XRevRange(%#v) failed:%v", []string{stream, start, stop}, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XRevRange(context.Background(), stream, start, stop)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XRevRange(%#v) failed:%v", []string{stream, start, stop}, cmd.Err())
	// 		return nil, fmt.Errorf("redis XRevRange(%#v) failed:%v", []string{stream, start, stop}, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}
func (c *Cache) XRevRangeN(stream, start, stop string, count int64) ([]redis.XMessage, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return nil, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.XMessageSliceCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XRevRangeN(context.Background(), stream, start, stop, count)
	if cmd.Err() != nil {
		mlog.Errorf("redis XRevRangeN(%#v) failed:%v", []interface{}{stream, start, stop, count}, cmd.Err())
		return nil, fmt.Errorf("redis XRevRangeN(%#v) failed:%v", []interface{}{stream, start, stop, count}, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XRevRangeN(context.Background(), stream, start, stop, count)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XRevRangeN(%#v) failed:%v", []interface{}{stream, start, stop, count}, cmd.Err())
	// 		return nil, fmt.Errorf("redis XRevRangeN(%#v) failed:%v", []interface{}{stream, start, stop, count}, cmd.Err())
	// 	}
	// }
	return cmd.Val(), nil
}
func (c *Cache) XTrimMaxLen(stream string, count int64) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XTrimMaxLen(context.Background(), stream, count)
	if cmd.Err() != nil {
		mlog.Errorf("redis XTrimMaxLen(%#v) failed:%v", []interface{}{stream, count}, cmd.Err())
		return fmt.Errorf("redis XTrimMaxLen(%#v) failed:%v", []interface{}{stream, count}, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XTrimMaxLen(context.Background(), stream, count)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XTrimMaxLen(%#v) failed:%v", []interface{}{stream, count}, cmd.Err())
	// 		return fmt.Errorf("redis XTrimMaxLen(%#v) failed:%v", []interface{}{stream, count}, cmd.Err())
	// 	}
	// }
	return nil
}
func (c *Cache) XTrimMaxLenApprox(stream string, maxLen, limit int64) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XTrimMaxLenApprox(context.Background(), stream, maxLen, limit)
	if cmd.Err() != nil {
		mlog.Errorf("redis XTrimMaxLenApprox(%#v) failed:%v", []interface{}{stream, maxLen, limit}, cmd.Err())
		return fmt.Errorf("redis XTrimMaxLenApprox(%#v) failed:%v", []interface{}{stream, maxLen, limit}, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XTrimMaxLenApprox(context.Background(), stream, maxLen, limit)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XTrimMaxLenApprox(%#v) failed:%v", []interface{}{stream, maxLen, limit}, cmd.Err())
	// 		return fmt.Errorf("redis XTrimMaxLenApprox(%#v) failed:%v", []interface{}{stream, maxLen, limit}, cmd.Err())
	// 	}
	// }
	return nil
}
func (c *Cache) XTrimMinID(stream, minID string) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XTrimMinID(context.Background(), stream, minID)
	if cmd.Err() != nil {
		mlog.Errorf("redis XTrimMinID(%#v) failed:%v", []interface{}{stream, minID}, cmd.Err())
		return fmt.Errorf("redis XTrimMinID(%#v) failed:%v", []interface{}{stream, minID}, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XTrimMinID(context.Background(), stream, minID)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XTrimMinID(%#v) failed:%v", []interface{}{stream, minID}, cmd.Err())
	// 		return fmt.Errorf("redis XTrimMinID(%#v) failed:%v", []interface{}{stream, minID}, cmd.Err())
	// 	}
	// }
	return nil
}
func (c *Cache) XTrimMinIDApprox(stream, minID string, limit int64) error {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return fmt.Errorf("redis don't have client")
	}

	var cmd *redis.IntCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.XTrimMinIDApprox(context.Background(), stream, minID, limit)
	if cmd.Err() != nil {
		mlog.Errorf("redis XTrimMinIDApprox(%#v) failed:%v", []interface{}{stream, minID, limit}, cmd.Err())
		return fmt.Errorf("redis XTrimMinIDApprox(%#v) failed:%v", []interface{}{stream, minID, limit}, cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.XTrimMinIDApprox(context.Background(), stream, minID, limit)
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis XTrimMinIDApprox(%#v) failed:%v", []interface{}{stream, minID, limit}, cmd.Err())
	// 		return fmt.Errorf("redis XTrimMinIDApprox(%#v) failed:%v", []interface{}{stream, minID, limit}, cmd.Err())
	// 	}
	// }
	return nil
}
func (c *Cache) Time() (time.Time, error) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have client")
		return time.Time{}, fmt.Errorf("redis don't have client")
	}

	var cmd *redis.TimeCmd
	// if c.rdsCli != nil {
	cmd = c.rdsCli.Time(context.Background())
	if cmd.Err() != nil {
		mlog.Errorf("redis Time failed:%v", cmd.Err())
		return time.Time{}, fmt.Errorf("redis Time failed:%v", cmd.Err())
	}
	// } else {
	// 	cmd = c.rdsClusterCli.Time(context.Background())
	// 	if cmd.Err() != nil {
	// 		mlog.Errorf("redis Time failed:%v", cmd.Err())
	// 		return time.Time{}, fmt.Errorf("redis Time failed:%v", cmd.Err())
	// 	}
	// }

	return cmd.Val(), nil
}

func (c *Cache) DistributedLock_lock(key string, expiration time.Duration) bool {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		return false
	}
	// if c.rdsCli != nil {
	res, err := c.rdsCli.SetNX(context.Background(), key, "locked", expiration).Result()
	if err != nil {
		mlog.Errorf("redis SetNX:%s ,err:%s", key, err.Error())
		return false
	}
	return res
	// } else {
	// 	res, err := c.rdsClusterCli.SetNX(context.Background(), key, "locked", expiration).Result()
	// 	if err != nil {
	// 		mlog.Errorf("redis del:%s ,err:%s", key, err.Error())
	// 		return false
	// 	}
	// 	return res
	// }
}

func (c *Cache) DistributedLock_unlock(key string) {
	if c.rdsCli == nil {
		mlog.Error("redis don't have cluster client")
		return
	}
	// 尝试删除锁
	// if c.rdsCli != nil {
	_, err := c.rdsCli.Del(context.Background(), key).Result()
	if err != nil {
		mlog.Errorf("redis SetNX:%s ,err:%s", key, err.Error())
	}
	// } else {
	// 	_, err := c.rdsClusterCli.Del(context.Background(), key).Result()
	// 	if err != nil {
	// 		mlog.Errorf("redis del:%s ,err:%s", key, err.Error())
	// 	}
	// }
}
func (c *Cache) Client() redis.UniversalClient {
	return c.rdsCli
}

// ZRemRangeByScore removes members in the sorted set `key` whose score is
// within the given range. Scores are strings to allow "-inf"/"+inf"/numeric.
func (c *Cache) ZRemRangeByScore(key, min, max string) int64 {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return 0
	}
	n, err := c.rdsCli.ZRemRangeByScore(context.Background(), key, min, max).Result()
	if err != nil {
		mlog.Errorf("ZRemRangeByScore(%s,%s,%s) failed: %v", key, min, max, err)
		return 0
	}
	return n
}

// ZCard returns the cardinality of the sorted set stored at key.
func (c *Cache) ZCard(key string) int64 {
	if c.rdsCli == nil {
		mlog.Errorf("redis don't have client")
		return 0
	}
	n, err := c.rdsCli.ZCard(context.Background(), key).Result()
	if err != nil {
		mlog.Errorf("ZCard(%s) failed: %v", key, err)
		return 0
	}
	return n
}
