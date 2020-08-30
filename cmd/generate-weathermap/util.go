package main

import (
	"math"
	"math/rand"

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

func createAndFillImage(len int, val uint8) []uint8 {
	var image []uint8
	for i := 0; i < len; i++ {
		image = append(image, val)
	}
	return image
}

func createRandom(len int) []uint8 {
	var image []uint8
	for i := 0; i < len; i++ {
		val := uint8(rand.Intn(256))
		image = append(image, val)
	}
	return image
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

		val1 = cgm.Clamp(val1, val2, 1.0)
		res := cgm.Map(val1, 0.0, 1.0, val2, 1.0)
		result[i] = uint8(res * 255)
	}

	return result
}

func scale(image []uint8, scale float32) []uint8 {
	size := len(image)
	result := make([]uint8, size)

	for i := 0; i < size; i++ {
		val := float32(image[i]) * scale
		val = cgm.Clamp(val, 0, 255)
		result[i] = uint8(val)
	}

	return result
}

func invert(image []uint8) []uint8 {
	size := len(image)
	result := make([]uint8, size)

	for i := 0; i < size; i++ {
		result[i] = uint8(255 - image[i])
	}

	return result
}

func max(image []uint8) uint8 {
	size := len(image)

	var max uint8 = 0
	for i := 0; i < size; i++ {
		max = uint8(math.Max(float64(image[i]), float64(max)))
	}

	return max
}

func min(image []uint8) uint8 {
	size := len(image)

	var min uint8 = 255
	for i := 0; i < size; i++ {
		min = uint8(math.Min(float64(image[i]), float64(min)))
	}

	return min
}

func spread(image []uint8) []uint8 {
	min := min(image)
	max := max(image)

	size := len(image)
	result := make([]uint8, size)
	for i := 0; i < size; i++ {
		val := cgm.Map(float32(image[i]), float32(min), float32(max), float32(0), float32(255))
		result[i] = uint8(val)
	}

	return result
}

func threshold(image []uint8, t uint8) []uint8 {
	size := len(image)
	result := make([]uint8, size)

	for i := 0; i < size; i++ {
		res := image[i]
		if res < t {
			res = 0
		}
		result[i] = res
	}

	return result
}
