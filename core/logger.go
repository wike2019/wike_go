package core

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(cfg *viper.Viper) *zap.Logger {
	logPath := cfg.GetString("log_path")
	if logPath == "" {
		logPath = "./logs/app.log"
	}

	maxSize := cfg.GetInt("log_max_size")
	if maxSize <= 0 {
		maxSize = 100
	}

	maxBackups := cfg.GetInt("log_max_backups")
	if maxBackups <= 0 {
		maxBackups = 10
	}

	// 文件写入器，使用 lumberjack 做日志轮转
	fileWriter := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSize,    // MB
		MaxBackups: maxBackups, // 最多保留备份数
		MaxAge:     7,          // 保留7天
		Compress:   true,       // 压缩旧日志
	}

	// JSON 编码写文件
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), zapcore.InfoLevel)

	// Console 编码写 stdout
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleLevel := zapcore.InfoLevel
	if cfg.GetString("mode") == "debug" {
		consoleLevel = zapcore.DebugLevel
	}
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), consoleLevel)

	// 合并两个输出
	core := zapcore.NewTee(fileCore, consoleCore)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
