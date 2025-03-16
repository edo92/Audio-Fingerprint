package fingerprint

// HashFingerprint creates 32-bit hashes from pairs of audio peaks.
// Each hash combines:
// - 9 bits: anchor frequency
// - 9 bits: target frequency
// - 14 bits: time delta between peaks
func HashFingerprint(peaks []Peak, targetZone int) []uint32 {
	var hashes []uint32
	for i, anchor := range peaks {
		for j := i + 1; j < len(peaks); j++ {
			target := peaks[j]
			dt := target.FrameIndex - anchor.FrameIndex
			if dt < 0 {
				continue
			}
			if dt > targetZone {
				break
			}
			f1 := uint32(anchor.FreqBin)
			f2 := uint32(target.FreqBin)
			dtU := uint32(dt)

			if f1 > 0x1FF {
				f1 = 0x1FF
			}
			if f2 > 0x1FF {
				f2 = 0x1FF
			}
			if dtU > 0x3FFF {
				dtU = 0x3FFF
			}
			hash := (f1 << 23) | (f2 << 14) | dtU
			hashes = append(hashes, hash)
		}
	}
	return hashes
}
