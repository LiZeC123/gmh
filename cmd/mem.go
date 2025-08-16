package cmd

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
)

const gigabyte = 1 << 30

func MemCheck(memoryCount uint, loopCount uint) error {
	memory := make([]byte, memoryCount*gigabyte)

	for i := range loopCount {
		fmt.Printf("Start Loop %d Random Write\n", i)
		// 向内存中写入随机数据
		rand.Read(memory)

		checksum1 := md5.Sum(memory)
		fmt.Printf("\tCheckSum For Loop %d: %x\n", i, checksum1)

		checksum2 := md5.Sum(memory)
		fmt.Printf("\tCheckSum For Loop %d: %x\n", i, checksum2)

		// 连续读取两次, 检查内存数据是否一致
		if checksum1 == checksum2 {
			fmt.Printf("CheckSum is same\n\n")
		} else {
			fmt.Printf("Checksum is not same\n\n")
		}
	}

	return nil
}
