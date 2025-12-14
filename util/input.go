package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

func GetAllInput(c *cli.Command, argsName, fileName string) (lines []string, err error) {
	lines = c.StringArgs(argsName)
	inputFile := c.String(fileName)
	if inputFile != "" {
		fileUrls, err := GetFileInput(inputFile)
		if err != nil {
			return nil, err
		}

		lines = append(lines, fileUrls...)
	}

	return
}

func GetFileInput(filePath string) (lines []string, err error) {
	var reader io.Reader

	if filePath == "-" {
		// 从标准输入读取
		reader = os.Stdin
	} else {
		// 从文件读取
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer CloseWithLog(file)
		reader = file
	}

	// 读取
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return
}

func GetArgsOrStdinInput(c *cli.Command, argsName string) (lines []string, err error) {
	// 传入命令行参数, 则直接返回命令行参数
	inputs := c.StringArgs(argsName)
	if len(inputs) > 0 {
		return inputs, nil
	}

	// 否则从读取标准输入
	lines, err = GetFileInput("-")
	if err != nil {
		return []string{}, err
	}
	return lines, nil
}
