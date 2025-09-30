package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"strings"
)

func StringList() *cli.Command {
	return &cli.Command{
		Name:  "strlist",
		Usage: "Build a string list with specified separator and wrapping characters",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "separator",
				Aliases:  []string{"s"},
				Required: false,
				Value:    ",",
				Usage:    "Delimiter for joining strings",
			},
			&cli.StringFlag{
				Name:     "warp",
				Aliases:  []string{"w"},
				Required: false,
				Value:    "'",
				Usage:    "Character to wrap around each string",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "str",
				Min:  1,
				Max:  -1,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			inputs := c.StringArgs("str")
			separator := c.String("separator")
			warp := c.String("warp")

			for i, input := range inputs {
				inputs[i] = warp + input + warp
			}

			rst := strings.Join(inputs, separator)

			fmt.Println(rst)
			return nil
		},
	}
}
