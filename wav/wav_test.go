package audio_test

import (
	audio "fingerprint/wav"
	"os"
	"testing"

	audioWav "github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func TestReadWavFile_NonExistent(t *testing.T) {
	_, _, err := audio.ReadWavFile("nonexistent.wav")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

// Test for invalid WAV file content.
func TestReadWavFile_InvalidFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "invalid*.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write non-WAV content.
	_, err = tmpFile.WriteString("this is not a valid wav file")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	_, _, err = audio.ReadWavFile(tmpFile.Name())
	if err == nil {
		t.Error("expected error for invalid WAV file, got nil")
	}
}

// Test reading a valid mono 16-bit PCM WAV file.
func TestReadWavFile_Mono(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "mono*.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Set up a WAV encoder for a 16-bit, mono PCM file.
	enc := wav.NewEncoder(tmpFile, 44100, 16, 1, 1)
	// Create a simple mono buffer.
	samples := []int{100, 200, -100, -200}
	buf := &audioWav.IntBuffer{
		Format: &audioWav.Format{
			NumChannels: 1,
			SampleRate:  44100,
		},
		Data:           samples,
		SourceBitDepth: 16,
	}

	if err := enc.Write(buf); err != nil {
		t.Fatal(err)
	}
	if err := enc.Close(); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Read back the file.
	readSamples, sampleRate, err := audio.ReadWavFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if sampleRate != 44100 {
		t.Errorf("expected sample rate 44100, got %d", sampleRate)
	}
	if len(readSamples) != len(samples) {
		t.Errorf("expected %d samples, got %d", len(samples), len(readSamples))
	}
	for i, v := range readSamples {
		if int(v) != samples[i] {
			t.Errorf("sample %d: expected %d, got %d", i, samples[i], v)
		}
	}
}
