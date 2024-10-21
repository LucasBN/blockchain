`./blockchain.sh -blockchain-state ./data/blockchain.json.gz -mempool ./data/mempool.json.gz -command get-tx-hash -- -block 18 -tx 7`

`./blockchain.sh -blockchain-state ./data/blockchain.json.gz -mempool ./data/mempool.json.gz -command produce-blocks -- -n 1 -blockchain-output ./data/blockchain-new.json.gz`
