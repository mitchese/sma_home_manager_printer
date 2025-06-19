package handler

import (
	"encoding/binary"
	"log"
	"net"
	"os"

	"sma_home_manager_printer/phase"
	"sma_home_manager_printer/util"
)

func MsgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println("-----------------------------------------------------")
	if len(b) < 200 {
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

	if len(b) < 596 {
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

	var L1, L2, L3 *phase.SinglePhase

	if isEnergyMeter {
		L1 = phase.DecodePhaseChunk(b[160:304])
		L2 = phase.DecodePhaseChunk(b[304:448])
		L3 = phase.DecodePhaseChunk(b[448:592])
	} else {
		L1 = phase.DecodePhaseChunk(b[164:308])
		L2 = phase.DecodePhaseChunk(b[308:452])
		L3 = phase.DecodePhaseChunk(b[452:596])
	}

	util.PrintPhaseTable(L1, L2, L3)

	log.Printf("Packet (gzip+base64): %s", util.CompressAndEncode(b))
}
