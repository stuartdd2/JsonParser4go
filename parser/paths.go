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
	"strings"
)

var (
	parserEmptyPath = make([]string, 0)
)

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

func (p *Path) PathAppend(p2 *Path) *Path {
	return &Path{path: append(p.path, p2.path...), delim: p.delim}
}

func (p *Path) StringAppend(s string) *Path {
	sPath := NewPath(s, p.delim)
	return &Path{path: append(p.path, sPath.path...), delim: p.delim}
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
