package config

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
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
func NewParser(filename string) (*Parser, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		// Ok if file does not exists but report read errors.
		if os.IsNotExist(err) == false {
			return nil, err
		}
		return nil, nil
	}
	p := &Parser{
		lexer: NewLexer(filename, string(contents)),
	}
	return p, nil
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
		if p.lexer.NextToken(); p.lexer.Token.Kind == TkError {
			return &ConfigurationError{
				Filename: p.lexer.Filename,
				Line:     p.lexer.Token.Line,
				Column:   p.lexer.Token.Column,
				msg:      p.lexer.Token.Value.(string),
			}
		}
		for p.lexer.Token.Kind != TkEOF {
			if err := p.parseConfig(c, ""); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Parser) skipEmptyLines() {
	for {
		if p.lexer.Token.Kind != TkEOL {
			return
		}
		p.lexer.NextToken()
	}
}

func (p *Parser) parseConfig(c *Configuration, section string) (err *ConfigurationError) {
	p.skipEmptyLines()
	if p.lexer.Token.Kind == TkLBracket {
		return p.parseSection(c, section)
	}
	if p.lexer.Token.Kind == TkIdentifier {
		return p.parseOptions(c, section)
	}
	return p.expectedError()
}

func (p *Parser) parseSection(c *Configuration, section string) (err *ConfigurationError) {
	if p.lexer.NextToken(); p.lexer.Token.Kind == TkIdentifier {
		// Remember section name for error reporting.
		currentSection := p.lexer.Token
		section = p.formatOptionName(section, p.lexer.Token.Value.(string))
		optionCountSave := p.optionCount
		if p.lexer.NextToken(); p.lexer.Token.Kind == TkRBracket {
			// Set error to end of section declaration.
			currentSection.Column = p.lexer.Token.Column
			if p.lexer.NextToken(); p.lexer.Token.Kind == TkEOL {
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
	for p.lexer.Token.Kind == TkIdentifier {
		option := p.lexer.Token.Value.(string)
		p.lexer.NextToken()
		if err = p.parseOption(c, section, option); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseOption(c *Configuration, section, option string) (err *ConfigurationError) {
	if p.lexer.Token.Kind == TkEqual {
		if p.lexer.NextToken(); p.lexer.Token.Kind == TkLBracket {
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
		if p.lexer.Token.Kind == TkEOL {
			p.skipEmptyLines()
			return nil
		}
		if p.lexer.Token.Kind == TkEOF {
			return nil
		}
	}
	return p.unexpectedError()
}

func (p *Parser) parseValue(c *Configuration) (err *ConfigurationError) {
	switch p.lexer.Token.Kind {
	case TkBool, TkInt, TkFloat, TkDate, TkString:
		return nil
	default:
		return p.unexpectedError()
	}
}

func (p *Parser) parseArray(c *Configuration, section, option string) (err *ConfigurationError) {
	skipEOL := func(p *Parser) {
		for p.lexer.Token.Kind == TkEOL {
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
		if p.lexer.Token.Kind == TkComma {
			p.lexer.NextToken()
		} else if p.lexer.Token.Kind == TkRBracket {
			option = p.formatOptionName(section, option)
			c.setOption(option, array)
			// One option was created or updated.
			p.optionCount++
			p.lexer.NextToken()
			return nil
		}
		skipEOL(p)
	}
}

func (p *Parser) convertValue(dstValue interface{}, srcValue *token) (err *ConfigurationError) {
	if reflect.TypeOf(dstValue) == dateType {
		if srcValue.Kind == TkDate {
			return nil
		}
		return p.convertValueError("time.Time", srcValue.Kind)
	}

	switch reflect.ValueOf(dstValue).Kind() {
	case reflect.Bool:
		if srcValue.Kind != TkBool {
			return p.convertValueError("bool", srcValue.Kind)
		}
	case reflect.Int64:
		if srcValue.Kind == TkInt {
			return nil
		}
		if srcValue.Kind == TkFloat {
			// Floating point number can be safely converted to an integer?
			f := srcValue.Value.(float64)
			if math.Floor(f) == f {
				srcValue.Kind = TkInt
				srcValue.Value = int64(f)
				return nil
			}
		}
		return p.convertValueError("int64", srcValue.Kind)
	case reflect.Float64:
		if srcValue.Kind == TkFloat {
			return nil
		}
		if srcValue.Kind == TkInt {
			i := srcValue.Value.(int64)
			srcValue.Kind = TkFloat
			srcValue.Value = float64(i)
			return nil
		}
		return p.convertValueError("float64", srcValue.Kind)
	case reflect.String:
		if srcValue.Kind == TkString {
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
	case TkEOF:
		kind = "end-of-file"
	case TkEOL:
		kind = "end-of-line"
	case TkIdentifier:
		kind = fmt.Sprintf("identifier %s", p.lexer.Token.Value)
	case TkBool:
		kind = fmt.Sprintf("boolean %t", p.lexer.Token.Value)
	case TkString:
		kind = fmt.Sprintf("string \"%s\"", p.lexer.Token.Value)
	case TkInt:
		kind = fmt.Sprintf("integer %d", p.lexer.Token.Value)
	case TkFloat:
		kind = fmt.Sprintf("floating point number %f", p.lexer.Token.Value)
	case TkDate:
		date := (p.lexer.Token.Value.(time.Time)).Format(time.RFC3339)
		kind = fmt.Sprintf("date %s", date)
	case TkEqual, TkLBracket, TkRBracket, TkComma:
		kind = fmt.Sprintf("character '%s'", p.lexer.Token.Value)
	default:
		panic(fmt.Sprintf("unexpected kind %s", kind))
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
