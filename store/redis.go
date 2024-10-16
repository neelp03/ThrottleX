// store/redis.go

package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisStore is a Redis-based implementation of the Store interface.
// It allows rate limiters to use Redis for storing rate-limiting data.
// This implementation supports distributed rate limiting across multiple instances.
type RedisStore struct {
	client *redis.Client   // Redis client for connecting to the Redis server
	ctx    context.Context // Context for Redis operations
}

// NewRedisStore creates a new RedisStore with the given Redis client.
//
// Parameters:
//   - client: A Redis client instance (*redis.Client) configured with appropriate options
//
// Returns:
//   - A pointer to a RedisStore instance
func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
		ctx:    context.Background(),
	}
}

// Increment increments the counter for the given key by 1 in Redis.
// If the key does not exist, it initializes it with a count of 1 and sets the expiration.
//
// Parameters:
//   - key: The key to increment
//   - expiration: The duration after which the key should expire
//
// Returns:
//   - count: The new count after incrementing
//   - err: An error if the operation fails
func (r *RedisStore) Increment(key string, expiration time.Duration) (int64, error) {
	// Use a Lua script to ensure atomicity of increment and expiration setting
	script := redis.NewScript(`
        local count = redis.call('INCR', KEYS[1])
        if tonumber(count) == 1 then
            redis.call('EXPIRE', KEYS[1], ARGV[1])
        end
        return count
    `)

	result, err := script.Run(r.ctx, r.client, []string{key}, int64(expiration.Seconds())).Result()
	if err != nil {
		return 0, err
	}
	count, ok := result.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected result type: %T", result)
	}
	return count, nil
}

// AddTimestamp adds a timestamp to a sorted set associated with the key.
// It also removes any timestamps that are outside the window.
//
// Parameters:
//   - key: The key associated with the sorted set
//   - timestamp: The timestamp to add (in nanoseconds)
//   - expiration: The duration after which the key should expire
//
// Returns:
//   - err: An error if the operation fails
func (r *RedisStore) AddTimestamp(key string, timestamp int64, expiration time.Duration) error {
	// Use ZADD to add the timestamp to the sorted set
	err := r.client.ZAdd(r.ctx, key, &redis.Z{
		Score:  float64(timestamp),
		Member: timestamp,
	}).Err()
	if err != nil {
		return err
	}

	// Set the expiration on the key
	err = r.client.Expire(r.ctx, key, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// CountTimestamps counts the number of timestamps within the given range [start, end].
// It also removes any timestamps outside the window to keep the sorted set clean.
//
// Parameters:
//   - key: The key associated with the sorted set
//   - start: The start of the range (inclusive, in nanoseconds)
//   - end: The end of the range (inclusive, in nanoseconds)
//
// Returns:
//   - count: The number of timestamps within the range
//   - err: An error if the operation fails
func (r *RedisStore) CountTimestamps(key string, start int64, end int64) (int64, error) {
	// Remove timestamps that are older than the start time
	err := r.client.ZRemRangeByScore(r.ctx, key, "0", fmt.Sprintf("(%d", start)).Err()
	if err != nil {
		return 0, err
	}

	// Count the number of timestamps within the score range
	count, err := r.client.ZCount(r.ctx, key, fmt.Sprintf("%d", start), fmt.Sprintf("%d", end)).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTokenBucket retrieves the current state of the token bucket.
//
// Parameters:
//   - key: The key associated with the token bucket state
//
// Returns:
//   - state: The TokenBucketState if it exists, or nil if it does not
//   - err: An error if the operation fails
func (r *RedisStore) GetTokenBucket(key string) (*TokenBucketState, error) {
	// Use HGETALL to get all fields in the hash
	result, err := r.client.HGetAll(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		// Key does not exist
		return nil, nil
	}

	tokensStr, ok := result["tokens"]
	if !ok {
		return nil, fmt.Errorf("tokens field missing in Redis hash")
	}
	lastUpdateStr, ok := result["last_update"]
	if !ok {
		return nil, fmt.Errorf("last_update field missing in Redis hash")
	}

	tokens, err := strconv.ParseFloat(tokensStr, 64)
	if err != nil {
		return nil, err
	}
	lastUpdateTime, err := strconv.ParseInt(lastUpdateStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &TokenBucketState{
		Tokens:         tokens,
		LastUpdateTime: lastUpdateTime,
	}, nil
}

// SetTokenBucket updates the state of the token bucket.
//
// Parameters:
//   - key: The key associated with the token bucket state
//   - state: The TokenBucketState to set
//   - expiration: The duration after which the key should expire
//
// Returns:
//   - err: An error if the operation fails
func (r *RedisStore) SetTokenBucket(key string, state *TokenBucketState, expiration time.Duration) error {
	// Use HMSET to set multiple fields in the hash
	err := r.client.HMSet(r.ctx, key, map[string]interface{}{
		"tokens":      state.Tokens,
		"last_update": state.LastUpdateTime,
	}).Err()
	if err != nil {
		return err
	}

	// Set the expiration on the key
	err = r.client.Expire(r.ctx, key, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetLeakyBucket retrieves the current state of the leaky bucket.
//
// Parameters:
//   - key: The key associated with the leaky bucket state
//
// Returns:
//   - state: The LeakyBucketState if it exists, or nil if it does not
//   - err: An error if the operation fails
func (r *RedisStore) GetLeakyBucket(key string) (*LeakyBucketState, error) {
	// Use HGETALL to get all fields in the hash
	result, err := r.client.HGetAll(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		// Key does not exist
		return nil, nil
	}

	queueStr, ok := result["queue"]
	if !ok {
		return nil, fmt.Errorf("queue field missing in Redis hash")
	}
	lastLeakTimeStr, ok := result["last_leak_time"]
	if !ok {
		return nil, fmt.Errorf("last_leak_time field missing in Redis hash")
	}

	queue, err := strconv.Atoi(queueStr)
	if err != nil {
		return nil, err
	}
	lastLeakTimeInt, err := strconv.ParseInt(lastLeakTimeStr, 10, 64)
	if err != nil {
		return nil, err
	}
	lastLeakTime := time.Unix(0, lastLeakTimeInt)

	return &LeakyBucketState{
		Queue:        queue,
		LastLeakTime: lastLeakTime,
	}, nil
}

// SetLeakyBucket updates the state of the leaky bucket.
//
// Parameters:
//   - key: The key associated with the leaky bucket state
//   - state: The LeakyBucketState to set
//   - expiration: The duration after which the key should expire
//
// Returns:
//   - err: An error if the operation fails
func (r *RedisStore) SetLeakyBucket(key string, state *LeakyBucketState, expiration time.Duration) error {
	// Use HMSET to set multiple fields in the hash
	err := r.client.HMSet(r.ctx, key, map[string]interface{}{
		"queue":          state.Queue,
		"last_leak_time": state.LastLeakTime.UnixNano(),
	}).Err()
	if err != nil {
		return err
	}

	// Set the expiration on the key
	err = r.client.Expire(r.ctx, key, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}
