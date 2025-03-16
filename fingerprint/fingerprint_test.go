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
	frames := [][]float64{
		{1, 1, 1, 1},
		{2, 2, 2, 2},
	}
	window := []float64{2, 2, 2, 2}

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
				{1.0, 2.0, 3.0, 2.5},
			},
			numBands: 2,
			expected: []Peak{
				{FrameIndex: 0, FreqBin: 1, Magnitude: 2.0},
				{FrameIndex: 0, FreqBin: 2, Magnitude: 3.0},
			},
		},
		{
			name: "two frames, three bands, frame length 6",
			spectrogram: [][]float64{
				{5, 4, 1, 7, 3, 6},
				{2, 8, 9, 3, 10, 1},
			},
			numBands: 3,
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
				{3, 1, 2, 4, 1},
			},
			numBands: 2,
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
			targetZone: 10,
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
				{FrameIndex: 20000, FreqBin: 20, Magnitude: 0},
			},
			targetZone: 30000,
			expected: []uint32{
				(10 << 23) | (20 << 14) | 0x3FFF,
			},
		},
		{
			name: "frequency clipping test",
			peaks: []Peak{
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
