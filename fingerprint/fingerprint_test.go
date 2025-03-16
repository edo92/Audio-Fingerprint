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

func TestHashFingerprint(t *testing.T) {
	tests := []struct {
		name       string
		peaks      []Peak
		targetZone int
		expected   []uint32
	}{
		{
			name:       "empty peaks",
			peaks:      []Peak{},
			targetZone: 5,
			expected:   []uint32{},
		},
		{
			name: "single peak, no pair",
			peaks: []Peak{
				{FrameIndex: 10, FreqBin: 100, Magnitude: 0},
			},
			targetZone: 5,
			expected:   []uint32{},
		},
		{
			name: "two peaks valid pair",
			peaks: []Peak{
				{FrameIndex: 10, FreqBin: 100, Magnitude: 0},
				{FrameIndex: 15, FreqBin: 200, Magnitude: 0},
			},
			targetZone: 10, // dt = 5, within targetZone.
			expected: []uint32{
				(100 << 23) | (200 << 14) | 5,
			},
		},
		{
			name: "three peaks multiple pairs",
			peaks: []Peak{
				{FrameIndex: 10, FreqBin: 50, Magnitude: 0},
				{FrameIndex: 12, FreqBin: 300, Magnitude: 0},
				{FrameIndex: 13, FreqBin: 400, Magnitude: 0},
			},
			targetZone: 5,
			// Expected pairs:
			// From peak 0 (frame 10): pair with peak 1 (dt=2) and peak 2 (dt=3).
			// From peak 1 (frame 12): pair with peak 2 (dt=1).
			expected: []uint32{
				(50 << 23) | (300 << 14) | 2,
				(50 << 23) | (400 << 14) | 3,
				(300 << 23) | (400 << 14) | 1,
			},
		},
		{
			name: "target zone filter, no pair if dt > targetZone",
			peaks: []Peak{
				{FrameIndex: 10, FreqBin: 100, Magnitude: 0},
				{FrameIndex: 30, FreqBin: 200, Magnitude: 0},
			},
			targetZone: 5, // dt = 20, exceeds targetZone.
			expected:   []uint32{},
		},
		{
			name: "dt clipping test",
			peaks: []Peak{
				{FrameIndex: 0, FreqBin: 10, Magnitude: 0},
				// dt = 20000, which exceeds 0x3FFF (16383) so should be clipped.
				{FrameIndex: 20000, FreqBin: 20, Magnitude: 0},
			},
			targetZone: 30000, // Allow dt > 0x3FFF so clipping occurs.
			expected: []uint32{
				(10 << 23) | (20 << 14) | 0x3FFF,
			},
		},
		{
			name: "frequency clipping test",
			peaks: []Peak{
				// Frequencies above 511 (0x1FF) should be clipped.
				{FrameIndex: 0, FreqBin: 600, Magnitude: 0},
				{FrameIndex: 1, FreqBin: 700, Magnitude: 0},
			},
			targetZone: 5, // dt = 1
			expected: []uint32{
				(0x1FF << 23) | (0x1FF << 14) | 1,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hashes := HashFingerprint(tc.peaks, tc.targetZone)
			if len(hashes) != len(tc.expected) {
				t.Errorf("expected %d hashes, got %d", len(tc.expected), len(hashes))
				return
			}
			for i, h := range hashes {
				if h != tc.expected[i] {
					t.Errorf("hash %d: expected 0x%08X, got 0x%08X", i, tc.expected[i], h)
				}
			}
		})
	}
}
