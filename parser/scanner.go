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

type Token struct {
	text string
	pos  int
	tok  TokenType
}

type Scanner struct {
	text []byte
	pos  int
	max  int
}

func NewScanner(s []byte) *Scanner {
	return &Scanner{text: s, pos: 0, max: len(s)}
}

func (s *Scanner) Diag(tok string) string {
	f := s.pos - 20
	if f < 0 {
		f = 0
	}
	t := f + 40
	if t > s.max {
		t = s.max
	}
	p := s.pos - len(tok)
	return fmt.Sprintf("Scanner: pos: %d len: %d. About here >>>%s|%s<<<", s.pos, s.max, s.text[f:p], s.text[p:t])
}

func (s *Scanner) Next() byte {
	if s.pos < s.max {
		c := s.text[s.pos]
		s.pos++
		return c
	}
	return 0
}

func (s *Scanner) HasNext() bool {
	return (s.pos < s.max)
}

func (s *Scanner) Back() *Scanner {
	if s.pos > 0 {
		s.pos--
	}
	return s
}

func (s *Scanner) SkipSpace() *Scanner {
	for s.HasNext() {
		if s.Next() > ' ' {
			s.Back()
			return s
		}
	}
	return s
}

func (s *Scanner) SkipToNext(c byte) *Scanner {
	for s.HasNext() {
		if s.Next() == c {
			s.Back()
			return s
		}
	}
	return s
}

func (s *Scanner) IsNext(mask uint16) bool {
	return CharIsAny(s.text[s.pos], mask)
}

func (s *Scanner) Reset() *Scanner {
	s.pos = 0
	return s
}
func (s *Scanner) PeekToken() *Token {
	p := s.pos
	t := s.NextToken()
	s.pos = p
	return t
}

func (s *Scanner) NextToken() *Token {
	s.SkipSpace()
	p := s.pos
	if s.HasNext() {
		c := s.Next()
		if c == '{' {
			return NewToken("", p, TT_OBJECT_OPEN)
		}
		if c == '}' {
			return NewToken("", p, TT_OBJECT_CLOSE)
		}
		if c == ',' {
			return NewToken("", p, TT_COMMA)
		}
		if c == ':' {
			return NewToken("", p, TT_COLON)
		}
		if c == '"' {
			return NewToken(s.scanQuotedString(c), p, TT_QUOTED_STRING)
		}
		if c == '[' {
			return NewToken("", p, TT_ARRAY_OPEN)
		}
		if c == ']' {
			return NewToken("", p, TT_ARRAY_CLOSE)
		}
		if c == 't' {
			i := s.skipValueWithMask(TRUE_C)
			if i != 4 {
				panic(fmt.Sprintf("Boolean 'true' Must be 4 chars long starting with 't' containing 'r' 'u' and 'e'. %s", s.Diag("")))
			}
			return NewToken("true", p, TT_BOOL_TRUE)
		}
		if c == 'f' {
			i := s.skipValueWithMask(FALSE_C)
			if i != 5 {
				panic(fmt.Sprintf("Boolean 'false' Must be 5 chars long starting with 'f' containing 'a' 'l' 's' and 'e'. %s", s.Diag("")))
			}
			return NewToken("false", p, TT_BOOL_FALSE)
		}
		if c == 'n' {
			i := s.skipValueWithMask(NULL_C)
			if i != 4 {
				panic(fmt.Sprintf("'null' Must be 4 chars long starting with 'n' containing 'a' 'l' 's' and 'e'. %s", s.Diag("")))
			}
			return NewToken("null", p, TT_NULL)
		}
		if CharIsAny(c, NUM) {
			s.Back()
			return NewToken(s.scanValueWithMask(NUM), p, TT_NUMBER)
		}
		panic(fmt.Sprintf("unrecognised token. '%c'. %s", rune(c), s.Diag(" ")))
	}
	panic(fmt.Sprintf("unexpected end of input. %s", s.Diag("")))
}

func EncodeQuotedString(inStr string) string {
	var sb strings.Builder
	for _, c := range inStr {
		switch c {
		case 0x0A: // Line Feed
			sb.WriteString("\\n")
		case 0x0D: // Carrage return
		case 0x08: // Back Space
			sb.WriteString("\\b")
		case 0x0C: // Form feed
			sb.WriteString("\\f")
		case 0x09: // Horizontal tab
			sb.WriteString("\\t")
		case '\\': // Horizontal tab
			sb.WriteString("\\\\")
		case '"': //Double quotes need to be escaped
			sb.WriteString("\\\"")
		default:
			if c > 255 {
				r := IntToHexChar4(uint(c))
				sb.WriteString("\\u")
				sb.WriteRune(r[0])
				sb.WriteRune(r[1])
				sb.WriteRune(r[2])
				sb.WriteRune(r[3])
			} else {
				if c > 127 {
					r := IntToHexChar2(uint(c))
					sb.WriteString("\\x")
					sb.WriteRune(r[0])
					sb.WriteRune(r[1])
				} else {
					sb.WriteRune(c)
				}
			}
		}
	}
	return sb.String()
}

func (s *Scanner) scanQuotedString(delim byte) string {
	var sb strings.Builder
	for s.HasNext() {
		c := s.Next()
		if c == '\\' {
			c = s.Next()
			switch c {
			case 'b': // 0x08
				sb.WriteString("\b")
			case 'r': // 0x0D
				sb.WriteString("\r")
			case 'n': // 0x0A
				sb.WriteString("\n")
			case 'f': // 0x0C
				sb.WriteString("\f")
			case 't': // 0x09
				sb.WriteString("\t")
			case 'u':
				b1 := s.readUInt16()
				b2 := s.readUInt16()
				sb.WriteRune(rune((b1 << 8) | b2))
			case 'x':
				b1 := s.readUInt16()
				sb.WriteRune(rune(b1))
			default:
				sb.WriteByte(c)
			}
		} else {
			if c == delim {
				return sb.String()
			}
			if s.HasNext() {
				if s.text[s.pos] == delim {
					s.Next()
					sb.WriteByte(c)
					return sb.String()
				}
			} else {
				panic(fmt.Sprintf("unterminated quoted String. %s", s.Diag(" ")))
			}
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

func (s *Scanner) readUInt16() uint16 {
	var b0 uint16
	var b1 uint16
	if s.HasNext() {
		b0 = HexCharToInt(s.Next())
	} else {
		panic(fmt.Sprintf("Unexpected end of input. %s", s.Diag(" ")))
	}
	if s.HasNext() {
		b1 = HexCharToInt(s.Next())
	} else {
		panic(fmt.Sprintf("Unexpected end of input. %s", s.Diag(" ")))
	}
	return (b0 << 4) + b1
}

func (s *Scanner) skipValueWithMask(mask uint16) int {
	i := 1
	for s.HasNext() {
		c := s.Next()
		if !CharIsAny(c, mask) {
			s.Back()
			return i
		}
		i++
	}
	return i
}

func (s *Scanner) scanValueWithMask(mask uint16) string {
	var sb strings.Builder
	for s.HasNext() {
		c := s.Next()
		if CharIsAny(c, mask) {
			sb.WriteByte(c)
		} else {
			s.Back()
			return sb.String()
		}
	}
	return sb.String()
}

// -------------------------------------------------------------
// Token logic
// -------------------------------------------------------------
func NewToken(txt string, pos int, tokenType TokenType) *Token {
	return &Token{text: txt, pos: pos, tok: tokenType}
}

func (t *Token) GetStringValue() string {
	return t.text
}

func (t *Token) GetNumberValue() float64 {
	if s, err := strconv.ParseFloat(t.text, 64); err == nil {
		return s
	}
	panic("Number conversion error")
}

func (t *Token) GetType() TokenType {
	return t.tok
}

func (t *Token) IsType(ofType TokenType) bool {
	return t.tok == ofType
}

func (t *Token) IsArrayOpen() bool {
	return t.tok == TT_ARRAY_OPEN
}

func (t *Token) IsArrayClose() bool {
	return t.tok == TT_ARRAY_CLOSE
}

func (t *Token) IsObjectOpen() bool {
	return t.tok == TT_OBJECT_OPEN
}

func (t *Token) IsObjectClose() bool {
	return t.tok == TT_OBJECT_CLOSE
}

func (t *Token) IsComma() bool {
	return t.tok == TT_COMMA
}

func (t *Token) IsColon() bool {
	return t.tok == TT_COLON
}

func (t *Token) IsQuotedString() bool {
	return t.tok == TT_QUOTED_STRING
}

func (t *Token) IsBool() bool {
	return t.tok == TT_BOOL_TRUE || t.tok == TT_BOOL_FALSE
}

func (t *Token) IsNumber() bool {
	return t.tok == TT_NUMBER
}

func (t *Token) String() string {
	return fmt.Sprintf("Token : (%d) %s : %s", t.tok, GetTokenTypeName(t.tok), t.GetStringValue())
}
