package msgqueue

import (
	"cache"
	"context"
	"sync"

	"mlib.com/mlog"
)

type Topic struct {
	Name     string `json:"name"`
	RedisKey string `json:"redis_key"`
	channels sync.Map
}

func NewTopic(topic_name, channel_name string) *Topic {
	t := &Topic{
		Name: topic_name,
	}
	t.RedisKey = DEFAULT_MSG_QUEUE_TOPIC_PREFIX + topic_name
	c := NewChannel(channel_name)
	t.AddChannel(c)

	return t
}

func (t *Topic) AddChannel(c *Channel) {
	t.channels.Store(c.Name, c)
	cache.Instance().SAdd(t.RedisKey, c.Name)
}

func (t *Topic) Init(args ...any) error {
	return nil
}

func (t *Topic) RunOnce(ctx context.Context) error {
	t.channels.Range(func(key any, value any) bool {
		c := value.(*Channel)
		c.RunOnce(ctx)
		return true
	})
	return nil
}

func (t *Topic) Destroy() {
	mlog.Debugf("topic(%s) destroy", t.Name)
	t.channels.Range(func(key any, value any) bool {
		c := value.(*Channel)
		c.Destroy()
		if !cache.Instance().ExistsKey(c.RedisKey) {
			cache.Instance().SRem(t.RedisKey, c.Name)
		}
		return true
	})
}
func (t *Topic) UserData() any {
	return nil
}
