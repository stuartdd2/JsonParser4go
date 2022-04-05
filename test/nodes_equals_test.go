package test

import (
	"testing"

	"github.com/stuartdd2/JsonParser4go/parser"
)

func TestEqualList(t *testing.T) {
	na := parser.NewJsonList("a")
	nb := parser.NewJsonList("a")
	if na.Equal(parser.NewJsonList("b")) {
		t.Error("failed: List 'a' should not equal List 'b'")
	}
	if na.Equal(parser.NewJsonList("")) {
		t.Error("failed: List 'a' should not equal List ''")
	}
	if !na.Equal(parser.NewJsonList("a")) {
		t.Error("failed: List 'a' should equal List 'a'")
	}

	na.Add(parser.NewJsonString("", "list 1"))
	nb.Add(parser.NewJsonString("", "list 1"))
	if na.Equal(parser.NewJsonList("b")) {
		t.Errorf("failed: List '%s' should not equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	if !na.Equal(nb) {
		t.Errorf("failed: List '%s' should equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	nb.Add(parser.NewJsonBool("", true))
	if na.Equal(nb) {
		t.Errorf("failed: List '%s' should not equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	na.Add(parser.NewJsonBool("", true))
	if !na.Equal(nb) {
		t.Errorf("failed: List '%s' should equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	na.Add(parser.NewJsonBool("ob 1", true))
	if na.Equal(nb) {
		t.Errorf("failed: List '%s' should not equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	nb.Add(parser.NewJsonBool("ob 1", true))
	if !na.Equal(nb) {
		t.Errorf("failed: List '%s' should equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	na.Add(parser.NewJsonBool("ob 2", true))
	if na.Equal(nb) {
		t.Errorf("failed: List '%s' should not equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	ob3, _ := nb.Add(parser.NewJsonBool("ob 3", true))
	if na.Equal(nb) {
		t.Errorf("failed: List '%s' should not equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	na.Add(parser.NewJsonBool("ob 3", true))
	ob2, _ := nb.Add(parser.NewJsonBool("ob 2", true))
	if na.Equal(nb) {
		t.Errorf("failed: List '%s' should not equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	parser.Rename(nb, ob2, "ob 3")
	parser.Rename(nb, ob3, "ob 2")
	if !na.Equal(nb) {
		t.Errorf("failed: List '%s' should equal List '%s'", na.JsonValue(), nb.JsonValue())
	}
	// t.Errorf("%s\n", na.JsonValue())
	// t.Errorf("%s\n", nb.JsonValue())
}

func TestEqualString(t *testing.T) {
	n1 := parser.NewJsonString("a", "abc")
	if n1.Equal(parser.NewJsonString("a", "def")) {
		t.Error("failed: String 'a:abc' should not equal String 'a:def'")
	}
	if n1.Equal(parser.NewJsonString("b", "def")) {
		t.Error("failed: String 'a:abc' should not equal String 'b:def'")
	}
	if n1.Equal(parser.NewJsonString("", "abc")) {
		t.Error("failed: String 'a:abc' should not equal String ':abc'")
	}
	if !n1.Equal(parser.NewJsonString("a", "abc")) {
		t.Error("failed: String 'a:abc' should equal String 'a:abc'")
	}
	n := parser.NewJsonString("", "def")
	if n.Equal(parser.NewJsonString("a", "def")) {
		t.Error("failed: String ':def' should not equal String 'a:def'")
	}
	if n.Equal(parser.NewJsonString("", "abc")) {
		t.Error("failed: String ':def' should not equal String ':abc'")
	}
	if !n.Equal(parser.NewJsonString("", "def")) {
		t.Error("failed: String ':def' should equal String ':def'")
	}
	n3 := parser.NewJsonString("a", "true")
	if n3.Equal(parser.NewJsonBool("a", true)) {
		t.Error("failed: String 'a:true' should not equal Bool 'a:true'")
	}
	n4 := parser.NewJsonString("a", "123")
	if n4.Equal(parser.NewJsonNumber("a", 123)) {
		t.Error("failed: String 'a:123' should not equal Number 'a:123'")
	}
}

func TestEqualBool(t *testing.T) {
	n1 := parser.NewJsonBool("a", true)
	if n1.Equal(parser.NewJsonBool("a", false)) {
		t.Error("failed: Bool 'a:true' should not equal Bool 'a:false'")
	}
	if n1.Equal(parser.NewJsonBool("b", false)) {
		t.Error("failed: Bool 'a:true' should not equal Bool 'b:false'")
	}
	if n1.Equal(parser.NewJsonBool("", true)) {
		t.Error("failed: Bool 'a:true' should not equal Bool ':true'")
	}
	if !n1.Equal(parser.NewJsonBool("a", true)) {
		t.Error("failed: Bool 'a:true' should equal Bool 'a:true'")
	}
	n := parser.NewJsonBool("", false)
	if n.Equal(parser.NewJsonBool("a", false)) {
		t.Error("failed: Bool ':false' should not equal Bool 'a:false'")
	}
	if n.Equal(parser.NewJsonBool("", true)) {
		t.Error("failed: Bool ':false' should not equal Bool ':true'")
	}
	if !n.Equal(parser.NewJsonBool("", false)) {
		t.Error("failed: Bool ':false' should equal Bool ':false'")
	}

}

func TestEqualNull(t *testing.T) {
	na := parser.NewJsonNull("a")
	if na.Equal(parser.NewJsonNull("b")) {
		t.Error("failed: Null 'a' should not equal Null 'b'")
	}
	if na.Equal(parser.NewJsonNumber("b", 10)) {
		t.Error("failed: Null 'a' should not equal Number 'b'")
	}
	if !na.Equal(parser.NewJsonNull("a")) {
		t.Error("failed: Null 'a' should equal 'a'")
	}
	n := parser.NewJsonNull("")
	if n.Equal(parser.NewJsonNull("b")) {
		t.Error("failed: Null '' should not equal Null 'b'")
	}
	if n.Equal(parser.NewJsonNumber("b", 10)) {
		t.Error("failed: Null '' should not equal Number 'b'")
	}
	if !n.Equal(parser.NewJsonNull("")) {
		t.Error("failed: Null '' should equal ''")
	}
}

func TestEqualNumber(t *testing.T) {
	n1 := parser.NewJsonNumber("a", 123)
	if n1.Equal(parser.NewJsonNumber("b", 124)) {
		t.Error("failed: Number 'a:123' should not equal 'b:123'")
	}
	if n1.Equal(parser.NewJsonNumber("a", 124)) {
		t.Error("failed: Number 'a:123' should not equal 'a:124'")
	}
	if n1.Equal(parser.NewJsonNull("a")) {
		t.Error("failed: Number 'a:123' should not equal Null 'a'")
	}
	if !n1.Equal(parser.NewJsonNumber("a", 123)) {
		t.Error("failed: Number 'a:123' should equal 'a:123'")
	}
	n := parser.NewJsonNumber("", 123)
	if n.Equal(parser.NewJsonNumber("", 124)) {
		t.Error("failed: Number ':123' should not equal ':124'")
	}
	if !n.Equal(parser.NewJsonNumber("", 123)) {
		t.Error("failed: Number ':123' should equal ':123'")
	}
}
