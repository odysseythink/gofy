package cache

import (
	"testing"

	"github.com/odysseythink/confy"
)

func TestCluster(t *testing.T) {
	confy.SetConfigFile("redis.yml")
	err := confy.ReadInConfig()
	if err != nil {
		t.Errorf("read config file(redis.yml) failed: %v", err)
		return
	}
	c := &Cache{}
	c.Init()
	keys := c.Keys("bot:*")
	t.Logf("bot:*=%#v", len(keys))
	t.Logf("ttl=%#v", c.TTL("Termcall-Event-Detail:20231018:0-1--ff11bc55-fe6a-47d9-9707-139579875157"))
}
