package cmd

import (
	"context"
	"errors"
	"fmt"
	"unicode"

	"github.com/LiZeC123/gmh/util"
	"github.com/urfave/cli/v3"
)

func CaseCommand() *cli.Command {
	return &cli.Command{
		Name:  "case",
		Usage: "Convert variable names between different case styles",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "snake",
				Aliases:  []string{"s"},
				Required: false,
				Usage:    "Convert to snake_case (lowercase_with_underscores)",
			},
			&cli.BoolFlag{
				Name:     "camel",
				Aliases:  []string{"c"},
				Required: false,
				Usage:    "Convert to camelCase (lowerCamelCase)",
			},
			&cli.BoolFlag{
				Name:     "pascal",
				Aliases:  []string{"p"},
				Required: false,
				Usage:    "Convert to PascalCase (UpperCamelCase)",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:      "names",
				Min:       0,
				Max:       -1,
				UsageText: "Variable names to convert (read from stdin if omitted)",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			inputs, err := util.GetArgsOrStdinInput(c, "names")
			if err != nil {
				return err
			}

			var f func(string) string
			if c.Bool("snake") {
				f = CamelToSnake
			}
			if c.Bool("camel") {
				f = SnakeToCamel
			}
			if c.Bool("pascal") {
				f = SnakeToPascal
			}

			if f == nil {
				return errors.New("no conversion format specified (use --snake, --camel or --pascal")
			}

			for _, varName := range inputs {
				fmt.Printf("%s\n", f(varName))
			}

			return nil
		},
	}
}

// CamelToSnake 将驼峰命名转换为蛇形命名
// helloWorld -> hello_world, MyHTTPRequest -> my_http_request
func CamelToSnake(s string) string {
	if s == "" {
		return ""
	}

	var result []rune
	runes := []rune(s)

	for i, r := range runes {
		// 如果是大写字母
		if unicode.IsUpper(r) {
			// 不是第一个字符，且前一个字符不是大写，或者下一个字符是小写时，添加下划线
			if i > 0 {
				prev := runes[i-1]
				// 当前是大写，前一个不是大写，或者（前一个是大写且下一个是小写）时加下划线
				if !unicode.IsUpper(prev) || (i+1 < len(runes) && unicode.IsLower(runes[i+1])) {
					result = append(result, '_')
				}
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}

// SnakeToCamel 将蛇形命名转换为驼峰命名
// hello_world -> helloWorld, my_http_request -> myHTTPRequest
func SnakeToCamel(s string) string {
	if s == "" {
		return ""
	}

	var result []rune
	runes := []rune(s)
	toUpper := false

	for i, r := range runes {
		if r == '_' {
			// 遇到下划线，标记下一个字符需要大写
			toUpper = true
			continue
		}

		if toUpper {
			// 将字符转换为大写，但如果是第一个字符且需要小写驼峰，则保持原样
			if i == 1 {
				// 小写驼峰：第一个单词首字母小写
				result = append(result, unicode.ToUpper(r))
			} else {
				result = append(result, unicode.ToUpper(r))
			}
			toUpper = false
		} else {
			// 保持原样
			result = append(result, r)
		}
	}

	return string(result)
}

// SnakeToPascal 将蛇形命名转换为帕斯卡命名（大驼峰）
// hello_world -> HelloWorld, my_http_request -> MyHTTPRequest
func SnakeToPascal(s string) string {
	if s == "" {
		return ""
	}

	var result []rune
	runes := []rune(s)
	toUpper := true // 帕斯卡命名第一个字符就大写

	for _, r := range runes {
		if r == '_' {
			toUpper = true
			continue
		}

		if toUpper {
			result = append(result, unicode.ToUpper(r))
			toUpper = false
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}
