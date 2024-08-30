package utils

import "fmt"

// Given a map, returns the value of a key as a string
func GetMapStringValue(metadata map[string]interface{}, key string) string {
	if value, ok := metadata[key].(string); ok {
		return value
	} else if value, ok := metadata[key].(bool); ok {
		return fmt.Sprintf("%v", value)
	} else if value, ok := metadata[key].(int); ok {
		return fmt.Sprintf("%v", value)
	}
	return ""
}
