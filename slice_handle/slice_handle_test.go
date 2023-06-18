package slice_handle

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindIndex(t *testing.T) {
	unitTest := []struct {
		Name       string
		InputArr   []interface{}
		FHandleArr func(e interface{}) bool
		Expected   int
	}{
		{
			Name:     "UT_IntArrayResult0",
			InputArr: []interface{}{1, 4, 7, 8, -14, 44},
			FHandleArr: func(e interface{}) bool {
				return e == 1
			},
			Expected: 0,
		},
		{
			Name:     "UT_IntArrayResultNotFound",
			InputArr: []interface{}{1, 4, 7, 8, -14, 44},
			FHandleArr: func(e interface{}) bool {
				return e == 5
			},
			Expected: -1,
		},
		{
			Name:     "UT_IntArrayResult3",
			InputArr: []interface{}{1, 4, 7, 8, -14, 44},
			FHandleArr: func(e interface{}) bool {
				return e == 8
			},
			Expected: 3,
		},
		{
			Name:     "UT_IntArrayResultNearest",
			InputArr: []interface{}{1, 4, 7, 8, -14, 44, 4},
			FHandleArr: func(e interface{}) bool {
				return e == 4
			},
			Expected: 1,
		},
		{
			Name:     "UT_IntArrayGt4",
			InputArr: []interface{}{1, 4, 7, 8, -14, 44, 4},
			FHandleArr: func(e interface{}) bool {
				return e.(int) > 4
			},
			Expected: 2,
		},
		{
			Name:     "UT_IntArrayLt0",
			InputArr: []interface{}{1, 4, 7, 8, -14, 44, 4},
			FHandleArr: func(e interface{}) bool {
				return e.(int) < 0
			},
			Expected: 4,
		},
	}
	for _, ut := range unitTest {
		actual := FindIndex(ut.InputArr, ut.FHandleArr)
		if assert.Equal(t, ut.Expected, actual) == false {
			t.Errorf(fmt.Sprintf("Run %s fail", ut.Name))
		}
	}
}

func BenchmarkSliceHandle(b *testing.B) {
	InputArr := []interface{}{1, 4, 7, 8, -14, 44}
	FHandleArr := func(e interface{}) bool {
		return e == 1
	}
	for index := 0; index < b.N; index++ {
		FindIndex(InputArr, FHandleArr)
	}
}

func TestWhere(t *testing.T) {
	baseInputInt := []interface{}{1, 4, 7, 8, -14, 44}
	baseInputStruct := []interface{}{struct {
		V int
	}{V: 1}, struct {
		V int
	}{V: 4}, struct {
		V int
	}{V: 7}, struct {
		V int
	}{V: 8}, struct {
		V int
	}{V: -14}, struct {
		V int
	}{V: 44}, struct {
		V int
	}{V: 4}}
	unitTest := []struct {
		Name       string
		InputArr   []interface{}
		FHandleArr func(e interface{}) bool
		Expected   []interface{}
	}{
		{
			Name:     "UT_IntArray1",
			InputArr: baseInputInt,
			FHandleArr: func(e interface{}) bool {
				return e == 1
			},
			Expected: []interface{}{1},
		},
		{
			Name:     "UT_IntArray2",
			InputArr: baseInputInt,
			FHandleArr: func(e interface{}) bool {
				return e.(int) > 5
			},
			Expected: []interface{}{7, 8, 44},
		},
		{
			Name:     "UT_IntArray3",
			InputArr: baseInputInt,
			FHandleArr: func(e interface{}) bool {
				return e.(int) < 5
			},
			Expected: []interface{}{1, 4, -14},
		},
		{
			Name:     "UT_StructArray3",
			InputArr: baseInputStruct,
			FHandleArr: func(e interface{}) bool {
				s := e.(struct{ V int })
				return s.V == 4
			},
			Expected: []interface{}{struct{ V int }{V: 4}, struct{ V int }{V: 4}},
		},
	}
	for _, ut := range unitTest {
		actual := Where(ut.InputArr, ut.FHandleArr)
		if assert.Equal(t, ut.Expected, actual) == false {
			t.Errorf(fmt.Sprintf("Run %s fail", ut.Name))
		}
	}
}
