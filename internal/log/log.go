package log

import (
	kitlog "github.com/go-kit/kit/log"
	zaplog "github.com/go-kit/kit/log/zap"
	kitlevel "github.com/go-kit/log/level"
	"github.com/jacexh/goutil/zaphelper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option struct {
	Name       string
	Level      string `default:"info"`
	Filename   string
	MaxSize    int `default:"100"`
	MaxAge     int `default:"7"`
	MaxBackups int `default:"10"`
	LocalTime  bool
	Compress   bool
}

type (
	Logger interface {
		Debug(keyvals ...interface{})
		Info(keyvals ...interface{})
		Warn(keyvals ...interface{})
		Error(keyvals ...interface{})
	}

	KitZapLogger struct {
		Zap *zap.Logger
		Kit kitlog.Logger
	}
)

var (
	levelMapper = map[string]zapcore.Level{
		"info":  zapcore.InfoLevel,
		"debug": zapcore.DebugLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}

	globalLogger Logger
)

// BuildLogger 构建全局日志
func BuildLogger(opt Option) Logger {
	conf := zap.NewProductionConfig()
	conf.Sampling = nil
	conf.EncoderConfig.TimeKey = "@timestamp"
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.Level = zap.NewAtomicLevelAt(levelMapper[opt.Level])
	zl := zaphelper.BuildRotateLogger(conf, zaphelper.RotatingFileConfig{
		LoggerName: opt.Name,
		Filename:   opt.Filename,
		MaxSize:    opt.MaxSize,
		MaxAge:     opt.MaxAge,
		MaxBackups: opt.MaxBackups,
		LocalTime:  opt.LocalTime,
		Compress:   opt.Compress,
	})

	sl := zl.Named("system")
	sl = sl.WithOptions(zap.AddCallerSkip(2))
	kl := zaplog.NewZapSugarLogger(sl, conf.Level.Level())
	globalLogger = &KitZapLogger{Zap: zl, Kit: kl}
	return globalLogger
}

func (kz *KitZapLogger) log(level string, keyvals ...interface{}) {
	switch level {
	case "debug":
		_ = kitlevel.Debug(kz.Kit).Log(keyvals...)

	case "info":
		_ = kitlevel.Info(kz.Kit).Log(keyvals...)

	case "warn":
		_ = kitlevel.Warn(kz.Kit).Log(keyvals...)

	case "error":
		_ = kitlevel.Error(kz.Kit).Log(keyvals...)
	}
}

func (kz *KitZapLogger) Debug(keyvals ...interface{}) {
	kz.log("debug", keyvals...)
}

func (kz *KitZapLogger) Info(keyvals ...interface{}) {
	kz.log("info", keyvals...)
}

func (kz *KitZapLogger) Warn(keyvals ...interface{}) {
	kz.log("warn", keyvals...)
}

func (kz *KitZapLogger) Error(keyvals ...interface{}) {
	kz.log("error", keyvals...)
}

func Debug(keyvals ...interface{}) {
	globalLogger.Debug(keyvals...)
}

func Info(keyvals ...interface{}) {
	globalLogger.Info(keyvals...)
}

func Warn(keyvals ...interface{}) {
	globalLogger.Warn(keyvals...)
}

func Error(keyvals ...interface{}) {
	globalLogger.Error(keyvals...)
}
