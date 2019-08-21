package logger

import "github.com/sirupsen/logrus"

// logger options
type Options struct {
	Severity string
}

// logger initialization
func Init(options Options, debug bool) {
	var (
		lvl logrus.Level
		fmt logrus.Formatter
		err error
	)

	if lvl, err = logrus.ParseLevel(options.Severity); err != nil {
		lvl = logrus.DebugLevel
	}

	if debug {
		fmt = &logrus.TextFormatter{
			ForceColors:      true,
			FullTimestamp:    true,
			QuoteEmptyFields: true,
		}
	} else {
		fmt = &logrus.JSONFormatter{}
	}

	logrus.SetLevel(lvl)
	logrus.SetFormatter(fmt)
}
