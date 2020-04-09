package survey

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/brad-jones/dotfiles/pkg/tools/gopass"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goerr/v2"
	"github.com/gosimple/slug"
	"github.com/urfave/cli/v2"
)

type Answers struct {
	GithubPassword   string
	GitlabPassword   string
	SudoPassword     string
	VaultKeyPassword string
	Reset            bool
}

func AskQuestions(c *cli.Context) *Answers {
	answers := &Answers{}

	// Some answers come from the CLI
	answers.Reset = c.Bool("reset")

	// If our secret vault (currently gopass) is unlocked
	// we will get these secrets from the vault.
	if gopass.VaultUnlocked() {
		answers.GithubPassword = gopass.GetSecret("websites/github.com/brad@bjc.id.au")
		answers.GitlabPassword = gopass.GetSecret("websites/gitlab.com/brad@bjc.id.au")
		answers.SudoPassword = gopass.GetSecret(fmt.Sprintf("sudo/%s", slug.Make(utils.GetComputerName())))
		answers.VaultKeyPassword = utils.GetSecretFromKeychain("passphrase", "vault")
		return answers
	}

	// If we have been piped data, read it
	if !utils.IsATerminal() {
		scanner := bufio.NewScanner(os.Stdin)

		scanner.Scan()
		answers.GithubPassword = scanner.Text()
		goerr.Check(scanner.Err())

		scanner.Scan()
		answers.GitlabPassword = scanner.Text()
		goerr.Check(scanner.Err())

		scanner.Scan()
		answers.SudoPassword = scanner.Text()
		goerr.Check(scanner.Err())

		scanner.Scan()
		answers.VaultKeyPassword = scanner.Text()
		goerr.Check(scanner.Err())

		return answers
	}

	// Otherwise interactively ask the user for this information.
	if err := survey.Ask([]*survey.Question{
		{
			Name:   "GithubPassword",
			Prompt: &survey.Password{Message: "The password used to clone repos from github?"},
		},
		{
			Name:   "GitlabPassword",
			Prompt: &survey.Password{Message: "The password used to clone repos from gitlab?"},
		},
		{
			Name:   "VaultKeyPassword",
			Prompt: &survey.Password{Message: "The password used to unlock the gopass vault gpg key?"},
		},
		{
			Name:   "SudoPassword",
			Prompt: &survey.Password{Message: "Your password for sudo? (on Windows this will be used for WSL)"},
		},
	}, answers); err != nil {
		if err == terminal.InterruptErr {
			fmt.Println("interrupted")
			os.Exit(0)
		}
		goerr.Check(err)
	}
	fmt.Println("")

	return answers
}
