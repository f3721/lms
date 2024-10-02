package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

// SubFloat64 return a - b
func SubFloat64(a, b float64) float64 {
	return decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b)).InexactFloat64()
}

// AddFloat64 return a + b
func AddFloat64(a, b float64) float64 {
	return decimal.NewFromFloat(a).Add(decimal.NewFromFloat(b)).InexactFloat64()
}

// MulFloat64AndInt return a * b
func MulFloat64AndInt(a float64, b int) float64 {
	return decimal.NewFromFloat(a).Mul(decimal.NewFromFloat(float64(b))).InexactFloat64()
}

// DivFloat64 return a / b
func DivFloat64(a, b float64) float64 {
	return decimal.NewFromFloat(a).Div(decimal.NewFromFloat(b)).InexactFloat64()
}

// DivFloat64 return a / b
func DivFloat64Fixed(a, b float64, fixed bool) float64 {
	f64 := decimal.NewFromFloat(a).Div(decimal.NewFromFloat(b)).InexactFloat64()
	if fixed {
		number, err := strconv.ParseFloat(fmt.Sprintf("%.2f", f64), 64)
		if err == nil {
			f64 = number
		}
	}
	return f64
}

// MulFloat64AndInt return a * b

func MulFloat64AndIntFixed(a float64, b int, fixed bool) float64 {
	f64 := decimal.NewFromFloat(a).Mul(decimal.NewFromFloat(float64(b))).InexactFloat64()
	if fixed {
		number, err := strconv.ParseFloat(fmt.Sprintf("%.2f", f64), 64)
		if err == nil {
			f64 = number
		}
	}
	return f64
}

func SubFloat64Fixed(a float64, b float64, fixed bool) float64 {
	f64 := decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b)).InexactFloat64()
	if fixed {
		number, err := strconv.ParseFloat(fmt.Sprintf("%.2f", f64), 64)
		if err == nil {
			f64 = number
		}
	}
	return f64
}

func AddFloat64Fixed(a float64, b float64, fixed bool) float64 {
	f64 := decimal.NewFromFloat(a).Add(decimal.NewFromFloat(b)).InexactFloat64()
	if fixed {
		number, err := strconv.ParseFloat(fmt.Sprintf("%.2f", f64), 64)
		if err == nil {
			f64 = number
		}
	}
	return f64
}
