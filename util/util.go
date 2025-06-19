package util

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"log"

	"sma_home_manager_printer/phase"
)

func PrintPhaseTable(L1, L2, L3 *phase.SinglePhase) {
	log.Println("+-----+-------------+---------------+---------------+")
	log.Println("|value|   L1 \t|     L2  \t|   L3  \t|")
	log.Println("+-----+-------------+---------------+---------------+")
	log.Printf("|  V  | %8.2f \t| %8.2f \t| %8.2f \t|", L1.Voltage, L2.Voltage, L3.Voltage)
	log.Printf("|  A  | %8.2f \t| %8.2f \t| %8.2f \t|", L1.A, L2.A, L3.A)
	log.Printf("|  W  | %8.2f \t| %8.2f \t| %8.2f \t|", L1.Power, L2.Power, L3.Power)
	log.Printf("| kWh | %8.2f \t| %8.2f \t| %8.2f \t|", L1.Forward, L2.Forward, L3.Forward)
	log.Printf("| kWh | %8.2f \t| %8.2f \t| %8.2f \t|", L1.Reverse, L2.Reverse, L3.Reverse)
	log.Println("+-----+-------------+---------------+---------------+")
}

func CompressAndEncode(b []byte) string {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(b)
	if err != nil {
		return "compression error"
	}
	gz.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
