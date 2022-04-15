package test

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stuartdd2/JsonParser4go/parser"
)

func rLog(s string) {
	http.Post("http://localhost:9998/log", "text/plain", bytes.NewBufferString(time.Now().Format("15:04:05.000 ")+s))
}

func rLogH(s string) {
	rLog(fmt.Sprintf("\n***\n*** Run Test %s\n***", s))
}

func rLogE(t *testing.T, s string) {
	rLog(s)
	t.Errorf(s)
}

func TestWalkNodeTreeForTrail3(t *testing.T) {
	root := InitParser(t, "", obj3)
	testWNTFP(t, root, parser.NewBarPath("address|business"), "true")
	testWNTFP(t, root, parser.NewBarPath("address|phoneNumbers|number"), "7349282382")
	testWNTFP(t, root, parser.NewBarPath("address|phoneNumbers|no"), "null")
}

func TestWalkNodeTreeForTrail4(t *testing.T) {
	rLogH(" TestWalkNodeTreeForTrail")
	root := InitParser(t, "", obj4)
	testWNTFP(t, root, parser.NewBarPath("list1||list3|2"), "20")
	// testWNTFP(t, root, parser.NewBarPath("list1|lastName"), "Jackson")
	//testWNTFP(t, root, parser.NewBarPath("list1|3|A"), "10")
}

func testWNTFP(t *testing.T, root parser.NodeC, req *parser.Path, val string) {
	trail, ok := parser.WalkNodeTreeForTrail(root, func(trail *parser.Trail, index int) bool {
		s := trail.String()
		if s == req.String() {
			rLog(fmt.Sprintf("HIT  %s == %s", s, req.String()))
			return true
		}
		rLog(fmt.Sprintf("MISS %s == %s", s, req.String()))
		return false
	})
	if !ok {
		rLogE(t, fmt.Sprintf("WalkNodeTreeForPath: could not find %s", req.String()))
		return
	}
	v := trail.GetLast()
	if v.String() != val && val != "*" {
		rLogE(t, fmt.Sprintf("WalkNodeTreeForPath: incorrect value returned for last node name [%s]. '%s' != '%s'", v.GetName(), v.String(), val))
	}
}

func TestNoCopyOnFind(t *testing.T) {
	root := InitParser(t, "", obj3)
	n := parser.NewJsonString("Name", "Fred")
	if n.GetName() != "Name" {
		t.Errorf("Node must be the Name node")
	}
	if n.String() != "Fred" {
		t.Errorf("Node must be the 'Name' node with value 'Fred'")
	}
	// add it and find it again!
	root.Add(n)
	nf, err := parser.Find(root, parser.NewDotPath("Name"))
	if err != nil {
		t.Errorf("Node 'Name' not found")
	}
	if nf.GetName() != "Name" {
		t.Errorf("Node must be the 'Name' node")
	}
	if nf.String() != "Fred" {
		t.Errorf("Node must be the 'Name' node with value 'Fred'")
	}

	// Chage value of original node
	n.SetValue("Freda")
	if n.GetName() != "Name" {
		t.Errorf("Node must be the Name node")
	}
	if n.String() != "Freda" {
		t.Errorf("Node must be the 'Name' node with value 'Freda'")
	}
	if nf.GetName() != "Name" {
		t.Errorf("Node must be the 'Name' node")
	}
	if nf.String() != "Freda" {
		t.Errorf("Node must be the 'Name' node with value 'Freda'")
	}

	nf2, err2 := parser.Find(root, parser.NewDotPath("Name"))
	if err2 != nil {
		t.Errorf("Node 'Name' not found")
	}
	if nf2.GetName() != "Name" {
		t.Errorf("Node must be the Name node")
	}
	if nf2.String() != "Freda" {
		t.Errorf("Node must be the Name node with value Freda")
	}

	nf2.(*parser.JsonString).SetValue("Boggle")
	if n.GetName() != "Name" {
		t.Errorf("Node must be the Name node")
	}
	if n.String() != "Boggle" {
		t.Errorf("Node must be the Name node with value Boggle")
	}
}

func TestWalkNodesUntillFound(t *testing.T) {
	root := InitParser(t, "", obj3)
	target, _ := parser.Find(root, parser.NewPath("address|streetAddress", "|"))
	node, parent, ok := parser.WalkNodeTree(root, target, func(node, parent, target parser.NodeI) bool {
		return node.GetName() == target.GetName()
	})
	if !ok {
		t.Errorf("Did not find target in the map")
	} else {
		if parent.GetName() != "address" {
			t.Errorf("Parent node should hve name 'address'. Actual:%s", parent.GetName())
		}
		if node.GetName() != "streetAddress" {
			t.Errorf("Parent node should hve name 'streetAddress'. Actual:%s", node.GetName())
		}
	}
}

func TestWalkAllNodes(t *testing.T) {
	m := make(map[string]string)
	root := InitParser(t, "", obj3)
	target, _ := parser.Find(root, parser.NewPath("address|phoneNumbers", "|"))
	parser.WalkNodeTree(root, target, func(node, parent, target parser.NodeI) bool {
		m[node.GetName()] = node.GetName()
		return false
	})
	var sb strings.Builder

	testStr := "no address city state business phoneNumbers string number boolean dupe1 dupe1 no streetAddress firstName lastName gender age bo"
	testList := strings.Split(testStr, " ")
	for _, v := range testList {
		_, ok := m[v]
		if !ok {
			t.Errorf("Did not find '%s' in the map", v)
		}
		sb.WriteString(v)
		sb.WriteString(" ")
	}
	if strings.Trim(sb.String(), " ") != testStr {
		t.Errorf("Strings should be the same \ntestStr:'%s'\nactual :'%s'", testStr, sb.String())
	}
}
func TestFindNodeInList(t *testing.T) {
	root := InitParser(t, "", nList2)
	target, _ := parser.Find(root, parser.NewPath("4.f", "."))

	p1, ok := parser.FindParentNode(root, target)
	CheckNodeParentTarget(t, p1, target, ok, "f", true, "")

	p2, ok := parser.FindParentNode(root, p1)
	CheckNodeParentTarget(t, p2, p1, ok, "", true, "")

	p3, ok := parser.FindParentNode(root, p2)
	CheckNodeParentTarget(t, p3, p2, ok, "", false, "")

}

func TestFindNodeInListComplex(t *testing.T) {
	root := InitParser(t, "", obj3)
	target, _ := parser.Find(root, parser.NewBarPath("address|phoneNumbers|number"))
	p1, ok := parser.FindParentNode(root, target)
	CheckNodeParentTarget(t, p1, target, ok, "number", true, "phoneNumbers")

	target, _ = parser.Find(root, parser.NewDotPath("address.phoneNumbers.no"))
	p1, ok = parser.FindParentNode(root, target)
	CheckNodeParentTarget(t, p1, target, ok, "no", true, "phoneNumbers")

}

func TestFindNodeInObjects(t *testing.T) {
	root := InitParser(t, "", obj2)
	target, _ := parser.Find(root, parser.NewPath("address.state", "."))
	as1 := target.String()

	p, ok := parser.FindParentNode(root, target)
	CheckNodeParentTarget(t, p, target, ok, "state", true, "address")
	target, _ = parser.Find(root, parser.NewPath("address|phoneNumbers|0|number", "|"))
	p, ok = parser.FindParentNode(root, target)
	CheckNodeParentTarget(t, p, target, ok, "number", true, "")
	target, _ = parser.Find(root, parser.NewPath("address.phoneNumbers", "."))
	p, ok = parser.FindParentNode(root, target)
	CheckNodeParentTarget(t, p, target, ok, "phoneNumbers", true, "address")
	target, _ = parser.Find(root, parser.NewPath("age", "."))
	p, ok = parser.FindParentNode(root, target)
	CheckNodeParentTarget(t, p, target, ok, "age", true, "")
	p, ok = parser.FindParentNode(root, root)
	CheckNodeParentTarget(t, p, root, ok, "", false, "")

	target = parser.NewJsonString("state", "CA")
	as2 := target.String()

	if as2 != as1 {
		t.Errorf("Node address.state and external node address.sate string should be the same '%s' != '%s", as1, as2)
	}
	p, ok = parser.FindParentNode(root, target)
	if ok {
		t.Errorf("Just because the String() returns the same doesent make then the same!")
	}
	if p != nil {
		t.Errorf("returned parent node should be nil")
	}
}

func CheckNodeParentTarget(t *testing.T, parent, target parser.NodeI, found bool, expected string, shouldHaveParent bool, parentStr string) {
	if shouldHaveParent {
		if !found {
			t.Errorf("returned node was not found. Expected '%s'", expected)
		}
		if parent == nil {
			t.Errorf("Returned parent is nil, Expected: '%s'", parentStr)
		}
		if parent.GetName() != parentStr {
			t.Errorf("Returned parent incorrect name '%s', Expected: '%s'", parent.GetName(), parentStr)
		}
		if target.GetName() != expected {
			t.Errorf("Returned target incorrect name '%s', Expected: '%s'", target.GetName(), expected)
		}
	} else {
		if parent != nil {
			t.Errorf("Returned parent should be nil")
		}
	}
}

func TestFindInList1(t *testing.T) {
	node := InitParser(t, "", nList1)
	CheckFindNode(t, node, "obj", "literal")
	CheckFindNode(t, node, "num", "99.9")
}

func TestFindInMap2(t *testing.T) {
	node := InitParser(t, "", obj2)
	CheckFindNode(t, node, "address.phoneNumbers.0.type", "home")
	CheckFindNode(t, node, "address.phoneNumbers.0.number", "7349282382")
	CheckFindNode(t, node, "address.phoneNumbers.0", "7349282382")
	CheckFindNode(t, node, "firstName", "Joe")
	CheckFindNode(t, node, "age", "28")
	CheckFindNode(t, node, "address.state", "CA")
}

func TestFindInMap1(t *testing.T) {
	node := InitParser(t, "", obj1)
	CheckFindNode(t, node, "lastName", "Jackson")
}

func TestFindNotInMap(t *testing.T) {
	node := InitParser(t, "", obj1)
	_, err := parser.Find(node, parser.NewDotPath("5"))
	if err == nil {
		t.Errorf("Find test failed. Should have thrown err")
	}
	CheckErr(t, err, "was not found")
}

func TestFindFirstEle(t *testing.T) {
	node := InitParser(t, "", nList1)
	CheckFindNode(t, node, "0", "literal")
	CheckFindNode(t, node, "4", "false")
}

func TestFindNodeHighIndexInList(t *testing.T) {
	node := InitParser(t, "", nList1)
	_, err := parser.Find(node, parser.NewDotPath("5"))
	if err == nil {
		t.Errorf("Find test failed. Should have thrown err")
	}
	CheckErr(t, err, "Index out of bounds")
}

func TestFindNodeNegativeIndexInList(t *testing.T) {
	node := InitParser(t, "", nList1)
	_, err := parser.Find(node, parser.NewDotPath("-1"))
	if err == nil {
		t.Errorf("Find test failed. Should have thrown err")
	}
	CheckErr(t, err, "Index out of bounds")
}

func TestFindNodeInvalidIndexInList(t *testing.T) {
	node := InitParser(t, "", nList1)
	_, err := parser.Find(node, parser.NewDotPath("x"))
	if err == nil {
		t.Errorf("Find test failed. Should have thrown err")
	}
	CheckErr(t, err, "element: 'x' was not found")
}

func TestFindNodeEmptySearch(t *testing.T) {
	node := InitParser(t, "", nList1)
	_, err := parser.Find(node, parser.NewPath("", ""))
	if err == nil {
		t.Errorf("Find test failed. Should have thrown err")
	}
	CheckErr(t, err, "No search paths")
}
