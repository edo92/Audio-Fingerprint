package dsp_test

import (
	"fingerprint/dsp"
	"math"
	"testing"
)

// check if two float64 numbers are equal within a small tolerance.
func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestGenerateLowPassKernel(t *testing.T) {
	cutoffFreq := 1000.0
	sampleRate := 48000
	numTaps := 11

	kernel := dsp.GenerateLowPassKernel(cutoffFreq, sampleRate, numTaps)

	// Test 1: Check kernel length
	if len(kernel) != numTaps {
		t.Errorf("expected kernel length %d, got %d", numTaps, len(kernel))
	}

	// Test 2: Check that the kernel elements sum to 1 (normalized)
	sum := 0.0
	for _, v := range kernel {
		sum += v
	}
	if !almostEqual(sum, 1.0, 1e-6) {
		t.Errorf("kernel sum = %f; want 1.0", sum)
	}

	// Test 3: Check that the kernel is symmetric
	for i := 0; i < numTaps/2; i++ {
		if !almostEqual(kernel[i], kernel[numTaps-1-i], 1e-6) {
			t.Errorf("kernel not symmetric at index %d: %f != %f", i, kernel[i], kernel[numTaps-1-i])
		}
	}
}
