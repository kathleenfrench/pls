package utils

// GetKeysFromMapStringInterface parses keys from a map[string]interface
func GetKeysFromMapStringInterface(m map[string]interface{}) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

// GetKeysFromMapString returns the values of a map[string]string keys as a slice of strings
func GetKeysFromMapString(m map[string]string) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
