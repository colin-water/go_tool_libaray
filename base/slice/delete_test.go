package slice

import (
	"github.com/yishengzhishui/library/base/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	testCases := []struct {
		name      string
		slice     []int
		index     int
		wantSlice []int
		wantVal   int
		wantErr   error
	}{
		{
			name:      "index 0",
			slice:     []int{123, 100},
			index:     0,
			wantSlice: []int{100},
			wantVal:   123,
		},
		{
			name:      "index middle",
			slice:     []int{123, 124, 125},
			index:     1,
			wantSlice: []int{123, 125},
			wantVal:   124,
		},
		{
			name:    "index out of range",
			slice:   []int{123, 100},
			index:   12,
			wantErr: common.NewErrIndexOutOfRange(2, 12),
		},
		{
			name:    "index less than 0",
			slice:   []int{123, 100},
			index:   -1,
			wantErr: common.NewErrIndexOutOfRange(2, -1),
		},
		{
			name:      "index last",
			slice:     []int{123, 100, 101, 102, 102, 102},
			index:     5,
			wantSlice: []int{123, 100, 101, 102, 102},
			wantVal:   102,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, val, err := Delete(tc.slice, tc.index)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSlice, res)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}

func TestFilterDelete(t *testing.T) {

	testCases := []struct {
		name            string
		src             []int
		deleteCondition func(index int, value int) bool

		wantRes []int
	}{
		{
			name: "空切片",
			src:  []int{},
			deleteCondition: func(index int, value int) bool {
				return false
			},

			wantRes: []int{},
		},
		{
			name: "不删除元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			deleteCondition: func(index int, value int) bool {
				return false
			},

			wantRes: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "删除首位元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6},
			deleteCondition: func(index int, value int) bool {
				return index == 0
			},

			wantRes: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "删除前面两个元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			deleteCondition: func(index int, value int) bool {
				return index == 0 || index == 1
			},

			wantRes: []int{2, 3, 4, 5, 6, 7},
		},
		{
			name: "删除中间单个元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			deleteCondition: func(index int, value int) bool {
				return index == 3
			},

			wantRes: []int{0, 1, 2, 4, 5, 6, 7},
		},
		{
			name: "删除中间多个不连续元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			deleteCondition: func(index int, value int) bool {
				return index == 2 || index == 4
			},

			wantRes: []int{0, 1, 3, 5, 6, 7},
		},
		{
			name: "删除中间多个连续元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			deleteCondition: func(index int, value int) bool {
				return index == 3 || index == 4
			},

			wantRes: []int{0, 1, 2, 5, 6, 7},
		},
		{
			name: "删除中间多个元素，第一部分为一个元素，第二部分为连续元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			deleteCondition: func(index int, value int) bool {
				return index == 2 || index == 4 || index == 5
			},

			wantRes: []int{0, 1, 3, 6, 7},
		},
		{
			name: "删除中间多个元素，第一部分为连续元素，第二部分为一个元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			deleteCondition: func(index int, value int) bool {
				return index == 2 || index == 3 || index == 5
			},

			wantRes: []int{0, 1, 4, 6, 7},
		},
		{
			name: "删除后面两个元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			deleteCondition: func(index int, value int) bool {
				return index == 6 || index == 7
			},

			wantRes: []int{0, 1, 2, 3, 4, 5},
		},
		{
			name: "删除末尾元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			deleteCondition: func(index int, value int) bool {
				return index == 7
			},

			wantRes: []int{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name: "删除所有元素",
			src:  []int{0, 1, 2, 3, 4, 5, 6, 7},
			deleteCondition: func(index int, value int) bool {
				return true
			},

			wantRes: []int{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := FilterDelete(tc.src, tc.deleteCondition)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
