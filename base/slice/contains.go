package slice

// Contains 判断 src 里面是否存在 val
func Contains[T comparable](src []T, val T) bool {
	return ContainsFunc[T](src, func(param T) bool {
		return param == val
	})
}

// ContainsFunc 判断 src 里面是否存在 符合equal函数的
// equal是func类型，可以是任何函数，
// 这个是为了判定是否满足一定条件的方法，可以是等于，小于等等
// 可以是类似下面这个样的
// func(val int) bool {
//		return val == 3
//	}
func ContainsFunc[T any](src []T, equal func(val T) bool) bool {
	// 遍历调用equal函数进行判断
	for _, v := range src {
		if equal(v) {
			return true
		}
	}
	return false
}

// ContainsAny 判断 sliceA, sliceB 两个切片，是否有一个元素一样
// map 方便快速比较，因为这个map的key就是切片的元素
func ContainsAny[T comparable](sliceA, sliceB []T) bool {
	srcA := toMap(sliceA)
	for _, value := range sliceB {
		if _, ok := srcA[value]; ok {
			return true
		}
	}
	return false
}

// ContainsAll  sliceA是否包含sliceB
func ContainsAll[T comparable](sliceA, sliceB []T) bool {
	srcMap := toMap(sliceA)
	for _, v := range sliceB {
		if _, exist := srcMap[v]; !exist {
			return false
		}
	}
	return true
}

// ContainsAnyFunc 判断 sliceA 里面是否存在 sliceB 中的任何一个元素
// type equalFunc[T any] func(src, dst T) bool， 通常用来判定两个元素是否相等
func ContainsAnyFunc[T any](sliceA, sliceB []T, equal equalFunc[T]) bool {
	for _, valA := range sliceA {
		for _, valB := range sliceB {
			if equal(valA, valB) {
				return true
			}
		}
	}
	return false
}

// ContainsAllFunc dev_test
// ContainsAllFunc 判断 sliceA 里面是否存在 sliceB 中的所有元素
func ContainsAllFunc[T any](sliceA, sliceB []T, equal equalFunc[T]) bool {
	for _, valB := range sliceB {
		if !ContainsFunc(sliceA, func(src T) bool {
			return equal(src, valB)
		}) {
			return false
		}
	}
	return true
}
