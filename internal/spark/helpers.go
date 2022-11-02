package spark

func appendIfNotNil[T any](array []*T, items ...*T) []*T {
	for _, item := range items {
		if item != nil {
			array = append(array, item)
		}
	}
	return array
}
