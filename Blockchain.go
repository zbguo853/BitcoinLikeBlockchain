package main

import (
	"github.com/boltdb/bolt"
)

const blockchainName = "blockchain"

// BlockChain结构体不仅包含链,还包含方法
type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

// 因为要用一个迭代器来爬链, 重新设置一个类型比较好
type BlockchainInterator struct {
	hashofBlocktoRead []byte
	db          *bolt.DB
}

func (chain *BlockChain) AddBlock(data string) {

	chain.db.Update(func(tx *bolt.Tx) error {
		bcbucket := tx.Bucket([]byte(blockchainName))
		newblock := NewBlock(data, chain.tip)
		_ = bcbucket.Put(newblock.Hash, newblock.Serialize())
		_ = bcbucket.Put([]byte("lastHash"), newblock.Hash)
		chain.tip = newblock.Hash
		return nil
	})
}

// 创建一个区块链实例, 这不一定是创建一条新链, 如果数据库文件已经有区块链, 那这个实例就代表文件中的区块链
func NewBlockchain() *BlockChain {
	tip := []byte{}

	db, _ := bolt.Open("my.db", 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		bcbucket := tx.Bucket([]byte(blockchainName))

		if bcbucket == nil {
			bcbucket, _ := tx.CreateBucket([]byte(blockchainName))
			genesisblock := NewGenesisBlock()
			_ = bcbucket.Put(genesisblock.Hash, genesisblock.Serialize())
			_ = bcbucket.Put([]byte("lastHash"), genesisblock.Hash)
			tip = bcbucket.Get([]byte("lastHash"))
		} else {
			tip = bcbucket.Get([]byte("lastHash"))
		}

		return nil
	})

	bc := BlockChain{tip, db}

	return &bc
}

func (bi *BlockchainInterator) Readblock() *Block {
	//curblock := Block{}
	var curblock *Block

	bi.db.View(func(tx *bolt.Tx) error {
		bcbucket := tx.Bucket([]byte(blockchainName))
		curblockBeforeDec := bcbucket.Get([]byte(bi.hashofBlocktoRead))
		curblock = Deserialize(curblockBeforeDec)
		bi.hashofBlocktoRead = curblock.PreBlockHash
		return nil
	})
	return curblock
}
