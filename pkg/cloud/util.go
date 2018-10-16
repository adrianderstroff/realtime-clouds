package cloud

import "github.com/adrianderstroff/realtime-clouds/pkg/cgm"

// mergeColorChannels interleaves the pixels of all provided one-channel images.
func mergeColorChannels(images ...[]uint8) []uint8 {
	channels := len(images)
	size := len(images[0])

	result := make([]uint8, channels*size)
	for pixel := 0; pixel < size; pixel++ {
		for imageidx, image := range images {
			result[pixel*channels+imageidx] = image[pixel]
		}
	}

	return result
}

func combine(slice1, slice2 []uint8) []uint8 {
	result := make([]uint8, len(slice1))

	for i := 0; i < len(slice1); i++ {
		result[i] = uint8(cgm.Lerp(float32(slice1[i]), float32(slice2[i]), 0.5))
	}

	return result
}
