package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/eoof/go-syslog"
	"github.com/eoof/q"
)

func main() {
	udpAddr := flag.String("udp", "0.0.0.0:514", "udp addr")
	tcpAddr := flag.String("tcp", "0.0.0.0:514", "tcp addr")
	flag.Parse()
	logChan := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(logChan)

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

	go func(logChan syslog.LogPartsChannel) {
		for logParts := range logChan {
			q.Q(logParts)
		}
	}(logChan)

	fmt.Println("waiting")
	server.Wait()
	fmt.Println("exiting")
}
