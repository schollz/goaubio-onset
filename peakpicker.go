package onset

// PeakPicker represents a peak picking object for onset detection
type PeakPicker struct {
	Threshold      float64
	WinPost        uint
	WinPre         uint
	Biquad         *Filter
	OnsetKeep      *Fvec
	OnsetProc      *Fvec
	OnsetPeek      *Fvec
	Thresholded    *Fvec
	Scratch        *Fvec
}

// NewPeakPicker creates a new peak picker
func NewPeakPicker() *PeakPicker {
	p := &PeakPicker{
		Threshold: 0.1,
		WinPost:   5,
		WinPre:    1,
	}

	bufSize := p.WinPost + p.WinPre + 1
	p.Scratch = NewFvec(bufSize)
	p.OnsetKeep = NewFvec(bufSize)
	p.OnsetProc = NewFvec(bufSize)
	p.OnsetPeek = NewFvec(3)
	p.Thresholded = NewFvec(1)

	// Create biquad lowpass filter
	// Coefficients from aubio: butter(2, 0.34)
	p.Biquad = NewBiquadFilter(0.15998789, 0.31997577, 0.15998789, 0.23484048, 0)

	return p
}

// Do performs peak picking on the onset detection function
func (p *PeakPicker) Do(onset *Fvec, out *Fvec) {
	// Push new novelty to the end
	FvecPush(p.OnsetKeep, onset.Data[0])

	// Store a copy
	p.OnsetProc.Copy(p.OnsetKeep)

	// Filter this copy
	p.Biquad.DoFiltFilt(p.OnsetProc, p.Scratch)

	// Calculate mean
	mean := FvecMean(p.OnsetProc)

	// Calculate median
	p.Scratch.Copy(p.OnsetProc)
	median := FvecMedian(p.Scratch)

	// Shift peek array
	for j := uint(0); j < 2; j++ {
		p.OnsetPeek.Data[j] = p.OnsetPeek.Data[j+1]
	}

	// Calculate new thresholded value
	p.Thresholded.Data[0] = p.OnsetProc.Data[p.WinPost] - median - mean*p.Threshold
	p.OnsetPeek.Data[2] = p.Thresholded.Data[0]

	// Check for peak
	if FvecPeakPick(p.OnsetPeek, 1) {
		out.Data[0] = FvecQuadraticPeakPos(p.OnsetPeek, 1)
	} else {
		out.Data[0] = 0
	}
}

// SetThreshold sets the peak picking threshold
func (p *PeakPicker) SetThreshold(threshold float64) {
	p.Threshold = threshold
}

// GetThreshold gets the peak picking threshold
func (p *PeakPicker) GetThreshold() float64 {
	return p.Threshold
}

// GetThresholdedInput returns the thresholded input
func (p *PeakPicker) GetThresholdedInput() *Fvec {
	return p.Thresholded
}
