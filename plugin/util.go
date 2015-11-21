package plugin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
)

func decodeChecksum(data string) ([]byte, error) {
	b, err := hex.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("invalid checksum: %s", err)
	}
	return b, nil
}

func checksum(source string, h hash.Hash, v []byte) error {
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("Failed to open file for checksum: %s", err)
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("Failed to hash: %s", err)
	}

	if actual := h.Sum(nil); !bytes.Equal(actual, v) {
		return fmt.Errorf(
			"Checksums did not match.\nExpected: %s\nGot: %s",
			hex.EncodeToString(v),
			hex.EncodeToString(actual),
		)
	}

	return nil
}
