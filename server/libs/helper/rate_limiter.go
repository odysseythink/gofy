package helper

import (
	"fmt"
	"strconv"
	"time"

	"mlib.com/gofy/server/cache"
)

type RateLimiter struct {
	prefix       string
	max_attempts int
	time_window  int64
}

func NewRateLimiter(prefix string, max_attempts int, time_window int64) *RateLimiter {
	return &RateLimiter{
		prefix:       prefix,
		max_attempts: max_attempts,
		time_window:  time_window,
	}
}

func (rl *RateLimiter) get_key(email string) string {
	return fmt.Sprintf("%s:%s", rl.prefix, email)
}

func (rl *RateLimiter) IsRateLimited(email string) bool {
	key := rl.get_key(email)
	current_time := time.Now().Unix()
	window_start_time := current_time - rl.time_window

	cache.Instance().ZRemRangeByScore(key, "-inf", strconv.Itoa(int(window_start_time)))
	attempts := cache.Instance().ZCard(key)

	return int(attempts) >= rl.max_attempts
}
func (rl *RateLimiter) IncrementRateLimit(email string) {
	key := rl.get_key(email)
	current_time := time.Now().Unix()

	cache.Instance().ZAdd(key, current_time, float64(current_time))
	cache.Instance().ExpireKey(key, int(rl.time_window*2))
}
