package common

// 返回指针
func ToPtr[T any](t T) *T {
	return &t
}
