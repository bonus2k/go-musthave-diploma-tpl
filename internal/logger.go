package internal

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var Log = zap.NewNop()
var Logf *zap.SugaredLogger

var customTimeEncoder = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("02.01.2006 15:04:05.000"))
}

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     customTimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func InitLogger(level string) error {
	logLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("can't parse log level; %w", err)
	}
	config := zap.Config{
		Level:            logLevel,
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	config.EncoderConfig = encoderConfig
	Log, err = config.Build()
	if err != nil {
		return fmt.Errorf("can't create config of logger; %w", err)
	}
	Logf = Log.Sugar()
	return nil
}
