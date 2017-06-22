package main

import "golang.org/x/net/websocket"

var pingCodec = websocket.Codec{Marshal: ping, Unmarshal: nil}

func ping(v interface{}) (msg []byte, payloadType byte, err error) {
	return nil, websocket.PingFrame, nil
}
