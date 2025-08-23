package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/LiZeC123/gmh/cmd"
	"github.com/LiZeC123/gmh/util"
	"github.com/urfave/cli/v3"
)

const defaultPort = 8080
const defaultTimeout = 1

func main() {
	cmd := &cli.Command{
		Usage: "All-in-one HTTP utility toolkit",
		Commands: []*cli.Command{
			{
				Name:    "server",
				Usage:   "Start an HTTP echo server",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.Uint16Flag{
						Name:     "port",
						Aliases:  []string{"p"},
						Value:    defaultPort,
						Required: false,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return cmd.StartServer(c.Uint16("port"))
				},
			},
			{
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

					task := cmd.Task{
						Urls:        urls,
						Concurrency: c.Uint16("concurrency"),
						Timeout:     c.Uint8("timeout"),
						Retry:       c.Uint8("retry"),
						UrlOnly:     c.Bool("url-only"),
					}

					out := cmd.DoCurlTask(task)
					total := len(task.Urls)
					count := 0
					succCount := 0
					failCount := 0
					for rst := range out {
						count++
						if rst.Err == nil {
							succCount ++
							if filter == "success" {
								fmt.Fprintln(writer, rst.Data)
							}
						} else {
							failCount++
							if filter ==  "failure" {
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
			},
			{
				Name:  "dns",
				Usage: "Perform DNS lookup",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "url",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					url := c.StringArg("url")
					if url == "" {
						return errors.New("url cannot be empty")
					}

					return cmd.DoDNS(url)
				},
			},
			{
				Name:  "tcping",
				Usage: "Probe TCP port connectivity",
				Arguments: []cli.Argument{
					&cli.StringArgs{
						Name: "host",
						Min: 1,
						Max: -1,
					},
				},
				Flags: []cli.Flag{
					&cli.Uint16Flag{
						Name:     "port",
						Aliases:  []string{"p"},
						Required: true,
					},
					&cli.Uint8Flag{
						Name:     "timeout",
						Aliases:  []string{"t"},
						Value:    defaultTimeout,
						Required: false,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					hosts := c.StringArgs("host")
					
					port := c.Uint16("port")
					timeout := c.Uint8("timeout")

					for _, host := range hosts {
						err := cmd.Tcping(host, port, timeout)
						if err != nil {
							return err
						}
						
					}

					return nil
				},
			},
			{
				Name:  "uuid",
				Usage: "Generate a UUID",
				Action: func(ctx context.Context, c *cli.Command) error {
					return cmd.UUID()
				},
			},
			{
				Name:      "json",
				Usage:     "Validate and format JSON",
				UsageText: "End input: Ctrl+Z (Windows) or Ctrl+D (Linux)",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "compress",
						Aliases:  []string{"c"},
						Required: false,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {

					content, err := io.ReadAll(os.Stdin)
					if err != nil {
						return err
					}
					if len(content) == 0 {
						return errors.New("empty JSON input")
					}

					rawJSON := string(content)
					hasReplace := false
					maxIterations := 10
					for range maxIterations {
						rawJSON, hasReplace = cmd.Unescape(rawJSON)
						if !hasReplace {
							break
						}
					}

					ok := cmd.Validate(rawJSON)
					if !ok {
						return nil
					}

					compress := c.Bool("compress")
					if compress {
						return cmd.CompressJSON(rawJSON)
					} else {
						return cmd.FormatJSON(rawJSON)
					}
				},
			},
			{
				Name:  "mem",
				Usage: "Perform a memory stability test",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:     "maxMemory",
						Aliases:  []string{"m"},
						Usage:    "Maximum memory to allocate in gigabytes (GB)",
						Required: true,
					},
					&cli.UintFlag{
						Name:     "loopCount",
						Aliases:  []string{"c"},
						Usage:    "Number of test iterations",
						Required: true,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					maxMemory := c.Uint("maxMemory")
					loopCount := c.Uint("loopCount")

					return cmd.MemCheck(maxMemory, loopCount)
				},
			},
		},

		Authors: []any{"LiZeC"},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
