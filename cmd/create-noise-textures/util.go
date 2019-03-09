package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/cgm"
)

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

func combine(images ...[]uint8) []uint8 {
	size := len(images[0])
	imageCount := len(images)
	result := make([]uint8, size)

	for i := 0; i < size; i++ {
		sum := 0.0
		for j := 0; j < imageCount; j++ {
			sum += float64(images[j][i])
		}
		result[i] = uint8(sum / float64(imageCount))
	}

	return result
}

func remapAll(image1 []uint8, image2 []uint8) []uint8 {
	size := len(image1)
	result := make([]uint8, size)

	for i := 0; i < size; i++ {
		val1 := float32(image1[i]) / 255.0
		val2 := float32(image2[i]) / 255.0
		val1 = cgm.Clamp(val1, 1-val2, 1.0)
		res := cgm.Map(val1, 1-val2, 1.0, 0.0, 1.0)
		result[i] = uint8(res * 255)
	}

	return result
}
