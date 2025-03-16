package fingerprint

// Peak represents a spectral peak in a frame.
type Peak struct {
	FrameIndex int     // Index of the frame.
	FreqBin    int     // Frequency bin index.
	Magnitude  float64 // Magnitude of the peak.
}

// DetectPeaks finds the strongest frequency peaks in each band of the spectrogram
func DetectPeaks(spectrogram [][]float64, numBands int) []Peak {
	var peaks []Peak
	for i, frame := range spectrogram {
		numBins := len(frame)
		bandSize := numBins / numBands

		for band := 0; band < numBands; band++ {
			start := band * bandSize
			end := start + bandSize
			if band == numBands-1 {
				end = numBins
			}
			maxVal := -1.0
			maxBin := -1
			for j := start; j < end; j++ {
				if frame[j] > maxVal {
					maxVal = frame[j]
					maxBin = j
				}
			}
			if maxBin != -1 {
				peaks = append(peaks, Peak{
					FrameIndex: i,
					FreqBin:    maxBin,
					Magnitude:  maxVal,
				})
			}
		}
	}
	return peaks
}
