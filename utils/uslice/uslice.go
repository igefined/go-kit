package uslice

func Contains[T comparable](arr []T, elem T) bool {
	for _, el := range arr {
		if el == elem {
			return true
		}
	}

	return false
}
