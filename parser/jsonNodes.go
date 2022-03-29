/*
 * Copyright (C) 2021 Stuart Davies (stuartdd)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package parser

import (
	"fmt"
	"sort"
	"strings"
)

type NodeType int

const (
	NT_OBJECT NodeType = iota
	NT_LIST   NodeType = iota
	NT_NUMBER NodeType = iota
	NT_STRING NodeType = iota
	NT_BOOL   NodeType = iota
	NT_NULL   NodeType = iota

	INDENT_ON       = 0
	INDENT_OFF_ONCE = 1
	INDENT_OFF      = -1
)

type NodeI interface {
	GetName() string
	setName(newName string)
	GetNodeType() NodeType
	IsContainer() bool
	String() string
	JsonValue() string
	JsonValueIndented(tab int) string
}

type NodeC interface {
	// From parent node
	NodeI
	GetValues() []NodeI
	Len() int
	Clear()
	Add(node NodeI) error
	GetNodeWithName(name string) NodeI
	Remove(nodeRemove NodeI) error
}

var (
	_ NodeC = (*JsonObject)(nil)
	_ NodeC = (*JsonList)(nil)
	_ NodeI = (*JsonString)(nil)
	_ NodeI = (*JsonNumber)(nil)
	_ NodeI = (*JsonBool)(nil)
	_ NodeI = (*JsonNull)(nil)
)

//
// Base node (parent) interface (NodeI) and properties
//
type jsonParentNode struct {
	name string
	nt   NodeType
}

func NewJsonParentNode(name string, nt NodeType) jsonParentNode {
	return jsonParentNode{name: name, nt: nt}
}

func (n *jsonParentNode) GetName() string {
	return n.name
}

func (n *jsonParentNode) setName(name string) {
	n.name = name
}

func (n *jsonParentNode) GetNodeType() NodeType {
	return n.nt
}

func (n *jsonParentNode) IsContainer() bool {
	return n.nt == NT_LIST || n.nt == NT_OBJECT
}

//
// Objects node is a ParentNode and a value of type map[string]*NodeI
//
type JsonObject struct {
	jsonParentNode
	value map[string]*NodeI
}

func NewJsonType(name string, nodeType NodeType) NodeI {
	switch nodeType {
	case NT_OBJECT:
		return NewJsonObject(name)
	case NT_LIST:
		return NewJsonList(name)
	case NT_NUMBER:
		return NewJsonNumber(name, 0)
	case NT_STRING:
		return NewJsonString(name, "")
	case NT_BOOL:
		return NewJsonBool(name, false)
	case NT_NULL:
		return NewJsonNull(name)
	}
	panic(fmt.Sprintf("NewJsonType: Invalid node type '%d'. Should never get here!", nodeType))
}

func NewJsonObject(name string) *JsonObject {
	return &JsonObject{jsonParentNode: NewJsonParentNode(name, NT_OBJECT), value: make(map[string]*NodeI)}
}

func (n *JsonObject) GetNodeWithName(name string) NodeI {
	v := n.value[name]
	if v == nil {
		return nil
	}
	return *v
}

func (n *JsonObject) GetSortedKeys() []string {
	keys := make([]string, 0, len(n.value))
	for k := range n.value {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (n *JsonObject) GetValues() []NodeI {
	values := make([]NodeI, 0, len(n.value))
	for _, v := range n.value {
		values = append(values, *v)
	}
	return values
}

func (n *JsonObject) GetValuesSorted() []NodeI {
	values := make([]NodeI, 0, len(n.value))
	for _, k := range n.GetSortedKeys() {
		values = append(values, *n.value[k])
	}
	return values
}

func (n *JsonObject) Clear() {
	n.value = make(map[string]*NodeI)
}

func (n *JsonObject) Len() int {
	return len(n.value)
}

func (n *JsonObject) Add(node NodeI) error {
	if node.GetName() == "" {
		return fmt.Errorf("a node in a JsonObject container must have a name")
	}
	if n.value[node.GetName()] != nil {
		return fmt.Errorf("duplicate name [%s] in JsonObject container with name [%s]", node.GetName(), n.name)
	}
	n.value[node.GetName()] = &node
	return nil
}

func (n *JsonObject) JsonValueIndented(tab int) string {
	return stringValueTabIndent(n, tab, 1, INDENT_ON)
}

func (n *JsonObject) JsonValue() string {
	return stringValueTabIndent(n, 0, 0, INDENT_OFF)
}

func (n *JsonObject) String() string {
	return n.JsonValue()
}

func (n *JsonObject) Remove(nodeRemove NodeI) error {
	_, ok := n.value[nodeRemove.GetName()]
	if ok {
		delete(n.value, nodeRemove.GetName())
	} else {
		return fmt.Errorf("no matching node [%s] found in parent object node [%s]", nodeRemove.GetName(), n.name)
	}
	return nil
}

//
// List node is a ParentNode and a value of type []NodeI
//
type JsonList struct {
	jsonParentNode
	value []*NodeI
}

func NewJsonList(name string) *JsonList {
	return &JsonList{jsonParentNode: NewJsonParentNode(name, NT_LIST), value: make([]*NodeI, 0)}
}

func (n *JsonList) GetNodeAt(i int) NodeI {
	return *n.value[i]
}

func (n *JsonList) GetNodeWithName(name string) NodeI {
	for _, v := range n.value {
		if (*v).GetName() != "" && (*v).GetName() == name {
			return *v
		}
	}
	return nil
}

func (n *JsonList) Add(node NodeI) error {
	if node.GetNodeType() == NT_OBJECT {
		obj := node.(*JsonObject)
		if obj.GetName() == "" && obj.Len() == 1 {
			node = obj.GetValues()[0]
		}
	}
	n.value = append(n.value, &node)
	return nil
}

func (n *JsonList) Clear() {
	n.value = make([]*NodeI, 0)
}

func (n *JsonList) Len() int {
	return len(n.value)
}

func (n *JsonList) GetValues() []NodeI {
	values := make([]NodeI, 0, len(n.value))
	for _, v := range n.value {
		values = append(values, *v)
	}
	return values
}

func (n *JsonList) JsonValueIndented(tab int) string {
	return stringValueTabIndent(n, tab, 1, INDENT_ON)
}

func (n *JsonList) JsonValue() string {
	return stringValueTabIndent(n, 0, 0, INDENT_OFF)
}

func (n *JsonList) String() string {
	return n.JsonValue()
}

func (n *JsonList) Remove(nodeRemove NodeI) error {
	newList := make([]*NodeI, len(n.value)-1)
	newPos := 0
	for _, nx := range n.value {
		if nodeRemove != *nx {
			newList[newPos] = nx
			newPos++
			if newPos >= len(n.value) {
				return fmt.Errorf("no matching node found in parent list node")
			}
		}
	}
	n.value = newList
	return nil
}

//
// String node is a ParentNode and a value of type string
//
type JsonString struct {
	jsonParentNode
	value string
}

func NewJsonString(name string, value string) *JsonString {
	return &JsonString{jsonParentNode: NewJsonParentNode(name, NT_STRING), value: value}
}

func (n *JsonString) GetValue() string {
	return n.value
}

func (n *JsonString) SetValue(newValue string) {
	n.value = newValue
}

func (n *JsonString) JsonValueIndented(tab int) string {
	return stringValueTabIndent(n, tab, 1, INDENT_ON)
}

func (n *JsonString) JsonValue() string {
	return stringValueTabIndent(n, 0, 0, INDENT_OFF)
}

func (n *JsonString) String() string {
	return n.value
}

//
// Number node is a ParentNode and a value of type float64
//
type JsonNumber struct {
	jsonParentNode
	value float64
}

func NewJsonNumber(name string, value float64) *JsonNumber {
	return &JsonNumber{jsonParentNode: NewJsonParentNode(name, NT_NUMBER), value: value}
}

func (n *JsonNumber) GetValue() float64 {
	return n.value
}

func (n *JsonNumber) GetIntValue() int64 {
	return int64(n.value)
}

func (n *JsonNumber) SetValue(newValue float64) {
	n.value = newValue
}

func (n *JsonNumber) SetIntValue(newValue int64) {
	n.value = float64(newValue)
}

func (n *JsonNumber) JsonValueIndented(tab int) string {
	return stringValueTabIndent(n, tab, 1, INDENT_ON)
}

func (n *JsonNumber) JsonValue() string {
	return stringValueTabIndent(n, 0, 0, INDENT_OFF)
}

func (n *JsonNumber) String() string {
	s := fmt.Sprintf("%f", n.value)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

//
// Bool node is a ParentNode and a value of type bool
//
type JsonBool struct {
	jsonParentNode
	value bool
}

func NewJsonBool(name string, value bool) *JsonBool {
	return &JsonBool{jsonParentNode: NewJsonParentNode(name, NT_BOOL), value: value}
}

func (n *JsonBool) GetValue() bool {
	return n.value
}

func (n *JsonBool) SetValue(newValue bool) {
	n.value = newValue
}

func (n *JsonBool) JsonValueIndented(tab int) string {
	return stringValueTabIndent(n, tab, 1, INDENT_ON)
}

func (n *JsonBool) JsonValue() string {
	return stringValueTabIndent(n, 0, 0, INDENT_OFF)
}

func (n *JsonBool) String() string {
	return fmt.Sprintf("%t", n.value)
}

//
// Nill node is a ParentNode that has no value
//
type JsonNull struct {
	jsonParentNode
}

func NewJsonNull(name string) *JsonNull {
	return &JsonNull{jsonParentNode: NewJsonParentNode(name, NT_NULL)}
}

func (n *JsonNull) JsonValueIndented(tab int) string {
	return stringValueTabIndent(n, tab, 1, INDENT_ON)
}

func (n *JsonNull) JsonValue() string {
	return stringValueTabIndent(n, 0, 0, INDENT_OFF)
}

func (n *JsonNull) String() string {
	return "null"
}

func CreateAndReturnNodeAtPath(root NodeI, path *Path, nodeType NodeType) (NodeI, error) {
	if path.IsEmpty() {
		return nil, fmt.Errorf("cannot create a node from an empty path")
	}
	if !root.IsContainer() {
		return nil, fmt.Errorf("cannot create a node root node is not a container")
	}
	cNode := root.(NodeC)

	rootPath := path.PathParent()
	name := path.StringLast()
	if rootPath.IsEmpty() {
		n := cNode.GetNodeWithName(name)
		if n != nil {
			return n, nil
		}
		ret := NewJsonType(name, nodeType)
		err := cNode.Add(ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	}

	for _, nn := range rootPath.Paths() {
		n := cNode.GetNodeWithName(nn)
		if n == nil {
			n = NewJsonObject(nn)
			err := cNode.Add(n)
			if err != nil {
				return nil, err
			}
		}
		if !n.IsContainer() {
			return nil, fmt.Errorf("found node at [%s] but it is not a container node", nn)
		}
		cNode = n.(NodeC)
	}

	ret := cNode.GetNodeWithName(name)
	if ret == nil {
		ret = NewJsonType(name, nodeType)
		cNode.Add(ret)
	} else {
		if ret.GetNodeType() != nodeType {
			return nil, fmt.Errorf("found node at [%s] but it is not a %s node", path, GetNodeTypeName(nodeType))
		}
	}
	return ret, nil
}

func Remove(root, node NodeI) error {
	parentNode, found := FindParentNode(root, node)
	if found {
		if parentNode == nil {
			return fmt.Errorf("cannot remove node as it does not have a parent")
		} else {
			switch parentNode.GetNodeType() {
			case NT_LIST:
				err := parentNode.(*JsonList).Remove(node)
				if err != nil {
					return err
				}
			case NT_OBJECT:
				err := parentNode.(*JsonObject).Remove(node)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("cannot remove node as its parent is not a container node")
			}
		}
	} else {
		return fmt.Errorf("cannot remove node as it was not found in root tree")
	}
	return nil
}

func Clone(n NodeI, newName string, cloneLeafNodeData bool) NodeI {
	if n.IsContainer() {
		cl := NewJsonType(newName, n.GetNodeType())
		for _, v := range n.(NodeC).GetValues() {
			nn := Clone(v, v.GetName(), cloneLeafNodeData)
			cl.(NodeC).Add(nn)
		}
		return cl
	} else {
		nn := NewJsonType(newName, n.GetNodeType())
		if cloneLeafNodeData {
			switch n.GetNodeType() {
			case NT_BOOL:
				nn.(*JsonBool).SetValue(n.(*JsonBool).GetValue())
			case NT_NUMBER:
				nn.(*JsonNumber).SetValue(n.(*JsonNumber).GetValue())
			case NT_STRING:
				nn.(*JsonString).SetValue(n.(*JsonString).GetValue())
			}
		}
		return nn
	}
}

func Rename(root, node NodeI, newName string) error {
	if node.GetName() == "" {
		return fmt.Errorf("cannot rename. This node has no name!")
	}
	parentNode, found := FindParentNode(root, node)
	if found {
		if parentNode == nil {
			node.setName(newName)
		} else {
			switch parentNode.GetNodeType() {
			case NT_LIST:
				node.setName(newName)
			case NT_OBJECT:
				po := parentNode.(*JsonObject)
				if po.GetNodeWithName(newName) != nil {
					return fmt.Errorf("cannot rename node. Parent already has a node with the new name")
				}
				delete(parentNode.(*JsonObject).value, node.GetName())
				node.setName(newName)
				err := parentNode.(*JsonObject).Add(node)
				if err != nil {
					return err
				}

			default:
				return fmt.Errorf("cannot rename node as its parent is not a container node")
			}
		}
	} else {
		return fmt.Errorf("cannot rename node as it was not found in root tree")
	}
	return nil
}

func GetNodeTypeName(tt NodeType) string {
	switch tt {
	case NT_OBJECT:
		return "OBJECT"
	case NT_LIST:
		return "LIST"
	case NT_NUMBER:
		return "NUMBER"
	case NT_STRING:
		return "STRING"
	case NT_BOOL:
		return "BOOL"
	case NT_NULL:
		return "NULL"
	}
	return "UNKNOWN"
}

func stringValueTabIndent(n NodeI, tab, indent int, useIndent int) string {
	var sb strings.Builder
	p := Padding(tab, indent, useIndent)
	indent++
	if useIndent > 0 {
		useIndent--
	}
	sb.WriteString(p)
	if n.GetName() != "" {
		sb.WriteByte('"')
		sb.WriteString(n.GetName())
		sb.WriteByte('"')
		sb.WriteByte(':')
		sb.WriteByte(' ')
	}
	switch n.GetNodeType() {
	case NT_LIST:
		sb.WriteByte('[')
		nL := n.(*JsonList)
		c := len(nL.value) - 1
		for i, v := range nL.value {
			if (*v).GetName() == "" {
				sb.WriteString(stringValueTabIndent(*v, tab, indent, INDENT_ON))
			} else {
				sb.WriteString(Padding(tab, indent, useIndent))
				sb.WriteByte('{')
				sb.WriteString(stringValueTabIndent(*v, tab, indent, INDENT_OFF_ONCE))
				sb.WriteByte('}')
			}
			if i < c {
				sb.WriteByte(',')
			}
		}
		sb.WriteString(p)
		sb.WriteByte(']')
	case NT_OBJECT:
		sb.WriteByte('{')
		nO := n.(*JsonObject)
		c := len(nO.value) - 1
		i := 0
		for _, v := range nO.value {
			sb.WriteString(stringValueTabIndent(*v, tab, indent, useIndent))
			if i < c {
				sb.WriteByte(',')
			}
			i++
		}
		sb.WriteString(p)
		sb.WriteByte('}')
	case NT_STRING:
		sb.WriteRune('"')
		sb.WriteString(EncodeQuotedString(n.String()))
		sb.WriteRune('"')
	default:
		sb.WriteString(n.String())
	}
	return sb.String()
}

func baseEquals(a, b NodeC) bool {
	if a.GetNodeType() != b.GetNodeType() {
		return false
	}
	if a.GetName() != b.GetName() {
		return false
	}
	if a.IsContainer() != b.IsContainer() {
		return false
	}
	return true
}
