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
	"strconv"
	"strings"
)

func Find(node NodeI, path *Path) (NodeI, error) {
	if path.IsEmpty() {
		return nil, fmt.Errorf("no search paths were resolved. Path provided: '%s'", path)
	}
	for _, v := range path.path {
		if node.GetNodeType() == NT_LIST {
			ln := (node.(*JsonList))
			i, err := strconv.Atoi(v)
			if err != nil {
				n := ln.GetNodeWithName(v)
				if n == nil {
					return nil, fmt.Errorf("node for path: '%s' element: '%s' was not found", path, v)
				}
				node = n
			} else {
				l := ln.Len()
				if i < 0 || i >= l {
					return nil, fmt.Errorf("index out of bounds. Range: 0..%d Path provided: '%s' Actual:%d", l-1, path, i)
				}
				node = ln.GetNodeAt(i)
			}
		} else {
			if node.GetNodeType() == NT_OBJECT {
				ob := (node.(*JsonObject))
				n := ob.GetNodeWithName(v)
				if n == nil {
					return nil, fmt.Errorf("node for path: '%s' element: '%s' was not found", path, v)
				}
				node = n
			} else {
				return nil, fmt.Errorf("node for path: '%s' element: '%s' was not found", path, v)
			}
		}
	}
	return node, nil
}

func WalkNodeTreeForTrail(root NodeC, visitWithTrail func(*Trail, int) bool) (*Trail, bool) {
	trail := NewTrail(10, "|")
	i := walkNodeTreeForPaths(root, trail, 0, visitWithTrail)
	if i < 0 {
		return nil, false
	}
	return trail, true
}

func walkNodeTreeForPaths(node NodeC, trail *Trail, dep int, visitWithTrail func(*Trail, int) bool) int {
	listNode := node.GetNodeType() == NT_LIST
	for i, v := range node.GetValues() {
		if listNode {
			trail.Push(v, i)
			if visitWithTrail(trail, i) {
				return dep
			}
		} else {
			trail.Push(v, -i)
			if visitWithTrail(trail, -1) {
				return dep
			}
		}
		if v.IsContainer() {
			n := walkNodeTreeForPaths(v.(NodeC), trail, dep+1, visitWithTrail)
			if n >= 0 {
				return n
			}
		}
		trail.Pop()
	}
	return -1
}

func WalkNodeTree(root, target NodeI, onEachNode func(NodeI, NodeI, NodeI) bool) (NodeI, NodeI, bool) {
	return walkNodes(root, nil, target, onEachNode)
}

func walkNodes(walkFrom, walkFromParent, target NodeI, visitNode func(NodeI, NodeI, NodeI) bool) (NodeI, NodeI, bool) {
	if visitNode(walkFrom, walkFromParent, target) {
		return walkFrom, walkFromParent, true
	}
	switch walkFrom.GetNodeType() {
	case NT_LIST:
		walkFromList := (walkFrom.(*JsonList))
		for i := 0; i < walkFromList.Len(); i++ {
			n := walkFromList.GetNodeAt(i)
			wn, wf, ok := walkNodes(n, walkFrom, target, visitNode)
			if ok {
				return wn, wf, ok
			}
		}
	case NT_OBJECT:
		walkFromObj := (walkFrom.(*JsonObject))
		for _, n := range walkFromObj.GetValues() {
			wn, wf, ok := walkNodes(n, walkFrom, target, visitNode)
			if ok {
				return wn, wf, ok
			}
		}
	}
	return nil, nil, false
}

func DiagnosticList(node NodeI) string {
	sb := strings.Builder{}
	sb.WriteString("Diag\n")
	sb.WriteString(GetNodeTypeName(node.GetNodeType()))
	sb.WriteString(": N:'")
	sb.WriteString(node.GetName())
	sb.WriteString("'\n")
	diagWalk(node, 2, &sb)
	return strings.TrimSpace(sb.String())
}

func diagStr(node NodeI, nt NodeType, indent int, sb *strings.Builder) {
	sb.WriteString(pad[:indent])
	sb.WriteString(GetNodeTypeName(nt))
	sb.WriteString(": N:'")
	sb.WriteString(node.GetName())
	sb.WriteString("' ")
	switch nt {
	case NT_STRING:
		sb.WriteString(fmt.Sprintf("V:'%s'\n", (node.(*JsonString)).GetValue()))
	case NT_BOOL:
		sb.WriteString(fmt.Sprintf("V:'%t'\n", (node.(*JsonBool)).GetValue()))
	case NT_NUMBER:
		sb.WriteString(fmt.Sprintf("V:'%f'\n", (node.(*JsonNumber)).GetValue()))
	case NT_NULL:
		sb.WriteString("\n")
	case NT_OBJECT:
		sb.WriteString("\n")
		diagWalk(node, indent+2, sb)
	case NT_LIST:
		sb.WriteString("\n")
		diagWalk(node, indent+2, sb)
	}
}

func diagWalk(node NodeI, indent int, sb *strings.Builder) {
	switch node.GetNodeType() {
	case NT_LIST:
		cn := (node.(*JsonList))
		for i := 0; i < cn.Len(); i++ {
			n := cn.GetNodeAt(i)
			diagStr(n, n.GetNodeType(), indent, sb)
		}
	case NT_OBJECT:
		cn := (node.(*JsonObject))
		for _, key := range cn.GetSortedKeys() {
			n := cn.GetNodeWithName(key)
			diagStr(n, n.GetNodeType(), indent, sb)
		}
	}
}
