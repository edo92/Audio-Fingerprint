package dsp

import (
	"math/cmplx"

	"gonum.org/v1/gonum/dsp/fourier"
)

// computes the FFT of a real valued frame. returns the magnitude spectrum
// (only the first half is returned since the input is real-valued).
func ComputeFFT(frame []float64) []float64 {
	n := len(frame)
	fft := fourier.NewFFT(n)

	spectrum := fft.Coefficients(nil, frame)
	half := n/2 + 1
	magnitudes := make([]float64, half)
	for i := 0; i < half; i++ {
		magnitudes[i] = cmplx.Abs(spectrum[i])
	}
	return magnitudes
}
