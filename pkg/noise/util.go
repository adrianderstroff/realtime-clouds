package noise

// loop loops the value val between 0 and res-1
func loop(val, res int) int {
	newval := val % res
	if newval < 0 {
		newval = res + newval
	}
	return newval
}
