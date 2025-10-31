// logger/logger_test.go
package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZapLogger(t *testing.T) {

	t.Run("Test File Rotation", func(t *testing.T) {
		config := LogConfig{
			Path:    t.TempDir(),
			Module:  "test",
			Level:   "debug",
			Size:    3, // 3 MB for testing
			Age:     7,
			Backups: 3,
			Console: true,
			Format:  "console",
		}

		logger := NewZapLogger(config)
		logDir := filepath.Join(config.Path, "test")
		defer logger.Sync()

		logger.Sugar().Info("Starting file rotation test ", logDir)

		// 生成大量日志以触发文件轮转
		for range 12 {
			logger.Sugar().Info("自定义配置日志")
			logger.Sugar().Debug("这是一条调试日志")
			logger.Sugar().Warn("这是一条警告日志")
			logger.Sugar().Error("这是一条错误日志")
		}

		// 同步确保所有日志写入
		logger.Sync()

		// 验证文件是否创建

		files, err := os.ReadDir(logDir)
		assert.NoError(t, err)
		assert.NotEmpty(t, files, "应该有日志文件被创建")

		// 删除临时文件夹
		os.RemoveAll(logDir)
	})
}
