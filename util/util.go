package util

func TernaryIf[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}

	return falseVal
}
