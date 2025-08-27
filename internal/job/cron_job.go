package job

import (
	"github.com/hanxiang-s/delay_queue/internal/logger"
	"github.com/hanxiang-s/delay_queue/internal/redis"
	"github.com/hanxiang-s/delay_queue/pkg"
)

type CronJob struct {
	Logger   logger.Logger
	RedisCli *redis.Client
	Action   pkg.JobBaseAction
}

// Run 定时执行任务
// cron每次执行时都会开启一个协程，cron执行内部出现执行耗时很长的情况（GetBatch、Execute）会导致开启大量协程（内存消耗），因此支持了ticker方式
// 同时去掉了defer ClearBatch，每Execute一次就ZRem一个member，避免cron间隔很短（每秒执行1次）每次执行时取到已经Execute的member, 还未defer clear
func (j *CronJob) Run() {
	key := j.RedisCli.FormatKey(j.Action.ID())
	batch, _, err := j.RedisCli.GetBatch(key)
	if err != nil {
		j.Logger.Errorf("get batch failed: %v", err)
		return
	}
	for _, v := range batch {
		if err = j.Action.Execute(v.Member); err == nil {
			if err = j.RedisCli.ZRem(key, v.Member); err != nil {
				j.Logger.Errorf("job zrem failed: %v", err)
			}
		} else {
			j.Logger.Errorf("job execute failed: %v", err)
		}
	}
	//defer func() {
	//	if len(batch) != 0 {
	//		j.RedisCli.ClearBatch(key, lastScore)
	//	}
	//}()
}
