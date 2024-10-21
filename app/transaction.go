package main

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"
)

type Transaction struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Amount    int    `json:"amount"`
	Fee       int    `json:"transaction_fee"`
	LockTime  int    `json:"lock_time"`
	Signature string `json:"signature"`
}

type Transactions []Transaction

func (t Transaction) Serialise() string {
	keys := []string{
		strconv.Itoa(t.Amount),
		strconv.Itoa(t.LockTime),
		t.Receiver,
		t.Sender,
		t.Signature,
		strconv.Itoa(t.Fee),
	}

	// Hash the joined keys
	hash := sha256.New()
	hash.Write([]byte(strings.Join(keys, ",")))
	return "0x" + hex.EncodeToString(hash.Sum(nil))
}

func (ts Transactions) TopN(n int, timestamp int) Transactions {
	// Filter out transactions that are locked
	var unlockedTxs Transactions
	for _, tx := range ts {
		if tx.LockTime <= timestamp {
			unlockedTxs = append(unlockedTxs, tx)
		}
	}

	// Sort the transactions by the fee (highest fee first)
	sort.Slice(unlockedTxs, func(i, j int) bool {
		return unlockedTxs[i].Fee > unlockedTxs[j].Fee
	})

	// Take the first n transactions
	var selectedTxs Transactions
	for i := 0; i < n; i++ {
		if i >= len(unlockedTxs) {
			break
		}
		selectedTxs = append(selectedTxs, unlockedTxs[i])
	}

	return selectedTxs
}

func (ts Transactions) MerkleRoot() (string, map[string]string) {
	var hashes []string
	for _, tx := range ts {
		hashes = append(hashes, tx.Serialise())
	}

	return generateMerkleRoot(hashes, make(map[string]string))
}

func generateMerkleRoot(hashes []string, proofMap map[string]string) (string, map[string]string) {
	// Base case: there's only one hash, so return it
	if len(hashes) == 1 {
		return hashes[0], proofMap
	}

	// Otherwise, recursively calculate the Merkle root
	if len(hashes)%2 != 0 {
		hashes = append(hashes, NullHash)
	}

	var newHashes []string
	for i := 0; i < len(hashes); i += 2 {
		// Add the hashes to the proof map
		proofMap[hashes[i]] = hashes[i+1]
		proofMap[hashes[i+1]] = hashes[i]

		// Combine the hashes into a new single hash
		hash := sha256.New()
		if hashes[i] < hashes[i+1] {
			hash.Write([]byte(hashes[i] + hashes[i+1]))
		} else {
			hash.Write([]byte(hashes[i+1] + hashes[i]))
		}
		newHashes = append(newHashes, "0x"+hex.EncodeToString(hash.Sum(nil)))
	}

	return generateMerkleRoot(newHashes, proofMap)
}
