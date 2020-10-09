package utils

func Uniq(values []string) []string {
	keys := make(map[string]bool)
	result := make([]string, 0)
	for _, entry := range values {
		if _, exists := keys[entry]; !exists {
			keys[entry] = true
			result = append(result, entry)
		}
	}
	return result
}

func Reverse(values []map[string]string) []map[string]string {
	result := make([]map[string]string, len(values))
	for i, v := range values {
		result[len(values)-i-1] = v
	}
	return result
}
