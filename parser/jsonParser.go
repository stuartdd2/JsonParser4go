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
)

const (
	pad string = "                  "
)

func Parse(json []byte) (node NodeC, err error) {
	defer func() {
		r := recover()
		if r != nil {
			node = nil
			err = fmt.Errorf("parser Error: %v", r)
		}
	}()
	sc := NewScanner(json).SkipSpace()
	var root NodeC
	tok := sc.Next()
	switch tok {
	case '[':
		root = parseList(sc, "")
	case '{':
		root = parseObject(sc, "")
	default:
		err = fmt.Errorf("parser Error: %s", sc.Diag("?"))
		return
	}
	node = root
	err = nil
	return
}

func parseObject(sc *Scanner, name string) NodeC {
	root := NewJsonObject(name)
	var err error
	for {
		toc := sc.NextToken()
		if toc.IsObjectClose() {
			return NewJsonObject(name)
		}
		if !toc.IsQuotedString() {
			panic(fmt.Sprintf("object name is invalid. Found '%s'. %s ", toc.GetStringValue(), sc.Diag(toc.GetStringValue())))
		}
		name := toc.GetStringValue()
		toc = sc.NextToken()
		if !toc.IsColon() {
			panic(fmt.Sprintf("object name not followed by a ':'. Found '%s'. %s ", toc.GetStringValue(), sc.Diag(toc.GetStringValue())))
		}
		toc = sc.NextToken()
		switch toc.GetType() {
		case TT_QUOTED_STRING:
			_, err = root.Add(NewJsonString(name, toc.GetStringValue()))
		case TT_NUMBER:
			_, err = root.Add(NewJsonNumber(name, toc.GetNumberValue()))
		case TT_BOOL_TRUE:
			_, err = root.Add(NewJsonBool(name, true))
		case TT_BOOL_FALSE:
			_, err = root.Add(NewJsonBool(name, false))
		case TT_NULL:
			_, err = root.Add(NewJsonNull(name))
		case TT_ARRAY_OPEN:
			_, err = root.Add(parseList(sc, name))
		case TT_OBJECT_OPEN:
			_, err = root.Add(parseObject(sc, name))
		default:
			panic(fmt.Sprintf("unrecognised token '%s'. %s ", toc.GetStringValue(), sc.Diag(toc.GetStringValue())))
		}
		if err != nil {
			panic(err.Error())
		}
		toc = sc.NextToken()
		if toc.IsObjectClose() {
			return root
		}
		if !toc.IsComma() {
			panic(fmt.Sprintf("expected a ',' seperator. Found '%s'. %s ", toc.GetStringValue(), sc.Diag(toc.GetStringValue())))
		}
	}
}

func parseList(sc *Scanner, name string) NodeC {
	root := NewJsonList(name)
	for {
		toc := sc.NextToken()
		if toc.IsArrayClose() {
			return root
		}
		switch toc.GetType() {
		case TT_QUOTED_STRING:
			root.Add(NewJsonString("", toc.GetStringValue()))
		case TT_NUMBER:
			root.Add(NewJsonNumber("", toc.GetNumberValue()))
		case TT_BOOL_TRUE:
			root.Add(NewJsonBool("", true))
		case TT_BOOL_FALSE:
			root.Add(NewJsonBool("", false))
		case TT_NULL:
			root.Add(NewJsonNull(""))
		case TT_OBJECT_OPEN:
			root.Add(parseObject(sc, ""))
		default:
			panic(fmt.Sprintf("unrecognised token '%s'. %s ", toc.GetStringValue(), sc.Diag(toc.GetStringValue())))
		}
		toc = sc.NextToken()
		if toc.IsArrayClose() {
			return root
		}
		if !toc.IsComma() {
			panic(fmt.Sprintf("expected a ',' seperator. Found '%s'. %s ", toc.GetStringValue(), sc.Diag(toc.GetStringValue())))
		}
	}
}
