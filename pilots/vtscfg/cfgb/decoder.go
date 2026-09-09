package cfgb

import (
	"os"
)

func Decode(data []byte) []byte {
	out := make([]byte, len(data))

	for i, b := range data {
		out[i] = byte((int(b) - 88 + 256) % 256)
	}

	return out
}

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(Decode(data)), nil
}

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
