package main

import (
	"context"
	"log"
	"os"

	"github.com/LiZeC123/gmh/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	c := &cli.Command{
		Usage:                 "All-in-one HTTP utility toolkit",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			cmd.ServerCommand(),
			cmd.CurlCommand(),
			cmd.DNSCommand(),
			cmd.TcpingCommand(),
			cmd.UUIDCommand(),
			cmd.JsonCommand(),
			cmd.MemCommand(),
			cmd.StringCount(),
			cmd.SplitFileCommand(),
			cmd.UniqueCommand(),
			cmd.StringList(),
			cmd.Replace(),
			cmd.CaseCommand(),
		},

		Authors: []any{"LiZeC"},
	}

	if err := c.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
