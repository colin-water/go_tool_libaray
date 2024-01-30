package slice

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		want []string
	}{
		{
			name: "src nil",
			want: []string{},
		},
		{
			name: "src empty",
			src:  []int{},
			want: []string{},
		},
		{
			name: "src has element",
			src:  []int{1, 2, 3},
			want: []string{"1", "2", "3"},
		},
	}

	//strconv.Itoa 将整数 src 转换为对应的字符串表示。1，2，3 => '1'，'2'，'3'
	//这是为了确保不同的整数值在字符串形式上是相等的。
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Map(tt.src, func(idx int, src int) string {
				return strconv.Itoa(src)
			})
			assert.Equal(t, res, tt.want)
		})
	}
}
func TestFilterMap(t *testing.T) {
	tests := []struct {
		name string
		src  []int
		want []string
	}{
		{
			name: "src nil",
			want: []string{},
		},
		{
			name: "src empty",
			src:  []int{},
			want: []string{},
		},
		{
			name: "src has element",
			src:  []int{1, -2, 3},
			want: []string{"1", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := FilterMap(tt.src, func(idx int, src int) (string, bool) {
				return strconv.Itoa(src), src >= 0
			})
			assert.Equal(t, res, tt.want)
		})
	}
}