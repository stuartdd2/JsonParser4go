package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stuartdd2/JsonParser4go/parser"
)

var (
	// Watch out for the '|' separators used to create a []string below.
	// Cannot use a fixed string as the order of NON list nodes can change and that is still correct.
	// See: strings.Split(string(testRootMapWithLists), "|")
	testRootMapWithLists         = []byte(`"OBJ.FL.1": 123|"LI.1": [{"LI.ST.1": "ABC"},{"LI.BO.1": false}]|"LI.2": ["ABC",198,false,{"NUM": 123}]|"OBJ.ST.1": "ABC"|"OBJ.BO.1": true`)
	testRootListWithMapsAndLists = []byte(`{"OBJ.BO.1": true},{"EMPTY": []},{"FRED": ["Help"]}|"BLUE":|"OBJ.1.ST.1": "ABC"|"OBJ.1.BO.1": true|"OBJ.1.FL.1": 123`)
)

func TestEqualString(t *testing.T) {
	n1 := parser.NewJsonString("a", "abc")
	if n1.Equal(parser.NewJsonString("a", "def")) {
		t.Error("failed: String 'a:abc' should not equal String 'a:def'")
	}
	if n1.Equal(parser.NewJsonString("b", "def")) {
		t.Error("failed: String 'a:abc' should not equal String 'b:def'")
	}
	if n1.Equal(parser.NewJsonString("", "abc")) {
		t.Error("failed: String 'a:abc' should not equal String ':abc'")
	}
	if !n1.Equal(parser.NewJsonString("a", "abc")) {
		t.Error("failed: String 'a:abc' should equal String 'a:abc'")
	}
	n := parser.NewJsonString("", "def")
	if n.Equal(parser.NewJsonString("a", "def")) {
		t.Error("failed: String ':def' should not equal String 'a:def'")
	}
	if n.Equal(parser.NewJsonString("", "abc")) {
		t.Error("failed: String ':def' should not equal String ':abc'")
	}
	if !n.Equal(parser.NewJsonString("", "def")) {
		t.Error("failed: String ':def' should equal String ':def'")
	}
	n3 := parser.NewJsonString("a", "true")
	if n3.Equal(parser.NewJsonBool("a", true)) {
		t.Error("failed: String 'a:true' should not equal Bool 'a:true'")
	}
	n4 := parser.NewJsonString("a", "123")
	if n4.Equal(parser.NewJsonNumber("a", 123)) {
		t.Error("failed: String 'a:123' should not equal Number 'a:123'")
	}
}

func TestEqualBool(t *testing.T) {
	n1 := parser.NewJsonBool("a", true)
	if n1.Equal(parser.NewJsonBool("a", false)) {
		t.Error("failed: Bool 'a:true' should not equal Bool 'a:false'")
	}
	if n1.Equal(parser.NewJsonBool("b", false)) {
		t.Error("failed: Bool 'a:true' should not equal Bool 'b:false'")
	}
	if n1.Equal(parser.NewJsonBool("", true)) {
		t.Error("failed: Bool 'a:true' should not equal Bool ':true'")
	}
	if !n1.Equal(parser.NewJsonBool("a", true)) {
		t.Error("failed: Bool 'a:true' should equal Bool 'a:true'")
	}
	n := parser.NewJsonBool("", false)
	if n.Equal(parser.NewJsonBool("a", false)) {
		t.Error("failed: Bool ':false' should not equal Bool 'a:false'")
	}
	if n.Equal(parser.NewJsonBool("", true)) {
		t.Error("failed: Bool ':false' should not equal Bool ':true'")
	}
	if !n.Equal(parser.NewJsonBool("", false)) {
		t.Error("failed: Bool ':false' should equal Bool ':false'")
	}

}

func TestEqualNull(t *testing.T) {
	na := parser.NewJsonNull("a")
	if na.Equal(parser.NewJsonNull("b")) {
		t.Error("failed: Null 'a' should not equal Null 'b'")
	}
	if na.Equal(parser.NewJsonNumber("b", 10)) {
		t.Error("failed: Null 'a' should not equal Number 'b'")
	}
	if !na.Equal(parser.NewJsonNull("a")) {
		t.Error("failed: Null 'a' should equal 'a'")
	}
	n := parser.NewJsonNull("")
	if n.Equal(parser.NewJsonNull("b")) {
		t.Error("failed: Null '' should not equal Null 'b'")
	}
	if n.Equal(parser.NewJsonNumber("b", 10)) {
		t.Error("failed: Null '' should not equal Number 'b'")
	}
	if !n.Equal(parser.NewJsonNull("")) {
		t.Error("failed: Null '' should equal ''")
	}
}

func TestEqualNumber(t *testing.T) {
	n1 := parser.NewJsonNumber("a", 123)
	if n1.Equal(parser.NewJsonNumber("b", 124)) {
		t.Error("failed: Number 'a:123' should not equal 'b:123'")
	}
	if n1.Equal(parser.NewJsonNumber("a", 124)) {
		t.Error("failed: Number 'a:123' should not equal 'a:124'")
	}
	if n1.Equal(parser.NewJsonNull("a")) {
		t.Error("failed: Number 'a:123' should not equal Null 'a'")
	}
	if !n1.Equal(parser.NewJsonNumber("a", 123)) {
		t.Error("failed: Number 'a:123' should equal 'a:123'")
	}
	n := parser.NewJsonNumber("", 123)
	if n.Equal(parser.NewJsonNumber("", 124)) {
		t.Error("failed: Number ':123' should not equal ':124'")
	}
	if !n.Equal(parser.NewJsonNumber("", 123)) {
		t.Error("failed: Number ':123' should equal ':123'")
	}
}

func TestObjectClear(t *testing.T) {
	root1, err := parser.Parse(obj5)
	if err != nil {
		t.Errorf("failed: to parse obj5: %s", err.Error())
		return
	}
	a, err := parser.Find(root1, parser.NewDotPath("address"))
	if err != nil {
		t.Errorf("failed: to find address: %s", err.Error())
		return
	}
	if a.GetNodeType() != parser.NT_OBJECT {
		t.Error("address should be an object")
		return
	}
	aO := a.(*parser.JsonObject)
	if aO.Len() != 3 {
		t.Errorf("address should be 3 not %d", aO.Len())
		return
	}
	CheckContains(t, aO, "\"list2\": [{\"string\": \"home\"}")
	aO.Clear()
	if aO.Len() != 0 {
		t.Errorf("address should be 0 not %d", aO.Len())
		return
	}
	CheckContains(t, aO, "{}")
}

func TestListClear(t *testing.T) {
	root1, err := parser.Parse(obj5)
	if err != nil {
		t.Errorf("failed: to parse obj5: %s", err.Error())
		return
	}
	l, err := parser.Find(root1, parser.NewDotPath("list1"))
	if err != nil {
		t.Errorf("failed: to find list1: %s", err.Error())
		return
	}
	if l.GetNodeType() != parser.NT_LIST {
		t.Error("List1 should be a list")
		return
	}
	lO := l.(*parser.JsonList)
	if lO.Len() != 4 {
		t.Errorf("List1 should be 4 not %d", lO.Len())
		return
	}
	CheckContains(t, lO, "\"Joe\",{\"lastName\": \"Jackson\"}")
	lO.Clear()
	if lO.Len() != 0 {
		t.Errorf("List1 should be 0 not %d", lO.Len())
		return
	}
	CheckContains(t, lO, "[]")

	l2, err := parser.Find(root1, parser.NewDotPath("list2"))
	if err != nil {
		t.Errorf("failed: to find list2: %s", err.Error())
		return
	}
	if l2.GetNodeType() != parser.NT_LIST {
		t.Error("List2 should be a list")
		return
	}
	l2O := l2.(*parser.JsonList)
	if l2O.Len() != 3 {
		t.Errorf("List2 should be 3 not %d", l2O.Len())
		return
	}
	CheckContains(t, l2, "{\"lastName\": \"Jackson\"}")
	l2O.Clear()
	if l2O.Len() != 0 {
		t.Errorf("List2 should be 0 not %d", l2O.Len())
		return
	}
	CheckContains(t, l2, "[]")
}

func TestClone(t *testing.T) {
	root1, err := parser.Parse(obj2)
	if err != nil {
		t.Errorf("failed: to parse obj2: %s", err.Error())
		return
	}
	root2 := parser.Clone(root1, root1.GetName(), true)
	CompareTrees(root1, root2)

	root1, err = parser.Parse(nList2)
	if err != nil {
		t.Errorf("failed: to parse obj2: %s", err.Error())
		return
	}
	root2 = parser.Clone(root1, root1.GetName(), true)
	CompareTrees(root1, root2)

	root1, err = parser.Parse(obj5)
	if err != nil {
		t.Errorf("failed: to parse obj2: %s", err.Error())
		return
	}
	root2 = parser.Clone(root1, root1.GetName(), true)
	CompareTrees(root1, root2)
}

func TestRemoveFromObject(t *testing.T) {
	root, err := parser.Parse(obj2)
	if err != nil {
		t.Errorf("failed: to parse obj2: %s", err.Error())
		return
	}
	n, err := parser.Find(root, parser.NewDotPath("gender"))
	if err != nil {
		t.Errorf("failed: to find gender: %s", err.Error())
		return
	}
	rootL := root.(*parser.JsonObject)
	rootLen := rootL.Len()
	rootL.Remove(n)
	_, err = parser.Find(root, parser.NewDotPath("gender"))
	if err == nil {
		t.Errorf("failed: to remove gender: %s", err.Error())
		return
	}
	if rootLen != (rootL.Len() + 1) {
		t.Errorf("failed: Len should be reduced by 1: %d %d (-1)", rootLen, rootL.Len())
		return
	}

	address, err := parser.Find(root, parser.NewDotPath("address"))
	if err != nil {
		t.Errorf("failed: to find address: %s", err.Error())
		return
	}
	city, err := parser.Find(root, parser.NewDotPath("address.city"))
	if err != nil {
		t.Errorf("failed: to find address: %s", err.Error())
		return
	}
	address.(*parser.JsonObject).Remove(city)
	_, err = parser.Find(root, parser.NewDotPath("address.city"))
	if err == nil {
		t.Errorf("failed: to remove address.city: %s", err.Error())
		return
	}
}

func TestRemoveFromList(t *testing.T) {
	root, err := parser.Parse(nList1)
	if err != nil {
		t.Errorf("failed: to parse nList1: %s", err.Error())
		return
	}
	n, err := parser.Find(root, parser.NewDotPath("obj"))
	if err != nil {
		t.Errorf("failed: to find obj: %s", err.Error())
		return
	}
	rootL := root.(*parser.JsonList)
	rootLen := rootL.Len()
	rootL.Remove(n)
	_, err = parser.Find(root, parser.NewDotPath("obj"))
	if err == nil {
		t.Errorf("failed: to remove obj: %s", err.Error())
		return
	}
	if rootLen != (rootL.Len() + 1) {
		t.Errorf("failed: Len should be reduced by 1: %d %d (-1)", rootLen, rootL.Len())
		return
	}

	objZero := rootL.GetValues()[0]
	rootL.Remove(objZero)
	if rootL.GetValues()[0] == objZero {
		t.Error("failed: to remove literal at [0]")
		return
	}
	if rootLen != (rootL.Len() + 2) {
		t.Errorf("failed: Len should be reduced by 2: %d %d (-2)", rootLen, rootL.Len())
		return
	}

}

func TestObjectCreateAtPath(t *testing.T) {
	root, err := parser.Parse(obj1)
	if err != nil {
		t.Errorf("failed: to parse obj3: %s", err.Error())
		return
	}
	oc, err := parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("title"), parser.NT_STRING)
	if err != nil {
		t.Error("Should return err")
	}
	oc.(*parser.JsonString).SetValue("mr")
	if !strings.Contains(root.JsonValue(), "\"title\": \"mr\"") {
		t.Error("Should contain '\"title\": \"mr\"'")
	}

	_, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("no"), parser.NT_NUMBER)
	if err != nil {
		t.Error("Should not return err")
	}
	// Try it twice to ensure it is ok if node already exists
	oc, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("no"), parser.NT_NUMBER)
	if err != nil {
		t.Error("Should not return err")
	}
	oc.(*parser.JsonNumber).SetValue(128)
	if !strings.Contains(root.JsonValue(), "\"no\": 128") {
		t.Errorf("Json %s Should contain '\"no\": 128'", root.JsonValue())
	}

	oc, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("alive"), parser.NT_BOOL)
	if err != nil {
		t.Error("Should return err")
	}
	oc.(*parser.JsonBool).SetValue(true)
	if !strings.Contains(root.JsonValue(), "\"alive\": true") {
		t.Errorf("Json %s Should contain '\"alive\": true'", root.JsonValue())
	}

	oc, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("list"), parser.NT_LIST)
	if err != nil {
		t.Error("Should return err")
	}
	oc.(*parser.JsonList).Add(parser.NewJsonNull("listnul"))
	if !strings.Contains(root.JsonValue(), "\"list\": [{\"listnul\": null}]") {
		t.Errorf("Json %s Should contain '\"list\": [{\"listnul\": null}]'", root.JsonValue())
	}

	oc, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("obj"), parser.NT_OBJECT)
	if err != nil {
		t.Error("Should return err")
	}
	oc.(*parser.JsonObject).Add(parser.NewJsonNull("objnul"))
	if !strings.Contains(root.JsonValue(), "\"obj\": {\"objnul\": null}") {
		t.Errorf("Json %s Should contain '\"obj\": {\"objnul\": null}'", root.JsonValue())
	}

	oc, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("obj.name"), parser.NT_STRING)
	if err != nil {
		t.Error("Should return err")
	}
	oc.(*parser.JsonString).SetValue("fred")
	if !strings.Contains(root.JsonValue(), "\"name\": \"fred\"") {
		t.Errorf("Json %s Should contain '\"name\": \"fred\"'", root.JsonValue())
	}

	oc, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("list.name"), parser.NT_STRING)
	if err != nil {
		t.Error("Should return err")
	}
	oc.(*parser.JsonString).SetValue("listfred")
	if !strings.Contains(root.JsonValue(), "\"name\": \"listfred\"") {
		t.Errorf("Json %s Should contain '\"name\": \"listfred\"'", root.JsonValue())
	}

	oc, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("obj.a.b.name"), parser.NT_STRING)
	if err != nil {
		t.Error("Should return err")
	}
	oc.(*parser.JsonString).SetValue("abname")

	if !strings.Contains(root.JsonValue(), "\"a\": {\"b\": {\"name\": \"abname\"}}") {
		t.Errorf("Json %s Should contain '\"a\": {\"b\": {\"name\": \"abname\"}}'", root.JsonValue())
	}

	//	t.Errorf(root.JsonValue())

}
func TestObjectCreateAtPathErrors(t *testing.T) {
	root, err := parser.Parse(obj1)
	if err != nil {
		t.Errorf("failed: to parse obj3: %s", err.Error())
		return
	}
	_, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath(""), parser.NT_NULL)
	if err == nil {
		t.Error("Should return err")
	}
	_, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("."), parser.NT_NULL)
	if err == nil {
		t.Error("Should return err")
	}
	_, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("gender.x"), parser.NT_NULL)
	if err == nil {
		t.Error("Should return err")
	}

	g, err := parser.Find(root, parser.NewDotPath("firstName"))
	if err != nil || g == nil {
		t.Error("could not find 'gender' node")
	}
	_, err = parser.CreateAndReturnNodeAtPath(g, parser.NewDotPath("x"), parser.NT_NULL)
	if err == nil {
		t.Error("Should return err")
	}
	_, err = parser.CreateAndReturnNodeAtPath(root, parser.NewDotPath("abc"), parser.NT_NULL)
	if err != nil {
		t.Error(err.Error())
	}
	if !strings.Contains(root.JsonValue(), "\"abc\": null") {
		t.Error("Should contain '\"abc\": null'")
	}
	g, err = parser.Find(root, parser.NewDotPath("abc"))
	if err != nil || g == nil {
		t.Error("could not find 'abc' node")
	}
	// t.Errorf(root.JsonValue())
}

func TestObjectCreateType(t *testing.T) {
	o := testCreateType(t, "OBJ", parser.NT_OBJECT)
	oO := o.(*parser.JsonObject)
	if oO.JsonValue() != "\"OBJ\": {}" {
		t.Errorf("failed: to produce correct JSON for the node : %s", oO.JsonValue())
	}
	l := testCreateType(t, "LIST", parser.NT_LIST)
	lO := l.(*parser.JsonList)
	if lO.JsonValue() != "\"LIST\": []" {
		t.Errorf("failed: to produce correct JSON for the node : %s", lO.JsonValue())
	}
	n := testCreateType(t, "NUM", parser.NT_NUMBER)
	nO := n.(*parser.JsonNumber)
	if nO.JsonValue() != "\"NUM\": 0" {
		t.Errorf("failed: to produce correct JSON for the node : %s", nO.JsonValue())
	}
	s := testCreateType(t, "STR", parser.NT_STRING)
	sO := s.(*parser.JsonString)
	if sO.JsonValue() != "\"STR\": \"\"" {
		t.Errorf("failed: to produce correct JSON for the node : %s", sO.JsonValue())
	}
	b := testCreateType(t, "BOOL", parser.NT_BOOL)
	bO := b.(*parser.JsonBool)
	if bO.JsonValue() != "\"BOOL\": false" {
		t.Errorf("failed: to produce correct JSON for the node : %s", bO.JsonValue())
	}
	x := testCreateType(t, "NULL", parser.NT_NULL)
	xO := x.(*parser.JsonNull)
	if xO.JsonValue() != "\"NULL\": null" {
		t.Errorf("failed: to produce correct JSON for the node : %s", xO.JsonValue())
	}
}

func testCreateType(t *testing.T, name string, nt parser.NodeType) parser.NodeI {
	o := parser.NewJsonType(name, nt)
	if o.GetNodeType() != nt {
		t.Errorf("failed: to create node with type : %s", parser.GetNodeTypeName(nt))
		return nil
	}
	if o.GetName() != name {
		t.Errorf("failed: to create node with type : %s and name %s", parser.GetNodeTypeName(nt), name)
		return nil
	}
	return o
}

func TestObjectRenameToDuplicate(t *testing.T) {
	root, err := parser.Parse(obj3)
	if err != nil {
		t.Errorf("failed: to parse obj3: %s", err.Error())
		return
	}
	errm := "Parent already has a node with the new name"
	testRename(t, root, "age", "address", "age", "address", errm, false)
	testRename(t, root, "gender", "address", "gender", "address", errm, false)
	testRename(t, root, "bo", "address", "bo", "address", errm, false)
	testRename(t, root, "no", "address", "no", "address", errm, false)
}

func TestObjectRename(t *testing.T) {
	root, err := parser.Parse(obj3)
	if err != nil {
		t.Errorf("failed: to parse obj3: %s", err.Error())
		return
	}

	// Test object in object
	testRename(t, root, "address.streetAddress", "address.road", "streetAddress", "road", "", false)
	testRename(t, root, "address.road", "address.streetAddress", "road", "streetAddress", "", false)
	// Test String Object in root
	testRename(t, root, "address", "address1", "address", "address1", "", false)
	testRename(t, root, "address1", "address", "address1", "address", "", false)
	// Test Number Object in root
	testRename(t, root, "age", "old", "age", "old", "", false)
	testRename(t, root, "old", "age", "old", "age", "", false)
	// Test Bool Object in root
	testRename(t, root, "bo", "boooo", "bo", "boooo", "", false)
	testRename(t, root, "boooo", "bo", "boooo", "bo", "", false)
	// Test Null Object in root
	testRename(t, root, "no", "noooo", "no", "noooo", "", false)
	testRename(t, root, "noooo", "no", "noooo", "no", "", false)
	// Test Number object in list
	testRename(t, root, "address.phoneNumbers.number", "address.phoneNumbers.num", "number", "num", "", false)
	testRename(t, root, "address.phoneNumbers.num", "address.phoneNumbers.number", "num", "number", "", false)
	// Test string object in list
	testRename(t, root, "address.phoneNumbers.string", "address.phoneNumbers.str", "string", "str", "", false)
	testRename(t, root, "address.phoneNumbers.str", "address.phoneNumbers.string", "str", "string", "", false)
	// Test bool object in list
	testRename(t, root, "address.phoneNumbers.boolean", "address.phoneNumbers.bool", "boolean", "bool", "", false)
	testRename(t, root, "address.phoneNumbers.bool", "address.phoneNumbers.boolean", "bool", "boolean", "", false)

	testRename(t, root, "address.phoneNumbers.number", "address.phoneNumbers.boolean", "number", "boolean", "", true)
	testRename(t, root, "address.phoneNumbers.no", "address.phoneNumbers.dupe1", "no", "dupe1", "", true)

}

func testRename(t *testing.T, root parser.NodeI, pathBefore, pathAfter, nameBefore, nameAfter string, failContains string, reqDuplicateAfter bool) {
	nb, err := parser.Find(root, parser.NewDotPath(pathBefore))
	if err != nil {
		t.Errorf("failed: Not Found %s [%s]", pathBefore, nameBefore)
		return
	}
	if nb.GetName() != nameBefore {
		t.Errorf("failed: Found but name is wrong %s [%s] --> actual name:%s ", pathBefore, nameBefore, nb.GetName())
		return
	}
	nbp, ok := parser.FindParentNode(root, nb)
	if !ok {
		t.Errorf("failed: Before rename. Failed to find parent of %s [%s]", pathBefore, nameBefore)
		return
	}

	if nbp.GetNodeType() == parser.NT_OBJECT {
		nb1 := nbp.(*parser.JsonObject).GetNodeWithName(nameBefore)
		if nb1 == nil {
			t.Errorf("failed: Before rename node not found in parent map %s [%s]", pathBefore, nameBefore)
			return
		}
		if nb1.GetName() != nameBefore {
			t.Errorf("failed: Before rename node found in parent map has wrong name %s [%s] --> actual name:%s", pathBefore, nameBefore, nb1.GetName())
			return
		}
	} else {
		nb2 := nbp.(*parser.JsonList).GetNodeWithName(nameBefore)
		if nb2 == nil {
			t.Errorf("failed: Before rename node not found in parent list %s [%s]", pathBefore, nameBefore)
			return
		}
		if nb2.GetName() != nameBefore {
			t.Errorf("failed: Before rename node found in parent list has wrong name %s [%s] --> actual name:%s", pathBefore, nameBefore, nb2.GetName())
			return
		}
	}

	err = parser.Rename(root, nb, nameAfter)
	if failContains != "" {
		if err == nil {
			t.Errorf("failed: Rename did not return an error")
			return
		}
		if !strings.Contains(err.Error(), failContains) {
			t.Errorf("failed: Rename error should contain '%s', actual: '%s'", failContains, err.Error())
		}
		return
	} else {
		if err != nil {
			t.Errorf("failed: Rename %s [%s] --> %s", pathBefore, nameBefore, nameAfter)
			return
		}
	}

	na, err := parser.Find(root, parser.NewDotPath(pathAfter))
	if err != nil {
		t.Errorf("failed: Not Found after rename %s [%s] --> %s [%s] ", pathBefore, nameBefore, pathAfter, nameAfter)
		return
	}
	if na.GetName() != nameAfter {
		t.Errorf("failed: Renamed but new name not set in object %s [%s] --> actual name:%s != %s", pathBefore, nameBefore, na.GetName(), nameAfter)
		return
	}

	nap, ok := parser.FindParentNode(root, na)
	if !ok {
		t.Errorf("failed: After rename. Failed to find parent of %s [%s]", pathAfter, nameAfter)
		return
	}
	if nap.GetNodeType() == parser.NT_OBJECT {
		na1 := nap.(*parser.JsonObject).GetNodeWithName(nameAfter)
		if na1 == nil {
			t.Errorf("failed: After node not found in parent map %s [%s]", pathAfter, nameAfter)
			return
		}
		if na1.GetName() != nameAfter {
			t.Errorf("failed: After node found in parent map has wrong name %s [%s] --> actual name:%s", pathAfter, nameAfter, na1.GetName())
			return
		}
	} else {
		na2 := nap.(*parser.JsonList).GetNodeWithName(nameAfter)
		if na2 == nil {
			t.Errorf("failed: After rename node not found in parent list %s [%s]", pathAfter, nameAfter)
			return
		}
		if na2.GetName() != nameAfter {
			t.Errorf("failed: Before rename node found in parent list has wrong name %s [%s] --> actual name:%s", pathAfter, nameAfter, na2.GetName())
			return
		}
	}

	if nbp != nap {
		t.Errorf("failed: After rename parent nodes are not the same %s [%s] --> %s [%s]", pathBefore, nameBefore, pathAfter, nameAfter)
		return
	}

	if reqDuplicateAfter {
		_, err = parser.Find(root, parser.NewDotPath(pathAfter))
		if err != nil {
			t.Errorf("failed: After rename, duplicate was not found %s [%s]", pathAfter, nameAfter)
			return
		}
	}
	_, err = parser.Find(root, parser.NewDotPath(pathBefore))
	if err == nil {
		t.Errorf("failed: After rename, before still found %s [%s]", pathBefore, nameBefore)
		return
	}
}

func TestObjectDuplicateName(t *testing.T) {
	root := parser.NewJsonObject("")
	err := root.Add(parser.NewJsonBool("B_OBJ", true))
	if err != nil {
		t.Errorf("Should not return an error")
	}
	err = root.Add(parser.NewJsonNumber("B_OBJ", 12345))
	if err == nil {
		t.Errorf("Should return an error")
	}
	if !strings.Contains(err.Error(), "duplicate name") {
		t.Errorf("Should return an error that contains 'duplicate name'")
	}
}

func TestListDuplicateName(t *testing.T) {
	root := parser.NewJsonList("")
	root.Add(parser.NewJsonBool("B_OBJ", true))
	root.Add(parser.NewJsonString("", "string"))
	root.Add(parser.NewJsonNumber("B_OBJ", 12345))
}

func TestListDuplicateNameInWrapper(t *testing.T) {
	root := parser.NewJsonList("")
	root.Add(parser.NewJsonBool("B_OBJ", true))
	root.Add(parser.NewJsonString("", "string"))
	w := parser.NewJsonObject("")
	w.Add(parser.NewJsonNumber("B_OBJ", 12345))
	root.Add(w)
}

func TestObjectValuesSorted(t *testing.T) {
	root := parser.NewJsonObject("")
	root.Add(parser.NewJsonBool("B_OBJ", true))
	subObj1 := parser.NewJsonObject("A_OBJ")
	subObj1.Add(parser.NewJsonString("3_OBJ", "ABC"))
	subObj1.Add(parser.NewJsonBool("2_OBJ", true))
	subObj1.Add(parser.NewJsonNumber("1_OBJ", 123))
	root.Add(subObj1)
	values := root.GetValuesSorted()
	if len(values) != 2 {
		t.Errorf("Should be two values. Found %d", len(values))
	}
	if values[0].GetName() != "A_OBJ" {
		t.Errorf("Value 0 should not be = %s", values[0].GetName())
	}
	if values[1].GetName() != "B_OBJ" {
		t.Errorf("Value 1 should not be = %s", values[1].GetName())
	}
	values2 := values[0].(*parser.JsonObject).GetValuesSorted()
	if len(values2) != 3 {
		t.Errorf("Should be two values. Found %d", len(values2))
	}
	if values2[0].GetName() != "1_OBJ" {
		t.Errorf("Value2 0 should not be = %s", values2[0].GetName())
	}
	if values2[1].GetName() != "2_OBJ" {
		t.Errorf("Value2 1 should not be = %s", values2[1].GetName())
	}
	if values2[2].GetName() != "3_OBJ" {
		t.Errorf("Value2 2 should not be = %s", values2[2].GetName())
	}
}
func TestListValues(t *testing.T) {
	root := parser.NewJsonList("")
	root.Add(parser.NewJsonBool("B_OBJ", true))
	subObj1 := parser.NewJsonList("A_OBJ")
	subObj1.Add(parser.NewJsonString("3_OBJ", "ABC"))
	subObj1.Add(parser.NewJsonBool("2_OBJ", true))
	subObj1.Add(parser.NewJsonNumber("1_OBJ", 123))
	root.Add(subObj1)
	values := root.GetValues()
	if len(values) != 2 {
		t.Errorf("Should be two values. Found %d", len(values))
	}
	if values[0].GetName() != "B_OBJ" {
		t.Errorf("Value 0 should not be = %s", values[0].GetName())
	}
	if values[1].GetName() != "A_OBJ" {
		t.Errorf("Value 1 should not be = %s", values[1].GetName())
	}
	values2 := values[1].(*parser.JsonList).GetValues()
	if len(values2) != 3 {
		t.Errorf("Should be three values. Found %d", len(values))
	}
	if values2[0].GetName() != "3_OBJ" {
		t.Errorf("Value2 0 should not be = %s", values2[3].GetName())
	}
	if values2[1].GetName() != "2_OBJ" {
		t.Errorf("Value2 1 should not be = %s", values2[1].GetName())
	}
	if values2[2].GetName() != "1_OBJ" {
		t.Errorf("Value2 2 should not be = %s", values2[2].GetName())
	}
}
func TestObjectValues(t *testing.T) {
	root := parser.NewJsonObject("")
	root.Add(parser.NewJsonBool("B_OBJ", true))
	subObj1 := parser.NewJsonObject("A_OBJ")
	subObj1.Add(parser.NewJsonString("3_OBJ", "ABC"))
	subObj1.Add(parser.NewJsonBool("2_OBJ", true))
	subObj1.Add(parser.NewJsonNumber("1_OBJ", 123))
	root.Add(subObj1)
	values := root.GetValuesSorted()
	if len(values) != 2 {
		t.Errorf("Should be two values. Found %d", len(values))
	}
	if values[0].GetName() != "A_OBJ" {
		t.Errorf("Value 0 should not be = %s", values[0].GetName())
	}
	if values[1].GetName() != "B_OBJ" {
		t.Errorf("Value 1 should not be = %s", values[1].GetName())
	}
	values2 := values[0].(*parser.JsonObject).GetValuesSorted()
	if len(values2) != 3 {
		t.Errorf("Should be three values. Found %d", len(values))
	}
	if values2[0].GetName() != "1_OBJ" {
		t.Errorf("Value2 0 should not be = %s", values2[0].GetName())
	}
	if values2[1].GetName() != "2_OBJ" {
		t.Errorf("Value2 1 should not be = %s", values2[1].GetName())
	}
	if values2[2].GetName() != "3_OBJ" {
		t.Errorf("Value2 2 should not be = %s", values2[2].GetName())
	}
}

func TestRootListWithMapsAndLists(t *testing.T) {
	root := parser.NewJsonList("")
	root.Add(parser.NewJsonBool("OBJ.BO.1", true))

	list1 := parser.NewJsonList("EMPTY")
	root.Add(list1)
	list2 := parser.NewJsonList("FRED")
	list2.Add(parser.NewJsonString("", "Help"))
	root.Add(list2)

	subObj1 := parser.NewJsonObject("OBJ.1")
	subObj1.Add(parser.NewJsonString("OBJ.1.ST.1", "ABC"))
	subObj1.Add(parser.NewJsonBool("OBJ.1.BO.1", true))
	subObj1.Add(parser.NewJsonNumber("OBJ.1.FL.1", 123))

	root.Add(subObj1)

	parent, ok := parser.FindParentNode(root, subObj1)
	if !ok {
		t.Errorf("FindNode error: Did not fine subObj1")
	}
	if parent.GetName() != root.GetName() {
		t.Errorf("FindNode Value error: parent '%s' root %s", parent.GetName(), root.GetName())
	}
	parser.Rename(root, subObj1, "BLUE")
	l := strings.Split(string(testRootListWithMapsAndLists), "|")
	s := root.JsonValue()
	for _, ts := range l {
		if !strings.Contains(s, ts) {
			t.Errorf("Node Value error: \nactual %s \nshould contain %s", s, ts)
		}
	}
}

func TestRootMapWithLists(t *testing.T) {
	root := parser.NewJsonObject("")
	root.Add(parser.NewJsonString("OBJ.ST.1", "ABC"))
	root.Add(parser.NewJsonBool("OBJ.BO.1", true))
	root.Add(parser.NewJsonNumber("OBJ.FL.1", 123))

	subList1 := parser.NewJsonList("LI.1")
	subList1.Add(parser.NewJsonString("LI.ST.1", "ABC"))
	subList1.Add(parser.NewJsonBool("LI.BO.1", false))

	subList2 := parser.NewJsonList("LI.2")
	subList2.Add(parser.NewJsonString("", "ABC"))
	subList2.Add(parser.NewJsonNumber("", 198))
	subList2.Add(parser.NewJsonBool("", false))
	subList2.Add(parser.NewJsonNumber("NUM", 123))
	root.Add(subList1)
	root.Add(subList2)

	l := strings.Split(string(testRootMapWithLists), "|")
	s := root.JsonValue()
	for _, ts := range l {
		if !strings.Contains(s, ts) {
			t.Errorf("Node Value error: \nactual %s \nshould contain %s", s, ts)
		}
	}

}
func TestNodes(t *testing.T) {
	pl := parser.NewJsonList("")
	pl.Add(parser.NewJsonString("LI.ST.1", "ABC"))
	pl.Add(parser.NewJsonBool("LI.BO.1", false))
	pl.Add(parser.NewJsonNumber("LI.FL.1", 123))

	ob := parser.NewJsonObject("OBJ")
	ob.Add(parser.NewJsonString("OBJ.ST.1", "ABC"))
	ob.Add(parser.NewJsonBool("OBJ.BO.1", true))
	ob.Add(parser.NewJsonNumber("OBJ.FL.1", 123))
	err := ob.Add(parser.NewJsonNumber("OBJ.BO.1", 123))
	if err == nil {
		t.Errorf("Cannot add duplicate nodes to objects")
	}
	err = ob.Add(parser.NewJsonNumber("", 123))
	if err == nil {
		t.Errorf("Cannot add nodes to objects without names")
	}

	pll := parser.NewJsonList("")
	pll.Add(parser.NewJsonBool("", false))
	pll.Add(parser.NewJsonNumber("", 999.1))
	pll.Add(parser.NewJsonString("", "LIT.ST.2"))
	pll.Add(parser.NewJsonString("A", "B"))

	nilNode := parser.NewJsonNull("NN")
	n := []parser.NodeI{
		parser.NewJsonString("ST", "ABC"),
		parser.NewJsonBool("BO", true),
		parser.NewJsonNumber("FL", 123.55),
		pl,
		ob,
		pll,
		nilNode,
	}

	assertType(t, n[0], parser.NT_STRING, []string{"ABC"}, 0)
	assertType(t, n[1], parser.NT_BOOL, []string{"true"}, 0)
	assertType(t, n[2], parser.NT_NUMBER, []string{"123.55"}, 123)
	assertType(t, n[3], parser.NT_LIST, []string{"[{\"LI.ST.1\": \"ABC\"},{\"LI.BO.1\": false},{\"LI.FL.1\": 123}]"}, 3)
	assertType(t, n[4], parser.NT_OBJECT, []string{"\"OBJ\": {", "\"OBJ.ST.1\": \"ABC\"", "\"OBJ.BO.1\": true", "\"OBJ.FL.1\": 123"}, 3)
	assertType(t, n[5], parser.NT_LIST, []string{"[false,999.1,\"LIT.ST.2\",{\"A\": \"B\"}]"}, 4)
	assertType(t, n[6], parser.NT_NULL, []string{"\"NN\": null"}, 0)

	assertType(t, ob.GetNodeWithName("OBJ.BO.1"), parser.NT_BOOL, []string{"true"}, 0)
	b := pl.GetNodeWithName("LI.BO.1")
	assertType(t, b, parser.NT_BOOL, []string{"false"}, 0)
	fmt.Println("END")
}

func assertType(t *testing.T, node parser.NodeI, nt parser.NodeType, strValue []string, eleCount int) {
	if node.GetNodeType() != nt {
		t.Errorf("Node Type error: expected %s actual %s", parser.GetNodeTypeName(nt), parser.GetNodeTypeName(node.GetNodeType()))
	}
	switch node.GetNodeType() {
	case parser.NT_OBJECT:
		n := (node.(*parser.JsonObject))
		if n.GetName() != node.GetName() {
			t.Errorf("Object Node String error: expected '%s' actual '%s'", n.GetName(), node.String())
		}
		for _, ts := range strValue {
			if !strings.Contains(n.JsonValue(), ts) {
				t.Errorf("Node JsonValue error: \nactual %s \nshould contain %s", n, ts)
			}
		}
		if n.Len() != eleCount {
			t.Errorf("List Count error: expected %d actual %d", eleCount, n.Len())
		}
	case parser.NT_LIST:
		n := (node.(*parser.JsonList))
		if n.GetName() != node.GetName() {
			t.Errorf("String Node String error: expected '%s' actual '%s'", n.GetName(), node.GetName())
		}
		if n.JsonValue() != strValue[0] {
			t.Errorf("Node JsonValue error: \nexpected '%s' \n  actual '%s'", strValue[0], n.JsonValue())
		}
		if n.Len() != eleCount {
			t.Errorf("Node Count error: expected %d actual %d", eleCount, n.Len())
		}
	case parser.NT_NUMBER:
		n := (node.(*parser.JsonNumber))
		sf := fmt.Sprintf("%f", n.GetValue())
		if !strings.HasPrefix(sf, node.String()) {
			t.Errorf("Number Node String error: expected %s actual %s", sf, node.String())
		}
		if !strings.HasPrefix(sf, strValue[0]) {
			t.Errorf("Number Node Value error: expected %s actual %s", strValue[0], sf)
		}
		if n.GetIntValue() != int64(eleCount) {
			t.Errorf("Number Node Value error: expected %d actual %d", n.GetIntValue(), eleCount)
		}
	case parser.NT_STRING:
		v := (node.(*parser.JsonString)).GetValue()
		if v != node.String() {
			t.Errorf("String Node String error: expected %s actual %s", v, node.String())
		}
		if v != strValue[0] {
			t.Errorf("String Node Value error: expected %s actual %s", strValue[0], v)
		}
	case parser.NT_BOOL:
		v := fmt.Sprintf("%t", (node.(*parser.JsonBool)).GetValue())
		if v != node.String() {
			t.Errorf("Bool Node String error: expected %s actual %s", v, node.String())
		}
		if v != strValue[0] {
			t.Errorf("Bool Node Value error: expected %s actual %s", strValue[0], v)
		}
	case parser.NT_NULL:
		if node.String() != "null" {
			t.Errorf("Null Node String error: expected %s actual %s", "null", node.String())
		}
		if node.JsonValue() != strValue[0] {
			t.Errorf("Null Node.JsonValue() Value error: expected %s actual %s", strValue[0], node.JsonValue())
		}
	}

}
