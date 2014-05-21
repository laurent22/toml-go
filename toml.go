// A Go parser library for the [Toml format](https://github.com/mojombo/toml).

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
	Children map[string]*Node
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

func cleanRawValue(s string) string {
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

func parseValue(s string) (Value, int, bool) {
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
		parsed, index, ok := parseString(s)
		if !ok { return v, 0, false }
		v.asString = parsed
		v.kind = kindString
		v.raw = s[0:index]
		return v, index, true
	}
	
	if s[0] == '[' {
		parsed, index, ok := parseArray(s)
		if !ok { return v, 0, false }
		v.asArray = parsed
		v.kind = kindArray
		v.raw = s[0:index]
		return v, index, true
	}
	
	if len(s) >= 20 && s[19] == 'Z' {
		parsed, index, ok := parseDate(s)
		if !ok { return v, 0, false }
		v.asDate = parsed
		v.kind = kindDate
		v.raw = s[0:index]
		return v, index, true
	}
	
	numString, index, ok := parseNumber(s)
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

func parseDate(s string) (time.Time, int, bool) {
	timeString := s[0:20]
	output, err := time.Parse(time.RFC3339, timeString)
	if err != nil { return output, 0, false }
	return output, 20, true
}

func parseNumber(s string) (string, int, bool) {
	numberString := ""
	allowedChars := "0123456789.-"
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

func parseArray(s string) ([]Value, int, bool) {
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
			v, index, ok := parseValue(s[i:len(s)])
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

func parseString(s string) (string, int, bool) {
	if len(s) <= 0 { return "", 0, false }
	
	index := 0
	escape := false
	output := ""
	state := 0 // 0 = left, 1 = inside
	i := -1
	for _, c := range s {
		i++
		
		if state == 0 {
			if c != '"' { return "", 0, false }
			state = 1
			continue
		}
		
		if state == 1 {
			if escape {
				if c == '0' {
					output += "\x00"
				} else if (c == 't') {
					output += "\t"
				} else if (c == 'n') {
					output += "\n"
				} else if (c == 'r') {
					output += "\r"
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
			
			if c == '\\' {
				escape = true
				continue
			} 
			
			if c == '"' && !escape {
				index = i + 1
				break
			}
			
			output += string(c)
		}
	}
		
	return output, index, true	
}

func (this Value) AsArray() []Value {
	return this.asArray
}

func (this Value) AsString() string {
	return this.asString
}

func (this Value) AsInt() int {
	return int(this.asInt)
}

func (this Value) AsInt8() int8 {
	return int8(this.asInt)
}

func (this Value) AsInt16() int16 {
	return int16(this.asInt)
}

func (this Value) AsInt32() int32 {
	return int32(this.asInt)
}

func (this Value) AsInt64() int64 {
	return this.asInt
}

func (this Value) AsFloat() float64 {
	return this.asFloat
}

func (this Value) AsFloat32() float32 {
	return float32(this.asFloat)
}

func (this Value) AsFloat64() float64 {
	return this.asFloat
}

func (this Value) AsBool() bool {
	return this.asBool	
}

func (this Value) AsDate() time.Time {
	return this.asDate
}

func (this *Node) createChildren() {
	if this.Children != nil { return }
	this.Children = make(map[string]*Node)
}

func (this *Node) child(name string) (*Node, bool) {
	if !this.hasChildren() { return nil, false }
	node, ok := this.Children[name]
	return node, ok
}

func (this *Node) setChild(name string, node *Node) {
	this.createChildren()
	this.Children[name] = node
	node.parent = this
}

func (this *Node) hasChildren() bool {
	return this.Children != nil
}

func (this Value) String() string {
	if this.kind == kindString {
		s := this.asString
		s = strings.Replace(s, "\n", "\\n", -1)
		s = strings.Replace(s, "\x00", "\\0", -1)
		s = strings.Replace(s, "\t", "\\t", -1)
		s = strings.Replace(s, "\r", "\\r", -1)
		s = strings.Replace(s, "\"", "\\\"", -1)
		s = strings.Replace(s, "\\", "\\\\", -1)
		return "\"" + s + "\""
	}
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
		
	if (this.kind == kindRoot && this.hasChildren()) {
		for _, node := range this.Children {
			output += node.String()
			output += "\n"
		}
	}
	
	if (this.kind == kindSection) {
		output += "[" + this.FullName() + "]"
		output += "\n"
		if (this.hasChildren()) {
			for _, node := range this.Children {
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

func newNodePointer() *Node {
	output := new(Node) 
	output.Children = nil
	output.parent = nil
	return output;
}

func newDocument() Document {
	var output Document
	output.root = newNodePointer()
	output.root.kind = kindRoot
	return output;
}

func (this Parser) parseKey(line string) (string, int) {
	index := strings.Index(line, "=")
	if index < 0 { return "", -1 }
	return strings.Trim(line[0:index], " \t\n\r"), index
}

func (this *Node) loadValues() {
	this.value, _, _ = parseValue(this.value.raw)
	
	for _, node := range this.Children {
		node.loadValues()
	}
}

func (this *Node) GetSection(path string) (*Node, bool) {
	names := strings.Split(path, ".")
	current := this
	nameIndex := 0
	
	for {
		for _, node := range current.Children {
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
		node, ok := this.child(path)
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
	
	output := newDocument()
	
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
				node, ok := current.child(name)
				if !ok {
					section := newNodePointer()
					section.name = name
					section.kind = kindSection
					current.setChild(name, section)
					current = section
				} else {
					current = node
				}
				currentSection = current
			}
			continue
		}
		
		// VALUE
		
		key, index := this.parseKey(line)
		if index < 0 {
			if currentValue != nil {
				currentValue.value.raw += line
			} else {
				panic("Invalid value: " + line)
			}
		} else {
			node := newNodePointer()
			node.name = key
			node.value.raw = cleanRawValue(line[index + 1:len(line)])
			node.kind = kindValue
			currentSection.setChild(key, node)
			currentValue = node
		}
	}
	
	output.root.loadValues()
		
	return output
}

func (this Parser) ParseFile(tomlFilePath string) Document {
	content, err := ioutil.ReadFile(tomlFilePath)
	if err != nil {
		panic(err.Error())
	}
	return this.Parse(string(content))
}
