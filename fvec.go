package onset

import "math"

// Fvec represents a vector of real-valued data (floating point)
type Fvec struct {
	Length uint
	Data   []float64
}

// NewFvec creates a new fvec of specified length
func NewFvec(length uint) *Fvec {
	return &Fvec{
		Length: length,
		Data:   make([]float64, length),
	}
}

// Zeros sets all values in the vector to zero
func (f *Fvec) Zeros() {
	for i := range f.Data {
		f.Data[i] = 0
	}
}

// Set sets a sample value at a given position
func (f *Fvec) Set(position uint, value float64) {
	if position < f.Length {
		f.Data[position] = value
	}
}

// Get gets a sample value at a given position
func (f *Fvec) Get(position uint) float64 {
	if position < f.Length {
		return f.Data[position]
	}
	return 0
}

// Copy copies data from source to this fvec
func (f *Fvec) Copy(source *Fvec) {
	length := f.Length
	if source.Length < length {
		length = source.Length
	}
	copy(f.Data[:length], source.Data[:length])
}

// Mean calculates the mean of the vector
func (f *Fvec) Mean() float64 {
	if f.Length == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range f.Data {
		sum += v
	}
	return sum / float64(f.Length)
}

// Max returns the maximum value in the vector
func (f *Fvec) Max() float64 {
	if f.Length == 0 {
		return 0
	}
	max := f.Data[0]
	for _, v := range f.Data[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Min returns the minimum value in the vector
func (f *Fvec) Min() float64 {
	if f.Length == 0 {
		return 0
	}
	min := f.Data[0]
	for _, v := range f.Data[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Weight multiplies all elements by a scalar
func (f *Fvec) Weight(weight float64) {
	for i := range f.Data {
		f.Data[i] *= weight
	}
}

// WeightedCopy copies data with a weight factor
func (f *Fvec) WeightedCopy(source *Fvec, weight float64) {
	length := f.Length
	if source.Length < length {
		length = source.Length
	}
	for i := uint(0); i < length; i++ {
		f.Data[i] = source.Data[i] * weight
	}
}

// LocalEnergyDB calculates local energy in dB
func (f *Fvec) LocalEnergyDB() float64 {
	energy := 0.0
	for _, v := range f.Data {
		energy += v * v
	}
	if energy > 0 {
		return 10.0 * math.Log10(energy/float64(f.Length))
	}
	return -90.0
}
