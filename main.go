package main

import (
	"context"
	"log"
	"os"

	"github.com/LiZeC123/gmh/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Usage: "All-in-one HTTP utility toolkit",
		Commands: []*cli.Command{
			cmd.ServerCommand(),
			cmd.CurlCommand(),
			cmd.DNSCommand(),
			cmd.TcpingCommand(),
			cmd.UUIDCommand(),
			cmd.JsonCommand(),
			cmd.MemCommand(),
			cmd.JoinCommand(),
			cmd.StringCount(),
			cmd.SplitFileCommand(),
			cmd.UniqueCommand(),
		},

		Authors: []any{"LiZeC"},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
