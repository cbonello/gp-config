package config_test

import (
	"github.com/cbonello/gp-config"
	. "launchpad.net/gocheck"
	"time"
)

type (
	DecodeTests struct {
		config *config.Configuration
	}
)

var (
	_ = Suite(&DecodeTests{})
)

func (ct *DecodeTests) cleanTestEnv(c *C) {
	ct.config = nil
}

// Decode(): nil structure.
func (ct *DecodeTests) TestDecode1(c *C) {
	var cfg *config.Configuration = nil

	err := cfg.Decode("", nil)
	c.Check(err, IsNil)
}

// Decode(): undefined section.
func (ct *DecodeTests) TestDecode2(c *C) {
	contents := `foo = "bar"`

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	err := ct.config.Decode("values", nil)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "'values': unknown section")
}

// Decode(): invalid 2nd argument; nil value.
func (ct *DecodeTests) TestDecode3(c *C) {
	contents := `
[values]
	boolean = false`

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	err := ct.config.Decode("values", nil)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "structure argument cannot be a nil value")
}

// Decode(): invalid 2nd argument; not a pointer.
func (ct *DecodeTests) TestDecode4(c *C) {
	contents := `
[values]
	boolean = false`

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	s := 1
	err := ct.config.Decode("values", s)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "structure argument is not a pointer")
}

// Decode(): invalid 2nd argument; nil pointer to a structure.
func (ct *DecodeTests) TestDecode5(c *C) {
	contents := `
[values]
	boolean = false`

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	type foo struct{}
	var bar *foo = nil

	err := ct.config.Decode("values", bar)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "structure argument cannot be a nil pointer")
}

// Decode(): invalid 2nd argument; not a pointer to a structure.
func (ct *DecodeTests) TestDecode6(c *C) {
	contents := `
[values]
	boolean = false`

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	s := 1
	err := ct.config.Decode("values", &s)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches, "structure argument is not a pointer to a structure")
}

// Decode(): support for embedded struct fields.
func (ct *DecodeTests) TestDecode7(c *C) {
	contents := `
[values]
	integer = 56
	fp = -3.14`

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	type (
		foo struct {
			Integer int64
		}
		bar struct {
			foo
			Fp float64
		}
	)
	var (
		b bar = bar{
			Fp: 314,
		}
	)

	b.foo.Integer = 12
	err := ct.config.Decode("values", &b)
	c.Check(err, IsNil)
	c.Check(b.foo.Integer, Equals, int64(56))
	c.Check(b.Fp, Equals, -3.14)
}

// Decode(): support for multiple embedded struct fields.
func (ct *DecodeTests) TestDecode8(c *C) {
	contents := `
[values]
	boolean = true
	integer = 56
	fp = -3.14`

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 3)

	type (
		foo0 struct {
			Boolean bool
		}
		foo1 struct {
			foo0
			Integer int64
		}
		bar struct {
			foo1
			Fp float64
		}
	)
	var (
		b bar = bar{
			Fp: 314,
		}
	)

	b.foo1.foo0.Boolean = false
	b.foo1.Integer = 78
	b.Fp = 4.5
	err := ct.config.Decode("values", &b)
	c.Check(err, IsNil)
	c.Check(b.foo1.foo0.Boolean, Equals, true)
	c.Check(b.foo1.Integer, Equals, int64(56))
	c.Check(b.Fp, Equals, -3.14)
}

// Decode(): support for multiple embedded struct fields.
func (ct *DecodeTests) TestDecode9(c *C) {
	contents := `
[values]
	boolean = true
	integer = 56
	fp = -3.14`

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 3)

	type (
		foo0 struct {
			Boolean bool
		}
		foo1 struct {
			Integer int64
			foo0
			dummy int
		}
		bar struct {
			Fp float64
			foo1
		}
	)
	var (
		b bar = bar{
			Fp: 314,
		}
	)

	b.foo1.foo0.Boolean = false
	b.foo1.Integer = 78
	b.Fp = 4.5
	err := ct.config.Decode("values", &b)
	c.Check(err, IsNil)
	c.Check(b.foo1.foo0.Boolean, Equals, true)
	c.Check(b.foo1.Integer, Equals, int64(56))
	c.Check(b.Fp, Equals, -3.14)
}

// Decode(): check for anonymous structure fields.
func (ct *DecodeTests) TestDecode10(c *C) {
	contents := `
[values]
	boolean = false`

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	type foo struct {
		_       int64
		Boolean bool
	}

	bar := foo{}
	err := ct.config.Decode("values", &bar)
	c.Check(err, IsNil)
	c.Check(bar.Boolean, Equals, false)
}

// Decode(): decoding of global options.
func (ct *DecodeTests) TestDecode11(c *C) {
	contents := `
boolean = false
integer = 12
fp = 3.1415
date = 2013-10-25T16:22:00Z
string = "Hello World!"

booleans = [false, true]
integers = [12, 34]
fps = [3.1415, 5.1413]
dates = [2012-10-25T16:22:00Z, 2013-10-25T16:22:00Z]
strings = ["Hello World!", "foo bar"]`

	type (
		values struct {
			Boolean bool
			Integer int64
			Fp      float64
			Date    time.Time
			String  string

			Booleans []bool
			Integers []int64
			Fps      []float64
			Dates    []time.Time
			Strings  []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 10)

	v := values{}
	expected0, _ := time.Parse(time.RFC3339, "2012-10-25T16:22:00Z")
	expected1, _ := time.Parse(time.RFC3339, "2013-10-25T16:22:00Z")
	err1 := ct.config.Decode("", &v)
	c.Check(err1, IsNil)
	c.Check(v.Boolean, Equals, false)
	c.Check(v.Integer, Equals, int64(12))
	c.Check(v.Fp, Equals, float64(3.1415))
	c.Check(v.Date, Equals, expected1)
	c.Check(v.String, Equals, "Hello World!")
	c.Check(v.Booleans, EqualSlice, []bool{false, true})
	c.Check(v.Integers, EqualSlice, []int64{12, 34})
	c.Check(v.Fps, EqualSlice, []float64{3.1415, 5.1413})
	c.Check(v.Dates, EqualSlice, []time.Time{expected0, expected1})
	c.Check(v.Strings, EqualSlice, []string{"Hello World!", "foo bar"})
}

// Decode(): decoding of local options.
func (ct *DecodeTests) TestDecode12(c *C) {
	contents := `
[values]
	boolean = false
	integer = 12
	fp = 3.1415
	date = 2012-10-25T16:22:00Z
	string = "Hello World!"
[arrays]
	booleans = [false, true]
	integers = [12, 34]
	fps = [3.1415, 5.1413]
	dates = [2012-10-25T16:22:00Z, 2013-10-25T16:22:00Z]
	strings = ["Hello World!", "foo bar"]`

	type (
		values struct {
			Boolean bool
			Integer int64
			Fp      float64
			Date    time.Time
			String  string
		}
		arrays struct {
			Booleans []bool
			Integers []int64
			Fps      []float64
			Dates    []time.Time
			Strings  []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 10)

	v := values{}
	expected0, _ := time.Parse(time.RFC3339, "2012-10-25T16:22:00Z")
	err1 := ct.config.Decode("values", &v)
	c.Check(err1, IsNil)
	c.Check(v.Boolean, Equals, false)
	c.Check(v.Integer, Equals, int64(12))
	c.Check(v.Fp, Equals, float64(3.1415))
	c.Check(v.Date, Equals, expected0)
	c.Check(v.String, Equals, "Hello World!")

	a := arrays{}
	expected1, _ := time.Parse(time.RFC3339, "2013-10-25T16:22:00Z")
	err2 := ct.config.Decode("arrays", &a)
	c.Check(err2, IsNil)
	c.Check(a.Booleans, EqualSlice, []bool{false, true})
	c.Check(a.Integers, EqualSlice, []int64{12, 34})
	c.Check(a.Fps, EqualSlice, []float64{3.1415, 5.1413})
	c.Check(a.Dates, EqualSlice, []time.Time{expected0, expected1})
	c.Check(a.Strings, EqualSlice, []string{"Hello World!", "foo bar"})
}

// Decode(): wrong field type.
func (ct *DecodeTests) TestDecode13(c *C) {
	contents := `port = 8080`

	type (
		values struct {
			Port *int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Port': value of type ptr is not assignable to type int64")
}

// Decode(): wrong field type.
func (ct *DecodeTests) TestDecode14(c *C) {
	contents := `port = 8080`

	type (
		values struct {
			Port interface{}
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Port': value of type interface is not assignable to type int64")
}

// Decode(): wrong field type.
func (ct *DecodeTests) TestDecode15(c *C) {
	contents := `port = 8080`

	type (
		otherValues struct {
			Port int64
		}
		values struct {
			*otherValues
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	v.otherValues = &otherValues{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'otherValues': embedded pointer fields are not yet supported!")
}

// Decode(): StructTags.
func (ct *DecodeTests) TestDecode16(c *C) {
	contents := `
addr = "localhost"
port = 8080`

	type (
		server struct {
			URL  string `option:"addr"`
			Port int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	s := server{}
	err := ct.config.Decode("", &s)
	c.Check(err, IsNil)
	c.Check(s.URL, Equals, "localhost")
	c.Check(s.Port, Equals, int64(8080))
}

// Decode(): StructTags.
func (ct *DecodeTests) TestDecode17(c *C) {
	contents := `
addr = "localhost"
port = 8080`

	type (
		server struct {
			addr string
			Port int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 2)

	s := server{}
	err := ct.config.Decode("", &s)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'addr': cannot set value of unexported struct field")
}

// Decode(): wrong boolean type.
func (ct *DecodeTests) TestDecodeBool1(c *C) {
	contents := `boolean = 9`

	type (
		values struct {
			Boolean bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type int64 is not assignable to type bool")
}

// Decode(): wrong boolean type.
func (ct *DecodeTests) TestDecodeBool2(c *C) {
	contents := `boolean = 9.8`

	type (
		values struct {
			Boolean bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type float64 is not assignable to type bool")
}

// Decode(): wrong boolean type.
func (ct *DecodeTests) TestDecodeBool3(c *C) {
	contents := `boolean = 2013-10-25T16:22:00Z`

	type (
		values struct {
			Boolean bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type time.Time is not assignable to type bool")
}

// Decode(): wrong boolean type.
func (ct *DecodeTests) TestDecodeBool4(c *C) {
	contents := `boolean = "2013-10-25T16:22:00Z"`

	type (
		values struct {
			Boolean bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type string is not assignable to type bool")
}

// Decode(): wrong boolean type.
func (ct *DecodeTests) TestDecodeBool5(c *C) {
	contents := `boolean = [true, false]`

	type (
		values struct {
			Boolean bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type \\[\\]bool is not assignable to type bool")
}

// Decode(): wrong boolean type.
func (ct *DecodeTests) TestDecodeBool6(c *C) {
	contents := `boolean = [9, 8]`

	type (
		values struct {
			Boolean bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type \\[\\]int64 is not assignable to type bool")
}

// Decode(): wrong boolean type.
func (ct *DecodeTests) TestDecodeBool7(c *C) {
	contents := `boolean = [9.8, 8.9]`

	type (
		values struct {
			Boolean bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type \\[\\]float64 is not assignable to type bool")
}

// Decode(): wrong boolean type.
func (ct *DecodeTests) TestDecodeBool8(c *C) {
	contents := `boolean = [2013-10-25T16:22:00Z, 2014-10-25T16:22:00Z]`

	type (
		values struct {
			Boolean bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type \\[\\]time.Time is not assignable to type bool")
}

// Decode(): wrong boolean type.
func (ct *DecodeTests) TestDecodeBool9(c *C) {
	contents := `boolean = ["foo", "bar"]`

	type (
		values struct {
			Boolean bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type \\[\\]string is not assignable to type bool")
}

// Decode(): wrong integer type.
func (ct *DecodeTests) TestDecodeInt1(c *C) {
	contents := `integer = false`

	type (
		values struct {
			Integer int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type bool is not assignable to type int64")
}

// Decode(): wrong integer type.
func (ct *DecodeTests) TestDecodeInt2(c *C) {
	contents := `integer = 9.8`

	type (
		values struct {
			Integer int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type float64 is not assignable to type int64")
}

// Decode(): wrong integer type.
func (ct *DecodeTests) TestDecodeInt3(c *C) {
	contents := `integer = 2013-10-25T16:22:00Z`

	type (
		values struct {
			Integer int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type time.Time is not assignable to type int64")
}

// Decode(): wrong integer type.
func (ct *DecodeTests) TestDecodeInt4(c *C) {
	contents := `integer = "2013-10-25T16:22:00Z"`

	type (
		values struct {
			Integer int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type string is not assignable to type int64")
}

// Decode(): wrong integer type.
func (ct *DecodeTests) TestDecodeInt5(c *C) {
	contents := `integer = [false, false]`

	type (
		values struct {
			Integer int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type \\[\\]bool is not assignable to type int64")
}

// Decode(): wrong integer type.
func (ct *DecodeTests) TestDecodeInt6(c *C) {
	contents := `integer = [55, 30]`

	type (
		values struct {
			Integer int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type \\[\\]int64 is not assignable to type int64")
}

// Decode(): wrong integer type.
func (ct *DecodeTests) TestDecodeInt7(c *C) {
	contents := `integer = [9.8, 8.9]`

	type (
		values struct {
			Integer int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type \\[\\]float64 is not assignable to type int64")
}

// Decode(): wrong integer type.
func (ct *DecodeTests) TestDecodeInt8(c *C) {
	contents := `integer = [2012-10-25T16:22:00Z, 2013-10-25T16:22:00Z]`

	type (
		values struct {
			Integer int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type \\[\\]time.Time is not assignable to type int64")
}

// Decode(): wrong integer type.
func (ct *DecodeTests) TestDecodeInt9(c *C) {
	contents := `integer = ["foo", "bar"]`

	type (
		values struct {
			Integer int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type \\[\\]string is not assignable to type int64")
}

// Decode(): wrong floating-point type.
func (ct *DecodeTests) TestDecodeFp1(c *C) {
	contents := `fp = false`

	type (
		values struct {
			Fp float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type bool is not assignable to type float64")
}

// Decode(): wrong floating-point type.
func (ct *DecodeTests) TestDecodeFp2(c *C) {
	contents := `fp = 9`

	type (
		values struct {
			Fp float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type int64 is not assignable to type float64")
}

// Decode(): wrong floating-point type.
func (ct *DecodeTests) TestDecodeFp3(c *C) {
	contents := `fp = 2013-10-25T16:22:00Z`

	type (
		values struct {
			Fp float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type time.Time is not assignable to type float64")
}

// Decode(): wrong floating-point type.
func (ct *DecodeTests) TestDecodeFp4(c *C) {
	contents := `fp = "2013-10-25T16:22:00Z"`

	type (
		values struct {
			Fp float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type string is not assignable to type float64")
}

// Decode(): wrong floating-point type.
func (ct *DecodeTests) TestDecodeFp5(c *C) {
	contents := `fp = [true, false]`

	type (
		values struct {
			Fp float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type \\[\\]bool is not assignable to type float64")
}

// Decode(): wrong floating-point type.
func (ct *DecodeTests) TestDecodeFp6(c *C) {
	contents := `fp = [9, 2]`

	type (
		values struct {
			Fp float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type \\[\\]int64 is not assignable to type float64")
}

// Decode(): wrong floating-point type.
func (ct *DecodeTests) TestDecodeFp7(c *C) {
	contents := `fp = [9.8, 7.2]`

	type (
		values struct {
			Fp float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type \\[\\]float64 is not assignable to type float64")
}

// Decode(): wrong floating-point type.
func (ct *DecodeTests) TestDecodeFp8(c *C) {
	contents := `fp = [2013-10-25T16:22:00Z, 2000-06-02T02:00:00Z]`

	type (
		values struct {
			Fp float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type \\[\\]time.Time is not assignable to type float64")
}

// Decode(): wrong floating-point type.
func (ct *DecodeTests) TestDecodeFp9(c *C) {
	contents := `fp = ["A", "L"]`

	type (
		values struct {
			Fp float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type \\[\\]string is not assignable to type float64")
}

// Decode(): wrong date type.
func (ct *DecodeTests) TestDecodeDate1(c *C) {
	contents := `date = true`

	type (
		values struct {
			Date time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type bool is not assignable to type time.Time")
}

// Decode(): wrong date type.
func (ct *DecodeTests) TestDecodeDate2(c *C) {
	contents := `date = 87`

	type (
		values struct {
			Date time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type int64 is not assignable to type time.Time")
}

// Decode(): wrong date type.
func (ct *DecodeTests) TestDecodeDate3(c *C) {
	contents := `date = 9.8`

	type (
		values struct {
			Date time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type float64 is not assignable to type time.Time")
}

// Decode(): wrong date type.
func (ct *DecodeTests) TestDecodeDate4(c *C) {
	contents := `date = "2013-10-25T16:22:00Z"`

	type (
		values struct {
			Date time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type string is not assignable to type time.Time")
}

// Decode(): wrong date type.
func (ct *DecodeTests) TestDecodeDate5(c *C) {
	contents := `date = [true, true]`

	type (
		values struct {
			Date time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type \\[\\]bool is not assignable to type time.Time")
}

// Decode(): wrong date type.
func (ct *DecodeTests) TestDecodeDate6(c *C) {
	contents := `date = [1, 87]`

	type (
		values struct {
			Date time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type \\[\\]int64 is not assignable to type time.Time")
}

// Decode(): wrong date type.
func (ct *DecodeTests) TestDecodeDate7(c *C) {
	contents := `date = [9.8, 5.63]`

	type (
		values struct {
			Date time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type \\[\\]float64 is not assignable to type time.Time")
}

// Decode(): wrong date type.
func (ct *DecodeTests) TestDecodeDate8(c *C) {
	contents := `date = [2013-10-24T16:22:00Z, 2013-10-25T16:22:00Z]`

	type (
		values struct {
			Date time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type \\[\\]time.Time is not assignable to type time.Time")
}

// Decode(): wrong date type.
func (ct *DecodeTests) TestDecodeDate9(c *C) {
	contents := `date = ["a", "bcd"]`

	type (
		values struct {
			Date time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type \\[\\]string is not assignable to type time.Time")
}

// Decode(): wrong string type.
func (ct *DecodeTests) TestDecodeString1(c *C) {
	contents := `string = true`

	type (
		values struct {
			String string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type bool is not assignable to type string")
}

// Decode(): wrong string type.
func (ct *DecodeTests) TestDecodeString2(c *C) {
	contents := `string = 87`

	type (
		values struct {
			String string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type int64 is not assignable to type string")
}

// Decode(): wrong string type.
func (ct *DecodeTests) TestDecodeString3(c *C) {
	contents := `string = 9.8`

	type (
		values struct {
			String string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type float64 is not assignable to type string")
}

// Decode(): wrong string type.
func (ct *DecodeTests) TestDecodeString4(c *C) {
	contents := `string = 2013-10-25T16:22:00Z`

	type (
		values struct {
			String string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type time.Time is not assignable to type string")
}

// Decode(): wrong string type.
func (ct *DecodeTests) TestDecodeString5(c *C) {
	contents := `string = [false, true]`

	type (
		values struct {
			String string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type \\[\\]bool is not assignable to type string")
}

// Decode(): wrong string type.
func (ct *DecodeTests) TestDecodeString6(c *C) {
	contents := `string = [12, 34]`

	type (
		values struct {
			String string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type \\[\\]int64 is not assignable to type string")
}

// Decode(): wrong string type.
func (ct *DecodeTests) TestDecodeString7(c *C) {
	contents := `string = [3.2, 9.8]`

	type (
		values struct {
			String string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type \\[\\]float64 is not assignable to type string")
}

// Decode(): wrong string type.
func (ct *DecodeTests) TestDecodeString8(c *C) {
	contents := `string = [2013-10-25T16:22:00Z, 2016-10-25T16:22:00Z]`

	type (
		values struct {
			String string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type \\[\\]time.Time is not assignable to type string")
}

// Decode(): wrong string type.
func (ct *DecodeTests) TestDecodeString9(c *C) {
	contents := `string = ["abcd", "efgh"]`

	type (
		values struct {
			String string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type \\[\\]string is not assignable to type string")
}

// Decode(): wrong boolean array type.
func (ct *DecodeTests) TestDecodeBoolArray1(c *C) {
	contents := `boolean = true`

	type (
		values struct {
			Boolean []bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type bool is not assignable to type \\[\\]bool")
}

// Decode(): wrong boolean array type.
func (ct *DecodeTests) TestDecodeBoolArray2(c *C) {
	contents := `boolean = 9`

	type (
		values struct {
			Boolean []bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type int64 is not assignable to type \\[\\]bool")
}

// Decode(): wrong boolean array type.
func (ct *DecodeTests) TestDecodeBoolArray3(c *C) {
	contents := `boolean = 9.8`

	type (
		values struct {
			Boolean []bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type float64 is not assignable to type \\[\\]bool")
}

// Decode(): wrong boolean array type.
func (ct *DecodeTests) TestDecodeBoolArray4(c *C) {
	contents := `boolean = 2013-10-25T16:22:00Z`

	type (
		values struct {
			Boolean []bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type time.Time is not assignable to type \\[\\]bool")
}

// Decode(): wrong boolean array type.
func (ct *DecodeTests) TestDecodeBoolArray5(c *C) {
	contents := `boolean = "foo"`

	type (
		values struct {
			Boolean []bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type string is not assignable to type \\[\\]bool")
}

// Decode(): wrong boolean array type.
func (ct *DecodeTests) TestDecodeBoolArray6(c *C) {
	contents := `boolean = [1, 2, 3, 4]`

	type (
		values struct {
			Boolean []bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type \\[\\]int64 is not assignable to type \\[\\]bool")
}

// Decode(): wrong boolean array type.
func (ct *DecodeTests) TestDecodeBoolArray7(c *C) {
	contents := `boolean = [9.8, 8.9]`

	type (
		values struct {
			Boolean []bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type \\[\\]float64 is not assignable to type \\[\\]bool")
}

// Decode(): wrong boolean array type.
func (ct *DecodeTests) TestDecodeBoolArray8(c *C) {
	contents := `boolean = [2013-10-25T16:22:00Z, 2014-10-25T16:22:00Z]`

	type (
		values struct {
			Boolean []bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type \\[\\]time.Time is not assignable to type \\[\\]bool")
}

// Decode(): wrong boolean array type.
func (ct *DecodeTests) TestDecodeBoolArray9(c *C) {
	contents := `boolean = ["foo", "bar"]`

	type (
		values struct {
			Boolean []bool
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Boolean': value of type \\[\\]string is not assignable to type \\[\\]bool")
}

// Decode(): wrong integer array type.
func (ct *DecodeTests) TestDecodeIntArray1(c *C) {
	contents := `integer = true`

	type (
		values struct {
			Integer []int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type bool is not assignable to type \\[\\]int64")
}

// Decode(): wrong integer array type.
func (ct *DecodeTests) TestDecodeIntArray2(c *C) {
	contents := `integer = 9`

	type (
		values struct {
			Integer []int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type int64 is not assignable to type \\[\\]int64")
}

// Decode(): wrong integer array type.
func (ct *DecodeTests) TestDecodeIntArray3(c *C) {
	contents := `integer = 9.8`

	type (
		values struct {
			Integer []int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type float64 is not assignable to type \\[\\]int64")
}

// Decode(): wrong integer array type.
func (ct *DecodeTests) TestDecodeIntArray4(c *C) {
	contents := `integer = 2013-10-25T16:22:00Z`

	type (
		values struct {
			Integer []int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type time.Time is not assignable to type \\[\\]int64")
}

// Decode(): wrong integer array type.
func (ct *DecodeTests) TestDecodeIntArray5(c *C) {
	contents := `integer = "foo"`

	type (
		values struct {
			Integer []int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type string is not assignable to type \\[\\]int64")
}

// Decode(): wrong integer array type.
func (ct *DecodeTests) TestDecodeIntArray6(c *C) {
	contents := `integer = [true, false, false]`

	type (
		values struct {
			Integer []int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type \\[\\]bool is not assignable to type \\[\\]int64")
}

// Decode(): wrong integer array type.
func (ct *DecodeTests) TestDecodeIntArray7(c *C) {
	contents := `integer = [9.8, 8.9]`

	type (
		values struct {
			Integer []int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type \\[\\]float64 is not assignable to type \\[\\]int64")
}

// Decode(): wrong integer array type.
func (ct *DecodeTests) TestDecodeIntArray8(c *C) {
	contents := `integer = [2013-10-25T16:22:00Z, 2014-10-25T16:22:00Z]`

	type (
		values struct {
			Integer []int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type \\[\\]time.Time is not assignable to type \\[\\]int64")
}

// Decode(): wrong integer array type.
func (ct *DecodeTests) TestDecodeIntArray9(c *C) {
	contents := `integer = ["foo", "bar"]`

	type (
		values struct {
			Integer []int64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Integer': value of type \\[\\]string is not assignable to type \\[\\]int64")
}

// Decode(): wrong floating-point array type.
func (ct *DecodeTests) TestDecodeFloatArray1(c *C) {
	contents := `fp = true`

	type (
		values struct {
			Fp []float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type bool is not assignable to type \\[\\]float64")
}

// Decode(): wrong floating-point array type.
func (ct *DecodeTests) TestDecodeFloatArray2(c *C) {
	contents := `fp = 9`

	type (
		values struct {
			Fp []float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type int64 is not assignable to type \\[\\]float64")
}

// Decode(): wrong floating-point array type.
func (ct *DecodeTests) TestDecodeFloatArray3(c *C) {
	contents := `fp = 9.8`

	type (
		values struct {
			Fp []float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type float64 is not assignable to type \\[\\]float64")
}

// Decode(): wrong floating-point array type.
func (ct *DecodeTests) TestDecodeFloatArray4(c *C) {
	contents := `fp = 2013-10-25T16:22:00Z`

	type (
		values struct {
			Fp []float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type time.Time is not assignable to type \\[\\]float64")
}

// Decode(): wrong floating-point array type.
func (ct *DecodeTests) TestDecodeFloatArray5(c *C) {
	contents := `fp = "foo"`

	type (
		values struct {
			Fp []float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type string is not assignable to type \\[\\]float64")
}

// Decode(): wrong floating-point array type.
func (ct *DecodeTests) TestDecodeFloatArray6(c *C) {
	contents := `fp = [true, false, false]`

	type (
		values struct {
			Fp []float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type \\[\\]bool is not assignable to type \\[\\]float64")
}

// Decode(): wrong floating-point array type.
func (ct *DecodeTests) TestDecodeFloatArray7(c *C) {
	contents := `fp = [9, 8, 8, 9]`

	type (
		values struct {
			Fp []float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type \\[\\]int64 is not assignable to type \\[\\]float64")
}

// Decode(): wrong floating-point array type.
func (ct *DecodeTests) TestDecodeFloatArray8(c *C) {
	contents := `fp = [2013-10-25T16:22:00Z, 2014-10-25T16:22:00Z]`

	type (
		values struct {
			Fp []float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type \\[\\]time.Time is not assignable to type \\[\\]float64")
}

// Decode(): wrong floating-point array type.
func (ct *DecodeTests) TestDecodeFloatArray9(c *C) {
	contents := `fp = ["foo", "bar"]`

	type (
		values struct {
			Fp []float64
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Fp': value of type \\[\\]string is not assignable to type \\[\\]float64")
}

// Decode(): wrong time.Time array type.
func (ct *DecodeTests) TestDecodeDateArray1(c *C) {
	contents := `date = true`

	type (
		values struct {
			Date []time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type bool is not assignable to type \\[\\]time.Time")
}

// Decode(): wrong date array type.
func (ct *DecodeTests) TestDecodeDateArray2(c *C) {
	contents := `date = 9`

	type (
		values struct {
			Date []time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type int64 is not assignable to type \\[\\]time.Time")
}

// Decode(): wrong date array type.
func (ct *DecodeTests) TestDecodeDateArray3(c *C) {
	contents := `date = 9.8`

	type (
		values struct {
			Date []time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type float64 is not assignable to type \\[\\]time.Time")
}

// Decode(): wrong date array type.
func (ct *DecodeTests) TestDecodeDateArray4(c *C) {
	contents := `date = 2013-10-25T16:22:00Z`

	type (
		values struct {
			Date []time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type time.Time is not assignable to type \\[\\]time.Time")
}

// Decode(): wrong date array type.
func (ct *DecodeTests) TestDecodeDateArray5(c *C) {
	contents := `date = "foo"`

	type (
		values struct {
			Date []time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type string is not assignable to type \\[\\]time.Time")
}

// Decode(): wrong date array type.
func (ct *DecodeTests) TestDecodeDateArray6(c *C) {
	contents := `date = [true, false, false]`

	type (
		values struct {
			Date []time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type \\[\\]bool is not assignable to type \\[\\]time.Time")
}

// Decode(): wrong floating-point array type.
func (ct *DecodeTests) TestDecodeDateArray7(c *C) {
	contents := `date = [9, 8, 8, 9]`

	type (
		values struct {
			Date []time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type \\[\\]int64 is not assignable to type \\[\\]time.Time")
}

// Decode(): wrong date array type.
func (ct *DecodeTests) TestDecodeDateArray8(c *C) {
	contents := `date = ["X", "Y", "Z"]`

	type (
		values struct {
			Date []time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type \\[\\]string is not assignable to type \\[\\]time.Time")
}

// Decode(): wrong date array type.
func (ct *DecodeTests) TestDecodeDateArray9(c *C) {
	contents := `date = ["foo", "bar"]`

	type (
		values struct {
			Date []time.Time
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'Date': value of type \\[\\]string is not assignable to type \\[\\]time.Time")
}

// Decode(): wrong string array type.
func (ct *DecodeTests) TestDecodeStringArray1(c *C) {
	contents := `string = true`

	type (
		values struct {
			String []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type bool is not assignable to type \\[\\]string")
}

// Decode(): wrong string array type.
func (ct *DecodeTests) TestDecodeStringArray2(c *C) {
	contents := `string = 9`

	type (
		values struct {
			String []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type int64 is not assignable to type \\[\\]string")
}

// Decode(): wrong string array type.
func (ct *DecodeTests) TestDecodeStringArray3(c *C) {
	contents := `string = 9.8`

	type (
		values struct {
			String []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type float64 is not assignable to type \\[\\]string")
}

// Decode(): wrong string array type.
func (ct *DecodeTests) TestDecodeStringArray4(c *C) {
	contents := `string = 2013-10-25T16:22:00Z`

	type (
		values struct {
			String []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type time.Time is not assignable to type \\[\\]string")
}

// Decode(): wrong string array type.
func (ct *DecodeTests) TestDecodeStringArray5(c *C) {
	contents := `string = "foo"`

	type (
		values struct {
			String []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type string is not assignable to type \\[\\]string")
}

// Decode(): wrong string array type.
func (ct *DecodeTests) TestDecodeStringArray6(c *C) {
	contents := `string = [true, false, false]`

	type (
		values struct {
			String []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type \\[\\]bool is not assignable to type \\[\\]string")
}

// Decode(): wrong string array type.
func (ct *DecodeTests) TestDecodeStringArray7(c *C) {
	contents := `string = [2013, 2014, 2015]`

	type (
		values struct {
			String []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type \\[\\]int64 is not assignable to type \\[\\]string")
}

// Decode(): wrong string array type.
func (ct *DecodeTests) TestDecodeStringArray8(c *C) {
	contents := `string = [9.8, 8.9]`

	type (
		values struct {
			String []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type \\[\\]float64 is not assignable to type \\[\\]string")
}

// Decode(): wrong string array type.
func (ct *DecodeTests) TestDecodeStringArray9(c *C) {
	contents := `string = [2013-10-25T16:22:00Z, 2014-10-25T16:22:00Z]`

	type (
		values struct {
			String []string
		}
	)

	ct.config = config.NewConfiguration()
	err0 := ct.config.LoadString(contents)
	c.Check(err0, IsNil)
	defer ct.cleanTestEnv(c)

	c.Check(ct.config.Len(), Equals, 1)

	v := values{}
	err := ct.config.Decode("", &v)
	c.Check(err, NotNil)
	c.Check(err, ErrorMatches,
		"'String': value of type \\[\\]time.Time is not assignable to type \\[\\]string")
}
