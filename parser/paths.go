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
	"strings"
)

var (
	parserEmptyPath = make([]string, 0)
)

type Trail struct {
	trail []*NodeI
	index []int
	len   int
	size  int
	delim string
}

func NewTrail(size int, delim string) *Trail {
	return &Trail{trail: make([]*NodeI, size), index: make([]int, size), len: 0, size: size, delim: delim}
}

func (p *Trail) GetList() []*NodeI {
	l := make([]*NodeI, p.len)
	for i := 0; i < p.len; i++ {
		l[i] = p.trail[i]
	}
	return l
}

func (p *Trail) GetLast() NodeI {
	if p.len == 0 {
		return nil
	}
	return *p.trail[p.len-1]
}

func (p *Trail) GetIndex(i int) int {
	return p.index[i]
}

func (p *Trail) SetDelim(delim string) {
	p.delim = delim
}

func (p *Trail) GetDelim() string {
	return p.delim
}

func (p *Trail) Clear() {
	p.len = 0
}

func (p *Trail) Push(n NodeI, index int) bool {
	if p.len < p.size {
		p.trail[p.len] = &n
		p.index[p.len] = index
		p.len++
		return true
	}
	return false
}

func (p *Trail) Pop() (NodeI, int) {
	if p.len == 0 {
		return nil, -1
	}
	p.len--
	return *p.trail[p.len], p.index[p.len]
}

func (p *Trail) String() string {
	var sb strings.Builder
	for i, v := range p.GetList() {
		name := (*v).GetName()
		if name != "" {
			sb.WriteString((*v).GetName())
		} else {
			ind := p.GetIndex(i)
			if ind > 0 {
				sb.WriteString(fmt.Sprintf("%d", ind))
			}
		}
		if i < p.len-1 {
			sb.WriteString(p.GetDelim())
		}
	}
	return sb.String()
}

type Path struct {
	path  []string
	delim string
}

func NewDotPath(path string) *Path {
	return NewPath(path, ".")
}

func NewBarPath(path string) *Path {
	return NewPath(path, "|")
}

func NewPath(path, delim string) *Path {
	if path == "" {
		if delim == "" {
			return &Path{path: parserEmptyPath, delim: "."}
		}
		return &Path{path: parserEmptyPath, delim: delim}
	} else {
		if delim == "" {
			return &Path{path: strings.Split(path, "."), delim: "."}
		}
		return &Path{path: strings.Split(path, delim), delim: delim}
	}
}

func (p *Path) String() string {
	if p.IsEmpty() {
		return ""
	}
	var sb strings.Builder
	max := len(p.path) - 1
	for i, v := range p.path {
		sb.WriteString(v)
		if i < max {
			sb.WriteString(p.delim)
		}
	}
	return sb.String()
}

//
// Equal does NOT depend on the delim, only on the nodes
//
func (p *Path) Equal(anyPath *Path) bool {
	if anyPath == nil {
		return false
	}
	if len(p.path) == len(anyPath.path) {
		for i, v := range p.path {
			if v != anyPath.path[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (p *Path) PathAppend(p2 *Path) {
	p.path = append(p.path, p2.path...)
}

func (p *Path) BackToFront() *Path {
	rp := NewPath("", p.delim)
	for i := len(p.path) - 1; i >= 0; i-- {
		rp.path = append(rp.path, p.path[i])
	}
	return rp
}

func (p *Path) StringAppend(s string) {
	sPath := NewPath(s, p.delim)
	p.path = append(p.path, sPath.path...)
}

func (p *Path) Paths() []string {
	return p.path
}

func (p *Path) PathFirst(n int) *Path {
	if n >= len(p.path) {
		return &Path{p.path, p.delim}
	}
	if n <= 0 {
		return &Path{parserEmptyPath, p.delim}
	}
	return &Path{p.path[:n], p.delim}
}

func (p *Path) PathLast(n int) *Path {
	if n >= len(p.path) {
		return &Path{p.path, p.delim}
	}
	if n <= 0 {
		return &Path{parserEmptyPath, p.delim}
	}
	return &Path{p.path[len(p.path)-n:], p.delim}
}

func (p *Path) PathParent() *Path {
	return p.PathFirst(p.Len() - 1)
}

func (p *Path) PathAt(i int) *Path {
	return NewPath(p.StringAt(i), p.delim)
}

func (p *Path) StringAt(i int) string {
	if (i < 0) || p.IsEmpty() || i >= p.Len() {
		return ""
	}
	return p.path[i]
}

func (p *Path) StringFirst() string {
	if p.IsEmpty() {
		return ""
	}
	return p.path[0]
}

func (p *Path) StringLast() string {
	if p.IsEmpty() {
		return ""
	}
	return p.path[len(p.path)-1]
}

func (p *Path) GetDelim() string {
	return p.delim
}

func (p *Path) Len() int {
	return len(p.path)
}

func (p *Path) IsEmpty() bool {
	return len(p.path) == 0
}
