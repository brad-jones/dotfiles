package assets

import (
	"bytes"
	"runtime"
	"text/template"

	"github.com/brad-jones/dotfiles/pkg/tools/gopass"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goerr/v2"
)

var tmplOptions = []string{
	"missingkey=error",
}

var tmplFuncMap = template.FuncMap{
	"GetSecret": gopass.GetSecret,
}

type tmplData struct {
	Runtime struct {
		GOOS   string
		GOARCH string
		IsWSL  bool
	}
}

func buildTmplData() *tmplData {
	dat := &tmplData{}
	dat.Runtime.GOOS = runtime.GOOS
	dat.Runtime.GOARCH = runtime.GOARCH
	dat.Runtime.IsWSL = utils.IsWSL()
	return dat
}

func ExecuteTemplate(in []byte) (out []byte, err error) {
	defer goerr.Handle(func(e error) { err = e })
	var b bytes.Buffer
	t, err := template.New("tmpl").
		Option(tmplOptions...).
		Funcs(tmplFuncMap).
		Parse(string(in))
	goerr.Check(err, "failed to parse")
	goerr.Check(t.Execute(&b, buildTmplData()), "failed to execute")
	out = b.Bytes()
	return
}
