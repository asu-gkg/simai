package main

import (
	"fmt"
	"os"
	"simai/param_parse"
)

func main() {
	fmt.Println("SimAi starts!")
	param := param_parse.NewUserParam()
	if err := param.Parse(os.Args[1:]); err != nil {
		panic("parse err")
	}
	
}
