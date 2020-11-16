package steps

import (
	"fmt"
	"os"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/brad-jones/goerr/v2"
)

// BootstrapSurveyAnswers contains answers from the survey
type BootstrapSurveyAnswers struct {
	GithubPassword   string
	GitlabPassword   string
	VaultKeyPassword string
	SudoPassword     string
}

// BootstrapSurvey will prompt the user for answers that we need to know ahead of time.
func BootstrapSurvey() *BootstrapSurveyAnswers {
	answers := &BootstrapSurveyAnswers{}
	questions := []*survey.Question{
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
	}

	if runtime.GOOS == "linux" {
		questions = append(questions, &survey.Question{
			Name:   "SudoPassword",
			Prompt: &survey.Password{Message: "Your password for sudo?"},
		})
	}

	if len(questions) > 0 {
		err := survey.Ask(questions, answers)
		if err != nil {
			if err == terminal.InterruptErr {
				fmt.Println("interrupted")
				os.Exit(0)
			}
			goerr.Check(err)
		}
		fmt.Println("")
	}

	return answers
}
