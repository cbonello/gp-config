package config

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

type (
	// ConfigurationError records parsing errors.
	ConfigurationError struct {
		Filename     string // Filename.
		Line, Column int    // Line and column.
		msg          string
	}

	configurationSections map[string]struct{}

	configurationOptions map[string]configurationValue

	configurationValue struct {
		ctype configurationType // Data's type.
		value interface{}       // Data.
	}

	configurationType uint8

	// Configuration context.
	Configuration struct {
		sync.RWMutex
		// List of sections declared.
		sections configurationSections
		// List of options declared.
		options configurationOptions
		// Given following configuration:
		//
		// foo = "bar"
		// [values]
		// 	 boolean = false
		// 	 integer = 12
		//
		// sections and options will be set as:
		//
		// sections: ["", values"]
		// options:	 ["foo", values.boolean", "values.integer"]
	}
)

const (
	_ARRAY_TYPE  configurationType = 1
	_BOOL_TYPE                     = 2
	_INT_TYPE                      = 4
	_FLOAT_TYPE                    = 8
	_DATE_TYPE                     = 16
	_STRING_TYPE                   = 32
)

var (
	sliceType  = reflect.Slice
	boolType   = reflect.Bool
	intType    = reflect.Int64
	floatType  = reflect.Float64
	dateType   = reflect.TypeOf((*time.Time)(nil)).Elem()
	stringType = reflect.String
	ptrType    = reflect.Ptr
	structType = reflect.Struct
)

// NewConfiguration creates a new configuration context.
func NewConfiguration() (c *Configuration) {
	return &Configuration{
		sections: configurationSections{},
		options:  configurationOptions{},
	}
}

// LoadFile loads the configuration stored in given file.
func (c *Configuration) LoadFile(filename string) (err *ConfigurationError) {
	p, err0 := NewParser(filename)
	if err0 != nil {
		return &ConfigurationError{
			Filename: filename,
			Line:     0,
			Column:   0,
			msg:      err0.Error(),
		}
	}
	if p != nil {
		if err1 := p.Parse(c); err1 != nil {
			return err1
		}
	}
	return nil
}

// LoadString loads the configuration stored in given string.
func (c *Configuration) LoadString(contents string) (err *ConfigurationError) {
	if p := NewStringParser(contents); p != nil {
		if err = p.Parse(c); err != nil {
			return err
		}
	}
	return nil
}

// Len returns the number of options defined.
func (c *Configuration) Len() (length int) {
	if c != nil {
		c.RLock()
		defer c.RUnlock()
		length = len(c.options)
	}
	return length
}

// HasOption returns true if given option is defined, false otherwise.
func (c *Configuration) HasOption(option string) (result bool) {
	if c != nil {
		c.RLock()
		defer c.RUnlock()
		_, result = c.options[strings.ToLower(option)]
	}
	return result
}

// Sections returns the sections defined sorted in ascending order.
func (c *Configuration) Sections() (sections []string) {
	sections = []string{}
	if c != nil {
		c.RLock()
		defer c.RUnlock()
		for s := range c.sections {
			sections = append(sections, s)
		}
		sort.Strings(sections)
	}
	return sections
}

// IsSection returns true if given section exists, otherwise false.
func (c *Configuration) IsSection(section string) bool {
	if c != nil {
		c.RLock()
		defer c.RUnlock()
		section = strings.ToLower(section)
		for s := range c.sections {
			if section == s {
				return true
			}
		}
	}
	return false
}

// Options returns the options defined in given section sorted in ascending
// order. Optionsglobally  defined are returned if given section is an empty
// string.
func (c *Configuration) Options(section string) (options []string) {
	options = []string{}
	if c != nil {
		c.RLock()
		defer c.RUnlock()
		if section = strings.ToLower(section); section == "" {
			// Look for options defined outside of a section.
			for o := range c.options {
				if strings.Index(o, ".") == -1 {
					options = append(options, o)
				}
			}
		} else {
			for o := range c.options {
				if strings.HasPrefix(o, section) {
					options = append(options, o)
				}
			}
		}
		sort.Strings(options)
	}
	return options
}

// Get returns the value (scalar or array) associated with given option name,
// or nil if option is undefined. option is a dot-separated path (e.g.
// section.option).
func (c *Configuration) Get(option string) (interface{}, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			return opt.value, nil
		}
	}
	return nil, fmt.Errorf("'%s': unknown option", option)
}

// GetBool returns the boolean value associated with given option name. An
// error is flagged if option does not record a boolean value or it is
// undefined.
func (c *Configuration) GetBool(option string) (bool, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _BOOL_TYPE {
				return opt.value.(bool), nil
			}
			return false, fmt.Errorf("'%s': not a boolean", option)
		}
	}
	return false, fmt.Errorf("'%s': unknown option", option)
}

// GetBoolDefault is similar to GetBool but given default value is returned
// if option does not exist or is of wrong type.
func (c *Configuration) GetBoolDefault(option string, dfault bool) bool {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _BOOL_TYPE {
				return opt.value.(bool)
			}
		}
	}
	return dfault
}

// GetInt returns the integer associated with given option name. An error
// is flagged if option does not record an integer or it is undefined.
func (c *Configuration) GetInt(option string) (int64, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _INT_TYPE {
				return opt.value.(int64), nil
			}
			return 0, fmt.Errorf("'%s': not an integer", option)
		}
	}
	return 0, fmt.Errorf("'%s': unknown option", option)
}

// GetIntDefault is similar to GetInt but given default value is returned
// if option does not exist or is of wrong type.
func (c *Configuration) GetIntDefault(option string, dfault int64) int64 {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _INT_TYPE {
				return opt.value.(int64)
			}
		}
	}
	return dfault
}

// GetFloat returns the floating-point number associated with given option
// name. An error is flagged if option does not record a floating-point number
// or it is undefined.
func (c *Configuration) GetFloat(option string) (float64, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _FLOAT_TYPE {
				return opt.value.(float64), nil
			}
			return 0, fmt.Errorf("'%s': not a floating-point number", option)
		}
	}
	return 0, fmt.Errorf("'%s': unknown option", option)
}

// GetFloatDefault is similar to GetFloat but given default value is returned
// if option does not exist or is of wrong type.
func (c *Configuration) GetFloatDefault(option string, dfault float64) float64 {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _FLOAT_TYPE {
				return opt.value.(float64)
			}
		}
	}
	return dfault
}

// GetDate returns the date associated with given option name. An error
// is flagged if option does not record a date or it is undefined.
func (c *Configuration) GetDate(option string) (time.Time, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _DATE_TYPE {
				return opt.value.(time.Time), nil
			}
			return time.Now(), fmt.Errorf("'%s': not a date", option)
		}
	}
	return time.Now(), fmt.Errorf("'%s': unknown option", option)
}

// GetDateDefault is similar to GetDate but given default value is returned
// if option does not exist or is of wrong type.
func (c *Configuration) GetDateDefault(option string, dfault time.Time) time.Time {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _DATE_TYPE {
				return opt.value.(time.Time)
			}
		}
	}
	return dfault
}

// GetString returns the string associated with given option name. An error
// is flagged if option does not record a string or it is undefined.
func (c *Configuration) GetString(option string) (string, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _STRING_TYPE {
				return opt.value.(string), nil
			}
			return "", fmt.Errorf("'%s': not a string", option)
		}
	}
	return "", fmt.Errorf("'%s': unknown option", option)
}

// GetStringDefault is similar to GetString but given default value is returned
// if option does not exist or is of wrong type.
func (c *Configuration) GetStringDefault(option string, dfault string) string {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _STRING_TYPE {
				return opt.value.(string)
			}
		}
	}
	return dfault
}

// GetBoolArray returns the array of booleans associated with given option
// name. An error is flagged if option does not record an array of booleans
// or it is undefined.
func (c *Configuration) GetBoolArray(option string) ([]bool, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _ARRAY_TYPE|_BOOL_TYPE {
				return opt.value.([]bool), nil
			}
			return nil, fmt.Errorf("'%s': not an array of booleans", option)
		}
	}
	return nil, fmt.Errorf("'%s': unknown option", option)
}

// GetBoolArrayDefault is similar to GetBoolArray but given default value is
// returned if option does not exist or is of wrong type.
func (c *Configuration) GetBoolArrayDefault(option string, dfault []bool) []bool {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _ARRAY_TYPE|_BOOL_TYPE {
				return opt.value.([]bool)
			}
		}
	}
	return dfault
}

// GetIntArray returns the array of integers associated with given option
// name. An error is flagged if option does not record an array of integers
// or it is undefined.
func (c *Configuration) GetIntArray(option string) ([]int64, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _ARRAY_TYPE|_INT_TYPE {
				return opt.value.([]int64), nil
			}
			return nil, fmt.Errorf("'%s': not an array of integers", option)
		}
	}
	return nil, fmt.Errorf("'%s': unknown option", option)
}

// GetIntArrayDefault is similar to GetIntArray but given default value is
// returned if option does not exist or is of wrong type.
func (c *Configuration) GetIntArrayDefault(option string, dfault []int64) []int64 {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _ARRAY_TYPE|_INT_TYPE {
				return opt.value.([]int64)
			}
		}
	}
	return dfault
}

// GetFloatArray returns the array of floating-point values associated with
// given option name. An error is flagged if option does not record an array
// of floating-point values or it is undefined.
func (c *Configuration) GetFloatArray(option string) ([]float64, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _ARRAY_TYPE|_FLOAT_TYPE {
				return opt.value.([]float64), nil
			}
			return nil, fmt.Errorf("'%s': not an array of floating-point numbers", option)
		}
	}
	return nil, fmt.Errorf("'%s': unknown option", option)
}

// GetFloatArrayDefault is similar to GetFloatArray but given default value
// is returned if option does not exist or is of wrong type.
func (c *Configuration) GetFloatArrayDefault(option string, dfault []float64) []float64 {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _ARRAY_TYPE|_FLOAT_TYPE {
				return opt.value.([]float64)
			}
		}
	}
	return dfault
}

// GetDateArray returns the array of dates associated with given option
// name. An error is flagged if option does not record an array of dates
// or it is undefined.
func (c *Configuration) GetDateArray(option string) ([]time.Time, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _ARRAY_TYPE|_DATE_TYPE {
				return opt.value.([]time.Time), nil
			}
			return nil, fmt.Errorf("'%s': not an array of dates", option)
		}
	}
	return nil, fmt.Errorf("'%s': unknown option", option)
}

// GetDateArrayDefault is similar to GetDateArray but given default value
// is returned if option does not exist or is of wrong type.
func (c *Configuration) GetDateArrayDefault(option string, dfault []time.Time) []time.Time {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _ARRAY_TYPE|_DATE_TYPE {
				return opt.value.([]time.Time)
			}
		}
	}
	return dfault
}

// GetStringArray returns the array of strings associated with given option
// name. An error is flagged if option does not record an array of strings
// or it is undefined.
func (c *Configuration) GetStringArray(option string) ([]string, error) {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _ARRAY_TYPE|_STRING_TYPE {
				return opt.value.([]string), nil
			}
			return nil, fmt.Errorf("'%s': not an array of strings", option)
		}
	}
	return nil, fmt.Errorf("'%s': unknown option", option)
}

// GetStringArrayDefault is similar to GetStringArray but given default value
// is returned if option does not exist or is of wrong type.
func (c *Configuration) GetStringArrayDefault(option string, dfault []string) []string {
	if c != nil {
		if opt := c.getOption(option); opt != nil {
			if opt.ctype == _ARRAY_TYPE|_STRING_TYPE {
				return opt.value.([]string)
			}
		}
	}
	return dfault
}

func buildOptionPath(section, option string) string {
	if section == "" {
		return option
	}
	return section + "." + option
}

func (c *Configuration) getOption(key string) *configurationValue {
	key = strings.ToLower(key)

	c.RLock()
	defer c.RUnlock()
	if value, found := c.options[key]; found == true {
		return &value
	}
	return nil
}

func (c *Configuration) setOption(key string, value interface{}) {
	key = strings.ToLower(key)

	// May be nil if an error was detected during parsing. Can be safely
	// ignored since an error was (or will be) generated by the parser.
	if value != nil {
		// Records section.
		c.Lock()
		defer c.Unlock()

		// Adds to list of sections if new one.
		s := c.getSection(key)
		if _, exists := c.sections[s]; exists == false {
			c.sections[s] = struct{}{}
		}

		// Records option. Internal representation used by the parser is
		// updated so that no reflection is required when accessing the
		// options to speed up read operations. Reflection is fine during
		// parsing since it is usually only performed during startup.
		rv := reflect.ValueOf(value)
		if rv.Kind() == sliceType {
			c.options[key] = c.setArray(key, rv)

		} else {
			c.options[key] = c.setValue(key, rv)
		}
	}
}

func (c *Configuration) setValue(key string, rv reflect.Value) configurationValue {
	value := configurationValue{}

	if rv.Type() == dateType {
		value.ctype = _DATE_TYPE
		value.value = rv.Interface().(time.Time)
	} else {
		switch rv.Kind() {
		case boolType:
			value.ctype = _BOOL_TYPE
			value.value = rv.Bool()
		case intType:
			value.ctype = _INT_TYPE
			value.value = rv.Int()
		case floatType:
			value.ctype = _FLOAT_TYPE
			value.value = rv.Float()
		case stringType:
			value.ctype = _STRING_TYPE
			value.value = rv.String()
		default:
			panic(fmt.Sprintf("unexpected type '%s'", reflect.TypeOf(rv).Kind()))
		}
	}
	return value
}

func (c *Configuration) setArray(key string, rv reflect.Value) configurationValue {
	value := configurationValue{}
	if reflect.TypeOf(rv.Index(0).Interface()) == dateType {
		value.ctype = _ARRAY_TYPE | _DATE_TYPE
		a := []time.Time{}
		for i := 0; i < rv.Len(); i++ {
			a = append(a, rv.Index(i).Interface().(time.Time))
		}
		value.value = a
	} else {
		switch reflect.ValueOf(rv.Index(0).Interface()).Kind() {
		case boolType:
			value.ctype = _ARRAY_TYPE | _BOOL_TYPE
			a := []bool{}
			for i := 0; i < rv.Len(); i++ {
				a = append(a, reflect.ValueOf(rv.Index(i).Interface()).Bool())
			}
			value.value = a
		case intType:
			value.ctype = _ARRAY_TYPE | _INT_TYPE
			a := []int64{}
			for i := 0; i < rv.Len(); i++ {
				a = append(a, reflect.ValueOf(rv.Index(i).Interface()).Int())
			}
			value.value = a
		case floatType:
			value.ctype = _ARRAY_TYPE | _FLOAT_TYPE
			a := []float64{}
			for i := 0; i < rv.Len(); i++ {
				a = append(a, reflect.ValueOf(rv.Index(i).Interface()).Float())
			}
			value.value = a
		case stringType:
			value.ctype = _ARRAY_TYPE | _STRING_TYPE
			a := []string{}
			for i := 0; i < rv.Len(); i++ {
				a = append(a, reflect.ValueOf(rv.Index(i).Interface()).String())
			}
			value.value = a
		default:
			panic(fmt.Sprintf("unexpected type '%s'", reflect.TypeOf(rv).Kind()))
		}
	}
	return value
}

func (c *Configuration) getSection(option string) string {
	s := strings.Split(option, ".")
	// Size is 1 if option was declared outside of a section.
	if len(s) > 1 {
		return s[0]
	}
	return ""
}

// Error dumps a configuration-error to a string.
func (c *ConfigurationError) Error() string {
	if c != nil {
		return c.msg
	}
	return ""
}

// String dumps a configuration to a string.
func (c *Configuration) String() string {
	buf := ""
	for k, v := range c.options {
		rv := reflect.ValueOf(v.value)
		if rv.Type() == dateType {
			buf += fmt.Sprintf("%s = %s\n", k, valueToString(rv))
		} else {
			switch rv.Kind() {
			case boolType:
				buf += fmt.Sprintf("%s = %s\n", k, valueToString(rv))
			case intType:
				buf += fmt.Sprintf("%s = %s\n", k, valueToString(rv))
			case floatType:
				buf += fmt.Sprintf("%s = %s\n", k, valueToString(rv))
			case stringType:
				buf += fmt.Sprintf("%s = %s\n", k, valueToString(rv))
			case sliceType:
				buf += fmt.Sprintf("%s = [", k)
				array := rv
				for i := 0; i < array.Len(); i++ {
					if i > 0 {
						buf += ", "
					}
					buf += fmt.Sprintf("%s",
						valueToString(reflect.ValueOf(array.Index(i).Interface())))
				}
				buf += fmt.Sprintf("]\n")
			default:
				panic(fmt.Sprintf("unexpected type '%s'",
					reflect.TypeOf(v).Kind()))
			}
		}
	}
	return buf
}

func valueToString(v reflect.Value) string {
	if v.Type() == dateType {
		date := v.Interface().(time.Time).Format(time.RFC3339)
		return fmt.Sprintf("%s", date)
	}
	switch v.Kind() {
	case boolType:
		return fmt.Sprintf("%t", v.Bool())
	case intType:
		return fmt.Sprintf("%d", v.Int())
	case floatType:
		return fmt.Sprintf("%f", v.Float())
	case stringType:
		return fmt.Sprintf("\"%s\"", v.String())
	default:
		panic(fmt.Sprintf("unexpected type '%s'", reflect.TypeOf(v).Kind()))
	}
}

func (c configurationType) String() string {
	prefix := ""
	if c&_ARRAY_TYPE != 0 {
		prefix = "[]"
		c = c ^ _ARRAY_TYPE
	}
	typ := ""
	switch c {
	case _BOOL_TYPE:
		typ = "bool"
	case _INT_TYPE:
		typ = "int64"
	case _FLOAT_TYPE:
		typ = "float64"
	case _DATE_TYPE:
		typ = "time.Time"
	case _STRING_TYPE:
		typ = "string"
	}
	return prefix + typ
}
