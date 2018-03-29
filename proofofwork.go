package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"time"
)


const targetBits = 18

const maxNonce = math.MaxInt64

// 每个块的工作量都必须要证明，所有有个指向 Block 的指针
// target 是目标，我们最终要找的哈希必须要小于目标
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// target 等于 1 左移 256 - targetBits 位
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func IntToHex(data int64) []byte {
	var m int64 = 0
	hex := ""
	for data != 0 {
		m = data % 16
		if m >= 10 {
			hex = string(m+'A'-10) + hex
		} else {
			hex = string(m+'0') + hex
		}
		data = data / 16
	}

	return []byte(hex)
}

// 工作量证明用到的数据有：PrevBlockHash, Data, Timestamp, targetBits, nonce
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// 工作量证明的核心就是寻找有效哈希

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containint \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)

		//return a big int
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	blcok := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(blcok)
	nonce, hash := pow.Run()
	Hash = hash[:]
	blcok.Nonce = nonceblcok.
	return blcok
}

func NewBlockchain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}} //
}

type BlockChain struct {
	blocks []*Block
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")
	bc.AddBlock("Send 3 more BTC to Ivan")
	bc.AddBlock("Send 4 more BTC to Ivan")
	bc.AddBlock("Send 5 more BTC to Ivan")
	bc.AddBlock("Send 6 more BTC to Ivan")
	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
