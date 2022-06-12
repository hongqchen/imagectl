package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	Logger       *zap.SugaredLogger
	levelEnabler zapcore.LevelEnabler
)

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func customLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + l.CapitalString() + "]")
}

func InitLogger(enableDebug string) {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = customLevelEncoder

	if enableDebug == "true" {
		levelEnabler = zapcore.DebugLevel
	} else {
		levelEnabler = zapcore.InfoLevel
	}

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), os.Stdout, levelEnabler)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.FatalLevel)).Sugar()
}
