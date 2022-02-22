package misc

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func initLogger() {
	hook := &lumberjack.Logger{
		Filename: "./log/zu_logic",
		MaxSize:  128,
		Compress: true,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(hook),
		zap.DebugLevel)

	Logger = zap.New(core, zap.AddCaller())
}
