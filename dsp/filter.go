package dsp

import (
	"math"
)

// GenerateLowPassKernel generates a low-pass FIR filter kernel using a windowed-sinc method.
// cutoffFreq is in Hz, sampleRate in Hz, and numTaps is the filter length (should be odd).
func GenerateLowPassKernel(cutoffFreq float64, sampleRate int, numTaps int) []float64 {
	kernel := make([]float64, numTaps)
	fc := cutoffFreq / float64(sampleRate)
	m := float64(numTaps - 1)

	for i := 0; i < numTaps; i++ {
		n := float64(i)
		if n == m/2 {
			kernel[i] = 2 * fc
		} else {
			kernel[i] = math.Sin(2*math.Pi*fc*(n-m/2)) / (math.Pi * (n - m/2))
		}
		kernel[i] *= 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/m)
	}

	sum := 0.0
	for i := 0; i < numTaps; i++ {
		sum += kernel[i]
	}
	for i := 0; i < numTaps; i++ {
		kernel[i] /= sum
	}
	return kernel
}
