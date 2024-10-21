package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
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
			txs := mempool.Transactions.TopN(100, blockchain.Time())
			blockchain.NewBlock(txs)
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
	case "generate-proof":
		generateProofCmd := flag.NewFlagSet("generate-proof", flag.ExitOnError)
		block := generateProofCmd.Int("block", -1, "Block number")
		txHash := generateProofCmd.String("tx-hash", "", "Transaction hash")
		outputFile := generateProofCmd.String("o", "", "Output file")

		// Parse the subcommand-specific flags
		if err := generateProofCmd.Parse(flag.Args()); err != nil {
			fmt.Println("Error parsing get-tx-hash command:", err)
			os.Exit(1)
		}

		valid, proofHashes := blockchain.Blocks[*block].ProveIncludesTxHash(*txHash)
		if !valid {
			fmt.Println("Invalid proof: the transaction is not included in this block")
			os.Exit(1)
		}

		proof := Proof{
			TransactionHash: *txHash,
			MerkleRoot:      proofHashes[0],
			ProofHashes:     proofHashes[1:],
		}

		jsonData, err := json.MarshalIndent(proof, "", "  ") // Pretty print JSON
		if err != nil {
			log.Fatalf("Error marshaling to JSON: %v", err)
		}

		err = os.WriteFile(*outputFile, jsonData, 0644)
		if err != nil {
			log.Fatalf("Error writing JSON to file: %v", err)
		}
	case "verify-proof":
		verifyProofCmd := flag.NewFlagSet("generate-proof", flag.ExitOnError)
		proofFile := verifyProofCmd.String("f", "", "Proof file")

		// Parse the subcommand-specific flags
		if err := verifyProofCmd.Parse(flag.Args()); err != nil {
			fmt.Println("Error parsing get-tx-hash command:", err)
			os.Exit(1)
		}

		file, err := os.Open(*proofFile)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer file.Close()

		// Read the file contents
		byteValue, err := io.ReadAll(file)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		// Unmarshal the JSON into the Proof struct
		var proof Proof
		err = json.Unmarshal(byteValue, &proof)
		if err != nil {
			log.Fatalf("Error unmarshaling JSON: %v", err)
		}

		fmt.Println(proof.Verify())
	default:
		fmt.Println("Error: Unknown command", *command)
		flag.Usage()
		os.Exit(1)
	}
}
