package convert

func String(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func Float64(f *float64) float64 {
	if f == nil {
		return 0
	}

	return *f
}

func Int64(i *int64) int64 {
	if i == nil {
		return 0
	}

	return *i
}

func Int(i *int) int {
	if i == nil {
		return 0
	}

	return *i
}
