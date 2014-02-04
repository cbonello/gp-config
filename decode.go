package config

import (
	"fmt"
	"reflect"
	"time"
)

// Decode initializes structPtr with contents of given section. Function does
// not support circular types; it will loop forever.
func (c *Configuration) Decode(section string, structPtr interface{}) (err error) {
	if c != nil {
		if c.IsSection(section) == false {
			return fmt.Errorf("'%s': unknown section", section)
		}
		if structPtr == nil {
			return fmt.Errorf("structure argument cannot be a nil value")
		}
		structPtrType := reflect.TypeOf(structPtr)
		if structPtrType.Kind() != ptrType {
			return fmt.Errorf("structure argument is not a pointer")
		}
		structPtrVal := reflect.ValueOf(structPtr)
		if structPtrVal.IsNil() {
			return fmt.Errorf("structure argument cannot be a nil pointer")
		}
		if structPtrVal.Elem().Kind() != structType {
			return fmt.Errorf("structure argument is not a pointer to a structure")
		}
		err = c.doDecode(section, structPtrVal.Elem(), structPtrType.Elem())
	}
	return err
}

func (c *Configuration) doDecode(section string, val reflect.Value, typ reflect.Type) (err error) {
	sVal := val
	sType := typ
	numFields := sVal.NumField()
	// For each field of destination structure...
	for f := 0; f < numFields; f++ {
		fieldVal := sVal.Field(f)
		fieldType := sType.Field(f)
		// Embedded field?
		if fieldType.Anonymous {
			if fieldType.Type.Kind() == structType {
				if err := c.doDecode(section, fieldVal, fieldType.Type); err != nil {
					return err
				}

			} else {
				return fmt.Errorf(
					"'%s': embedded pointer fields are not supported yet",
					fieldType.Name)
			}

			// My failed attempt to support embedded pointer fields. For next
			// release I guess..
			/*
				} else {
					structPtrType := fieldType.Type
					if structPtrType.Kind() == ptrType {
						structPtrVal := reflect.ValueOf(fieldType.Type)
						if structPtrVal.IsNil() == false {
							if structPtrVal.Elem().Kind() == structType {
								err := c.doDecode(section, structPtrVal.Elem(), structPtrType.Elem())
								if err != nil {
									return err
								}
							}
						}
					}
				}
			*/
		} else {
			// Build option's path from section and either the StructTag or
			// the field name. See http://golang.org/pkg/reflect/#StructTag.
			tag := fieldType.Tag.Get("option")
			if tag == "" {
				tag = fieldType.Name
			}
			path := buildOptionPath(section, tag)
			// Path corresponds to an existing option?
			if src := c.getOption(path); src != nil {
				if fieldVal.IsValid() == false {
					return fmt.Errorf("'%s': cannot set field's value",
						fieldType.Name)
				}
				if fieldVal.CanSet() == false {
					return fmt.Errorf("'%s': cannot set value of unexported struct field",
						fieldType.Name)
				}
				if fieldVal.Type().Kind() == sliceType {
					eltType := fieldType.Type.Elem()
					if err := c.decodeSlice(path, src, fieldVal, eltType); err != nil {
						return err
					}
				} else {
					if err := c.decodeValue(path, src, fieldVal); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (c *Configuration) decodeSlice(path string, src *configurationValue,
	dst reflect.Value, eltType reflect.Type) error {

	srcVal := reflect.ValueOf(src.value)
	if eltType == dateType {
		if src.ctype != _ArrayType|_DateType {
			return fmt.Errorf(
				"'%s': value of type %s is not assignable to type []time.Time",
				path, src.ctype)
		}
		a := []time.Time{}
		for i := 0; i < srcVal.Len(); i++ {
			a = append(a, srcVal.Index(i).Interface().(time.Time))
		}
		dst.Set(reflect.ValueOf(a))
	} else {
		switch eltType.Kind() {
		case boolType:
			if src.ctype != _ArrayType|_BoolType {
				return fmt.Errorf("'%s': value of type %s is not assignable to type []bool",
					path, src.ctype)
			}
			a := []bool{}
			for i := 0; i < srcVal.Len(); i++ {
				a = append(a, reflect.ValueOf(srcVal.Index(i).Interface()).Bool())
			}
			dst.Set(reflect.ValueOf(a))
		case intType:
			if src.ctype != _ArrayType|_IntType {
				return fmt.Errorf("'%s': value of type %s is not assignable to type []int64", path,
					src.ctype)
			}
			a := []int64{}
			for i := 0; i < srcVal.Len(); i++ {
				a = append(a, reflect.ValueOf(srcVal.Index(i).Interface()).Int())
			}
			dst.Set(reflect.ValueOf(a))
		case floatType:
			if src.ctype != _ArrayType|_FloatType {
				return fmt.Errorf("'%s': value of type %s is not assignable to type []float64", path,
					src.ctype)
			}
			a := []float64{}
			for i := 0; i < srcVal.Len(); i++ {
				a = append(a, reflect.ValueOf(srcVal.Index(i).Interface()).Float())
			}
			dst.Set(reflect.ValueOf(a))
		case stringType:
			if src.ctype != _ArrayType|_StringType {
				return fmt.Errorf("'%s': value of type %s is not assignable to type []string", path,
					src.ctype)
			}
			a := []string{}
			for i := 0; i < srcVal.Len(); i++ {
				a = append(a, reflect.ValueOf(srcVal.Index(i).Interface()).String())
			}
			dst.Set(reflect.ValueOf(a))
		default:
			// Type not supported, '[]float' for instance. Default case is
			// unlikely to be executed since unsupported type errors are
			// trapped while configuration files are loaded.
			return fmt.Errorf("'%s': value of type %s is not assignable to type %s", path,
				dst.Kind(), src.ctype)
		}

	}
	return nil
}

func (c *Configuration) decodeValue(path string, src *configurationValue, dst reflect.Value) error {
	if dst.Type() == dateType {
		if src.ctype != _DateType {
			return fmt.Errorf(
				"'%s': value of type %s is not assignable to type time.Time", path,
				src.ctype)
		}
		dst.Set(reflect.ValueOf(src.value.(time.Time)))
	} else {
		switch dst.Kind() {
		case boolType:
			if src.ctype != _BoolType {
				return fmt.Errorf(
					"'%s': value of type %s is not assignable to type bool", path,
					src.ctype)
			}
			dst.SetBool(src.value.(bool))
		case intType:
			if src.ctype != _IntType {
				return fmt.Errorf(
					"'%s': value of type %s is not assignable to type int64", path,
					src.ctype)
			}
			dst.SetInt(src.value.(int64))
		case floatType:
			if src.ctype != _FloatType {
				return fmt.Errorf(
					"'%s': value of type %s is not assignable to type float64", path,
					src.ctype)
			}
			dst.SetFloat(src.value.(float64))
		case stringType:
			if src.ctype != _StringType {
				return fmt.Errorf(
					"'%s': value of type %s is not assignable to type string", path,
					src.ctype)
			}
			dst.SetString(src.value.(string))
		default:
			// Type not supported, 'int' for instance.
			return fmt.Errorf(
				"'%s': value of type %s is not assignable to type %s", path,
				dst.Kind(), src.ctype)
		}
	}
	return nil
}
