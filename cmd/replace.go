package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

func Replace() *cli.Command {
	return &cli.Command{
		Name:  "replace",
		Usage: "Replace all occurrences of OLD with NEW in input",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "old",
				Aliases:  []string{"o"},
				Required: true,
				Usage:    "String to be replaced",
			},
			&cli.StringFlag{
				Name:     "new",
				Aliases:  []string{"n"},
				Required: true,
				Usage:    "Replacement string",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			// 获取命令行参数
			oldStr := c.String("old")
			newStr := c.String("new")

			// 读取标准输入
			input, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}

			// 执行替换操作
			output := strings.ReplaceAll(string(input), oldStr, newStr)

			// 输出结果到标准输出
			fmt.Print(output)
			return nil
		},
	}
}
