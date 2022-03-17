package gmath

import (
	"math"
	"math/cmplx"
	"strconv"
)

// Decbin Returns a string containing a binary representation of the given num argument.
func Decbin(num int64) string {
	return strconv.FormatInt(num, 2)
}

// Dechex Returns a string containing a hexadecimal representation of the given unsigned num argument.
func Dechex(num int64) string {
	return strconv.FormatInt(num, 16)
}

// Decoct Returns a string containing an octal representation of the given num argument.
func Decoct(num int64) string {
	return strconv.FormatInt(num, 8)
}

// Max If the first and only parameter is an array, max() returns the highest value in that array.
// If at least two parameters are provided, max() returns the biggest of these values.
func Max(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		max = math.Max(max, nums[i])
	}
	return max
}

// Min If the first and only parameter is an array, min() returns the lowest value in that array.
// If at least two parameters are provided, min() returns the smallest of these values.
func Min(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	min := nums[0]
	for i := 1; i < len(nums); i++ {
		min = math.Min(min, nums[i])
	}
	return min
}

// Round Returns the rounded value of num to specified precision (number of digits after the decimal point). precision can also be negative or zero (default).
func Round(num float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Trunc((num+0.5/p)*p) / p
}

// Floor Returns the next lowest integer value (as float) by rounding down num if necessary.
func Floor(num float64) float64 {
	return math.Floor(num)
}

// Ceil Returns the next highest integer value by rounding up num if necessary.
func Ceil(num float64) float64 {
	return math.Ceil(num)
}

// Pi Returns an approximation of pi. Also, you can use the M_PI constant which yields identical results to pi().
func Pi() float64 {
	return math.Pi
}

// Abs Returns the absolute value of num.
func Abs(num float64) float64 {
	return math.Abs(num)
}

// Acos Returns the arc cosine of num in radians. acos() is the inverse function of cos(), which means that a==cos(acos(a)) for every value of a that is within acos()' range.
func Acos(num complex128) complex128 {
	return cmplx.Acos(num)
}

// Acosh Returns the inverse hyperbolic cosine of num, i.e. the value whose hyperbolic cosine is num.
func Acosh(num complex128) complex128 {
	return cmplx.Acosh(num)
}

// Asin Returns the arc sine of num in radians. asin() is the inverse function of sin(), which means that a==sin(asin(a)) for every value of a that is within asin()'s range.
func Asin(num complex128) complex128 {
	return cmplx.Asin(num)
}

// Asinh Returns the inverse hyperbolic sine of num, i.e. the value whose hyperbolic sine is num.
func Asinh(num complex128) complex128 {
	return cmplx.Asinh(num)
}

// Atan2 TThis function calculates the arc tangent of the two variables x and y.
// It is similar to calculating the arc tangent of y / x, except that the signs of both arguments are used to determine the quadrant of the result.
func Atan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

// Atan Returns the arc tangent of num in radians. atan() is the inverse function of tan(), which means that a==tan(atan(a)) for every value of a that is within atan()'s range.
func Atan(num complex128) complex128 {
	return cmplx.Atan(num)
}

// Atanh Returns the inverse hyperbolic tangent of num, i.e. the value whose hyperbolic tangent is num.
func Atanh(num complex128) complex128 {
	return cmplx.Atanh(num)
}

// Cos Returns the cosine of the num parameter. The num parameter is in radians.
func Cos(num float64) float64 {
	return math.Cos(num)
}

// Cosh Returns the hyperbolic cosine of num, defined as (exp(arg) + exp(-arg))/2.
func Cosh(num float64) float64 {
	return math.Cosh(num)
}

// Exp Returns e raised to the power of num.
func Exp(num float64) float64 {
	return math.Exp(num)
}

// Expm1 Returns the equivalent to 'exp(num) - 1' computed in a way that is accurate even if the value of num is near zero, a case where 'exp (num) - 1' would be inaccurate due to subtraction of two numbers that are nearly equal.
func Expm1(num float64) float64 {
	return math.Exp(num) - 1
}

// IsFinite Checks whether num is a legal finite on this platform.
func IsFinite(num float64, sign int) bool {
	return !math.IsInf(num, sign)
}

// IsInfinite Returns true if num is infinite (positive or negative), like the result of log(0) or any value too big to fit into a float on this platform.
func IsInfinite(num float64, sign int) bool {
	return math.IsInf(num, sign)
}

// IsNan Checks whether num is 'not a number', like the result of acos(1.01).
func IsNan(num float64) bool {
	return math.IsNaN(num)
}

// Log Returns the natural logarithm of num.
func Log(num float64) float64 {
	return math.Log(num)
}

// Log10 Returns the base-10 logarithm of num.
func Log10(num float64) float64 {
	return math.Log10(num)
}

// Log1p Returns log(1 + number), computed in a way that is accurate even when the value of number is close to zero.
func Log1p(num float64) float64 {
	return math.Log1p(num)
}

// Pow Returns num raised to the power of exponent.
func Pow(num, exponent float64) float64 {
	return math.Pow(num, exponent)
}

// Sin Returns the sine of the num parameter. The num parameter is in radians.
func Sin(num float64) float64 {
	return math.Sin(num)
}

// Sinh Returns the hyperbolic sine of num, defined as (exp(num) - exp(-num))/2.
func Sinh(num float64) float64 {
	return math.Sinh(num)
}

// Sqrt Returns the square root of num.
func Sqrt(num float64) float64 {
	return math.Sqrt(num)
}

// Tan Returns the tangent of the num parameter. The num parameter is in radians.
func Tan(num float64) float64 {
	return math.Tan(num)
}

// Tanh Returns the hyperbolic tangent of num, defined as sinh(num)/cosh(num).
func Tanh(num float64) float64 {
	return math.Tanh(num)
}

// BaseConvert Returns a string containing num represented in base to_base.
// The base in which num is given is specified in from_base.
// Both from_base and to_base have to be between 2 and 36, inclusive.
// Digits in numbers with a base higher than 10 will be represented with the letters a-z, with a meaning 10, b meaning 11 and z meaning 35.
// The case of the letters doesn't matter, i.e. num is interpreted case-insensitively.
func BaseConvert(x string, fromBase, toBase int) (string, error) {
	i, err := strconv.ParseInt(x, fromBase, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, toBase), nil
}
