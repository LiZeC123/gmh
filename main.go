package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/LiZeC123/gmh/cmd"
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
				Usage: "Send an HTTP request",
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

					return cmd.DoCurl(url)
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
