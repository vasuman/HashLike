package res

import (
	"html/template"

	"github.com/vasuman/HashLike/pow"
)

//go:generate go run tools/gen.go ./data/

func Setup() {
	sysDesc := func(sys string) string {
		return pow.Desc[sys]
	}
	yesOrNo := func(b bool) string {
		if b {
			return "Yes"
		}
		return "No"
	}
	// setup helper functions
	Template.Funcs(template.FuncMap{
		"sysDesc": sysDesc,
		"yesOrNo": yesOrNo,
	})
	genInit()
}
