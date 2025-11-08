package onset

import (
	"strings"
)

// Onset represents an onset detection object
type Onset struct {
	Pv                  *Pvoc
	Od                  *Specdesc
	Pp                  *PeakPicker
	Fftgrain            *Cvec
	Desc                *Fvec
	Silence             float64
	Minioi              uint
	Delay               uint
	Samplerate          uint
	HopSize             uint
	TotalFrames         uint
	LastOnset           uint
	ApplyCompression    bool
	LambdaCompression   float64
	ApplyAWhitening     bool
	SpectralWhitening   *SpectralWhitening
}

// NewOnset creates a new onset detection object
func NewOnset(onsetMode string, bufSize, hopSize, samplerate uint) *Onset {
	o := &Onset{
		Samplerate:        samplerate,
		HopSize:           hopSize,
		Pv:                NewPvoc(bufSize, hopSize),
		Pp:                NewPeakPicker(),
		Od:                NewSpecdesc(onsetMode, bufSize),
		Fftgrain:          NewCvec(bufSize),
		Desc:              NewFvec(1),
		SpectralWhitening: NewSpectralWhitening(bufSize, hopSize, samplerate),
	}

	o.SetDefaultParameters(onsetMode)
	o.Reset()

	return o
}

// Do processes input and detects onsets
func (o *Onset) Do(input *Fvec, onset *Fvec) {
	isonset := 0.0

	// Phase vocoder
	o.Pv.Do(input, o.Fftgrain)

	// Apply adaptive whitening if enabled
	if o.ApplyAWhitening {
		o.SpectralWhitening.Do(o.Fftgrain)
	}

	// Apply compression if enabled
	if o.ApplyCompression {
		o.Fftgrain.LogMag(o.LambdaCompression)
	}

	// Compute spectral descriptor
	o.Od.Do(o.Fftgrain, o.Desc)

	// Peak picking
	o.Pp.Do(o.Desc, onset)
	isonset = onset.Data[0]

	if isonset > 0 {
		if SilenceDetection(input, o.Silence) {
			// Silent onset, not marking
			isonset = 0
		} else {
			// We have an onset
			newOnset := o.TotalFrames + uint(Round(isonset*float64(o.HopSize)))

			// Check if last onset time was more than minioi ago
			if o.LastOnset+o.Minioi < newOnset {
				// Start of file: make sure (new_onset - delay) >= 0
				if o.LastOnset > 0 && o.Delay > newOnset {
					isonset = 0
				} else {
					o.LastOnset = Max(o.Delay, newOnset)
				}
			} else {
				// Doubled onset, not marking
				isonset = 0
			}
		}
	} else {
		// We are at the beginning of the file
		if o.TotalFrames <= o.Delay {
			// And we don't find silence
			if !SilenceDetection(input, o.Silence) {
				newOnset := o.TotalFrames
				if o.TotalFrames == 0 || o.LastOnset+o.Minioi < newOnset {
					isonset = float64(o.Delay) / float64(o.HopSize)
					o.LastOnset = o.TotalFrames + o.Delay
				}
			}
		}
	}

	onset.Data[0] = isonset
	o.TotalFrames += o.HopSize
}

// GetLast returns the time of the latest onset detected, in samples
func (o *Onset) GetLast() uint {
	if o.Delay > o.LastOnset {
		return 0
	}
	return o.LastOnset - o.Delay
}

// GetLastS returns the time of the latest onset detected, in seconds
func (o *Onset) GetLastS() float64 {
	return float64(o.GetLast()) / float64(o.Samplerate)
}

// GetLastMs returns the time of the latest onset detected, in milliseconds
func (o *Onset) GetLastMs() float64 {
	return o.GetLastS() * 1000.0
}

// SetAWhitening enables or disables adaptive whitening
func (o *Onset) SetAWhitening(enable bool) {
	o.ApplyAWhitening = enable
}

// GetAWhitening returns whether adaptive whitening is enabled
func (o *Onset) GetAWhitening() bool {
	return o.ApplyAWhitening
}

// SetCompression sets the compression lambda value
func (o *Onset) SetCompression(lambda float64) {
	if lambda < 0 {
		return
	}
	o.LambdaCompression = lambda
	o.ApplyCompression = lambda > 0
}

// GetCompression returns the compression lambda value
func (o *Onset) GetCompression() float64 {
	if o.ApplyCompression {
		return o.LambdaCompression
	}
	return 0
}

// SetSilence sets the silence threshold
func (o *Onset) SetSilence(silence float64) {
	o.Silence = silence
}

// GetSilence returns the silence threshold
func (o *Onset) GetSilence() float64 {
	return o.Silence
}

// SetThreshold sets the peak picking threshold
func (o *Onset) SetThreshold(threshold float64) {
	o.Pp.SetThreshold(threshold)
}

// GetThreshold returns the peak picking threshold
func (o *Onset) GetThreshold() float64 {
	return o.Pp.GetThreshold()
}

// SetMinioi sets the minimum inter-onset interval in samples
func (o *Onset) SetMinioi(minioi uint) {
	o.Minioi = minioi
}

// GetMinioi returns the minimum inter-onset interval in samples
func (o *Onset) GetMinioi() uint {
	return o.Minioi
}

// SetMinioiS sets the minimum inter-onset interval in seconds
func (o *Onset) SetMinioiS(minioi float64) {
	o.SetMinioi(uint(Round(minioi * float64(o.Samplerate))))
}

// GetMinioiS returns the minimum inter-onset interval in seconds
func (o *Onset) GetMinioiS() float64 {
	return float64(o.Minioi) / float64(o.Samplerate)
}

// SetMinioiMs sets the minimum inter-onset interval in milliseconds
func (o *Onset) SetMinioiMs(minioi float64) {
	o.SetMinioiS(minioi / 1000.0)
}

// GetMinioiMs returns the minimum inter-onset interval in milliseconds
func (o *Onset) GetMinioiMs() float64 {
	return o.GetMinioiS() * 1000.0
}

// SetDelay sets the constant delay in samples
func (o *Onset) SetDelay(delay uint) {
	o.Delay = delay
}

// GetDelay returns the constant delay in samples
func (o *Onset) GetDelay() uint {
	return o.Delay
}

// SetDelayS sets the constant delay in seconds
func (o *Onset) SetDelayS(delay float64) {
	o.SetDelay(uint(delay * float64(o.Samplerate)))
}

// GetDelayS returns the constant delay in seconds
func (o *Onset) GetDelayS() float64 {
	return float64(o.Delay) / float64(o.Samplerate)
}

// SetDelayMs sets the constant delay in milliseconds
func (o *Onset) SetDelayMs(delay float64) {
	o.SetDelayS(delay / 1000.0)
}

// GetDelayMs returns the constant delay in milliseconds
func (o *Onset) GetDelayMs() float64 {
	return o.GetDelayS() * 1000.0
}

// GetDescriptor returns the current value of the onset detection function
func (o *Onset) GetDescriptor() float64 {
	return o.Desc.Data[0]
}

// GetThresholdedDescriptor returns the thresholded value of the onset detection function
func (o *Onset) GetThresholdedDescriptor() float64 {
	thresholded := o.Pp.GetThresholdedInput()
	return thresholded.Data[0]
}

// Reset resets the onset detection state
func (o *Onset) Reset() {
	o.LastOnset = 0
	o.TotalFrames = 0
}

// SetDefaultParameters sets default parameters based on onset mode
func (o *Onset) SetDefaultParameters(onsetMode string) {
	// Set some default parameters
	o.SetThreshold(0.3)
	o.SetDelay(uint(4.3 * float64(o.HopSize)))
	o.SetMinioiMs(50.0)
	o.SetSilence(-70.0)
	o.SetAWhitening(false)
	o.SetCompression(0.0)

	// Method specific optimizations
	mode := strings.ToLower(onsetMode)
	switch mode {
	case "energy":
		// Use defaults
	case "hfc", "default":
		o.SetThreshold(0.058)
		o.SetCompression(1.0)
	case "complexdomain", "complex":
		o.SetDelay(uint(4.6 * float64(o.HopSize)))
		o.SetThreshold(0.15)
		o.SetAWhitening(true)
		o.SetCompression(1.0)
	case "phase":
		o.SetAWhitening(false)
		o.SetCompression(0.0)
	case "wphase":
		// Use defaults
	case "mkl":
		o.SetThreshold(0.05)
		o.SetAWhitening(true)
		o.SetCompression(0.02)
	case "kl":
		o.SetThreshold(0.35)
		o.SetAWhitening(true)
		o.SetCompression(0.02)
	case "specflux":
		o.SetThreshold(0.18)
		o.SetAWhitening(true)
		o.SpectralWhitening.SetRelaxTime(100)
		o.SpectralWhitening.SetFloor(1.0)
		o.SetCompression(10.0)
	case "specdiff":
		// Use defaults
	case "old_default":
		o.SetThreshold(0.3)
		o.SetMinioiMs(20.0)
		o.SetCompression(0.0)
	}
}
