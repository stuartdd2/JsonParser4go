package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stuartdd/jsonParserGo/parser"
)

var (
	text = []byte(`
	{
		"data": {
			"positional": true
		},
		"datafile":"../data.json","screen": {"fullScreen":false,"height":437.142883,"width": 1422.857178}
		,
		"search": {
			"case": false,
			"lastGoodList": [
				"Lloyds",
				"Bank"
			]
		},
		"nullval" : null
	}`)
)

func TestEncodeStr(t *testing.T) {
	testEncode(t, "A\\t\\n\\nY", "A\t\n\nY")
	testEncode(t, "A\\uFFFDZ", "A\x81Z")
	testEncode(t, "A\\b\\nY", "A\b\n\rY")
	testEncode(t, "A\\f\\nY", "A\f\r\nY")
	testEncode(t, "A\\b\\nY", "A\b\nY")
	testEncode(t, "A\\nD<Z\\u2605Y", "A\nD\x3CZ\u2605Y")
	testEncode(t, "A\\tD\\fZ\\u2605\\u2605Y", "A\tD\fZ\u2605\u2605Y")
}

func testEncode(t *testing.T, expected, testStr string) {
	s := parser.EncodeQuotedString(testStr)
	if expected != s {
		t.Errorf("FAIL. Not Equal\nExpected:%s\nActual  :%s", bytesString(expected), bytesString(s))
	}
}

func bytesString(inStr string) string {
	var sb strings.Builder
	for _, c := range inStr {
		if c < 32 || c > 127 {
			sb.WriteString(fmt.Sprintf("(%d)", c))
		} else {
			sb.WriteString(fmt.Sprintf("%c", c))
		}
	}
	return sb.String()
}

func TestHexCharToIntPanic1(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("FAIL. Should panic")
		}
		if !strings.Contains(r.(string), "Invalid hex digit") {
			t.Errorf("FAIL. Should contain 'Invalid hex digit' actual: '%s'", r.(string))
		}
	}()
	parser.HexCharToInt(47) // '0' -1
}
func TestHexCharToIntPanic2(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("FAIL. Should panic")
		}
		if !strings.Contains(r.(string), "Invalid hex digit") {
			t.Errorf("FAIL. Should contain 'Invalid hex digit' actual: '%s'", r.(string))
		}
	}()
	parser.HexCharToInt('G')
}

func TestIntToHexChar4(t *testing.T) {
	failIfNotEqualRune2(t, "0000", parser.IntToHexChar4(0))
	failIfNotEqualRune2(t, "0040", parser.IntToHexChar4(64))
	failIfNotEqualRune2(t, "007F", parser.IntToHexChar4(127))
	failIfNotEqualRune2(t, "0080", parser.IntToHexChar4(128))
	failIfNotEqualRune2(t, "00FE", parser.IntToHexChar4(254))
	failIfNotEqualRune2(t, "00FF", parser.IntToHexChar4(255))
	failIfNotEqualRune2(t, "FFFE", parser.IntToHexChar4(65534))
	failIfNotEqualRune2(t, "FFFF", parser.IntToHexChar4(65535))
	failIfNotEqualRune2(t, "AF12", parser.IntToHexChar4(44818))
	failIfNotEqualRune2(t, "0000", parser.IntToHexChar4(65536))
}

func TestIntToHexChar2(t *testing.T) {
	failIfNotEqualRune2(t, "00", parser.IntToHexChar2(0))
	failIfNotEqualRune2(t, "40", parser.IntToHexChar2(64))
	failIfNotEqualRune2(t, "7F", parser.IntToHexChar2(127))
	failIfNotEqualRune2(t, "80", parser.IntToHexChar2(128))
	failIfNotEqualRune2(t, "FE", parser.IntToHexChar2(254))
	failIfNotEqualRune2(t, "FF", parser.IntToHexChar2(255))
	failIfNotEqualRune2(t, "FE", parser.IntToHexChar2(65534))
	failIfNotEqualRune2(t, "FF", parser.IntToHexChar2(65535))
	failIfNotEqualRune2(t, "12", parser.IntToHexChar2(44818))
	failIfNotEqualRune2(t, "00", parser.IntToHexChar2(65536))
}

func failIfNotEqualRune2(t *testing.T, exp string, act []rune) {
	for i, c := range exp {
		if c != act[i] {
			t.Errorf("FAIL. runes are not equal. Expected %s, actual %s'", exp, string(act))
		}
	}
}
func TestIntToHexChar(t *testing.T) {
	failIfNotEqualRune(t, '0', parser.IntToHexChar(0))
	failIfNotEqualRune(t, '1', parser.IntToHexChar(1))
	failIfNotEqualRune(t, '2', parser.IntToHexChar(2))
	failIfNotEqualRune(t, '3', parser.IntToHexChar(3))
	failIfNotEqualRune(t, '4', parser.IntToHexChar(4))
	failIfNotEqualRune(t, '5', parser.IntToHexChar(5))
	failIfNotEqualRune(t, '6', parser.IntToHexChar(6))
	failIfNotEqualRune(t, '7', parser.IntToHexChar(7))
	failIfNotEqualRune(t, '8', parser.IntToHexChar(8))
	failIfNotEqualRune(t, '9', parser.IntToHexChar(9))
	failIfNotEqualRune(t, 'A', parser.IntToHexChar(10))
	failIfNotEqualRune(t, 'B', parser.IntToHexChar(11))
	failIfNotEqualRune(t, 'C', parser.IntToHexChar(12))
	failIfNotEqualRune(t, 'D', parser.IntToHexChar(13))
	failIfNotEqualRune(t, 'E', parser.IntToHexChar(14))
	failIfNotEqualRune(t, 'F', parser.IntToHexChar(15))
	failIfNotEqualRune(t, '0', parser.IntToHexChar(16))
	failIfNotEqualRune(t, 'F', parser.IntToHexChar(255))
	failIfNotEqualRune(t, 'F', parser.IntToHexChar(65535))
}

func failIfNotEqualRune(t *testing.T, exp, act rune) {
	if exp != act {
		t.Errorf("FAIL. runes are not equal. Expected %d, actual %d'", exp, act)
	}
}
func TestHexCharToInt(t *testing.T) {
	failIfNotEqualInt(t, parser.HexCharToInt('0'), 0)
	failIfNotEqualInt(t, parser.HexCharToInt('1'), 1)
	failIfNotEqualInt(t, parser.HexCharToInt('2'), 2)
	failIfNotEqualInt(t, parser.HexCharToInt('3'), 3)
	failIfNotEqualInt(t, parser.HexCharToInt('4'), 4)
	failIfNotEqualInt(t, parser.HexCharToInt('5'), 5)
	failIfNotEqualInt(t, parser.HexCharToInt('6'), 6)
	failIfNotEqualInt(t, parser.HexCharToInt('7'), 7)
	failIfNotEqualInt(t, parser.HexCharToInt('8'), 8)
	failIfNotEqualInt(t, parser.HexCharToInt('9'), 9)
	failIfNotEqualInt(t, parser.HexCharToInt('A'), 10)
	failIfNotEqualInt(t, parser.HexCharToInt('B'), 11)
	failIfNotEqualInt(t, parser.HexCharToInt('C'), 12)
	failIfNotEqualInt(t, parser.HexCharToInt('D'), 13)
	failIfNotEqualInt(t, parser.HexCharToInt('E'), 14)
	failIfNotEqualInt(t, parser.HexCharToInt('F'), 15)
}

func failIfNotEqualInt(t *testing.T, i, j uint16) {
	if i != j {
		t.Errorf("FAIL. ints are not equal. Expected %d, actual %d'", i, j)
	}
}

func TestPanicUnquotedString(t *testing.T) {
	s := parser.NewScanner([]byte(`["literal, 1234, true, false]`))
	s.NextToken()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("FAIL. Should throw 'unterminated quoted String'")
		}
	}()
	s.NextToken()
}
func TestPanicUnexpectedEOF(t *testing.T) {
	s := parser.NewScanner([]byte(`[ `))
	s.NextToken()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("FAIL. Should throw 'Unexpected end of input'")
		}
	}()
	s.NextToken()
}
func TestPanicUnrecognisedToken(t *testing.T) {
	s := parser.NewScanner([]byte(`[ bad ]`))
	s.NextToken()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("FAIL. Should throw 'Unrecognised token'")
		}
	}()
	s.NextToken()
}

func TestTokens(t *testing.T) {
	s := parser.NewScanner(text)
	token(t, "", s, parser.TT_OBJECT_OPEN)
	token(t, "data", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "", s, parser.TT_OBJECT_OPEN)
	token(t, "positional", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "true", s, parser.TT_BOOL_TRUE)
	token(t, "", s, parser.TT_OBJECT_CLOSE)
	token(t, "", s, parser.TT_COMMA)
	token(t, "datafile", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "../data.json", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COMMA)
	token(t, "screen", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "", s, parser.TT_OBJECT_OPEN)
	token(t, "fullScreen", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "false", s, parser.TT_BOOL_FALSE)
	token(t, "", s, parser.TT_COMMA)
	token(t, "height", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "437.142883", s, parser.TT_NUMBER)
	token(t, "", s, parser.TT_COMMA)
	token(t, "width", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "1422.857178", s, parser.TT_NUMBER)
	token(t, "", s, parser.TT_OBJECT_CLOSE)
	token(t, "", s, parser.TT_COMMA)
	token(t, "search", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "", s, parser.TT_OBJECT_OPEN)
	token(t, "case", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "false", s, parser.TT_BOOL_FALSE)
	token(t, "", s, parser.TT_COMMA)
	token(t, "lastGoodList", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "", s, parser.TT_ARRAY_OPEN)
	token(t, "Lloyds", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COMMA)
	token(t, "Bank", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_ARRAY_CLOSE)
	token(t, "", s, parser.TT_OBJECT_CLOSE)
	token(t, "", s, parser.TT_COMMA)
	token(t, "nullval", s, parser.TT_QUOTED_STRING)
	token(t, "", s, parser.TT_COLON)
	token(t, "null", s, parser.TT_NULL)
	token(t, "", s, parser.TT_OBJECT_CLOSE)

}
func TestTokens2(t *testing.T) {
	s := parser.NewScanner([]byte(" {true false}[\"ABC\", \"123\"]:1234"))
	s.SkipSpace()
	toc := token(t, "", s, parser.TT_OBJECT_OPEN)
	assertTrue(t, "OBJECT_OPEN", toc.IsObjectOpen())
	toc = token(t, "true", s, parser.TT_BOOL_TRUE)
	assertTrue(t, "TRUE", toc.IsBool())
	toc = token(t, "false", s, parser.TT_BOOL_FALSE)
	assertTrue(t, "FALSE", toc.IsBool())
	toc = token(t, "", s, parser.TT_OBJECT_CLOSE)
	assertTrue(t, "OBJECT_CLOSE", toc.IsObjectClose())
	toc = token(t, "", s, parser.TT_ARRAY_OPEN)
	assertTrue(t, "ARRAY_OPEN", toc.IsArrayOpen())
	toc = token(t, "ABC", s, parser.TT_QUOTED_STRING)
	assertTrue(t, "QUOTED_STRING", toc.IsQuotedString())
	toc = token(t, "", s, parser.TT_COMMA)
	assertTrue(t, "COMMA", toc.IsComma())
	toc = token(t, "123", s, parser.TT_QUOTED_STRING)
	assertTrue(t, "QUOTED_STRING", toc.IsQuotedString())
	toc = token(t, "", s, parser.TT_ARRAY_CLOSE)
	assertTrue(t, "ARRAY_CLOSE", toc.IsArrayClose())
	toc = token(t, "", s, parser.TT_COLON)
	assertTrue(t, "COLON", toc.IsColon())
	toc = token(t, "1234", s, parser.TT_NUMBER)
	assertTrue(t, "NUMBER", toc.IsNumber())
}

func TestBasicScanner(t *testing.T) {
	s := parser.NewScanner([]byte(" HI"))
	s.SkipSpace()
	assertChar(t, "Incorrect char from Next after SkipSpace", s.Next(), 'H')
	assertChar(t, "Incorrect char from Next", s.Next(), 'I')
	assertChar(t, "Incorrect char from Next", s.Next(), 0)

	s.Reset().SkipToNext('I')
	assertChar(t, "Incorrect char from Next", s.Next(), 'I')
	s.Reset().SkipToNext('X')
	assertChar(t, "Incorrect char from Next", s.Next(), 0)

	s = parser.NewScanner([]byte(" Hello *1 "))
	s.SkipToNext('*')
	assertChar(t, "Incorrect char from Next", s.Next(), '*')
	s.SkipToNext('1')
	assertChar(t, "Incorrect char from Next", s.Next(), '1')

}

func token(t *testing.T, expected string, s *parser.Scanner, tt parser.TokenType) *parser.Token {
	defer func() {
		r := recover()
		if r != nil {
			t.Errorf("FAIL. Should not panic")
		}
	}()
	toc := s.NextToken()
	if !toc.IsType(tt) {
		t.Errorf("NextToken error: Token [%s]: is not of type %s ", toc, parser.GetTokenTypeName(tt))
	}
	if toc.GetStringValue() != expected {
		t.Errorf("NextToken error: Token [%s]: does not contain text '%s' ", toc, expected)
	}
	return toc
}

func assertTrue(t *testing.T, m string, value bool) {
	if !value {
		t.Errorf("AssertTrue: %s: actual: %t expected: true", m, value)
	}
}

func assertChar(t *testing.T, m string, c, b byte) {
	if c != b {
		t.Errorf("AssertChar: %s: actual: %s expected: %s", m, displayChar(c), displayChar(b))
	}
}

func displayChar(c byte) string {
	if c < 32 {
		return fmt.Sprintf("[%d]", c)
	}
	return fmt.Sprintf("'%s'", string(c))
}
