package cmd

import (
	"fmt"

	"github.com/google/uuid"
)

func UUID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	fmt.Println(id.String())

	return nil
}
