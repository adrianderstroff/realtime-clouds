package cloud

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
