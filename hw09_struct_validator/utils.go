package hw09structvalidator

type ContainsConstains interface {
	~int | ~string
}

func Contains[T ContainsConstains](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
