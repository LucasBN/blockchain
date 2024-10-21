package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

type Block struct {
	Header   BlockHeader  `json:"header"`
	Txs      Transactions `json:"transactions"`
	ProofMap map[string]string
}

func (b Block) ProveIncludesTxHash(txHash string) (bool, []string) {
	// Check if the proof map contains an entry for the transcation hash
	sibling, exists := b.ProofMap[txHash]
	if !exists {
		// Compare the hash to the merkle
		return txHash == b.Header.MerkleRoot, []string{txHash}
	}

	// Hash the sibling and
	hash := sha256.New()
	if txHash < sibling {
		hash.Write([]byte(txHash + sibling))
	} else {
		hash.Write([]byte(sibling + txHash))
	}

	nextHash := "0x" + hex.EncodeToString(hash.Sum(nil))

	valid, proofHashes := b.ProveIncludesTxHash(nextHash)
	return valid, append(proofHashes, sibling)
}

type BlockHeader struct {
	Height            int    `json:"height"`
	PreviousBlockHash string `json:"previous_block_header_hash"`
	Timestamp         int    `json:"timestamp"`
	MerkleRoot        string `json:"transactions_merkle_root"`
	TransactionsCount int    `json:"transactions_count"`
	Miner             string `json:"miner"`
	Nonce             int    `json:"nonce"`
	Difficulty        int    `json:"difficulty"`
	Hash              string `json:"hash"`
}

func (h BlockHeader) Serialise() string {
	keys := []string{
		strconv.Itoa(h.Difficulty),
		strconv.Itoa(h.Height),
		h.Miner,
		strconv.Itoa(h.Nonce),
		h.PreviousBlockHash,
		strconv.Itoa(h.Timestamp),
		strconv.Itoa(h.TransactionsCount),
		h.MerkleRoot,
	}

	hash := sha256.New()
	hash.Write([]byte(strings.Join(keys, ",")))
	return "0x" + hex.EncodeToString(hash.Sum(nil))
}
