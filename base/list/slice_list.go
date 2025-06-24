package list

import (
	"github.com/colin-water/go_tool_libaray/base/common"
	"github.com/colin-water/go_tool_libaray/base/slice"
)

type ArrayList[T any] struct {
	vals []T
}

var _ List[any] = &ArrayList[any]{}

func NewArrayList[T any](capacity int) *ArrayList[T] {
	return &ArrayList[T]{vals: make([]T, 0, capacity)}
}
func NewArrayListWithData[T any](data []T) *ArrayList[T] {
	return &ArrayList[T]{vals: data}
}

func (a *ArrayList[T]) Get(index int) (t T, e error) {
	length := a.Len()
	if index < 0 || length <= index {
		return t, common.NewErrIndexOutOfRange(length, index)
	}
	return a.vals[index], nil
}

func (a *ArrayList[T]) Append(values ...T) error {
	a.vals = append(a.vals, values...)
	return nil
}

// Add 在列表下标为index的位置插入一个元素
func (a *ArrayList[T]) Add(index int, t T) error {
	if index < 0 || index > len(a.vals) {
		return common.NewErrIndexOutOfRange(len(a.vals), index)
	}
	result, err := slice.Add(a.vals, t, index)
	a.vals = result
	return err
}

// Set 设置列表中下标为index位置的元素值为t
func (a *ArrayList[T]) Set(index int, t T) error {
	length := len(a.vals)
	if index >= length || index < 0 {
		return common.NewErrIndexOutOfRange(length, index)
	}
	a.vals[index] = t
	return nil
}

// Delete 从列表中删除指定下标的元素，并返回删除的元素值
func (a *ArrayList[T]) Delete(index int) (T, error) {
	result, val, err := slice.Delete(a.vals, index)
	if err != nil {
		return val, err
	}
	a.vals = result
	a.shrink()
	return val, nil
}
func (a *ArrayList[T]) shrink() {
	a.vals = slice.Shrink(a.vals)
}

// Len 返回列表的长度
func (a *ArrayList[T]) Len() int {
	return len(a.vals)
}

// Cap 返回列表的容量
func (a *ArrayList[T]) Cap() int {
	return cap(a.vals)
}

// Range 遍历列表，对每个元素执行指定的函数 fn
func (a *ArrayList[T]) Range(fn func(index int, t T) error) error {
	for key, value := range a.vals {
		e := fn(key, value)
		if e != nil {
			return e
		}
	}
	return nil
}

// AsSlice 返回列表的副本作为切片
func (a *ArrayList[T]) AsSlice() []T {
	res := make([]T, len(a.vals))
	copy(res, a.vals)
	return res
}
