package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/LiZeC123/gmh/util"
	"github.com/urfave/cli/v3"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// 正则表达式匹配以 go-my-test: 开头的注释行
var goTestPattern = regexp.MustCompile(`(?m)^\s*//\s*go-my-test:\s*(.*)$`)

const description = `
Run with environment variables:
	// go-my-test: namespace=Test go test -run TestXXX

Run in debug mode:
	// go-my-test: go test -gcflags="all=-l -N" -run TestXXX

Run with protobuf conflictPolicy:
	// go-my-test: go test -ldflags " -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn" -run TestXXX   
`

func RunTestCommand() *cli.Command {
	return &cli.Command{
		Name:        "test",
		Usage:       "Run all tests with the `go-my-test:` tag.",
		Description: description,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "concurrent",
				Aliases: []string{"c"},
				Value:   false,
				Usage:   "Run test concurrently",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "o.txt",
				Usage:   "Output file",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			// 获取当前目录
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("error getting current directory: %v", err)
			}

			// 获取所有测试命令
			commands := walkTestFile(currentDir)
			if len(commands) == 0 {
				return nil
			}

			outputFile := c.String("output")
			file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("error opening output file: %v", err)
			}
			defer util.CloseWithLog(file)

			// 执行命令
			concurrent := c.Bool("concurrent")
			if concurrent {
				for _, c := range commands {
					command := c
					go execCmd(command, file)
				}
			} else {
				for _, cmd := range commands {
					execCmd(cmd, file)
				}
			}

			return nil
		},
	}
}

func walkTestFile(root string) []string {
	// 查找当前目录下的所有.go文件
	rst := make([]string, 0)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			// 读取文件并查找匹配的注释行
			commands, err := findComments(path)
			if err != nil {
				util.PrintErrorLog("Error finding comments in file: %v", err)
				// continue process next file
			}
			rst = append(rst, commands...)
		}
		return nil
	})

	if err != nil {
		util.PrintErrorLog("Error walking through directory: %v", err)
	}

	return rst
}

func findComments(filePath string) ([]string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer util.CloseWithLog(file)

	// 逐行读取文件
	rst := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := goTestPattern.FindStringSubmatch(line); matches != nil {
			rst = append(rst, matches[1])
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return rst, nil
}

func execCmd(command string, outFile *os.File) {
	// 创建一个执行shell命令的Cmd结构体
	cmd := exec.Command("sh", "-c", command)

	// 获取命令输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		util.PrintErrorLog("executing command: %v with error: %v\nOutput: %v\n\n", command, err, string(output))
	} else {
		fmt.Printf("executing command: %v Done\n", command)
		util.PrintToFile(outFile, "executing command: %v Done\nOutput: %v\n\n", command, string(output))
	}
}
