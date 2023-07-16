package aml2py

import (
	_ "embed"
	"fmt"
	"strings"

	aml "github.com/Drelf2018/go-api-markup-language"
	"github.com/Drelf2020/utils"
)

// 转 value 为 python 格式
func ValueToPython(typ, val string) string {
	val = strings.Replace(val, ",constant", "", 1)
	if typ == "str" {
		if val == "none" {
			return "\"\""
		}
		return "\"" + val + "\""
	} else if typ == "bool" {
		if val == "none" {
			return "False"
		}
		return utils.Capitalize(val)
	} else {
		if val == "none" {
			return "0"
		}
		return val
	}
}

// 转为 Python 格式
func SentenceToPython(sentence *aml.Sentence) (s string) {
	s = sentence.Name
	typ := strings.Replace(sentence.Type, "num", "int", 1)
	val := sentence.Value
	if typ != "" {
		s += ": " + typ
	}

	if val != "" {
		s += " = " + ValueToPython(typ, val)
	}
	return
}

func ToPythonFunc(format, name string, api *aml.Api) string {
	r, o, all := []string{}, []string{}, []string{}
	hint := api.Hint
	params := []string{}
	for _, sentence := range api.Params.Map {
		params = append(params, fmt.Sprintf("%v (%v): %v", sentence.Name, sentence.Type, sentence.Hint))
		if sentence.IsConstant() {
			typ := strings.Replace(sentence.Type, "num", "int", 1)
			all = append(all, sentence.Name+"="+ValueToPython(typ, sentence.Value))
			continue
		}
		if sentence.IsRequired() {
			r = append(r, SentenceToPython(sentence))
		} else if sentence.IsOptional() {
			o = append(o, SentenceToPython(sentence))
		}
		all = append(all, sentence.Name+"="+sentence.Name)
	}
	if len(api.Params.Map) != 0 {
		hint += "\n\n    Args:\n        " + strings.Join(params, "\n\n        ")
	}
	args := append(r, o...)

	format = strings.ReplaceAll(format, "demo", name)
	format = strings.ReplaceAll(format, "args", strings.Join(args, ", "))
	format = strings.ReplaceAll(format, "hint", hint)
	format = strings.ReplaceAll(format, "update(", "update("+strings.Join(all, ", "))
	return format
}

//go:embed func.py
var FUNC string

//go:embed api.py
var API string

func init() {
	aml.Plugin{
		Cmd:         "python",
		Author:      "Drelf2018",
		Version:     "1.2.0",
		Description: "将 aml 文件转为可执行 python 代码",
		Link:        "https://github.com/Drelf2018/aml2py",
		Generate: func(p *aml.Parser) (files []aml.File) {
			include := FUNC[:strings.Index(FUNC, "# loop")]
			include = strings.Replace(include, "path", p.NewExt(".json"), 1)
			function := utils.Slice(FUNC, "# loop", "# end", 0)
			for name, api := range p.Output {
				include += ToPythonFunc(function, name, api)
			}
			return []aml.File{
				{Name: "api.py", Content: API},
				{Name: p.NewExt(".py"), Content: include},
			}
		},
	}.Load()
}
