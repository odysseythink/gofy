package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"mlib.com/gofy/server/cache"
)

// ReadyNode represents a node ready for execution.
type ReadyNode struct {
	NodeID           string         `json:"node_id"`
	NodeType         string         `json:"node_type"`
	PreviousNodeID   string         `json:"previous_node_id,omitempty"`
	EdgeSourceHandle string         `json:"edge_source_handle,omitempty"`
	Priority         int            `json:"priority"`
	Data             map[string]any `json:"data,omitempty"`
}

// ReadyQueue schedules nodes for execution.
type ReadyQueue interface {
	Push(node *ReadyNode) error
	Pop(ctx context.Context) (*ReadyNode, error)
	Len() int
	State() *ReadyQueueState
	RestoreState(state *ReadyQueueState) error
}

// ReadyQueueState is serializable state for persistence.
type ReadyQueueState struct {
	PendingNodes []*ReadyNode `json:"pending_nodes"`
}

// InMemoryReadyQueue uses a slice-based queue (single process).
type InMemoryReadyQueue struct {
	mu    sync.Mutex
	nodes []*ReadyNode
	ch    chan struct{}
}

func NewInMemoryReadyQueue() *InMemoryReadyQueue {
	return &InMemoryReadyQueue{
		nodes: make([]*ReadyNode, 0),
		ch:    make(chan struct{}, 1),
	}
}

func (q *InMemoryReadyQueue) Push(node *ReadyNode) error {
	q.mu.Lock()
	q.nodes = append(q.nodes, node)
	q.mu.Unlock()
	select {
	case q.ch <- struct{}{}:
	default:
	}
	return nil
}

func (q *InMemoryReadyQueue) Pop(ctx context.Context) (*ReadyNode, error) {
	for {
		q.mu.Lock()
		if len(q.nodes) > 0 {
			node := q.nodes[0]
			q.nodes = q.nodes[1:]
			q.mu.Unlock()
			return node, nil
		}
		q.mu.Unlock()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-q.ch:
		}
	}
}

func (q *InMemoryReadyQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.nodes)
}

func (q *InMemoryReadyQueue) State() *ReadyQueueState {
	q.mu.Lock()
	defer q.mu.Unlock()
	nodes := make([]*ReadyNode, len(q.nodes))
	copy(nodes, q.nodes)
	return &ReadyQueueState{PendingNodes: nodes}
}

func (q *InMemoryReadyQueue) RestoreState(state *ReadyQueueState) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.nodes = state.PendingNodes
	return nil
}

// RedisReadyQueue uses Redis list for distributed node scheduling.
type RedisReadyQueue struct {
	key string
}

func NewRedisReadyQueue(workflowRunID string) *RedisReadyQueue {
	return &RedisReadyQueue{key: fmt.Sprintf("gofy:workflow:ready:%s", workflowRunID)}
}

func (q *RedisReadyQueue) Push(node *ReadyNode) error {
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}
	return cache.Instance().Client().RPush(context.Background(), q.key, string(data)).Err()
}

func (q *RedisReadyQueue) Pop(ctx context.Context) (*ReadyNode, error) {
	result, err := cache.Instance().Client().BLPop(ctx, 0, q.key).Result()
	if err != nil {
		return nil, err
	}
	if len(result) < 2 {
		return nil, fmt.Errorf("empty result")
	}
	var node ReadyNode
	if err := json.Unmarshal([]byte(result[1]), &node); err != nil {
		return nil, err
	}
	return &node, nil
}

func (q *RedisReadyQueue) Len() int {
	n, _ := cache.Instance().Client().LLen(context.Background(), q.key).Result()
	return int(n)
}

func (q *RedisReadyQueue) State() *ReadyQueueState {
	results, _ := cache.Instance().Client().LRange(context.Background(), q.key, 0, -1).Result()
	nodes := make([]*ReadyNode, 0, len(results))
	for _, r := range results {
		var node ReadyNode
		if json.Unmarshal([]byte(r), &node) == nil {
			nodes = append(nodes, &node)
		}
	}
	return &ReadyQueueState{PendingNodes: nodes}
}

func (q *RedisReadyQueue) RestoreState(state *ReadyQueueState) error {
	ctx := context.Background()
	cache.Instance().Client().Del(ctx, q.key)
	for _, node := range state.PendingNodes {
		if err := q.Push(node); err != nil {
			return err
		}
	}
	return nil
}
