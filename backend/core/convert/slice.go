package convert

import "github.com/lib/pq"

func GenericSlice[T any](s []T) []any {
	g := make([]any, len(s)) // `make` to avoid potential copies when appending

	for i := range s {
		g[i] = s[i]
	}

	return g
}

func PqStringArrayToStrSlice(a pq.StringArray) []string {
	return a
}
