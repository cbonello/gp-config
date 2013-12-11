// Copyright 2013 Christophe Bonello. All rights reserved.

package config

import (
	"fmt"
	"io/ioutil"
	"math"
	"reflect"
	"time"
)

type (
	// Parser context.
	Parser struct {
		lexer *lexer
		// Number of options created or updated in current section (to detect
		// empty sections).
		optionCount uint
	}
)

// NewParser instanciates a parser for given configuration file.
func NewParser(filename string) *Parser {
	if contents, err := ioutil.ReadFile(filename); err != nil {
		// Ok if file does not exists.
		return nil
	} else {
		p := &Parser{
			lexer: NewLexer(filename, string(contents)),
		}
		return p
	}
}

// NewStringParser instanciates a parser for given configuration string.
func NewStringParser(contents string) *Parser {
	p := &Parser{
		lexer: NewLexer(":string:", contents),
	}
	return p
}

// Parse parses a configuration stored either in a file or a string.
func (p *Parser) Parse(c *Configuration) (err *ConfigurationError) {
	if p != nil {
		if p.lexer.NextToken(); p.lexer.Token.Kind == TK_ERROR {
			return &ConfigurationError{
				Filename: p.lexer.Filename,
				Line:     p.lexer.Token.Line,
				Column:   p.lexer.Token.Column,
				msg:      p.lexer.Token.Value.(string),
			}
		} else {
			for p.lexer.Token.Kind != TK_EOF {
				if err := p.parseConfig(c, ""); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *Parser) skipEmptyLines() {
	for {
		if p.lexer.Token.Kind != TK_EOL {
			return
		}
		p.lexer.NextToken()
	}
}

func (p *Parser) parseConfig(c *Configuration, section string) (err *ConfigurationError) {
	p.skipEmptyLines()
	if p.lexer.Token.Kind == TK_LBRACKET {
		return p.parseSection(c, section)
	}
	if p.lexer.Token.Kind == TK_IDENTIFIER {
		return p.parseOptions(c, section)
	}
	return p.expectedError()
}

func (p *Parser) parseSection(c *Configuration, section string) (err *ConfigurationError) {
	if p.lexer.NextToken(); p.lexer.Token.Kind == TK_IDENTIFIER {
		// Remember section name for error reporting.
		currentSection := p.lexer.Token
		section = p.formatOptionName(section, p.lexer.Token.Value.(string))
		optionCountSave := p.optionCount
		if p.lexer.NextToken(); p.lexer.Token.Kind == TK_RBRACKET {
			// Set error to end of section declaration.
			currentSection.Column = p.lexer.Token.Column
			if p.lexer.NextToken(); p.lexer.Token.Kind == TK_EOL {
				p.lexer.NextToken()
				err := p.parseConfig(c, section)
				// No options were declared in section?
				if p.optionCount == optionCountSave {
					return p.emptySectionError(p.lexer.Filename,
						&currentSection)
				}
				return err
			}
		}
	}
	return p.unexpectedError()
}

func (p *Parser) parseOptions(c *Configuration, section string) (err *ConfigurationError) {
	for p.lexer.Token.Kind == TK_IDENTIFIER {
		option := p.lexer.Token.Value.(string)
		p.lexer.NextToken()
		if err = p.parseOption(c, section, option); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseOption(c *Configuration, section, option string) (err *ConfigurationError) {
	if p.lexer.Token.Kind == TK_EQUAL {
		if p.lexer.NextToken(); p.lexer.Token.Kind == TK_LBRACKET {
			p.lexer.NextToken()
			err = p.parseArray(c, section, option)
		} else {
			err = p.parseValue(c)
			option = p.formatOptionName(section, option)
			c.setOption(option, p.lexer.Token.Value)
			// One option was created or updated.
			p.optionCount++
			p.lexer.NextToken()
		}
		if err != nil {
			return err
		}
		if p.lexer.Token.Kind == TK_EOL {
			p.skipEmptyLines()
			return nil
		}
		if p.lexer.Token.Kind == TK_EOF {
			return nil
		}
	}
	return p.unexpectedError()
}

func (p *Parser) parseValue(c *Configuration) (err *ConfigurationError) {
	switch p.lexer.Token.Kind {
	case TK_BOOL, TK_INT, TK_FLOAT, TK_DATE, TK_STRING:
		return nil
	default:
		return p.unexpectedError()
	}
}

func (p *Parser) parseArray(c *Configuration, section, option string) (err *ConfigurationError) {
	skipEOL := func(p *Parser) {
		for p.lexer.Token.Kind == TK_EOL {
			p.lexer.NextToken()
		}
	}

	skipEOL(p)
	firstValue := p.lexer.Token
	array := []interface{}{}
	for {
		currentValue := p.lexer.Token
		if err = p.parseValue(c); err != nil {
			return err
		}
		if firstValue.Kind != currentValue.Kind {
			// Array elements must share the same underlying type.
			if err := p.convertValue(firstValue.Value, &currentValue); err != nil {
				return err
			}
		}
		array = append(array, currentValue.Value)
		p.lexer.NextToken()
		skipEOL(p)
		//p.lexer.NextToken()
		if p.lexer.Token.Kind == TK_COMMA {
			p.lexer.NextToken()
		} else if p.lexer.Token.Kind == TK_RBRACKET {
			option = p.formatOptionName(section, option)
			c.setOption(option, array)
			// One option was created or updated.
			p.optionCount++
			p.lexer.NextToken()
			return nil
		}
		skipEOL(p)
	}
	return nil
}

func (p *Parser) convertValue(dstValue interface{}, srcValue *token) (err *ConfigurationError) {
	if reflect.TypeOf(dstValue) == dateType {
		if srcValue.Kind == TK_DATE {
			return nil
		}
		return p.convertValueError("time.Time", srcValue.Kind)
	}

	switch reflect.ValueOf(dstValue).Kind() {
	case reflect.Bool:
		if srcValue.Kind != TK_BOOL {
			return p.convertValueError("bool", srcValue.Kind)
		}
	case reflect.Int64:
		if srcValue.Kind == TK_INT {
			return nil
		}
		if srcValue.Kind == TK_FLOAT {
			// Floating point number can be safely converted to an integer?
			f := srcValue.Value.(float64)
			if math.Floor(f) == f {
				srcValue.Kind = TK_INT
				srcValue.Value = int64(f)
				return nil
			}
		}
		return p.convertValueError("int64", srcValue.Kind)
	case reflect.Float64:
		if srcValue.Kind == TK_FLOAT {
			return nil
		}
		if srcValue.Kind == TK_INT {
			i := srcValue.Value.(int64)
			srcValue.Kind = TK_FLOAT
			srcValue.Value = float64(i)
			return nil
		}
		return p.convertValueError("float64", srcValue.Kind)
	case reflect.String:
		if srcValue.Kind == TK_STRING {
			return nil
		}
		return p.convertValueError("string", srcValue.Kind)
	default:
		panic(fmt.Sprintf("unexpected type '%s'",
			reflect.TypeOf(dstValue).Kind()))
	}
	return nil
}

func (p *Parser) formatOptionName(section, option string) string {
	if section == "" {
		return option
	}
	return section + "." + option
}

func (p *Parser) emptySectionError(filename string, t *token) (err *ConfigurationError) {
	return &ConfigurationError{
		Filename: filename,
		Line:     t.Line,
		Column:   t.Column,
		msg:      fmt.Sprintf("empty section %s", t.Value),
	}
}

func (p *Parser) expectedError() (err *ConfigurationError) {
	return &ConfigurationError{
		Filename: p.lexer.Filename,
		Line:     p.lexer.Token.Line,
		Column:   p.lexer.Token.Column,
		msg:      "expected section or option declaration",
	}
}

func (p *Parser) unexpectedError() (err *ConfigurationError) {
	var kind string

	switch p.lexer.Token.Kind {
	case TK_EOF:
		kind = "end-of-file"
	case TK_EOL:
		kind = "end-of-line"
	case TK_IDENTIFIER:
		kind = fmt.Sprintf("identifier %s", p.lexer.Token.Value)
	case TK_BOOL:
		kind = fmt.Sprintf("boolean %t", p.lexer.Token.Value)
	case TK_STRING:
		kind = fmt.Sprintf("string \"%s\"", p.lexer.Token.Value)
	case TK_INT:
		kind = fmt.Sprintf("integer %d", p.lexer.Token.Value)
	case TK_FLOAT:
		kind = fmt.Sprintf("floating point number %f", p.lexer.Token.Value)
	case TK_DATE:
		date := (p.lexer.Token.Value.(time.Time)).Format(time.RFC3339)
		kind = fmt.Sprintf("date %s", date)
	case TK_EQUAL, TK_LBRACKET, TK_RBRACKET, TK_COMMA:
		kind = fmt.Sprintf("character '%s'", p.lexer.Token.Value)
	default:
		panic(fmt.Sprintf("unexpected kind %d", kind))
	}

	return &ConfigurationError{
		Filename: p.lexer.Filename,
		Line:     p.lexer.Token.Line,
		Column:   p.lexer.Token.Column,
		msg:      fmt.Sprintf("unexpected %s", kind),
	}
}

func (p *Parser) convertValueError(srcKind string, dstKind kind) (err *ConfigurationError) {
	return &ConfigurationError{
		Filename: p.lexer.Filename,
		Line:     p.lexer.Token.Line,
		Column:   p.lexer.Token.Column,
		msg:      fmt.Sprintf("cannot use type %s as type %s", dstKind, srcKind),
	}
}
