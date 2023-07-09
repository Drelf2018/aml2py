package main

import (
	"fmt"

	"github.com/Drelf2018/aml2py"
	aml "github.com/Drelf2018/go-api-markup-language"
)

func main() {
	am := aml.NewParser("./tests/user.aml").Parse()
	api, file := aml2py.ToPython(am, "user.json")
	fmt.Printf("api: %v\n", api)
	fmt.Printf("file: %v\n", file)
}
