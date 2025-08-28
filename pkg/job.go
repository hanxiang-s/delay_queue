package pkg

type JobBaseAction interface {
	ID() string               //任务ID
	Scheduler() Scheduler     //任务定时执行Scheduler, 执行时从zset中获取0<score<=当前时间的member
	Execute(member any) error //任务执行方法
}

type DelayJob struct {
	ID        string    // 任务ID, 等于JobBaseAction.ID()
	Type      DelayType // 时间类型: 0-延迟N秒执行,1-具体执行时间
	DelayTime int64     // 延迟时间: type=0时为延迟秒数,type=1时为执行秒时间戳
	Member    any       // 任务执行参数，以其作为redis zset的member
}

type Scheduler struct {
	Type  SchedulerType // 定时执行类型: 0-cron, 1-ticker
	Value string        // 定时执行值: type=0时为cron表达式,type=1时为ticker时间间隔
}
