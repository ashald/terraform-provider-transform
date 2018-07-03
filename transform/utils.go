package transform

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func getSHA256(input interface{}) (string, error) {
	serialized, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(serialized)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
