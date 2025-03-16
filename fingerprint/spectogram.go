package fingerprint

import (
	"fingerprint/dsp"
	"sync"
)

func computeSpectrogram(frames [][]float64, window []float64) [][]float64 {
	numFrames := len(frames)
	spectrogram := make([][]float64, numFrames)
	var wg sync.WaitGroup

	for i, frame := range frames {
		wg.Add(1)
		go func(i int, frame []float64) {
			defer wg.Done()

			for j := range frame {
				frame[j] *= window[j]
			}

			magnitudes := dsp.ComputeFFT(frame)
			spectrogram[i] = magnitudes
		}(i, frame)
	}
	wg.Wait()
	return spectrogram
}
