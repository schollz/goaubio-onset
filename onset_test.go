package onset

import (
	"math"
	"testing"
)

func TestFvecCreation(t *testing.T) {
	v := NewFvec(10)
	if v.Length != 10 {
		t.Errorf("Expected length 10, got %d", v.Length)
	}
	if len(v.Data) != 10 {
		t.Errorf("Expected data length 10, got %d", len(v.Data))
	}
}

func TestFvecOperations(t *testing.T) {
	v := NewFvec(5)
	v.Data[0] = 1.0
	v.Data[1] = 2.0
	v.Data[2] = 3.0
	v.Data[3] = 4.0
	v.Data[4] = 5.0

	mean := v.Mean()
	if mean != 3.0 {
		t.Errorf("Expected mean 3.0, got %f", mean)
	}

	max := v.Max()
	if max != 5.0 {
		t.Errorf("Expected max 5.0, got %f", max)
	}

	min := v.Min()
	if min != 1.0 {
		t.Errorf("Expected min 1.0, got %f", min)
	}
}

func TestCvecCreation(t *testing.T) {
	c := NewCvec(512)
	expectedLength := uint(512/2 + 1)
	if c.Length != expectedLength {
		t.Errorf("Expected length %d, got %d", expectedLength, c.Length)
	}
}

func TestPeakPicker(t *testing.T) {
	pp := NewPeakPicker()
	if pp.Threshold != 0.1 {
		t.Errorf("Expected default threshold 0.1, got %f", pp.Threshold)
	}

	pp.SetThreshold(0.5)
	if pp.GetThreshold() != 0.5 {
		t.Errorf("Expected threshold 0.5, got %f", pp.GetThreshold())
	}
}

func TestSpecdesc(t *testing.T) {
	bufSize := uint(512)
	s := NewSpecdesc("hfc", bufSize)

	if s.OnsetType != OnsetHFC {
		t.Errorf("Expected HFC onset type")
	}

	// Test energy method
	s2 := NewSpecdesc("energy", bufSize)
	if s2.OnsetType != OnsetEnergy {
		t.Errorf("Expected Energy onset type")
	}
}

func TestOnsetCreation(t *testing.T) {
	bufSize := uint(512)
	hopSize := uint(256)
	samplerate := uint(44100)

	o := NewOnset("hfc", bufSize, hopSize, samplerate)

	if o.Samplerate != samplerate {
		t.Errorf("Expected samplerate %d, got %d", samplerate, o.Samplerate)
	}
	if o.HopSize != hopSize {
		t.Errorf("Expected hopSize %d, got %d", hopSize, o.HopSize)
	}
}

func TestOnsetDetection(t *testing.T) {
	bufSize := uint(512)
	hopSize := uint(256)
	samplerate := uint(44100)

	o := NewOnset("hfc", bufSize, hopSize, samplerate)
	input := NewFvec(hopSize)
	output := NewFvec(1)

	// Generate a test signal with a clear onset
	for i := uint(0); i < hopSize; i++ {
		t := float64(i) / float64(samplerate)
		input.Data[i] = math.Sin(2 * math.Pi * 440 * t)
	}

	// Process the input
	o.Do(input, output)

	// The output should be a value (onset detected or not)
	if output.Data[0] < 0 {
		t.Errorf("Onset value should not be negative, got %f", output.Data[0])
	}
}

func TestOnsetMethods(t *testing.T) {
	bufSize := uint(512)
	hopSize := uint(256)
	samplerate := uint(44100)

	methods := []string{"energy", "hfc", "complex", "phase", "specdiff", "kl", "mkl", "specflux"}

	for _, method := range methods {
		o := NewOnset(method, bufSize, hopSize, samplerate)
		input := NewFvec(hopSize)
		output := NewFvec(1)

		// Generate a simple test signal
		for i := uint(0); i < hopSize; i++ {
			input.Data[i] = math.Sin(2 * math.Pi * 440 * float64(i) / float64(samplerate))
		}

		// Should not panic
		o.Do(input, output)
	}
}

func TestOnsetThresholds(t *testing.T) {
	bufSize := uint(512)
	hopSize := uint(256)
	samplerate := uint(44100)

	o := NewOnset("hfc", bufSize, hopSize, samplerate)

	o.SetThreshold(0.5)
	if o.GetThreshold() != 0.5 {
		t.Errorf("Expected threshold 0.5, got %f", o.GetThreshold())
	}

	o.SetSilence(-80.0)
	if o.GetSilence() != -80.0 {
		t.Errorf("Expected silence -80.0, got %f", o.GetSilence())
	}

	o.SetMinioiMs(100.0)
	if o.GetMinioiMs() != 100.0 {
		t.Errorf("Expected minioi 100.0 ms, got %f", o.GetMinioiMs())
	}
}

func TestSpectralWhitening(t *testing.T) {
	bufSize := uint(512)
	hopSize := uint(256)
	samplerate := uint(44100)

	sw := NewSpectralWhitening(bufSize, hopSize, samplerate)
	if sw.BufSize != bufSize {
		t.Errorf("Expected bufSize %d, got %d", bufSize, sw.BufSize)
	}

	sw.SetRelaxTime(100.0)
	if sw.GetRelaxTime() != 100.0 {
		t.Errorf("Expected relax time 100.0, got %f", sw.GetRelaxTime())
	}

	sw.SetFloor(1e-3)
	if sw.GetFloor() != 1e-3 {
		t.Errorf("Expected floor 1e-3, got %f", sw.GetFloor())
	}
}

func TestMedian(t *testing.T) {
	v := NewFvec(5)
	v.Data = []float64{3, 1, 4, 1, 5}

	median := FvecMedian(v)
	if median != 3.0 {
		t.Errorf("Expected median 3.0, got %f", median)
	}
}

func TestPeakDetection(t *testing.T) {
	v := NewFvec(5)
	v.Data = []float64{1, 2, 5, 3, 1}

	if !FvecPeakPick(v, 2) {
		t.Error("Expected peak at position 2")
	}

	if FvecPeakPick(v, 1) {
		t.Error("Did not expect peak at position 1")
	}
}

func TestFilter(t *testing.T) {
	f := NewBiquadFilter(0.15998789, 0.31997577, 0.15998789, 0.23484048, 0)

	if f.Order != 3 {
		t.Errorf("Expected order 3, got %d", f.Order)
	}

	input := NewFvec(10)
	for i := range input.Data {
		input.Data[i] = 1.0
	}

	f.Do(input)

	// Filter should have modified the input
	// Just check it doesn't crash
}

func BenchmarkOnsetDetection(b *testing.B) {
	bufSize := uint(512)
	hopSize := uint(256)
	samplerate := uint(44100)

	o := NewOnset("hfc", bufSize, hopSize, samplerate)
	input := NewFvec(hopSize)
	output := NewFvec(1)

	// Generate test signal
	for i := uint(0); i < hopSize; i++ {
		input.Data[i] = math.Sin(2 * math.Pi * 440 * float64(i) / float64(samplerate))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.Do(input, output)
	}
}
