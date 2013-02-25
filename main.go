package main

import (
	"./toml"
	"fmt"
)

func main() {
	var test toml.Parser
	doc := test.ParseFile("example.toml")
	
	var value toml.Value
	var ok bool
	
	// Get a string:
	
	value, ok = doc.GetValue("servers.beta.ip")
	if !ok { // Optionally, you can check if the value exists or not
		panic("value doesn't exists")
	} else {
		fmt.Println(value)
	}
	
	// Get an array
	
	value, _ = doc.GetValue("clients.data")
	array := value.AsArray()
	fmt.Println(array[0].AsArray()[0].AsString())
	
	// Get an int
	
	value, _ = doc.GetValue("database.connection_max")
	fmt.Println(value.AsInt())
}