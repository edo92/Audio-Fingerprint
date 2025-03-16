package main

import (
	"fingerprint/fingerprint"
	audio "fingerprint/wav"
	"fmt"
	"log"
)

func main() {
	audioFile := "assets/audio.wav"

	samples, sampleRate, err := audio.ReadWavFile(audioFile)
	if err != nil {
		log.Fatalf("Failed to read WAV file: %v", err)
	}

	hashes, err := fingerprint.Fingerprint(samples, sampleRate)
	if err != nil {
		log.Fatalf("Failed to generate fingerprint: %v", err)
	}

	fmt.Println(hashes)
}
