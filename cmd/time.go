package cmd

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"context"
	"fmt"
	"github.com/LiZeC123/gmh/util"
	"github.com/urfave/cli/v3"
)


func TimeConvertCommand() *cli.Command {
	return &cli.Command{
		Name:  "time",
		Usage: "Time and timestamp conversion tool",
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:      "time",
				Min:       0,
				Max:       1,
				UsageText: "Time or timestamp to convert. If omitted, reads from standard input",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			inputs, err := util.GetArgsOrStdinInput(c, "time")
			if err != nil {
				return err
			}

			if len(inputs) == 0 {
				return errors.New("no input provided")
			}

			input := inputs[0]
			if input == "now" {
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
				return nil
			}

			result, err := ParseTimeOrTimestamp(input)
			if err != nil {
				return err
			}
			fmt.Println(result)

			return nil
		},
	}
}



func ParseTimeOrTimestamp(input string) (string, error) {
	// 移除首尾空格
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("input is empty")
	}

	// 检查是否为纯数字
	isNumeric := true
	for _, r := range input {
		if r < '0' || r > '9' {
			isNumeric = false
			break
		}
	}

	if isNumeric {
		// 将字符串转换为int64
		value, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return "", errors.New("invalid timestamp")
		}

		// 判断时间戳单位（秒、毫秒、微秒、纳秒）
		var t time.Time
		switch {
		case value < 0:
			return "", errors.New("timestamp cannot be negative")
		case value < 1e10: // 10位以内，认为是秒
			t = time.Unix(value, 0)
		case value < 1e13: // 13位以内，认为是毫秒
			t = time.UnixMilli(value)
		case value < 1e16: // 16位以内，认为是微秒
			t = time.UnixMicro(value)
		default: // 19位以内，认为是纳秒
			t = time.Unix(0, value)
		}

		// 检查时间是否有效
		if t.IsZero() || t.Unix() <= 0 {
			return "", errors.New("invalid timestamp")
		}

		// 返回格式化后的时间
		return t.Format("2006-01-02 15:04:05"), nil
	} else {
		// 尝试多种常见时间格式
		formats := []string{
			"2006-01-02 15:04:05",
			"2006/01/02 15:04:05",
			"2006-01-02T15:04:05Z07:00", // RFC3339格式
			"2006-01-02",
			"2006/01/02",
			"2006-01-02 15:04",
			"2006/01/02 15:04",
			"15:04:05 2006-01-02",
		}

		var t time.Time
		var err error
		
		for _, format := range formats {
			t, err = time.Parse(format, input)
			if err == nil {
				break
			}
		}

		if err != nil {
			// 如果所有格式都失败，尝试解析为RFC1123等格式
			t, err = time.Parse(time.RFC1123, input)
			if err != nil {
				t, err = time.Parse(time.RFC1123Z, input)
				if err != nil {
					return "", errors.New("unsupported time format")
				}
			}
		}

		// 检查时间是否有效
		if t.IsZero() || t.Unix() <= 0 {
			return "", errors.New("invalid time")
		}

		// 返回Unix时间戳（秒）
		return strconv.FormatInt(t.Unix(), 10), nil
	}
}