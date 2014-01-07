package config_test

import (
	"fmt"
	"github.com/cbonello/gp-config"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"time"
)

type (
	ParserTests struct {
		pathname string
		config   *config.Configuration
		parser   *config.Parser
	}
)

var (
	_ = Suite(&ParserTests{})
)

func (pt *ParserTests) createTmpFile(c *C, contents string) string {
	f, err := ioutil.TempFile("", "ParserTest")
	c.Check(err, IsNil)
	defer f.Close()

	pt.pathname = f.Name()
	n, err := f.Write([]byte(contents))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, len(contents))
	return f.Name()
}

func (pt *ParserTests) deleteTmpFile(c *C, pathname string) {
	err := os.Remove(pathname)
	c.Assert(err, IsNil)
}

func (pt *ParserTests) createTestEnv(c *C, contents string) {
	var err error
	fn := pt.createTmpFile(c, contents)
	pt.config = config.NewConfiguration()
	pt.parser, err = config.NewParser(fn)
	c.Assert(err, IsNil)
}

func (pt *ParserTests) cleanTestEnv(c *C) {
	pt.parser = nil
	pt.config = nil
	pt.deleteTmpFile(c, pt.pathname)
	pt.pathname = ""
}

// NewParser(): non-exisiting input file.
func (pt *ParserTests) TestNewParser1(c *C) {
	// Hopefully a non-existing file.
	parser, err := config.NewParser("rumpelstilzchen")

	c.Check(parser, IsNil)
	c.Check(err, IsNil)
}

// NewParser(): non-exisiting input file.
func (pt *ParserTests) TestNewParser2(c *C) {
	var err error
	fn := pt.createTmpFile(c, "foo")
	// Remove read permission to ensure that NewParser() will fail.
	os.Chmod(fn, 0211)
	pt.config = config.NewConfiguration()
	pt.parser, err = config.NewParser(fn)
	defer pt.cleanTestEnv(c)

	c.Assert(err, NotNil)
}

// Parse(): empty input file.
func (pt *ParserTests) TestPass1(c *C) {
	contents := ``

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 0)
}

// Parse(): global boolean option.
func (pt *ParserTests) TestPass2(c *C) {
	contents := `foo = true`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	foo, err := pt.config.GetBool("foo")
	c.Check(err, IsNil)
	c.Check(foo, Equals, true)
}

// Parse(): global integer option.
func (pt *ParserTests) TestPass3(c *C) {
	contents := `foo = 12`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	foo, err := pt.config.GetInt("foo")
	c.Check(err, IsNil)
	c.Check(foo, Equals, int64(12))
}

// Parse(): global floating-point option.
func (pt *ParserTests) TestPass4(c *C) {
	contents := `foo = +1.2`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	foo, err := pt.config.GetFloat("foo")
	c.Check(err, IsNil)
	c.Check(foo, Equals, 1.2)
}

// Parse(): global string option.
func (pt *ParserTests) TestPass5(c *C) {
	contents := `foo = "bar"`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	foo, err := pt.config.GetString("foo")
	c.Check(err, IsNil)
	c.Check(foo, Equals, "bar")
}

// Parse(): global date option.
func (pt *ParserTests) TestPass6(c *C) {
	contents := `foo = 1979-05-27T07:32:00Z`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	expected, _ := time.Parse(time.RFC3339, "1979-05-27T07:32:00Z")
	foo, err := pt.config.GetDate("foo")
	c.Check(err, IsNil)
	c.Check(foo, Equals, expected)
}

// Parse(): global boolean-array option.
func (pt *ParserTests) TestPass7(c *C) {
	contents := `foo = [true, false]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	foo, err := pt.config.GetBoolArray("foo")
	c.Check(err, IsNil)
	c.Check(foo, EqualSlice, []bool{true, false})
}

// Parse(): global integer-array option.
func (pt *ParserTests) TestPass8(c *C) {
	contents := `foo = [1, 2, 3, 4]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	foo, err := pt.config.GetIntArray("foo")
	c.Check(err, IsNil)
	c.Check(foo, EqualSlice, []int64{1, 2, 3, 4})
}

// Parse(): global floating-point array option.
func (pt *ParserTests) TestPass9(c *C) {
	contents := `foo = [1.2, -3.4]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	foo, err := pt.config.GetFloatArray("foo")
	c.Check(err, IsNil)
	c.Check(foo, EqualSlice, []float64{1.2, -3.4})
}

// Parse(): global date-array option.
func (pt *ParserTests) TestPass10(c *C) {
	contents := `foo = [1979-05-27T07:32:00Z, 1980-05-27T07:32:00Z]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	expected1, _ := time.Parse(time.RFC3339, "1979-05-27T07:32:00Z")
	expected2, _ := time.Parse(time.RFC3339, "1980-05-27T07:32:00Z")
	array := []time.Time{expected1, expected2}
	foo, err := pt.config.GetDateArray("foo")
	c.Check(err, IsNil)
	c.Check(foo, EqualSlice, array)
}

// Parse(): global string-array option.
func (pt *ParserTests) TestPass11(c *C) {
	contents := `foo = ["bar", "rab"]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	foo, err := pt.config.GetStringArray("foo")
	c.Check(err, IsNil)
	c.Check(foo, EqualSlice, []string{"bar", "rab"})
}

// Parse(): local boolean option.
func (pt *ParserTests) TestPass12(c *C) {
	contents := `
[test]
foo = true`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "test.foo")
	test_foo, err := pt.config.GetBool("test.foo")
	c.Check(err, IsNil)
	c.Check(test_foo, Equals, true)
}

// Parse(): local options.
func (pt *ParserTests) TestPass13(c *C) {
	contents := `
[test]
foo1 = true
foo1 = 12

foo2 = false`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 2)
	c.Check(pt.config, HasKey, "test.foo1")
	test_foo1, err := pt.config.GetInt("test.foo1")
	c.Check(err, IsNil)
	c.Check(test_foo1, Equals, int64(12))
	c.Check(pt.config, HasKey, "test.foo2")
	test_foo2, err := pt.config.GetBool("test.foo2")
	c.Check(err, IsNil)
	c.Check(test_foo2, Equals, false)
}

// Parse(): implicit float to int conmversion.
func (pt *ParserTests) TestPass14(c *C) {
	contents := `foo = [1, 2.0]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	foo, err := pt.config.GetIntArray("foo")
	c.Check(err, IsNil)
	c.Check(foo, EqualSlice, []int64{1, 2})
}

// Parse(): implicit int to float conmversion.
func (pt *ParserTests) TestPass15(c *C) {
	contents := `foo = [2.0, -3]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo")
	foo, err := pt.config.GetFloatArray("foo")
	c.Check(err, IsNil)
	c.Check(foo, EqualSlice, []float64{2.0, -3.0})
}

// Parse(): redeclaration of options.
func (pt *ParserTests) TestPass16(c *C) {
	contents := `
[foo]
a = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
[bar]
b = 3
[foo]
a = "new type and value"`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 2)
	c.Check(pt.config, HasKey, "foo.a")
	foo_a, err := pt.config.GetString("foo.a")
	c.Check(err, IsNil)
	c.Check(foo_a, Equals, "new type and value")
	c.Check(pt.config, HasKey, "bar.b")
	bar_b, err := pt.config.GetInt("bar.b")
	c.Check(err, IsNil)
	c.Check(bar_b, Equals, int64(3))
}

// Parse(): multiline arrays.
func (pt *ParserTests) TestPass17(c *C) {
	contents := `
[foo]
a = [
	0, 1,
	2, 3,
	4, 5,
	6, 7,
	8, 9
]
`
	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	c.Check(pt.parser.Parse(pt.config), IsNil)
	c.Check(pt.config.Len(), Equals, 1)
	c.Check(pt.config, HasKey, "foo.a")
	value, err := pt.config.GetIntArray("foo.a")
	c.Check(err, IsNil)
	c.Check(value, EqualSlice, []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

// Parse(): parser error.
func (pt *ParserTests) TestFail1(c *C) {
	contents := `[`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected end-of-file")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 2)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail2(c *C) {
	contents := `[   true`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected boolean true")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 5)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail3(c *C) {
	contents := `[12`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected integer 12")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 2)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail4(c *C) {
	contents := `[1.2`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected floating point number 1.200000")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 2)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail5(c *C) {
	contents := "[ \t 1979-05-27T07:32:00Z"

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected date 1979-05-27T07:32:00Z")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 5)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail6(c *C) {
	contents := `[ "foo"  `

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected string \"foo\"")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 3)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail7(c *C) {
	contents := `[ foo] `

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected end-of-file")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 8)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail8(c *C) {
	contents := `true`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "expected section or option declaration")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 1)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail9(c *C) {
	contents := `id`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected end-of-file")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 3)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail10(c *C) {
	contents := `id  ,   `

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected character ','")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 5)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail11(c *C) {
	contents := `id    =   `

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected end-of-file")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 11)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail12(c *C) {
	contents := `id    =   =`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected character '='")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 11)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail13(c *C) {
	contents := `id    =  [`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected end-of-file")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 11)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail14(c *C) {
	contents := `id = [,`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected character ','")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 7)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail15(c *C) {
	contents := `id = [1  ,    `

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected end-of-file")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 15)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail16(c *C) {
	contents := `id = [true, 1`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type int64 as type bool")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 13)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail17(c *C) {
	contents := `id = [true, 1.0 ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type float64 as type bool")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 13)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail18(c *C) {
	contents := `id = [true, 1979-05-27T09:32:00Z ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type time.Time as type bool")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 13)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail19(c *C) {
	contents := `id = [true, "abc" ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type string as type bool")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 13)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail20(c *C) {
	contents := `id = [1, true]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type bool as type int64")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 10)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail21(c *C) {
	contents := `id = [2, 3.4 ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type float64 as type int64")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 10)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail22(c *C) {
	contents := `id = [3, 1979-05-27T09:32:00Z ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type time.Time as type int64")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 10)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail23(c *C) {
	contents := `id = [4, "abc" ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type string as type int64")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 10)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail24(c *C) {
	contents := `id = [1.2, true]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type bool as type float64")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 12)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail25(c *C) {
	contents := `id = [3.4, 1979-05-27T09:32:00Z ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type time.Time as type float64")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 12)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail26(c *C) {
	contents := `id = [4.5, "abc" ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type string as type float64")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 12)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail27(c *C) {
	contents := `id = [1979-05-27T07:32:00Z, true]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type bool as type time.Time")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 29)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail28(c *C) {
	contents := `id = [1979-05-27T07:32:00Z, 3 ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type int64 as type time.Time")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 29)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail29(c *C) {
	contents := `id = [1979-05-27T07:32:00Z, 1979.05 ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type float64 as type time.Time")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 29)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail30(c *C) {
	contents := `id = [1979-05-27T07:32:00Z, "abc" ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type string as type time.Time")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 29)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail31(c *C) {
	contents := `id = ["true", false]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type bool as type string")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 15)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail32(c *C) {
	contents := `id = ["true", 1 ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type int64 as type string")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 15)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail33(c *C) {
	contents := `id = ["true", 1979.6 ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type float64 as type string")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 15)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail34(c *C) {
	contents := `id = ["true", 1979-05-27T09:32:00Z ]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "cannot use type time.Time as type string")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 15)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail35(c *C) {
	contents := `[a] b = "c"`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected identifier b")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 5)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail36(c *C) {
	contents := `
[a]
[b]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "empty section a")
	c.Check(err.Line, Equals, 2)
	c.Check(err.Column, Equals, 3)
}

// Parse(): parser error.
func (pt *ParserTests) TestFail37(c *C) {
	contents := `
a = 12
[a]
a = 13
[b]
`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "empty section b")
	c.Check(err.Line, Equals, 5)
	c.Check(err.Column, Equals, 3)
	a, err1 := pt.config.GetInt("a")
	c.Check(err1, IsNil)
	c.Check(a, Equals, int64(12))
	a_a, err1 := pt.config.GetInt("a.a")
	c.Check(err1, IsNil)
	c.Check(a_a, Equals, int64(13))
}

// Parse(): parser error.
func (pt *ParserTests) TestFail38(c *C) {
	contents := `[[]]`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, ErrorMatches, "unexpected character '\\['")
	c.Check(err.Line, Equals, 1)
	c.Check(err.Column, Equals, 2)
}

// Parse(): lexer error.
func (pt *ParserTests) TestFail39(c *C) {
	contents := `*`

	pt.createTestEnv(c, contents)
	defer pt.cleanTestEnv(c)
	c.Check(pt.parser, NotNil)
	err := pt.parser.Parse(pt.config)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "unexpected '\\*' character")
}

// HasKey checks whether the configuration dictionnary records a given key
// or not.
type hasKeyChecker struct {
	*CheckerInfo
}

var HasKey Checker = &hasKeyChecker{
	&CheckerInfo{Name: "HasKey", Params: []string{"map", "key"}},
}

func (checker *hasKeyChecker) Check(params []interface{}, names []string) (result bool, err string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			err = fmt.Sprint(v)
		}
	}()
	cfg := params[0].(*config.Configuration)
	k := params[1].(string)
	result = cfg.HasOption(k)
	return
}
