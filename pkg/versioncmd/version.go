package versioncmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

var version = "0.0.1"

func New() *ffcli.Command {
	flagSet := flag.NewFlagSet("quokka version", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "quokka version",
		ShortHelp:  "Get the current version of Quokka",
		FlagSet:    flagSet,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				return errors.New("command does not accept args")
			}
			fmt.Fprintf(os.Stdout, "\nQuokka v%s\n", version)

			return nil
		},
	}
}
