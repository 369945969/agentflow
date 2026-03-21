package agentcore

import (
	"crypto/rand"
	"encoding/hex"
)

func generateID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
