package main

import (
	"fmt"
)

type ErrNegativeSqrt float64
func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("\"cannot Sqrt negative number: %.f\"", float64(e));
}

func Sqrt(x float64) (float64, error) {
	if (x < 0){
		return 0, ErrNegativeSqrt(x)
	}
	var z float64 = 1
	for Abs(z * z - x) > 1e-6 {
		z -= (z * z - x) / (2 * z)
	}
	return z, nil
}

func Abs(x float64) float64 {
	if x > 0 {
		return x
	}
	return -x
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
