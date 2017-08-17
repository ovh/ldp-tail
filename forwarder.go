package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"strings"
)

func forwardToLDP(messageChan <-chan map[string]interface{}, ldpURL string, forwardToken string) {
	conn, err := tls.Dial("tcp", ldpURL, &tls.Config{})
	if err != nil {
		panic("failed to Connect")
	}
	defer conn.Close()

	for {
		message := <-messageChan
		//replace X-OVH-TOKEN
		for k, _ := range message {
			if strings.HasPrefix(k, "_X-OVH-TOKEN") {
				delete(message, k)
			}
		}
		message["_X-OVH-TOKEN"] = forwardToken
		data, err := json.Marshal(&message)
		data = append(data, byte(0))
		sent, err := conn.Write(data)
		if sent < len(data) || err != nil {
			log.Fatal(err)
		}
	}
}
