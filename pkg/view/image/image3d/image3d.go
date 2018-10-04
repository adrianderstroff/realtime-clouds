package image3d

import (
	"errors"
	"fmt"
	_ "image/jpeg"
	"path/filepath"
	"strings"
	"unsafe"

	gl "github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/image/image2d"
)

// Image stores the dimensions, data format and it's pixel data.
// It can be used to manipulate single pixels and is used to
// upload it's data to a texture.
type Image3D struct {
	width     int
	height    int
	slices    int
	channels  int
	pixelType uint32
	data      []image2d.Image2D
}

// Make constructs an image of the specified length x width x height and with all pixels
// set to the specified rgba value.
func Make(width, height, slices, channels int) (Image3D, error) {
	// create image data
	var data []image2d.Image2D
	for i := 0; i < slices; i++ {
		image, err := image2d.Make(width, height, channels)
		if err != nil {
			return Image3D{}, err
		}
		data = append(data, image)
	}

	return Image3D{
		width:     width,
		height:    height,
		slices:    slices,
		channels:  channels,
		pixelType: data[0].GetPixelType(),
		data:      data,
	}, nil
}

// MakeFromData constructs an image of the specified width, height, slices and the specified data.
func MakeFromData(width, height, slices int, data []uint8) (Image3D, error) {
	// determine number of channels
	channels := len(data) / (width * height * slices)

	// create the individual images
	var images []image2d.Image2D
	size := width * height
	for i := 0; i < slices; i++ {
		s, e := i*size, (i+1)*size
		image, err := image2d.MakeFromData(width, height, data[s:e])
		if err != nil {
			return Image3D{}, err
		}
		images = append(images, image)
	}

	return Image3D{
		width:     width,
		height:    height,
		slices:    slices,
		channels:  channels,
		pixelType: uint32(gl.UNSIGNED_BYTE),
		data:      images,
	}, nil
}

// MakeImageFromPath constructs the image data from the specified paths.
// If there is no image at the specified path an error is returned instead.
// The dimensions of all images must match.
func MakeFromPath(paths []string) (Image3D, error) {
	// early exit if no path had been provided
	if len(paths) == 0 {
		return Image3D{}, errors.New("No image paths provided")
	}

	// get the first path
	first, err := image2d.MakeFromPath(paths[0])
	if err != nil {
		return Image3D{}, err
	}

	// load multiple images for each path
	var images []image2d.Image2D
	images = append(images, first)
	for _, path := range paths[1:] {
		// load current image
		image, err := image2d.MakeFromPath(path)
		if err != nil {
			return Image3D{}, err
		}

		// check if dimensions and formats match
		if first.GetWidth() != image.GetWidth() ||
			first.GetHeight() != image.GetHeight() ||
			first.GetChannels() != image.GetChannels() ||
			first.GetPixelType() != image.GetPixelType() {
			return Image3D{}, errors.New("Image dimensions or formats don't match.")
		}

		// append images
		images = append(images, image)
	}

	return Image3D{
		width:     first.GetWidth(),
		height:    first.GetHeight(),
		slices:    len(paths),
		channels:  first.GetChannels(),
		pixelType: first.GetPixelType(),
		data:      images,
	}, nil
}

// SaveToPath saves all slices as png images to the specified path.
// All images will be enumerated starting with 0.
// The file names will look like the following: dir/filename<NUMBER>.png
func (image *Image3D) SaveToPath(path string) error {
	ext := filepath.Ext(path)
	pathnoext := strings.TrimSuffix(path, ext)
	for i := 0; i < image.slices; i++ {
		curpath := pathnoext + fmt.Sprint(i) + ext
		err := image.data[i].SaveToPath(curpath)
		if err != nil {
			return err
		}
	}

	return nil
}

// FlipX flips all slices horizontally.
func (image *Image3D) FlipX() {
	for _, slice := range image.data {
		slice.FlipX()
	}
}

// FlipY flips all slices vertically.
func (image *Image3D) FlipY() {
	for _, slice := range image.data {
		slice.FlipY()
	}
}

// GetWidth returns the width of the image.
func (image *Image3D) GetWidth() int {
	return image.width
}

// GetHeight returns the height of the image.
func (image *Image3D) GetHeight() int {
	return image.height
}

// GetSlices returns the number of slices of the image.
func (image *Image3D) GetSlices() int {
	return image.slices
}

// GetChannels return the number of the channels of the image.
func (image *Image3D) GetChannels() int {
	return image.channels
}

// GetPixelType gets the data type of the pixel data.
func (image *Image3D) GetPixelType() uint32 {
	return image.pixelType
}

// GetDataPointer returns an pointer to the beginning of the image data.
func (image *Image3D) GetDataPointer() unsafe.Pointer {
	return gl.Ptr(image.data)
}

// GetData returns a copy of the images data
func (image *Image3D) GetData() []uint8 {
	// collect data of all slices
	var data []uint8
	for _, slice := range image.data {
		data = append(data, slice.GetData()...)
	}

	return data
}

// GetR returns the red value of the pixel at (x,y) in slice z.
func (image *Image3D) GetR(x, y, z int) uint8 {
	return image.data[z].GetR(x, y)
}

// GetG returns the green value of the pixel at (x,y) in slice z.
func (image *Image3D) GetG(x, y, z int) uint8 {
	return image.data[z].GetG(x, y)
}

// GetB returns the blue value of the pixel at (x,y) in slice z.
func (image *Image3D) GetB(x, y, z int) uint8 {
	return image.data[z].GetB(x, y)
}

// GetA returns the alpha value of the pixel at (x,y) in slice z.
func (image *Image3D) GetA(x, y, z int) uint8 {
	return image.data[z].GetA(x, y)
}

// GetRGB returns the RGB values of the pixel at (x,y) in slice z.
func (image *Image3D) GetRGB(x, y, z int) (uint8, uint8, uint8) {
	return image.data[z].GetRGB(x, y)
}

// GetRGBA returns the RGBA value of the pixel at (x,y) in slice z.
func (image *Image3D) GetRGBA(x, y, z int) (uint8, uint8, uint8, uint8) {
	return image.data[z].GetRGBA(x, y)
}

// SetR sets the red value of the pixel at (x,y) in slice z.
func (image *Image3D) SetR(x, y, z int, r uint8) {
	image.data[z].SetR(x, y, r)
}

// SetG sets the green value of the pixel at (x,y) in slice z.
func (image *Image3D) SetG(x, y, z int, g uint8) {
	image.data[z].SetG(x, y, g)
}

// SetB sets the blue value of the pixel at (x,y) in slice z.
func (image *Image3D) SetB(x, y, z int, b uint8) {
	image.data[z].SetB(x, y, b)
}

// SetA sets the alpha value of the pixel at (x,y) in slice z.
func (image *Image3D) SetA(x, y, z int, a uint8) {
	image.data[z].SetA(x, y, a)
}

// SetRGB sets the RGB values of the pixel at (x,y) in slice z.
func (image *Image3D) SetRGB(x, y, z int, r, g, b uint8) {
	image.data[z].SetRGB(x, y, r, g, b)
}

// SetRGBA sets the RGBA values of the pixel at (x,y) in slice z.
func (image *Image3D) SetRGBA(x, y, z int, r, g, b, a uint8) {
	image.data[z].SetRGBA(x, y, r, g, b, a)
}

// String pretty prints information about the image.
func (image Image3D) String() string {
	return fmt.Sprintf("Image3D (%v,%v,%v) %v", image.width, image.height, image.slices, image.channels)
}
