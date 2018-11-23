package main

import (
	"bytes"
	"encoding/binary"
	"github.com/labstack/gommon/log"
)

//整数转化为16进制
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)                        //开辟内存,存储字节集
	err := binary.Write(buff, binary.BigEndian, num) //num转化字节集合写入, 大端字节序
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes() // .Bytes()方法会返回所有未读数据,但改变这些数据就会改变Buffer中的内容
}






