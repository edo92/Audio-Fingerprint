package wav_test

import (
	"fingerprint/wav"
	"os"
	"testing"
)

func TestReadWavFile_NonExistent(t *testing.T) {
	_, _, err := wav.ReadWavFile("nonexistent.wav")
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

	_, _, err = wav.ReadWavFile(tmpFile.Name())
	if err == nil {
		t.Error("expected error for invalid WAV file, got nil")
	}
}
