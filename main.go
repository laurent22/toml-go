package main

import (
	"./toml"
	"fmt"
)

func main() {
	//var test toml.Parser
	s := toml.RemoveComment("\"example # no that's not  comment\" # but here's one")
	fmt.Println(s)
	//test.ParseFile("S:/Docs/PROGS/go/toml/example.toml")
}