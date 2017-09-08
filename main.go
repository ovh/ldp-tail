package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"text/template"
	"time"

	"golang.org/x/net/websocket"
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {

	c := getConf()

	// Check args
	if c.Address == "" {
		log.Fatal("`address` is required")
	}

	t, err := template.New("template").Funcs(template.FuncMap{
		"color":    color,
		"bColor":   bColor,
		"noColor":  func() string { return color("reset") },
		"date":     date,
		"join":     join,
		"concat":   concat,
		"duration": duration,
		"int":      func(v string) (int64, error) { f, e := strconv.ParseFloat(v, 64); return int64(f), e },
		"float":    func(v string) (float64, error) { f, e := strconv.ParseFloat(v, 64); return f, e },
		"get":      get,
		"column":   column,
		"begin":    begin,
		"contain":  contain,
		"level":    level,
	}).Parse(c.Pattern)
	if err != nil {
		log.Fatalf("Failed to parse pattern: %s", err.Error())
	}

	var u *url.URL

	u, err = url.Parse(c.Address)
	if err != nil {
		log.Fatal(err)
	}

	// Display filters
	if len(c.Match) > 0 {
		log.Printf("Filters are:")
		for _, f := range c.Match {
			if f.Not {
				log.Printf("  "+supportedMatchOperatorsMap[f.Operator].descriptionNot, f.Key, f.Value)
			} else {
				log.Printf("  "+supportedMatchOperatorsMap[f.Operator].description, f.Key, f.Value)
			}

		}
	}
	var messageChannel chan map[string]interface{}
	forwardEnabled := false
	silent := false

	if c.ForwardURL != "" && c.ForwardToken != "" {
		messageChannel = make(chan map[string]interface{}, 100)
		go forwardToLDP(messageChannel, c.ForwardURL, c.ForwardToken)
		forwardEnabled = true
		if c.Silent {
			silent = true
		}
	}

	for {
		// Try to connect
		log.Printf("Connecting to %s...\n", u.Host)

		ws, err := websocket.Dial(u.String(), "", "http://mySelf")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Connected!")

		if silent {
			log.Println("Forwarding in silent mode")
		}

		var msg []byte
		for {
			ws.SetReadDeadline(time.Now().Add(5 * time.Second))
			err := websocket.Message.Receive(ws, &msg)

			if t, ok := err.(net.Error); ok && t.Timeout() {
				// Timeout, send a Pong && continue
				pingCodec.Send(ws, nil)
				continue
			}

			if err != nil {
				log.Printf("Error while reading from %q: %q. Will try to reconnect after 1s...\n", u.Host, err.Error())
				time.Sleep(1 * time.Second)
				break
			}

			// Extract Message
			var logMessage struct {
				Message string `json:"message"`
			}
			json.Unmarshal(msg, &logMessage)

			// Extract infos
			var message map[string]interface{}
			json.Unmarshal([]byte(logMessage.Message), &message)

			if !match(message, c.Match) {
				continue
			}

			if !silent {
				if c.Raw {
					fmt.Printf("%+v\n", message)
				} else {
					// Print them
					err = t.Execute(os.Stdout, message)
					os.Stdout.Write([]byte{'\n'})
					if err != nil {
						log.Printf("Error while executing template: %s", err.Error())
					}
				}
			}

			if forwardEnabled {
				messageChannel <- message
			}
		}
	}
}
