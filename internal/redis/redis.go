package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hanxiang-s/delay_queue/internal/logger"

	"github.com/hanxiang-s/delay_queue/pkg"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	keyPrefix  string
	batchLimit int64
	client     *redis.Client
	ctx        context.Context
	logger     logger.Logger
}

func New(keyPrefix string, batchLimit int64, opt *redis.Options) *Client {
	return &Client{
		keyPrefix:  keyPrefix,
		batchLimit: batchLimit,
		client:     redis.NewClient(opt),
		ctx:        context.Background(),
		logger:     logger.DefaultLogger,
	}
}

func (c *Client) ZAdd(job pkg.DelayJob) error {
	if job.Arg == nil {
		return errors.New("job arg is nil")
	}
	key := c.FormatKey(job.ID)
	delayTime := job.DelayTime
	job.DelayTime = -1
	var z redis.Z
	z.Member = job.Arg
	z.Score = float64(delayTime + time.Now().Unix())
	switch job.Type {
	case pkg.DelayTypeDuration:
		z.Score = float64(delayTime + time.Now().Unix())
	case pkg.DelayTypeDate:
		z.Score = float64(delayTime)
	default:
		return errors.New("job type is not supported")
	}
	return c.client.ZAdd(c.ctx, key, &z).Err()
}

func (c *Client) ZRem(key string, arg any) error {
	if arg == nil {
		return errors.New("job arg is nil")
	}
	return c.client.ZRem(c.ctx, key, arg).Err()
}

func (c *Client) GetBatch(key string) ([]redis.Z, int64, error) {
	var redisZs []redis.Z
	var lastScore int64
	var err error
	var opt redis.ZRangeBy
	opt.Min = "0"
	opt.Max = fmt.Sprintf("%d", time.Now().Unix())
	opt.Offset = 0
	opt.Count = c.batchLimit
	redisZs, err = c.client.ZRangeByScoreWithScores(c.ctx, key, &opt).Result()
	if len(redisZs) > 0 {
		lastScore = int64(redisZs[len(redisZs)-1].Score)
	}
	return redisZs, lastScore, err
}

func (c *Client) ClearBatch(key string, lastScore int64) {
	if err := c.client.ZRemRangeByScore(c.ctx, key, "0", fmt.Sprintf("%d", lastScore)).Err(); err != nil {
		c.logger.Errorf("clear batch failed: %v", err)
	}
}

func (c *Client) FormatKey(jobID string) string {
	return fmt.Sprintf("%s:%s", c.keyPrefix, jobID)
}

func (c *Client) SetLogger(logger logger.Logger) {
	c.logger = logger
}
