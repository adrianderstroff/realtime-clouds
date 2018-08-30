package math

import "math"

func Min32(a, b float32) float32 {
	return float32(math.Min(float64(a), float64(b)))
}

func Max32(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

func Sqrt32(a float32) float32 {
	return float32(math.Sqrt(float64(a)))
}
