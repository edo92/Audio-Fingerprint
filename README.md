# Audio Fingerprinting Library

A Go implementation of an audio fingerprinting system that generates robust fingerprints from audio files.

## Overview

This library implements an audio fingerprinting algorithm that can generate robust hashes from audio files. The generated fingerprints are designed to be:

- Robust: Resistant to noise, distortion, and other audio transformations
- Compact: Represented as 32-bit integers for efficient storage and lookup
- Distinctive: Able to uniquely identify a specific piece of audio

The fingerprinting approach is inspired by the Shazam algorithm, where constellations of audio peaks in the time-frequency domain are used to create fingerprints.

## Quick Start

[!NOTE] **audio file needs to be 16-bit PCM format**

1. Clone the repository:

```bash
git clone https://github.com/edo92/Audio-Fingerprint
cd audio-fingerprint
```

2. Build the application:

```bash
go build -o audio-fp ./cmd/main.go
```

3. Make the binary executable (Linux/Mac):

```bash
chmod +x audio-fp
```

4. Verify the installation:

```bash
./audio-fp --version
```

5. Run dev

```bash
go run ./cmd/main.go
```

## Architecture

The system is organized into three main components:

- Audio Processing: Reading and preprocessing audio files
- DSP (Digital Signal Processing): Signal processing functions like filtering and FFT
- Fingerprinting: The core fingerprinting algorithm

## Package Structure

```
audio-fingerprint/
├── cmd/
│   └── main.go           # Command line entry point
├── dsp/
│   ├── fft.go            # Fast Fourier Transform implementation
│   ├── filter.go         # FIR filter implementation
│   └── dsp_test.go       # DSP unit tests
├── fingerprint/
│   ├── fingerprint.go    # Main fingerprinting algorithm
│   ├── hash.go           # Hash generation from audio peaks
│   ├── peaks.go          # Peak detection in spectrogram
│   ├── spectogram.go     # Spectrogram computation
│   ├── utils.go          # Utility functions
│   └── fingerprint_test.go # Fingerprinting unit tests
├── wav/
│   ├── wav.go            # WAV file reading
│   └── wav_test.go       # WAV file reading unit tests
└── go.mod                # Go module definition
```

## How It Works

### The Fingerprinting Process

The audio fingerprinting process follows these steps:

1. **Read Audio File**: Read WAV file and convert to mono if necessary

2. **Preprocessing**:

- Low-pass filter the signal
- Downsample to a target sample rate (11025 Hz)

3. **Framing**: Divide the signal into overlapping frames

4. **Spectral Analysis**:

- Apply a Hamming window to each frame
- Compute FFT to get the frequency spectrum

5. **Peak Detection**:

- Divide the spectrum into frequency bands
- Find the strongest peak in each band

6. **Fingerprint Generation**:

- Create pairs of peaks within a target zone
- Generate 32-bit hash values from each pair

## Technical Details

### Signal Preprocessing

The system first reads the input WAV file (currently supporting 16-bit PCM) and converts it to a mono signal if necessary. A low-pass filter is applied to remove high-frequency components that might be affected by noise or compression artifacts. The signal is then downsampled to reduce computational complexity.

### Framing and Spectral Analysis

The preprocessed signal is divided into overlapping frames (1024 samples each with 512 sample overlap). Each frame is multiplied by a Hamming window to reduce spectral leakage, and then the FFT is computed to obtain the frequency spectrum.

### Peak Finding

For each frame, the frequency spectrum is divided into bands (default: 6 bands). Within each band, the strongest peak is identified. These peaks form a constellation of points in the time-frequency domain that characterize the audio.

### Hash Generation

Fingerprints are created by pairing peaks within a target zone (default: 20 frames). For each pair, a 32-bit hash is generated:

- 9 bits: Anchor peak frequency bin (0-511)
- 9 bits: Target peak frequency bin (0-511)
- 14 bits: Time delta between peaks (0-16383)

These hashes can be used for audio identification by matching against a database of known fingerprints.

## API Reference

### wav package

`ReadWavFile(path string) ([]int16, int, error)`

Reads a WAV file and returns the audio samples and sample rate.

- **Parameters**:

  - path: Path to the WAV file

- **Returns**:
  - []int16: Audio samples
  - int: Sample rate in Hz
  - error: Error if any

## dsp package

`GenerateLowPassKernel(cutoffFreq float64, sampleRate int, numTaps int) []float64`

Creates a low-pass FIR filter kernel.

- Parameters:

  - cutoffFreq: Cutoff frequency in Hz
  - sampleRate: Sample rate in Hz
  - numTaps: Filter length (must be odd)

- Returns:
  - []float64: Filter kernel coefficients

`ApplyFIRFilter(input []float64, kernel []float64) []float64`
Applies an FIR filter to the input signal.

- Parameters:

  - input: Input signal
  - kernel: Filter kernel

- Returns:
  - []float64: Filtered signal

`ComputeFFT(frame []float64) []float64`
Computes the FFT of a real-valued frame.

- Parameters:

  - frame: Input signal frame

- Returns:
  - []float64: Magnitude spectrum (first half only)

## fingerprint package

`Fingerprint(samples []int16, sampleRate int) ([]uint32, error)`
Generates fingerprint hashes from audio samples.

- Parameters:

  - samples: Audio samples
  - sampleRate: Sample rate in Hz

- Returns:
  - []uint32: Fingerprint hashes
  - error: Error if any

`DetectPeaks(spectrogram [][]float64, numBands int) []Peak`
Finds the strongest frequency peaks in each band of the spectrogram.

- Parameters:

  - spectrogram: 2D array of spectral magnitudes
  - numBands: Number of frequency bands

- Returns:
  - []Peak: Array of peak information

`HashFingerprint(peaks []Peak, targetZone int) []uint32`
Creates 32-bit hashes from pairs of audio peaks.

- Parameters:

  - peaks: Array of peak information
  - targetZone: Maximum frame difference for pairing peaks

- Returns:
  - []uint32: Fingerprint hashes

## Development

## Constants

The fingerprinting algorithm uses several constants that can be tuned:

- TargetSampleRate: Downsampled rate (11025 Hz)
- FilterTaps: Number of taps in FIR filter (101, should be odd)
- FrameSize: Samples per frame (1024)
- HopSize: Hop size for overlapping frames (512)
- NumBands: Number of frequency bands for peak detection (6)
- TargetZoneFrames: Maximum frame difference for pairing peaks (20)

## Testing

Each package includes a comprehensive test suite. Run the tests with:

```go
go test ./...
```

The tests cover various scenarios including:

- DSP functionality (filtering, FFT)
- Peak detection and hash generation
- WAV file reading with different formats
- Edge cases and error handling
