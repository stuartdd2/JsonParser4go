package test

import (
	"fmt"
	"github.com/stuartdd/jsonParserGo/parser"
	"strings"
	"testing"
)

var (
	sample4 = []byte(`{"age": 28,"number": "7349282382","firstName": "Joe","lastName": "Jackson","gender": "male"}`)
	sample5 = []byte(`[{"age": 28},{"number": "7349282382"},"firstName"]`)
)

func TestParserSample5(t *testing.T) {
	n, err := parser.Parse(sample5)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	l := (n).(*parser.JsonList)
	l.Add(parser.NewJsonString("A", "B"))
	l.Add(parser.NewJsonBool("T", true))
	l.Add(parser.NewJsonNumber("", 10.5))
	l.Add(parser.NewJsonObject(""))
	s := n.JsonValue()
	if s != "[{\"age\": 28},{\"number\": \"7349282382\"},\"firstName\",{\"A\": \"B\"},{\"T\": true},10.5,{}]" {
		t.Errorf("Result does not match expected")
	}
}
func TestParserNotJson(t *testing.T) {
	_, err := parser.Parse([]byte(`Hello world`))
	if err == nil {
		t.Errorf("Error: %s", err.Error())
	}
}

func TestParserSample4(t *testing.T) {
	_, err := parser.Parse(sample4)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

func TestParserWithObjects(t *testing.T) {
	_, err := parser.Parse([]byte(`{"lit":"literal", "num":-099.9, "t":true, "list": [true, false], "obj":{"a":1, "b":true, "c":{"t":true, "n":1.000} }, "f":false}`))
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

func TestParserWithListWithObjects(t *testing.T) {
	_, err := parser.Parse([]byte(`["literal", {"obj":"literal"}, {"num":99.9}, {"t":true}, {"f":false}]`))
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

func TestParserWithEmptyList(t *testing.T) {
	_, err := parser.Parse([]byte(`[{"eList":[]}, {"num":99.9}, {"t":true}, {"f":false}]`))
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}
func TestParserWithEmptyObject(t *testing.T) {
	_, err := parser.Parse([]byte(`[{"eList":{}}, {"num":99.9}, {"t":true}, {"f":false}]`))
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

func TestParserWithListNoError(t *testing.T) {
	node, err := parser.Parse([]byte(`["literal", 1234.5, true, false]`))
	if err != nil {
		t.Errorf("parser.Parse: Should not have returned as error: %s", err.Error())
	}
	if node.GetNodeType() != parser.NT_LIST {
		t.Errorf("parser.AsNodeType: Expected parser.NT_LIST")
	}
	if fmt.Sprintf("%T", node) != "*parser.JsonList" {
		t.Errorf("node Type: Expected *parser.JsonList actual %T", node)
	}
	testListNode(t, "PARSED_LIST", node.(*parser.JsonList))
	node = makeList()
	testListNode(t, "MAKE_LIST", node.(*parser.JsonList))
}

func testListNode(t *testing.T, id string, listNode *parser.JsonList) {
	if fmt.Sprintf("%T", listNode) != "*parser.JsonList" {
		t.Errorf("%s: listNode: Expected *parser.JsonList", id)
	}
	if listNode.GetNodeType() != parser.NT_LIST {
		t.Errorf("%s: listNode: Expected parser.NT_LIST", id)
	}

	n0 := listNode.GetNodeAt(0)
	if fmt.Sprintf("%T", n0) != "*parser.JsonString" {
		t.Errorf("%s: 0 Type: Expected *parser.JsonString, actual %T", id, n0)
	}
	if n0.GetNodeType() != parser.NT_STRING {
		t.Errorf("%s: 0 GetNodeType: Expected parser.NT_STRING", id)
	}
	if n0.GetName() != "" {
		t.Errorf("%s: 0 GetName(): Expected '' ", id)
	}
	if n0.JsonValue() != "\"literal\"" {
		t.Errorf("%s: 0 JsonValue(): Expected \"literal\" actual %s", id, n0.JsonValue())
	}
	l0 := n0.(*parser.JsonString)
	if l0.GetValue() != "literal" {
		t.Errorf("%s: 0 GetValue: Expected \"literal\"", id)
	}

	n1 := listNode.GetNodeAt(1)
	if fmt.Sprintf("%T", n1) != "*parser.JsonNumber" {
		t.Errorf("%s: 1 Type: Expected *parser.JsonNumber, actual %T", id, n0)
	}
	if n1.GetNodeType() != parser.NT_NUMBER {
		t.Errorf("%s: 1 GetNodeType: Expected parser.NT_NUMBER", id)
	}
	if n1.GetName() != "" {
		t.Errorf("%s: 1 GetName(): Expected '' ", id)
	}
	if n1.String() != "1234.5" {
		t.Errorf("%s: 1 String(): Expected \"1234.5\" actual %s", id, n1.String())
	}
	l1 := n1.(*parser.JsonNumber)
	if l1.GetValue() != 1234.5 {
		t.Errorf("%s: 1 GetValue: Expected \"1234.5\" actual %s", id, l1.String())
	}
	if l1.GetIntValue() != 1234 {
		t.Errorf("%s: 1 GetValue: Expected \"1234\" actual %d", id, l1.GetIntValue())
	}

	n2 := listNode.GetNodeAt(2)
	if fmt.Sprintf("%T", n2) != "*parser.JsonBool" {
		t.Errorf("%s: 2 Type: Expected *parser.JsonBool, actual %T", id, n0)
	}
	if n2.GetNodeType() != parser.NT_BOOL {
		t.Errorf("%s: 2 GetNodeType: Expected parser.NT_BOOL", id)
	}
	if n2.GetName() != "" {
		t.Errorf("%s: 2 GetName(): Expected '' ", id)
	}
	if n2.String() != "true" {
		t.Errorf("%s: 2 String(): Expected \"true\"", id)
	}
	l2 := n2.(*parser.JsonBool)
	if l2.GetValue() != true {
		t.Errorf("%s: 2 GetValue: Expected \"true\" actual %t", id, l2.GetValue())
	}

	n3 := listNode.GetNodeAt(3)
	if fmt.Sprintf("%T", n3) != "*parser.JsonBool" {
		t.Errorf("%s: 3 Type: Expected *parser.JsonBool, actual %T", id, n0)
	}
	if n3.GetNodeType() != parser.NT_BOOL {
		t.Errorf("%s: 3 GetNodeType: Expected parser.NT_BOOL", id)
	}
	if n3.GetName() != "" {
		t.Errorf("%s: 3 GetName(): Expected '' ", id)
	}
	if n3.String() != "false" {
		t.Errorf("%s: 3 String(): Expected \"false\"", id)
	}
	l3 := n3.(*parser.JsonBool)
	if l3.GetValue() != false {
		t.Errorf("%s: 3 GetValue: Expected \"false\" actual %t", id, l2.GetValue())
	}
}

func TestErrorNoComma(t *testing.T) {
	_, err := parser.Parse([]byte(`["literal" 1234, true]`))
	if err == nil {
		t.Errorf("Should have returned a 'Expected a ',' seperator'")
	}
}

func TestErrorUnterminatedLit(t *testing.T) {
	_, err := parser.Parse([]byte(`["literal, 1234, true]`))
	if err == nil {
		t.Errorf("Should have returned a 'unterminated quote'")
	}
}

func TestErrorBadNumber(t *testing.T) {
	_, err := parser.Parse([]byte(`["literal", 12x34, true]`))
	if err == nil {
		t.Errorf("Should have returned a 'unrecognised token. x'")
	}
	// t.Error(err.Error())
}

func TestErrorDuplicateNames(t *testing.T) {
	_, err := parser.Parse([]byte(`{"A":"a","A":true}`))
	if err == nil {
		t.Errorf("Should have returned an error")
	}
	if !strings.Contains(err.Error(), "duplicate name") {
		t.Errorf("Should have returned a 'duplicate name' error")
	}
}

func TestErrorEmptyNames(t *testing.T) {
	_, err := parser.Parse([]byte(`{"A":"a","":true}`))
	if err == nil {
		t.Errorf("Should have returned an error")
	}
	if !strings.Contains(err.Error(), "must have a name") {
		t.Errorf("Should have returned a 'must have a name' error")
	}
}

func TestErrorBadBool(t *testing.T) {
	_, err := parser.Parse([]byte(`["literal", 1234, falxe]`))
	if err == nil {
		t.Errorf("Should have returned a 'Expected a ',' seperator. Found 7")
	}
}

func TestErrorBadObjectName(t *testing.T) {
	_, err := parser.Parse([]byte(`["literal", 1234, {true}]`))
	if err == nil {
		t.Errorf("Should have returned a 'Expected a ',' seperator. Found 7")
	}
}

func TestErrorBadObjectMissingColon(t *testing.T) {
	_, err := parser.Parse([]byte(`["literal", 1234, {"true" true}]`))
	if err == nil {
		t.Errorf("Should have returned a 'Expected a ',' seperator. Found 7")
	}
	// t.Errorf("Error: %s", err.Error())
}
func makeList() parser.NodeI {
	p := parser.NewJsonList("")
	s := parser.NewJsonString("", "literal")
	n := parser.NewJsonNumber("", 1234.5)
	b1 := parser.NewJsonBool("", true)
	b2 := parser.NewJsonBool("", false)
	p.Add(s)
	p.Add(n)
	p.Add(b1)
	p.Add(b2)
	return p
}
