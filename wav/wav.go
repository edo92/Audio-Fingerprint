package audio

import (
	"errors"
	"os"

	"github.com/go-audio/wav"
)

// ReadWavFile reads a 16-bit PCM WAV file from the given path,
// returning mono samples as []int16 and the sample rate.
func ReadWavFile(path string) ([]int16, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, 0, errors.New("invalid WAV file")
	}
	// Only support PCM 16-bit.
	if decoder.BitDepth != 16 {
		return nil, 0, errors.New("only 16-bit PCM WAV files are supported")
	}

	// Decode the entire file.
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, 0, err
	}

	samples := make([]int16, 0, len(buf.Data)/buf.Format.NumChannels)
	if buf.Format.NumChannels == 1 {
		for _, v := range buf.Data {
			samples = append(samples, int16(v))
		}
	} else {
		// Downmix stereo to mono by averaging the channels.
		for i := 0; i < len(buf.Data); i += buf.Format.NumChannels {
			sum := 0
			for c := 0; c < buf.Format.NumChannels; c++ {
				sum += int(buf.Data[i+c])
			}
			avg := int16(sum / buf.Format.NumChannels)
			samples = append(samples, avg)
		}
	}
	return samples, buf.Format.SampleRate, nil
}
