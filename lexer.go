// Copyright 2013 Christophe Bonello. All rights reserved.

package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type (
	kind int8

	char struct {
		offset int
		line   int
		column int
		r      rune
	}

	token struct {
		Kind   kind
		Line   int
		Column int
		Value  interface{}
	}

	lexer struct {
		Filename string
		contents string
		inputLen int   // Length of file contents.
		offset   int   // Start of last rune read from input.
		width    int   // Width of last rune read from input.
		line     int   // Line of last rune declaration.
		column   int   // Column of last rune declaration.
		c        char  // Last rune read.
		Token    token // Last token read.
	}
)

// Tokens returned by lexer.
const (
	TK_EOF        kind = iota // End-of-file token.
	TK_EOL                    // End-of-line token.
	TK_ERROR                  // An error occurred; value is error message.
	TK_IDENTIFIER             // Identifier token.
	TK_BOOL                   // Boolean.
	TK_STRING                 // A string (does not include double quotes).
	TK_INT                    // An integer.
	TK_FLOAT                  // A floating point number.
	TK_DATE                   // A date.
	TK_EQUAL                  // '='.
	TK_LBRACKET               // '['.
	TK_RBRACKET               // ']'.
	TK_COMMA                  // ','.

	_EOF = -1
)

var (
	// To parse a date in zulu form (RFC 3339 format).
	dateRe = regexp.MustCompile("^\\d{1,4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z")
)

// NewLexer instanciates a new lexer.
func NewLexer(filename, contents string) *lexer {
	l := &lexer{
		Filename: filename,
		contents: contents,
		inputLen: len(contents),
		offset:   0,
		width:    0,
		line:     1,
		column:   1,
	}
	l.nextRune()
	return l
}

// NextToken returns the next token.
func (l *lexer) NextToken() {
	for {
		l.skipWhitespaces()
		if l.c.r == '#' {
			l.skipComment()
		} else {
			switch {
			case unicode.IsLetter(l.c.r) || l.c.r == '_':
				l.parseIdentifier()
			case l.c.r == '"':
				l.parseString()
			case unicode.IsDigit(l.c.r):
				if l.parseDate() == false {
					l.parseNumber()
				}
			case l.c.r == '-', l.c.r == '+':
				l.parseNumber()
			case l.c.r == '=':
				defer l.nextRune()
				l.setToken(TK_EQUAL, l.c.line, l.c.column, "=")
			case l.c.r == '[':
				defer l.nextRune()
				l.setToken(TK_LBRACKET, l.c.line, l.c.column, "[")
			case l.c.r == ']':
				defer l.nextRune()
				l.setToken(TK_RBRACKET, l.c.line, l.c.column, "]")
			case l.c.r == ',':
				defer l.nextRune()
				l.setToken(TK_COMMA, l.c.line, l.c.column, ",")
			case l.c.r == '\n':
				defer l.nextRune()
				l.setToken(TK_EOL, l.c.line, l.c.column, nil)
			case l.c.r == _EOF:
				l.setToken(TK_EOF, l.c.line, l.c.column, nil)
			default:
				if unicode.IsPrint(l.c.r) {
					l.setErrorToken("unexpected '%c' character",
						l.c.r)
				} else {
					l.setErrorToken("unexpected \\u%04X character",
						l.c.r)
				}
			}
			break
		}
	}
}

func isSpace(r rune) bool {
	switch r {
	case '\t', '\v', '\f', ' ', 0x85, 0xA0:
		return true
	}
	return false
}

func isEOL(r rune) bool {
	return r == '\n' || r == '\r'
}

func isHexadecimal(r rune) bool {
	return unicode.IsNumber(r) ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}

func isAlphaNumeric(r rune) bool {
	alphaDigitSet := []*unicode.RangeTable{unicode.L, unicode.N}
	return unicode.IsOneOf(alphaDigitSet, r) || r == '_' || r == '-'
}

// Returns and consumes the next rune.
func (l *lexer) nextRune() {
	l.c.offset = l.offset
	l.c.line = l.line
	l.c.column = l.column
	if l.offset >= l.inputLen {
		l.c.r = _EOF
	} else {
		l.c.r, l.width = utf8.DecodeRuneInString(l.contents[l.offset:])
		if l.c.r == '\n' {
			l.line++
			l.column = 0
		}
		l.offset += l.width
		l.column += l.width
	}
}

// Consumes the next rune if it's from the valid set.
func (l *lexer) acceptOneRune(valid string) bool {
	if strings.IndexRune(valid, l.c.r) >= 0 {
		l.nextRune()
		return true
	}
	return false
}

// Consumes as many runes as we can from the valid set.
// The count determines the number of runes to consume:
//   n < 0: consume as many runes as we can.
//	 n = 0: do nothing!
//   n > 0: consumes exactly n runes.
func (l *lexer) acceptManyRunes(valid string, n int) (consumed int) {
	switch {
	case n < 0:
		for l.acceptOneRune(valid) {
			consumed++
		}
	case n > 0:
		for l.acceptOneRune(valid) {
			consumed++
			if n = n - 1; n <= 0 {
				break
			}
		}
	}
	return consumed
}

// Returns but does not consume the next rune.
func (l *lexer) peekRune() (r rune) {
	r, _ = utf8.DecodeRuneInString(l.contents[l.offset:])
	return r
}

// Skip N runes.
func (l *lexer) skipRune(n int) {
	for n > 0 {
		l.nextRune()
		if l.c.r == _EOF {
			break
		}
		n--
	}
}

func (l *lexer) setErrorToken(format string, args ...interface{}) {
	l.Token.Kind = TK_ERROR
	l.Token.Line = l.c.line
	l.Token.Column = l.c.column
	l.Token.Value = fmt.Sprintf(format, args...)
	//fmt.Printf("TOKEN %s\n", l.Token)
}

func (l *lexer) setToken(k kind, line, column int, value interface{}) {
	l.Token.Kind = k
	l.Token.Line = line
	l.Token.Column = column
	l.Token.Value = value
	//fmt.Printf("TOKEN %s\n", l.Token)
}

func (l *lexer) skipWhitespaces() {
	if isSpace(l.c.r) {
		for {
			if l.nextRune(); isSpace(l.c.r) == false {
				break
			}
		}
	}
}

func (l *lexer) skipComment() {
	for {
		if l.nextRune(); isEOL(l.c.r) == true {
			break
		}
	}
}

func (l *lexer) parseIdentifier() {
	start := l.c
	for {
		if l.nextRune(); isAlphaNumeric(l.c.r) == false {
			break
		}
	}
	id := l.contents[start.offset:l.c.offset]
	switch id {
	case "true":
		l.setToken(TK_BOOL, start.line, start.column, true)
	case "false":
		l.setToken(TK_BOOL, start.line, start.column, false)
	default:
		l.setToken(TK_IDENTIFIER, start.line, start.column, id)
	}
}

func (l *lexer) parseDate() bool {
	start := l.c

	// Regexp: easiest way to parse a date.
	if d := dateRe.FindString(l.contents[start.offset:]); d != "" {
		date, err := time.Parse(time.RFC3339, d)
		if err != nil {
			l.setErrorToken(err.Error())
			return true
		}
		l.skipRune(len(d) - 1)
		l.setToken(TK_DATE, start.line, start.column, date)
		l.nextRune()
		return true
	}
	return false
}

func (l *lexer) parseNumber() {
	k := TK_INT
	start := l.c
	sign := ""
	digits := "0123456789"
	base := 10

	// STEP #1: parse a number and only handle the obvious syntax errors. Any
	// problem will be catched up later by strconv anyway.
	if l.c.r == '-' {
		sign = "-"
		l.nextRune()
	} else if l.c.r == '+' {
		sign = "+"
		l.nextRune()
	}

	x := l.peekRune()
	if l.c.r == '0' && (x == 'x' || x == 'X') {
		l.nextRune() // Skips 'x' or 'X'.
		l.nextRune()
		digits = "0123456789abcdefABCDEF"
		base = 16
		if l.acceptManyRunes(digits, -1) == 0 {
			l.setErrorToken("malformed hex constant %q",
				l.contents[start.offset:l.c.offset])
			return
		}
	} else {
		digitsBeforeDot := true
		if l.acceptManyRunes(digits, -1) == 0 {
			// .5 is a valid floating-point constant; we cannot flagged an
			// error yet ('.' may be followed by digits).
			if l.c.r != '.' {
				l.setErrorToken("malformed constant %q",
					l.contents[start.offset:l.c.offset])
				return
			}
			digitsBeforeDot = false
		}
		if l.acceptOneRune(".") {
			k = TK_FLOAT
			if l.acceptManyRunes(digits, -1) == 0 {
				// 5. is a valid floating point constant. However, it's an
				// error if there was no digits before the '.'.
				if digitsBeforeDot == false {
					l.setErrorToken("malformed floating-point constant %q",
						l.contents[start.offset:l.c.offset])
					return
				}
			}
		}
		if l.acceptOneRune("eE") {
			l.acceptOneRune("+-")
			if l.acceptManyRunes("0123456789", -1) == 0 {
				l.c = start
				l.setErrorToken("malformed floating-point constant exponent")
				return
			}
		}
	}
	s := l.contents[start.offset:l.c.offset]
	// STEP #2: use strconv to convert from string to int or float.
	if k == TK_INT {
		if base == 16 {
			// -0xFF is flagged as an error by ParseInt(). Number should be
			// written as -FF.
			if sign != "" {
				s = s[1:]
			}
			// Skip the 0x prefix.
			s = sign + s[2:]
		}
		if ival, err := strconv.ParseInt(s, base, 64); err == nil {
			l.setToken(k, start.line, start.column, ival)
		} else {
			msg := err.Error()
			msg = strings.TrimPrefix(msg, "strconv.ParseInt: parsing ")
			l.setErrorToken(msg)
		}
	} else {
		if fval, err := strconv.ParseFloat(s, 64); err == nil {
			l.setToken(k, start.line, start.column, fval)
		} else {
			msg := err.Error()
			msg = strings.TrimPrefix(msg, "strconv.ParseFloat: parsing ")
			l.setErrorToken(msg)
		}
	}
}

func (l *lexer) parseString() {
	start := l.c

	// Skips opening double quotes.
	l.nextRune()
	for {
		if l.isEOLOrEOF() {
			return
		}
		if l.c.r == '\\' {
			// Saves position of '\'; if an error occurs this is where we
			// should report the error.
			es := l.c
			l.nextRune()
			switch l.c.r {
			case 'b', 't', 'n', 'f', 'r', '"', '/', '\\':
			case 'u':
				for i := 0; i < 4; i++ {
					l.nextRune()
					if l.isEOLOrEOF() {
						return
					}
					if !isHexadecimal(l.c.r) {
						if l.c.r == '"' {
							// End of string.
							l.c.column = es.column
							l.setErrorToken("malformed hex escape sequence %s",
								l.contents[es.offset:l.c.offset])
							return
						}
						l.setErrorToken("non-hex character in escape sequence: %q",
							l.c.r)
						return
					}
				}
			default:
				l.setErrorToken("unknown escape sequence: %s",
					l.contents[es.offset:l.c.offset+utf8.RuneLen(l.c.r)])
				return
			}
		} else if l.c.r == '"' {
			break
		}
		l.nextRune()
	}
	// Skips closing double quotes.
	l.nextRune()

	// Removes double quotes from string.
	qlen := utf8.RuneLen('"')
	s := l.contents[start.offset+qlen : l.c.offset-qlen]
	l.setToken(TK_STRING, start.line, start.column, s)
}

func (l *lexer) isEOLOrEOF() bool {
	if l.c.r == _EOF {
		l.setErrorToken("end-of-file in string")
		return true
	}
	if isEOL(l.c.r) {
		l.setErrorToken("newline in string")
		return true
	}
	return false
}

func (k kind) String() string {
	switch k {
	case TK_EOF:
		return "eof"
	case TK_EOL:
		return "eol"
	case TK_ERROR:
		return "error"
	case TK_IDENTIFIER:
		return "identifier"
	case TK_BOOL:
		return "bool"
	case TK_STRING:
		return "string"
	case TK_INT:
		return "int64"
	case TK_FLOAT:
		return "float64"
	case TK_DATE:
		return "time.Time"
	case TK_EQUAL:
		return "equal"
	case TK_LBRACKET:
		return "lbracket"
	case TK_RBRACKET:
		return "rbracket"
	default:
		return "comma"
	}
}

func (c char) String() string {
	pos := fmt.Sprintf("[%d - %3d:%3d]", c.offset, c.line, c.column)
	return fmt.Sprintf("rune %s ('%c')", pos, c.r)
}

func (t token) String() string {
	pos := fmt.Sprintf("[%3d:%3d]", t.Line, t.Column)
	// %-10s: IDENTIFIER is the biggest type name with 10 letters.
	switch {
	case t.Kind == TK_EOF:
		return fmt.Sprintf("%-10s %s", t.Kind, pos)
	case t.Kind == TK_EOL:
		return fmt.Sprintf("%-10s %s", t.Kind, pos)
	case t.Kind == TK_ERROR:
		return fmt.Sprintf("%-10s %s \"%s\"", t.Kind, pos, t.Value.(string))
	case t.Kind == TK_IDENTIFIER:
		return fmt.Sprintf("%-10s %s \"%s\"", t.Kind, pos, t.Value.(string))
	case t.Kind == TK_BOOL:
		return fmt.Sprintf("%-10s %s %t", t.Kind, pos, t.Value.(bool))
	case t.Kind == TK_STRING:
		return fmt.Sprintf("%-10s %s \"%s\"", t.Kind, pos, t.Value.(string))
	case t.Kind == TK_INT:
		return fmt.Sprintf("%-10s %s %d", t.Kind, pos, t.Value.(int64))
	case t.Kind == TK_FLOAT:
		return fmt.Sprintf("%-10s %s %f)", t.Kind, pos, t.Value.(float64))
	case t.Kind == TK_DATE:
		date := (t.Value.(time.Time)).Format(time.RFC3339)
		return fmt.Sprintf("%-10s %s %s", t.Kind, pos, date)
	case t.Kind == TK_EQUAL:
		return fmt.Sprintf("%-10s %s '='", t.Kind, pos)
	case t.Kind == TK_LBRACKET:
		return fmt.Sprintf("%-10s %s '['", t.Kind, pos)
	case t.Kind == TK_RBRACKET:
		return fmt.Sprintf("%-10s %s ']'", t.Kind, pos)
	default:
		return fmt.Sprintf("%-10s %s ','", t.Kind, pos)
	}
}
