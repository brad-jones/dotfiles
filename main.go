package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/brad-jones/dotfiles/pkg/assets"
	"github.com/brad-jones/dotfiles/pkg/steps"
	"github.com/brad-jones/dotfiles/pkg/survey"
	"github.com/brad-jones/dotfiles/pkg/tools/scoop"
	"github.com/brad-jones/dotfiles/pkg/tools/winsudo"
	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goasync/v2/task"
	"github.com/brad-jones/goerr/v2"
	"github.com/urfave/cli/v2"
)

// Injected by ldflags
// Makefile release target does this
// see: https://stackoverflow.com/questions/11354518
var (
	versionNo = "0.0.0"
	commitUrl = "dev"
	buildDate = "unknown"
)

func main() {
	// All un handled errors should bubble all the way up to here
	// where we will spit out a stack trace for debugging purposes.
	defer goerr.Handle(func(err error) {
		goerr.PrintTrace(err)
		os.Exit(1)
	})

	// Output the ldflag values when the "--version" flag is supplied
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("versionNo=%s\ncommitUrl=%s\nbuildDate=%s",
			c.App.Version,
			commitUrl,
			buildDate,
		)
	}

	goerr.Check((&cli.App{
		Version: versionNo,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "reset",
				Usage: "If set then many things will be deleted & replaced instead of just checking for existence.",
			},
		},
		Action: func(c *cli.Context) (err error) {
			defer goerr.Handle(func(e error) { err = e })

			// Start by collecting some information from the user
			// On subsequent executions these answers will be filled
			// automatically from cache or the unlocked secrets vault.
			answers := survey.AskQuestions(c)

			if runtime.GOOS == "windows" {
				// Before we go any further we need a way to "elevate" on Windows.
				winsudo.MustInstall(answers.Reset)

				// And then we need the scoop package manager
				// to be able to install all the other things.
				scoop.MustInstall(answers.Reset, false)
			}

			await.MustFastAllOrError(
				// Update (or install) all our other software
				steps.UpdateAsync(),

				// Setup our dart scripts
				steps.InstallDartScriptsAsync(),

				// Unlock our secrets
				task.New(func() {
					steps.MustUnlockVault(answers)
					steps.MustUnlockKeys(answers)
				}),

				// This will make this binary self-update and run again on logon
				steps.InstallRunAtLogonScriptAsync(),
			)

			// Write all out other files
			if runtime.GOOS == "windows" {
				assets.WriteFolderToHome("AppData/Local/Microsoft/Windows Terminal")
				assets.WriteFolderToHome("AppData/Roaming/Code")
				assets.WriteFolderToHome("Documents")
			}
			assets.WriteFolderToHome(".aws")
			assets.WriteFolderToHome("Projects")
			assets.WriteFileToHome(".gitconfig")

			return
		},
	}).Run(os.Args))
}
