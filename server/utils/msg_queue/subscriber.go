package msgqueue

import (
	"cache"
	"context"
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"mlib.com/mlog"

	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	ID          string `json:"id"`
	Topic       string `json:"topic"`
	Channel     string `json:"channel"`
	RedisKey    string `json:"redis_key"`
	SubRedisKey string `json:"sub_redis_key"`
	pubsub      *redis.PubSub
	cb          func(topic, channel, message, msgid string)
}

func NewSubscriber(topic, channel string, cb func(topic, channel, message, msgid string)) (*Subscriber, error) {
	suber := &Subscriber{
		ID:      uuid.NewV4().String(),
		Topic:   topic,
		Channel: channel,
		cb:      cb,
	}
	suber.RedisKey = DEFAULT_MSG_QUEUE_SUBSCRIBER_PREFIX + suber.ID
	suber.SubRedisKey = DEFAULT_MSG_QUEUE_PREFIX + topic + ":" + channel + ":" + suber.ID
	bindata, _ := json.Marshal(suber)
	cache.Instance().SetEx(suber.RedisKey, string(bindata), DEFAULT_MSG_QUEUE_SUBSCRIBER_IDLE_MAX*time.Second)
	pubsub, err := cache.Instance().Subscribe(suber.SubRedisKey)
	if err != nil {
		mlog.Errorf("Subscribe(%#v) failed:%v", suber.SubRedisKey, err)
		return nil, fmt.Errorf("Subscribe(%#v) failed:%v", suber.SubRedisKey, err)
	}
	suber.pubsub = pubsub
	return suber, nil
}

func (suber *Subscriber) Init(args ...any) error {
	bindata, _ := json.Marshal(suber)
	cache.Instance().SetEx(suber.RedisKey, string(bindata), DEFAULT_MSG_QUEUE_SUBSCRIBER_IDLE_MAX*time.Second)
	return nil
}

func (suber *Subscriber) RunOnce(ctx context.Context) error {
	if suber.RedisKey != "" {
		cache.Instance().ExpireKey(suber.RedisKey, DEFAULT_MSG_QUEUE_SUBSCRIBER_IDLE_MAX)
	}
	return nil
}

func (suber *Subscriber) Destroy() {
	mlog.Debugf("Subscriber(%s) destroy", suber.ID)
	if suber.RedisKey != "" {
		cache.Instance().DelKey(suber.RedisKey)
	}
}
func (suber *Subscriber) UserData() any {
	return nil
}
