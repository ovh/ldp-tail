package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	flag "github.com/spf13/pflag"
)

type conf struct {
	Address string
	Match   []matchCriterion
	Pattern string
	Raw     bool
}

var defaultConf = conf{
	"",
	nil,
	"{{._appID}}> {{.short_message}}",
	false,
}

var operatorRegexp = regexp.MustCompile(`(.+?)\.(not\.)?(` + strings.Join(supportedMatchOperators, "|") + `)`)

func getConf() conf {

	configFile := flag.String("config", "", "Configuration file")

	address := flag.String("address", defaultConf.Address, "URI of the websocket")
	match := flag.StringArray("match", nil, "Fields to match")
	pattern := flag.String("pattern", defaultConf.Pattern, "Template to apply on each message to display it")
	raw := flag.Bool("raw", defaultConf.Raw, "Display raw message instead of parsing it")

	flag.Parse()

	c := defaultConf

	// Load Override default config with file
	if *configFile != "" {
		_, err := toml.DecodeFile(*configFile, &c)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Override configuration with flags
	if *address != defaultConf.Address {
		c.Address = *address
	}
	if *pattern != defaultConf.Pattern {
		c.Pattern = *pattern
	}
	if *raw != defaultConf.Raw {
		c.Raw = *raw
	}

	// Match Criteria
	for _, m := range *match {
		v := strings.SplitN(m, "=", 2)
		var key, operator, value string
		var not bool

		// Check if key match an operator
		if subMatch := operatorRegexp.FindStringSubmatch(v[0]); subMatch != nil {
			key = subMatch[1]
			not = subMatch[2] != ""
			operator = subMatch[3]
		} else {
			// Default operator
			key = v[0]
			operator = "eq"
		}

		if operator != "present" && operator != "missing" {
			if len(v) != 2 {
				log.Fatal(fmt.Errorf("Match should be in the form 'key(.operator)?=value', found %s", v))
			} else {
				value = v[1]
			}
		}

		c.Match = append(c.Match, matchCriterion{Key: key, Operator: operator, Not: not, Value: value})
	}
	if ok, err := isValidMatchCriteria(c.Match); !ok {
		log.Fatal(err)
	}

	return c
}
