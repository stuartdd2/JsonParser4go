package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stuartdd2/JsonParser4go/parser"
)

func TestWalkWithTrailForFile(t *testing.T) {
	rLogH("TestWalkWithTrailForFile")
	root := InitParserFromFile(t, "TestDataTypes.json")
	groups, _ := parser.Find(root, parser.NewDotPath("groups"))
	collect := make(map[string]string)
	var li strings.Builder
	parser.WalkNodeTreeForTrail(groups.(parser.NodeC), func(trail *parser.Trail, i int) bool {
		s := trail.String()
		rLog(s)
		_, ok := collect[s]
		if ok {
			li.WriteString(fmt.Sprintf("Duplicate:%s\n", s))
		} else {
			collect[s] = s
		}
		return false
	})
	if li.Len() > 0 {
		t.Errorf(li.String())
	}
	if len(collect) != 78 {
		t.Errorf("Should be 78 nodes not %d", len(collect))
	}
}

func TestWalkNodeTreeForTrail3(t *testing.T) {
	rLogH("TestWalkNodeTreeForTrail3")
	root := InitParser(t, "", obj3)
	testWNTFP(t, root, parser.NewBarPath("address|business"), "true")
	testWNTFP(t, root, parser.NewBarPath("address|phoneNumbers|1|number"), "7349282382")
	testWNTFP(t, root, parser.NewBarPath("address|phoneNumbers|5|no"), "null")
}

func TestWalkNodeTreeForTrail4(t *testing.T) {
	rLogH("TestWalkNodeTreeForTrail")
	root := InitParser(t, "", obj4)
	testWNTFP(t, root, parser.NewBarPath("list1|3|list3|2"), "20")
	testWNTFP(t, root, parser.NewBarPath("list1|2|lastName"), "Jackson")
	testWNTFP(t, root, parser.NewBarPath("list1|3|A"), "10")
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

	p1 := target.GetParent()
	CheckNodeParentTarget(t, p1, target, "f", true, "")

	p2 := target.GetParent()
	CheckNodeParentTarget(t, p2, p1, "", true, "")

}

func TestFindNodeInListComplex(t *testing.T) {
	root := InitParser(t, "", obj3)
	target, _ := parser.Find(root, parser.NewBarPath("address|phoneNumbers|1|number"))
	p1 := target.GetParent()

	CheckNodeParentTarget(t, p1, target, "number", true, "phoneNumbers")

	target, _ = parser.Find(root, parser.NewDotPath("address.phoneNumbers.no"))
	p1 = target.GetParent()
	CheckNodeParentTarget(t, p1, target, "no", true, "phoneNumbers")

}

func TestFindNodeInObjects(t *testing.T) {
	root := InitParser(t, "", obj2)
	target, _ := parser.Find(root, parser.NewPath("address.state", "."))
	as1 := target.String()

	p := target.GetParent()
	CheckNodeParentTarget(t, p, target, "state", true, "address")
	target, _ = parser.Find(root, parser.NewPath("address|phoneNumbers|0|number", "|"))
	p = target.GetParent()
	CheckNodeParentTarget(t, p, target, "number", true, "")
	target, _ = parser.Find(root, parser.NewPath("address.phoneNumbers", "."))
	p = target.GetParent()
	CheckNodeParentTarget(t, p, target, "phoneNumbers", true, "address")
	target, _ = parser.Find(root, parser.NewPath("age", "."))
	p = target.GetParent()
	CheckNodeParentTarget(t, p, target, "age", true, "")
	p = target.GetParent()
	CheckNodeParentTarget(t, p, root, "", false, "")

	target = parser.NewJsonString("state", "CA")
	as2 := target.String()

	if as2 != as1 {
		t.Errorf("Node address.state and external node address.sate string should be the same '%s' != '%s", as1, as2)
	}
	p = target.GetParent()
	if p != nil {
		t.Errorf("returned parent node should be nil")
	}
}

func CheckNodeParentTarget(t *testing.T, parent, target parser.NodeI, expected string, shouldHaveParent bool, parentStr string) {
	if shouldHaveParent {
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
