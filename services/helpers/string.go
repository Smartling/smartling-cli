package helpers

import (
	"fmt"
	"strings"
)

// MKeyValueToMap converts key/value pairs `key="value"` into map[string]string{key:value}
func MKeyValueToMap(ss []string) (map[string]string, error) {
	res := make(map[string]string)
	for _, s := range ss {
		pair, err := KeyValueToMap(s)
		if err != nil {
			return nil, err
		}
		for k, v := range pair {
			if _, exists := res[k]; exists {
				return nil, fmt.Errorf("duplicate key: %s", k)
			}
			res[k] = v
		}
	}
	return res, nil
}

// KeyValueToMap converts key/value pair `key="value"` into map[string]string{key:value}
func KeyValueToMap(s string) (map[string]string, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid key=value format: %s", s)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
	return map[string]string{key: value}, nil
}
