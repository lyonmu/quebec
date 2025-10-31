package config


import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogConfigCreatesZapInstance(t *testing.T) {
	conf := LogConfig{
		LogPath:   "./test.log",
		LogLevel:  "info",
		LogFormat: "json",
		LogMax:    24,
	}

	var level zapcore.Level
	err := level.UnmarshalText([]byte(conf.LogLevel))
	if err != nil {
		t.Fatalf("unable to parse log level: %v", err)
	}

	var encoderCfg zapcore.EncoderConfig
	var encoder zapcore.Encoder
	if conf.LogFormat == "json" {
		encoderCfg = zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	logFile, err := os.Create(conf.LogPath)
	if err != nil {
		t.Fatalf("unable to create log file: %v", err)
	}
	defer func() {
		logFile.Close()
		// _ = os.Remove(conf.LogPath)
	}()

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(logFile),
		level,
	)

	logger := zap.New(core)
	defer logger.Sync()

	logger.Info("logconfig zap instance test")
	fileInfo, err := os.Stat(conf.LogPath)
	if err != nil {
		t.Fatalf("could not stat log file: %v", err)
	}
	if fileInfo.Size() == 0 {
		t.Fatalf("log file exists but is empty, zap instance might have failed")
	}
}
