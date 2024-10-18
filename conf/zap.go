package conf

import (
	"go.uber.org/zap/zapcore"
	"time"
)

type Zap struct {
	Use           bool   `mapstructure:"use" json:"use" yaml:"use"`
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                            // 日志级别
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                         // 输出方式
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`                  // 存储目录
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级别
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 栈名
	ShowLine      bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 输出控制台
	RetentionDay  int    `mapstructure:"retention-day" json:"retention-day" yaml:"retention-day"`    // 日志保留天数
}

func (z *Zap) Levels() []zapcore.Level {
	levels := make([]zapcore.Level, 0, 7)
	level, err := zapcore.ParseLevel(z.Level)
	if err != nil {
		level = zapcore.DebugLevel
	}
	for ; level <= zapcore.FatalLevel; level++ {
		levels = append(levels, level)
	}
	return levels
}

func (z *Zap) Encoder() zapcore.Encoder {
	config := zapcore.EncoderConfig{
		TimeKey:       "time",
		NameKey:       "name",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: z.StacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(z.Prefix + t.Format(time.DateTime))
		},
		EncodeLevel:    z.LevelEncoder(),
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
	if z.Format == "json" {
		return zapcore.NewJSONEncoder(config)
	}
	return zapcore.NewConsoleEncoder(config)
}

func (z *Zap) LevelEncoder() zapcore.LevelEncoder {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder":
		return zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder":
		return zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder":
		return zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder":
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}
