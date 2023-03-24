package logger

import "github.com/rs/zerolog"

type ZeroLogAdapter struct {
	Logger *zerolog.Logger
}

func (z *ZeroLogAdapter) Debug(msg string) {
	z.Logger.Debug().Msg(msg)
}

func (z *ZeroLogAdapter) Info(msg string) {
	z.Logger.Info().Msg(msg)
}

func (z *ZeroLogAdapter) Warn(err error, msg string) {
	z.Logger.Warn().Err(err).Msg(msg)
}

func (z *ZeroLogAdapter) Err(err error, msg string) {
	z.Logger.Error().Err(err).Msg(msg)
}

func (z *ZeroLogAdapter) Fatal(err error, msg string) {
	z.Logger.Fatal().Err(err).Msg(msg)
}
