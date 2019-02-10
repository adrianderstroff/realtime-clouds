package main

import (
	"fmt"

	"github.com/adrianderstroff/realtime-clouds/pkg/noise"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/image/image2d"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/image/image3d"
)

const (
	TEX_PATH = "./assets/images/textures/"
)

func main() {

	// create cloud base texture
	fmt.Println("Creating Cloud Base Shape")
	p1 := noise.Perlin3D(128, 128, 128, 5)
	w1 := noise.Worley3D(128, 128, 128, 5)
	w2 := noise.Worley3D(128, 128, 128, 6)
	w3 := noise.Worley3D(128, 128, 128, 7)
	w4 := noise.Worley3D(128, 128, 128, 7)
	pw1 := combine(p1, w1)
	cloudBaseData := mergeColorChannels(pw1, w2, w3, w4)
	cloudBaseImage, err := image3d.MakeFromData(128, 128, 128, cloudBaseData)
	if err != nil {
		panic(err)
	}
	cloudBaseImage.SaveToPath(TEX_PATH + "cloud-base/base.png")

	// create cloud detail texture
	fmt.Println("Creating Cloud Detail")
	f1 := noise.Worley3D(32, 32, 32, 5)
	f2 := noise.Worley3D(32, 32, 32, 6)
	f3 := noise.Worley3D(32, 32, 32, 7)
	cloudDetailData := mergeColorChannels(f1, f2, f3)
	cloudDetailImage, err := image3d.MakeFromData(32, 32, 32, cloudDetailData)
	if err != nil {
		panic(err)
	}
	cloudDetailImage.SaveToPath(TEX_PATH + "cloud-detail/detail.png")

	// create cloud turbulence texture
	fmt.Println("Creating Cloud Turbulence")
	c1 := noise.Curl2D(128, 128, 5)
	c2 := noise.Curl2D(128, 128, 6)
	c3 := noise.Curl2D(128, 128, 7)
	cloudTurbulenceData := mergeColorChannels(c1, c2, c3)
	cloudTurbulenceImage, err := image2d.MakeFromData(128, 128, cloudTurbulenceData)
	if err != nil {
		panic(err)
	}
	cloudTurbulenceImage.SaveToPath(TEX_PATH + "cloud-turbulence/turbulence.png")
}
