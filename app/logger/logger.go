package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"
	
	commonvariable "github.com/lyonmu/quebec/app/common/variable"
	
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type LogConfig struct {
	Level    string `json:"level" yaml:"level" mapstructure:"level"`             // 日志级别
	SaveFile bool   `json:"save_file" yaml:"save_file" mapstructure:"save_file"` // 保存到文件
	Path     string `json:"path" yaml:"path" mapstructure:"path"`                // 日志路径，到文件名，Rotation 文件保存在同目录
	MaxSize  int64  `json:"max_size" yaml:"max_size" mapstructure:"max_size"`    // 文件最大大小单位为MB
	MaxDay   int    `json:"max_day" yaml:"max_day" mapstructure:"max_day"`       // 日志保留天数
	writer   io.Writer
}

// GetWriter 获取 io.writer
func (c *LogConfig) GetWriter() io.Writer {
	if c.writer == nil {
		c.writer = c.getWriter()
	}
	return c.writer
}

func (c *LogConfig) getWriter() io.Writer {
	if !c.SaveFile {
		return os.Stdout
	}
	logf, err := rotatelogs.New(
		filepath.Join(c.Path+".%Y%m%d"),
		rotatelogs.WithLinkName(c.Path),
		rotatelogs.WithMaxAge(time.Duration(c.MaxDay)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		rotatelogs.WithRotationSize(c.MaxSize*commonvariable.MiB),
	)
	if err != nil {
		return os.Stdout
	}
	return logf
}

// SetLogger 设置日志
func SetLogger(c *LogConfig) {
	if lvl, err := logrus.ParseLevel(c.Level); err == nil {
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	
	logrus.SetReportCaller(true)
	logrus.SetOutput(c.GetWriter())
	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: time.DateTime,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	})
}
