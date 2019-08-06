package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

//定义区块结构
type Block struct {
	//版本号
	Version uint64

	// 前区块哈希
	PrevHash []byte

	//交易的根哈希值
	MerkleRoot []byte

	//时间戳
	TimeStamp uint64

	//难度值, 系统提供一个数据，用于计算出一个哈希值
	Bits uint64

	//随机数，挖矿要求的数值
	Nonce uint64

	// 哈希, 为了方便，我们将当前区块的哈希放入Block中
	Hash []byte

	//数据
	
	Transactions []*Transaction
}

//创建一个区块（提供一个方法）
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	b := Block{
		Version:    0,
		PrevHash:   prevHash,
		MerkleRoot: nil, 
		TimeStamp:  uint64(time.Now().Unix()),

		Bits:  0, 
		Nonce: 0, 
		Hash:  nil,
		// Data:  
		Transactions: txs,
	}

	//填充梅克尔根值
	b.HashTransactionMerkleRoot()
	fmt.Printf("merkleRoot:%x\n", b.MerkleRoot)

	//将POW集成到Block中
	pow := NewProofOfWork(&b)
	hash, nonce := pow.Run()
	b.Hash = hash
	b.Nonce = nonce

	return &b
}

//绑定Serialize方法， gob编码
func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer
	
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(b)
	if err != nil {
		fmt.Printf("Encode err:", err)
		return nil
	}

	return buffer.Bytes()
}

//反序列化，输入[]byte，返回block
func Deserialize(src []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(src))
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Printf("decode err:", err)
		return nil
	}

	return &block
}

//简易梅克尔根，把所有的交易拼接到一起，做哈希处理，最终赋值给block.MerKleRoot
func (block *Block) HashTransactionMerkleRoot() {

	var info [][]byte

	for _, tx := range block.Transactions {
		
		txHashValue := tx.TXID //[]byte
		info = append(info, txHashValue)
	}

	value := bytes.Join(info, []byte{})
	hash := sha256.Sum256(value)

	//讲hash值赋值MerKleRoot字段
	block.MerkleRoot = hash[:]
}
