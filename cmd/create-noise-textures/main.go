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

func createCloudBaseTexture() {
	fmt.Println("Creating Cloud Base Shape")
	p1 := noise.Perlin3D(128, 128, 128, 5, 8)
	w1 := noise.Worley3D(128, 128, 128, 4)
	w2 := createAndFillImage(128*128*128, 0)   //noise.Worley3D(128, 128, 128, 8)
	w3 := createAndFillImage(128*128*128, 0)   //noise.Worley3D(128, 128, 128, 16)
	w4 := createAndFillImage(128*128*128, 255) //noise.Worley3D(128, 128, 128, 32)
	pw1 := remapAll(spread(p1), w1)
	pw1 = w1

	cloudBaseData := mergeColorChannels(pw1, w2, w3, w4)
	cloudBaseImage, err := image3d.MakeFromData(128, 128, 128, cloudBaseData)
	if err != nil {
		panic(err)
	}
	cloudBaseImage.SaveToPath(TEX_PATH + "cloud-base/base.png")
}

func createCloudDetailTexture() {
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
}

func createCloudTurbulenceTexture() {
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

func createCloudMapTexture() {
	fmt.Println("Creating cloud map")
	red := noise.Perlin2D(1024, 1024, 4, 8)
	red = spread(red)

	cp1 := noise.Perlin2D(1024, 1024, 16, 1)
	cp2 := noise.Perlin2D(1024, 1024, 32, 1)
	cp3 := noise.Perlin2D(1024, 1024, 64, 1)
	cw1 := noise.Worley2D(1024, 1024, 16)
	cw2 := noise.Worley2D(1024, 1024, 32)
	cw3 := noise.Worley2D(1024, 1024, 64)
	cpn := combine(cp1, cp2, cp3)
	cwn := combine(cw1, cw2, cw3)
	_ = cpn
	_ = cwn
	green := createAndFillImage(1024*1024, 0) //combine(cpn, cwn)

	blue := createAndFillImage(1024*1024, 0) //noise.Perlin2D(1024, 1024, 8)

	cloudMapData := mergeColorChannels(red, green, blue)
	cloudMapImage, err := image2d.MakeFromData(1024, 1024, cloudMapData)
	if err != nil {
		panic(err)
	}
	cloudMapImage.SaveToPath(TEX_PATH + "cloud-map/cloud-map.png")
}

func main() {
	// create cloud base texture
	createCloudBaseTexture()

	// create cloud detail texture
	createCloudDetailTexture()

	// create cloud turbulence texture
	createCloudTurbulenceTexture()

	// create cloud map
	createCloudMapTexture()
}
