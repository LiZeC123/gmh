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
		Usage: "Split input from stdin into multiple files by line count",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Required: true,
				Usage:    "Output file pattern",
			},
			&cli.IntFlag{
				Name:     "line",
				Aliases:  []string{"l"},
				Required: true,
				Usage:    "Number of lines per output file",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			outputPrefix := c.String("output")
			line := c.Int("line")

			return splitFileByLines(outputPrefix, line)
		},
	}
}
func splitFileByLines(outputPrefix string, linesPerFile int) (err error) {
	// 打开输入文件
	scanner := bufio.NewScanner(os.Stdin)
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
			if currentWriter != nil && currentFile != nil {
				flushAndClose(currentWriter, currentFile)
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

		// 写入行（如果不是第一行，先添加换行符）
		if lineCount%linesPerFile > 0 {
			_, err = currentWriter.WriteString("\n")
			if err != nil {
				return fmt.Errorf("failed to write newline: %w", err)
			}
		}

		_, err = currentWriter.WriteString(lineContent)
		if err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}
		lineCount++
	}

	// 刷新并关闭最后一个文件
	if currentWriter != nil && currentFile != nil {
		flushAndClose(currentWriter, currentFile)
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	fmt.Printf("Split %d lines into %d files\n", lineCount, fileCount-1)
	return nil
}

func flushAndClose(w *bufio.Writer, file *os.File) {
	err := w.Flush()
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}
}
