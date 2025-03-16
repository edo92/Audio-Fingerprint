package fingerprint

import (
	"math"
	"reflect"
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
	// Test data: two frames with window [2,2,2,2]
	// Frame 1: [1,1,1,1] -> [2,2,2,2]
	// Frame 2: [2,2,2,2] -> [4,4,4,4]
	frames := [][]float64{
		{1, 1, 1, 1},
		{2, 2, 2, 2},
	}
	window := []float64{2, 2, 2, 2}

	// Expected FFT magnitudes: [8,0,0] for frame 1, [16,0,0] for frame 2
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

func TestDetectPeaks(t *testing.T) {
	tests := []struct {
		name        string
		spectrogram [][]float64
		numBands    int
		expected    []Peak
	}{
		{
			name:        "empty spectrogram",
			spectrogram: [][]float64{},
			numBands:    2,
			expected:    []Peak{},
		},
		{
			name: "one frame, two bands, frame length 4",
			spectrogram: [][]float64{
				{1.0, 2.0, 3.0, 2.5}, // Single frame
			},
			numBands: 2,
			// Band 0 (0-1): max=2.0 at index 1
			// Band 1 (2-3): max=3.0 at index 2
			expected: []Peak{
				{FrameIndex: 0, FreqBin: 1, Magnitude: 2.0},
				{FrameIndex: 0, FreqBin: 2, Magnitude: 3.0},
			},
		},
		{
			name: "two frames, three bands, frame length 6",
			spectrogram: [][]float64{
				{5, 4, 1, 7, 3, 6},  // Frame 0
				{2, 8, 9, 3, 10, 1}, // Frame 1
			},
			numBands: 3,
			// Frame 0: Band 0 (0-1): max=5@0, Band 1 (2-3): max=7@3, Band 2 (4-5): max=6@5
			// Frame 1: Band 0 (0-1): max=8@1, Band 1 (2-3): max=9@2, Band 2 (4-5): max=10@4
			expected: []Peak{
				{FrameIndex: 0, FreqBin: 0, Magnitude: 5},
				{FrameIndex: 0, FreqBin: 3, Magnitude: 7},
				{FrameIndex: 0, FreqBin: 5, Magnitude: 6},
				{FrameIndex: 1, FreqBin: 1, Magnitude: 8},
				{FrameIndex: 1, FreqBin: 2, Magnitude: 9},
				{FrameIndex: 1, FreqBin: 4, Magnitude: 10},
			},
		},
		{
			name: "non-divisible frame length",
			spectrogram: [][]float64{
				{3, 1, 2, 4, 1}, // Length 5, 2 bands
			},
			numBands: 2,
			// Band 0 (0-1): max=3@0, Band 1 (2-4): max=4@3
			expected: []Peak{
				{FrameIndex: 0, FreqBin: 0, Magnitude: 3},
				{FrameIndex: 0, FreqBin: 3, Magnitude: 4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			peaks := DetectPeaks(tt.spectrogram, tt.numBands)
			if !reflect.DeepEqual(peaks, tt.expected) {
				t.Errorf("expected peaks %v, got %v", tt.expected, peaks)
			}
		})
	}
}
