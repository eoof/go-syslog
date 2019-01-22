package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/eoof/go-syslog"
)

func main() {
	udpAddr := flag.String("udp", "0.0.0.0:14514", "udp addr")
	tcpAddr := flag.String("tcp", "0.0.0.0:14514", "tcp addr")
	flag.Parse()
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	if err := server.ListenUDP(*udpAddr); err != nil {
		log.Fatal(err)
	}
	if err := server.ListenTCP(*tcpAddr); err != nil {
		log.Fatal(err)
	}

	if err := server.Boot(); err != nil {
		log.Fatal(err)
	}

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			fmt.Println(logParts)
		}
	}(channel)

	fmt.Println("waiting")
	server.Wait()
	fmt.Println("exiting")
}
