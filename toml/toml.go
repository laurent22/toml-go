package toml

import (
	"strings"
	"io/ioutil"
	"fmt"
	"strconv"
)

type Kind int32

const (
	kindRoot = 1
	kindSection = 2
	kindValue = 3
	kindBool = 4
	kindString = 5
	kindInt = 6
	kindFloat = 7
	kindArray = 8
	kindDate = 9
)

type Parser struct {
	
}

type Node struct {
	name string
	value Value
	kind Kind
	children map[string]*Node
	parent *Node
}

type Value struct {
	raw string
	kind Kind
}

type Document struct {
	root *Node
}

func CleanRawValue(s string) string {
	s = strings.Trim(s, " \t\n\r")
	if len(s) == 0 { return "" }
	
	inString := false
	escape := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if escape {
			escape = false
			continue
		}
		if c == '"' {
			inString = !inString
			continue
		}
		if inString && c == '\\' {
			escape = true
			continue
		}
		if c == '#' && !inString {
			return strings.Trim(s[0:i], " \t\n\r")
		}
	}
	
	return s
	
}

// func ParseValue(s string) Value, bool {
// 	var v Value
// 	v.raw = s
// 	if s == "true" || s == "false" {
// 		v.kind = kindBool
// 	}
// 	if s[0] == '"' && s[len(s) - 1] == '"' {
// 		v.kind = kindString
// 	}
// 	if 
// 	return v
// }

func ParseArray(s string) ([]Value, bool) {
	var output []Value
	if len(s) <= 0 { return output, false }

	//current := output
	//openingCount := 0
	state := 0 // 1 = start, 2 = inside value, 3 = end value
	for i := 1; i < len(s); i++ {
		c := s[i]

		if state == 0 {
			if c != '[' { return output, false }
			state = 1
			continue
		}
		
		if state == 1  {
			
			
			state = 2
			
		}
		// if c == '[' {
		// 	openingCount++
		// }
		// if c == ']' {
		// 	if openingCount == 0 { // Found the closing bracket
		// 		var value Value
		// 		value.raw = s[1:i]
		// 		append(output, value)
		// 	} else {
		// 		openingCount--
		// 	}
		// }
	}	
	
	return output, true
}

func (this Value) AsArray() ([]Value, bool) {
	return ParseArray(this.raw)
}

func (this Value) RawNoComment() string {
	index := strings.Index(this.raw, "#")
	if index < 0 { return this.raw }
	return strings.Trim(this.raw[0:index], " \t\n\r")
}

func ParseString(s string) (string, bool) {
	if len(s) <= 0 { return "", false }
	
	escape := false
	output := ""
	state := 0 // 0 = left, 1 = inside
	for i := 0; i < len(s); i++ {
		c := s[i]
		
		if state == 0 {
			if c != '"' { return "", false }
			state = 1
			continue
		}
		
		if state == 1 {
			if c == '\\' {
				escape = true
				continue
			} 
			
			if c == '"' && !escape {
				break
			}
			
			if escape {
				if c == '0' {
					
				} else if (c == 't') {
					output += "\t"
				} else if (c == 'n') {
					output += "\n"
				} else if (c == '"') {
					output += "\""
				} else if (c == '\\') {
					output += "\\"
				} else {
					return "", false // Or panic?
				}
				escape = false
				continue
			}
			
			output += string(c)
		}
	}
	
	return output, true	
}

func (this Value) AsString() (string, bool) {
	return ParseString(this.raw)
}

func (this Value) AsInt() (int64, bool) {
	output, err := strconv.ParseInt(this.RawNoComment(), 10, 64)
	if err != nil { return 0, false }
	return output, true
}

func (this *Node) CreateChildren_() {
	if this.children != nil { return }
	this.children = make(map[string]*Node)
}

func (this *Node) Child(name string) (*Node, bool) {
	if !this.HasChildren() { return nil, false }
	node, ok := this.children[name]
	return node, ok
}

func (this *Node) SetChild(name string, node *Node) {
	this.CreateChildren_()
	this.children[name] = node
	node.parent = this
}

func (this *Node) HasChildren() bool {
	return this.children != nil
}

func (this *Node) String() string {
	output := ""
		
	if (this.kind == kindRoot && this.HasChildren()) {
		for _, node := range this.children {
			output += node.String()
			output += "\n"
		}
	}
	
	if (this.kind == kindSection) {
		output += "[" + this.FullName() + "]"
		output += "\n"
		if (this.HasChildren()) {
			for _, node := range this.children {
				output += node.String()
			}
		}
	}
	
	if (this.kind == kindValue) {
		asInt, _ := this.value.AsInt()
		asString, _ := this.value.AsString()
		fmt.Println(this.name, "as int", asInt)
		fmt.Println(this.name, "as string", asString)
		output += this.name + " = " + this.value.raw
		output += "\n"
	}
	
	return output
}

func (this *Node) FullName() string {
	output := this.name
	current := this
	for {
		current = current.parent
		if current == nil || current.kind == kindRoot { break }
		output = current.name + "." + output
	}
	return output
}

func (this Document) String() string {
	return this.root.String()
}

func NewNodePointer() *Node {
	output := new(Node) 
	output.children = nil
	output.parent = nil
	return output;
}

func NewDocument() Document {
	var output Document
	output.root = NewNodePointer()
	output.root.kind = kindRoot
	return output;
}

// func (this Document) Get(path string) (*Node, bool) {
// 	names := strings.Split(path, ".")
// 	current := this.root
// 	for i := 0; i < len(names); i++ {
// 		node, ok := current.children[names[i]]
// 		if !ok { return current, false }
// 		current = node
// 	}
// 	return current, true
// }

func (this Parser) ParseKey(line string) (string, int) {
	index := strings.Index(line, "=")
	if index < 0 { return "", -1 }
	return strings.Trim(line[0:index], " \t\n\r"), index
}

func (this Parser) Parse(tomlString string) Document {
	_ = fmt.Println
	
	output := NewDocument()
	
	var currentValue *Node
	currentValue = nil 
	currentSection := output.root
	lines := strings.Split(tomlString, "\n")
	for i := 0; i < len(lines); i++ {
		var line = strings.Trim(lines[i], " \t\n\r")
		if len(line) == 0 { continue }
		
		// COMMENT
		
		if line[0] == '#' { continue }
		
		// SECTION
		
		if line[0] == '[' && line[len(line) - 1] == ']' {
			currentValue = nil
			path := line[1:len(line) - 1]
			names := strings.Split(path, ".")
			current := output.root
			for j := 0; j < len(names); j++ {
				name := names[j]
				node, ok := current.Child(name)
				if !ok {
					section := NewNodePointer()
					section.name = name
					section.kind = kindSection
					current.SetChild(name, section)
					current = section
				} else {
					current = node
				}
				currentSection = current
			}
			continue
		}
		
		// VALUE
		
		key, index := this.ParseKey(line)
		if index < 0 {
			if currentValue != nil {
				currentValue.value.raw += line
			} else {
				panic("Invalid value: " + line)
			}
		} else {
			node := NewNodePointer()
			node.name = key
			node.value.raw = CleanRawValue(line[index + 1:len(line)])
			node.kind = kindValue
			currentSection.SetChild(key, node)
			currentValue = node
		}
	}
	
	fmt.Println(output.String())
	
	return output
}

func (this Parser) ParseFile(tomlFilePath string) Document {
	content, err := ioutil.ReadFile(tomlFilePath)
	if err != nil {
		panic(err.Error())
	}
	return this.Parse(string(content))
}
