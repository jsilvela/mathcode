package main

import (
	"errors"
	"fmt"
	"math"
)

type realfunc func(float64) float64

func rateOfChange(f realfunc, x, h float64) (float64, error) {
	if h == 0 {
		return 0, errors.New("h cannot be zero")
	}
	return (f(x+h) - f(x)) / h, nil
}

func derivativeAt(f realfunc, x float64) (float64, error) {
	var limit float64 = 0
	for h := 0.1; h > 0.00000000001; h = h / 2 {
		rate, _ := rateOfChange(f, x, h)
		if math.Abs(limit-rate) < 0.000001 {
			return rate, nil
		}
		limit = rate
	}
	return 0, fmt.Errorf("function not differentiable at %f", x)
}

func derivative(f realfunc) realfunc {
	return func(x float64) float64 {
		r, _ := derivativeAt(f, x)
		return r
	}
}

type vector struct {
	X float64
	Y float64
	Z float64
}

type vectorFunc func(vector) float64

type linearMap func(vector) float64

func sumVect(a, b vector) vector {
	return vector{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

func scaleVect(k float64, a vector) vector {
	return vector{
		X: k * a.X,
		Y: k * a.Y,
		Z: k * a.Z,
	}
}

func lengthVect(a vector) float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

func linearProjectionAt(f vectorFunc, x vector) linearMap {
	return func(h vector) float64 {
		var limit float64 = 0
		hLength := lengthVect(h)
		for k := 0.1; k*hLength > 0.00000000001; k = k / 2 {
			projection := (f(sumVect(x, scaleVect(k, h))) - f(x)) / k
			if math.Abs(limit-projection) < 0.0000001 {
				return limit
			}
			limit = projection
		}
		return 0
	}
}

type formField func(x vector) linearMap

func derivative3D(f vectorFunc) formField {
	return func(x vector) linearMap {
		return linearProjectionAt(f, x)
	}
}

func main() {
	// non-differentiable function
	beast := func(x float64) float64 {
		if x == 0 {
			return 0
		}
		return x * math.Sin(1/x)
	}
	_, err := derivativeAt(beast, 0)
	fmt.Printf("derivativeAt(beast, 0) → %s\n", err.Error())
	fmt.Printf("beast(0) = %f\n", beast(0))

	// real functions
	fmt.Printf("exp(1) → %f\n", math.Exp(1))
	expPrime := derivative(math.Exp)
	fmt.Printf("expPrime(1) → %f\n", expPrime(1))
	fmt.Printf("cos(π) → %f\n", math.Cos(math.Pi))
	sinPrime := derivative(math.Sin)
	fmt.Printf("sinPrime(π) → %f\n", sinPrime(math.Pi))

	// vector functions
	sphereF := func(a vector) float64 {
		return a.X*a.X + a.Y*a.Y + a.Z*a.Z
	}
	spherePrime := derivative3D(sphereF)
	fmt.Printf("spherePrime(1, 0, 0)(1, 0, 0) → %f\n",
		spherePrime(vector{X: 1})(vector{X: 1}))
	fmt.Printf("spherePrime(1, 0, 0)(0, 1, 0) → %f\n",
		spherePrime(vector{X: 1})(vector{Y: 1}))
	fmt.Printf("spherePrime(1, 0, 0)(0, 0, 1) → %f\n",
		spherePrime(vector{X: 1})(vector{Z: 1}))
}
