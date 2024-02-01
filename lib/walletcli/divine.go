package walletcli

import "math/rand"

func GetDivine(min, max float64, count int) []float64 {
	hasil := make([]float64, count)

	sumcount := float64(0)

	for c := 0; c < count; c++ {
		ratio := (min + rand.Float64()*(max-min)) + float64(1)
		hasil[c] = ratio

		sumcount += ratio
	}

	fix := make([]float64, count)
	for c := 0; c < count; c++ {

		fix[c] = hasil[c] / sumcount
	}

	return fix
}
