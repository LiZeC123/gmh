package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func SplitFileCommand() *cli.Command {
	return &cli.Command{
		Name:  "spf",
		Usage: "Split File with given separator",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Required: true,
				Usage:    "Input file to split",
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Required: true,
				Usage:    "Output file prefix",
			},
			&cli.IntFlag{
				Name:     "line",
				Aliases:  []string{"l"},
				Required: true,
				Usage:    "Number of lines per output file",
			},
			&cli.StringFlag{
				Name:     "prefix",
				Aliases:  []string{"p"},
				Required: false,
				Usage:    "Prefix to add to each line",
			},
			&cli.StringFlag{
				Name:     "suffix",
				Aliases:  []string{"s"},
				Required: false,
				Usage:    "Suffix to add to each line",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			inputFile := c.String("input")
			outputPrefix := c.String("output")
			line := c.Int("line")
			linePrefix := c.String("prefix")
			lineSuffix := c.String("suffix")

			return splitFileByLines(inputFile, outputPrefix, line, linePrefix, lineSuffix)
		},
	}
}
func splitFileByLines(inputFile, outputPrefix string, linesPerFile int, linePrefix, lineSuffix string) error {
	// 打开输入文件
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fileCount := 1
	lineCount := 0
	var currentFile *os.File
	var currentWriter *bufio.Writer

	// 确保输出目录存在
	outputDir := filepath.Dir(outputPrefix)
	if outputDir != "" && outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// 获取文件扩展名
	ext := filepath.Ext(outputPrefix)
	base := outputPrefix
	if ext != "" {
		base = strings.TrimSuffix(outputPrefix, ext)
	}

	for scanner.Scan() {
		if lineCount%linesPerFile == 0 {
			// 关闭当前文件（如果存在）
			if currentFile != nil {
				currentWriter.Flush()
				currentFile.Close()
			}

			// 创建新文件
			outputFileName := fmt.Sprintf("%s_%04d%s", base, fileCount, ext)
			currentFile, err = os.Create(outputFileName)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			currentWriter = bufio.NewWriter(currentFile)
			fileCount++
		}

		// 处理行内容（添加前缀和后缀）
		lineContent := scanner.Text()
		if linePrefix != "" {
			lineContent = linePrefix + lineContent
		}
		if lineSuffix != "" {
			lineContent = lineContent + lineSuffix
		}

		// 写入行（如果不是第一行，先添加换行符）
		if lineCount%linesPerFile > 0 {
			_, err = currentWriter.WriteString("\n")
			if err != nil {
				return fmt.Errorf("failed to write newline: %w", err)
			}
		}

		_, err := currentWriter.WriteString(lineContent)
		if err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}
		lineCount++
	}

	// 刷新并关闭最后一个文件
	if currentWriter != nil {
		currentWriter.Flush()
	}
	if currentFile != nil {
		currentFile.Close()
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	fmt.Printf("Split %d lines into %d files\n", lineCount, fileCount-1)
	return nil
}
