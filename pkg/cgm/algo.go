package cgm

// Map maps a value val in the range (smin,smax) to the range (dmin,dmax)
func Map(val, smin, smax, dmin, dmax float32) float32 {
	return (val-smin)/(smax-smin)*(dmax-dmin) + dmin
}

// Clamp restricts a value val to the range (min, max).
// If val is small than min it's set to min if it's bigger than max it's
// set to max respectively.
func Clamp(val, min, max float32) float32 {
	return Min32(Max32(val, min), max)
}
