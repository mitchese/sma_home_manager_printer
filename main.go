package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"log"
	"net"
	"os"

	"github.com/dmichael/go-multicast/multicast"
	"github.com/klauspost/compress/gzip"
)

const (
	defaultAddress = "239.12.255.254:9522"
)

type singlePhase struct {
	voltage float32 // Volts
	a       float32 // Amps
	power   float32 // Watts
	forward float64 // kWh, purchased power
	reverse float64 // kWh, sold power
}

func main() {
	address := getMulticastAddress()
	log.Printf("Listening for SMA Home Manager packets on %s", address)
	multicast.Listen(address, msgHandler)
}

func getMulticastAddress() string {
	addr := os.Getenv("SMA_MULTICAST_ADDRESS")
	if addr == "" {
		return defaultAddress
	}
	return addr
}

func compressAndEncode(b []byte) string {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(b)
	if err != nil {
		return "compression error"
	}
	gz.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println("-----------------------------------------------------")
	if len(b) < 200 { // if it is too short, it is probably not a meter update
		log.Printf("Received packet too small from %v: size %d", src, n)
		return
	}

	protocolID := binary.BigEndian.Uint16(b[16:18])
	if protocolID != 24681 {
		log.Printf("Protocol ID mismatch (got 0x%04x), not a meter update", protocolID)
		return
	}

	serial := binary.BigEndian.Uint32(b[20:24])
	if serial == 0xffffffff {
		log.Println("Implausible serial, rejecting packet")
		return
	}

	if len(b) < 596 { // ensure enough data for all phase chunks
		log.Printf("Received packet is too small for phase data. Size: %d, Serial: %d", n, serial)
		return
	}
	firmware := binary.BigEndian.Uint32(b[600:604])

	log.Printf("Received datagram from meter at %v", src)
	log.Printf("Uid: %d", binary.BigEndian.Uint32(b[4:8]))
	log.Printf("Serial: %d", serial)
	log.Printf("Firmware: %d", firmware)

	powertot := (float32(binary.BigEndian.Uint32(b[32:36])) - float32(binary.BigEndian.Uint32(b[52:56]))) / 10.0
	forward := float64(binary.BigEndian.Uint64(b[40:48])) / 3600.0 / 1000.0
	reverse := float64(binary.BigEndian.Uint64(b[60:68])) / 3600.0 / 1000.0

	log.Printf("Total W: %.2f", powertot)
	log.Printf("Total Buy kWh:  %.4f", forward)
	log.Printf("Total Sell kWh: %.4f", reverse)

	isEnergyMeter := false
	if os.Getenv("SMA_ENERGY_METER") == "true" {
		isEnergyMeter = true
		log.Printf("SMA Energy Meter, NOT SHM2!")
	}

	var L1, L2, L3 *singlePhase

	if isEnergyMeter {
		L1 = decodePhaseChunk(b[160:304])
		L2 = decodePhaseChunk(b[304:448])
		L3 = decodePhaseChunk(b[448:592])
	} else {
		L1 = decodePhaseChunk(b[164:308])
		L2 = decodePhaseChunk(b[308:452])
		L3 = decodePhaseChunk(b[452:596])
	}

	printPhaseTable(L1, L2, L3)

	log.Printf("Packet (gzip+base64): %s", compressAndEncode(b))
}

func decodePhaseChunk(b []byte) *singlePhase {
	if len(b) < 136 { // ensure enough data for all fields
		return &singlePhase{}
	}
	forwardW := float32(binary.BigEndian.Uint32(b[4:8])) / 10
	reverseW := float32(binary.BigEndian.Uint32(b[24:28])) / 10

	L := singlePhase{}
	L.voltage = float32(binary.BigEndian.Uint32(b[132:136])) / 1000 // millivolts
	L.power = forwardW - reverseW
	if L.voltage != 0 {
		L.a = L.power / L.voltage
	}
	L.forward = float64(binary.BigEndian.Uint64(b[12:20])) / 3600.0 / 1000.0
	L.reverse = float64(binary.BigEndian.Uint64(b[32:40])) / 3600.0 / 1000.0

	return &L
}

func printPhaseTable(L1, L2, L3 *singlePhase) {
	log.Println("+-----+-------------+---------------+---------------+")
	log.Println("|value|   L1 \t|     L2  \t|   L3  \t|")
	log.Println("+-----+-------------+---------------+---------------+")
	log.Printf("|  V  | %8.2f \t| %8.2f \t| %8.2f \t|", L1.voltage, L2.voltage, L3.voltage)
	log.Printf("|  A  | %8.2f \t| %8.2f \t| %8.2f \t|", L1.a, L2.a, L3.a)
	log.Printf("|  W  | %8.2f \t| %8.2f \t| %8.2f \t|", L1.power, L2.power, L3.power)
	log.Printf("| kWh | %8.2f \t| %8.2f \t| %8.2f \t|", L1.forward, L2.forward, L3.forward)
	log.Printf("| kWh | %8.2f \t| %8.2f \t| %8.2f \t|", L1.reverse, L2.reverse, L3.reverse)
	log.Println("+-----+-------------+---------------+---------------+")
}
