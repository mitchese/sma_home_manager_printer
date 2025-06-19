package main

import (
	"log"
	"os"

	"sma_home_manager_printer/handler"

	"github.com/dmichael/go-multicast/multicast"
)

const (
	defaultAddress = "239.12.255.254:9522"
)

func main() {
	address := getMulticastAddress()
	log.Printf("Listening for SMA Home Manager packets on %s", address)
	multicast.Listen(address, handler.MsgHandler)
}

func getMulticastAddress() string {
	addr := os.Getenv("SMA_MULTICAST_ADDRESS")
	if addr == "" {
		return defaultAddress
	}
	return addr
}
