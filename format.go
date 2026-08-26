package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var supportedFormatMap = map[string]func(map[string]any) string{
	"json":         jsonFormater,
	"logrus":       logrusFormaterWithoutColor,
	"logrus-color": logrusFormaterWithColor,
}

var supportedFormat = func() []string {
	a := make([]string, 0, len(supportedFormatMap))

	for k := range supportedFormatMap {
		a = append(a, k)
	}

	return a
}()

var syslogLevelToLogrus = []logrus.Level{
	logrus.PanicLevel, // LOG_EMERG   = 0
	logrus.PanicLevel, // LOG_ALERT   = 1
	logrus.FatalLevel, // LOG_CRIT    = 2
	logrus.ErrorLevel, // LOG_ERR     = 3
	logrus.WarnLevel,  // LOG_WARNING = 4
	logrus.InfoLevel,  // LOG_NOTICE  = 5
	logrus.InfoLevel,  // LOG_INFO    = 6
	logrus.DebugLevel, // LOG_DEBUG   = 7
}

func logrusLevelToColor(l logrus.Level) int {
	switch l {
	case logrus.DebugLevel:
		return 37 // gray
	case logrus.WarnLevel:
		return 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return 31 // red
	default:
		return 36 // blue
	}
}

func jsonFormater(v map[string]any) string {
	payload, err := json.Marshal(v)
	if err != nil {
		return "could not format in JSON"
	}
	return string(payload)
}

func logrusFormaterWithColor(v map[string]any) string {
	return logrusFormater(v, true)
}

func logrusFormaterWithoutColor(v map[string]any) string {
	return logrusFormater(v, false)
}

func logrusFormater(v map[string]any, color bool) string {
	syslogLevel := 6
	if x, ok := v["level"]; ok {
		syslogLevel = int(x.(float64))
	} else if x, ok := v["_stream"]; ok {
		if x == "stderr" {
			syslogLevel = 3
		}
	}
	level := syslogLevelToLogrus[syslogLevel]

	timestamp := int64(v["timestamp"].(float64))

	// Fields
	keys := make([]string, 0, len(v))
	for k := range v {
		if k[0] == '_' && k != "_file" && k != "_line" && k != "_pid" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// Output
	b := &bytes.Buffer{}
	levelString := strings.ToUpper(level.String())[0:4]
	if color {
		levelString = fmt.Sprintf("\x1b[%dm%s\x1b[0m", logrusLevelToColor(level), levelString)
	}
	fmt.Fprintf(b, "%s[%s] %s", levelString, time.Unix(timestamp, 0).Format(time.RFC3339), v["short_message"])

	for _, k := range keys {
		fmt.Fprintf(b, ` %s="%v"`, k[1:], v[k])
	}

	return b.String()
}
