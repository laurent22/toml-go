package main

import (
	"../toml"
	"fmt"
	"strconv"
	"time"
)

/**
 * Some very rudimentary assert functions to easily test
 */

func assertTrue(desc string, v bool) {
	if !v { panic("Failed: " + desc) }
	fmt.Print(".")
}

func assertFalse(desc string, v bool) {
	if v { panic("Failed: " + desc) }
	fmt.Print(".")
}

func assertIntEqual(desc string, a int, b int) {
	if a != b { panic("Failed: " + desc + " - " + strconv.FormatInt(int64(a), 10) + " != " + strconv.FormatInt(int64(b), 10)) }
	fmt.Print(".")
}

func assertFloatEqual(desc string, a float64, b float64) {
	if a != b { panic("Failed: " + desc + " - " + strconv.FormatFloat(a, 'f', -1, 64) + " != " + strconv.FormatFloat(b, 'f', -1, 64)) }
	fmt.Print(".")
}

func assertStringEqual(desc string, a string, b string) {
	if a != b { panic("Failed: " + desc + " - " + a + " != " + b) }
	fmt.Print(".")
}

func assertTimeEqual(desc string, a time.Time, b time.Time) {
	if a != b { panic("Failed: " + desc + " - " + a.Format(time.RFC3339) + " != " + b.Format(time.RFC3339)) }
	fmt.Print(".")
}

func main() {
	var parser toml.Parser
	doc := parser.ParseFile("test1.toml")
	
	var v toml.Value
	var ok bool
	
	v, ok = doc.GetValue("title")
	assertTrue("Value exists", ok)
	assertStringEqual("String is correct", v.AsString(), "TOML Example")
	
	v, ok = doc.GetValue("doesntexist")
	assertFalse("Value doesn't exist", ok)
	
	v, ok = doc.GetValue("database.enabled")
	assertTrue("Boolean is valid", v.AsBool())
	
	v, ok = doc.GetValue("database.debug")
	assertFalse("Boolean is valid", v.AsBool())
	
	v, ok = doc.GetValue("owner.name")
	assertStringEqual("Nested value", v.AsString(), "Tom Preston-Werner")
	
	v, ok = doc.GetValue("owner.dob")
	expectedTime, _ := time.Parse(time.RFC3339, "1979-05-27T07:32:00Z")
	assertTimeEqual("Date is valid", v.AsDate(), expectedTime)
	
	v, ok = doc.GetValue("database.connection_max")
	assertIntEqual("Int is valid", v.AsInt(), 5000)
	
	v, ok = doc.GetValue("floats.pi")
	assertFloatEqual("Float is valid", v.AsFloat(), 3.14)
	
	v, ok = doc.GetValue("floats.minus")
	assertFloatEqual("Negative float is valid", v.AsFloat(), -10.001)
	
	v, ok = doc.GetValue("database.ports")
	assertIntEqual("Array size if correct", len(v.AsArray()), 3)
	assertIntEqual("Array content is correct", v.AsArray()[0].AsInt(), 8001)
	assertIntEqual("Array content is correct", v.AsArray()[1].AsInt(), 8001)
	assertIntEqual("Array content is correct", v.AsArray()[2].AsInt(), 8002)
	
	_, ok = doc.GetSection("servers.alpha")
	assertTrue("Nested section exists", ok)
	
	v, ok = doc.GetValue("clients.data")
	assertIntEqual("Array size is correct", len(v.AsArray()), 2)
	subArray0 := v.AsArray()[0].AsArray()
	subArray1 := v.AsArray()[1].AsArray()
	assertIntEqual("Sub-array size is correct", len(subArray0), 2)
	assertIntEqual("Sub-array size is correct", len(subArray1), 3)
	assertStringEqual("Sub-array content is correct", subArray0[0].AsString(), "gamma")
	assertStringEqual("Sub-array content is correct", subArray0[1].AsString(), "delta")
	assertIntEqual("Sub-array content is correct", subArray1[0].AsInt(), 1)
	assertIntEqual("Sub-array content is correct", subArray1[1].AsInt(), 2)
	assertIntEqual("Sub-array content is correct", subArray1[2].AsInt(), 123)
	
	v, ok = doc.GetValue("multilinearray.test")
	assertIntEqual("Array size is correct", len(v.AsArray()), 3)
	assertStringEqual("Array content is correct", v.AsArray()[0].AsString(), "one")
	assertStringEqual("Array content is correct", v.AsArray()[1].AsString(), "two")
	assertStringEqual("Array content is correct", v.AsArray()[2].AsString(), "three")
}