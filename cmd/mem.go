package cmd

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"

	"github.com/urfave/cli/v3"
)

const gigabyte = 1 << 30

func MemCommand() *cli.Command {
	return &cli.Command{
		Name:  "mem",
		Usage: "Perform a memory stability test",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:     "maxMemory",
				Aliases:  []string{"m"},
				Usage:    "Maximum memory to allocate in gigabytes (GB)",
				Required: true,
			},
			&cli.UintFlag{
				Name:     "loopCount",
				Aliases:  []string{"c"},
				Usage:    "Number of test iterations",
				Required: true,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			maxMemory := c.Uint("maxMemory")
			loopCount := c.Uint("loopCount")

			return MemCheck(maxMemory, loopCount)
		},
	}
}

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
