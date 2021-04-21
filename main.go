package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/brad-jones/dotfiles/pkg/assets"
	"github.com/brad-jones/dotfiles/pkg/steps"
	"github.com/brad-jones/dotfiles/pkg/survey"
	"github.com/brad-jones/dotfiles/pkg/tools/scoop"
	"github.com/brad-jones/dotfiles/pkg/tools/winsudo"
	"github.com/brad-jones/dotfiles/pkg/updater"
	"github.com/brad-jones/goasync/v2/await"
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
	// All unhandled errors should bubble all the way up to here
	// where we will spit out a stack trace for debugging purposes.
	defer goerr.Handle(func(err error) {
		goerr.PrintTrace(err)
		go func() { time.Sleep(time.Second * 120); os.Exit(1) }()
		fmt.Println("Press Enter to continue...")
		fmt.Println("NOTE: the program will terminate in 120s regardless")
		fmt.Scanln()
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

	// Define the main CLI application.
	// We are using: https://github.com/urfave/cli
	goerr.Check((&cli.App{
		Version: versionNo,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "reset",
				Usage: "If set then the setup process will start from the " +
					"very start even if it has already successfully run.",
			},
			&cli.BoolFlag{
				Name: "update",
				Usage: "If set then we will attempt to update ourselves, " +
					"before then running our updated self.",
			},
		},
		Action: func(c *cli.Context) (err error) {
			defer goerr.Handle(func(e error) { err = e })

			// Start by collecting some information from the user
			// On subsequent executions these answers will be filled
			// automatically from cache and/or the unlocked secrets vault.
			answers := survey.AskQuestions(c)

			// If the update flag has been provided then we will execute the
			// self update process. This will download a new version of this
			// tool and then execute the new version of the tool. If the new
			// version fails, a rollback will be performed.
			if answers.Update {
				updater.MustUpdate(versionNo, answers)
				os.Exit(0)
			}

			// Windows systems need a way to "elevate" & install packages.
			// Sudo for Windows: https://github.com/gerardog/gsudo
			// A command-line installer for Windows: https://scoop.sh
			if runtime.GOOS == "windows" {
				winsudo.MustInstall(answers.Reset)
				scoop.MustInstall(answers.Reset, false)
			}

			// The very next thing we do is unlock our secrets.
			// We do this early on so that we may refer to secrets
			// as required throughout the rest of the process.
			steps.MustUnlockVault(answers)
			steps.MustUnlockKeys(answers)

			await.MustFastAllOrError(
				// Update (or install) all our other software
				steps.UpdateAsync(),

				// These scripts automate some of my daily tasks.
				//
				// They either do very specific things for "me" and so make less
				// sense to release as standalone tools or they are simply at a
				// PoC stage and not ready for public consumption.
				//
				// TODO: I want to either re-write into Deno scripts so that
				// they are truly self contained scripts and then Deno handles
				// the installation of the dependencies for me. Or re-write with
				// into standalone Go tools if warrented.
				steps.InstallDartScriptsAsync(),

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

			// When run at logon we want to wait here for a bit so it gives the
			// user a chance to read what was printed to the console before the
			// console window closes.
			time.Sleep(time.Second * 3)

			return
		},
	}).Run(os.Args))
}
