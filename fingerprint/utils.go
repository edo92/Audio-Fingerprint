package fingerprint

import "math"

func frameSignal(signal []float64, frameSize int, hopSize int) [][]float64 {
	var frames [][]float64
	n := len(signal)
	for start := 0; start+frameSize <= n; start += hopSize {
		frame := make([]float64, frameSize)
		copy(frame, signal[start:start+frameSize])
		frames = append(frames, frame)
	}
	return frames
}

func hammingWindow(n int) []float64 {
	window := make([]float64, n)
	for i := 0; i < n; i++ {
		window[i] = 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(n-1))
	}
	return window
}
