package zaplog

import (
	"github.com/spf13/viper"
	"github.com/wike2019/wike_go/pkg/service/memorylog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var once sync.Once
var logger *zap.Logger

func GetLogger(*viper.Viper) *zap.Logger {
	once.Do(func() {
		encoderConfig := zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder, // ShortCallerEncoder will print function name and line number: pkg/file.go:line
		}
		Level := zap.InfoLevel
		if viper.GetBool("development") {
			Level = zap.DebugLevel
		}
		// 创建一个写入文件的 zapcore.Core
		fileSyncWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   viper.GetString("logPath"),
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
			LocalTime:  true,
		})
		fileCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			fileSyncWriter,
			Level,
		)

		// 创建一个写入控制台的 zapcore.Core
		consoleSyncWriter := zapcore.AddSync(os.Stdout)
		consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
		consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			consoleSyncWriter,
			Level,
		)
		memorylog.LogInfo = &memorylog.LogStruct{
			Log: make(memorylog.LogList, 0, 2000),
		}
		tempSyncWriter := zapcore.AddSync(memorylog.LogInfo)
		tempSyncCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			tempSyncWriter,
			Level,
		)

		// 使用 zapcore.NewTee 创建一个包含两个 zapcore.Core 的 zapcore.Core
		core := zapcore.NewTee(fileCore, consoleCore, tempSyncCore)

		// 创建 zap 日志记录器
		logger = zap.New(core, zap.AddCaller())
	})
	return logger
}
