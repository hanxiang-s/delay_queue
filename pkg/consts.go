package pkg

const (
	DefaultKeyPrefix  = "DelayQueue"
	DefaultBatchLimit = 10000
)

// DelayType 延迟任务类型
type DelayType int

const (
	DelayTypeDuration DelayType = iota // 延迟多少秒执行
	DelayTypeDate                      // 具体执行时间(时间戳:秒)
)

type SchedulerType int

const (
	SchedulerTypeCron   SchedulerType = iota // cron表达式
	SchedulerTypeTicker                      // ticker时间间隔
)
