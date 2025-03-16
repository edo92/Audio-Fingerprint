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

// Test reading a valid stereo 16-bit PCM WAV file and downmixing to mono.
func TestReadWavFile_Stereo(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "stereo*.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Set up a WAV encoder for a 16-bit, stereo PCM file.
	enc := wav.NewEncoder(tmpFile, 48000, 16, 2, 1)

	samples := []int{
		100, 300,
		200, 400,
		-100, -300,
		-200, -400,
	}
	buf := &audioWav.IntBuffer{
		Format: &audioWav.Format{
			NumChannels: 2,
			SampleRate:  48000,
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

	readSamples, sampleRate, err := audio.ReadWavFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if sampleRate != 48000 {
		t.Errorf("expected sample rate 48000, got %d", sampleRate)
	}

	// Expected downmixed samples.
	expected := []int16{200, 300, -200, -300}
	if len(readSamples) != len(expected) {
		t.Errorf("expected %d samples, got %d", len(expected), len(readSamples))
	}
	for i, v := range readSamples {
		if v != expected[i] {
			t.Errorf("sample %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

// Test reading a WAV file with an unsupported bit depth (e.g., 8-bit).
func TestReadWavFile_WrongBitDepth(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "assets/audio*.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Set up a WAV encoder with 8-bit depth.
	enc := wav.NewEncoder(tmpFile, 44100, 8, 1, 1)
	samples := []int{10, 20, 30, 40}
	buf := &audioWav.IntBuffer{
		Format: &audioWav.Format{
			NumChannels: 1,
			SampleRate:  44100,
		},
		Data:           samples,
		SourceBitDepth: 8,
	}

	if err := enc.Write(buf); err != nil {
		t.Fatal(err)
	}
	if err := enc.Close(); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	_, _, err = audio.ReadWavFile(tmpFile.Name())
	if err == nil {
		t.Error("expected error for 8-bit WAV file, got nil")
	}
}
