package msgqueue

import (
	"cache"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/odysseythink/mlog"
	"github.com/odysseythink/mrun"
	uuid "github.com/satori/go.uuid"
)

type MsgQueue struct {
	wg            sync.WaitGroup
	ID            string
	topics        sync.Map
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
}

func New() *MsgQueue {
	mq := &MsgQueue{
		ID: uuid.NewV4().String(),
	}
	mq.ctx, mq.ctxCancelFunc = context.WithCancel(context.Background())
	return mq
}

func (mq *MsgQueue) Publish(topic_name string, message any) error {
	if topic_name == "" || message == nil {
		mlog.Errorf("invalid arg")
		return fmt.Errorf("invalid arg")
	}
	payload := ""
	switch message := message.(type) {
	case string:
		payload = message
	default:
		data, err := json.Marshal(message)
		if err != nil {
			mlog.Errorf("json.Marshal(%#v) failed:%v", message, err)
			return fmt.Errorf("json.Marshal(%#v) failed:%v", message, err)
		}
		payload = string(data)
	}
	strTopicKey := DEFAULT_MSG_QUEUE_TOPIC_PREFIX + topic_name
	channel_names := cache.Instance().SMembers(strTopicKey)
	for _, channel_name := range channel_names {
		strChannelKey := DEFAULT_MSG_QUEUE_CHANNEL_PREFIX + channel_name
		for {
			subscriber_ids := cache.Instance().SMembers(strChannelKey)
			if len(subscriber_ids) > 0 {
				idx := rand.Intn(len(subscriber_ids))
				matched_subscriber_id := subscriber_ids[idx]
				suber_redis_key := DEFAULT_MSG_QUEUE_SUBSCRIBER_PREFIX + matched_subscriber_id
				if cache.Instance().ExistsKey(suber_redis_key) {
					sub_redis_key := DEFAULT_MSG_QUEUE_PREFIX + topic_name + ":" + channel_name + ":" + matched_subscriber_id
					mlog.Debugf("publish payload(%s) to key=%s", payload, sub_redis_key)
					cache.Instance().Publish(sub_redis_key, payload)
					break
				} else {
					cache.Instance().SRem(strChannelKey, matched_subscriber_id)
				}
			} else {
				break
			}
		}
	}
	return nil
}

func (mq *MsgQueue) Subscribe(topic_name, channel_name string, cb func(topic, channel, message, msgid string)) error {
	if topic_name == "" || channel_name == "" {
		mlog.Errorf("invalid arg")
		return fmt.Errorf("invalid arg")
	}
	var topic *Topic
	if val, ok := mq.topics.Load(topic_name); ok {
		topic = val.(*Topic)
	} else {
		topic = NewTopic(topic_name, channel_name)
	}
	defer mq.topics.Store(topic_name, topic)
	var channel *Channel
	if val, ok := topic.channels.Load(channel_name); ok {
		channel = val.(*Channel)
	} else {
		channel = NewChannel(channel_name)
	}

	subscriber, err := NewSubscriber(topic_name, channel_name, cb)
	if err != nil {
		mlog.Errorf("create subscriber failed:%v", err)
		return errors.New("create subscriber failed")
	}
	channel.AddSubscriber(subscriber)
	topic.AddChannel(channel)
	mq.wg.Add(1)
	go func() {
		defer func() {
			mq.wg.Done()
			mlog.Infof("subscriber exit")
		}()
		for {
			select {
			case <-mq.ctx.Done():
				return
			case msg := <-subscriber.pubsub.Channel():
				mlog.Infof("----------receive message=%s", msg)
				mrun.WorkerSubmit(func() {
					subscriber.cb(subscriber.Topic, subscriber.Channel, msg.Payload, "")
				})
			}
		}
	}()

	return nil
}

func (mq *MsgQueue) Init(args ...any) error {
	if cache.Instance().Client() == nil {
		mlog.Error("no redis client provide, can't use DistributeLock")
		return errors.New("no redis client provide, can't use DistributeLock")
	}
	if mq.ctx == nil {
		mq.ctx, mq.ctxCancelFunc = context.WithCancel(context.Background())
	}
	return nil
}

func (mq *MsgQueue) RunOnce(ctx context.Context) error {
	mq.topics.Range(func(key any, value any) bool {
		t := value.(*Topic)
		t.RunOnce(ctx)
		return true
	})
	return nil
}

func (mq *MsgQueue) Destroy() {
	mlog.Debugf("msg queue destroy")
	mq.ctxCancelFunc()
	mq.wg.Wait()
	mq.topics.Range(func(key any, value any) bool {
		t := value.(*Topic)
		t.Destroy()
		return true
	})
}
func (mq *MsgQueue) UserData() any {
	return nil
}

var (
	gOnce     sync.Once
	gInstance *MsgQueue
)

func Instance() *MsgQueue {
	gOnce.Do(func() {
		gInstance = New()
	})
	return gInstance
}

func init() {
	mrun.Register(Instance(), []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(5), mrun.NewModuleRunPeriodOption(1000)}, nil)
}
