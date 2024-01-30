package slice

// equalFunc 比较两个元素是否相等
type equalFunc[T any] func(paramA, paramB T) bool

// 比较元素是否是某一个值
type matchFunc[T any] func(param T) bool
