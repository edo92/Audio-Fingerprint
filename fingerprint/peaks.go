package fingerprint

type Peak struct {
	FrameIndex int
	FreqBin    int
	Magnitude  float64
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
