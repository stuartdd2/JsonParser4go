package test

import (
	"github.com/stuartdd/jsonParserGo/parser"
	"strings"
	"testing"
)

// func TestObjectTab(t *testing.T) {
// 	n, err := parser.Parse(obj5)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	t.Errorf("\n%s", n.JsonValue())
// 	t.Errorf("%s", n.JsonValueIndented(4))
// }

func TestNumberTab(t *testing.T) {
	expected := "\"name\": 123.5"
	n := parser.NewJsonNumber("name", 123.5)
	testTabValue(t, 0, expected, n.JsonValueIndented(0))
	testTabValue(t, 1, expected, n.JsonValueIndented(1))
	testTabValue(t, 2, expected, n.JsonValueIndented(2))
	n = parser.NewJsonNumber("", 100.0)
	testTabValue(t, 3, "100", n.JsonValueIndented(3))
	n = parser.NewJsonNumber("", 0)
	testTabValue(t, 3, "0", n.JsonValueIndented(3))
	n = parser.NewJsonNumber("", 0.999)
	testTabValue(t, 3, "0.999", n.JsonValueIndented(3))
}

func TestBoolTab(t *testing.T) {
	expected := "\"name\": true"
	n := parser.NewJsonBool("name", true)
	testTabValue(t, 0, expected, n.JsonValueIndented(0))
	testTabValue(t, 1, expected, n.JsonValueIndented(1))
	testTabValue(t, 2, expected, n.JsonValueIndented(2))
	testTabValue(t, 4, expected, n.JsonValueIndented(4))
	testTabValue(t, 6, expected, n.JsonValueIndented(6))
	n = parser.NewJsonBool("", true)
	testTabValue(t, 3, "true", n.JsonValueIndented(3))
	n = parser.NewJsonBool("", false)
	testTabValue(t, 3, "false", n.JsonValueIndented(3))
}

func TestNullTab(t *testing.T) {
	expected := "\"name\": null"
	n := parser.NewJsonNull("name")
	testTabValue(t, 0, expected, n.JsonValueIndented(0))
	testTabValue(t, 1, expected, n.JsonValueIndented(1))
	testTabValue(t, 2, expected, n.JsonValueIndented(2))
	n = parser.NewJsonNull("")
	testTabValue(t, 3, "null", n.JsonValueIndented(3))
}

func TestStringTab(t *testing.T) {
	expected := "\"name\": \"value\""
	n := parser.NewJsonString("name", "value")
	testTabValue(t, 0, expected, n.JsonValueIndented(0))
	testTabValue(t, 1, expected, n.JsonValueIndented(1))
	testTabValue(t, 2, expected, n.JsonValueIndented(2))
	n = parser.NewJsonString("", "value")
	testTabValue(t, 4, "\"value\"", n.JsonValueIndented(4))
	n = parser.NewJsonString("", "")
	testTabValue(t, 2, "\"\"", n.JsonValueIndented(2))
}

func testTabValue(t *testing.T, tabLen int, expected, actual string) {
	tStart := 0
	if tabLen > 0 {
		if !strings.HasPrefix(actual, "\n") {
			t.Errorf("Tab len incorrect. 1st char should be newline. value: '%s'", actual)
			return
		}
		tStart = len("\n")
	}
	for i, c := range actual[tStart:] {
		if c != ' ' {
			if i != tabLen {
				t.Errorf("Tab len incorrect. Expected: %d actual: %d. value: '%s'", tabLen, i, actual)
				return
			}
			break
		}
	}
	if expected == "" {
		return
	}
	if actual[tStart+tabLen:] != expected {
		t.Errorf("Value is incorrect. Expected: '%s' actual: '%s'", expected, actual[tabLen:])
		return
	}
}
