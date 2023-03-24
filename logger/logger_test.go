package logger_test

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/z-tasker/avail/logger"
	"testing"
)

type testConfigurer struct {
	logLevel  string
	logFormat string
}

func (tc testConfigurer) LogLevel() string {
	return tc.logLevel
}

func (tc testConfigurer) LogFormat() string {
	return tc.logFormat
}

func TestInit(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		origLogger := log.Logger

		testCases := []struct {
			level  string
			format string
		}{
			{"info", "console"},
			{"info", "json"},
			{"debug", "console"},
			{"debug", "json"},
		}

		for _, tc := range testCases {
			t.Run(tc.level+"_"+tc.format, func(t *testing.T) {
				config := testConfigurer{
					logLevel:  tc.level,
					logFormat: tc.format,
				}

				logger.Init(config)

			})
		}
		log.Logger = origLogger
	})

	t.Run("error test", func(t *testing.T) {
		config := testConfigurer{
			logLevel:  "invalid",
			logFormat: "console",
		}
		assert.PanicsWithValue(t, "Unknown LogLevel: invalid\n", func() { logger.Init(config) })

		config = testConfigurer{
			logLevel:  "info",
			logFormat: "invalid",
		}
		assert.PanicsWithValue(t, "Unknown LogFormat: invalid\n", func() { logger.Init(config) })
	})
}
