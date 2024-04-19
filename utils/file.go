package utils

import (
	"strings"
)

func GetExtensions(filename string) string {
	split := strings.Split(filename, ".")

	if len(split) <= 1 {
		return ""
	}

	return split[len(split)-1]
}
