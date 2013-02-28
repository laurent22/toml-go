package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
)

// Automatically build the `Document::GetXXX()` functions

func main() {
	types := [...]string{"[]Value", "string", "int", "int8", "int16", "int32", "int64", "float", "float32", "float64", "bool", "time.Time"}
	defaults := [...]string{"make([]Value, 0)", "\"\"", "0", "0", "0", "0", "0", "0.0", "0.0", "0.0", "false", "time.Now()"}
	
	output := ""
	
	for i := 0; i < len(types); i++ {
		typeName := types[i]
		typeTitle := strings.Title(types[i])
		
		if typeName == "[]Value" {
			typeTitle = "Array"	
		}
		
		if (typeName == "time.Time") {
			typeTitle = "Date"	
		}
		
		if (typeName == "float") {
			typeName = "float64"
		}
		
		s := ""
		s += "func (this Document) Get" + typeTitle + "(name string, defaultValue..." + typeName + ") " + typeName + " {\n"
		s += "\tv, ok := this.GetValue(name)\n"
		s += "\tif !ok {\n"
		s += "\t\tif len(defaultValue) >= 1 {\n"
		s += "\t\t\treturn defaultValue[0]\n"
		s += "\t\t} else {\n"
		s += "\t\t\treturn " + defaults[i] + "\n"
		s += "\t\t}\n"
		s += "\t}\n"
		s += "\treturn v.As" + typeTitle + "()\n"
		s += "}\n"
		output += s + "\n"
	}
	
	header := ""
	header += "package toml\n\n"
	header += "import (\n"
	header += "\t\"time\"\n"
	header += ")\n\n"

	output = header + output
	
	os.Remove("./toml/accessors.go")
	err := ioutil.WriteFile("./toml/accessors.go", []byte(output), os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	
	fmt.Println(output)
}