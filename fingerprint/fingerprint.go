package fingerprint

import "errors"

const (
	TargetSampleRate = 11025 // Downsampled rate.
	FilterTaps       = 101   // Number of taps in FIR filter (should be odd).
	FrameSize        = 1024  // Samples per frame.
	HopSize          = 512   // Hop size for overlapping frames.
	NumBands         = 6     // Number of frequency bands for peak detection.
	TargetZoneFrames = 20    // Maximum frame difference for pairing peaks.
)

func Fingerprint(samples []int16, sampleRate int) ([]uint32, error) {
	if sampleRate < TargetSampleRate {
		return nil, errors.New("sample rate is lower than target sample rate")
	}
	n := len(samples)
	floatSamples := make([]float64, n)
	for i, s := range samples {
		floatSamples[i] = float64(s) / 32768.0
	}

	return nil, nil
}
