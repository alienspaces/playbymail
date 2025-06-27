package convert

func Ptr[T any](v T) *T {
	return &v
}

func PtrStrict[T comparable](v T) *T {
	zero := new(T)
	if *zero == v {
		return nil
	}

	return &v
}

func Slicep[T any](s []T) *[]T {
	if len(s) == 0 {
		return nil
	}

	return &s
}

func Bool(b *bool) bool {
	return b != nil && *b
}
