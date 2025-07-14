package pkg

type DelayJob struct {
	ID        string    // 任务ID，需保证唯一，keyPrefix+ID作为cron的ID：BaseAction.ID
	Type      DelayType // 时间类型: 0-延迟N秒执行,1-具体执行时间
	DelayTime int64     // 延迟时间: type=0时为延迟秒数,type=1时为执行秒时间戳
	Arg       any       // 任务执行参数，以其作为redis zset的member
}
