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

type TokenType uint16

const (
	TT_ARRAY_OPEN    TokenType = iota
	TT_ARRAY_CLOSE   TokenType = iota
	TT_OBJECT_OPEN   TokenType = iota
	TT_OBJECT_CLOSE  TokenType = iota
	TT_COMMA         TokenType = iota
	TT_QUOTED_STRING TokenType = iota
	TT_BOOL_TRUE     TokenType = iota
	TT_BOOL_FALSE    TokenType = iota
	TT_NUMBER        TokenType = iota
	TT_COLON         TokenType = iota
	TT_NULL          TokenType = iota
)

var (
	padding   = string("                                                                                                    ")
	digitsHex = []rune{
		'0', // 0
		'1', // 1
		'2', // 2
		'3', // 3
		'4', // 4
		'5', // 5
		'6', // 6
		'7', // 7
		'8', // 8
		'9', // 9
		'A', // 10
		'B', // 11
		'C', // 12
		'D', // 13
		'E', // 14
		'F', // 15
	}

	hexDigits = []uint16{
		0,  // 0
		1,  // 1
		2,  // 2
		3,  // 3
		4,  // 4
		5,  // 5
		6,  // 6
		7,  // 7
		8,  // 8
		9,  // 9
		0,  // :    (colon)
		0,  // ;    (semi-colon)
		0,  // <    (less than)
		0,  // =    (equal sign)
		0,  // >    (greater than)
		0,  // ?    (question mark)
		0,  // @    (AT symbol)
		10, // A
		11, // B
		12, // C
		13, // D
		14, // E
		15, // F
	}

	NULL uint16 = 0x0000

	NUM    uint16 = 0x0001
	ALF_LC uint16 = 0x0002
	ALF_UC uint16 = 0x0004
	ALF    uint16 = 0x0008

	SPACE  uint16 = 0x0010 // A Space
	WS     uint16 = 0x0020 // White Space
	US     uint16 = 0x0040 // Underscore
	DQUOTE uint16 = 0x0080 // Double Quote

	SQUOTE uint16 = 0x0100 // Single Quote
	DOT    uint16 = 0x0200
	HYPHEN uint16 = 0x0400
	ESCAPE uint16 = 0x0800 // The \ character

	DOLLER  uint16 = 0x1000
	TRUE_C  uint16 = 0x2000 // Chars in 'true'
	FALSE_C uint16 = 0x4000 // Chars in 'false'
	NULL_C  uint16 = 0x8000 // Chars in 'null'

	NCNAME       uint16 = ALF | NUM | US | DOT | HYPHEN
	FIRST_NCNAME uint16 = ALF | US | DOLLER
	ALPHA_NUM    uint16 = ALF | NUM
	ALPHA        uint16 = ALF
	QUOTE        uint16 = DQUOTE | SQUOTE

	masks = []uint16{
		WS, WS, WS, WS, WS, WS, WS, WS,
		WS, WS, WS, WS, WS, WS, WS, WS,
		WS, WS, WS, WS, WS, WS, WS, WS,
		WS, WS, WS, WS, WS, WS, WS, WS,
		SPACE | WS,                      // SP   (Space)
		NULL,                            // !    (exclamation mark)
		DQUOTE,                          // "    (double quote)
		NULL,                            // #    (number sign)
		DOLLER,                          // $    (dollar sign)
		NULL,                            // %    (percent)
		NULL,                            // &    (ampersand)
		SQUOTE,                          // '    (single quote)
		NULL,                            // (  (left/open parenthesis)
		NULL,                            // )  (right/closing parenth.)
		NULL,                            // *    (asterisk)
		NULL | NUM,                      // +    (plus)
		NULL,                            // ,    (comma)
		HYPHEN | NUM,                    // -    (minus or dash)
		DOT | NUM,                       // .    (dot)
		NULL,                            // /    (forward slash)
		NUM,                             // 0
		NUM,                             // 1
		NUM,                             // 2
		NUM,                             // 3
		NUM,                             // 4
		NUM,                             // 5
		NUM,                             // 6
		NUM,                             // 7
		NUM,                             // 8
		NUM,                             // 9
		NULL,                            // :    (colon)
		NULL,                            // ;    (semi-colon)
		NULL,                            // <    (less than)
		NULL,                            // =    (equal sign)
		NULL,                            // >    (greater than)
		NULL,                            // ?    (question mark)
		NULL,                            // @    (AT symbol)
		ALF_UC | ALF,                    // A
		ALF_UC | ALF,                    // B
		ALF_UC | ALF,                    // C
		ALF_UC | ALF,                    // D
		ALF_UC | ALF,                    // E
		ALF_UC | ALF,                    // F
		ALF_UC | ALF,                    // G
		ALF_UC | ALF,                    // H
		ALF_UC | ALF,                    // I
		ALF_UC | ALF,                    // J
		ALF_UC | ALF,                    // K
		ALF_UC | ALF,                    // L
		ALF_UC | ALF,                    // M
		ALF_UC | ALF,                    // N
		ALF_UC | ALF,                    // O
		ALF_UC | ALF,                    // P
		ALF_UC | ALF,                    // Q
		ALF_UC | ALF,                    // R
		ALF_UC | ALF,                    // S
		ALF_UC | ALF,                    // T
		ALF_UC | ALF,                    // U
		ALF_UC | ALF,                    // V
		ALF_UC | ALF,                    // W
		ALF_UC | ALF,                    // X
		ALF_UC | ALF,                    // Y
		ALF_UC | ALF,                    // Z
		NULL,                            // [    (left/opening bracket)
		ESCAPE,                          // \    (back slash)
		NULL,                            // ]    (right/closing bracket)
		NULL,                            // ^    (caret/circumflex)
		US,                              // _    (underscore)
		SQUOTE,                          // `
		FALSE_C | ALF_LC | ALF,          // a
		ALF_LC | ALF,                    // b
		ALF_LC | ALF,                    // c
		ALF_LC | ALF,                    // d
		FALSE_C | TRUE_C | ALF_LC | ALF, // e
		FALSE_C | ALF_LC | ALF,          // f
		ALF_LC | ALF,                    // g
		ALF_LC | ALF,                    // h
		ALF_LC | ALF,                    // i
		ALF_LC | ALF,                    // j
		ALF_LC | ALF,                    // k
		NULL_C | FALSE_C | ALF_LC | ALF, // l
		ALF_LC | ALF,                    // m
		NULL_C | ALF_LC | ALF,           // n
		ALF_LC | ALF,                    // o
		ALF_LC | ALF,                    // p
		ALF_LC | ALF,                    // q
		TRUE_C | ALF_LC | ALF,           // r
		FALSE_C | ALF_LC | ALF,          // s
		TRUE_C | ALF_LC | ALF,           // t
		NULL_C | TRUE_C | ALF_LC | ALF,  // u
		ALF_LC | ALF,                    // v
		ALF_LC | ALF,                    // w
		ALF_LC | ALF,                    // x
		ALF_LC | ALF,                    // y
		ALF_LC | ALF,                    // z
		NULL,                            // {    (left/opening brace)
		NULL,                            // |    (vertical bar)
		NULL,                            // }    (right/closing brace)
		NULL,                            // ~    (tilde)
		NULL,                            // DEL    (delete)
		NULL,                            // Padding!
	}
)

func Padding(tab, indent int, useIndent int) string {
	if useIndent == 0 && tab > 0 {
		return "\n" + padding[:tab*indent]
	}
	return ""
}

func IntToHexChar(c uint) rune {
	return digitsHex[c&0x00000F]
}

func IntToHexChar2(c uint) []rune {
	rr := make([]rune, 2)
	rr[0] = IntToHexChar(c >> 4)
	rr[1] = IntToHexChar(c)
	return rr
}

func IntToHexChar4(c uint) []rune {
	rr := make([]rune, 4)
	rr[0] = IntToHexChar(c >> 12)
	rr[1] = IntToHexChar(c >> 8)
	rr[2] = IntToHexChar(c >> 4)
	rr[3] = IntToHexChar(c)
	return rr
}

func HexCharToInt(c byte) uint16 {
	if c >= 'a' {
		c = c - 32 // Convert to uppercase
	}
	if c < '0' || c >= 'G' { // Range check!
		panic(fmt.Sprintf("parser:scanner:HexToInt: Invalid hex digit '%c'", c))
	}
	return hexDigits[c-'0'] // return value
}

func GetTokenTypeName(tt TokenType) string {
	switch tt {
	case TT_ARRAY_OPEN:
		return "ARRAY_OPEN"
	case TT_ARRAY_CLOSE:
		return "ARRAY_CLOSE"
	case TT_OBJECT_OPEN:
		return "OBJECT_OPEN"
	case TT_OBJECT_CLOSE:
		return "OBJECT_CLOSE"
	case TT_COMMA:
		return "COMMA"
	case TT_QUOTED_STRING:
		return "QUOTED_STRING"
	case TT_NUMBER:
		return "NUMBER"
	case TT_COLON:
		return "COLON"
	case TT_BOOL_TRUE:
		return "TRUE"
	case TT_BOOL_FALSE:
		return "FALSE"
	case TT_NULL:
		return "NULL"
	}
	return "UNKNOWN"
}

func CharIsAll(ch byte, mask uint16) bool {
	return ((masks[ch&0x007F] & mask) == mask)
}

func CharIsAny(ch byte, mask uint16) bool {
	return ((masks[ch&0x007F] & mask) != 0)
}
