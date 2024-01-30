package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: 0,
			name: "first one",
		},
		{
			src:  []int{},
			dst:  1,
			want: -1,
			name: "the length of src is 0",
		},
		{
			dst:  1,
			want: -1,
			name: "src nil",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: -1,
			name: "dst not exist",
		},
		{
			src:  []int{1, 3, 4, 2, 0},
			dst:  0,
			want: 4,
			name: "last one",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, Index[int](test.src, test.dst))
		})
	}
}

func TestIndexFunc(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: 0,
			name: "first one",
		},
		{
			src:  []int{},
			dst:  1,
			want: -1,
			name: "the length of src is 0",
		},
		{
			dst:  1,
			want: -1,
			name: "src nil",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: -1,
			name: "dst not exist",
		},
		{
			src:  []int{1, 3, 4, 2, 0},
			dst:  0,
			want: 4,
			name: "last one",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, IndexFunc[int](test.src, func(src int) bool {
				return src == test.dst
			}))
		})
	}
}

func TestLastIndex(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: 1,
			name: "first one",
		},
		{
			src:  []int{},
			dst:  1,
			want: -1,
			name: "the length of src is 0",
		},
		{
			dst:  1,
			want: -1,
			name: "src nil",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: -1,
			name: "dst not exist",
		},
		{
			src:  []int{0, 1, 3, 4, 2, 0},
			dst:  0,
			want: 5,
			name: "last one",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, LastIndex[int](test.src, test.dst))
	}
}

func TestLastIndexFunc(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: 1,
			name: "first one",
		},
		{
			src:  []int{},
			dst:  1,
			want: -1,
			name: "the length of src is 0",
		},
		{
			dst:  1,
			want: -1,
			name: "src nil",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: -1,
			name: "dst not exist",
		},
		{
			src:  []int{0, 1, 3, 4, 2, 0},
			dst:  0,
			want: 5,
			name: "last one",
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, LastIndexFunc[int](test.src, func(src int) bool {
			return src == test.dst
		}))
	}
}

func TestIndexAll(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want []int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: []int{0, 1},
			name: "normal test",
		},
		{
			src:  []int{},
			dst:  1,
			want: []int{},
			name: "the length of src is 0",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: []int{},
			name: "dst not exist",
		},
		{
			src:  []int{0, 1, 3, 4, 2, 0},
			dst:  0,
			want: []int{0, 5},
			name: "normal test",
		},
	}
	for _, test := range tests {
		res := IndexAll[int](test.src, test.dst)
		assert.ElementsMatch(t, test.want, res)
	}
}

func TestIndexAllFunc(t *testing.T) {
	tests := []struct {
		src  []int
		dst  int
		want []int
		name string
	}{
		{
			src:  []int{1, 1, 3, 5},
			dst:  1,
			want: []int{0, 1},
			name: "normal test",
		},
		{
			src:  []int{},
			dst:  1,
			want: []int{},
			name: "the length of src is 0",
		},
		{
			src:  []int{1, 4, 6},
			dst:  7,
			want: []int{},
			name: "dst not exist",
		},
		{
			src:  []int{0, 1, 3, 4, 2, 0},
			dst:  0,
			want: []int{0, 5},
			name: "normal test",
		},
	}
	for _, test := range tests {
		res := IndexAllFunc[int](test.src, func(src int) bool {
			return src == test.dst
		})
		assert.ElementsMatch(t, test.want, res)
	}
}
