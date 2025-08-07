package typetools

import (
	"reflect"
	"time"
)

// String returns a pointer to the given string.
// Parameter: s string - the original string
// Return: *string - pointer to the string
func String(s string) *string {
	return &s
}

// Bool returns a pointer to the given bool.
// Parameter: b bool - the original boolean value
// Return: *bool - pointer to the boolean value
func Bool(b bool) *bool {
	return &b
}

// Time returns a pointer to the given time.Time.
// Parameter: t time.Time - the original time
// Return: *time.Time - pointer to the time
func Time(t time.Time) *time.Time {
	return &t
}

// TimeDuration returns a pointer to the given time.Duration.
// Parameter: duration time.Duration - the original duration
// Return: *time.Duration - pointer to the duration
func TimeDuration(duration time.Duration) *time.Duration {
	return &duration
}

// Int returns a pointer to the given int.
// Parameter: i int - the original integer
// Return: *int - pointer to the integer
func Int(i int) *int {
	return &i
}

// UInt returns a pointer to the given uint.
// Parameter: i uint - the original unsigned integer
// Return: *uint - pointer to the unsigned integer
func UInt(i uint) *uint {
	return &i
}

// Int8 returns a pointer to the given int8.
// Parameter: i int8 - the original int8 value
// Return: *int8 - pointer to the int8 value
func Int8(i int8) *int8 { return &i }

// UInt8 returns a pointer to the given uint8.
// Parameter: i uint8 - the original uint8 value
// Return: *uint8 - pointer to the uint8 value
func UInt8(i uint8) *uint8 { return &i }

// Int16 returns a pointer to the given int16.
// Parameter: i int16 - the original int16 value
// Return: *int16 - pointer to the int16 value
func Int16(i int16) *int16 { return &i }

// UInt16 returns a pointer to the given uint16.
// Parameter: i uint16 - the original uint16 value
// Return: *uint16 - pointer to the uint16 value
func UInt16(i uint16) *uint16 { return &i }

// Int32 returns a pointer to the given int32.
// Parameter: i int32 - the original int32 value
// Return: *int32 - pointer to the int32 value
func Int32(i int32) *int32 { return &i }

// UInt32 returns a pointer to the given uint32.
// Parameter: i uint32 - the original uint32 value
// Return: *uint32 - pointer to the uint32 value
func UInt32(i uint32) *uint32 { return &i }

// Int64 returns a pointer to the given int64.
// Parameter: i int64 - the original int64 value
// Return: *int64 - pointer to the int64 value
func Int64(i int64) *int64 { return &i }

// UInt64 returns a pointer to the given uint64.
// Parameter: i uint64 - the original uint64 value
// Return: *uint64 - pointer to the uint64 value
func UInt64(i uint64) *uint64 { return &i }

// GoString returns the value of the string pointer or "" if nil.
// Parameter: s *string - pointer to string
// Return: string - value or empty string if nil
func GoString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// GoBool returns the value of the bool pointer or false if nil.
// Parameter: b *bool - pointer to bool
// Return: bool - value or false if nil
func GoBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// GoInt returns the value of the int pointer or 0 if nil.
// Parameter: i *int - pointer to int
// Return: int - value or 0 if nil
func GoInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// GoUInt returns the value of the uint pointer or 0 if nil.
// Parameter: i *uint - pointer to uint
// Return: uint - value or 0 if nil
func GoUInt(i *uint) uint {
	if i == nil {
		return 0
	}
	return *i
}

// GoInt8 returns the value of the int8 pointer or 0 if nil.
// Parameter: i *int8 - pointer to int8
// Return: int8 - value or 0 if nil
func GoInt8(i *int8) int8 {
	if i == nil {
		return 0
	}
	return *i
}

// GoUInt8 returns the value of the uint8 pointer or 0 if nil.
// Parameter: i *uint8 - pointer to uint8
// Return: uint8 - value or 0 if nil
func GoUInt8(i *uint8) uint8 {
	if i == nil {
		return 0
	}
	return *i
}

// GoInt16 returns the value of the int16 pointer or 0 if nil.
// Parameter: i *int16 - pointer to int16
// Return: int16 - value or 0 if nil
func GoInt16(i *int16) int16 {
	if i == nil {
		return 0
	}
	return *i
}

// GoUInt16 returns the value of the uint16 pointer or 0 if nil.
// Parameter: i *uint16 - pointer to uint16
// Return: uint16 - value or 0 if nil
func GoUInt16(i *uint16) uint16 {
	if i == nil {
		return 0
	}
	return *i
}

// GoInt32 returns the value of the int32 pointer or 0 if nil.
// Parameter: i *int32 - pointer to int32
// Return: int32 - value or 0 if nil
func GoInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

// GoUInt32 returns the value of the uint32 pointer or 0 if nil.
// Parameter: i *uint32 - pointer to uint32
// Return: uint32 - value or 0 if nil
func GoUInt32(i *uint32) uint32 {
	if i == nil {
		return 0
	}
	return *i
}

// GoInt64 returns the value of the int64 pointer or 0 if nil.
// Parameter: i *int64 - pointer to int64
// Return: int64 - value or 0 if nil
func GoInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// GoUInt64 returns the value of the uint64 pointer or 0 if nil.
// Parameter: i *uint64 - pointer to uint64
// Return: uint64 - value or 0 if nil
func GoUInt64(i *uint64) uint64 {
	if i == nil {
		return 0
	}
	return *i
}

// GoTime returns the value of the time.Time pointer or zero value if nil.
// Parameter: t *time.Time - pointer to time.Time
// Return: time.Time - value or zero if nil
func GoTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// GoTimeDuration returns the value of the time.Duration pointer or 0 if nil.
// Parameter: duration *time.Duration - pointer to time.Duration
// Return: time.Duration - value or 0 if nil
func GoTimeDuration(duration *time.Duration) time.Duration {
	if duration == nil {
		return 0
	}
	return *duration
}

// IsNil checks if the given interface is nil or points to a nil value.
// Parameter: v interface{} - any value
// Return: bool - true if nil or points to nil, false otherwise
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	}
	return false
}
