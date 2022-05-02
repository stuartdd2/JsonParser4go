package test

import (
	"fmt"
	"testing"

	"github.com/stuartdd2/JsonParser4go/parser"
)

func TestObjectKeySort(t *testing.T) {
	root1, err := parser.Parse(obj2)
	if err != nil {
		t.Errorf("failed: to parse obj5: %s", err.Error())
		return
	}
	o, err := parser.Find(root1, parser.NewDotPath("address"))
	if err != nil {
		t.Errorf("failed: to parse obj5: %s", err.Error())
		return
	}
	x := o.(*parser.JsonObject).GetSortedKeys()
	if fmt.Sprintf("%s", x) != "[business city phoneNumbers state streetAddress]" {
		t.Errorf("failed: tosort: %s", x)
		return
	}
}

func TestObjectDataSort(t *testing.T) {
	root1, err := parser.Parse(obj2)
	if err != nil {
		t.Errorf("failed: to parse obj5: %s", err.Error())
		return
	}
	o, err := parser.Find(root1, parser.NewDotPath("address"))
	if err != nil {
		t.Errorf("failed: to parse obj5: %s", err.Error())
		return
	}
	x := o.(*parser.JsonObject).GetValuesSorted()
	if fmt.Sprintf("%s", x) != "[true San Diego \"phoneNumbers\": [{\"type\": \"home\",\"number\": \"7349282382\"}] CA 101]" {
		t.Errorf("failed: tosort: %s", x)
		return
	}
}
