package main

import (
	"./toml"
	"fmt"
)

func main() {
	// Create a Toml parser
	var parser toml.Parser
	
	// Parse a file
	doc := parser.ParseFile("example.toml")
	
	// Or parse a string directly:
	// doc := parser.Parse(someTomlString)
	
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
	fmt.Println(value.AsArray()[0].AsArray()[0].AsString())
	
	// Get an int
	
	value, _ = doc.GetValue("database.connection_max")
	fmt.Println(value.AsInt())
	
	// Get a float
	
	value, _ = doc.GetValue("floats.pi")
	fmt.Println(value.AsFloat())
	
	// Get a negative float
	
	value, _ = doc.GetValue("floats.minus")
	fmt.Println(value.AsFloat())
	
	// Get a boolean
	
	value, _ = doc.GetValue("database.enabled")
	fmt.Println(value.AsBool())
	
	// Get a date
	
	value, _ = doc.GetValue("owner.dob")
	fmt.Println(value.AsDate())
	
	// Get a section
	
	section, _ := doc.GetSection("owner")
	fmt.Println(section)
	
	// Get title
	value, _ = doc.GetValue("title")
	fmt.Println(value.AsString())
}