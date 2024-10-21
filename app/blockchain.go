package main

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"strings"
)

type Blockchain struct {
	Blocks     []Block
	Difficulty int
	Miner      string
}

func (bc *Blockchain) LoadFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer gzReader.Close()

	var blocks []Block
	if err := json.NewDecoder(gzReader).Decode(&blocks); err != nil {
		panic(err)
	}

	// Calculate the merkle root and transaction proof map for each block
	updatedBlocks := []Block{}
	for _, block := range blocks {
		merkleRoot, proofMap := block.Txs.MerkleRoot()
		block.Header.MerkleRoot = merkleRoot
		block.ProofMap = proofMap
		updatedBlocks = append(updatedBlocks, block)
	}

	bc.Blocks = updatedBlocks
}

func (bc *Blockchain) WriteToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	if err := json.NewEncoder(gzWriter).Encode(bc.Blocks); err != nil {
		panic(err)
	}
}

func (bc *Blockchain) PreviousBlock() Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) Time() int {
	return bc.PreviousBlock().Header.Timestamp + 10
}

func (bc *Blockchain) NewBlock(txs Transactions) Block {
	// Calculate the Merkle root and the proof map
	merkleRoot, proofMap := txs.MerkleRoot()

	// Create the header
	blockHeader := BlockHeader{
		Height:            bc.PreviousBlock().Header.Height + 1,
		PreviousBlockHash: bc.PreviousBlock().Header.Hash,
		Timestamp:         bc.Time(),
		MerkleRoot:        merkleRoot,
		TransactionsCount: len(txs),
		Miner:             bc.Miner,
		Nonce:             0,
		Difficulty:        bc.Difficulty,
	}

	// Repeatedly increment the Nonce until the hash has the required number of
	// leading zeros, which is given by the difficulty
	for {
		hash := blockHeader.Serialise()
		if hash[2:2+bc.Difficulty] == strings.Repeat("0", bc.Difficulty) {
			blockHeader.Hash = hash
			break
		}
		blockHeader.Nonce++
	}

	newBlock := Block{
		Header:   blockHeader,
		Txs:      txs,
		ProofMap: proofMap,
	}

	bc.Blocks = append(bc.Blocks, newBlock)

	return newBlock
}
