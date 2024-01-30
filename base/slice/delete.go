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
