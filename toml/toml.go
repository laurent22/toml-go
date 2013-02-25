package toml

import (
	"strings"
	"io/ioutil"
	"fmt"
	"strconv"
	"time"
)

type Kind int

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
	asBool bool
	asInt int64
	asFloat float64
	asString string
	asArray []Value
	asDate time.Time
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

func ParseValue(s string) (Value, int, bool) {
	var v Value
	if len(s) == 0 { return v, 0, false }
	
	if strings.Index(s, "true") == 0 || strings.Index(s, "false") == 0 {
		v.kind = kindBool
		v.asBool = strings.Index(s, "true") == 0
		index := 4
		if !v.asBool { index = 5 }
		v.raw = s[0:index]
		return v, index, true
	}
	
	if s[0] == '"' {
		parsed, index, ok := ParseString(s)
		if !ok { return v, 0, false }
		v.asString = parsed
		v.kind = kindString
		v.raw = s[0:index]
		return v, index, true
	}
	
	if s[0] == '[' {
		parsed, index, ok := ParseArray(s)
		if !ok { return v, 0, false }
		v.asArray = parsed
		v.kind = kindArray
		v.raw = s[0:index]
		return v, index, true
	}
	
	if len(s) >= 20 && s[19] == 'Z' {
		parsed, index, ok := ParseDate(s)
		if !ok { return v, 0, false }
		v.asDate = parsed
		v.kind = kindDate
		v.raw = s[0:index]
		return v, index, true
	}
	
	numString, index, ok := ParseNumber(s)
	if !ok { return v, 0, false }
	
	parsedInt, err := strconv.ParseInt(numString, 10, 64)
	if err == nil {
		v.asInt = parsedInt
		v.kind = kindInt
		v.raw = s[0:index]
		return v, index, true
	}
	
	parsedFloat, err := strconv.ParseFloat(numString, 64)
	if err == nil {
		v.asFloat = parsedFloat
		v.kind = kindFloat
		v.raw = s[0:index]
		return v, index, true
	}
	
	return v, 0, false
}

func ParseDate(s string) (time.Time, int, bool) {
	timeString := s[0:20]
	output, err := time.Parse(time.RFC3339, timeString)
	if err != nil { return output, 0, false }
	return output, 20, true
}

func ParseNumber(s string) (string, int, bool) {
	numberString := ""
	allowedChars := "0123456789."
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !strings.Contains(allowedChars, string(c)) {
			break
		}
		numberString += string(c)
	}
	if len(numberString) <= 0 { return "", 0, false }
	return numberString, len(numberString), true
}

func ParseArray(s string) ([]Value, int, bool) {
	var output []Value
	if len(s) <= 0 { return output, 0, false }

	endIndex := 0
	state := 0 // 0 = start, 1 = before value, 2 = end value
	for i := 0; i < len(s); i++ {
		c := s[i]
						
		if state == 0 {
			if c != '[' { return output, endIndex, false }
			state = 1
			continue
		}
		
		if state == 1  {
			v, index, ok := ParseValue(s[i:len(s)])
			if !ok {
				continue
			} else {
				output = append(output, v)
				i = i + index - 1 
				state = 2
			}
		}
		
		if state == 2 {
			if c == ',' {
				state = 1
				continue
			}
			if c == ']' {
				endIndex = i + 1
				break
			}
		}
	}
	
	return output, endIndex, true
}

func (this Value) AsArray() []Value {
	return this.asArray
}

func (this Value) RawNoComment() string {
	index := strings.Index(this.raw, "#")
	if index < 0 { return this.raw }
	return strings.Trim(this.raw[0:index], " \t\n\r")
}

func ParseString(s string) (string, int, bool) {
	if len(s) <= 0 { return "", 0, false }
	
	index := 0
	escape := false
	output := ""
	state := 0 // 0 = left, 1 = inside
	for i := 0; i < len(s); i++ {
		c := s[i]
		
		if state == 0 {
			if c != '"' { return "", 0, false }
			state = 1
			continue
		}
		
		if state == 1 {
			if c == '\\' {
				escape = true
				continue
			} 
			
			if c == '"' && !escape {
				index = i + 1
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
					return "", 0, false // Or panic?
				}
				escape = false
				continue
			}
			
			output += string(c)
		}
	}
		
	return output, index, true	
}

func (this Value) AsString() string {
	return this.asString
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

func (this Value) String() string {
	if this.kind == kindString { return "\"" + this.asString + "\"" }
	if this.kind == kindInt { return strconv.FormatInt(this.asInt, 10) }
	if this.kind == kindFloat { return strconv.FormatFloat(this.asFloat, 'f', -1, 64); }
	if this.kind == kindBool { if this.asBool { return "true" } else { return "false" } }
	if this.kind == kindDate { return this.asDate.Format(time.RFC3339) }
	if this.kind == kindArray {
		array := this.asArray
		output := ""
		for i := 0; i < len(array); i++ {
			if output != "" { output += ", " }
			output += array[i].String()
		}
		return "[" + output + "]"
	}
	return "undefined"
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
		output += this.name + " = " + this.value.String()
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

func (this Parser) ParseKey(line string) (string, int) {
	index := strings.Index(line, "=")
	if index < 0 { return "", -1 }
	return strings.Trim(line[0:index], " \t\n\r"), index
}

func (this *Node) LoadValues() {
	this.value, _, _ = ParseValue(this.value.raw)
	
	for _, node := range this.children {
		node.LoadValues()
	}
}

func (this *Node) GetSection(path string) (*Node, bool) {
	names := strings.Split(path, ".")
	current := this
	nameIndex := 0
	
	for {
		for _, node := range current.children {
			 if node.kind != kindSection { continue } 
			 if node.name == names[nameIndex] {
			 	current = node
			 	nameIndex++	
			 	if nameIndex >= len(names) {
			 		return current, true
			 	}
			 	break
			 }
		}
	}
	
	return current, false
}

func (this *Node) GetValue(path string) (Value, bool) {
	var output Value
	names := strings.Split(path, ".")
	if len(names) == 1 {
		node, ok := this.Child(path)
		if !ok { return output, false }
		return node.value, true
	}
	
	sectionPath := strings.Join(names[0:len(names) - 1], ".")
	section, ok := this.GetSection(sectionPath)
	if !ok { return output, false }
	return section.GetValue(names[len(names) - 1])
}

func (this Document) GetSection(path string) (*Node, bool) {
	return this.root.GetSection(path)	
}

func (this Document) GetValue(path string) (Value, bool) {
	return this.root.GetValue(path)
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
	
	output.root.LoadValues()
	
	//fmt.Println(output.root)
	
	return output
}

func (this Parser) ParseFile(tomlFilePath string) Document {
	content, err := ioutil.ReadFile(tomlFilePath)
	if err != nil {
		panic(err.Error())
	}
	return this.Parse(string(content))
}
