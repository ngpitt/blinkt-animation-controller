package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"sync"
	"time"
)

type Line struct {
	Id    string
	Leds  map[int]string
	Delay int
}

type Client struct {
	Conn    net.Conn
	Encoder *json.Encoder
	Decoder *json.Decoder
}

func main() {
	brightness := flag.Float64("brightness", 0.25, "brightness")
	listen := flag.String("listen", ":8888", "listen")
	config := flag.String("config", "[]", "config")
	flag.Parse()
	log.Printf("Starting server on \"%s\"...\n", *listen)
	lines := []Line{}
	err := json.Unmarshal([]byte(*config), &lines)
	if err != nil {
		log.Panicln(err)
	}
	clients := make(map[string]*Client)
	mutex := &sync.RWMutex{}
	listener, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Panicln(err)
	}
	defer listener.Close()
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			client := &Client{
				conn,
				json.NewEncoder(conn),
				json.NewDecoder(conn),
			}
			id := ""
			err = client.Decoder.Decode(&id)
			if err != nil {
				conn.Close()
				continue
			}
			err = client.Encoder.Encode(brightness)
			if err != nil {
				conn.Close()
				continue
			}
			mutex.Lock()
			clients[id] = client
			mutex.Unlock()
			log.Printf("Client \"%s\" connected...\n", id)
		}
	}()
	for {
		for _, line := range lines {
			mutex.RLock()
			client := clients[line.Id]
			mutex.RUnlock()
			if client != nil {
				err := client.Encoder.Encode(line.Leds)
				if err != nil {
					client.Conn.Close()
					mutex.Lock()
					delete(clients, line.Id)
					mutex.Unlock()
					log.Printf("Client \"%s\" disconnected...\n", line.Id)
				}
			}
			time.Sleep(time.Duration(line.Delay) * time.Millisecond)
		}
	}
}
