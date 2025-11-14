// logger/logger.go (完整版)
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	Path    string `name:"path" env:"LOG_PATH" default:"/var/log" help:"日志文件路径" mapstructure:"path" yaml:"path" json:"path"`
	Module  string `name:"module" env:"LOG_MODULE" default:"quebec" help:"模块名称" mapstructure:"module" yaml:"module" json:"module"`
	Level   string `enum:"debug,info,warn,error" name:"level" env:"LOG_LEVEL" default:"info" help:"日志级别 [debug, info, warn, error]" mapstructure:"level" yaml:"level" json:"level"`
	Size    int    `name:"size" env:"LOG_SIZE" default:"100" help:"日志文件大小[MB]" mapstructure:"size" yaml:"size" json:"size"`
	Age     int    `name:"age" env:"LOG_AGE" default:"168" help:"日志文件保留时间[小时]" mapstructure:"age" yaml:"age" json:"age"`
	Backups int    `name:"backups" env:"LOG_BACKUPS" default:"10" help:"日志文件保留数量" mapstructure:"backups" yaml:"backups" json:"backups"`
	Console bool   `name:"console" env:"LOG_CONSOLE" default:"false" help:"是否开启控制台输出" mapstructure:"console" yaml:"console" json:"console"`
	Format  string `enum:"console,json" name:"format" env:"LOG_FORMAT" default:"console" help:"日志输出类型" mapstructure:"format" yaml:"format" json:"format"`
}

// 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05"))
}

// 创建带有模块名的编码器配置
func createModuleEncoderConfig(moduleName string, baseConfig zapcore.EncoderConfig) zapcore.EncoderConfig {
	// 复制基础配置
	config := baseConfig

	// 自定义级别编码器，在级别前面添加模块名
	config.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("[%s] %s", moduleName, l.CapitalString()))
	}

	return config
}

// NewZapLogger 创建配置好的Zap logger
func NewZapLogger(config LogConfig) *zap.Logger {
	// 确保日志目录存在
	if config.Path != "" {
		if err := os.MkdirAll(filepath.Join(config.Path, config.Module), 0755); err != nil {
			panic(err)
		}
	}

	// 设置日志级别
	var level zapcore.Level
	switch strings.ToLower(config.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 基础编码器配置
	baseEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	// 文件输出核心

	logFileName := filepath.Join(config.Path, config.Module, fmt.Sprintf("%s.log", config.Module))

	var fileEncoder zapcore.Encoder
	if config.Format == "console" {
		// 为文件console格式创建带模块名的编码器配置
		fileEncoder = zapcore.NewConsoleEncoder(createModuleEncoderConfig(config.Module, baseEncoderConfig))
	} else {
		fileEncoder = zapcore.NewJSONEncoder(baseEncoderConfig)
	}

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    config.Size,
		MaxAge:     config.Age,
		MaxBackups: config.Backups,
		Compress:   true,
		LocalTime:  true,
	})

	fileCore := zapcore.NewCore(fileEncoder, fileWriter, level)
	cores = append(cores, fileCore)

	// 控制台输出核心
	if config.Console {
		var consoleEncoder zapcore.Encoder
		if config.Format == "console" {
			// 为控制台console格式创建带模块名的编码器配置
			consoleEncoder = zapcore.NewConsoleEncoder(createModuleEncoderConfig(config.Module, baseEncoderConfig))
		} else {
			consoleEncoder = zapcore.NewJSONEncoder(baseEncoderConfig)
		}

		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 如果没有启用任何输出，启用控制台输出作为默认
	if len(cores) == 0 {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(createModuleEncoderConfig(config.Module, baseEncoderConfig)),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 创建多输出核心
	core := zapcore.NewTee(cores...)

	// 创建logger
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger
}
