package logger

import (
	constants "aweme_kitex/pkg/constant"
	"os"
	"time"

	"github.com/Shopify/sarama"
	logkafka "github.com/kenjones-cisco/logrus-kafka-hook"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func DoInit(url, topic string, level logrus.Level) error {
	logger = logrus.New()
	if url != "" && topic != "" {
		// 需要写入kafka，否则只需要写本地
		producer, err := logkafka.SimpleProducer([]string{url}, sarama.CompressionSnappy, sarama.WaitForLocal, nil)
		if err != nil {
			return err
		}
		hook := logkafka.New().WithFormatter(logkafka.DefaultFormatter(logrus.Fields{"appid": "aweme.kitex"})).
			WithProducer(producer).WithLevels([]logrus.Level{logrus.DebugLevel, logrus.ErrorLevel, logrus.InfoLevel, logrus.PanicLevel}).
			WithTopic(topic)
		logger.Hooks.Add(hook)
	}

	src, err := os.OpenFile(constants.LogFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	logger.Out = src
	// 设置日志级别
	logger.SetLevel(level)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		constants.LogFileName+".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(constants.LogFileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.Hooks.Add(lfHook)
	return nil
}

func Error(args ...interface{}) {
	logger.Error(args)
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}

func Warn(args ...interface{}) {
	logger.Warn(args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Panic(args ...interface{}) {
	logger.Panic(args)
}
