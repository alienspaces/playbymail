package slice

import "testing"

func TestEqual(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		a    []T
		b    []T
		want bool
	}

	strTests := []testCase[string]{
		{
			name: "string slices \\ equal ordered \\ true",
			a:    []string{"a", "b"},
			b:    []string{"a", "b"},
			want: true,
		},
		{
			name: "string slices \\ equal unordered \\ false",
			a:    []string{"b", "a"},
			b:    []string{"a", "b"},
			want: false,
		},
		{
			name: "string slices \\ unequal ordered \\ false",
			a:    []string{"a", "b"},
			b:    []string{"a", "b", "c"},
			want: false,
		},
	}

	for _, tt := range strTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.a, tt.b); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}

	intTests := []testCase[int]{
		{
			name: "int slices \\ equal ordered \\ true",
			a:    []int{1, 2},
			b:    []int{1, 2},
			want: true,
		},
		{
			name: "int slices \\ equal unordered \\ false",
			a:    []int{2, 1},
			b:    []int{1, 2},
			want: false,
		},
		{
			name: "int slices \\ unequal ordered \\ false",
			a:    []int{1, 2},
			b:    []int{1, 2, 3},
			want: false,
		},
	}

	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.a, tt.b); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
