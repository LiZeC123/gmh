package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func UniqueCommand() *cli.Command {
	return &cli.Command{
		Name:  "unique",
		Usage: "Remove duplicates lines from a file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Required: true,
				Usage:    "Input file to split",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			filePath := c.String("input")

			file, err := os.Open(filePath)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Failed to open file: %v", err), 1)
			}
			defer file.Close()

			// 使用map记录已出现的行
			seen := make(map[string]bool)
			scanner := bufio.NewScanner(file)

			// 逐行扫描文件
			for scanner.Scan() {
				line := scanner.Text()
				if !seen[line] {
					seen[line] = true
					fmt.Println(line) // 输出唯一行
				}
			}

			// 检查扫描错误
			if err := scanner.Err(); err != nil {
				return err
			}

			return nil
		},
	}
}
