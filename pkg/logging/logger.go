package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init(level, logFile string, useJSON bool) error{
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	Log.SetLevel(lvl)
	
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil{
			return err
		}
		Log.SetOutput(file)
	}else {
		Log.SetOutput(os.Stdout)
	}

	if useJSON{
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
	return nil
} 