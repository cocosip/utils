package util

import "strconv"

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
