package main

import (
	"context"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/raycatso/quokka/pkg/rootcmd"
	"github.com/raycatso/quokka/pkg/versioncmd"
)

func main() {
	var (
		rootCmd    = rootcmd.New()
		versionCmd = versioncmd.New()
	)

	rootCmd.Subcommands = []*ffcli.Command{
		/* runCmd */
		versionCmd,
	}

	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error during Parse: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
