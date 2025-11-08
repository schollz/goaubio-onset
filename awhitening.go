package onset

import "math"

const (
	spectralWhiteningDefaultRelaxTime = 250.0  // in seconds
	spectralWhiteningDefaultDecay     = 0.001  // -60dB attenuation
	spectralWhiteningDefaultFloor     = 1.0e-4 // floor value
)

// SpectralWhitening represents an adaptive spectral whitening object
type SpectralWhitening struct {
	BufSize    uint
	HopSize    uint
	Samplerate uint
	RelaxTime  float64
	RDecay     float64
	Floor      float64
	PeakValues *Fvec
}

// NewSpectralWhitening creates a new spectral whitening object
func NewSpectralWhitening(bufSize, hopSize, samplerate uint) *SpectralWhitening {
	s := &SpectralWhitening{
		BufSize:    bufSize,
		HopSize:    hopSize,
		Samplerate: samplerate,
		Floor:      spectralWhiteningDefaultFloor,
		PeakValues: NewFvec(bufSize/2 + 1),
	}
	s.SetRelaxTime(spectralWhiteningDefaultRelaxTime)
	s.Reset()
	return s
}

// Do applies spectral whitening to the FFT grain
func (s *SpectralWhitening) Do(fftgrain *Cvec) {
	length := fftgrain.Length
	if s.PeakValues.Length < length {
		length = s.PeakValues.Length
	}

	for i := uint(0); i < length; i++ {
		tmp := math.Max(s.RDecay*s.PeakValues.Data[i], s.Floor)
		s.PeakValues.Data[i] = math.Max(fftgrain.Norm[i], tmp)
		if s.PeakValues.Data[i] > 0 {
			fftgrain.Norm[i] /= s.PeakValues.Data[i]
		}
	}
}

// SetRelaxTime sets the relax time for spectral whitening
func (s *SpectralWhitening) SetRelaxTime(relaxTime float64) {
	s.RelaxTime = relaxTime
	s.RDecay = math.Pow(spectralWhiteningDefaultDecay,
		(float64(s.HopSize)/float64(s.Samplerate))/s.RelaxTime)
}

// GetRelaxTime gets the relax time
func (s *SpectralWhitening) GetRelaxTime() float64 {
	return s.RelaxTime
}

// SetFloor sets the floor value
func (s *SpectralWhitening) SetFloor(floor float64) {
	s.Floor = floor
}

// GetFloor gets the floor value
func (s *SpectralWhitening) GetFloor() float64 {
	return s.Floor
}

// Reset resets the spectral whitening state
func (s *SpectralWhitening) Reset() {
	for i := range s.PeakValues.Data {
		s.PeakValues.Data[i] = s.Floor
	}
}
