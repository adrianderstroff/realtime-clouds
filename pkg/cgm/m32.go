package cgm

import "math"

// Min32 is a float32 variant of math.Min which is float64
func Min32(a, b float32) float32 {
	return float32(math.Min(float64(a), float64(b)))
}

// Max32 is a float32 variant of math.Max which is float64
func Max32(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

// Sqrt32 is a float32 variant of math.Sqrt which is float64
func Sqrt32(a float32) float32 {
	return float32(math.Sqrt(float64(a)))
}

// Floor32 is a float32 variant of math.Floor which is float64
func Floor32(a float32) float32 {
	return float32(math.Floor(float64(a)))
}

// Mod32 is a float32 variant of math.Mod which is float64
func Mod32(a, b float32) float32 {
	return float32(math.Mod(float64(a), float64(b)))
}
