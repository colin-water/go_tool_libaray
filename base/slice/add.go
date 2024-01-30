package slice

import "github.com/yishengzhishui/library/base/common"

func Add[T any](src []T, element T, index int) ([]T, error) {
	length := len(src)
	if index < 0 || length < index {
		return nil, common.NewErrIndexOutOfRange(length, index)
	}
	//先将src扩展一个元素
	var zeroValue T
	src = append(src, zeroValue)
	// 遍历切片，知道index
	for i := len(src) - 1; i > index; i-- {
		src[i] = src[i-1]
	}
	src[index] = element
	return src, nil
}

func AddV1[T any](src []T, element T, index int) ([]T, error) {
	length := len(src)
	if index < 0 || length < index {
		return nil, common.NewErrIndexOutOfRange(length, index)
	}
	var zeroValue T
	src = append(src, zeroValue)
	copy(src[index+1:], src[index:])
	src[index] = element
	return src, nil
}
