package config_test

import (
	"fmt"
	. "launchpad.net/gocheck"
	"testing"
)

func TestAll(t *testing.T) {
	TestingT(t)
}

// EqualSlice checks whether two arrays are equal or not.
type equalSliceChecker struct {
	*CheckerInfo
}

var EqualSlice = &equalSliceChecker{
	&CheckerInfo{Name: "EqualSlice", Params: []string{"sliceObtained", "sliceExpected"}},
}

func (checker *equalSliceChecker) Check(params []interface{}, names []string) (result bool, err string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			err = fmt.Sprint(v)
		}
	}()
	if fmt.Sprintf("'%v'", params[0]) == fmt.Sprintf("'%v'", params[1]) {
		return true, err
	}
	return false, err
}
