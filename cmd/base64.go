package cmd

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/LiZeC123/gmh/util"
	"github.com/urfave/cli/v3"
)

// Base64Command 返回一个支持 Base64 编码/解码的 CLI 命令
func Base64Command() *cli.Command {
	return &cli.Command{
		Name:  "base64",
		Usage: "Base64 encode or decode data",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "encode",
				Aliases: []string{"e"},
				Usage:   "Encode input as base64 (default is decode)",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:      "data",
				Min:       0,
				Max:       1,
				UsageText: "Data to encode/decode. If omitted, reads from standard input",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			// 获取输入数据
			inputs, err := util.GetArgsOrStdinInput(c, "data")
			if err != nil {
				return err
			}
			
			count := len(inputs)
			if count == 0 {
				return errors.New("no input provided")
			}
			if count != 1 {
				return errors.New("only one input is allowed")
			}


			input := inputs[0]
			encode := c.Bool("encode")
			var result string
			if encode {
				result = base64.StdEncoding.EncodeToString([]byte(input))
			} else {
				decoded, err := base64.StdEncoding.DecodeString(input)
				if err != nil {
					return fmt.Errorf("invalid base64 input: %w", err)
				}
				result = string(decoded)
			}

			fmt.Println(result)
			return nil
		},
	}
}