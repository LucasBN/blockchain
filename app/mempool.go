package main

import (
	"compress/gzip"
	"encoding/json"
	"os"
)

type Mempool struct {
	Transactions Transactions
}

func (m *Mempool) LoadFromFile(filename string) {
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

	var transactions []Transaction
	if err := json.NewDecoder(gzReader).Decode(&transactions); err != nil {
		panic(err)
	}

	m.Transactions = transactions
}
