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

func (ts Transactions) TopN(n int) Transactions {
	// Sort the transactions by the fee (highest fee first)
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Fee > ts[j].Fee
	})

	// Take the first n transactions
	var selectedTxs Transactions
	for i := 0; i < n; i++ {
		if i >= len(ts) {
			break
		}
		selectedTxs = append(selectedTxs, ts[i])
	}

	return selectedTxs
}

func (ts Transactions) MerkleRoot() string {
	var hashes []string
	for _, tx := range ts {
		hashes = append(hashes, tx.Serialise())
	}

	for len(hashes) > 1 {
		if len(hashes)%2 != 0 {
			hashes = append(hashes, NullHash)
		}

		var newHashes []string
		for i := 0; i < len(hashes); i += 2 {
			hash := sha256.New()
			if hashes[i] < hashes[i+1] {
				hash.Write([]byte(hashes[i] + hashes[i+1]))
			} else {
				hash.Write([]byte(hashes[i+1] + hashes[i]))
			}
			newHashes = append(newHashes, "0x"+hex.EncodeToString(hash.Sum(nil)))
		}
		hashes = newHashes
	}

	return hashes[0]
}
