package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"testing"
)

func TestMsgHandler_TotalBuyKWh(t *testing.T) {
	testCases := []struct {
		name        string
		b64         string
		envValue    string // "true", "false", or ""
		expectedLog string
	}{
		{
			name:        "SMA_ENERGY_METER true",
			b64:         `H4sIAAAJbogA/+zQvW9SYRTH8d9z7xPDvYmGwQssXCE4yMsgb6OAIiAaY6IOToZRnNQJjTE38SUuJrhrvHZom05t04FuTE0KDXRqkw4NYwlDhzadG05P/4VO57uc5MN5csh98fQ+oA0fgDIeIthqqxvvq40BsDyD0gCQgArM591f0y0fBhlgkDnr7768gslmkqm1bx970HNTb6DJrI2dzx4ssrewyPLJSXsCm9/aZKHTg9dHuE7WgsP3HfrNGviBVYR4P0Tm2rNuEGG2MBke36z2EOH7ETL9zPvRQZTvR8ns48NxHy6/dcnclw82fdwi+4+YBszvXxHn/5Pk3STtWreb5gQpthSZ8el30UOaLU0GZ/AvgAxbhkwv7G33kGXLktmj8ZM+cmw5MvNaIr6LPFuBZ5HnvctJu6nI2dJPlNhKZGpxutJBma1MZqpRwUOFrUJmfPg7rKDGViNLN/b/3EGdrU6mhsGTGBpsj3g2L2Z3/nUN/RySJEmSJEmSJEmSJEmSJEmSJEmSJEnSlXc+APoUjAgAIAAA`,
			envValue:    "true",
			expectedLog: "|  W  |     3.50 	|     0.00 	|     0.00 	|",
		},
		{
			name:        "SMA_ENERGY_METER false",
			b64:         `H4sIAAAJbogA/+zOy2sTURTH8d/MXG1nZNpUmSa2TU19kUm76LRRFEkbldRHImiD+yx1EQgi4saQZTQuZiPoxhc+Fu7MRoRIVIQsBQMGA0JWBlwpLgbcyDncv8HV+S7mcj+cM9zyhZOAMh8DMMwSYpWrxvX2/Ilnt4fzVRgKnDFJ30PRlx8NmGzHYLLZhbDfhKXnLDbsbv3tQpEdbkOx7fd/f41g6zmbLbt3cKkJR885bFPDGT+ES4admKLz43F4es/jmT13Tn+6glmyVAqzbE66fr6HuJ6LsyFz9GETCTL7JRJs9tOf22Mk9VySLZ4+UBxjkWyphUW2XcHlbgP7yGIHkVKA9fw9luhuvYNPp3MKPs/arVKnhoz+Z4Ztole+H2FZ2zIbFpKrwAqZ+w0rbPZdFeURkM3kELC5b/sJYE3vrrFNvn5UD7FONvEBWXrPvRGO0N0sI0en9wI5nvVqN7/3sKH3N9is0Z/qG2xq22TDjle/VpEnmw6QZ5vO3iimUSCbu4YC21zxVukztvTuFpsadC5WcIbMXcBZes+DJzhHd2s9pFfF1DYkSZIkSZIkSZIkSZIkSZIkSZIkSZKk/9a/AQAEzZAOACAAAA==`,
			envValue:    "false",
			expectedLog: "|  W  |  -822.40 	|   262.60 	|   554.20 	|",
		},
		// Add more test cases here as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				os.Setenv("SMA_ENERGY_METER", tc.envValue)
			} else {
				os.Unsetenv("SMA_ENERGY_METER")
			}

			data, err := base64.StdEncoding.DecodeString(tc.b64)
			if err != nil {
				t.Fatalf("base64 decode failed: %v", err)
			}
			zr, err := gzip.NewReader(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("gzip reader failed: %v", err)
			}
			decompressed, err := io.ReadAll(zr)
			if err != nil {
				t.Fatalf("gzip decompress failed: %v", err)
			}
			zr.Close()

			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			// Call MsgHandler with a dummy UDPAddr
			MsgHandler(&net.UDPAddr{}, len(decompressed), decompressed)

			// Print the phase table and all log output
			output := buf.String()
			if testing.Verbose() {
				os.Stdout.WriteString(output)
			}

			if !strings.Contains(output, tc.expectedLog) {
				t.Errorf("log output does not contain expected value: %s", tc.expectedLog)
			}
		})
	}
}
