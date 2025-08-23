package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/LiZeC123/gmh/util"
	"github.com/urfave/cli/v3"
)

func CurlCommand() *cli.Command {
	return &cli.Command{
		Name:  "curl",
		Usage: "Send HTTP requests to one or multiple URLs",
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "url",
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
			&cli.BoolFlag{
				Name:     "url-only",
				Aliases:  []string{"u"},
				Required: false,
				Usage:    "Output only URLs instead of full HTTP responses",
			},
			&cli.StringFlag{
				Name:     "filter",
				Aliases:  []string{"f"},
				Value:    "all",
				Required: false,
				Usage:    "Filter output: (a)ll, (s)uccess, or (f)ailure",
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Required: false,
				Usage:    "Write results to the specified file",
			},
			&cli.BoolFlag{
				Name:     "progress",
				Aliases:  []string{"p"},
				Required: false,
				Usage:    "Show progress bar when processing multiple URLs",
			},
			&cli.Uint16Flag{
				Name:    "concurrency",
				Aliases: []string{"c"},
				Value:   200,
				Usage:   "Maximum number of concurrent requests",
			},
			&cli.Uint8Flag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Value:   10,
				Usage:   "Timeout in seconds for each request",
			},
			&cli.Uint8Flag{
				Name:    "retry",
				Aliases: []string{"r"},
				Value:   0,
				Usage:   "Number of times to retry failed requests",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			// 准备输入
			urls, err := util.GetAllInput(c, "url", "input")
			if err != nil {
				return err
			}
			if len(urls) == 0 {
				return fmt.Errorf("no URLs provided. Use command arguments, --file, or stdin")
			}

			// 准备输出
			outputFile := c.String("output")
			var writer io.Writer = os.Stdout
			if outputFile != "" {
				f, err := os.Create(outputFile)
				if err != nil {
					return err
				}
				defer f.Close()
				writer = f
			}
			showProgress := c.Bool("progress") && outputFile != ""

			// 输出参数处理
			filter := c.String("filter")
			switch filter {
			case "a", "all":
				filter = "all"
			case "s", "success":
				filter = "success"
			case "f", "failure":
				filter = "failure"
			default:
				return fmt.Errorf("invalid filter value: %s. Use (a)ll, (s)uccess, or (f)ailure", filter)
			}

			// 执行并发检测
			task := Task{
				Urls:        urls,
				Concurrency: c.Uint16("concurrency"),
				Timeout:     c.Uint8("timeout"),
				Retry:       c.Uint8("retry"),
				UrlOnly:     c.Bool("url-only"),
			}
			out := DoCurlTask(task)

			// 收集执行结果
			total := len(task.Urls)
			count := 0
			succCount := 0
			failCount := 0
			for rst := range out {
				count++
				if rst.Err == nil {
					succCount++
					if filter == "success" {
						fmt.Fprintln(writer, rst.Data)
					}
				} else {
					failCount++
					if filter == "failure" {
						fmt.Fprintln(writer, rst.Data)
					}
				}

				if filter == "all" {
					fmt.Fprintln(writer, rst.Data)
				}

				if showProgress {
					fmt.Printf("Total %d Done %d (%.2f%%): Succ: %d Fail: %d (%.2f%%)\n", total, count, 100*float32(count)/float32(total), succCount, failCount, float32(succCount)/float32(total))
				}
			}
			return nil
		},
	}
}

type Task struct {
	Urls        []string
	Concurrency uint16
	Timeout     uint8
	Retry       uint8
	UrlOnly     bool
}

type TaskRst struct {
	Data string
	Err  error
}


func DoCurlTask(task Task) (out chan TaskRst) {
	out = make(chan TaskRst)
	sem := make(chan struct{}, task.Concurrency)
	var wg sync.WaitGroup

	go func() {
		defer close(out)

		for _, url := range task.Urls {
			sem <- struct{}{}
			wg.Add(1)

			go func(u string) {
				defer func() {
					<-sem
					wg.Done()
				}()

				data, err := DoCurl(url, task.Timeout, task.Retry)

				if task.UrlOnly {
					data = url
				}

				// 详细的错误信息输出到标准错误, 可重定向到文件
				if err != nil {
					fmt.Fprintf(os.Stderr, "Curl %s failed with err: %v\n", url, err)
				}

				out <- TaskRst{
					Data: data,
					Err:  err,
				}
			}(url)
		}
		wg.Wait()
	}()

	return out

}


func DoCurl(url string, timeout uint8, retry uint8) (body string, err error) {
	req, _ := http.NewRequest("GET", url, nil)

	// 覆盖默认标识，添加浏览器特征头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	return string(bytes), nil
}
