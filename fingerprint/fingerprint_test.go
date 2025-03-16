package fingerprint

import (
	"math"
	"testing"
)

// almostEqualSlices compares two float64 slices element by element within a tolerance.
func almostEqualSlices(a, b []float64, tol float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > tol {
			return false
		}
	}
	return true
}

func TestComputeSpectrogram(t *testing.T) {
	// Two frames of length 4, with window [2,2,2,2]:
	// Frame 1: [1,1,1,1] -> [2,2,2,2]
	// Frame 2: [2,2,2,2] -> [4,4,4,4]
	frames := [][]float64{
		{1, 1, 1, 1},
		{2, 2, 2, 2},
	}
	window := []float64{2, 2, 2, 2}

	// For constant signals, FFT gives magnitude spectrum of length 3:
	// Frame 1: [2,2,2,2] -> FFT magnitudes: [8,0,0] (DC=8)
	// Frame 2: [4,4,4,4] -> FFT magnitudes: [16,0,0] (DC=16)
	expected := [][]float64{
		{8, 0, 0},
		{16, 0, 0},
	}

	spectrogram := computeSpectrogram(frames, window)

	if len(spectrogram) != len(expected) {
		t.Fatalf("expected %d frames, got %d", len(expected), len(spectrogram))
	}

	tol := 1e-6
	for i := range expected {
		if !almostEqualSlices(spectrogram[i], expected[i], tol) {
			t.Errorf("frame %d: expected %v, got %v", i, expected[i], spectrogram[i])
		}
	}
}
