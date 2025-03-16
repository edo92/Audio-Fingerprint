package main

import (
	"fingerprint/fingerprint"
	"fingerprint/wav"
	"fmt"
	"log"
)

func main() {
	audioFile := "audio.wav"

	samples, sampleRate, err := wav.ReadWavFile(audioFile)
	if err != nil {
		log.Fatalf("Failed to read WAV file: %v", err)
	}

	hashes, err := fingerprint.Fingerprint(samples, sampleRate)
	if err != nil {
		log.Fatalf("Failed to generate fingerprint: %v", err)
	}

	fmt.Println(hashes)
}
