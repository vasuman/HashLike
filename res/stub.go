package res

import (
	"html/template"

	"github.com/vasuman/HashLike/pow"
)

//go:generate go-res-pack ./data/ res_gen.go

func Setup() {
	// setup helper functions
	sysDesc := func(sys string) string {
		return pow.Desc[sys]
	}
	yesOrNo := func(b bool) string {
		if b {
			return "Yes"
		}
		return "No"
	}
	Template.Funcs(template.FuncMap{
		"sysDesc": sysDesc,
		"yesOrNo": yesOrNo,
	})
	genInit()
}
