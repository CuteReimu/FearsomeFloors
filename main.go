package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

func main() {
	g := newBoard(2)
	ebiten.SetWindowSize(1024, 768)
	if err := ebiten.RunGame(g); err != nil {
		logger.Fatal(err)
	}
}

var logger = &errorEntryWithStack{logrus.NewEntry(logrus.New())}

type errorEntryWithStack struct {
	*logrus.Entry
}

func (e *errorEntryWithStack) WithError(err error) *logrus.Entry {
	return e.Logger.WithError(fmt.Errorf("%+v", err))
}

func init() {
	logrus.SetReportCaller(true)

	writerError, err := rotatelogs.New(
		path.Join("logs", "error-%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logrus.WithError(err).Fatalln("unable to write logs")
	}
	logrus.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.WarnLevel:  writerError,
			logrus.ErrorLevel: writerError,
			logrus.FatalLevel: writerError,
			logrus.PanicLevel: writerError,
		}, &logrus.TextFormatter{DisableQuote: true},
	))

	writerConsole, err := rotatelogs.New(
		path.Join("logs", "console-%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logrus.WithError(err).Fatalln("unable to write logs")
	}
	logrus.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writerConsole,
			logrus.WarnLevel:  writerConsole,
			logrus.ErrorLevel: writerConsole,
			logrus.FatalLevel: writerConsole,
			logrus.PanicLevel: writerConsole,
		}, &logrus.TextFormatter{DisableQuote: true},
	))
}
