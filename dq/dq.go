package dq

import (
	"github.com/go-redis/redis/v8"

	"github.com/hanxiang-s/delay_queue/internal/cron"
	rd "github.com/hanxiang-s/delay_queue/internal/redis"

	"github.com/hanxiang-s/delay_queue/pkg"

	"github.com/hanxiang-s/delay_queue/internal/logger"
)

type DelayQueue struct {
	logger   logger.Logger
	redisCli *rd.Client
	cron     *cron.Cron
}

// New 创建延迟队列
// redis key=keyPrefix:jobID
func New(keyPrefix string, batchLimit int64, opt *redis.Options) *DelayQueue {
	if keyPrefix == "" {
		keyPrefix = pkg.DefaultKeyPrefix
	}
	if batchLimit == 0 {
		batchLimit = pkg.DefaultBatchLimit
	}
	delayQueue := &DelayQueue{
		logger:   logger.DefaultLogger,
		redisCli: rd.New(keyPrefix, batchLimit, opt),
		cron:     cron.New(),
	}
	return delayQueue
}

func (dq *DelayQueue) Register(action pkg.JobBaseAction) error {
	cronJob := &cron.Job{
		Logger:   dq.logger,
		RedisCli: dq.redisCli,
		Action:   action,
	}
	_, err := dq.cron.Add(dq.redisCli.FormatKey(action.ID()), action.Cron(), cronJob)
	return err
}

func (dq *DelayQueue) AddJob(job pkg.DelayJob) error {
	return dq.redisCli.ZAdd(job)
}

func (dq *DelayQueue) RemoveJob(job pkg.DelayJob) error {
	return dq.redisCli.ZRem(job)
}

func (dq *DelayQueue) SetLogger(logger logger.Logger) {
	if logger != nil {
		dq.logger = logger
		dq.redisCli.SetLogger(logger)
	}
}
