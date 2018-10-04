package main

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/noise"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/image/image3d"
)

func main() {
	worley := noise.Worley3D(128, 128, 128, 10)
	image, err := image3d.MakeFromData(128, 128, 128, worley)
	if err != nil {
		panic(err)
	}
	println(image.String())
	err = image.SaveToPath("./cmd/create-noises/worley/worley.png")
	if err != nil {
		panic(err)
	}
}
