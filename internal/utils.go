package internal

func Distinct[T comparable](tags []T) []T {
	unique := make(map[T]struct{}, 2)

	for _, tag := range tags {
		unique[tag] = struct{}{}
	}

	tags = tags[:0]

	for tag := range unique {
		tags = append(tags, tag)
	}

	return tags
}
