package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	// 其他导入语句
)

func InitLogger(path string) *zap.Logger {
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    100,  // 每个日志文件最大大小（MB）
		MaxBackups: 10,   // 最多保留旧文件个数
		MaxAge:     30,   // 保留旧文件的最大天数
		Compress:   true, // 是否压缩旧文件
	})
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	logger.Sync()

	logger.Info("日志初始化成功")
	return logger
}
