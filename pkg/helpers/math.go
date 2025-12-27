package helpers

func Abs(a, b int64) int64 {
	if a > b {
		return a - b
	}
	return b - a
}
