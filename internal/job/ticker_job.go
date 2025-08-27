package job

import (
	"time"

	"github.com/hanxiang-s/delay_queue/internal/logger"
	"github.com/hanxiang-s/delay_queue/internal/redis"
	"github.com/hanxiang-s/delay_queue/pkg"
)

type TickerJob struct {
	Logger   logger.Logger
	RedisCli *redis.Client
	Action   pkg.JobBaseAction
	Interval int
}

func (j *TickerJob) Run() {
	ticker := time.NewTicker(time.Second * time.Duration(j.Interval))
	for range ticker.C {
		key := j.RedisCli.FormatKey(j.Action.ID())
		batch, lastScore, err := j.RedisCli.GetBatch(key)
		if err != nil {
			j.Logger.Errorf("get batch failed: %v", err)
			return
		}
		isClear := len(batch) != 0
		for _, v := range batch {
			if err = j.Action.Execute(v.Member); err != nil {
				j.Logger.Errorf("job execute failed: %v", err)
			}
		}
		if isClear {
			if err = j.RedisCli.ClearBatch(key, lastScore); err != nil {
				j.Logger.Errorf("job clear batch failed: %v", err)
			}
		}
	}
}
