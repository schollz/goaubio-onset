package onset

import (
	"math"
	"strings"
)

// SpecdescType represents the type of spectral descriptor
type SpecdescType int

const (
	OnsetEnergy SpecdescType = iota
	OnsetSpecdiff
	OnsetHFC
	OnsetComplex
	OnsetPhase
	OnsetWPhase
	OnsetKL
	OnsetMKL
	OnsetSpecflux
)

// Specdesc represents a spectral descriptor for onset detection
type Specdesc struct {
	OnsetType SpecdescType
	Threshold float64
	OldMag    *Fvec
	Dev1      *Fvec
	Theta1    *Fvec
	Theta2    *Fvec
}

// NewSpecdesc creates a new spectral descriptor
func NewSpecdesc(onsetMode string, size uint) *Specdesc {
	rsize := size/2 + 1
	s := &Specdesc{
		Threshold: 0.1,
		OldMag:    NewFvec(rsize),
		Dev1:      NewFvec(rsize),
		Theta1:    NewFvec(rsize),
		Theta2:    NewFvec(rsize),
	}

	// Determine onset type from mode string
	mode := strings.ToLower(onsetMode)
	switch mode {
	case "energy":
		s.OnsetType = OnsetEnergy
	case "specdiff":
		s.OnsetType = OnsetSpecdiff
	case "hfc", "default":
		s.OnsetType = OnsetHFC
	case "complexdomain", "complex":
		s.OnsetType = OnsetComplex
	case "phase":
		s.OnsetType = OnsetPhase
	case "wphase":
		s.OnsetType = OnsetWPhase
	case "kl":
		s.OnsetType = OnsetKL
	case "mkl":
		s.OnsetType = OnsetMKL
	case "specflux":
		s.OnsetType = OnsetSpecflux
	default:
		s.OnsetType = OnsetHFC
	}

	return s
}

// Do computes the spectral descriptor
func (s *Specdesc) Do(fftgrain *Cvec, onset *Fvec) {
	switch s.OnsetType {
	case OnsetEnergy:
		s.energy(fftgrain, onset)
	case OnsetHFC:
		s.hfc(fftgrain, onset)
	case OnsetComplex:
		s.complex(fftgrain, onset)
	case OnsetPhase:
		s.phase(fftgrain, onset)
	case OnsetWPhase:
		s.wphase(fftgrain, onset)
	case OnsetSpecdiff:
		s.specdiff(fftgrain, onset)
	case OnsetKL:
		s.kl(fftgrain, onset)
	case OnsetMKL:
		s.mkl(fftgrain, onset)
	case OnsetSpecflux:
		s.specflux(fftgrain, onset)
	default:
		s.hfc(fftgrain, onset)
	}
}

// energy computes energy-based onset detection
func (s *Specdesc) energy(fftgrain *Cvec, onset *Fvec) {
	onset.Data[0] = 0.0
	for j := uint(0); j < fftgrain.Length; j++ {
		onset.Data[0] += fftgrain.Norm[j] * fftgrain.Norm[j]
	}
}

// hfc computes High Frequency Content onset detection
func (s *Specdesc) hfc(fftgrain *Cvec, onset *Fvec) {
	onset.Data[0] = 0.0
	for j := uint(0); j < fftgrain.Length; j++ {
		onset.Data[0] += float64(j+1) * fftgrain.Norm[j]
	}
}

// complex computes Complex Domain onset detection
func (s *Specdesc) complex(fftgrain *Cvec, onset *Fvec) {
	onset.Data[0] = 0.0
	for j := uint(0); j < fftgrain.Length; j++ {
		// Predict phase
		s.Dev1.Data[j] = 2.0*s.Theta1.Data[j] - s.Theta2.Data[j]

		// Euclidean distance in complex domain
		dev := s.Dev1.Data[j] - fftgrain.Phas[j]
		val := s.OldMag.Data[j]*s.OldMag.Data[j] +
			fftgrain.Norm[j]*fftgrain.Norm[j] -
			2.0*s.OldMag.Data[j]*fftgrain.Norm[j]*math.Cos(dev)

		if val > 0 {
			onset.Data[0] += math.Sqrt(val)
		}

		// Store old phase data
		s.Theta2.Data[j] = s.Theta1.Data[j]
		s.Theta1.Data[j] = fftgrain.Phas[j]
		s.OldMag.Data[j] = fftgrain.Norm[j]
	}
}

// phase computes Phase-based onset detection
func (s *Specdesc) phase(fftgrain *Cvec, onset *Fvec) {
	onset.Data[0] = 0.0
	for j := uint(0); j < fftgrain.Length; j++ {
		dev := math.Abs(fftgrain.Phas[j] - s.Theta1.Data[j])
		if s.Threshold < fftgrain.Norm[j] {
			onset.Data[0] += dev
		}
		s.Theta1.Data[j] = fftgrain.Phas[j]
	}
}

// wphase computes Weighted Phase Deviation onset detection
func (s *Specdesc) wphase(fftgrain *Cvec, onset *Fvec) {
	onset.Data[0] = 0.0
	for j := uint(0); j < fftgrain.Length; j++ {
		dev := math.Abs(fftgrain.Phas[j] - s.Theta1.Data[j])
		if s.Threshold < fftgrain.Norm[j] {
			onset.Data[0] += fftgrain.Norm[j] * dev
		}
		s.Theta1.Data[j] = fftgrain.Phas[j]
	}
}

// specdiff computes Spectral Difference onset detection
func (s *Specdesc) specdiff(fftgrain *Cvec, onset *Fvec) {
	onset.Data[0] = 0.0
	for j := uint(0); j < fftgrain.Length; j++ {
		val := fftgrain.Norm[j]*fftgrain.Norm[j] - s.OldMag.Data[j]*s.OldMag.Data[j]
		if val > 0 {
			s.Dev1.Data[j] = math.Sqrt(val)
		} else {
			s.Dev1.Data[j] = 0.0
		}

		if s.Threshold < fftgrain.Norm[j] {
			onset.Data[0] += math.Abs(s.Dev1.Data[j])
		}
		s.OldMag.Data[j] = fftgrain.Norm[j]
	}
}

// kl computes Kullback-Liebler onset detection
func (s *Specdesc) kl(fftgrain *Cvec, onset *Fvec) {
	onset.Data[0] = 0.0
	for j := uint(0); j < fftgrain.Length; j++ {
		onset.Data[0] += fftgrain.Norm[j] *
			math.Log(1.0+fftgrain.Norm[j]/(s.OldMag.Data[j]+1e-1))
		s.OldMag.Data[j] = fftgrain.Norm[j]
	}
	if math.IsNaN(onset.Data[0]) {
		onset.Data[0] = 0.0
	}
}

// mkl computes Modified Kullback-Liebler onset detection
func (s *Specdesc) mkl(fftgrain *Cvec, onset *Fvec) {
	onset.Data[0] = 0.0
	for j := uint(0); j < fftgrain.Length; j++ {
		onset.Data[0] += math.Log(1.0 + fftgrain.Norm[j]/(s.OldMag.Data[j]+1e-1))
		s.OldMag.Data[j] = fftgrain.Norm[j]
	}
	if math.IsNaN(onset.Data[0]) {
		onset.Data[0] = 0.0
	}
}

// specflux computes Spectral Flux onset detection
func (s *Specdesc) specflux(fftgrain *Cvec, onset *Fvec) {
	onset.Data[0] = 0.0
	for j := uint(0); j < fftgrain.Length; j++ {
		if fftgrain.Norm[j] > s.OldMag.Data[j] {
			onset.Data[0] += fftgrain.Norm[j] - s.OldMag.Data[j]
		}
		s.OldMag.Data[j] = fftgrain.Norm[j]
	}
}
