package main

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestFormat(t *testing.T) {
	const (
		LOG_EMERG   = 0
		LOG_ALERT   = 1
		LOG_CRIT    = 2
		LOG_ERR     = 3
		LOG_WARNING = 4
		LOG_NOTICE  = 5
		LOG_INFO    = 6
		LOG_DEBUG   = 7
	)

	expected := strings.ToUpper(logrus.PanicLevel.String())[:4]
	got := syslogLevelToLogrusString(LOG_EMERG)
	if got != expected {
		t.Errorf("LOG_EMERG: expected `%s', got `%s'", expected, got)
	}
	expected = strings.ToUpper(logrus.PanicLevel.String())[:4]
	got = syslogLevelToLogrusString(LOG_ALERT)
	if got != expected {
		t.Errorf("LOG_ALERT: expected `%s', got `%s'", expected, got)
	}
	expected = strings.ToUpper(logrus.FatalLevel.String())[:4]
	got = syslogLevelToLogrusString(LOG_CRIT)
	if got != expected {
		t.Errorf("LOG_CRIT: expected `%s', got `%s'", expected, got)
	}
	expected = strings.ToUpper(logrus.ErrorLevel.String())[:4]
	got = syslogLevelToLogrusString(LOG_ERR)
	if got != expected {
		t.Errorf("LOG_ERR: expected `%s', got `%s'", expected, got)
	}
	expected = strings.ToUpper(logrus.WarnLevel.String())[:4]
	got = syslogLevelToLogrusString(LOG_WARNING)
	if got != expected {
		t.Errorf("LOG_WARNING: expected `%s', got `%s'", expected, got)
	}
	expected = strings.ToUpper(logrus.InfoLevel.String())[:4]
	got = syslogLevelToLogrusString(LOG_NOTICE)
	if got != expected {
		t.Errorf("LOG_NOTICE: expected `%s', got `%s'", expected, got)
	}
	expected = strings.ToUpper(logrus.InfoLevel.String())[:4]
	got = syslogLevelToLogrusString(LOG_INFO)
	if got != expected {
		t.Errorf("LOG_INFO: expected `%s', got `%s'", expected, got)
	}
	expected = strings.ToUpper(logrus.DebugLevel.String())[:4]
	got = syslogLevelToLogrusString(LOG_DEBUG)
	if got != expected {
		t.Errorf("LOG_DEBUG: expected `%s', got `%s'", expected, got)
	}

	expected = strings.ToUpper(logrus.PanicLevel.String())[:4]
	got = syslogLevelToLogrusString(LOG_EMERG - 1)
	if got != expected {
		t.Errorf("<LOG_EMERG: expected `%s', got `%s'", expected, got)
	}
	expected = strings.ToUpper(logrus.DebugLevel.String())[:4]
	got = syslogLevelToLogrusString(LOG_DEBUG + 1)
	if got != expected {
		t.Errorf(">LOG_DEBUG: expected `%s', got `%s'", expected, got)
	}
}
