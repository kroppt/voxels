package settings

// SetFOV sets the vertical field of view.
func (r *Repository) SetFOV(degY float64) {
	r.c.setFOV(degY)
}

// GetFOV gets the vertical field of view.
func (r *Repository) GetFOV() float64 {
	return r.c.getFOV()
}

// SetResolution sets the width and height of the window in pixels.
func (r *Repository) SetResolution(width, height int32) {
	r.c.setResolution(width, height)
}

// GetResolution gets the width and height of the window in pixels.
func (r *Repository) GetResolution() (int32, int32) {
	return r.c.getResolution()
}
