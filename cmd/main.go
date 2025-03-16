package main

import (
	audio "fingerprint/wav"
	"fmt"
	"log"
)

func main() {
	audioFile := "audio.wav"

	samples, _, err := audio.ReadWavFile(audioFile)
	if err != nil {
		log.Fatalf("Failed to read WAV file: %v", err)
	}

	fmt.Println(samples)
}
