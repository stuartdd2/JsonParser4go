package test

import (
	"testing"

	"github.com/stuartdd2/JsonParser4go/parser"
)

func TestParentOfObjectsInRoot(t *testing.T) {
	root1, err := parser.Parse(singleListData)
	if err != nil {
		t.Errorf("failed: to parse singleListData: %s", err.Error())
		return
	}
	if root1.GetParent() != nil {
		t.Errorf("failed: Rootnode parent must be nil")
		return
	}
	if root1 != root1.GetNodeWithName("config").GetParent() {
		t.Errorf("failed: config nod must have root as parent")
		return
	}
	if root1 != root1.GetNodeWithName("actions").GetParent() {
		t.Errorf("failed: actions nod must have root as parent")
		return
	}
}

func TestParentOfObjectsInList(t *testing.T) {
	root1, err := parser.Parse(singleListData)
	if err != nil {
		t.Errorf("failed: to parse singleListData: %s", err.Error())
		return
	}
	if root1.GetParent() != nil {
		t.Errorf("failed: Rootnode parent must be nil")
		return
	}
	list, err := parser.Find(root1, parser.NewDotPath("actions.0.list"))
	if err != nil {
		t.Errorf("failed: to find actions.0.list: %s", err.Error())
		return
	}
	for _, o := range list.(*parser.JsonList).GetValues() {
		if o.GetParent() != list {
			t.Errorf("failed: object: %s does not have list as a parent", err.Error())
		}
	}
}
func TestParentOfObjectsInObject(t *testing.T) {
	root1, err := parser.Parse(singleListData)
	if err != nil {
		t.Errorf("failed: to parse singleListData: %s", err.Error())
		return
	}
	list, err := parser.Find(root1, parser.NewDotPath("actions.0"))
	if err != nil {
		t.Errorf("failed: to find actions.0: %s", err.Error())
		return
	}
	for _, o := range list.(*parser.JsonObject).GetValues() {
		if o.GetParent() != list {
			t.Errorf("failed: object: %s does not have list as a parent", err.Error())
		}
	}
}

func TestParentOfObjectAndNumber(t *testing.T) {
	root1, err := parser.Parse(singleListData)
	if err != nil {
		t.Errorf("failed: to parse singleListData: %s", err.Error())
		return
	}
	obj, err := parser.Find(root1, parser.NewDotPath("actions.0"))
	if err != nil {
		t.Errorf("failed: to find actions: %s", err.Error())
		return
	}

	nm := parser.NewJsonNumber("number", 10)
	obj.(*parser.JsonObject).Add(nm)

	if nm.GetParent() != obj {
		t.Errorf("failed: to set parent of Number: %s", err.Error())
	}
}

func TestParentOfObjectAndString(t *testing.T) {
	root1, err := parser.Parse(singleListData)
	if err != nil {
		t.Errorf("failed: to parse singleListData: %s", err.Error())
		return
	}
	list, err := parser.Find(root1, parser.NewDotPath("actions.0.list"))
	if err != nil {
		t.Errorf("failed: to find actions.0.list: %s", err.Error())
		return
	}

	ob := parser.NewJsonObject("container")
	st := parser.NewJsonString("string", "value")
	_, err = ob.Add(st)
	if err != nil {
		t.Errorf("failed: to add node to object: %s", err.Error())
		return
	}

	_, err = list.(*parser.JsonList).Add(ob)
	if err != nil {
		t.Errorf("failed: to add node to list: %s", err.Error())
		return
	}

	if ob.GetParent() != list {
		t.Errorf("failed: to set parent of Object: %s", err.Error())
	}
	if st.GetParent() != ob {
		t.Errorf("failed: to set parent of String: %s", err.Error())
	}
}
