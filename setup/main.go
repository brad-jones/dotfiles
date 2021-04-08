package main

import (
	"os"

	"github.com/brad-jones/dotfiles/setup/tasks"
	"github.com/brad-jones/goerr/v2"
	"github.com/urfave/cli/v2"
)

func main() {
	defer goerr.Handle(func(err error) {
		goerr.PrintTrace(err)
		os.Exit(1)
	})

	goerr.Check((&cli.App{
		Action: func(c *cli.Context) (err error) {
			defer goerr.Handle(func(e error) { err = e })

			if c.Args().Get(0) == "chezmoi-apply" {
				goerr.Check(tasks.ChezmoiApply())
				return
			}

			goerr.Check(tasks.Bootstrap())
			return
		},
	}).Run(os.Args))
}

/*
	TODO

	install xop

	Fix Dart warnings

	Fix aws-vault rotate

	The rdp shortcuts

	Docker Desktop

	WSL setup, which then leads into Linux setup

	Vscode extensions sync
*/
