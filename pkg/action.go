package pkg

type BaseAction interface {
	ID() string
	Cron() string
	Execute(arg any) error
}
