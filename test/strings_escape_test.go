package test

import (
	"github.com/stuartdd/jsonParserGo/parser"
	"testing"
)

var (
	testNoEscInn = []byte(`{"str":"No Escape String"}`)
	testNoEscOut = []byte(`No Escape String`)

	testEscInn = []byte(`{"str":"No\nEsc\\ape\tSt\"ring\""}`)
	testEscOut = []byte(`No
Esc\ape	St"ring"`)
	testCodeInn3C         = []byte(`{"str":"\u003C"}`)
	testCodeOut3C         = []byte(`<`)
	testCodeInnBlackStars = []byte(`{"str":"\u2605\u2605\u2605"}`)
	testCodeOutBlackStars = []byte(`★★★`)
	testCodeInnX3C        = []byte(`{"str":"\x3C"}`)
	testCodeOutX3C        = []byte(`<`)
	testCodeInnXX         = []byte(`{"str":"some_\u00fc\u00f1\u00eec\u00f8d\u00e9_and_\u007f"}`)
	testCodeOutXX         = []byte(`some_üñîcødé_and_`)
)

func TestCodeXX(t *testing.T) {
	p, err := parser.Parse(testCodeInnXX)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	assertExact(t, "01", p, "str", string(testCodeOutXX))
}

func TestCodeX3C(t *testing.T) {
	p, err := parser.Parse(testCodeInnX3C)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	assertExact(t, "01", p, "str", string(testCodeOutX3C))
}

func TestCodeBlackStars(t *testing.T) {
	p, err := parser.Parse(testCodeInnBlackStars)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	assertExact(t, "01", p, "str", string(testCodeOutBlackStars))
}

func TestCode3C(t *testing.T) {
	p, err := parser.Parse(testCodeInn3C)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	assertExact(t, "01", p, "str", string(testCodeOut3C))
}

func TestEsc(t *testing.T) {
	p, err := parser.Parse(testEscInn)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	assertExact(t, "01", p, "str", string(testEscOut))
}

func TestNoEsc(t *testing.T) {
	p, err := parser.Parse(testNoEscInn)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	assertExact(t, "01", p, "str", string(testNoEscOut))
}

func assertExact(t *testing.T, id string, n parser.NodeI, name, s2 string) {
	if n.GetNodeType() != parser.NT_OBJECT {
		t.Errorf("assertExact:%s FAIL: node is not an object", id)
	}
	sn := n.(*parser.JsonObject).GetNodeWithName(name)
	if sn == nil {
		t.Errorf("assertExact:%s FAIL: cannot find %s", id, name)
	}
	s := sn.String()
	if s != s2 {
		t.Errorf("assertExact:%s FAIL: Strings do not match\nActual  :%s\nExpected:%s", id, s, s2)
	}
}
