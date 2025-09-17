package msgqueue

import (
	"cache"
	"fmt"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestMessageQueue(t *testing.T) {
	viper.SetConfigFile("redis.yml")
	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("read config file(redis.yml) failed: %v", err)
		return
	}
	c := cache.Instance()
	c.Init()
	group1Cnt1 := 0
	Instance().Subscribe("bot:test-stream", "group-1", func(topic, channel, message, msgid string) {
		t.Logf("topic(%s)  channel(%s) receieve message:%s", topic, channel, message)
		group1Cnt1++
	})
	group1Cnt2 := 0
	Instance().Subscribe("bot:test-stream", "group-1", func(topic, channel, message, msgid string) {
		t.Logf("topic(%s)  channel(%s) receieve message:%s", topic, channel, message)
		group1Cnt2++
	})
	group2Cnt1 := 0
	Instance().Subscribe("bot:test-stream", "group-2", func(topic, channel, message, msgid string) {
		t.Logf("topic(%s)  channel(%s) receieve message:%s", topic, channel, message)
		group2Cnt1++
	})

	for idx := 0; idx < 100; idx++ {
		Instance().Publish("bot:test-stream", fmt.Sprintf("hello %d", idx))
	}
	time.Sleep(15 * time.Second)
	t.Logf("------------group1(%d %d) group2(%d)", group1Cnt1, group1Cnt2, group2Cnt1)
	keys := c.Keys("bot:*")
	t.Logf("bot:*=%#v", len(keys))
	t.Logf("ttl=%#v", c.TTL("Termcall-Event-Detail:20231018:0-1--ff11bc55-fe6a-47d9-9707-139579875157"))
	c.Destroy()
}
