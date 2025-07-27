package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/LiZeC123/go-my-http/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:      "server",
				Usage:     "启动HTTP服务",
				UsageText: "启动HTTP服务并打印完整的HTTP报文",
				Aliases:   []string{"s"},
				Flags: []cli.Flag{
					&cli.Int16Flag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "监听端口",
						Value:   8080,
						Required: false,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					cmd.StartServer(c.Int16("port"))
					return nil
				},
			},
			{
				Name: "curl",
				Usage: "发送HTTP请求",
				UsageText: "发送HTTP请求",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "url",
						Aliases: []string{"u"},
						Usage: "待探测的URL",
						Required: true,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					url := c.String("url")
					if url == "" {
						return errors.New("url is empty")
					}

					cmd.DoCurl(url)
					return nil
				},
			},
		},
		Authors: []any{"LiZeC"},
	}



	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
