TOML-GO
=======

An easy-to-use Go parser for the [Toml format](https://github.com/mojombo/toml).

This parser was last tested on TOML version [3f4224ecdc4a65fdd28b4fb70d46f4c0bd3700aa](https://github.com/mojombo/toml/tree/00f11b019406531c8c7989846b1c1a54e9b8d8bb)

It correctly parses both the [example.toml](https://github.com/mojombo/toml/blob/master/tests/example.toml) and [hard_example.toml](https://github.com/mojombo/toml/blob/master/tests/hard_example.toml) files of the official Toml repository.

Usage
-----

See [main.go](main.go) for some examples or the [API doc](http://godoc.org/github.com/laurent22/toml-go/toml) for the full public API. Basically, you create a parser and give it a file or string to parse. You can then access the values using the `Get` accessors (eg. `GetString()`, `GetInt()`, etc.). You can also provide an optional default value to any of these accessors.

```go
var parser toml.Parser
doc := parser.ParseFile("example.toml")

// Or parse a string directly:
// doc := parser.Parse(someTomlString)

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
```

Advanced usage
--------------

```go
// ==================================================
// The GetValue() / As<Type>() pattern provides a bit
// more flexibility but is more verbose
// ==================================================

var value toml.Value
var ok bool

// Get a string:

value, ok = doc.GetValue("servers.beta.ip")
if !ok { // Optionally, you can check if the value exists or not
  panic("value doesn't exist")
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

// Get the title

value, _ = doc.GetValue("title")
fmt.Println(value.AsString())
```
