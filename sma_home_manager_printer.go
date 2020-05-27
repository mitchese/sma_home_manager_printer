package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/dmichael/go-multicast/multicast"
)

const (
	// TODO: This is always 239.12.255.254:9522, but maybe get it from env
	//       with a default of this address, in case it changes at some point
	address = "239.12.255.254:9522"
)

type singlePhase struct {
	voltage float32 // Volts: 230,0
	a       float32 // Amps: 8,3
	power   float32 // Watts: 1909
	forward float64 // kWh, purchased power
	reverse float64 // kWh, sold power
}

func main() {
	multicast.Listen(address, msgHandler)
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	// 0-28: SMA/SUSyID/SN/Uptime
	log.Println("-----------------------------------------------------")
	log.Println("Received datagram from meter")
	log.Println("Uid: ", binary.BigEndian.Uint32(b[4:8]))
	log.Println("Serial: ", binary.BigEndian.Uint32(b[20:24]))

	//              ...forward....                                 ...reverse...  both in 0.1W, converted to W
	powertot := ((float32(binary.BigEndian.Uint32(b[32:36])) - float32(binary.BigEndian.Uint32(b[52:56]))) / 10.0)

	// in watt seconds, convert to kWh
	forward := float64(binary.BigEndian.Uint64(b[40:48])) / 3600.0 / 1000.0
	reverse := float64(binary.BigEndian.Uint64(b[60:68])) / 3600.0 / 1000.0

	log.Println("Total W: ", powertot)
	log.Println("Total Buy kWh:  ", forward)
	log.Println("Total Sell kWh: ", reverse)

	L1 := decodePhaseChunk(b[164:308])
	L2 := decodePhaseChunk(b[308:452])
	L3 := decodePhaseChunk(b[452:596])

	log.Println("+-----+-------------+---------------+---------------+")
	log.Println("|value|   L1 \t|     L2  \t|   L3  \t|")
	log.Println("+-----+-------------+---------------+---------------+")
	log.Println(fmt.Sprintf("|  V  | %8.2f \t| %8.2f \t| %8.2f \t|", L1.voltage, L2.voltage, L3.voltage))
	log.Println(fmt.Sprintf("|  A  | %8.2f \t| %8.2f \t| %8.2f \t|", L1.a, L2.a, L3.a))
	log.Println(fmt.Sprintf("|  W  | %8.2f \t| %8.2f \t| %8.2f \t|", L1.power, L2.power, L3.power))
	log.Println(fmt.Sprintf("| kWh | %8.2f \t| %8.2f \t| %8.2f \t|", L1.forward, L2.forward, L3.forward))
	log.Println(fmt.Sprintf("| kWh | %8.2f \t| %8.2f \t| %8.2f \t|", L1.reverse, L2.reverse, L3.reverse))
	log.Println("+-----+-------------+---------------+---------------+")
}

func decodePhaseChunk(b []byte) *singlePhase {
	// why does this measure in 1/10 of watts?!
	forwardW := float32(binary.BigEndian.Uint32(b[4:8])) / 10
	reverseW := float32(binary.BigEndian.Uint32(b[24:28])) / 10

	// not used, but leaving here for future in case we need VA at some point
	//bezugVA := float32(binary.BigEndian.Uint32(b[84:88])) / 10
	//einspeiseVA := float32(binary.BigEndian.Uint32(b[104:108])) / 10

	L := singlePhase{}
	L.voltage = float32(binary.BigEndian.Uint32(b[132:136])) / 1000 // millivolts!
	L.power = forwardW - reverseW
	L.a = L.power / L.voltage                                                // victron needs A as well, make future me's life easier
	L.forward = float64(binary.BigEndian.Uint64(b[12:20])) / 3600.0 / 1000.0 // watt-seconds -> kWh
	L.reverse = float64(binary.BigEndian.Uint64(b[32:40])) / 3600.0 / 1000.0 // watt-seconds -> kWh

	return &L
}
