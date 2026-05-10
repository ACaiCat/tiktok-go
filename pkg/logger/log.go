package logger

import (
	"io"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func InitLogger() {
	writer, err := rotatelogs.New(
		"./logs/app.%Y-%m-%d.log",
		rotatelogs.WithLinkName("./logs/app.log"),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)

	if err != nil {
		hlog.Fatal(err)
	}

	logger := hertzlogrus.NewLogger()
	logger.Logger().SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	logger.SetOutput(io.MultiWriter(writer, os.Stdout))
	logger.SetLevel(hlog.LevelInfo)

	hlog.SetLogger(logger)
}
