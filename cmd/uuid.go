package cmd

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/urfave/cli/v3"
)

func UUIDCommand() *cli.Command {
	return &cli.Command{
		Name:  "uuid",
		Usage: "Generate a UUID",
		Action: func(ctx context.Context, c *cli.Command) error {
			return UUID()
		},
	}
}

func UUID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	fmt.Println(id.String())

	return nil
}
