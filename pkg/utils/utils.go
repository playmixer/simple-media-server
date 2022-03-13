package utils

import "strings"

func CheckExtensions(filename string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(filename, "."+ext) {
			return true
		}
	}
	return false
}
