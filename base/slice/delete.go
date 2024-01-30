package slice

import (
	"github.com/yishengzhishui/library/base/common"
)

// Delete 删除切片元素，并将删除的元素值返回
func Delete[T any](src []T, index int) ([]T, T, error) {
	length := len(src)
	if index < 0 || index >= length {
		//先将src扩展一个元素
		var zeroValue T
		return nil, zeroValue, common.NewErrIndexOutOfRange(length, index)
	}
	value := src[index]
	// 从index开始，所有的元素向前移动一位
	for i := index; i < length-1; i++ {
		src[i] = src[i+1]
	}
	//剔除队尾重复数据
	src = src[:length-1]
	return src, value, nil
}

// FilterDelete 删除符合条件的元素
// 考虑到性能问题，所有操作都会在原切片上进行
// 满足条件的元素删除后，其他剩余的元素会往前移动，有且只会移动一次
func FilterDelete[T any](src []T, f func(index int, value T) bool) []T {
	// 删除符合条件的元素后的 index
	newIndex := 0
	for index, value := range src {
		// 判断是否满足删除的条件
		//满足条件，跳过当前循环，即不将该元素包含在新的切片中
		if f(index, value) {
			continue
		}
		src[newIndex] = value
		newIndex++
	}
	return src[:newIndex]
}
