package util

import (
	"reflect"
	"time"
)

func GoString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func GoBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func GoTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func GoTimeDuration(duration *time.Duration) time.Duration {
	if duration == nil {
		return 0
	}
	return *duration
}

func GoInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func GoInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func GoUInt32(i *uint32) uint32 {
	if i == nil {
		return 0
	}
	return *i
}

func GoInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func GoUInt64(i *uint64) uint64 {
	if i == nil {
		return 0
	}
	return *i
}

func IsNil(v interface{}) bool {
	return v == nil || reflect.ValueOf(v).IsNil()
}

func GoPtrString(s string) *string {
	return &s
}

func GoPtrBool(b bool) *bool {
	return &b
}

func GoPtrTime(t time.Time) *time.Time {
	return &t
}

func GoPtrTimeDuration(duration time.Duration) *time.Duration {
	return &duration
}

func GoPtrInt(i int) *int {
	return &i
}

func GoPtrInt32(i int32) *int32 {
	return &i
}

func GoPtrUInt32(i uint32) *uint32 {
	return &i
}

func GoPtrInt64(i int64) *int64 {
	return &i
}

func GoPtrUInt64(i uint64) *uint64 {
	return &i
}
