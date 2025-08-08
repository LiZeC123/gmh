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
		Usage: "GO语言实现的HTTP实用工具",
		Commands: []*cli.Command{
			{
				Name:      "server",
				Usage:     "启动Echo HTTP服务",
				UsageText: "启动一个HTTP服务 此服务打印并返回完整的HTTP报文",
				Aliases:   []string{"s"},
				Flags: []cli.Flag{
					&cli.Int16Flag{
						Name:     "port",
						Aliases:  []string{"p"},
						Usage:    "监听端口",
						Value:    8080,
						Required: false,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					cmd.StartServer(c.Int16("port"))
					return nil
				},
			},
			{
				Name:      "curl",
				Usage:     "发送HTTP请求",
				UsageText: "发送HTTP请求",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "url",
						Aliases:  []string{"u"},
						Usage:    "待探测的URL",
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
			{
				Name:      "dns",
				Usage:     "执行DNS解析",
				UsageText: "执行DNS解析",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "url",
						Aliases:  []string{"u"},
						Usage:    "待解析的URL",
						Required: true,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					url := c.String("url")
					if url == "" {
						return errors.New("url is empty")
					}

					cmd.DoDNS(url)
					return nil
				},
			},
			{
				Name:      "tcping",
				Usage:     "TCP端口探测",
				UsageText: "探测指定的目标能否建立TCP链接",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "host",
						Aliases:  []string{"h"},
						Usage:    "目标主机",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "port",
						Aliases:  []string{"p"},
						Usage:    "目标端口",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "timeout",
						Aliases:  []string{"t"},
						Usage:    "最大超时时间",
						Value:    1,
						Required: false,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					host := c.String("host")
					port := c.String("port")
					timeout := c.Int("timeout")
					cmd.Tcping(host, port, timeout)
					return nil
				},
			},
			{
				Name: "uuid",
				Usage: "生成UUID",
				UsageText: "生成一个随机UUID",
				Action: func(ctx context.Context, c *cli.Command) error {
					return cmd.UUID()
				},
			},
		},
		Authors: []any{"LiZeC"},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
