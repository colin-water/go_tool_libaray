package slice

import (
	"github.com/yishengzhishui/library/base/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "value",
			input: []int{1},
			want:  1,
		},
		{
			name:  "values",
			input: []int{2, 3, 1},
			want:  3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Max[int](tc.input)
			assert.Equal(t, tc.want, res)
		})
	}
	// 测试是否引发panics
	//如果执行这个代码块时确实引发了 panic，assert.Panics 就会通过测试
	assert.Panics(t, func() {
		Max[int](nil)
	})
	assert.Panics(t, func() {
		Max[int]([]int{})
	})

	maxTypesTest[uint](t)
	maxTypesTest[uint8](t)
	maxTypesTest[uint16](t)
	maxTypesTest[uint32](t)
	maxTypesTest[uint64](t)
	maxTypesTest[int](t)
	maxTypesTest[int8](t)
	maxTypesTest[int16](t)
	maxTypesTest[int32](t)
	maxTypesTest[int64](t)
	maxTypesTest[float32](t)
	maxTypesTest[float64](t)
}

func TestMin(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name:  "value",
			input: []int{3},
			want:  3,
		},
		{
			name:  "values",
			input: []int{3, 1, 2},
			want:  1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Min[int](tc.input)
			assert.Equal(t, tc.want, res)
		})
	}

	assert.Panics(t, func() {
		Min[int](nil)
	})
	assert.Panics(t, func() {
		Min[int]([]int{})
	})

	minTypesTest[uint](t)
	minTypesTest[uint8](t)
	minTypesTest[uint16](t)
	minTypesTest[uint32](t)
	minTypesTest[uint64](t)
	minTypesTest[int](t)
	minTypesTest[int8](t)
	minTypesTest[int16](t)
	minTypesTest[int32](t)
	minTypesTest[int64](t)
	minTypesTest[float32](t)
	minTypesTest[float64](t)
}

func TestSum(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
		want  int
	}{
		{
			name: "nil",
		},
		{
			name:  "empty",
			input: []int{},
		},
		{
			name:  "value",
			input: []int{1},
			want:  1,
		},
		{
			name:  "values",
			input: []int{1, 2, 3},
			want:  6,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Sum[int](tc.input)
			assert.Equal(t, tc.want, res)
		})
	}

	sumTypesTest[uint](t)
	sumTypesTest[uint8](t)
	sumTypesTest[uint16](t)
	sumTypesTest[uint32](t)
	sumTypesTest[uint64](t)
	sumTypesTest[int](t)
	sumTypesTest[int8](t)
	sumTypesTest[int16](t)
	sumTypesTest[int32](t)
	sumTypesTest[int64](t)
	sumTypesTest[float32](t)
	sumTypesTest[float64](t)
}

// maxTypesTest 只是用来测试一下满足 Max 方法约束的所有类型
func maxTypesTest[T common.RealNumber](t *testing.T) {
	res := Max[T]([]T{1, 2, 3})
	assert.Equal(t, T(3), res)
}

// minTypesTest 只是用来测试一下满足 Min 方法约束的所有类型
func minTypesTest[T common.RealNumber](t *testing.T) {
	res := Min[T]([]T{1, 2, 3})
	assert.Equal(t, T(1), res)
}

// sumTypesTest 只是用来测试一下满足 Sum 方法约束的所有类型
func sumTypesTest[T common.RealNumber](t *testing.T) {
	res := Sum[T]([]T{1, 2, 3})
	assert.Equal(t, T(6), res)
}
