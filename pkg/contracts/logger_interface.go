package contracts

type LoggerInterface interface {
	Debug(msg string)
	Info(msg string)
	Warn(err error)
	Err(err error)
	Fatal(err error)
}
