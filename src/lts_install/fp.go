package lts_install

func Find[T any](collection []T, predicate func(item T) bool) (T, bool) {
	for i := range collection {
		if predicate(collection[i]) {
			return collection[i], true
		}
	}

	var result T
	return result, false
}
