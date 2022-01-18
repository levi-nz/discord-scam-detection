package main

import (
	"crypto/sha256"
	"encoding/hex"
)

var (
	// maliciousImages contains hex-encoded SHA-256 hashes of
	// malicious images to check for.
	maliciousImages = []string{
		"04be4e2d5e0d02319440715328149cf0347d62e3b51a19b07cddfbacb6542468", // https://discord.com/assets/652f40427e1f5186ad54836074898279.png
		"890a7bb7147f5a5b105e6ea82d31ccaa397f67fcc6b3efa19b9cbb2018b99653", // Nitro "get chat perks + 2 boosts"
	}
)

// IsImageMalicious checks if the given image is considered malicious
// based on its SHA-256 checksum.
//
// If the image is malicious, the encoded checksum is returned.
func IsImageMalicious(image []byte) (encodedChecksum string, isMalicious bool) {
	checksum := sha256.Sum256(image)
	encodedChecksum = hex.EncodeToString(checksum[:])

	for _, maliciousChecksum := range maliciousImages {
		if encodedChecksum == maliciousChecksum {
			return encodedChecksum, true
		}
	}

	return "", false
}
