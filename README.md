[![OSCS Status](https://www.oscs1024.com/platform/badge/yasin-wu/delay_queue.svg?size=small)](https://www.murphysec.com/dr/kFJ0vHLhJQTz8wiubq)
## 介绍

delay queue是基于Redis Zset+Cron实现的Golang版延时队列。
实现方案是任务cron定时执行时主动轮询小于当前时间的元素, 取出符合条件元素执行任务，完成任务后删除该元素。
支持延迟多少秒和延迟到具体时间执行。
说明：redis zset key和cron id都是keyPrefix:jobID
## 安装

```
go get -u github.com/hanxiang-s/delay_queue
```

推荐使用go.mod

## 使用

```go
type JobActionSMS struct{}

// ID 任务ID
func (j *JobActionSMS) ID() string {
    return "JobActionSMS"
}

// Scheduler 任务定时执行，执行时从zset中获取0<score<=当前时间的member去执行任务
func (j *JobActionSMS) Scheduler() pkg.Scheduler {
    return pkg.Scheduler{
        Type:  pkg.SchedulerTypeCron,
        Value: "@every 1s",
    }
}

/*
// Scheduler 任务定时执行，执行时从zset中获取0<score<=当前时间的member去执行任务
func (j *JobActionSMS) Scheduler() pkg.Scheduler {
    return pkg.Scheduler{
        Type:  pkg.SchedulerTypeTicker,
        Value: "1",
    }
}
 */

// Execute 任务执行方法
func (j *JobActionSMS) Execute(arg any) error {
    phone, _ := arg.(string)
    fmt.Printf("sending sms to %s,time:%v\n", phone, time.Now().Format("2006-01-02 15:04:05"))
    return nil
}

func main() {
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

```
