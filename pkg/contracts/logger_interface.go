package contracts

type LoggerInterface interface {
	Debug(msg string)
	Info(msg string)
	Warn(err error, msg string)
	Err(err error, msg string)
	Fatal(err error, msg string)
}
