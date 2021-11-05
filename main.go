package main

import (
	"context"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	root := &ffcli.Command{
		Exec: func(ctx context.Context, args []string) error {
			println("Quokka")
			return nil
		},
	}

	root.ParseAndRun(context.Background(), os.Args[1:])
}
