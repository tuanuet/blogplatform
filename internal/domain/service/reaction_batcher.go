package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/infrastructure/cache"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	RedisKeyDirtyBlogs = "blog:reaction:dirty"
	RedisKeyDeltasUp   = "blog:reaction:deltas:up"
	RedisKeyDeltasDown = "blog:reaction:deltas:down"
)

// ReactionBatcher handles batching of reaction count updates to the database using Redis
type ReactionBatcher struct {
	blogRepo repository.BlogRepository
	redis    *cache.RedisClient
	interval time.Duration
	stopCh   chan struct{}
}

// NewReactionBatcher creates a new reaction batcher
func NewReactionBatcher(blogRepo repository.BlogRepository, redis *cache.RedisClient, interval time.Duration) *ReactionBatcher {
	b := &ReactionBatcher{
		blogRepo: blogRepo,
		redis:    redis,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
	return b
}

// Start begins the batch processing loop
func (b *ReactionBatcher) Start() {
	go func() {
		ticker := time.NewTicker(b.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				b.flush()
			case <-b.stopCh:
				b.flush()
				return
			}
		}
	}()
}

// Stop stops the batcher
func (b *ReactionBatcher) Stop() {
	close(b.stopCh)
}

// Add queues a reaction update in Redis
func (b *ReactionBatcher) Add(blogID uuid.UUID, upDelta, downDelta int) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pipe := b.redis.Client().Pipeline()
	idStr := blogID.String()

	if upDelta != 0 {
		pipe.HIncrBy(ctx, RedisKeyDeltasUp, idStr, int64(upDelta))
	}
	if downDelta != 0 {
		pipe.HIncrBy(ctx, RedisKeyDeltasDown, idStr, int64(downDelta))
	}
	pipe.SAdd(ctx, RedisKeyDirtyBlogs, idStr)

	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("Failed to add reaction delta to Redis for blog %s: %v", idStr, err)
	}
}

var getAndDeleteScript = redis.NewScript(`
	local up = redis.call('HGET', KEYS[1], ARGV[1])
	local down = redis.call('HGET', KEYS[2], ARGV[1])
	redis.call('HDEL', KEYS[1], ARGV[1])
	redis.call('HDEL', KEYS[2], ARGV[1])
	return {up, down}
`)

// flush writes all pending updates from Redis to the database
func (b *ReactionBatcher) flush() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Get a batch of dirty blogs
	// Use SPOP to atomically get and remove from the dirty set
	blogs, err := b.redis.Client().SPopN(ctx, RedisKeyDirtyBlogs, 100).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Failed to pop dirty blogs from Redis: %v", err)
		return
	}

	if len(blogs) == 0 {
		return
	}

	for _, blogIDStr := range blogs {
		blogID, err := uuid.Parse(blogIDStr)
		if err != nil {
			continue
		}

		// 2. Atomically get and delete deltas using Lua script
		res, err := getAndDeleteScript.Run(ctx, b.redis.Client(), []string{RedisKeyDeltasUp, RedisKeyDeltasDown}, blogIDStr).Slice()
		if err != nil {
			log.Printf("Failed to get deltas from Redis for blog %s: %v", blogIDStr, err)
			continue
		}

		upStr, _ := res[0].(string)
		downStr, _ := res[1].(string)

		upDelta, _ := strconv.Atoi(upStr)
		downDelta, _ := strconv.Atoi(downStr)

		if upDelta == 0 && downDelta == 0 {
			continue
		}

		// 3. Update DB
		if err := b.blogRepo.UpdateCounts(ctx, blogID, upDelta, downDelta); err != nil {
			log.Printf("Failed to batch update counts for blog %s: %v", blogID, err)
			// On failure, we should ideally put it back in Redis, but for now simple logging is okay.
			// To be extra safe, we could SADD it back to DirtyBlogs.
			b.Add(blogID, upDelta, downDelta)
		}
	}
}
