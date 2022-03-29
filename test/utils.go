package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stuartdd2/JsonParser4go/parser"
)

var (
	nList1 = []byte(`["literal", {"obj":"literal"}, {"num":99.9}, {"t":true}, {"f":false}]`)
	nList2 = []byte(`["literal", {"obj":"literal"}, {"num":99.9}, {"t":true}, {"f":false, "t":true}]`)
	obj1   = []byte(`{
		"firstName": "Joe",
		"lastName": "Jackson",
		"gender": "male"
	 }`)
	obj2 = []byte(`{
		"firstName": "Joe",
		"lastName": "Jackson",
		"gender": "male",
		"age": 28,
		"address": {
			"streetAddress": "101",
			"city": "San Diego",
			"state": "CA",
			"business": true,
			"phoneNumbers": [{ "type": "home", "number": "7349282382" }]
		}
	 }`)
	obj3 = []byte(`{
		"firstName": "Joe",
		"lastName": "Jackson",
		"gender": "male",
		"age": 28,
		"bo": true,
		"no": null,
		"address": {
			"streetAddress": "101",
			"city": "San Diego",
			"state": "CA",
			"business": true,
			"phoneNumbers": [{
				"string": "home"
			}, {
				"number": 7349282382
			}, {
				"boolean": true
			}, {
				"dupe1": false
			}, {
				"dupe1": 123456
			}, {
				"no": null
			}]
		}
	}`)
	obj4 = []byte(`{
		"list1": [
			"Joe",
			null,
			{
				"lastName": "Jackson"
			},
			{
				"A": 10,
				"B": true,
				"C": null,
				"list3": [
					true,
					null,
					20
				],
				"D": false
			}
		]
	}`)
	obj5 = []byte(`{
		"list1": ["Joe",{"lastName": "Jackson"},{"A":10, "B":true, "C":null, "list3":[true,null,20]},"Jim"],
		"gender": "male",
		"list2": [{"lastName": "Jackson"},"Jack",{"streetAddress": "101"}],
		"age": 28,
		"bo": true,
		"no": null,
		"address": {
			"streetAddress": "101",
			"business": true,
			"list2": [{
				"string": "home"
			}, {
				"number": 7349282382
			}, {
				"boolean": true
			}, {
				"dupe1": false
			}, {
				"dupe1": 123456
			}, {
				"no": null
			}]
		}
	}`)
)

func InitParser(t *testing.T, sourceName string, dat []byte) parser.NodeC {
	node, err := parser.Parse(dat)
	if err != nil {
		t.Errorf("Failed to parse file %s. Error %s\n", sourceName, err.Error())
		return nil
	}
	return node
}

func CheckErr(t *testing.T, err error, cont string) {
	if err == nil {
		t.Errorf("Test Failed. Err is null. Contains was %s\n", cont)
		return
	}
	if strings.Contains(strings.ToLower(err.Error()), strings.ToLower(cont)) {
		return
	}
	t.Errorf("Test Failed. Error does not contain '%s'. Error %s\n", cont, err.Error())
}

func CheckFindNode(t *testing.T, node parser.NodeI, path string, cont string) parser.NodeI {
	if node == nil {
		t.Errorf("Test Failed. Node is nil")
		return nil
	}
	n, err := parser.Find(node, parser.NewDotPath(path))
	if err != nil {
		t.Errorf("Should NOT have thrown %s", err.Error())
		return nil
	}
	if CheckContains(t, n, cont) {
		return n
	}
	return nil
}

func CheckContains(t *testing.T, node parser.NodeI, cont string) bool {
	if node == nil {
		t.Errorf("n should never be nil")
		return false
	}
	if strings.Contains(strings.ToLower(node.JsonValue()), strings.ToLower(cont)) {
		return true
	}
	t.Errorf("Test Failed. Node JsonValue() does not contain '%s'. Node %s\n", cont, node.JsonValue())
	return false
}

func Diagnostic(t *testing.T, node parser.NodeI) {
	t.Errorf(parser.DiagnosticList(node))
	t.Errorf(node.String())
}

func CompareTrees(n1, n2 parser.NodeI) {
	if n1 == nil && n2 != nil {
		panic(fmt.Sprintf("CompareTrees - Node 1 is nil != Node 2 '%s' is not", n2.GetName()))
	}
	if n2 == nil && n1 != nil {
		panic(fmt.Sprintf("CompareTrees - Node 2 is nil != Node 1 '%s' is not", n1.GetName()))
	}
	if n1.IsContainer() {
		if n1.(parser.NodeC).Len() != n2.(parser.NodeC).Len() {
			panic(fmt.Sprintf("CompareTrees - Node 1 '%s' len %d != Node 2 '%s' len %d", n1.GetName(), n1.(parser.NodeC).Len(), n2.GetName(), n2.(parser.NodeC).Len()))
		}
		if n1.GetNodeType() == parser.NT_LIST {
			for i, v1 := range n1.(parser.NodeC).GetValues() {
				v2 := n2.(*parser.JsonList).GetNodeAt(i)
				CompareTrees(v1, v2)
			}
		} else {
			for _, v1 := range n1.(parser.NodeC).GetValues() {
				v2 := n2.(*parser.JsonObject).GetNodeWithName(v1.GetName())
				CompareTrees(v1, v2)
			}
		}
	} else {
		compareNode(n1, n2)
	}
}

func compareNode(n1, n2 parser.NodeI) {
	if n1.GetNodeType() != n2.GetNodeType() {
		panic(fmt.Sprintf("CompareTrees - Node 1 '%s' type %s != Node 2 '%s' type %s", n1.GetName(), parser.GetNodeTypeName(n1.GetNodeType()), n2.GetName(), parser.GetNodeTypeName(n2.GetNodeType())))
	}
	if n1.String() != n2.String() {
		panic(fmt.Sprintf("CompareTrees - Node 1 '%s' value %s != Node 2 '%s' value %s", n1.GetName(), n1.String(), n2.GetName(), n2.String()))
	}

}
