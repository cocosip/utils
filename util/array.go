package util

import "strings"

func Distinct(values []string) []string {
	var result []string
	for i := range values {
		if strings.TrimSpace(values[i]) != "" {
			for j := range result {
				if values[i] != result[j] {
					result = append(result, values[i])
				}
			}
		}
	}
	return result
}
