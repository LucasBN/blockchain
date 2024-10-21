`./blockchain.sh -blockchain-state ./data/blockchain.json.gz -mempool ./data/mempool.json.gz -command get-tx-hash -- -block 18 -tx 7`

`./blockchain.sh -blockchain-state ./data/blockchain.json.gz -mempool ./data/mempool.json.gz -command produce-blocks -- -n 1 -blockchain-output ./data/blockchain-new.json.gz`

`./blockchain.sh -blockchain-state ./data/blockchain.json.gz -mempool ./data/mempool.json.gz -command generate-proof -- -block 18 -tx-hash 0x49276072ad4033ec9644d8831167125d6deb16747c927f206329efb5de62b77c -o data/proof.json`

`./blockchain.sh -blockchain-state ./data/blockchain.json.gz -mempool ./data/mempool.json.gz -command verify-proof -- -f data/proof.json`
