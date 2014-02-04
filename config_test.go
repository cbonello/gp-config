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
	ConfigTests struct {
		pathnames []string
		config    *config.Configuration
	}
)

var (
	_ = Suite(&ConfigTests{})
)

func (ct *ConfigTests) createTmpFile(c *C, contents string) string {
	f, err := ioutil.TempFile("", "ConfigTest")
	c.Check(err, IsNil)
	defer f.Close()

	ct.pathnames = append(ct.pathnames, f.Name())
	n, err := f.Write([]byte(contents))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, len(contents))
	return f.Name()
}

func (ct *ConfigTests) deleteTmpFile(c *C, pathname string) {
	err := os.Remove(pathname)
	c.Assert(err, IsNil)
}

func (ct *ConfigTests) createTestEnv(c *C, contents []string) {
	ct.pathnames = []string{}
	for _, cnt := range contents {
		ct.createTmpFile(c, cnt)
	}
	ct.config = config.NewConfiguration()
	for _, fn := range ct.pathnames {
		err := ct.config.LoadFile(fn)
		c.Assert(err, IsNil)
	}
}

func (ct *ConfigTests) cleanTestEnv(c *C) {
	ct.config = nil
	for _, p := range ct.pathnames {
		ct.deleteTmpFile(c, p)
	}
	ct.pathnames = nil
}

// Len().
func (ct *ConfigTests) TestLen1(c *C) {
	var cfg *config.Configuration

	l := cfg.Len()
	c.Check(l, Equals, 0)
}

// LoadFile(): load one file.
// Len().
func (ct *ConfigTests) TestLoadFile1(c *C) {
	contents := `
foo = "bar"

[values]
	boolean = false
	integer = 12
	fp = 3.1415
	date = 2013-10-25T16:22:00Z
	string = "Hello World!"

[arrays]
	boolean = [false, true]
	integer = [12, 34]
	fp = [3.1415, 5.1413]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]
	string = ["Hello World!", "foo bar"]
`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 11)
}

// LoadFile(): load more than one file.
func (ct *ConfigTests) TestLoadFile2(c *C) {
	contents1 := `
[values]
	boolean = false
	integer = 12`

	contents2 := `
[values]
	integer = 34
	fp = 3.1415
	date = 2013-10-25T16:22:00Z
	string = "Hello World!"
`
	contents3 := `
[arrays]
	boolean = [false, true]
	integer = [12, 34]
	fp = [3.1415, 5.1413]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]
	string = ["Hello World!", "foo bar"]
`

	ct.createTestEnv(c, []string{contents1, contents2, contents3})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 10)
}

// LoadFile(): load file containing a syntax error.
// Len().
func (ct *ConfigTests) TestLoadFile3(c *C) {
	contents := `[values`

	ct.pathnames = []string{}
	ct.createTmpFile(c, contents)

	ct.config = config.NewConfiguration()
	defer ct.cleanTestEnv(c)

	err := ct.config.LoadFile(ct.pathnames[0])
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "unexpected end-of-file")
}

// LoadString(): load one string.
func (ct *ConfigTests) TestLoadString1(c *C) {
	contents := `
[values]
	boolean = false
	integer = 12
	fp = 3.1415
	date = 2013-10-25T16:22:00Z
	string = "Hello World!"

[arrays]
	boolean = [false, true]
	integer = [12, 34]
	fp = [3.1415, 5.1413]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]
	string = ["Hello World!", "foo bar"]`

	ct.config = config.NewConfiguration()
	defer ct.cleanTestEnv(c)
	err := ct.config.LoadString(contents)
	c.Check(err, IsNil)

	c.Check(ct.config.Len(), Equals, 10)
}

// LoadString(): load more than one string.
func (ct *ConfigTests) TestLoadString2(c *C) {
	contents1 := `
[values]
	boolean = false
	integer = 12
	fp = 3.1415
	date = 2013-10-25T16:22:00Z
	string = "Hello World!"`

	contents2 := `
[values]
	boolean = true

[arrays]
	boolean = [false, true]
	integer = [12, 34]
`
	contents3 := `
[arrays]
	fp = [3.1415, 5.1413]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]
	string = ["Hello World!", "foo bar"]`

	ct.config = config.NewConfiguration()
	defer ct.cleanTestEnv(c)
	err := ct.config.LoadString(contents1)
	c.Check(err, IsNil)
	err = ct.config.LoadString(contents2)
	c.Check(err, IsNil)
	err = ct.config.LoadString(contents3)
	c.Check(err, IsNil)

	c.Check(ct.config.Len(), Equals, 10)
}

// LoadString(): load string containing a syntax error.
// Len().
func (ct *ConfigTests) TestLoadString3(c *C) {
	contents := `[`

	ct.config = config.NewConfiguration()
	defer ct.cleanTestEnv(c)
	err := ct.config.LoadString(contents)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "unexpected end-of-file")
}

// HasOption(): nil structure.
func (ct *ConfigTests) TestHasOption1(c *C) {
	var cfg *config.Configuration

	exists := cfg.HasOption("foo")
	c.Check(exists, Equals, false)
}

// HasOption().
func (ct *ConfigTests) TestHasOption2(c *C) {
	contents := `
[values]
	boolean = false
	integer = 12
	fp = 3.1415
	date = 2013-10-25T16:22:00Z
	string = "Hello World!"

[arrays]
	boolean = [false, true]
	integer = [12, 34]
	fp = [3.1415, 5.1413]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]
	string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 10)

	exists := ct.config.HasOption("values.boolean")
	c.Check(exists, Equals, true)
	// Options are case insensitive.
	exists = ct.config.HasOption("VALUES.BOOLean")
	c.Check(exists, Equals, true)
	exists = ct.config.HasOption("arrays.date")
	c.Check(exists, Equals, true)
	exists = ct.config.HasOption("values.birthday")
	c.Check(exists, Equals, false)
}

// Sections(): nil structure.
func (ct *ConfigTests) TestSections1(c *C) {
	ct.config = config.NewConfiguration()
	defer ct.cleanTestEnv(c)

	sections := ct.config.Sections()
	c.Check(len(sections), Equals, 0)
}

// Sections(): no section defined.
func (ct *ConfigTests) TestSections2(c *C) {
	contents := `
string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	sections := ct.config.Sections()
	c.Check(len(sections), Equals, 1)
	c.Check(sections, EqualSlice, []string{""})
}

// Sections(): global options only.
func (ct *ConfigTests) TestSections3(c *C) {
	contents := `
boolean = false
integer = 12
fp = 3.1415
date = 2013-10-25T16:22:00Z
string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 5)

	sections := ct.config.Sections()
	c.Check(len(sections), Equals, 1)
	c.Check(sections, EqualSlice, []string{""})

	c.Check(ct.config.Len(), Equals, 5)
}

// Sections(): local options only.
func (ct *ConfigTests) TestSections4(c *C) {
	contents := `
[values]
	boolean = false
	integer = 12
	fp = 3.1415
	date = 2013-10-25T16:22:00Z
	string = "Hello World!"

[arrays]
	boolean = [false, true]
	integer = [12, 34]
	fp = [3.1415, 5.1413]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]
	string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 10)

	sections := ct.config.Sections()
	c.Check(len(sections), Equals, 2)
	c.Check(sections, EqualSlice, []string{"arrays", "values"})
}

// Sections(): global and local options.
func (ct *ConfigTests) TestSections5(c *C) {
	contents := `
string = ["Hello World!", "foo bar"]

[values]
	boolean = false

[arrays]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 3)

	sections := ct.config.Sections()
	c.Check(len(sections), Equals, 3)
	c.Check(sections, EqualSlice, []string{"", "arrays", "values"})
}

// IsSection(): nil structure.
func (ct *ConfigTests) TestIsSection1(c *C) {
	var cfg *config.Configuration

	section := cfg.IsSection("")
	c.Check(section, Equals, false)
}

// IsSection(): global and local options.
func (ct *ConfigTests) TestIsSection2(c *C) {
	contents := `
date = 2013-10-25T16:22:00Z
[values]
	boolean = false
	integer = 12
	fp = 3.1415
	date = 2013-10-25T16:22:00Z
	string = "Hello World!"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 6)

	section := ct.config.IsSection("")
	c.Check(section, Equals, true)
	section = ct.config.IsSection("values")
	c.Check(section, Equals, true)
	section = ct.config.IsSection("vALUEs")
	c.Check(section, Equals, true)
}

// IsSection(): undefined section.
func (ct *ConfigTests) TestIsSection3(c *C) {
	contents := `fp = 3.1415`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	section := ct.config.IsSection("values")
	c.Check(section, Equals, false)
}

// Options(): nil structure.
func (ct *ConfigTests) TestOptions1(c *C) {
	var cfg *config.Configuration

	options := cfg.Options("")
	c.Check(len(options), Equals, 0)
	c.Check(options, EqualSlice, []string{})
}

// Options(): no options.
func (ct *ConfigTests) TestOptions2(c *C) {
	contents := ``

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 0)

	options := ct.config.Options("")
	c.Check(len(options), Equals, 0)
	c.Check(options, EqualSlice, []string{})
}

// Options(): global and local options.
func (ct *ConfigTests) TestOptions3(c *C) {
	contents := `
foo = 1.0

[v]
	boolean = false
	integer = 12
	fp = 3.1415
[w]
	date = 1515-10-25T16:22:00Z
	string = "Hello World!"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 6)

	options0 := ct.config.Options("")
	c.Check(len(options0), Equals, 1)
	c.Check(options0, EqualSlice, []string{"foo"})

	options1 := ct.config.Options("v")
	c.Check(len(options1), Equals, 3)
	c.Check(options1, EqualSlice, []string{"v.boolean", "v.fp", "v.integer"})

	options2 := ct.config.Options("W")
	c.Check(len(options2), Equals, 2)
	c.Check(options2, EqualSlice, []string{"w.date", "w.string"})
}

// Get(): nil structure.
func (ct *ConfigTests) TestGet1(c *C) {
	var cfg *config.Configuration

	value, err := cfg.Get("")
	c.Check(value, Equals, nil)
	c.Check(err, ErrorMatches, "'': unknown option")
}

// Get(): unknown option.
func (ct *ConfigTests) TestGet2(c *C) {
	contents := `
[values]
	boolean = false
[arrays]
	boolean = [false, true]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value, err := ct.config.Get("v.boolean")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.boolean': unknown option")
	c.Check(value, IsNil)
}

// Get(): scalars (outside a section).
func (ct *ConfigTests) TestGet3(c *C) {
	contents := `
boolean = false
integer = 12
fp = 3.1415
date = 2013-10-25T16:22:00Z
string = "Hello World!"`

	ct.config = config.NewConfiguration()
	err := ct.config.LoadString(contents)
	c.Check(err, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 5)

	value0, err0 := ct.config.Get("boolean")
	c.Check(err0, IsNil)
	c.Check(value0.(bool), Equals, false)
	value1, err1 := ct.config.Get("integer")
	c.Check(err1, IsNil)
	c.Check(value1.(int64), Equals, int64(12))
	value2, err2 := ct.config.Get("fp")
	c.Check(err2, IsNil)
	c.Check(value2.(float64), Equals, float64(3.1415))
	expected3, _ := time.Parse(time.RFC3339, "2013-10-25T16:22:00Z")
	value3, err3 := ct.config.Get("date")
	c.Check(err3, IsNil)
	c.Check(value3.(time.Time), Equals, expected3)
	value4, err4 := ct.config.Get("string")
	c.Check(err4, IsNil)
	c.Check(value4.(string), Equals, "Hello World!")
}

// Get(): scalars (inside a section).
func (ct *ConfigTests) TestGet4(c *C) {
	contents := `
[values]
	boolean = false
	integer = 12
	fp = 3.1415
	date = 2013-10-25T16:22:00Z
	string = "Hello World!"`

	ct.config = config.NewConfiguration()
	err := ct.config.LoadString(contents)
	c.Check(err, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 5)

	value0, err0 := ct.config.Get("values.boolean")
	c.Check(err0, IsNil)
	c.Check(value0.(bool), Equals, false)
	value1, err1 := ct.config.Get("values.integer")
	c.Check(err1, IsNil)
	c.Check(value1.(int64), Equals, int64(12))
	value2, err2 := ct.config.Get("values.fp")
	c.Check(err2, IsNil)
	c.Check(value2.(float64), Equals, float64(3.1415))
	expected3, _ := time.Parse(time.RFC3339, "2013-10-25T16:22:00Z")
	value3, err3 := ct.config.Get("values.date")
	c.Check(err3, IsNil)
	c.Check(value3.(time.Time), Equals, expected3)
	value4, err4 := ct.config.Get("values.string")
	c.Check(err4, IsNil)
	c.Check(value4.(string), Equals, "Hello World!")
}

// Get(): arrays (outside a section).
func (ct *ConfigTests) TestGet5(c *C) {
	contents := `
boolean = [false, true]
integer = [12, 34]
fp = [3.1415, 5.1413]
date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]
string = ["Hello World!", "foo bar"]`

	ct.config = config.NewConfiguration()
	err := ct.config.LoadString(contents)
	c.Check(err, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 5)

	arr0, err0 := ct.config.Get("boolean")
	c.Check(err0, IsNil)
	c.Check(arr0, EqualSlice, []bool{false, true})
	arr1, err1 := ct.config.Get("integer")
	c.Check(err1, IsNil)
	c.Check(arr1, EqualSlice, []int64{12, 34})
	arr2, err2 := ct.config.Get("fp")
	c.Check(err2, IsNil)
	c.Check(arr2, EqualSlice, []float64{3.1415, 5.1413})
	expected31, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	expected32, _ := time.Parse(time.RFC3339, "2013-10-25T16:22:00Z")
	arr3, err3 := ct.config.Get("date")
	c.Check(err3, IsNil)
	c.Check(arr3, EqualSlice, []time.Time{expected31, expected32})
	arr4, err4 := ct.config.Get("string")
	c.Check(err4, IsNil)
	c.Check(arr4, EqualSlice, []string{"Hello World!", "foo bar"})
}

// Get(): arrays (outside a section).
func (ct *ConfigTests) TestGet6(c *C) {
	contents := `
[arrays]
	boolean = [false, true]
	integer = [12, 34]
	fp = [3.1415, 5.1413]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]
	string = ["Hello World!", "foo bar"]`

	ct.config = config.NewConfiguration()
	err := ct.config.LoadString(contents)
	c.Check(err, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 5)

	arr0, err0 := ct.config.Get("arrays.boolean")
	c.Check(err0, IsNil)
	c.Check(arr0, EqualSlice, []bool{false, true})
	arr1, err1 := ct.config.Get("arrays.integer")
	c.Check(err1, IsNil)
	c.Check(arr1, EqualSlice, []int64{12, 34})
	arr2, err2 := ct.config.Get("arrays.fp")
	c.Check(err2, IsNil)
	c.Check(arr2, EqualSlice, []float64{3.1415, 5.1413})
	expected31, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	expected32, _ := time.Parse(time.RFC3339, "2013-10-25T16:22:00Z")
	arr3, err3 := ct.config.Get("arrays.date")
	c.Check(err3, IsNil)
	c.Check(arr3, EqualSlice, []time.Time{expected31, expected32})
	arr4, err4 := ct.config.Get("arrays.string")
	c.Check(err4, IsNil)
	c.Check(arr4, EqualSlice, []string{"Hello World!", "foo bar"})
}

// GetBool(): nil structure.
func (ct *ConfigTests) TestGetBool1(c *C) {
	var cfg *config.Configuration

	value, err := cfg.GetBool("foo")
	c.Check(value, Equals, bool(false))
	c.Check(err, ErrorMatches, "'foo': unknown option")
}

// GetBool(): unknown option.
func (ct *ConfigTests) TestGetBool2(c *C) {
	contents := `
[values]
	integer = 1`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err := ct.config.GetBool("v.integer")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.integer': unknown option")
}

// GetBool(): wrong type.
func (ct *ConfigTests) TestGetBool3(c *C) {
	contents := `
[values]
	integer = 12`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetBool("values.integer")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.integer': not a boolean")
}

// GetBool(): wrong type.
func (ct *ConfigTests) TestGetBool4(c *C) {
	contents := `
[values]
	fp = 3.1415`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetBool("values.fp")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.fp': not a boolean")
}

// GetBool(): wrong type.
func (ct *ConfigTests) TestGetBool5(c *C) {
	contents := `
[values]
	date = 2013-10-25T16:22:00Z`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetBool("values.date")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.date': not a boolean")
}

// GetBool(): wrong type.
func (ct *ConfigTests) TestGetBool6(c *C) {
	contents := `
[values]
	string = "Hello World!"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetBool("values.string")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.string': not a boolean")
}

// GetBool(): local and global options.
func (ct *ConfigTests) TestGetBool7(c *C) {
	contents := `
boolean = true
[arrays]
	boolean = true`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0, err0 := ct.config.GetBool("boolean")
	c.Check(err0, IsNil)
	c.Check(value0, EqualSlice, true)
	value1, err1 := ct.config.GetBool("arrays.boolean")
	c.Check(err1, IsNil)
	c.Check(value1, EqualSlice, true)
}

// GetBoolDefault(): nil structure.
func (ct *ConfigTests) TestGetBoolDefault1(c *C) {
	var cfg *config.Configuration

	value := cfg.GetBoolDefault("foo", false)
	c.Check(value, Equals, false)
}

// GetBoolDefault(): unknown option.
func (ct *ConfigTests) TestGetBoolDefault2(c *C) {
	contents := `
[values]
	boolean = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetBoolDefault("foo", true)
	c.Check(value, Equals, true)
}

// GetBoolDefault(): wrong type.
func (ct *ConfigTests) TestGetBoolDefault3(c *C) {
	contents := `
[values]
	integer = 12`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetBoolDefault("values.integer", true)
	c.Check(value, Equals, true)
}

// GetBoolDefault(): wrong type.
func (ct *ConfigTests) TestGetBoolDefault4(c *C) {
	contents := `
[values]
	fp = 3.1415`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetBoolDefault("values.fp", true)
	c.Check(value, Equals, true)
}

// GetBoolDefault(): wrong type.
func (ct *ConfigTests) TestGetBoolDefault5(c *C) {
	contents := `
[values]
	date = 2013-10-25T16:22:00Z`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetBoolDefault("values.date", true)
	c.Check(value, Equals, true)
}

// GetBoolDefault(): wrong type.
func (ct *ConfigTests) TestGetBoolDefault6(c *C) {
	contents := `
[values]
	string = "Hello World!"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetBoolDefault("values.string", true)
	c.Check(value, Equals, true)
}

// GetBoolDefault(): local and global options.
func (ct *ConfigTests) TestGetBoolDefault7(c *C) {
	contents := `
boolean = true
[arrays]
	boolean = true`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0 := ct.config.GetBoolDefault("boolean", false)
	c.Check(value0, Equals, true)
	value1 := ct.config.GetBoolDefault("arrays.boolean", false)
	c.Check(value1, Equals, true)
}

// GetInt(): nil structure.
func (ct *ConfigTests) TestGetInt1(c *C) {
	var cfg *config.Configuration

	_, err := cfg.GetInt("foo")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'foo': unknown option")
}

// GetInt(): unknown option.
func (ct *ConfigTests) TestGetInt2(c *C) {
	contents := `
[values]
	integer = 5
[arrays]
	integer = [89, 4]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value, err := ct.config.Get("v.integer")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.integer': unknown option")
	c.Check(value, IsNil)
}

// GetInt(): wrong type.
func (ct *ConfigTests) TestGetInt3(c *C) {
	contents := `
[values]
	boolean = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetInt("values.boolean")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.boolean': not an integer")
}

// GetInt(): wrong type.
func (ct *ConfigTests) TestGetInt4(c *C) {
	contents := `
[values]
	fp = 3.1415`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetInt("values.fp")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.fp': not an integer")
}

// GetInt(): wrong type.
func (ct *ConfigTests) TestGetInt5(c *C) {
	contents := `
[values]
	date = 2013-10-25T16:22:00Z`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetInt("values.date")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.date': not an integer")
}

// GetInt(): wrong type.
func (ct *ConfigTests) TestGetInt6(c *C) {
	contents := `
[values]
	string = "Hello World!"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetInt("values.string")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.string': not an integer")
}

// GetInt(): local and global options.
func (ct *ConfigTests) TestGetInt7(c *C) {
	contents := `
integer = 12
[arrays]
	integer = 98`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0, err0 := ct.config.GetInt("integer")
	c.Check(err0, IsNil)
	c.Check(value0, Equals, int64(12))
	value1, err1 := ct.config.GetInt("arrays.integer")
	c.Check(err1, IsNil)
	c.Check(value1, Equals, int64(98))
}

// GetIntDefault(): nil structure.
func (ct *ConfigTests) TestGetIntDefault1(c *C) {
	var cfg *config.Configuration

	value := cfg.GetIntDefault("foo", 9)
	c.Check(value, Equals, int64(9))
}

// GetIntDefault(): unknown option.
func (ct *ConfigTests) TestGetIntDefault2(c *C) {
	contents := `
[values]
	integer = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetIntDefault("integer", 1234)
	c.Check(value, Equals, int64(1234))
}

// GetIntDefault(): wrong type.
func (ct *ConfigTests) TestGetIntDefault3(c *C) {
	contents := `
[values]
	boolean = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetIntDefault("values.integer", 10)
	c.Check(value, Equals, int64(10))
}

// GetIntDefault(): wrong type.
func (ct *ConfigTests) TestGetIntDefault4(c *C) {
	contents := `
[values]
	fp = 3.1415`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetIntDefault("values.fp", 20)
	c.Check(value, Equals, int64(20))
}

// GetIntDefault(): wrong type.
func (ct *ConfigTests) TestGetIntDefault5(c *C) {
	contents := `
[values]
	date = 2013-10-25T16:22:00Z`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetIntDefault("values.date", 30)
	c.Check(value, Equals, int64(30))
}

// GetIntDefault(): wrong type.
func (ct *ConfigTests) TestGetIntDefault6(c *C) {
	contents := `
[values]
	string = "Hello World!"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetIntDefault("values.string", 40)
	c.Check(value, Equals, int64(40))
}

// GetIntDefault(): local and global options.
func (ct *ConfigTests) TestGetIntDefault7(c *C) {
	contents := `
integer = 12
[arrays]
	integer = 98`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0 := ct.config.GetIntDefault("integer", 0)
	c.Check(value0, Equals, int64(12))
	value1 := ct.config.GetIntDefault("arrays.integer", 1)
	c.Check(value1, Equals, int64(98))
}

// GetFloat(): nil structure.
func (ct *ConfigTests) TestGetFloat1(c *C) {
	var cfg *config.Configuration

	value, err := cfg.GetFloat("foo")
	c.Check(value, Equals, float64(0))
	c.Check(err, ErrorMatches, "'foo': unknown option")
}

// GetFloat(): unknown option.
func (ct *ConfigTests) TestGetFloat2(c *C) {
	contents := `
[values]
	fp = 1`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err := ct.config.GetFloat("v.fp")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.fp': unknown option")
}

// GetFloat(): wrong type.
func (ct *ConfigTests) TestGetFloat3(c *C) {
	contents := `
[values]
	boolean = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetFloat("values.boolean")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.boolean': not a floating-point number")
}

// GetFloat(): wrong type.
func (ct *ConfigTests) TestGetFloat4(c *C) {
	contents := `
[values]
	integer = 12`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetFloat("values.integer")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.integer': not a floating-point number")
}

// GetFloat(): wrong type.
func (ct *ConfigTests) TestGetFloat5(c *C) {
	contents := `
[values]
	date = 2013-10-25T16:22:00Z`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetFloat("values.date")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.date': not a floating-point number")
}

// GetFloat(): wrong type.
func (ct *ConfigTests) TestGetFloat6(c *C) {
	contents := `
[values]
	string = "Hello World!"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err := ct.config.GetFloat("values.string")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'values.string': not a floating-point number")
}

// GetFloat(): local and global options.
func (ct *ConfigTests) TestGetFloat7(c *C) {
	contents := `
fp = 12.21
[arrays]
	fp = 98.5`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0, err0 := ct.config.GetFloat("fp")
	c.Check(err0, IsNil)
	c.Check(value0, Equals, float64(12.21))
	value1, err1 := ct.config.GetFloat("arrays.fp")
	c.Check(err1, IsNil)
	c.Check(value1, Equals, float64(98.5))
}

// GetFloatDefault(): nil structure.
func (ct *ConfigTests) TestGetFloatDefault1(c *C) {
	var cfg *config.Configuration

	value := cfg.GetFloatDefault("foo", 9.)
	c.Check(value, Equals, float64(9))
}

// GetFloatDefault(): unknown option.
func (ct *ConfigTests) TestGetFloatDefault2(c *C) {
	contents := `
[values]
	integer = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetFloatDefault("integer", 1234.)
	c.Check(value, Equals, float64(1234))
}

// GetFloatDefault(): wrong type.
func (ct *ConfigTests) TestGetFloatDefault3(c *C) {
	contents := `
[values]
	boolean = true`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetFloatDefault("values.boolean", 20.0)
	c.Check(value, Equals, float64(20))
}

// GetFloatDefault(): wrong type.
func (ct *ConfigTests) TestGetFloatDefault4(c *C) {
	contents := `
[values]
	boolean = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetFloatDefault("values.integer", 1.0)
	c.Check(value, Equals, float64(1.0))
}

// GetFloatDefault(): wrong type.
func (ct *ConfigTests) TestGetFloatDefault5(c *C) {
	contents := `
[values]
	date = 2013-10-25T16:22:00Z`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetFloatDefault("values.date", 30.03)
	c.Check(value, Equals, float64(30.03))
}

// GetFloatDefault(): wrong type.
func (ct *ConfigTests) TestGetFloatDefault6(c *C) {
	contents := `
[values]
	string = "Hello World!"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetFloatDefault("values.string", 4.4)
	c.Check(value, Equals, float64(4.4))
}

// GetFloatDefault(): local and global options.
func (ct *ConfigTests) TestGetFloatDefault7(c *C) {
	contents := `
fp = 12.21
[arrays]
	fp = 98.5`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0 := ct.config.GetFloatDefault("fp", 14.7)
	c.Check(value0, Equals, float64(12.21))
	value1 := ct.config.GetFloatDefault("arrays.fp", 65.988)
	c.Check(value1, Equals, float64(98.5))
}

// GetDate(): nil structure.
func (ct *ConfigTests) TestGetDate1(c *C) {
	var cfg *config.Configuration

	_, err := cfg.GetDate("foo")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'foo': unknown option")
}

// GetDate(): unknown option.
func (ct *ConfigTests) TestGetDate2(c *C) {
	contents := `
[values]
	fp = 1`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err := ct.config.GetFloat("v.fp")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.fp': unknown option")
}

// GetDate(): unknown option.
func (ct *ConfigTests) TestGetDate3(c *C) {
	contents := `
[values]
	boolean = true`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetDate("values.boolean")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.boolean': not a date")
}

// GetDate(): unknown option.
func (ct *ConfigTests) TestGetDate4(c *C) {
	contents := `
[values]
	integer = 12`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetDate("values.integer")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.integer': not a date")
}

// GetDate(): unknown option.
func (ct *ConfigTests) TestGetDate5(c *C) {
	contents := `
[values]
	fp = 3.1415`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetDate("values.fp")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.fp': not a date")
}

// GetDate(): unknown option.
func (ct *ConfigTests) TestGetDate6(c *C) {
	contents := `
[values]
	string = "Hello World!"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetDate("values.string")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.string': not a date")
}

// GetDate(): local and global options.
func (ct *ConfigTests) TestGetDate7(c *C) {
	contents := `
date = 2013-10-24T16:22:00Z
[arrays]
	date = 2013-12-03T16:22:00Z`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	date0, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	value0, err0 := ct.config.GetDate("date")
	c.Check(err0, IsNil)
	c.Check(value0, Equals, date0)
	date1, _ := time.Parse(time.RFC3339, "2013-12-03T16:22:00Z")
	value1, err1 := ct.config.GetDate("arrays.date")
	c.Check(err1, IsNil)
	c.Check(value1, Equals, date1)
}

// GetFloatDefault(): nil structure.
func (ct *ConfigTests) TestGetDateDefault1(c *C) {
	var cfg *config.Configuration

	date, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	value := cfg.GetDateDefault("foo", date)
	c.Check(value, Equals, date)
}

// GetDateDefault(): unknown option.
func (ct *ConfigTests) TestGetDateDefault2(c *C) {
	contents := `
[values]
	integer = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	date, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	value := ct.config.GetDateDefault("integer", date)
	c.Check(value, Equals, date)
}

// GetDateDefault(): wrong type.
func (ct *ConfigTests) TestGetDateDefault3(c *C) {
	contents := `
[values]
	boolean = true`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	date, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	value := ct.config.GetDateDefault("values.boolean", date)
	c.Check(value, Equals, date)
}

// GetDateDefault(): wrong type.
func (ct *ConfigTests) TestGetDateDefault4(c *C) {
	contents := `
[values]
	boolean = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	date, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	value := ct.config.GetDateDefault("values.boolean", date)
	c.Check(value, Equals, date)
}

// GetDateDefault(): wrong type.
func (ct *ConfigTests) TestGetDateDefault5(c *C) {
	contents := `
[values]
	integer = 2013`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	date, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	value := ct.config.GetDateDefault("values.integer", date)
	c.Check(value, Equals, date)
}

// GetDateDefault(): wrong type.
func (ct *ConfigTests) TestGetDateDefault6(c *C) {
	contents := `
[values]
	string = "Hello World!"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	date, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	value := ct.config.GetDateDefault("values.string", date)
	c.Check(value, Equals, date)
}

// GetDateDefault(): local and global options.
func (ct *ConfigTests) TestGetDateDefault7(c *C) {
	contents := `
date = 2013-10-24T16:22:00Z
[arrays]
	date = 2013-12-03T16:22:00Z`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	date0, _ := time.Parse(time.RFC3339, "1900-10-24T16:22:00Z")
	expected0, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	value0 := ct.config.GetDateDefault("date", date0)
	c.Check(value0, Equals, expected0)
	date1, _ := time.Parse(time.RFC3339, "1800-12-03T16:22:00Z")
	expected1, _ := time.Parse(time.RFC3339, "2013-12-03T16:22:00Z")
	value1 := ct.config.GetDateDefault("arrays.date", date1)
	c.Check(value1, Equals, expected1)
}

// GetString(): nil structure.
func (ct *ConfigTests) TestGetString1(c *C) {
	var cfg *config.Configuration

	_, err := cfg.GetString("foo")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'foo': unknown option")
}

// GetString(): unknown option.
func (ct *ConfigTests) TestGetString2(c *C) {
	contents := `
[values]
	string = "foo-bar"
`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err := ct.config.GetString("v.string")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.string': unknown option")
}

// GetString(): wrong type.
func (ct *ConfigTests) TestGetString3(c *C) {
	contents := `
[values]
	boolean = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetString("values.boolean")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.boolean': not a string")
}

// GetString(): wrong type.
func (ct *ConfigTests) TestGetString4(c *C) {
	contents := `
[values]
	integer = 12`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetString("values.integer")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.integer': not a string")
}

// GetString(): wrong type.
func (ct *ConfigTests) TestGetString5(c *C) {
	contents := `
[values]
	fp = 3.1415`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetString("values.fp")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.fp': not a string")
}

// GetString(): wrong type.
func (ct *ConfigTests) TestGetString6(c *C) {
	contents := `
[values]
	date = 2013-10-25T16:22:00Z`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetString("values.date")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'values.date': not a string")
}

// GetString(): local and global options.
func (ct *ConfigTests) TestGetString7(c *C) {
	contents := `
string = "2013-10-24T16:22:00Z"
[arrays]
	string = "foo"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0, err0 := ct.config.GetString("string")
	c.Check(err0, IsNil)
	c.Check(value0, Equals, "2013-10-24T16:22:00Z")
	value1, err1 := ct.config.GetString("arrays.string")
	c.Check(err1, IsNil)
	c.Check(value1, Equals, "foo")
}

// GetStringDefault(): nil structure.
func (ct *ConfigTests) GetStringDefault(c *C) {
	var cfg *config.Configuration

	value := cfg.GetStringDefault("foo", "foo")
	c.Check(value, Equals, "foo")
}

// GetStringDefault(): unknown option.
func (ct *ConfigTests) TestGetStringDefault2(c *C) {
	contents := `
[values]
	integer = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetStringDefault("integer", "bar")
	c.Check(value, Equals, "bar")
}

// GetStringDefault(): wrong type.
func (ct *ConfigTests) TestGetStringDefault3(c *C) {
	contents := `
[values]
	boolean = true`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetStringDefault("values.boolean", "foobar")
	c.Check(value, Equals, "foobar")
}

// GetStringDefault(): wrong type.
func (ct *ConfigTests) TestGetStringDefault4(c *C) {
	contents := `
[values]
	integer = false`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetStringDefault("values.integer", "barfoo")
	c.Check(value, Equals, "barfoo")
}

// GetStringDefault(): wrong type.
func (ct *ConfigTests) TestGetStringDefault5(c *C) {
	contents := `
[values]
	fp = 98.7`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetStringDefault("values.string", "rab")
	c.Check(value, Equals, "rab")
}

// GetStringDefault(): wrong type.
func (ct *ConfigTests) TestGetStringDefault6(c *C) {
	contents := `
[values]
	date = 2013-10-25T16:22:00Z`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	value := ct.config.GetStringDefault("values.date", "oof")
	c.Check(value, Equals, "oof")
}

// GetStringDefault(): local and global options.
func (ct *ConfigTests) TestGetStringDefault7(c *C) {
	contents := `
string = "2013-10-24T16:22:00Z"
[arrays]
	string = "foo"`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0 := ct.config.GetStringDefault("string", "abcd")
	c.Check(value0, Equals, "2013-10-24T16:22:00Z")
	value1 := ct.config.GetStringDefault("arrays.string", "efgh")
	c.Check(value1, Equals, "foo")
}

// GetBoolArray(): nil structure.
func (ct *ConfigTests) TestGetBoolArray1(c *C) {
	var cfg *config.Configuration

	_, err := cfg.GetBoolArray("foo")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'foo': unknown option")
}

// GetBoolArray(): unknown option.
func (ct *ConfigTests) TestGetBoolArray2(c *C) {
	contents := `
[values]
	string = ["foo", "bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err := ct.config.GetBoolArray("v.string")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.string': unknown option")
}

// GetBoolArray(): wrong type.
func (ct *ConfigTests) TestGetBoolArray3(c *C) {
	contents := `
[arrays]
	integer = [12, 34]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetBoolArray("arrays.integer")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.integer': not an array of booleans")
}

// GetBoolArray(): wrong type.
func (ct *ConfigTests) TestGetBoolArray4(c *C) {
	contents := `
[arrays]
	fp = [3.1415, 5.1413]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetBoolArray("arrays.fp")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.fp': not an array of booleans")
}

// GetBoolArray(): wrong type.
func (ct *ConfigTests) TestGetBoolArray5(c *C) {
	contents := `
[arrays]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetBoolArray("arrays.date")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.date': not an array of booleans")
}

// GetBoolArray(): wrong type.
func (ct *ConfigTests) TestGetBoolArray6(c *C) {
	contents := `
[arrays]
	string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetBoolArray("arrays.string")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.string': not an array of booleans")
}

// GetBoolArray(): local and global options.
func (ct *ConfigTests) TestGetBoolArray7(c *C) {
	contents := `
booleans = [false, false]
[arrays]
	booleans = [true, true]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0, err0 := ct.config.GetBoolArray("booleans")
	c.Check(err0, IsNil)
	c.Check(value0, EqualSlice, []bool{false, false})
	value1, err1 := ct.config.GetBoolArray("arrays.booleans")
	c.Check(err1, IsNil)
	c.Check(value1, EqualSlice, []bool{true, true})
}

// GetBoolArrayDefault(): nil structure.
func (ct *ConfigTests) TestGGetBoolArrayDefault1(c *C) {
	var cfg *config.Configuration

	arr := []bool{true}
	value := cfg.GetBoolArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetBoolArrayDefault(): unknown option.
func (ct *ConfigTests) TestGetBoolArrayDefault2(c *C) {
	contents := `
[values]
	boolean = ["foo", "bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []bool{true}
	value := ct.config.GetBoolArrayDefault("v.boolean", arr)
	c.Check(value, EqualSlice, arr)
}

// GetBoolArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetBoolArrayDefault3(c *C) {
	contents := `
[arrays]
	integer = [12, 34]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []bool{true}
	value := ct.config.GetBoolArrayDefault("arrays.integer", arr)
	c.Check(value, EqualSlice, arr)
}

// GetBoolArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetBoolArrayDefault4(c *C) {
	contents := `
[arrays]
	fp = [3.1415, 5.1413]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []bool{true}
	value := ct.config.GetBoolArrayDefault("arrays.fp", arr)
	c.Check(value, EqualSlice, arr)
}

// GetBoolArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetBoolArrayDefault5(c *C) {
	contents := `
[arrays]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []bool{true}
	value := ct.config.GetBoolArrayDefault("arrays.date", arr)
	c.Check(value, EqualSlice, arr)
}

// GetBoolArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetBoolArrayDefault6(c *C) {
	contents := `
[arrays]
	string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []bool{true}
	value := ct.config.GetBoolArrayDefault("arrays.string", arr)
	c.Check(value, EqualSlice, arr)
}

// GetBoolArrayDefault(): local and global options.
func (ct *ConfigTests) TestGetBoolArrayDefault7(c *C) {
	contents := `
booleans = [true, true]
[arrays]
	booleans = [true, true]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0 := ct.config.GetBoolArrayDefault("booleans", []bool{false, false})
	c.Check(value0, EqualSlice, []bool{true, true})
	value1 := ct.config.GetBoolArrayDefault("arrays.booleans", []bool{false, false})
	c.Check(value1, EqualSlice, []bool{true, true})
}

// GetIntArray(): nil structure.
func (ct *ConfigTests) TestGetIntArray1(c *C) {
	var cfg *config.Configuration

	_, err := cfg.GetIntArray("foo")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'foo': unknown option")
}

// GetIntArray(): unknown option.
func (ct *ConfigTests) TestGetIntArray2(c *C) {
	contents := `
[values]
	string = ["foo", "bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err := ct.config.GetIntArray("v.string")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.string': unknown option")
}

// GetIntArray(): wrong type.
func (ct *ConfigTests) TestGetIntArray3(c *C) {
	contents := `
[arrays]
	boolean = [false, true]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetIntArray("arrays.boolean")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.boolean': not an array of integers")
}

// GetIntArray(): wrong type.
func (ct *ConfigTests) TestGetIntArray4(c *C) {
	contents := `
[arrays]
	fp = [3.1415, 5.1413]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetIntArray("arrays.fp")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.fp': not an array of integers")
}

// GetIntArray(): wrong type.
func (ct *ConfigTests) TestGetIntArray5(c *C) {
	contents := `
[arrays]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetIntArray("arrays.date")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.date': not an array of integers")
}

// GetIntArray(): wrong type.
func (ct *ConfigTests) TestGetIntArray6(c *C) {
	contents := `
[arrays]
	string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetIntArray("arrays.string")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.string': not an array of integers")
}

// GetIntArray(): local and global options.
func (ct *ConfigTests) TestGetIntArray7(c *C) {
	contents := `
integers = [1, 2]
[arrays]
	integers = [3, 4]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0, err0 := ct.config.GetIntArray("integers")
	c.Check(err0, IsNil)
	c.Check(value0, EqualSlice, []int64{1, 2})
	value1, err1 := ct.config.GetIntArray("arrays.integers")
	c.Check(err1, IsNil)
	c.Check(value1, EqualSlice, []int64{3, 4})
}

// GetIntArrayDefault(): nil structure.
func (ct *ConfigTests) TestGetIntArrayDefault1(c *C) {
	var cfg *config.Configuration

	arr := []int64{9, 8}
	value := cfg.GetIntArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetIntArrayDefault(): unknown option.
func (ct *ConfigTests) TestGetIntArrayDefault2(c *C) {
	contents := `
[values]
	integers = [17, 93]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []int64{9, -8}
	value := ct.config.GetIntArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetIntArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetIntArrayDefault3(c *C) {
	contents := `
[arrays]
	booleans = [true, false]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []int64{9, 8}
	value := ct.config.GetIntArrayDefault("arrays.booleans", arr)
	c.Check(value, EqualSlice, arr)
}

// GetIntArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetIntArrayDefault4(c *C) {
	contents := `
[arrays]
	fp = [3.1415, 5.1413]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []int64{9, 8}
	value := ct.config.GetIntArrayDefault("arrays.fp", arr)
	c.Check(value, EqualSlice, arr)
}

// GetIntArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetIntArrayDefault5(c *C) {
	contents := `
[arrays]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []int64{9, 8}
	value := ct.config.GetIntArrayDefault("arrays.date", arr)
	c.Check(value, EqualSlice, arr)
}

// GetIntArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetIntArrayDefault6(c *C) {
	contents := `
[arrays]
	string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []int64{9, -8}
	value := ct.config.GetIntArrayDefault("arrays.string", arr)
	c.Check(value, EqualSlice, arr)
}

// GetIntArrayDefault(): local and global options.
func (ct *ConfigTests) TestGetIntArrayDefault7(c *C) {
	contents := `
integers = [1, -2]
[arrays]
	integers = [-3, 4]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	arr0 := []int64{9, -8}
	expected0 := []int64{1, -2}
	value0 := ct.config.GetIntArrayDefault("integers", arr0)
	c.Check(value0, EqualSlice, expected0)
	arr1 := []int64{9, -8}
	expected1 := []int64{-3, 4}
	value1 := ct.config.GetIntArrayDefault("arrays.integers", arr1)
	c.Check(value1, EqualSlice, expected1)
}

// GetFloatArray(): nil structure.
func (ct *ConfigTests) TestGetFloatArray1(c *C) {
	var cfg *config.Configuration

	_, err := cfg.GetFloatArray("foo")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'foo': unknown option")
}

// GetFloatArray(): unknown option.
func (ct *ConfigTests) TestGetFloatArray2(c *C) {
	contents := `
[values]
	string = ["foo", "bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err := ct.config.GetFloatArray("v.string")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.string': unknown option")
}

// GetFloatArray(): wrong type.
func (ct *ConfigTests) TestGetFloatArray3(c *C) {
	contents := `
[arrays]
	boolean = [false, true]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetFloatArray("arrays.boolean")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches,
		"'arrays.boolean': not an array of floating-point numbers")
}

// GetFloatArray(): wrong type.
func (ct *ConfigTests) TestGetFloatArray4(c *C) {
	contents := `
[arrays]
	integer = [12, 34]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetFloatArray("arrays.integer")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches,
		"'arrays.integer': not an array of floating-point numbers")
}

// GetFloatArray(): wrong type.
func (ct *ConfigTests) TestGetFloatArray5(c *C) {
	contents := `
[arrays]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetFloatArray("arrays.date")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches,
		"'arrays.date': not an array of floating-point numbers")
}

// GetFloatArray(): wrong type.
func (ct *ConfigTests) TestGetFloatArray6(c *C) {
	contents := `
[arrays]
	string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetFloatArray("arrays.string")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches,
		"'arrays.string': not an array of floating-point numbers")
}

// GetFloatArray(): local and global options.
func (ct *ConfigTests) TestGetFloatArray7(c *C) {
	contents := `
fps = [1.1, 2.2]
[arrays]
	fps = [3.3, 4.4]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0, err0 := ct.config.GetFloatArray("fps")
	c.Check(err0, IsNil)
	c.Check(value0, EqualSlice, []float64{1.1, 2.2})
	value1, err1 := ct.config.GetFloatArray("arrays.fps")
	c.Check(err1, IsNil)
	c.Check(value1, EqualSlice, []float64{3.3, 4.4})
}

// GetFloatArrayDefault(): nil structure.
func (ct *ConfigTests) TestGetFloatArrayDefault1(c *C) {
	var cfg *config.Configuration

	arr := []float64{3.1415}
	value := cfg.GetFloatArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetFloatArrayDefault(): unknown option.
func (ct *ConfigTests) TestGetFloatArrayDefault2(c *C) {
	contents := `
[values]
	string = ["foo", "bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []float64{3.1415}
	value := ct.config.GetFloatArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetFloatArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetFloatArrayDefault3(c *C) {
	contents := `
[arrays]
	boolean = [false, true]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []float64{3.1415}
	value := ct.config.GetFloatArrayDefault("arrays.boolean", arr)
	c.Check(value, EqualSlice, arr)
}

// GetFloatArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetFloatArrayDefault4(c *C) {
	contents := `
[arrays]
	integer = [12, 34]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []float64{3.1415}
	value := ct.config.GetFloatArrayDefault("arrays.boolean", arr)
	c.Check(value, EqualSlice, arr)
}

// GetFloatArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetFloatArrayDefault5(c *C) {
	contents := `
[arrays]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []float64{3.1415}
	value := ct.config.GetFloatArrayDefault("arrays.boolean", arr)
	c.Check(value, EqualSlice, arr)
}

// GetFloatArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetFloatArrayDefault6(c *C) {
	contents := `
[arrays]
	string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []float64{-3.1415}
	value := ct.config.GetFloatArrayDefault("arrays.boolean", arr)
	c.Check(value, EqualSlice, arr)
}

// GetFloatArrayDefault(): local and global options.
func (ct *ConfigTests) TestGetFloatArrayDefault7(c *C) {
	contents := `
fps = [-1.1, +2.2]
[arrays]
	fps = [+3.3, +4.4]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	arr0 := []float64{-3.1415}
	expected0 := []float64{-1.1, +2.2}
	value0 := ct.config.GetFloatArrayDefault("fps", arr0)
	c.Check(value0, EqualSlice, expected0)
	arr1 := []float64{3.1415}
	expected1 := []float64{3.3, 4.4}
	value1 := ct.config.GetFloatArrayDefault("arrays.fps", arr1)
	c.Check(value1, EqualSlice, expected1)
}

// GetDateArray(): nil structure.
func (ct *ConfigTests) TestGetDateArray1(c *C) {
	var cfg *config.Configuration

	_, err := cfg.GetDateArray("foo")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'foo': unknown option")
}

// GetDateArray(): unknown option.
func (ct *ConfigTests) TestGetDateArray2(c *C) {
	contents := `
[values]
	string = ["foo", "bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err := ct.config.GetDateArray("v.string")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.string': unknown option")
}

// GetDateArray(): wrong type.
func (ct *ConfigTests) TestGetDateArray3(c *C) {
	contents := `
[arrays]
	boolean = [false, true]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetDateArray("arrays.boolean")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.boolean': not an array of dates")
}

// GetDateArray(): wrong type.
func (ct *ConfigTests) TestGetDateArray4(c *C) {
	contents := `
[arrays]
	integer = [12, 34]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetDateArray("arrays.integer")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.integer': not an array of dates")
}

// GetDateArray(): wrong type.
func (ct *ConfigTests) TestGetDateArray5(c *C) {
	contents := `
[arrays]
	fp = [3.1415, 5.1413]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetDateArray("arrays.fp")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.fp': not an array of dates")
}

// GetDateArray(): wrong type.
func (ct *ConfigTests) TestGetDateArray6(c *C) {
	contents := `
[arrays]
	string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetDateArray("arrays.string")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.string': not an array of dates")
}

// GetDateArray(): local and global options.
func (ct *ConfigTests) TestGetDateArray7(c *C) {
	contents := `
dates = [2012-10-24T16:22:00Z, 2013-10-24T16:22:00Z]
[arrays]
	dates = [2013-10-24T16:22:00Z, 2014-10-24T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	expected0, _ := time.Parse(time.RFC3339, "2012-10-24T16:22:00Z")
	expected1, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	expected2, _ := time.Parse(time.RFC3339, "2014-10-24T16:22:00Z")

	value0, err0 := ct.config.GetDateArray("dates")
	c.Check(err0, IsNil)
	c.Check(value0, EqualSlice, []time.Time{expected0, expected1})
	value1, err1 := ct.config.GetDateArray("arrays.dates")
	c.Check(err1, IsNil)
	c.Check(value1, EqualSlice, []time.Time{expected1, expected2})
}

// GetDateArrayDefault(): nil structure.
func (ct *ConfigTests) TestGetDateArrayDefault1(c *C) {
	var cfg *config.Configuration

	expected, _ := time.Parse(time.RFC3339, "2012-10-24T16:22:00Z")
	arr := []time.Time{expected}
	value := cfg.GetDateArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetDateArrayDefault(): unknown option.
func (ct *ConfigTests) TestGetDateArrayDefault2(c *C) {
	contents := `
[values]
	string = ["foo", "bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	expected, _ := time.Parse(time.RFC3339, "2012-10-24T16:22:00Z")
	arr := []time.Time{expected}
	value := ct.config.GetDateArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetDateArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetDateArrayDefault3(c *C) {
	contents := `
[arrays]
	boolean = [false, true]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	expected, _ := time.Parse(time.RFC3339, "2012-10-24T16:22:00Z")
	arr := []time.Time{expected}
	value := ct.config.GetDateArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetDateArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetDateArrayDefault4(c *C) {
	contents := `
[arrays]
	integer = [12, 34]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	expected, _ := time.Parse(time.RFC3339, "2012-10-24T16:22:00Z")
	arr := []time.Time{expected}
	value := ct.config.GetDateArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetDateArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetDateArrayDefault5(c *C) {
	contents := `
[arrays]
	fp = [3.1415, 5.1413]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	expected, _ := time.Parse(time.RFC3339, "2012-10-24T16:22:00Z")
	arr := []time.Time{expected}
	value := ct.config.GetDateArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetDateArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetDateArrayDefault6(c *C) {
	contents := `
[arrays]
	string = ["Hello World!", "foo bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	expected, _ := time.Parse(time.RFC3339, "2012-10-24T16:22:00Z")
	arr := []time.Time{expected}
	value := ct.config.GetDateArrayDefault("foo", arr)
	c.Check(value, EqualSlice, arr)
}

// GetDateArrayDefault(): local and global options.
func (ct *ConfigTests) TestGetDateArrayDefault7(c *C) {
	contents := `
dates = [2012-10-24T16:22:00Z, 2013-10-24T16:22:00Z]
[arrays]
	dates = [2013-10-24T16:22:00Z, 2014-10-24T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	expected0, _ := time.Parse(time.RFC3339, "2012-10-24T16:22:00Z")
	expected1, _ := time.Parse(time.RFC3339, "2013-10-24T16:22:00Z")
	expected2, _ := time.Parse(time.RFC3339, "2014-10-24T16:22:00Z")

	arr0 := []time.Time{expected0}
	value0 := ct.config.GetDateArrayDefault("dates", arr0)
	c.Check(value0, EqualSlice, []time.Time{expected0, expected1})
	arr1 := []time.Time{expected1}
	value1 := ct.config.GetDateArrayDefault("arrays.dates", arr1)
	c.Check(value1, EqualSlice, []time.Time{expected1, expected2})
}

// GetStringArray(): nil structure.
func (ct *ConfigTests) TestGetStringArray1(c *C) {
	var cfg *config.Configuration

	_, err := cfg.GetStringArray("foo")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'foo': unknown option")
}

// GetStringArray(): unknown option.
func (ct *ConfigTests) TestGetStringArray2(c *C) {
	contents := `
[values]
	string = ["foo", "bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err := ct.config.GetStringArray("v.string")
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'v.string': unknown option")
}

// GetStringArray(): wrong type.
func (ct *ConfigTests) TestGetStringArray3(c *C) {
	contents := `
[arrays]
	boolean = [false, true]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetStringArray("arrays.boolean")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.boolean': not an array of strings")
}

// GetStringArray(): wrong type.
func (ct *ConfigTests) TestGetStringArray4(c *C) {
	contents := `
[arrays]
	integer = [12, 34]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetStringArray("arrays.integer")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.integer': not an array of strings")
}

// GetStringArray(): wrong type.
func (ct *ConfigTests) TestGetStringArray5(c *C) {
	contents := `
[arrays]
	fp = [3.1415, 5.1413]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetStringArray("arrays.fp")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.fp': not an array of strings")
}

// GetStringArray(): wrong type.
func (ct *ConfigTests) TestGetStringArray6(c *C) {
	contents := `
[arrays]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	_, err0 := ct.config.GetStringArray("arrays.date")
	c.Check(err0, NotNil)
	c.Check(err0, ErrorMatches, "'arrays.date': not an array of strings")
}

// GetStringArray(): local and global options.
func (ct *ConfigTests) TestGetStringArray7(c *C) {
	contents := `
strings = ["foo", "bar"]
[arrays]
	strings = ["bar", "foo"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	value0, err0 := ct.config.GetStringArray("strings")
	c.Check(err0, IsNil)
	c.Check(value0, EqualSlice, []string{"foo", "bar"})
	value1, err1 := ct.config.GetStringArray("arrays.strings")
	c.Check(err1, IsNil)
	c.Check(value1, EqualSlice, []string{"bar", "foo"})
}

// GetStringArrayDefault(): nil structure.
func (ct *ConfigTests) TestGetStringArrayDefault1(c *C) {
	var cfg *config.Configuration

	arr := []string{"foo"}
	value := cfg.GetStringArrayDefault("bar", arr)
	c.Check(value, EqualSlice, arr)
}

// GetStringArrayDefault(): unknown option.
func (ct *ConfigTests) TestGetStringArrayDefault2(c *C) {
	contents := `
[values]
	string = ["foo", "bar"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []string{"foo"}
	value := ct.config.GetStringArrayDefault("bar", arr)
	c.Check(value, EqualSlice, arr)
}

// GetStringArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetStringArrayDefault3(c *C) {
	contents := `
[arrays]
	boolean = [false, true]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []string{"foo"}
	value := ct.config.GetStringArrayDefault("bar", arr)
	c.Check(value, EqualSlice, arr)
}

// GetStringArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetStringArrayDefault4(c *C) {
	contents := `
[arrays]
	integer = [12, 34]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []string{"foo"}
	value := ct.config.GetStringArrayDefault("bar", arr)
	c.Check(value, EqualSlice, arr)
}

// GetStringArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetStringArrayDefault5(c *C) {
	contents := `
[arrays]
	fp = [3.1415, 5.1413]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []string{"foo"}
	value := ct.config.GetStringArrayDefault("bar", arr)
	c.Check(value, EqualSlice, arr)
}

// GetStringArrayDefault(): wrong type.
func (ct *ConfigTests) TestGetStringArrayDefault6(c *C) {
	contents := `
[arrays]
	date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	arr := []string{"foo"}
	value := ct.config.GetStringArrayDefault("bar", arr)
	c.Check(value, EqualSlice, arr)
}

// GetStringArrayDefault(): local and global options.
func (ct *ConfigTests) TestGetStringArrayDefault7(c *C) {
	contents := `
strings = ["foo", "bar"]
[arrays]
	strings = ["bar", "foo"]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	arr0 := []string{"foo"}
	value0 := ct.config.GetStringArrayDefault("strings", arr0)
	c.Check(value0, EqualSlice, []string{"foo", "bar"})
	arr1 := []string{"bar"}
	value1 := ct.config.GetStringArrayDefault("arrays.strings", arr1)
	c.Check(value1, EqualSlice, []string{"bar", "foo"})
}

// String(): Dump function.
func (ct *ConfigTests) TestGetDump1(c *C) {
	contents := `boolean = true
integer = 1
fp = 2.300000
date = 2013-10-25T16:22:00Z
string = "foo"
booleans = [true, false]`

	ct.createTestEnv(c, []string{contents})
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 6)

	str := fmt.Sprintf("%s", ct.config)
	c.Check(str, Equals, contents+"\n")
}
