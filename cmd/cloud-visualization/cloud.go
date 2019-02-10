package main

import (
	"strconv"

	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/image/image3d"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

func Make3DCloudTexture(basePath string, size int) (texture.Texture, error) {
	imagePaths := make([]string, size)
	for i := 0; i < size; i++ {
		imagePaths[i] = basePath + strconv.Itoa(i) + ".png"
	}
	cloudBaseImage, err := image3d.MakeFromPath(imagePaths)
	if err != nil {
		return texture.Texture{}, err
	}
	return texture.Make3DFromImage(&cloudBaseImage, gl.RGBA, gl.RGBA)
}
