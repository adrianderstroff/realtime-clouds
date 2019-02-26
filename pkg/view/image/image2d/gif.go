package image2d

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
)

func SaveGif(filepath string, images []Image2D) error {
	outGif := &gif.GIF{}
	for _, image2D := range images {
		simage, err := image2D.ToImage()
		if err != nil {
			return err
		}
		palettedImage := image.NewPaletted(simage.Bounds(), palette.Plan9)
		draw.Draw(palettedImage, palettedImage.Rect, simage, simage.Bounds().Min, draw.Over)

		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 0)
	}

	f, _ := os.OpenFile("rgb.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)

	return nil
}
