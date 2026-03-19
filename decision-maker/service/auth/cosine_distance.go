package auth

import (
	"errors"
	"math"
)

func dotProduct(a []float64, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, errors.New("error: unable to find dot product of vectors with different dimensions")
	}
	var dp float64 = 0
	for i := range len(a) {
		dp += a[i] * b[i]
	}

	return dp, nil
}

func vectorMagnitude(a []float64) float64 {
	var sum float64 = 0
	for _, n := range a {
		sum += math.Pow(n, 2)
	}

	return math.Sqrt(sum)
}

func cosineDistance(u, v []float64) (float64, error) {
	dp, err := dotProduct(u, v)
	if err != nil {
		return 0, err
	}

	magU := vectorMagnitude(u)
	magV := vectorMagnitude(v)

	if magU == 0 || magV == 0 {
		return 0, errors.New("error: cosine distance is undefined for zero vectors")
	}

	return 1 - dp/(magU*magV), nil
}
