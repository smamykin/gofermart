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

func (z *ZeroLogAdapter) Warn(err error) {
	z.Logger.Warn().Stack().Err(err).Msg("")
}

func (z *ZeroLogAdapter) Err(err error) {
	z.Logger.Error().Stack().Err(err).Msg("")
}

func (z *ZeroLogAdapter) Fatal(err error) {
	z.Logger.Fatal().Stack().Err(err).Msg("")
}
