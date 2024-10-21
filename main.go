package main

import (
	"fmt"
)

func main() {
	command := "mine-block"

	// Load the blockchain from the file
	blockchain := &Blockchain{
		Difficulty: 3,
		Miner:      "0xca4388fb6d0ee25d59c24360e49c2dd4c9d02727",
	}
	blockchain.LoadFromFile("data/blockchain.json.gz")

	// Load the mempool from the file
	mempool := &Mempool{}
	mempool.LoadFromFile("data/mempool.json.gz")

	switch command {
	case "ld-most-rct-block-hash":
		fmt.Println(blockchain.PreviousBlock().Header.Hash)
	case "mine-block":
		// Fetch up to 100 of the most profitable transactions from the mempool
		// and create a new block
		newBlock := blockchain.NewBlock(mempool.Transactions.TopN(100))
		fmt.Println(newBlock.Header.Hash)

		blockchain.WriteToFile("data/blockchain-new.json.gz")
	default:
		panic(fmt.Sprintf("Unimplemented: command %s not supported", command))
	}
}
