package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapCore struct {
	conf *ZapCoreConfig
	zapcore.Core
}

type ZapCoreConfig struct {
	Level        zapcore.Level
	Encoder      zapcore.Encoder
	Director     string
	RetentionDay int
	LogInConsole bool
}

func NewZapCore(conf *ZapCoreConfig) *ZapCore {
	entity := &ZapCore{conf: conf}
	syncer := entity.WriteSyncer()
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == conf.Level
	})
	entity.Core = zapcore.NewCore(conf.Encoder, syncer, levelEnabler)
	return entity
}

func (z *ZapCore) WriteSyncer(formats ...string) zapcore.WriteSyncer {
	cutter := NewCutter(
		z.conf.Director,
		z.conf.Level.String(),
		WithExpireDay(z.conf.RetentionDay),
		WithFormats(formats...),
		WithLayout(time.DateOnly),
	)
	if z.conf.LogInConsole {
		multiSyncer := zapcore.NewMultiWriteSyncer(os.Stdout, cutter)
		return zapcore.AddSync(multiSyncer)
	}
	return zapcore.AddSync(cutter)
}
