package test

import (
	"testing"

	"github.com/stuartdd2/JsonParser4go/parser"
)

func TestListRemoveNodeInObjects(t *testing.T) {
	root := InitParser(t, "", obj2)
	n := CheckFindNode(t, root, "address.phoneNumbers", "7349282382")
	parent := n.GetParent()
	if parent == nil {
		t.Errorf("FindParentNode failed to find the node or the parent is nil")
	}
	if n != nil {
		parser.Remove(n)
		if n.GetParent() != nil {
			t.Errorf("Removed node still has a parent")
		}
	} else {
		t.Errorf("Node to remove has no parent")
	}
}

func TestListRemoveNodeInObjectsInList(t *testing.T) {
	root := InitParser(t, "", obj2)
	n := CheckFindNode(t, root, "address.phoneNumbers.0.type", "home")
	parent := n.GetParent()
	if parent == nil {
		t.Errorf("FindParentNode failed to find the node or the parent is nil")
	}
	if n != nil {
		parser.Remove(n)
		if n.GetParent() != nil {
			t.Errorf("removed node still has parent")
		}
	} else {
		t.Errorf("Node to remove has no parent")
	}
}
func TestListRemoveObjectInObjects(t *testing.T) {
	root := InitParser(t, "", obj2)
	n := CheckFindNode(t, root, "address", "\"streetAddress\": \"101\"")
	parent := n.GetParent()
	if parent == nil {
		t.Errorf("Node has no parent")
	}
	if n != nil {
		parser.Remove(n)
		if n.GetParent() != nil {
			t.Errorf("removed node still has parent")
		}
	} else {
		t.Errorf("Node to remove has no parent")
	}
}
func TestListRemoveNodeInList(t *testing.T) {
	//	nList1 = []byte(`["literal", {"obj":"literal"}, {"num":99.9}, {"t":true}, {"f":false}]`)
	root := InitParser(t, "", nList1)
	n := CheckFindNode(t, root, "obj", "literal")
	parent := n.GetParent()
	if parent == nil {
		t.Errorf("Node to remove has no parent")
	}
	if n != nil {
		parser.Remove(n)
	}
	p := n.GetParent()
	if p != nil {
		t.Errorf("removed node still has parent")
	}
	rootL := (root).(*parser.JsonList)
	newNode := parser.NewJsonString("Fred", "Jones")
	rootL.Add(newNode)
	listP := newNode.GetParent()
	if listP == nil {
		t.Errorf("New Node should have a parent")
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
