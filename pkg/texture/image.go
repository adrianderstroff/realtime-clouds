// Package texture provides classes for creating and storing images and textures.
package texture

import "unsafe"

type Image struct {
	format         uint32
	internalFormat int32
	width          int32
	height         int32
	pixelType      uint32
	data           unsafe.Pointer
}

func MakeImage(path string) (Image, error) {
	return Image{}, nil
}
