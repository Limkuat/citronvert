package citronvert

import (
	"encoding/binary"
	"os"
	"testing"
)

// Helpers

func readPCMFromWAV(f *os.File, durationMsecs int) []float64 {
	f.Seek(44, 0) // WAV header
	size := samplingRate * durationMsecs / 1000
	samples := make([]int16, size)
	err := binary.Read(f, binary.LittleEndian, samples)
	if err != nil {
		panic(err)
	}
	return F64(samples)
}

func samplesFromFile(filename string) []float64 {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	return readPCMFromWAV(f, 512)
}

// Test

func TestDominantFreq(t *testing.T) {
	f, _ := os.Open("./test_data/sine_440Hz.wav")
	defer f.Close()
	samples := readPCMFromWAV(f, 512)
	spectrum := Spectrum(samples)
	domFreq := DominantFreq(spectrum)
	if !(435 <= domFreq && domFreq <= 445) {
		t.Errorf("Expected 440 Hz (±5 Hz tolerance), got %d Hz", domFreq)
	}
}

func TestSpectralFlatness(t *testing.T) {
	samplesSet := [][]float64{
		samplesFromFile("./test_data/white_noise.wav"),            // very flat (only noise)
		samplesFromFile("./test_data/low_volume_white_noise.wav"), // same, attenuated
		samplesFromFile("./test_data/white_plus_440.wav"),         // less flat (noise + tone, same volume)
		samplesFromFile("./test_data/less_white_plus_440.wav"),    // even less flat (white + tone, attenuated noise)
		samplesFromFile("./test_data/sine_440Hz.wav"),             // not flat (pure tone, no noise)

		samplesFromFile("./test_data/voice.wav"),            // not flat (pure tone, no noise)
		samplesFromFile("./test_data/white_plus_voice.wav"), // not flat (pure tone, no noise)
	}

	for _, samples := range samplesSet {
		//spectrum := Spectrum(samples)
		//sf := SpectralFlatness(spectrum)
		//df := DominantFreq(spectrum)
		//t.Log("DF =", df, "Hz SF =", sf)
		VADScore(samples)
	}
}
