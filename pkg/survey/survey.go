package survey

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goerr/v2"
	"github.com/urfave/cli/v2"
)

type Answers struct {
	GithubPassword   string
	GitlabPassword   string
	SudoPassword     string
	VaultKeyPassword string
	Reset            bool
	Update           bool
	UpdateToVersion  string
}

func AskQuestions(c *cli.Context) *Answers {
	answers := &Answers{}

	// Some answers come from the CLI & environment
	answers.Reset = c.Bool("reset")
	answers.Update = c.Bool("update")
	if answers.Update {
		answers.UpdateToVersion = os.Getenv("DOTFILES_VERSION")
	}

	// If we haven't been told to do a full reset and the things needed for the
	// vault are already installed, we are just going to assume that the vault
	// needs unlocking. The other answers we can get from the vault when we unlock it.
	if !answers.Reset && utils.CommandExists("gpg") && utils.CommandExists("gopass") {
		answers.VaultKeyPassword = utils.GetSecretFromKeychain("passphrase", "vault")
		return answers
	}

	// If we have been piped data, read it
	// Expects a single JSON line.
	if !utils.IsATerminal() {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		j := scanner.Bytes()
		goerr.Check(scanner.Err(), "failed to read JSON answers")
		goerr.Check(json.Unmarshal(j, &answers), "failed to parse JSON answers", string(j))
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
