package citronvert

import (
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

const (
	frameDurationMsecs = 256
	samplingRate       = 16000
)

type Power struct {
	F    int
	XHat complex128
}

type VADResult struct {
	SF     float64
	E      float64
	DF     int
	Voiced bool
}

func NormalizedF64(samples []int16) []float64 {
	normalizedSamples := make([]float64, len(samples))
	var maxPeak int16
	for _, level := range samples {
		if level > maxPeak {
			maxPeak = level
		}
	}
	if maxPeak == 0 {
		maxPeak = 1
	}
	var factor float64 = 1 / float64(maxPeak)
	for i, level := range samples {
		normalizedSamples[i] = float64(level) * factor
	}
	return normalizedSamples
}

func F64(samples []int16) []float64 {
	const maxInt16AsF64 float64 = 32768.0
	fsamples := make([]float64, len(samples))
	for i, level := range samples {
		fsamples[i] = float64(level) / maxInt16AsF64
	}
	return fsamples
}

func arithmeticMean(spectrum []Power) float64 {
	acc := 0.0
	for _, s := range spectrum {
		mod := cmplx.Abs(s.XHat)
		acc += mod * mod
	}
	return acc / float64(len(spectrum))
}

func geometricMean(spectrum []Power) float64 {
	acc := 0.0
	for _, s := range spectrum {
		mod := cmplx.Abs(s.XHat)
		acc += math.Log(mod * mod)
	}
	return math.Exp(acc / float64(len(spectrum)))
}

func SpectralFlatness(spectrum []Power) float64 {
	if len(spectrum) == 0 {
		return 0
	}
	Am := arithmeticMean(spectrum)
	Gm := geometricMean(spectrum)
	if Am == 0 {
		return math.Inf(1)
	}
	return Gm / Am
}

func DominantFreq(spectrum []Power) (dominant int) {
	lastHighest := 0.0
	for _, p := range spectrum {
		mod := cmplx.Abs(p.XHat)
		if mod > lastHighest {
			lastHighest = mod
			dominant = p.F
		}
	}
	return
}

func Spectrum(samples []float64) []Power {
	fftResult := fft.FFTReal(samples)
	N := len(samples)
	spectrum := make([]Power, N)
	for k, z := range fftResult {
		spectrum[k].F = k * samplingRate / N
		spectrum[k].XHat = z
	}
	return spectrum
}

func Energy(spectrum []Power) float64 {
	acc := 0.0
	for _, s := range spectrum {
		mod := cmplx.Abs(s.XHat)
		acc += mod * mod
	}
	return acc / float64(len(spectrum))
}

func VADScore(samples []float64) VADResult {
	spectrum := Spectrum(samples)
	SF := SpectralFlatness(spectrum)
	DF := DominantFreq(spectrum)
	E := Energy(spectrum)
	return VADResult{
		DF:     DF,
		SF:     SF,
		E:      E,
		Voiced: false,
	}
}
