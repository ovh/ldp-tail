package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var supportedFormatMap = map[string]func(map[string]interface{}) string{
	"logrus": logrusFormater,
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

func logrusFormater(v map[string]interface{}) string {
	syslogLevel := int(v["level"].(float64))
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
	fmt.Fprintf(b, "%s[%s] %-44s", levelString, time.Unix(timestamp, 0).Format(time.RFC3339), v["short_message"])

	for _, k := range keys {
		fmt.Fprintf(b, ` %s="%v"`, k[1:], v[k])
	}

	return b.String()
}
