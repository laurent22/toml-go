package main

import (
	toml ".."
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
	
	// ==================================================
	// Get some values:
	// ==================================================
	
	fmt.Println(doc.GetString("servers.beta.ip"))
	fmt.Println(doc.GetArray("clients.data"))
	fmt.Println(doc.GetInt("doesntexist", 123)) // Optionally, a default value can be provided
	fmt.Println(doc.GetString("doesntexisteither", "some default"))
	fmt.Println(doc.GetFloat("floats.pi"))
	fmt.Println(doc.GetBool("database.enabled"))
	fmt.Println(doc.GetDate("owner.dob"))
	
	// ==================================================
	// Or using the GetValue() / As<Type>() pattern:
	// ==================================================
	
	value, ok = doc.GetValue("servers.beta.ip")
	if !ok { // Optionally, you can check if the value exists or not
		panic("value doesn't exists")
	} else {
		fmt.Println(value.AsString())
	}
	
	value, _ = doc.GetValue("clients.data")
	fmt.Println(value.AsArray()[0].AsArray()[0].AsString())
	
	value, _ = doc.GetValue("database.connection_max")
	fmt.Println(value.AsInt())
	
	value, _ = doc.GetValue("floats.pi")
	fmt.Println(value.AsFloat())
	
	value, _ = doc.GetValue("floats.minus")
	fmt.Println(value.AsFloat())
	
	value, _ = doc.GetValue("database.enabled")
	fmt.Println(value.AsBool())
	
	value, _ = doc.GetValue("owner.dob")
	fmt.Println(value.AsDate())
	
	// Get a section
	
	section, _ := doc.GetSection("owner")
	fmt.Println(section)
	
	// Get title
	
	value, _ = doc.GetValue("title")
	fmt.Println(value.AsString())
}
