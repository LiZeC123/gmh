package cmd

import (
	"context"
	"os"
	"text/tabwriter"

	"github.com/LiZeC123/gmh/util"

	"github.com/urfave/cli/v3"
)

type CounterInfo struct {
	Name            string
	Total           int
	EnglishCount    int
	NonEnglishCount int
}

func StringCount() *cli.Command {
	return &cli.Command{
		Name:  "strcnt",
		Usage: "Count string length, English and non-English characters",
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:      "str",
				Min:       0,
				Max:       -1,
				UsageText: "Input strings (if omitted, reads from stdin)",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			inputs, err := util.GetArgsOrStdinInput(c, "str")
			if err != nil {
				return err
			}

			countInfos := make([]*CounterInfo, 0, len(inputs))
			for _, input := range inputs {
				countInfos = append(countInfos, countOne(input))
			}

			return doPrint(countInfos)
		},
	}
}

func countOne(input string) *CounterInfo {
	total := len([]rune(input))
	englishCount := 0

	for _, char := range input {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			englishCount++
		}
	}

	nonEnglishCount := total - englishCount

	return &CounterInfo{
		Name:            truncateSubString(input, 32),
		Total:           total,
		EnglishCount:    englishCount,
		NonEnglishCount: nonEnglishCount,
	}
}

func truncateSubString(s string, m int) string {
	// 异常输入返回空字符串
	if m <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) > m {
		return string(runes[:m])
	}
	return s
}

func doPrint(info []*CounterInfo) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)

	util.PrintToFile(w, "Count\tEnglish\tNon-English\tInput-Content\t")

	for _, i := range info {
		util.PrintToFile(w, "%d\t%d\t%d\t%s\t\n", i.Total, i.EnglishCount, i.NonEnglishCount, i.Name)
	}

	// 刷新缓冲区，确保输出
	err := w.Flush()
	if err != nil {
		return err
	}

	return nil

}
