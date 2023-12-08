package util

import (
	"strconv"
	"strings"
)

func ParseInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func ParseUInt(s string) uint {
	v := ParseInt(s)
	return uint(v)
}

func ParseInt64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

func ParseUInt64(s string) uint64 {
	v := ParseInt64(s)
	return uint64(v)
}

func Anonymous(v string, max int, rep rune) string {
	if strings.TrimSpace(v) == "" {
		return ""
	}
	if max < 2 {
		max = 2
	}
	r := make([]rune, max)
	for i := range r {
		r[i] = rep
	}
	source := []rune(v)
	r[0] = source[0]
	if len(source) > 2 {
		r[len(r)-1] = source[len(source)-1]
	}
	return string(r)
}
