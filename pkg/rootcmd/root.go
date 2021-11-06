package rootcmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func New() *ffcli.Command {
	var (
		command     = "quokka"
		rootFlagSet = flag.NewFlagSet(command, flag.ExitOnError)
		usage       = fmt.Sprintf("%s [flags] [<arg>...]", command)
	)

	return &ffcli.Command{
		Name:       command,
		ShortUsage: usage,
		FlagSet:    rootFlagSet,
		Exec: func(ctx context.Context, args []string) error {
			// root command serves no purpose other
			// than displaying the usage text to the user.
			return flag.ErrHelp
		},
	}
}
