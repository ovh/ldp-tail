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

// From https://github.com/gemnasium/logrus-graylog-hook/blob/master/graylog_hook.go#L166
func logrusLevelToString(v int) string {
	return logrus.Level(v - 2).String()
}

func logrusFormater(v map[string]interface{}) string {
	level := logrusLevelToString(int(v["level"].(float64)))
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
	fmt.Fprintf(b, "%s[%s] %-44s", strings.ToUpper(level)[0:4], time.Unix(timestamp, 0).Format(time.RFC3339), v["short_message"])

	for _, k := range keys {
		fmt.Fprintf(b, ` %s="%v"`, k[1:], v[k])
	}

	return b.String()
}
