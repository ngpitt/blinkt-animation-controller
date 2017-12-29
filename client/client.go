package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"

	"github.com/ngpitt/blinkt"
)

func main() {
	server := flag.String("server", "", "server")
	flag.Parse()
	id := os.Getenv("NODE_NAME")
	log.Printf("Connecting to \"%s\"...\n", *server)
	conn, err := net.Dial("tcp", *server)
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	err = encoder.Encode(id)
	if err != nil {
		log.Panicln(err)
	}
	brightness := 0.0
	err = decoder.Decode(&brightness)
	if err != nil {
		log.Panicln(err)
	}
	strip := blinkt.NewBlinkt(blinkt.Blue, brightness)
	defer strip.Cleanup(blinkt.Blue, brightness)
	leds := make(map[int]string)
	for {
		err = decoder.Decode(&leds)
		if err != nil {
			log.Panicln(err)
		}
		for led, color := range leds {
			strip.Set(led, color, brightness)
		}
		strip.Show()
	}
}
