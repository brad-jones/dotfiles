package main

import (
	"os"

	"github.com/brad-jones/dotfiles/setup/tasks"
	"github.com/brad-jones/goerr/v2"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := (&cli.App{
		Action: func(c *cli.Context) error {
			if c.Args().Get(0) == "chezmoi-apply" {
				if err := tasks.ChezmoiApply(); err != nil {
					return goerr.Wrap(err)
				}
			} else {
				if err := tasks.Bootstrap(); err != nil {
					return goerr.Wrap(err)
				}
			}
			return nil
		},
	}).Run(os.Args); err != nil {
		goerr.PrintTrace(err)
		os.Exit(1)
	}
}
