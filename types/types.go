package types

import (
	"reflect"
	"time"
)

func String(s string) *string {
	return &s
}

func Bool(b bool) *bool {
	return &b
}

func Time(t time.Time) *time.Time {
	return &t
}

func TimeDuration(duration time.Duration) *time.Duration {
	return &duration
}

func Int(i int) *int {
	return &i
}

func UInt(i uint) *uint {
	return &i
}

func Int8(i int8) *int8 {
	return &i
}

func UInt8(i uint8) *uint8 {
	return &i
}

func Int16(i int16) *int16 {
	return &i
}

func UInt16(i uint16) *uint16 {
	return &i
}

func Int32(i int32) *int32 {
	return &i
}

func UInt32(i uint32) *uint32 {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}

func UInt64(i uint64) *uint64 {
	return &i
}

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

func GoInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func GoUInt(i *uint) uint {
	if i == nil {
		return 0
	}
	return *i
}

func GoInt8(i *int8) int8 {
	if i == nil {
		return 0
	}
	return *i
}

func GoUInt8(i *uint8) uint8 {
	if i == nil {
		return 0
	}
	return *i
}

func GoInt16(i *int16) int16 {
	if i == nil {
		return 0
	}
	return *i
}

func GoUInt16(i *uint16) uint16 {
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

func IsNil(v interface{}) bool {
	return v == nil || reflect.ValueOf(v).IsNil()
}
