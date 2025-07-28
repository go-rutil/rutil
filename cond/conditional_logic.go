package cond

// If implements a generic ternary / conditional operator, similar to what is
// found in many other languages (condition ? trueValue : falseValue).
// T can be any type - the values are returned based on the boolean condition.
// Returns ifTrue when true, otherwise returns ifFalse.
// Note that both arguments are evaluated when passed, so there is no short-circuit here.
func If[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}
