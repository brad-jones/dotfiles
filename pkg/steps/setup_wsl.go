package steps

import (
	"github.com/brad-jones/dotfiles/pkg/survey"
	"github.com/brad-jones/dotfiles/pkg/tools"
	"github.com/brad-jones/dotfiles/pkg/tools/wsl"
	"github.com/brad-jones/goerr/v2"
)

func SetupWSL(answers *survey.Answers) (err error) {
	defer goerr.Handle(func(e error) { err = e })
	wsl.MustInstall(answers.Reset)
	v := tools.GetVersion("fedora")
	name := wsl.MustInstallFedora(v.No, v.Hash, true, answers.Reset)
	wsl.MustInstallDotfiles(name, answers)
	return
}
