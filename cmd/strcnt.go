package cmd

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

func StringCount() *cli.Command {
	return &cli.Command{
		Name:  "strcnt",
		Usage: "Count string length, English and non-English characters",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "str",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {

			input := c.StringArg("str")
			total := len([]rune(input))
			englishCount := 0

			for _, char := range input {
				if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
					englishCount++
				}
			}

			nonEnglishCount := total - englishCount

			fmt.Printf("Total characters: %d\n", total)
			fmt.Printf("English characters: %d\n", englishCount)
			fmt.Printf("Non-English characters: %d\n", nonEnglishCount)

			return nil
		},
	}
}
