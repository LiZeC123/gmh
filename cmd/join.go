package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/LiZeC123/gmh/util"
	"github.com/urfave/cli/v3"
)

func JoinCommand() *cli.Command {
	return &cli.Command{
		Name:  "join",
		Usage: "Join strings with given separator",
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "string",
				Min:  0,
				Max:  -1,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Required: false,
				Usage:    "Input file containing URLs (one per line). Use '-' for stdin",
			},
			&cli.StringFlag{
				Name:     "separator",
				Aliases:  []string{"s"},
				Required: true,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			data, err := util.GetAllInput(c, "string", "input")
			if err != nil {
				return err
			}
			if len(data) == 0 {
				return fmt.Errorf("no string provided. Use command arguments, --input, or stdin")
			}

			separator := c.String("separator")
			output := strings.Join(data, separator)
			fmt.Println(output)
			return nil
		},
	}
}
