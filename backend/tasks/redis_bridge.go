package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/odysseythink/gofy/backend/cache"
	"github.com/odysseythink/mlog"
	"github.com/redis/go-redis/v9"
)

// redisCli returns the project's shared Redis UniversalClient.
func redisCli() redis.UniversalClient {
	return cache.Instance().Client()
}

// pushToRedis appends a value to a Redis list (RPUSH).
func pushToRedis(key string, value string) error {
	cli := redisCli()
	if cli == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return cli.RPush(context.Background(), key, value).Err()
}

// popFromRedis removes and returns the first element of a Redis list,
// blocking up to timeout. Returns ("", nil) on timeout.
func popFromRedis(key string, timeout time.Duration) (string, error) {
	cli := redisCli()
	if cli == nil {
		return "", fmt.Errorf("redis client not initialized")
	}
	result, err := cli.BLPop(context.Background(), timeout, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	// BLPop returns [key, value]
	if len(result) < 2 {
		return "", nil
	}
	return result[1], nil
}

// pushToRedisDelayed adds a value to a Redis sorted set with the given score (ZADD).
func pushToRedisDelayed(key string, value string, score float64) error {
	cli := redisCli()
	if cli == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return cli.ZAdd(context.Background(), key, redis.Z{
		Score:  score,
		Member: value,
	}).Err()
}

// moveDelayedToReady moves tasks whose scheduled time has arrived
// from the delayed sorted set to the ready list.
func moveDelayedToReady(delayedKey string, readyKey string) {
	cli := redisCli()
	if cli == nil {
		return
	}
	ctx := context.Background()
	now := fmt.Sprintf("%d", time.Now().Unix())

	// Fetch tasks with score <= now
	results, err := cli.ZRangeByScore(ctx, delayedKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: now,
	}).Result()
	if err != nil {
		mlog.Errorf("moveDelayedToReady: ZRangeByScore failed: %v", err)
		return
	}

	for _, val := range results {
		// Push to the ready queue
		if err := cli.RPush(ctx, readyKey, val).Err(); err != nil {
			mlog.Errorf("moveDelayedToReady: RPush failed: %v", err)
			continue
		}
		// Remove from the delayed set
		if err := cli.ZRem(ctx, delayedKey, val).Err(); err != nil {
			mlog.Errorf("moveDelayedToReady: ZRem failed: %v", err)
		}
	}
}
