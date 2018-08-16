// Package geometry provides geometric primitives that can be used in meshes.
// It also offers a way to create custom geometric shapes.
package geometry

// combine merges multiple slices into one
func Combine(slices ...[]float32) []float32 {
	var result []float32
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// repeat creates a slice that consists of the provided slices multiple times repeated.
func Repeat(slice []float32, number int) []float32 {
	var result []float32
	for i := 0; i < number; i++ {
		result = append(result, slice...)
	}
	return result
}
