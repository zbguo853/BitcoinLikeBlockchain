package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/labstack/gommon/log"
	"strconv"
	"time"
)

//定义区块
type Block struct {
	Timestamp    int64  //时间线, 1970.1.1 00.00.00
	Data         []byte //交易数据,byte类型类似 uint8,实际上[]byte是哈希计算函数sha256.sum256()输入数据的格式
	PreBlockHash []byte
	Hash         []byte
	Nonce        int
}

// 由于要存进数据库, 得把整个block格式化为[]byte类型, 然后以hash-block键值对存进boltDB里面
func (b *Block) Serialize() []byte {

	blockBuffer := new(bytes.Buffer)
	enc := gob.NewEncoder(blockBuffer)

	// Encode (send) some values.
	err := enc.Encode(b)
	if err != nil {
		log.Panic("encode error:", err)
	}
	return blockBuffer.Bytes()
}

// decode 回一个Block出来
func Deserialize(b []byte) *Block {
	decBlock := Block{}

	dec := gob.NewDecoder(bytes.NewReader(b)) //

	err := dec.Decode(&decBlock)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}
	return &decBlock

}

//设定结构体对象哈希
//这是一个方法, 就是属于结构体的函数
func (block *Block) SetHash() {
	//处理当前的时间,转化为10进制的字符串,再转化为字节集
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	//叠加要哈希的数据,就是串在一起,
	// bytes.Join的使用格式就是这样, 第一个参数是[][]byte,里面是多个需要被串接的[]byte,
	// 第二个参数也是[]byte类型的,是插在多个被串接的[]byte之间的值
	headers := bytes.Join([][]byte{block.PreBlockHash, block.Data, timestamp}, []byte{})
	//计算出哈希
	hash := sha256.Sum256(headers) //返回的是 [size]byte 一般是[32]byte, 32*8=256位
	block.Hash = hash[:]           //[32]byte是一个数组,不能直接赋值给切片,要取所有元素再赋值
}

//创建一个区块
func NewBlock(data string, preBlockHash []byte) *Block {
	//block 是一个指针, 取得一个对象初始化之后的地址
	block := &Block{time.Now().Unix(), []byte(data), preBlockHash, []byte{}, 0}

	//block.SetHash() //SetHash是属于block的方法
	pow := NewProofOfWork(block) //挖矿附加这个区块
	nonce, hash := pow.mining()  //开始挖矿
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

//创建创世区块
func NewGenesisBlock() *Block {
	return NewBlock("创世区块", []byte{})
}
