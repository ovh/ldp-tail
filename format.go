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

var syslogLevelToLogrus = [...]string{
	strings.ToUpper(logrus.PanicLevel.String())[:4], // LOG_EMERG   = 0
	strings.ToUpper(logrus.PanicLevel.String())[:4], // LOG_ALERT   = 1
	strings.ToUpper(logrus.FatalLevel.String())[:4], // LOG_CRIT    = 2
	strings.ToUpper(logrus.ErrorLevel.String())[:4], // LOG_ERR     = 3
	strings.ToUpper(logrus.WarnLevel.String())[:4],  // LOG_WARNING = 4
	strings.ToUpper(logrus.InfoLevel.String())[:4],  // LOG_NOTICE  = 5
	strings.ToUpper(logrus.InfoLevel.String())[:4],  // LOG_INFO    = 6
	strings.ToUpper(logrus.DebugLevel.String())[:4], // LOG_DEBUG   = 7
}

// String reverse of
// https://github.com/gemnasium/logrus-graylog-hook/blob/master/graylog_hook.go#L144
func syslogLevelToLogrusString(v int) string {
	if v < 0 {
		v = 0
	} else if v >= len(syslogLevelToLogrus) {
		v = len(syslogLevelToLogrus) - 1
	}
	return syslogLevelToLogrus[v]
}

func logrusFormater(v map[string]interface{}) string {
	level := syslogLevelToLogrusString(int(v["level"].(float64)))
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
	fmt.Fprintf(b, "%s[%s] %-44s", level, time.Unix(timestamp, 0).Format(time.RFC3339), v["short_message"])

	for _, k := range keys {
		fmt.Fprintf(b, ` %s="%v"`, k[1:], v[k])
	}

	return b.String()
}
