package onset

import "math"

// Cvec represents a vector of complex-valued data stored in polar coordinates
// (magnitude and phase)
type Cvec struct {
	Length uint      // length = (original length)/2 + 1
	Norm   []float64 // magnitude array
	Phas   []float64 // phase array
}

// NewCvec creates a new cvec of specified length
// The actual arrays will have size (length/2 + 1)
func NewCvec(length uint) *Cvec {
	size := length/2 + 1
	return &Cvec{
		Length: size,
		Norm:   make([]float64, size),
		Phas:   make([]float64, size),
	}
}

// Zeros sets all norm and phase values to zero
func (c *Cvec) Zeros() {
	for i := range c.Norm {
		c.Norm[i] = 0
		c.Phas[i] = 0
	}
}

// SetNorm sets the norm (magnitude) at a given position
func (c *Cvec) SetNorm(position uint, value float64) {
	if position < c.Length {
		c.Norm[position] = value
	}
}

// GetNorm gets the norm (magnitude) at a given position
func (c *Cvec) GetNorm(position uint) float64 {
	if position < c.Length {
		return c.Norm[position]
	}
	return 0
}

// SetPhas sets the phase at a given position
func (c *Cvec) SetPhas(position uint, value float64) {
	if position < c.Length {
		c.Phas[position] = value
	}
}

// GetPhas gets the phase at a given position
func (c *Cvec) GetPhas(position uint) float64 {
	if position < c.Length {
		return c.Phas[position]
	}
	return 0
}

// Copy copies data from source to this cvec
func (c *Cvec) Copy(source *Cvec) {
	length := c.Length
	if source.Length < length {
		length = source.Length
	}
	copy(c.Norm[:length], source.Norm[:length])
	copy(c.Phas[:length], source.Phas[:length])
}

// LogMag applies logarithmic compression to magnitudes
func (c *Cvec) LogMag(lambda float64) {
	if lambda > 0 {
		for i := range c.Norm {
			c.Norm[i] = math.Log(1.0 + lambda*c.Norm[i])
		}
	}
}
