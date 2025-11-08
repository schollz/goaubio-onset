package onset

// Filter represents a digital filter
type Filter struct {
	Order uint
	A     []float64 // feedback coefficients
	B     []float64 // feedforward coefficients
	X     []float64 // input history
	Y     []float64 // output history
}

// NewFilter creates a new filter with given order
func NewFilter(order uint) *Filter {
	f := &Filter{
		Order: order,
		A:     make([]float64, order),
		B:     make([]float64, order),
		X:     make([]float64, order),
		Y:     make([]float64, order),
	}
	// Set default to identity
	f.A[0] = 1.0
	f.B[0] = 1.0
	return f
}

// NewBiquadFilter creates a biquad filter with given coefficients
func NewBiquadFilter(b0, b1, b2, a1, a2 float64) *Filter {
	f := NewFilter(3)
	f.B[0] = b0
	f.B[1] = b1
	f.B[2] = b2
	f.A[0] = 1.0
	f.A[1] = a1
	f.A[2] = a2
	return f
}

// Do applies the filter to the input vector in-place
func (f *Filter) Do(in *Fvec) {
	for j := uint(0); j < in.Length; j++ {
		// New input
		f.X[0] = in.Data[j]
		f.Y[0] = f.B[0] * f.X[0]

		// Apply filter
		for l := uint(1); l < f.Order; l++ {
			f.Y[0] += f.B[l] * f.X[l]
			f.Y[0] -= f.A[l] * f.Y[l]
		}

		// New output
		in.Data[j] = f.Y[0]

		// Store for next sample
		for l := f.Order - 1; l > 0; l-- {
			f.X[l] = f.X[l-1]
			f.Y[l] = f.Y[l-1]
		}
	}
}

// DoFiltFilt applies the filter forward and backward to avoid phase distortion
func (f *Filter) DoFiltFilt(in *Fvec, tmp *Fvec) {
	length := in.Length

	// Apply filtering forward
	f.Do(in)
	f.Reset()

	// Mirror the signal
	for j := uint(0); j < length; j++ {
		tmp.Data[length-j-1] = in.Data[j]
	}

	// Apply filtering on mirrored signal
	f.Do(tmp)
	f.Reset()

	// Invert back
	for j := uint(0); j < length; j++ {
		in.Data[j] = tmp.Data[length-j-1]
	}
}

// Reset clears the filter history
func (f *Filter) Reset() {
	for i := range f.X {
		f.X[i] = 0
		f.Y[i] = 0
	}
}
