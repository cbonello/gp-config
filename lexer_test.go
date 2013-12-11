package config_test

import (
	"fmt"
	"github.com/cbonello/gp-config"
	. "launchpad.net/gocheck"
)

type (
	LexerTests struct{}
)

var (
	_ = Suite(&LexerTests{})
)

func (lt *LexerTests) TestPass1(c *C) {
	contents := `
iden_ti-fier true false "abcd" "\u123456" 1234 +1 -210
0xAb -0xFFee 5. 1.2 -2.3456 1.5E5 -1.4e-4
# Comment
1979-05-27T07:32:00Z = [ ] ,`

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_EOL)
	c.Check(l.Token.Line, Equals, 1)
	c.Check(l.Token.Column, Equals, 1)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_IDENTIFIER)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 1)
	c.Check(l.Token.Value, Equals, "iden_ti-fier")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_BOOL)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 14)
	c.Check(l.Token.Value, Equals, true)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_BOOL)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 19)
	c.Check(l.Token.Value, Equals, false)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_STRING)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 25)
	c.Check(l.Token.Value, Equals, "abcd")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_STRING)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 32)
	c.Check(l.Token.Value, Equals, "\\u123456")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_INT)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 43)
	c.Check(l.Token.Value, Equals, int64(1234))
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_INT)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 48)
	c.Check(l.Token.Value, Equals, int64(1))
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_INT)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 51)
	c.Check(l.Token.Value, Equals, int64(-210))
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_EOL)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 55)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_INT)
	c.Check(l.Token.Line, Equals, 3)
	c.Check(l.Token.Column, Equals, 1)
	c.Check(l.Token.Value, Equals, int64(171))
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_INT)
	c.Check(l.Token.Line, Equals, 3)
	c.Check(l.Token.Column, Equals, 6)
	c.Check(l.Token.Value, Equals, int64(-65518))
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_FLOAT)
	c.Check(l.Token.Line, Equals, 3)
	c.Check(l.Token.Column, Equals, 14)
	c.Check(l.Token.Value, Equals, float64(5.0))
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_FLOAT)
	c.Check(l.Token.Line, Equals, 3)
	c.Check(l.Token.Column, Equals, 17)
	c.Check(l.Token.Value, Equals, float64(1.2))
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_FLOAT)
	c.Check(l.Token.Line, Equals, 3)
	c.Check(l.Token.Column, Equals, 21)
	c.Check(l.Token.Value, Equals, float64(-2.3456))
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_FLOAT)
	c.Check(l.Token.Line, Equals, 3)
	c.Check(l.Token.Column, Equals, 29)
	c.Check(l.Token.Value, Equals, float64(1.5E5))
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_FLOAT)
	c.Check(l.Token.Line, Equals, 3)
	c.Check(l.Token.Column, Equals, 35)
	c.Check(l.Token.Value, Equals, float64(-1.4e-4))
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_EOL)
	c.Check(l.Token.Line, Equals, 3)
	c.Check(l.Token.Column, Equals, 42)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_EOL)
	c.Check(l.Token.Line, Equals, 4)
	c.Check(l.Token.Column, Equals, 10)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_DATE)
	c.Check(l.Token.Line, Equals, 5)
	c.Check(l.Token.Column, Equals, 1)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_EQUAL)
	c.Check(l.Token.Line, Equals, 5)
	c.Check(l.Token.Column, Equals, 22)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_LBRACKET)
	c.Check(l.Token.Line, Equals, 5)
	c.Check(l.Token.Column, Equals, 24)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_RBRACKET)
	c.Check(l.Token.Line, Equals, 5)
	c.Check(l.Token.Column, Equals, 26)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_COMMA)
	c.Check(l.Token.Line, Equals, 5)
	c.Check(l.Token.Column, Equals, 28)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_EOF)
	c.Check(l.Token.Line, Equals, 5)
	c.Check(l.Token.Column, Equals, 29)
}

// Invalid symbol.
func (lt *LexerTests) TestSymbol1(c *C) {
	contents := "("

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "unexpected '(' character")
}

// Non-print character.
func (lt *LexerTests) TestSymbol2(c *C) {
	contents := fmt.Sprintf("%c", 0x007)

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "unexpected \\u0007 character")
}

// Malformed date.
func (lt *LexerTests) TestDate1(c *C) {
	// 1975/05/27 99:32:00; hour is out of range.
	contents := "1979-05-27T99:32:00Z"

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals,
		"parsing time \"1979-05-27T99:32:00Z\": hour out of range")
}

// Malformed constant.
func (lt *LexerTests) TestNumber1(c *C) {
	contents := "+"

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "malformed constant \"+\"")
}

// Malformed hexadecimal constant.
func (lt *LexerTests) TestNumber2(c *C) {
	contents := "0Xz"

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "malformed hex constant \"0X\"")
}

// Malformed floating-point number.
func (lt *LexerTests) TestNumber3(c *C) {
	contents := "123456789e123456789"

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "\"123456789e123456789\": invalid syntax")
}

// Malformed floating-point contant.
func (lt *LexerTests) TestNumber4(c *C) {
	contents := "+."

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "malformed floating-point constant \"+.\"")
}

// Malformed floating-point contant.
func (lt *LexerTests) TestNumber5(c *C) {
	contents := "+.1e"

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "malformed floating-point constant exponent")
}

// Malformed floating-point contant.
func (lt *LexerTests) TestNumber6(c *C) {
	contents := "+9.1e+"

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "malformed floating-point constant exponent")
}

// Out-of-range floating-point contant.
func (lt *LexerTests) TestNumber7(c *C) {
	contents := "0.123456789e123456789"

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "\"0.123456789e123456789\": value out of range")
}

// Syntax error.
func (lt *LexerTests) TestString1(c *C) {
	contents := `"abcd`

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "end-of-file in string")
}

// Syntax error.
func (lt *LexerTests) TestString2(c *C) {
	contents := `"a
`

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "newline in string")
}

// Malformed escape sequence.
func (lt *LexerTests) TestString3(c *C) {
	contents := `"a\ `

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "unknown escape sequence: \\ ")
}

// Malformed escape sequence.
func (lt *LexerTests) TestString4(c *C) {
	contents := `"a\u8`

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "end-of-file in string")
}

// Malformed escape sequence.
func (lt *LexerTests) TestString5(c *C) {
	contents := `"a\u8A
`

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "newline in string")
}

// Malformed escape sequence.
func (lt *LexerTests) TestString6(c *C) {
	contents := `"a\u0ab"`

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "malformed hex escape sequence \\u0ab")
}

// Malformed escape sequence.
func (lt *LexerTests) TestString7(c *C) {
	contents := `"a\u0aby"`

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_ERROR)
	c.Check(l.Token.Value, Equals, "non-hex character in escape sequence: 'y'")
}

// String(): Dump function.
func (lt *LexerTests) TestGetDump1(c *C) {
	contents := `
a true 1 2.3 2013-10-25T16:22:00Z "foo" = [ ] ,`

	l := config.NewLexer("dummy.conf", contents)
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_EOL)
	c.Check(l.Token.Line, Equals, 1)
	c.Check(l.Token.Column, Equals, 1)
	str := fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "eol        [  1:  1]")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_IDENTIFIER)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 1)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "identifier [  2:  1] \"a\"")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_BOOL)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 3)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "bool       [  2:  3] true")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_INT)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 8)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "int64      [  2:  8] 1")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_FLOAT)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 10)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "float64    [  2: 10] 2.300000)")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_DATE)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 14)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "time.Time  [  2: 14] 2013-10-25T16:22:00Z")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_STRING)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 35)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "string     [  2: 35] \"foo\"")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_EQUAL)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 41)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "equal      [  2: 41] '='")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_LBRACKET)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 43)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "lbracket   [  2: 43] '['")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_RBRACKET)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 45)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "rbracket   [  2: 45] ']'")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_COMMA)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 47)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "comma      [  2: 47] ','")
	l.NextToken()
	c.Check(l.Token.Kind, Equals, config.TK_EOF)
	c.Check(l.Token.Line, Equals, 2)
	c.Check(l.Token.Column, Equals, 48)
	str = fmt.Sprintf("%s", l.Token)
	c.Check(str, Equals, "eof        [  2: 48]")
}
