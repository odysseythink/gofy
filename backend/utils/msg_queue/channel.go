package msgqueue

import (
	"cache"
	"context"
	"sync"

	"github.com/odysseythink/mlog"
)

type Channel struct {
	Name        string `json:"name"`
	RedisKey    string `json:"redis_key"`
	subscribers sync.Map
}

func NewChannel(name string) *Channel {
	c := &Channel{
		Name: name,
	}
	c.RedisKey = DEFAULT_MSG_QUEUE_CHANNEL_PREFIX + c.Name
	return c
}

func (c *Channel) AddSubscriber(suber *Subscriber) {
	mlog.Debugf("channel add subscriber=%#v", suber)
	c.subscribers.Store(suber.ID, suber)
	cache.Instance().SAdd(c.RedisKey, suber.ID)
}

func (c *Channel) Init(args ...any) error {
	return nil
}

func (c *Channel) RunOnce(ctx context.Context) error {
	c.subscribers.Range(func(key any, value any) bool {
		subscriber := value.(*Subscriber)
		subscriber.RunOnce(ctx)
		return true
	})
	return nil
}

func (c *Channel) Destroy() {
	c.subscribers.Range(func(key any, value any) bool {
		subscriber := value.(*Subscriber)
		subscriber.Destroy()
		cache.Instance().SRem(c.RedisKey, subscriber.ID)
		return true
	})
}
func (c *Channel) UserData() any {
	return nil
}
