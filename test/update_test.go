package test

import (
	"github.com/stuartdd/jsonParserGo/parser"
	"testing"
)

func TestListRemoveNodeInObjects(t *testing.T) {
	root := InitParser(t, "", obj2)
	n := CheckFindNode(t, root, "address.phoneNumbers", "7349282382")
	parent, ok := parser.FindParentNode(root, n)
	if !ok || parent == nil {
		t.Errorf("FindParentNode failed to find the node or the parent is nil")
	}
	if n != nil {
		parser.Remove(root, n)
		_, ok := parser.FindParentNode(root, n)
		if ok {
			t.Errorf("FindParentNode found the removed node")
		}
	} else {
		t.Errorf("FindParentNode failed to find the node")
	}
}

func TestListRemoveNodeInObjectsInList(t *testing.T) {
	root := InitParser(t, "", obj2)
	n := CheckFindNode(t, root, "address.phoneNumbers.0.type", "home")
	parent, ok := parser.FindParentNode(root, n)
	if !ok || parent == nil {
		t.Errorf("FindParentNode failed to find the node or the parent is nil")
	}
	if n != nil {
		parser.Remove(root, n)
		_, ok := parser.FindParentNode(root, n)
		if ok {
			t.Errorf("FindParentNode found the removed node")
		}
	} else {
		t.Errorf("FindParentNode failed to find the node")
	}
}
func TestListRemoveObjectInObjects(t *testing.T) {
	root := InitParser(t, "", obj2)
	n := CheckFindNode(t, root, "address", "\"streetAddress\": \"101\"")
	parent, ok := parser.FindParentNode(root, n)
	if !ok || parent == nil {
		t.Errorf("FindParentNode failed to find the node or the parent is nil")
	}
	if n != nil {
		parser.Remove(root, n)
		_, ok := parser.FindParentNode(root, n)
		if ok {
			t.Errorf("FindParentNode found the removed node")
		}
	} else {
		t.Errorf("FindParentNode failed to find the node")
	}
}
func TestListRemoveNodeInList(t *testing.T) {
	//	nList1 = []byte(`["literal", {"obj":"literal"}, {"num":99.9}, {"t":true}, {"f":false}]`)
	root := InitParser(t, "", nList1)
	n := CheckFindNode(t, root, "obj", "literal")
	parent, ok := parser.FindParentNode(root, n)
	if !ok || parent == nil {
		t.Errorf("FindParentNode failed to find the node or the parent is nil")
	}
	if n != nil {
		parser.Remove(root, n)
	}
	p, ok := parser.FindParentNode(root, n)
	if ok {
		t.Errorf("FindParentNode should NOT find the node")
	}
	if p != nil {
		t.Errorf("FindParentNode should NOT find the node so parent should be null")
	}
	rootL := (root).(*parser.JsonList)
	newNode := parser.NewJsonString("Fred", "Jones")
	rootL.Add(newNode)
	listP, ok := parser.FindParentNode(root, newNode)
	if !ok || listP == nil {
		t.Errorf("FindParentNode failed to find the NEW node or the parent is nil")
	}
}
func TestUpdate(t *testing.T) {
	node := InitParser(t, "", obj2)

	n := CheckFindNode(t, node, "address.phoneNumbers.0.type", "home")
	if n != nil {
		ns := (n).(*parser.JsonString)
		ns.SetValue("away")
		CheckFindNode(t, node, "address.phoneNumbers.0.type", "away")
	}

	n = CheckFindNode(t, node, "age", "28")
	if n != nil {
		ns := (n).(*parser.JsonNumber)
		ns.SetValue(99)
		CheckFindNode(t, node, "age", "99")
	}

	n = CheckFindNode(t, node, "address.business", "true")
	if n != nil {
		ns := (n).(*parser.JsonBool)
		ns.SetValue(false)
		CheckFindNode(t, node, "address.business", "false")
	}

	n = CheckFindNode(t, node, "address.phoneNumbers", "away")
	if n != nil {
		ns := (n).(*parser.JsonList)
		ns.Add(parser.NewJsonNumber("N2", 12345))
		CheckFindNode(t, node, "address.phoneNumbers.N2", "\"N2\": 12345")
	}

	n = CheckFindNode(t, node, "address.phoneNumbers.N2", "\"N2\": 12345")
	if n != nil {
		ns := (n).(*parser.JsonNumber)
		ns.SetValue(99999)
		CheckFindNode(t, node, "address.phoneNumbers.N2", "\"N2\": 99999")
	}

}
