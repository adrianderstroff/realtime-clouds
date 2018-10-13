package cloud

import (
	"github.com/adrianderstroff/realtime-clouds/pkg/core/gl"
	"github.com/adrianderstroff/realtime-clouds/pkg/noise"
	"github.com/adrianderstroff/realtime-clouds/pkg/view/texture"
)

func CloudDetail(width, height, slices int) (texture.Texture, error) {
	f1 := noise.Worley3D(width, height, slices, 5)
	f2 := noise.Worley3D(width, height, slices, 6)
	f3 := noise.Worley3D(width, height, slices, 7)

	data := mergeColorChannels(f1, f2, f3)

	return texture.Make3DFromData(data, width, height, slices, gl.RGB, gl.RGB)
}
