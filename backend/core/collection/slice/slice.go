package slice

// Equal returns true when a and b are of equal length, equal elements and equal element order
func Equal[T comparable](a []T, b []T) bool {

	if len(a) != len(b) {
		return false
	}

	for idx := range a {
		if a[idx] != b[idx] {
			return false
		}
	}

	return true
}

// FromMap converts a map to a slice where the map values are the slice elements.
func FromMap[K comparable, V any](m map[K]V) []V {
	var s []V

	for _, v := range m {
		s = append(s, v)
	}

	return s
}

// FromMapKeys converts a map to a slice where the map keys are the slice elements.
func FromMapKeys[K comparable, V any](m map[K]V) []K {
	var s []K

	for k := range m {
		s = append(s, k)
	}

	return s
}

// Map maps the slice values using the mapFn.
func Map[T any, R any](mapFn func(T) R, s ...T) []R {
	var r []R

	for _, t := range s {
		r = append(r, mapFn(t))
	}

	return r
}

func ToSliceWithPtrs[T any](s []T) []*T {
	var ptrs []*T

	for i := range s {
		ptrs = append(ptrs, &s[i])
	}

	return ptrs
}

func ToSliceWithoutPtrs[T any](s []*T) []T {
	var res []T

	for i := range s {
		res = append(res, *s[i])
	}

	return res
}
