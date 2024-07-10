package log

import (
	"os"

	pluginZap "github.com/micro/plugins/v5/logger/zap"
	"go-micro.dev/v5/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	rotator "gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultPath       = "/var/log/cube-api.log"
	defaultLevel      = 2
	defaultMaxSize    = 10
	defaultMaxBackups = 3
	defaultMaxAge     = 28
	defaultCompress   = true
)

var (
	Opts *Options
)

type Option func(*Options)

func newMultiWriteSyncer(rotator zapcore.WriteSyncer) zapcore.WriteSyncer {
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(rotator),
		zapcore.AddSync(os.Stderr),
	)
}

func newEncoder() zapcore.Encoder {
	conf := zap.NewProductionEncoderConfig()
	conf.ConsoleSeparator = "  "
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(conf)
}

func newLogRotator() zapcore.WriteSyncer {
	return zapcore.AddSync(
		&rotator.Logger{
			Filename:   Opts.File,
			MaxSize:    Opts.Rotation.Size,
			MaxBackups: Opts.Rotation.Backups,
			MaxAge:     Opts.Rotation.TTL,
			Compress:   Opts.Rotation.Compress,
		},
	)
}

func newLogger(rotator zapcore.WriteSyncer) (logger.Logger, error) {
	logger := zap.New(zapcore.NewCore(
		newEncoder(),
		newMultiWriteSyncer(rotator),
		zapcore.InfoLevel,
	))

	return pluginZap.NewLogger(
		pluginZap.WithLogger(logger),
	)
}

func initOptions(opts []Option) {
	Opts = &Options{
		File:  defaultPath,
		Level: defaultLevel,
		Rotation: Rotation{
			Backups:  defaultMaxBackups,
			Size:     defaultMaxSize,
			TTL:      defaultMaxAge,
			Compress: defaultCompress,
		},
	}

	for _, o := range opts {
		o(Opts)
	}
}

func NewCentralLogger(opts ...Option) (logger.Logger, error) {
	initOptions(opts)
	return newLogger(newLogRotator())
}
