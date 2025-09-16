package distributelock

import (
	"context"
	"errors"
	"sync"
	"time"

	"mlib.com/gofy/server/cache"

	uuid "github.com/satori/go.uuid"
	"mlib.com/confy"
	"mlib.com/mlog"
)

const (
	DISTRIBUTE_LOCK_KEY_PREFIX      = "distribute_lock:"
	DEFAULT_DISTRIBUTE_LOCK_TIMEOUT = 30
)

type DistributeLock struct {
	id        string
	keys      sync.Map
	refreshAt time.Time
}

func New() *DistributeLock {
	return &DistributeLock{
		id: uuid.NewV4().String(),
	}
}

func (lock *DistributeLock) TryLock(key string, timeout time.Duration) bool {
	if lock.id == "" {
		mlog.Error("distributeLockID not init")
		return false
	}
	if key == "" {
		mlog.Error("invalid arg")
		return false
	}
	strLockKey := DISTRIBUTE_LOCK_KEY_PREFIX + key
	err, exist := cache.Instance().SetNX(strLockKey, lock.id, time.Duration(DEFAULT_DISTRIBUTE_LOCK_TIMEOUT)*time.Second)
	if err != nil {
		mlog.Errorf("SetNX failed:%v", err)
		return false
	}

	if !exist {
		if timeout <= 0 {
			return false
		} else {
			time.Sleep(timeout)
			err, exist = cache.Instance().SetNX(strLockKey, lock.id, time.Duration(DEFAULT_DISTRIBUTE_LOCK_TIMEOUT)*time.Second)
			if err != nil {
				mlog.Errorf("SetNX failed:%v", err)
				return false
			}

			if exist {
				lock.keys.Store(strLockKey, true)
				return true
			} else {
				return false
			}
		}
	} else {
		lock.keys.Store(strLockKey, true)
		return true
	}

	// otherID := cache.Instance().GetString(strLockKey)
	// if otherID != lock.id {
	// 	if timeout <= 0 {
	// 		return false
	// 	} else {
	// 		time.Sleep(timeout)
	// 		err = cache.Instance().SetNX(strLockKey, lock.id, time.Duration(DEFAULT_DISTRIBUTE_LOCK_TIMEOUT)*time.Second)
	// 		if err != nil {
	// 			mlog.Errorf("SetNX failed:%v", err)
	// 			return false
	// 		}
	// 		otherID = cache.Instance().GetString(strLockKey)
	// 		if otherID == lock.id {
	// 			lock.keys.Store(strLockKey, true)
	// 			return true
	// 		} else {
	// 			return false
	// 		}
	// 	}
	// } else {
	// 	lock.keys.Store(strLockKey, true)
	// 	return true
	// }
}

func (lock *DistributeLock) UnLock(key string) {
	if lock.id == "" {
		mlog.Error("distributeLockID not init")
		return
	}
	if key == "" {
		mlog.Error("invalid arg")
		return
	}

	strLockKey := DISTRIBUTE_LOCK_KEY_PREFIX + key
	otherID := cache.Instance().GetString(strLockKey)
	if otherID == lock.id {
		lock.keys.Delete(strLockKey)
		cache.Instance().DelKey(strLockKey)
	}
}

func (lock *DistributeLock) Init(args ...any) error {
	if cache.Instance().Client() == nil {
		mlog.Error("no redis client provide, can't use DistributeLock")
		return errors.New("no redis client provide, can't use DistributeLock")
	}
	return nil
}

func (lock *DistributeLock) RunOnce(ctx context.Context) error {
	if lock.id == "" {
		mlog.Error("distributeLockID not init")
		return nil
	}
	now := time.Now()
	if lock.refreshAt.IsZero() {
		lock.refreshAt = now
	}
	if lock.refreshAt.Sub(now).Seconds() <= 0 {
		lock.keys.Range(func(key, value any) bool {
			if strLockKey, ok := key.(string); ok {
				cache.Instance().ExpireKey(strLockKey, DEFAULT_DISTRIBUTE_LOCK_TIMEOUT)
			}
			return true
		})
		lock.refreshAt = now.Add(time.Duration(confy.GetWithDefault[int]("distribute_lock.refresh_period", 1)) * time.Second)
	}

	return nil
}

func (lock *DistributeLock) Destroy() {

}
func (lock *DistributeLock) UserData() any {
	return nil
}

var (
	gOnce     sync.Once
	gInstance *DistributeLock
)

func Instance() *DistributeLock {
	gOnce.Do(func() {
		gInstance = New()
	})
	return gInstance
}
