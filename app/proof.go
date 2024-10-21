package main

import (
	"crypto/sha256"
	"encoding/hex"
)

type Proof struct {
	TransactionHash string   `json:"transaction_hash"`
	MerkleRoot      string   `json:"merkle_root"`
	ProofHashes     []string `json:"proof_hashes"`
}

func (p Proof) Verify() bool {
	currentHash := p.TransactionHash

	// Iterate over the hashes in reverse order
	for i := len(p.ProofHashes) - 1; i >= 0; i-- {
		proofHash := p.ProofHashes[i]
		hash := sha256.New()
		if currentHash < proofHash {
			hash.Write([]byte(currentHash + proofHash))
		} else {
			hash.Write([]byte(proofHash + currentHash))
		}
		currentHash = "0x" + hex.EncodeToString(hash.Sum(nil))
	}

	return currentHash == p.MerkleRoot
}
