package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/odysseythink/gofy/backend/cache"
)

// CommandType defines workflow control commands.
type CommandType string

const (
	CommandAbort          CommandType = "abort"
	CommandPause          CommandType = "pause"
	CommandResume         CommandType = "resume"
	CommandUpdateVariable CommandType = "update_variable"
)

// Command is a control message sent to a running workflow.
type Command struct {
	Type      CommandType    `json:"type"`
	Data      map[string]any `json:"data,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// CommandChannel allows sending control commands to a workflow execution.
type CommandChannel interface {
	Send(cmd *Command) error
	Receive(ctx context.Context) (*Command, error)
	Close() error
}

// InMemoryCommandChannel uses Go channels (same process only).
type InMemoryCommandChannel struct {
	ch   chan *Command
	once sync.Once
}

func NewInMemoryCommandChannel(bufferSize int) *InMemoryCommandChannel {
	return &InMemoryCommandChannel{ch: make(chan *Command, bufferSize)}
}

func (c *InMemoryCommandChannel) Send(cmd *Command) error {
	select {
	case c.ch <- cmd:
		return nil
	default:
		return fmt.Errorf("command channel full")
	}
}

func (c *InMemoryCommandChannel) Receive(ctx context.Context) (*Command, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case cmd := <-c.ch:
		return cmd, nil
	}
}

func (c *InMemoryCommandChannel) Close() error {
	c.once.Do(func() { close(c.ch) })
	return nil
}

// RedisCommandChannel uses Redis pub/sub for cross-process commands.
type RedisCommandChannel struct {
	key string
}

func NewRedisCommandChannel(workflowRunID string) *RedisCommandChannel {
	return &RedisCommandChannel{key: fmt.Sprintf("gofy:workflow:cmd:%s", workflowRunID)}
}

func (c *RedisCommandChannel) Send(cmd *Command) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return cache.Instance().Client().Publish(context.Background(), c.key, string(data)).Err()
}

func (c *RedisCommandChannel) Receive(ctx context.Context) (*Command, error) {
	sub := cache.Instance().Client().Subscribe(ctx, c.key)
	defer sub.Close()

	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return nil, err
	}

	var cmd Command
	if err := json.Unmarshal([]byte(msg.Payload), &cmd); err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (c *RedisCommandChannel) Close() error {
	return nil
}
