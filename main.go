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

		// When this is executed as a login script, this just gives the
		// user a chance to read the error before the Window closes.
		fmt.Println("Press Enter to continue...")
		fmt.Println("NOTE: the program will terminate in 30s regardless")
		go func() { time.Sleep(time.Second * 30); os.Exit(1) }()
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
				if updater.MustUpdate(versionNo, answers) {
					// If we did update successfully we can exit here as the
					// updater runs the new version of this tool for us, in
					// effect making the rest of the instructions in this
					// function redundant and possibly incompatible with the
					// what the new version of the tool just executed.
					time.Sleep(time.Second * 3)
					os.Exit(0)
				}
				// However if the updater fails we can then run this version
				// of the tool in the hope that it will leave the system in a
				// useable state, in effect performing a rollback.
			}

			// There are small selection of configuration files that are
			// required to unlock the secrets vault & maybe others in the
			// future (eg: scoop config say). Most dotfiles are installed
			// after the vault in unlocked so that they can consume secrets.
			assets.MustWriteDotfiles(false)

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

			// This is the whole point right :)
			// Lets write our actual dotfiles to the filesystem.
			assets.MustWriteDotfiles(true)

			// Here we install (or update) most of the rest of our software
			steps.MustInstallOrUpdate(answers)
			await.MustFastAllOrError(
				steps.InstallDartScriptsAsync(answers),
				steps.InstallRunAtLogonScriptAsync(),
			)

			// Finally, on Windows systems we perform some recursion and we
			// create a WSL instance in which we then execute ourselves or
			// rather the linux version, embedded with-in, inside the new
			// WSL instance.
			if runtime.GOOS == "windows" {
				steps.SetupWSL(answers)
			}

			// When run at logon we want to wait here for a bit so it gives the
			// user a chance to read what was printed to the console before the
			// console window closes.
			time.Sleep(time.Second * 3)

			return
		},
	}).Run(os.Args))
}
