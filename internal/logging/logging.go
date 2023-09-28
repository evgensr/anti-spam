package logging

import (
	"github.com/evgensr/anti-spam/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var e *zap.Logger

type Logger struct {
	*zap.Logger
}

func NewLogger(cfg *config.Config) *Logger {
	Init(cfg)
	return &Logger{
		e,
	}
}

func Init(cfg *config.Config) {

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Log.Path,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
		LocalTime:  cfg.Log.LocalTime,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	e = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		zap.InfoLevel,
	),
		zap.AddCaller(),
	)

}
