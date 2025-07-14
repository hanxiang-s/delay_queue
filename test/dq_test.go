package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/hanxiang-s/delay_queue/dq"
	"github.com/hanxiang-s/delay_queue/pkg"
)

type JobActionSMS struct{}

func (j *JobActionSMS) ID() string {
	return "JobActionSMS"
}

func (j *JobActionSMS) Cron() string {
	return "@every 1s"
}

func (j *JobActionSMS) Execute(arg any) error {
	phone, _ := arg.(string)
	fmt.Printf("sending sms to %s,time: %v\n", phone, time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func TestDelayQueue(t *testing.T) {
	redisOpt := &redis.Options{Addr: "127.0.0.1:6379", Password: "password"}
	cli := dq.New("test", 0, redisOpt)
	if err := cli.Register(&JobActionSMS{}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("add job: ", time.Now().Format("2006-01-02 15:04:05"))
	if err := cli.AddJob(pkg.DelayJob{
		ID:        (&JobActionSMS{}).ID(),
		Type:      pkg.DelayTypeDuration, //延迟N秒执行
		DelayTime: 10,                    //延迟秒数
		Arg:       "138****0000",
	}); err != nil {
		log.Fatal(err)
	}
	if err := cli.AddJob(pkg.DelayJob{
		ID:        (&JobActionSMS{}).ID(),
		Type:      pkg.DelayTypeDate,      //延迟到具体时间执行
		DelayTime: time.Now().Unix() + 10, //执行时间的秒时间戳
		Arg:       "138****1111",
	}); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 30)
}
