package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	from   string
	to     string
	amount float64
}

type Block struct {
	data         Data
	hash         string
	previousHash string
	timestamp    time.Time
	pow          uint64
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   uint64
}

func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.data)
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.FormatUint(b.pow, 10)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) mine(difficulty uint64) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", int(difficulty))) {
		b.pow++
		b.hash = b.calculateHash()
	}
}

func CreateBlockchain(difficulty uint64) Blockchain {
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}
	return Blockchain{
		genesisBlock: genesisBlock,
		chain:        []Block{genesisBlock},
		difficulty:   difficulty,
	}
}

func (b *Blockchain) addBlock(from, to string, amount float64) {
	blockData := Data{from, to, amount}
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}
	newBlock.mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
}

func (b Blockchain) isValid() bool {
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.hash != currentBlock.calculateHash() ||
			currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}

	return true
}

func (b Blockchain) GetBlockData(i int) (Data, error) {
	if i < 0 || i >= len(b.chain) {
		return Data{}, fmt.Errorf("index out of range")
	}
	return b.chain[i].data, nil
}

func main() {
	blockchain := CreateBlockchain(2)

	blockchain.addBlock("Alice", "Bob", 5)
	blockchain.addBlock("John", "Bob", 2)

	fmt.Println(blockchain.GetBlockData(1))

	fmt.Println(blockchain.isValid())
}
