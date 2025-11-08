package onset

import (
	"math"
	"sort"
)

// SilenceDetection checks if the input is below silence threshold
func SilenceDetection(input *Fvec, threshold float64) bool {
	db := input.LocalEnergyDB()
	return db < threshold
}

// FvecPush pushes a new element to the end of vector, shifting all elements left
func FvecPush(v *Fvec, newElem float64) {
	for i := uint(0); i < v.Length-1; i++ {
		v.Data[i] = v.Data[i+1]
	}
	v.Data[v.Length-1] = newElem
}

// FvecMedian computes the median of a vector (modifies the input)
func FvecMedian(input *Fvec) float64 {
	if input.Length == 0 {
		return 0
	}

	// Create a copy to avoid modifying original
	arr := make([]float64, input.Length)
	copy(arr, input.Data)

	n := len(arr)
	low := 0
	high := n - 1
	median := (low + high) / 2

	for {
		if high <= low {
			return arr[median]
		}

		if high == low+1 {
			if arr[low] > arr[high] {
				arr[low], arr[high] = arr[high], arr[low]
			}
			return arr[median]
		}

		// Find median of low, middle and high items
		middle := (low + high) / 2
		if arr[middle] > arr[high] {
			arr[middle], arr[high] = arr[high], arr[middle]
		}
		if arr[low] > arr[high] {
			arr[low], arr[high] = arr[high], arr[low]
		}
		if arr[middle] > arr[low] {
			arr[middle], arr[low] = arr[low], arr[middle]
		}

		// Swap low item into position (low+1)
		arr[middle], arr[low+1] = arr[low+1], arr[middle]

		// Partition
		ll := low + 1
		hh := high
		for {
			for ll++; arr[low] > arr[ll]; ll++ {
			}
			for hh--; arr[hh] > arr[low]; hh-- {
			}

			if hh < ll {
				break
			}

			arr[ll], arr[hh] = arr[hh], arr[ll]
		}

		// Swap middle item back
		arr[low], arr[hh] = arr[hh], arr[low]

		// Re-set active partition
		if hh <= median {
			low = ll
		}
		if hh >= median {
			high = hh - 1
		}
	}
}

// FvecMean computes the mean of a vector
func FvecMean(input *Fvec) float64 {
	return input.Mean()
}

// FvecPeakPick checks if position pos is a peak
func FvecPeakPick(onset *Fvec, pos uint) bool {
	if pos == 0 || pos >= onset.Length-1 {
		return false
	}
	return onset.Data[pos] > onset.Data[pos-1] &&
		onset.Data[pos] > onset.Data[pos+1] &&
		onset.Data[pos] > 0
}

// FvecQuadraticPeakPos finds the quadratic interpolated peak position
func FvecQuadraticPeakPos(x *Fvec, pos uint) float64 {
	if pos == 0 || pos == x.Length-1 {
		return float64(pos)
	}

	x0 := pos
	if pos >= 1 {
		x0 = pos - 1
	}

	x2 := pos + 1
	if x2 >= x.Length {
		x2 = pos
	}

	if x0 == pos {
		if x.Data[pos] <= x.Data[x2] {
			return float64(pos)
		}
		return float64(x2)
	}

	if x2 == pos {
		if x.Data[pos] <= x.Data[x0] {
			return float64(pos)
		}
		return float64(x0)
	}

	s0 := x.Data[x0]
	s1 := x.Data[pos]
	s2 := x.Data[x2]

	return float64(pos) + 0.5*(s0-s2)/(s0-2.0*s1+s2)
}

// MedianSimple is a simpler median implementation using sort
func MedianSimple(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

// Max returns the maximum of two values
func Max(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

// Round rounds a float64 to the nearest integer
func Round(x float64) int {
	return int(math.Floor(x + 0.5))
}
