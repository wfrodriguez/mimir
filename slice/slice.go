package slice

type FnSliceCond func(i int) bool

func Index(limit int, cond FnSliceCond) int {
	for i := 0; i < limit; i++ {
		if cond(i) {
			return i
		}
	}

	return -1
}

func InSlice[T comparable](slice []T, val T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}
