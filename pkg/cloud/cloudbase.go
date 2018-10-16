package cloud

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/noise"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

func CloudBase(width, height, slices int) (texture.Texture, error) {
	p1 := noise.Perlin3D(width, height, slices, 5)
	w1 := noise.Worley3D(width, height, slices, 5)
	w2 := noise.Worley3D(width, height, slices, 6)
	w3 := noise.Worley3D(width, height, slices, 7)
	w4 := noise.Worley3D(width, height, slices, 7)

	pw1 := combine(p1, w1)

	data := mergeColorChannels(pw1, w2, w3, w4)

	return texture.Make3DFromData(data, width, height, slices, gl.RGBA, gl.RGBA)
}
