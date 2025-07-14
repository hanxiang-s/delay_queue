package pkg

type JobBaseAction interface {
	ID() string            //任务ID
	Cron() string          //任务定时表达式，表示多久执行一次从redis zset中获取0<score<当前秒时间戳member
	Execute(arg any) error //任务执行方法
}
