// logger/logger_test.go
package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZapLogger(t *testing.T) {

	t.Run("Test File Rotation", func(t *testing.T) {
		tempLogPath := t.TempDir()

		config := LogConfig{
			Path:    tempLogPath,
			Module:  "test",
			Level:   "debug",
			MaxSize: 1, // 1 MB for testing
			MaxAge:  1,
			Backups: 3,
			Console: false,
			Format:  "console",
		}

		logger := NewZapLogger(config)
		sugar := logger.Sugar()
		logDir := filepath.Join(config.Path, config.Module)

		// 生成大量日志以触发文件轮转
		payload := strings.Repeat("日志内容 ", 256)
		for i := 0; i < 5000; i++ {
			sugar.Infof("自定义配置日志 %d %s", i, payload)
			sugar.Debugf("这是一条调试日志 %d %s", i, payload)
			sugar.Warnf("这是一条警告日志 %d %s", i, payload)
			sugar.Errorf("这是一条错误日志 %d %s", i, payload)
		}

		// 同步确保所有日志写入
		require.NoError(t, logger.Sync())

		// 验证文件是否创建
		files, err := os.ReadDir(logDir)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(files), 2, "应该存在至少一个轮转日志")

		hasCompressed := false
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".gz") {
				hasCompressed = true
				break
			}
		}
		assert.True(t, hasCompressed, "轮转日志应该被压缩为.gz文件")
	})
}
