package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define the main arguments
	blockchainStateFile := flag.String("blockchain-state", "", "Path to the blockchain state file")
	mempoolPath := flag.String("mempool", "", "Path to the mempool file")
	command := flag.String("command", "", "The command to execute (produce-blocks or get-tx-hash)")

	// Parse the main arguments
	flag.Parse()

	// Check if command is provided
	if *command == "" {
		fmt.Println("Error: Command is required")
		flag.Usage()
		os.Exit(1)
	}

	// Load the blockchain from the file
	blockchain := &Blockchain{
		Difficulty: 3,
		Miner:      "0xca4388fb6d0ee25d59c24360e49c2dd4c9d02727",
	}
	blockchain.LoadFromFile(*blockchainStateFile)

	// Load the mempool from the file
	mempool := &Mempool{}
	mempool.LoadFromFile(*mempoolPath)

	switch *command {
	case "produce-blocks":
		produceBlocksCmd := flag.NewFlagSet("produce-blocks", flag.ExitOnError)
		nBlocks := produceBlocksCmd.Int("n", 1, "Number of blocks to produce")
		blockchainOutput := produceBlocksCmd.String("blockchain-output", "", "Path to the blockchain output file")
		// mempoolOutput := produceBlocksCmd.String("mempool-output", "", "Path to the mempool output file")

		// Parse the subcommand-specific flags
		if err := produceBlocksCmd.Parse(flag.Args()); err != nil {
			fmt.Println("Error parsing get-tx-hash command:", err)
			os.Exit(1)
		}

		for i := 0; i < *nBlocks; i++ {
			blockchain.NewBlock(mempool.Transactions.TopN(100))
		}

		blockchain.WriteToFile(*blockchainOutput)
	case "get-tx-hash":
		getTxHashCmd := flag.NewFlagSet("get-tx-hash", flag.ExitOnError)
		block := getTxHashCmd.Int("block", -1, "Block number")
		tx := getTxHashCmd.Int("tx", -1, "Transaction index")

		// Just parse the remaining arguments directly since -- is already stripped
		if err := getTxHashCmd.Parse(flag.Args()); err != nil {
			fmt.Println("Error parsing get-tx-hash command:", err)
			os.Exit(1)
		}

		// Check if block and tx are provided
		if *block == -1 || *tx == -1 {
			fmt.Println("Error: block and tx parameters are required for get-tx-hash")
			getTxHashCmd.Usage()
			os.Exit(1)
		}

		fmt.Println(blockchain.Blocks[*block].Txs[*tx].Serialise())
	default:
		fmt.Println("Error: Unknown command", *command)
		flag.Usage()
		os.Exit(1)
	}
}
