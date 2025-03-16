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

func TestApplyFIRFilterIdentity(t *testing.T) {
	// Identity filter: using kernel [0, 1, 0] should return the input unchanged.
	input := []float64{1, 2, 3, 4, 5}
	kernel := []float64{0, 1, 0}
	expected := []float64{1, 2, 3, 4, 5}

	output := dsp.ApplyFIRFilter(input, kernel)
	if len(output) != len(expected) {
		t.Fatalf("expected output length %d, got %d", len(expected), len(output))
	}
	for i, v := range output {
		if !almostEqual(v, expected[i], 1e-6) {
			t.Errorf("at index %d: expected %f, got %f", i, expected[i], v)
		}
	}
}

func TestApplyFIRFilterSmoothing(t *testing.T) {
	// Smoothing filter: a 3-point moving average kernel [1/3, 1/3, 1/3].
	input := []float64{1, 2, 3, 4, 5}
	kernel := []float64{1.0 / 3, 1.0 / 3, 1.0 / 3}

	expected := []float64{1.0, 2.0, 3.0, 4.0, 3.0}

	output := dsp.ApplyFIRFilter(input, kernel)
	if len(output) != len(expected) {
		t.Fatalf("expected output length %d, got %d", len(expected), len(output))
	}
	for i, v := range output {
		if !almostEqual(v, expected[i], 1e-6) {
			t.Errorf("at index %d: expected %f, got %f", i, expected[i], v)
		}
	}
}

func TestApplyFIRFilterEmptyInput(t *testing.T) {
	// When input is empty, the output should be an empty slice.
	input := []float64{}
	kernel := []float64{0, 1, 0}
	output := dsp.ApplyFIRFilter(input, kernel)
	if len(output) != 0 {
		t.Errorf("expected empty output, got length %d", len(output))
	}
}

func TestApplyFIRFilterEmptyKernel(t *testing.T) {
	// When the kernel is empty, no accumulation is performed,
	// so the output should be a slice of zeros with the same length as the input.
	input := []float64{1, 2, 3}
	kernel := []float64{}
	expected := []float64{0, 0, 0}

	output := dsp.ApplyFIRFilter(input, kernel)
	if len(output) != len(expected) {
		t.Fatalf("expected output length %d, got %d", len(expected), len(output))
	}
	for i, v := range output {
		if !almostEqual(v, expected[i], 1e-6) {
			t.Errorf("at index %d: expected %f, got %f", i, expected[i], v)
		}
	}
}
