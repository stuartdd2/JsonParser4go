package test

import (
	"reflect"
	"testing"

	"github.com/stuartdd2/JsonParser4go/parser"
)

func TestTrailString(t *testing.T) {
	root := InitParser(t, "", obj4)
	sa, _ := parser.Find(root, parser.NewDotPath("list1.0"))
	ln, _ := parser.Find(root, parser.NewDotPath("list1.lastName"))
	aa, _ := parser.Find(root, parser.NewDotPath("list1.3.A"))
	trail := parser.NewTrail(4, "|")
	if trail.String() != "" {
		t.Errorf("Empty Trail empty string")
	}
	trail.Push(sa, 0)
	if trail.String() != "" {
		t.Errorf("Empty because sa has no name")
	}
	l := trail.GetList()
	if len(l) != 1 {
		t.Errorf("List should be 1 long")
	}
	if (*l[0]).String() != "Joe" {
		t.Errorf("Emprt node retirned ion list shoule have value Joe")
	}
	trail.Clear()
	trail.Push(ln, 1)
	if trail.String() != "lastName" {
		t.Errorf("String should be 'lastname' == '%s'", trail.String())
	}
	trail.Push(sa, 2)
	if trail.String() != "lastName|2" {
		t.Errorf("lastname + (no name) == '%s'", trail.String())
	}
	trail.Push(aa, 4)
	if trail.String() != "lastName|2|A" {
		t.Errorf("lastname + (no name) + A == '%s'", trail.String())
	}
}

func TestTrailClear(t *testing.T) {
	root := InitParser(t, "", obj3)
	trail := parser.NewTrail(2, "|")
	sa, _ := parser.Find(root, parser.NewDotPath("address.streetAddress"))
	ln, _ := parser.Find(root, parser.NewDotPath("lastName"))
	trail.Push(sa, 0)
	trail.Push(ln, 1)
	l := trail.GetList()
	if len(l) != 2 {
		t.Errorf("List should be 2 long")
	}
	trail.Clear()
	ll := trail.GetList()
	if len(ll) != 0 {
		t.Errorf("List should be 0 long")
	}
}

func TestTrailList(t *testing.T) {
	root := InitParser(t, "", obj3)
	trail := parser.NewTrail(2, "|")
	sa, _ := parser.Find(root, parser.NewDotPath("address.streetAddress"))
	ln, _ := parser.Find(root, parser.NewDotPath("lastName"))
	trail.Push(sa, 0)
	trail.Push(ln, 1)
	l := trail.GetList()
	if len(l) != 2 {
		t.Errorf("List should be 2 long")
	}
	l0 := *l[0]
	l1 := *l[1]
	if sa.String() != l0.String() {
		t.Errorf("List [0] should be sa")
	}
	if ln.String() != l1.String() {
		t.Errorf("List [1] should be ln")
	}
	l1.(*parser.JsonString).SetValue("HOCK")
	if l1.String() != "HOCK" {
		t.Errorf("ln is same as l1 so value should be same")
	}
	if ln.String() != "HOCK" {
		t.Errorf("ln is same as l1 so value should be same")
	}
	sa.(*parser.JsonString).SetValue("999")
	if sa.String() != "999" {
		t.Errorf("sa is same as l0 so value should be same")
	}
	if l0.String() != "999" {
		t.Errorf("ln is same as l1 so value should be same")
	}
	trail.Pop()
	ll := trail.GetList()
	if len(ll) != 1 {
		t.Errorf("List should be 1 long")
	}
	ll0 := *ll[0]
	if ll0 != l0 {
		t.Errorf("Should be the same object")
	}
	if ll0 != sa {
		t.Errorf("Should be the same object")
	}
	trail.Pop()
	lll := trail.GetList()
	if len(lll) != 0 {
		t.Errorf("List should be empty")
	}
	if len(ll) != 1 {
		t.Errorf("List should be 1 long")
	}
	if len(l) != 2 {
		t.Errorf("List should be 2 long")
	}
}

func TestTrailPop(t *testing.T) {
	root := InitParser(t, "", obj3)
	trail := parser.NewTrail(2, "|")
	sa, _ := parser.Find(root, parser.NewDotPath("address.streetAddress"))
	ln, _ := parser.Find(root, parser.NewDotPath("lastName"))
	trail.Push(ln, 0)
	trail.Push(sa, 1)

	if ln.String() != "Jackson" {
		t.Errorf("Eat my hat!")
	}
	if sa.String() != "101" {
		t.Errorf("Eat my hat!")
	}
	sa.(*parser.JsonString).SetValue("999")
	sax, _ := trail.Pop()
	if sa.String() != "999" {
		t.Errorf("Eat my hat!")
	}
	if sax.String() != "999" {
		t.Errorf("Not the same object!!")
	}
	sax.(*parser.JsonString).SetValue("888")
	if sax.String() != "888" {
		t.Errorf("Eat my hat!")
	}
	if sa.String() != "888" {
		t.Errorf("Not the same object!!")
	}
	lnx, _ := trail.Pop()
	if lnx == nil {
		t.Errorf("Adderss should not be nil")
	}
	if lnx.String() != "Jackson" {
		t.Errorf("Eat my hat!")
	}
}

func TestTrailPush(t *testing.T) {
	root := InitParser(t, "", obj3)
	trail := parser.NewTrail(2, "|")
	n1, _ := parser.Find(root, parser.NewDotPath("firstName"))
	n2, _ := parser.Find(root, parser.NewDotPath("address"))
	n3, _ := parser.Find(root, parser.NewDotPath("address.streetAddress"))
	pop, _ := trail.Pop()
	if pop != nil {
		t.Errorf("Trail should pop nil if empty")
	}
	if len(trail.GetList()) != 0 {
		t.Errorf("Trail list should be empty")
	}
	ok := trail.Push(n1, 1)
	if !ok {
		t.Errorf("Push should return true")
	}
	if len(trail.GetList()) != 1 {
		t.Errorf("Trail list should contain 1 item")
	}
	n := trail.GetList()[0]
	if *n != n1 {
		t.Errorf("Trail node should be the same node")
	}

	ok = trail.Push(n2, 4)
	if !ok {
		t.Errorf("Push should return true")
	}
	if len(trail.GetList()) != 2 {
		t.Errorf("Trail list should contain 2 items")
	}
	n = trail.GetList()[1]
	if *n != n2 {
		t.Errorf("Trail node should equal node 2 from tree")
	}

	ok = trail.Push(n3, 5)
	if ok {
		t.Errorf("Trail list should not have any more room")
	}
	if len(trail.GetList()) != 2 {
		t.Errorf("Trail list should contain 2 items")
	}
	n = trail.GetList()[1]
	if *n != n2 {
		t.Errorf("Trail node should equal node 2 from tree")
	}
}

func TestEqual(t *testing.T) {
	p := parser.NewBarPath("a|b|c")
	pp := parser.NewBarPath("a|b|c")
	pp1 := parser.NewBarPath("a|b")
	pp2 := parser.NewBarPath("a|b|2")
	ppD := parser.NewDotPath("a.b.c")

	if p.String() != pp.String() {
		t.Errorf("TestEqual failed String() compare")
	}
	if p.String() == ppD.String() {
		t.Errorf("TestEqual failed String() compare dot and bar")
	}
	if !reflect.DeepEqual(p, pp) {
		t.Errorf("TestEqual failed DeepEqual")
	}
	if reflect.DeepEqual(p, ppD) {
		t.Errorf("TestEqual failed DeepEqual dot and bar")
	}
	if !p.Equal(pp) {
		t.Errorf("TestEqual failed should be Equal")
	}
	if !p.Equal(ppD) {
		t.Errorf("TestEqual failed dot and bar should be Equal")
	}
	if p.Equal(nil) {
		t.Errorf("TestEqual failed Equals(nil) should never be equal")
	}
	if p.Equal(pp1) {
		t.Errorf("TestEqual failed Equals(pp1) not same len")
	}
	if p.Equal(pp2) {
		t.Errorf("TestEqual failed Equals(pp2) not same content")
	}
	if p == pp {
		t.Errorf("TestEqual failed == never return true because of the array in Path")
	}
}
func TestStringFirst(t *testing.T) {
	p1 := parser.NewPath("a|b|c", "|")
	if p1.StringFirst() != "a" {
		t.Errorf("TestStringFirst Path StringFirst '%s' != 'a'", p1.StringFirst())
	}
	p2 := parser.NewPath("b", "|")
	if p2.StringFirst() != "b" {
		t.Errorf("TestStringFirst Path StringFirst '%s' != 'b'", p2.StringFirst())
	}
	p4 := parser.NewPath("", "|")
	if p4.StringFirst() != "" {
		t.Errorf("TestStringFirst Path StringFirst '%s' != ''", p4.StringFirst())
	}
}

func TestStringAppend(t *testing.T) {
	p1 := parser.NewPath("a|b|c", "|")
	testPath(t, 2, p1, "a|b|c", "c", "|", 3)
	p1.StringAppend("1|2")
	testPath(t, 2, p1, "a|b|c|1|2", "2", "|", 5)
}

func TestPathAppend(t *testing.T) {
	p1 := parser.NewPath("a|b|c", "|")
	testPath(t, 2, p1, "a|b|c", "c", "|", 3)
	p2 := parser.NewPath("1|2", "|")
	testPath(t, 2, p2, "1|2", "2", "|", 2)
	p1.PathAppend(p2)
	testPath(t, 2, p1, "a|b|c|1|2", "2", "|", 5)
	testPath(t, 2, p2, "1|2", "2", "|", 2)

	p4 := parser.NewPath("", "|")
	testEmpty(t, 1, p4, "|")
	p2.PathAppend(p4)
	testPath(t, 2, p2, "1|2", "2", "|", 2)
	p4.PathAppend(p1)
	testPath(t, 2, p4, "a|b|c|1|2", "2", "|", 5)
}
func TestPathLast(t *testing.T) {
	p := parser.NewPath("a|b|c", "|")
	testEmpty(t, 1, p.PathLast(0), "|")
	testPath(t, 2, p.PathLast(1), "c", "c", "|", 1)
	testPath(t, 3, p.PathLast(2), "b|c", "c", "|", 2)
	testPath(t, 4, p.PathLast(3), "a|b|c", "c", "|", 3)
	testPath(t, 5, p.PathLast(4), "a|b|c", "c", "|", 3)
	p = parser.NewPath("a|b", "|")
	testEmpty(t, 6, p.PathLast(0), "|")
	testPath(t, 7, p.PathLast(1), "b", "b", "|", 1)
	testPath(t, 8, p.PathLast(2), "a|b", "b", "|", 2)
	testPath(t, 9, p.PathLast(3), "a|b", "b", "|", 2)
	p = parser.NewPath("a", "|")
	testEmpty(t, 8, p.PathLast(0), "|")
	testPath(t, 9, p.PathLast(1), "a", "a", "|", 1)
	testPath(t, 10, p.PathLast(2), "a", "a", "|", 1)
	p = parser.NewPath("", "|")
	testEmpty(t, 11, p.PathLast(0), "|")
	testEmpty(t, 12, p.PathLast(1), "|")
}
func TestPathFirst(t *testing.T) {
	p := parser.NewPath("a|b|c", "|")
	testEmpty(t, 1, p.PathFirst(0), "|")
	testPath(t, 2, p.PathFirst(1), "a", "a", "|", 1)
	testPath(t, 3, p.PathFirst(2), "a|b", "b", "|", 2)
	testPath(t, 4, p.PathFirst(3), "a|b|c", "c", "|", 3)
	p = parser.NewPath("a|b", "|")
	testEmpty(t, 5, p.PathFirst(0), "|")
	testPath(t, 6, p.PathFirst(1), "a", "a", "|", 1)
	testPath(t, 7, p.PathFirst(2), "a|b", "b", "|", 2)
	testPath(t, 8, p.PathFirst(3), "a|b", "b", "|", 2)
	p = parser.NewPath("a", "|")
	testEmpty(t, 9, p.PathFirst(0), "|")
	p = parser.NewPath("", "|")
	testEmpty(t, 10, p.PathFirst(0), "|")
	testEmpty(t, 11, p.PathFirst(1), "|")
}
func TestEmptyPaths(t *testing.T) {
	p := parser.NewPath("", "|")
	testEmpty(t, 1, p, "|")
	p = parser.NewDotPath("")
	testEmpty(t, 2, p, ".")
	p = parser.NewPath("", "")
	testEmpty(t, 3, p, ".")
	p = parser.NewPath("", ".")
	testEmpty(t, 4, p, ".")
}
func TestPaths(t *testing.T) {
	p := parser.NewPath("a", "|")
	testPathAt(t, 1, p, 0, "a")
	testPath(t, 2, p, "a", "a", "|", 1)
	testEmpty(t, 3, p.PathParent(), "|")

	p = parser.NewPath("a.b", "|")
	testPathAt(t, 4, p, 0, "a.b")
	testPath(t, 5, p, "a.b", "a.b", "|", 1)
	testEmpty(t, 6, p.PathParent(), "|")

	p = parser.NewPath("a.b|c", "|")
	testPathAt(t, 7, p, 0, "a.b")
	testPathAt(t, 8, p, 1, "c")
	testPath(t, 9, p, "a.b|c", "c", "|", 2)
	testPath(t, 10, p.PathParent(), "a.b", "a.b", "|", 1)
	testEmpty(t, 11, p.PathParent().PathParent(), "|")

	p = parser.NewDotPath("a.b|c")
	testPathAt(t, 12, p, 0, "a")
	testPathAt(t, 13, p, 1, "b|c")
	testPath(t, 14, p, "a.b|c", "b|c", ".", 2)
	testPath(t, 15, p.PathParent(), "a", "a", ".", 1)
	testEmpty(t, 15, p.PathParent().PathParent(), ".")

	p = parser.NewBarPath("a.b|c|1|2")
	testPathAt(t, 16, p, 0, "a.b")
	testPathAt(t, 17, p, 1, "c")
	testPathAt(t, 18, p, 2, "1")
	testPathAt(t, 19, p, 3, "2")
	testPathAt(t, 20, p, 4, "")
	testPathAt(t, 21, p, -1, "")
	testPath(t, 22, p, "a.b|c|1|2", "2", "|", 4)
	testPath(t, 23, p.PathParent(), "a.b|c|1", "1", "|", 3)
	testPath(t, 24, p.PathParent().PathParent(), "a.b|c", "c", "|", 2)
	testPath(t, 25, p.PathParent().PathParent().PathParent(), "a.b", "a.b", "|", 1)
	testEmpty(t, 26, p.PathParent().PathParent().PathParent().PathParent(), "|")
	p = parser.NewBarPath("")
	testPathAt(t, 27, p, -2, "")
	testPathAt(t, 28, p, -1, "")
	testPathAt(t, 29, p, 0, "")
	testPathAt(t, 30, p, 1, "")
	testPathAt(t, 31, p, 2, "")
}

func testPath(t *testing.T, id int, p *parser.Path, str, last, delim string, len int) {
	if p.GetDelim() != delim {
		t.Errorf("%d Path Delim '%s' != '%s'", id, p.GetDelim(), delim)
		return
	}
	if p.StringLast() != last {
		t.Errorf("%d Path path GetLast() '%s' is not '%s'", id, p.StringLast(), last)
		return
	}
	if p.String() != str {
		t.Errorf("%d Path path String() '%s' is not '%s'", id, p.String(), str)
		return
	}
	if p.Len() != len {
		t.Errorf("%d Path Len() '%d' is not '%d'", id, p.Len(), len)
		return
	}
}

func testEmpty(t *testing.T, id int, p *parser.Path, delim string) {
	if p.GetDelim() != delim {
		t.Errorf("%d Empty Delim %s != %s", id, p.GetDelim(), delim)
		return
	}
	if p.String() != "" {
		t.Errorf("%d Empty path String() '%s' is not \"\"", id, p.String())
		return
	}
	if !p.IsEmpty() {
		t.Errorf("%d Empty path IsEmpty did not return true", id)
	}
	if p.Len() != 0 {
		t.Errorf("%d Empty path Len did not return 0", id)
	}
	if len(p.Paths()) != 0 {
		t.Errorf("%d Empty path Paths len did not return 0", id)
	}
	if !p.PathParent().IsEmpty() {
		t.Errorf("%d Empty path PathParent IsEmpty did not return true", id)
	}
	if p.StringLast() != "" {
		t.Errorf("%d Empty path GetLast did not return \"\"", id)
	}
}

func testPathAt(t *testing.T, id int, p *parser.Path, at int, str string) {
	if p.StringAt(at) != str {
		t.Errorf("%d Path StringAt(%d) '%s' != '%s'", id, at, p.StringAt(at), str)
		return
	}
	if p.PathAt(at).String() != str {
		t.Errorf("%d Path PathAt(%d) '%s' != '%s'", id, at, p.StringAt(at), str)
		return
	}
}
