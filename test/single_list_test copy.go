package test

import (
	"testing"

	"github.com/stuartdd2/JsonParser4go/parser"
)

func TestSingleList(t *testing.T) {
	root1, err := parser.Parse(singleListData)
	if err != nil {
		t.Errorf("failed: to parse singleListData: %s", err.Error())
		return
	}
	list, err := parser.Find(root1, parser.NewDotPath("actions.0.list"))
	if err != nil {
		t.Errorf("failed: to find list: %s", err.Error())
		return
	}
	for _, v := range list.(*parser.JsonList).GetValues() {
		if v.GetNodeType() != parser.NT_OBJECT {
			t.Errorf("Node %s=%s should be an object\n", v.GetName(), v.String())
		}
	}
}
