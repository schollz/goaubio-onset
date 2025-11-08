package onset

import (
	"math"

	"github.com/mjibson/go-dsp/fft"
)

// Pvoc represents a phase vocoder
type Pvoc struct {
	WinSize  uint      // window size
	HopSize  uint      // hop size
	Fft      *Fvec     // FFT object
	Window   *Fvec     // analysis window
	Synth    *Fvec     // synthesis window
	In       *Fvec     // input buffer
	Out      *Fvec     // output buffer
	Grain    *Cvec     // current grain (FFT output)
	OldGrain *Cvec     // previous grain
	PrevPhas []float64 // previous phase values
}

// NewPvoc creates a new phase vocoder
func NewPvoc(winSize, hopSize uint) *Pvoc {
	p := &Pvoc{
		WinSize:  winSize,
		HopSize:  hopSize,
		Fft:      NewFvec(winSize),
		Window:   NewFvec(winSize),
		In:       NewFvec(hopSize),
		Grain:    NewCvec(winSize),
		OldGrain: NewCvec(winSize),
		PrevPhas: make([]float64, winSize/2+1),
	}

	// Create Hann window
	for i := uint(0); i < winSize; i++ {
		p.Window.Data[i] = 0.5 - 0.5*math.Cos(2.0*math.Pi*float64(i)/float64(winSize))
	}

	return p
}

// Do processes input through phase vocoder
func (p *Pvoc) Do(input *Fvec, fftgrain *Cvec) {
	// Copy input to FFT buffer with windowing
	for i := uint(0); i < p.WinSize; i++ {
		if i < input.Length {
			p.Fft.Data[i] = input.Data[i] * p.Window.Data[i]
		} else {
			p.Fft.Data[i] = 0
		}
	}

	// Perform FFT
	fftResult := fft.FFTReal(p.Fft.Data)

	// Convert to polar form (magnitude and phase)
	for i := uint(0); i < fftgrain.Length; i++ {
		real := real(fftResult[i])
		imag := imag(fftResult[i])
		fftgrain.Norm[i] = math.Sqrt(real*real + imag*imag)
		fftgrain.Phas[i] = math.Atan2(imag, real)
	}
}

// RDo performs inverse phase vocoder operation (not needed for onset detection)
func (p *Pvoc) RDo(fftgrain *Cvec, output *Fvec) {
	// Not implemented as it's not needed for onset detection
}
