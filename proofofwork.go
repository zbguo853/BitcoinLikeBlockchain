package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64 // 最大的64位整数, 挖矿时nonce的最大值
)

const targetBits = 24 // 对比的位数,当hash值的前24位均为0即满足要求,

type ProofOfWork struct {
	block  *Block   // 区块
	target *big.Int // 存储计算哈希对比的特定整数, 就是一个非常大的整形
	// 之所以使用big结构是为了方便和hash值做比较检查是否满足要求, 直接在数学上比大小
}

// 创建一个工作量证明的挖矿对象
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)                  // 初始化目标整数, 初始化一个大整数
	target.Lsh(target, uint(256-targetBits)) // left shift 左移(256-targetBits)位, 估计这样创建大数开销更小
	pow := &ProofOfWork{block, target}       // 创建对象
	return pow
}

// 准备数据进行挖矿计算,
// 属于ProofOfWork的方法
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PreBlockHash,        // 上一块哈希
			pow.block.Data,                // 当前数据
			IntToHex(pow.block.Timestamp), // 时间十六进制
			IntToHex(int64(targetBits)),   // 位数,16进制
			IntToHex(int64(nonce)),        // 保存工作量证明的 nonce
		}, []byte{}, // 为啥这里要逗号?
	)
	return data
}

// 挖矿执行过程
func (pow *ProofOfWork) mining() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	fmt.Printf("Mining the block containing: %s\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce) // 准备好数据
		hash = sha256.Sum256(data)     // 计算出哈希
		//fmt.Printf("\r%x", hash)       // 打印显示哈希
		hashInt.SetBytes(hash[:])      // 把hash转换成大整形, 获取要对比的数据
		if hashInt.Cmp(pow.target) == -1 { // 检查nonce是否符合挖到矿的条件
			fmt.Printf("\r%x", hash)
			fmt.Println("\n\n")
			break
		} else {
			nonce++
		}
		//fmt.Println("\n\n")
	}
	return nonce, hash[:] // nonce相当于解题的答案
}

// 同样是属于ProofOfWork的方法, 但这是用来校验新块的, 在接收别人的新块的时候用
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)   // 准备好的数据
	hash := sha256.Sum256(data)                // 计算出哈希
	hashInt.SetBytes(hash[:])                  // 把hash转换成大整形, 获取要对比的数据
	isValid := (hashInt.Cmp(pow.target) == -1) // 校验数据
	return isValid
}

