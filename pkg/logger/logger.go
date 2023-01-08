package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func InitLogger() error {
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	infoWriteSyncer := getLogWriter("./utils/log/info.log")
	errorWriteSyncer := getLogWriter("./utils/log/error.log")
	encoder := getEncoder()
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoWriteSyncer, infoLevel),
		zapcore.NewCore(encoder, errorWriteSyncer, errorLevel),
	)
	log := zap.New(core, zap.AddCaller())
	logger = log.Sugar()
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderconfig := zap.NewProductionEncoderConfig()
	encoderconfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderconfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderconfig)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	lumberLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1,  // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxAge:     5,  // 保留旧文件的最大个数
		MaxBackups: 30, // 保留旧文件的最大天数
		LocalTime:  false,
		Compress:   false,
	}
	return zapcore.AddSync(lumberLogger)
}

// 外部直接访问
func Sync() {
	logger.Sync()
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
