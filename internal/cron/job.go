package cron

import (
	"github.com/hanxiang-s/delay_queue/internal/logger"
	"github.com/hanxiang-s/delay_queue/internal/redis"
	"github.com/hanxiang-s/delay_queue/pkg"
)

type Job struct {
	Logger   logger.Logger
	RedisCli *redis.Client
	Action   pkg.BaseAction
}

func (j *Job) Run() {
	key := j.RedisCli.FormatKey(j.Action.ID())
	batch, lastScore, err := j.RedisCli.GetBatch(key)
	if err != nil {
		j.Logger.Errorf("get batch failed: %v", err)
		return
	}
	for _, v := range batch {
		if err = j.Action.Execute(v.Member); err != nil {
			j.Logger.Errorf("job execute failed: %v", err)
		}
	}
	defer func() {
		if len(batch) != 0 {
			j.RedisCli.ClearBatch(key, lastScore)
		}
	}()
}
