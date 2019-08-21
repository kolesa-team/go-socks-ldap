package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	t.Run("ValidLevel", func(t *testing.T) {
		o := Options{
			Severity: "fatal",
		}

		Init(o, false)
		assert.Equal(t, logrus.GetLevel(), logrus.FatalLevel)
	})
	t.Run("InvalidLevel", func(t *testing.T) {
		o := Options{
			Severity: "=(",
		}

		Init(o, false)
		assert.Equal(t, logrus.GetLevel(), logrus.DebugLevel)
	})
	t.Run("DebugEnabled", func(t *testing.T) {
		o := Options{
			Severity: "=(",
		}

		Init(o, true)
		assert.Equal(t, logrus.StandardLogger().Formatter, &logrus.TextFormatter{
			ForceColors:      true,
			FullTimestamp:    true,
			QuoteEmptyFields: true,
		})
	})
}
