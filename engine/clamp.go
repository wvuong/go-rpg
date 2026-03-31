package engine

func ClampFloat64(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func CycleNext[T any](index int, slice []T) int {
	v := index + 1
	if v >= len(slice) {
		return 0
	}

	return v
}

func CyclePrevious[T any](index int, slice []T) int {
	v := index - 1
	if v < 0 {
		return len(slice) - 1
	}

	return v
}
