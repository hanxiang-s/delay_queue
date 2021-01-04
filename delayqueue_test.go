package delay_queue

import (
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/yasin-wu/delay-queue/delayqueue"
)

func TestDelayQueue(t *testing.T) {
	delayCmd := &delayqueue.DelayCmd{
		Config: &delayqueue.Config{
			BucketSize: 10,
			Redis: delayqueue.RedisConfig{
				Host:     "192.168.131.135:6379",
				Password: "1qazxsw21201",
			},
		},
	}
	delayCmd.Init()
	topic := "test_topic"
	for i := 1; i < 4; i++ {
		job := delayqueue.Job{
			Topic: topic,
			Id:    fmt.Sprintf("%d", i),
			Delay: int64(10),
			TTR:   86400,
		}
		delayCmd.Push(job)
	}
	for {
		job, err := delayCmd.Pop(topic)
		if err != nil {
			t.Error(err.Error())
		}
		fmt.Println("=====time", time.Now())
		spew.Dump(job)
		delayCmd.Delete(job.Id)
	}
}
