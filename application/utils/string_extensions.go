package utils

func Contains(items []string, item string) bool {
	for _, p := range items {
		if p == item {
			return true
		}
	}
	return false
}
