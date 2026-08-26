package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	flag "github.com/spf13/pflag"
)

const config = ".ldp-tail.toml"

type conf struct {
	Address     string
	Match       []matchCriterion
	ExcludeKeys []string
	Raw         bool
	Format      string
	formatFunc  func(map[string]any) string
	Pattern     string
}

var defaultConf = conf{
	Pattern: "{{._appID}}> {{.short_message}}",
}

var operatorRegexp = regexp.MustCompile(`(.+?)\.(not\.)?(` + strings.Join(supportedMatchOperators, "|") + `)`)

// Build details
var (
	buildVersion = "dev"
	buildCommit  = "unknown"
	buildDate    = "unknown"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s (Version %s):\n", os.Args[0], buildVersion)
		flag.PrintDefaults()
	}
}

func getConf() conf {
	configFile := flag.String("config", "", "Configuration file")

	address := flag.String("address", defaultConf.Address, "URI of the websocket")
	match := flag.StringArray("match", nil, "Fields to match")
	raw := flag.Bool("raw", defaultConf.Raw, "Display raw message instead of parsing it")
	format := flag.String("format", defaultConf.Format, fmt.Sprintf("Display messages using a pre-defined format. Valid values: (%s)", strings.Join(supportedFormat, ", ")))
	pattern := flag.String("pattern", defaultConf.Pattern, "Template to apply on each message to display it")

	flag.Parse()

	c := defaultConf

	// Load Override default config with file
	if *configFile != "" {
		_, err := toml.DecodeFile(*configFile, &c)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *configFile == "" {
		usr, err := user.Current()
		if err == nil {
			filename := filepath.Join(usr.HomeDir, config)
			_, err = toml.DecodeFile(filename, &c)
			if err != nil {
				log.Printf("Could not load global %s: %s", filename, err)
			}
		}
	}

	// Override configuration with flags
	if *address != defaultConf.Address {
		c.Address = *address
	}
	if *raw != defaultConf.Raw {
		c.Raw = *raw
	}
	if *format != defaultConf.Format {
		c.Format = *format
	}
	if *pattern != defaultConf.Pattern {
		c.Pattern = *pattern
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

	if flag.NArg() > 0 {
		if flag.Arg(0) == "version" {
			fmt.Fprintf(os.Stderr, "ldp-tail version %s (%s - %s)\n", buildVersion, buildCommit, buildDate)
			os.Exit(0)
		} else if flag.Arg(0) == "help" {
			flag.Usage()
			os.Exit(0)
		} else {
			fmt.Printf("Invalid command %q\n", flag.Arg(0))
			flag.Usage()
			os.Exit(-1)
		}
	}

	// Check format helper
	if c.Format != "" {
		f, ok := supportedFormatMap[c.Format]
		if !ok {
			fmt.Fprintf(os.Stderr, "Invalid `format`: %q\n", c.Format)
			flag.Usage()
			os.Exit(-1)
		}
		c.formatFunc = f
	}

	if c.Address == "" {
		fmt.Fprintf(os.Stderr, "No `address` specified. Please specify it with --address or thru a config file\n")
		flag.Usage()
		os.Exit(-1)
	}

	return c
}
