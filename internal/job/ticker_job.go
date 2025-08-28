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
	j.Logger.Infof("delay queue ticker job start")
	ticker := time.NewTicker(time.Second * time.Duration(j.Interval))
	for range ticker.C {
		key := j.RedisCli.FormatKey(j.Action.ID())
		batch, _, err := j.RedisCli.GetBatch(key)
		if err != nil {
			j.Logger.Errorf("get batch failed: %v", err)
			continue
		}
		var members []any
		for _, v := range batch {
			if err = j.Action.Execute(v.Member); err == nil {
				members = append(members, v.Member)
			} else {
				j.Logger.Errorf("job execute failed: %v", err)
			}
		}
		if len(members) > 0 {
			if err = j.RedisCli.ZRem(key, members); err != nil {
				j.Logger.Errorf("job zrem members failed: %v", err)
			}
		}
	}
	j.Logger.Infof("delay queue ticker job stop")
}
