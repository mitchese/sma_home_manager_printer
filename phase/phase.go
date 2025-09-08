package phase

import (
	"encoding/binary"
)

type SinglePhase struct {
	Voltage float32 // Volts
	A       float32 // Amps
	Power   float32 // Watts
	Forward float64 // kWh, purchased power
	Reverse float64 // kWh, sold power
}

func DecodePhaseChunk(b []byte) *SinglePhase {
	if len(b) < 136 {
		return &SinglePhase{}
	}
	forwardW := float32(binary.BigEndian.Uint32(b[4:8])) / 10
	reverseW := float32(binary.BigEndian.Uint32(b[24:28])) / 10

	L := SinglePhase{}
	L.Voltage = float32(binary.BigEndian.Uint32(b[132:136])) / 1000 // millivolts
	L.Power = forwardW - reverseW
	if L.Voltage != 0 {
		L.A = L.Power / L.Voltage
	}
	L.Forward = float64(binary.BigEndian.Uint64(b[12:20])) / 3600.0 / 1000.0
	L.Reverse = float64(binary.BigEndian.Uint64(b[32:40])) / 3600.0 / 1000.0

	return &L
}
