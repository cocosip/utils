// Package util provides a collection of utility functions.
package util

import (
	"strconv"
	"strings"
)

// ParseIntOrDefault parses a string to an int. It returns 0 if the parsing fails.
func ParseIntOrDefault(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

// ParseIntE parses a string to an int, returning an error if parsing fails.
func ParseIntE(s string) (int, error) {
	return strconv.Atoi(s)
}

// ParseUintOrDefault parses a string to a uint. It returns 0 if the parsing fails.
// This function correctly handles parsing of unsigned integers.
func ParseUintOrDefault(s string) uint {
	v, err := strconv.ParseUint(s, 10, 0) // bitSize 0 means it infers from the system's int size
	if err != nil {
		return 0
	}
	return uint(v)
}

// ParseUintE parses a string to a uint, returning an error if parsing fails.
func ParseUintE(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

// ParseInt64OrDefault parses a string to an int64. It returns 0 if the parsing fails.
func ParseInt64OrDefault(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// ParseInt64E parses a string to an int64, returning an error if parsing fails.
func ParseInt64E(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseUint64OrDefault parses a string to a uint64. It returns 0 if the parsing fails.
// This function correctly handles parsing of unsigned integers.
func ParseUint64OrDefault(s string) uint64 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// ParseUint64E parses a string to a uint64, returning an error if parsing fails.
func ParseUint64E(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// AnonymizeString masks a string, showing only the first and last characters.
// It is useful for displaying sensitive information like phone numbers or names.
// v: The input string to anonymize.
// max: The maximum length of the output string.
// rep: The rune used for masking the middle part.
func AnonymizeString(v string, max int, rep rune) string {
	if strings.TrimSpace(v) == "" {
		return ""
	}

	source := []rune(v)
	sourceLen := len(source)

	if max < 2 {
		max = 2 // Ensure the output has at least two characters for first/last.
	}

	result := make([]rune, max)
	for i := range result {
		result[i] = rep
	}

	if sourceLen > 0 {
		result[0] = source[0]
	}

	// If the source string has 2 or more characters, and the result has space for a last character.
	if sourceLen >= 2 && max > 1 {
		result[len(result)-1] = source[sourceLen-1]
	}

	return string(result)
}