package dq

import (
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"

	"github.com/hanxiang-s/delay_queue/internal/cron"
	"github.com/hanxiang-s/delay_queue/internal/job"
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
	var err error
	switch {
	case action.Scheduler().Type == pkg.SchedulerTypeCron:
		cronJob := &job.CronJob{
			Logger:   dq.logger,
			RedisCli: dq.redisCli,
			Action:   action,
		}
		_, err = dq.cron.Add(dq.redisCli.FormatKey(action.ID()), action.Scheduler().Value, cronJob)
	case action.Scheduler().Type == pkg.SchedulerTypeTicker:
		interval, err := strconv.Atoi(action.Scheduler().Value)
		if err != nil {
			return err
		}
		if interval == 0 {
			return errors.New("invalid job interval: 0")
		}
		tickerJob := &job.TickerJob{
			Logger:   dq.logger,
			RedisCli: dq.redisCli,
			Action:   action,
			Interval: interval,
		}
		go tickerJob.Run()
	}
	return err
}

func (dq *DelayQueue) AddJob(job pkg.DelayJob) error {
	return dq.redisCli.ZAdd(job)
}

func (dq *DelayQueue) RemoveJob(job pkg.DelayJob) error {
	return dq.redisCli.ZRem(dq.redisCli.FormatKey(job.ID), job.Arg)
}

func (dq *DelayQueue) SetLogger(logger logger.Logger) {
	if logger != nil {
		dq.logger = logger
		dq.redisCli.SetLogger(logger)
	}
}
