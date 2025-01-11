package log

import (
	"os"

	pluginZap "github.com/micro/plugins/v5/logger/zap"
	"go-micro.dev/v5/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	rotator "gopkg.in/natefinch/lumberjack.v2"
)

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

func newLogRotator(opts *Options) zapcore.WriteSyncer {
	return zapcore.AddSync(
		&rotator.Logger{
			Filename:   opts.File,
			MaxSize:    opts.Rotation.Size,
			MaxBackups: opts.Rotation.Backups,
			MaxAge:     opts.Rotation.TTL,
			Compress:   opts.Rotation.Compress,
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

func initOptions(opts []Option) *Options {
	options := &Options{
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
		o(options)
	}

	return options
}

func NewGlobalHelper(opts ...Option) error {
	initedOpts := initOptions(opts)

	var err error
	logger.DefaultLogger, err = newLogger(newLogRotator(initedOpts))
	if err != nil {
		return err
	}

	return nil
}
