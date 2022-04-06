package test

import (
	"fmt"
	"testing"

	"github.com/stuartdd2/JsonParser4go/parser"
)

func TestEqualWithParse(t *testing.T) {
	na, err := parser.Parse(obj3)
	if err != nil {
		t.Errorf("Failed to parse obj3 to na")
	}
	nb, err := parser.Parse(obj3)
	if err != nil {
		t.Errorf("Failed to parse obj3 to nb")
	}
	if !na.Equal(nb) {
		t.Errorf("failed: Object '%s' should equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	naState, err := parser.Find(na, parser.NewDotPath("address.state"))
	if err != nil {
		t.Errorf("Failed to find 'state' in na")
	}
	nbState, err := parser.Find(nb, parser.NewDotPath("address.state"))
	if err != nil {
		t.Errorf("Failed to find 'state' in nb")
	}
	naState.(*parser.JsonString).SetValue("DA")
	if na.Equal(nb) {
		t.Errorf("failed: Object '%s' should NOT equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	nbState.(*parser.JsonString).SetValue("DA")
	if !na.Equal(nb) {
		t.Errorf("failed: Object '%s' should equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	naPn, err := parser.Find(na, parser.NewDotPath("address.phoneNumbers.1"))
	if err != nil {
		t.Errorf("Failed to find 'address.phoneNumbers[1]' in na. err:%s", err)
	}
	nbPn, err := parser.Find(nb, parser.NewDotPath("address.phoneNumbers.1"))
	if err != nil {
		t.Errorf("Failed to find 'address.phoneNumbers[1]' in nb. err:%s", err)
	}
	naPn.(*parser.JsonNumber).SetValue(100)
	if na.Equal(nb) {
		t.Errorf("failed: Object '%s' should NOT equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	nbPn.(*parser.JsonNumber).SetValue(100)
	if !na.Equal(nb) {
		t.Errorf("failed: Object '%s' should equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	parser.Remove(na, naPn)
	if na.Equal(nb) {
		t.Errorf("failed: Object '%s' should NOT equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	parser.Remove(nb, nbPn)
	if !na.Equal(nb) {
		t.Errorf("failed: Object '%s' should equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}

}
func TestEqualObject(t *testing.T) {
	na := parser.NewJsonObject("a")
	nb := parser.NewJsonObject("a")
	if na.Equal(parser.NewJsonObject("b")) {
		t.Error("failed: Object 'a' should NOT equal Object 'b'")
	}
	if na.Equal(parser.NewJsonObject("")) {
		t.Error("failed: Object 'a' should NOT equal Object ''")
	}
	if !na.Equal(parser.NewJsonObject("a")) {
		t.Error("failed: Object 'a' should equal Object 'a'")
	}
	l1a, err := na.Add(parser.NewJsonString("l1", "list 1"))
	if err != nil {
		t.Errorf("failed: Object add na 'l1:List 1' -> %s", err)
	}
	l1aStr := l1a.(*parser.JsonString)
	l1b, err := nb.Add(parser.NewJsonString("l1", "list 1"))
	if err != nil {
		t.Errorf("failed: Object add nb 'l1:List 1' -> %s", err)
	}
	l1bStr := l1b.(*parser.JsonString)
	if !na.Equal(nb) {
		t.Errorf("failed: Object '%s' should equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	parser.Rename(na, l1a, "l2")
	if na.Equal(nb) {
		t.Errorf("failed: Object '%s' should NOT equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	parser.Rename(na, l1a, "l1")
	if !na.Equal(nb) {
		t.Errorf("failed: Object '%s' should equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	l1aStr.SetValue("List 2")
	if na.Equal(nb) {
		t.Errorf("failed: Object '%s' should NOT equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	l1bStr.SetValue("List 2")
	if !na.Equal(nb) {
		t.Errorf("failed: Object '%s' should equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	// Test order of add does impact Equal
	for i := 0; i < 10; i++ { // add 10 objects
		name := fmt.Sprintf("L%d", i)
		valu := fmt.Sprintf("List %d", i)
		nb.Add(parser.NewJsonString(name, valu))
	}
	for i := 9; i >= 0; i-- { // Add 10 objects in reverse order
		name := fmt.Sprintf("L%d", i)
		valu := fmt.Sprintf("List %d", i)
		na.Add(parser.NewJsonString(name, valu))
	}
	if !na.Equal(nb) {
		t.Errorf("failed: Object '%s' should equal Object '%s'", na.JsonValue(), nb.JsonValue())
	}
	t.Errorf("%s\n", na.JsonValue())
	t.Errorf("%s\n", nb.JsonValue())

}
func TestEqualList(t *testing.T) {
	na := parser.NewJsonList("a")
	nb := parser.NewJsonList("a")
	if na.Equal(parser.NewJsonList("b")) {
		t.Error("failed: List 'a' should NOT equal List 'b'")
	}
	if na.Equal(parser.NewJsonList("")) {
		t.Error("failed: List 'a' should NOT equal List ''")
	}
	if !na.Equal(parser.NewJsonList("a")) {
		t.Error("failed: List 'a' should equal List 'a'")
	}

	na.Add(parser.NewJsonString("", "list 1"))
	nb.Add(parser.NewJsonString("", "list 1"))
	if na.Equal(parser.NewJsonList("b")) {
		t.Errorf("failed: List '%s' should NOT equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	if !na.Equal(nb) {
		t.Errorf("failed: List '%s' should equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	nb.Add(parser.NewJsonBool("", true))
	if na.Equal(nb) {
		t.Errorf("failed: List '%s' should NOT equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	na.Add(parser.NewJsonBool("", true))
	if !na.Equal(nb) {
		t.Errorf("failed: List '%s' should equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	na.Add(parser.NewJsonBool("ob 1", true))
	if na.Equal(nb) {
		t.Errorf("failed: List '%s' should NOT equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	nb.Add(parser.NewJsonBool("ob 1", true))
	if !na.Equal(nb) {
		t.Errorf("failed: List '%s' should equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	na.Add(parser.NewJsonBool("ob 2", true))
	if na.Equal(nb) {
		t.Errorf("failed: List '%s' should NOT equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	ob3, _ := nb.Add(parser.NewJsonBool("ob 3", true))
	if na.Equal(nb) {
		t.Errorf("failed: List '%s' should NOT equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	na.Add(parser.NewJsonBool("ob 3", true))
	ob2, _ := nb.Add(parser.NewJsonBool("ob 2", true))
	if na.Equal(nb) {
		t.Errorf("failed: List '%s' should NOT equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	parser.Rename(nb, ob2, "ob 3")
	parser.Rename(nb, ob3, "ob 2")
	if !na.Equal(nb) {
		t.Errorf("failed: List '%s' should equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
}

func TestEqualString(t *testing.T) {
	n1 := parser.NewJsonString("a", "abc")
	if n1.Equal(parser.NewJsonString("a", "def")) {
		t.Error("failed: String 'a:abc' should NOT equal String 'a:def'")
	}
	if n1.Equal(parser.NewJsonString("b", "def")) {
		t.Error("failed: String 'a:abc' should NOT equal String 'b:def'")
	}
	if n1.Equal(parser.NewJsonString("", "abc")) {
		t.Error("failed: String 'a:abc' should NOT equal String ':abc'")
	}
	if !n1.Equal(parser.NewJsonString("a", "abc")) {
		t.Error("failed: String 'a:abc' should equal String 'a:abc'")
	}
	n := parser.NewJsonString("", "def")
	if n.Equal(parser.NewJsonString("a", "def")) {
		t.Error("failed: String ':def' should NOT equal String 'a:def'")
	}
	if n.Equal(parser.NewJsonString("", "abc")) {
		t.Error("failed: String ':def' should NOT equal String ':abc'")
	}
	if !n.Equal(parser.NewJsonString("", "def")) {
		t.Error("failed: String ':def' should equal String ':def'")
	}
	n3 := parser.NewJsonString("a", "true")
	if n3.Equal(parser.NewJsonBool("a", true)) {
		t.Error("failed: String 'a:true' should NOT equal Bool 'a:true'")
	}
	n4 := parser.NewJsonString("a", "123")
	if n4.Equal(parser.NewJsonNumber("a", 123)) {
		t.Error("failed: String 'a:123' should NOT equal Number 'a:123'")
	}
}

func TestEqualBool(t *testing.T) {
	n1 := parser.NewJsonBool("a", true)
	if n1.Equal(parser.NewJsonBool("a", false)) {
		t.Error("failed: Bool 'a:true' should NOT equal Bool 'a:false'")
	}
	if n1.Equal(parser.NewJsonBool("b", false)) {
		t.Error("failed: Bool 'a:true' should NOT equal Bool 'b:false'")
	}
	if n1.Equal(parser.NewJsonBool("", true)) {
		t.Error("failed: Bool 'a:true' should NOT equal Bool ':true'")
	}
	if !n1.Equal(parser.NewJsonBool("a", true)) {
		t.Error("failed: Bool 'a:true' should equal Bool 'a:true'")
	}
	n := parser.NewJsonBool("", false)
	if n.Equal(parser.NewJsonBool("a", false)) {
		t.Error("failed: Bool ':false' should NOT equal Bool 'a:false'")
	}
	if n.Equal(parser.NewJsonBool("", true)) {
		t.Error("failed: Bool ':false' should NOT equal Bool ':true'")
	}
	if !n.Equal(parser.NewJsonBool("", false)) {
		t.Error("failed: Bool ':false' should equal Bool ':false'")
	}

}

func TestEqualNull(t *testing.T) {
	na := parser.NewJsonNull("a")
	if na.Equal(parser.NewJsonNull("b")) {
		t.Error("failed: Null 'a' should NOT equal Null 'b'")
	}
	if na.Equal(parser.NewJsonNumber("b", 10)) {
		t.Error("failed: Null 'a' should NOT equal Number 'b'")
	}
	if !na.Equal(parser.NewJsonNull("a")) {
		t.Error("failed: Null 'a' should equal 'a'")
	}
	n := parser.NewJsonNull("")
	if n.Equal(parser.NewJsonNull("b")) {
		t.Error("failed: Null '' should NOT equal Null 'b'")
	}
	if n.Equal(parser.NewJsonNumber("b", 10)) {
		t.Error("failed: Null '' should NOT equal Number 'b'")
	}
	if !n.Equal(parser.NewJsonNull("")) {
		t.Error("failed: Null '' should equal ''")
	}
}

func TestEqualNumber(t *testing.T) {
	n1 := parser.NewJsonNumber("a", 123)
	if n1.Equal(parser.NewJsonNumber("b", 124)) {
		t.Error("failed: Number 'a:123' should NOT equal 'b:123'")
	}
	if n1.Equal(parser.NewJsonNumber("a", 124)) {
		t.Error("failed: Number 'a:123' should NOT equal 'a:124'")
	}
	if n1.Equal(parser.NewJsonNull("a")) {
		t.Error("failed: Number 'a:123' should NOT equal Null 'a'")
	}
	if !n1.Equal(parser.NewJsonNumber("a", 123)) {
		t.Error("failed: Number 'a:123' should equal 'a:123'")
	}
	n := parser.NewJsonNumber("", 123)
	if n.Equal(parser.NewJsonNumber("", 124)) {
		t.Error("failed: Number ':123' should NOT equal ':124'")
	}
	if !n.Equal(parser.NewJsonNumber("", 123)) {
		t.Error("failed: Number ':123' should equal ':123'")
	}
}
