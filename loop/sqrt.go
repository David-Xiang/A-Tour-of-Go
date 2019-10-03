package main

import (
	"fmt"
)

func Sqrt(x float64) (z float64) {
	z = 1
	for Abs(z * z - x) > 1e-6 {
		z -= (z * z - x) / (2 * z)
	}
	return
}

func Abs(x float64) float64 {
	if x > 0 {
		return x
	}
	return -x
}

func main() {
	fmt.Println(Sqrt(2))
}