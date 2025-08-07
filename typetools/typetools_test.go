package typetools

import (
	"testing"
	"time"
)

func TestString(t *testing.T) {
	v := "hello"
	ptr := String(v)
	if ptr == nil || *ptr != v {
		t.Errorf("String failed: got %v", ptr)
	}
	if GoString(ptr) != v {
		t.Errorf("GoString failed: got %v", GoString(ptr))
	}
	if GoString(nil) != "" {
		t.Errorf("GoString(nil) should be empty string")
	}
}

func TestBool(t *testing.T) {
	v := true
	ptr := Bool(v)
	if ptr == nil || *ptr != v {
		t.Errorf("Bool failed: got %v", ptr)
	}
	if !GoBool(ptr) {
		t.Errorf("GoBool failed: got %v", GoBool(ptr))
	}
	if GoBool(nil) {
		t.Errorf("GoBool(nil) should be false")
	}
}

func TestIntTypes(t *testing.T) {
	if GoInt(Int(42)) != 42 { t.Error("GoInt failed") }
	if GoInt(nil) != 0 { t.Error("GoInt(nil) failed") }
	if GoUInt(UInt(42)) != 42 { t.Error("GoUInt failed") }
	if GoUInt(nil) != 0 { t.Error("GoUInt(nil) failed") }
	if GoInt8(Int8(8)) != 8 { t.Error("GoInt8 failed") }
	if GoInt8(nil) != 0 { t.Error("GoInt8(nil) failed") }
	if GoUInt8(UInt8(8)) != 8 { t.Error("GoUInt8 failed") }
	if GoUInt8(nil) != 0 { t.Error("GoUInt8(nil) failed") }
	if GoInt16(Int16(16)) != 16 { t.Error("GoInt16 failed") }
	if GoInt16(nil) != 0 { t.Error("GoInt16(nil) failed") }
	if GoUInt16(UInt16(16)) != 16 { t.Error("GoUInt16 failed") }
	if GoUInt16(nil) != 0 { t.Error("GoUInt16(nil) failed") }
	if GoInt32(Int32(32)) != 32 { t.Error("GoInt32 failed") }
	if GoInt32(nil) != 0 { t.Error("GoInt32(nil) failed") }
	if GoUInt32(UInt32(32)) != 32 { t.Error("GoUInt32 failed") }
	if GoUInt32(nil) != 0 { t.Error("GoUInt32(nil) failed") }
	if GoInt64(Int64(64)) != 64 { t.Error("GoInt64 failed") }
	if GoInt64(nil) != 0 { t.Error("GoInt64(nil) failed") }
	if GoUInt64(UInt64(64)) != 64 { t.Error("GoUInt64 failed") }
	if GoUInt64(nil) != 0 { t.Error("GoUInt64(nil) failed") }
}

func TestTimeTypes(t *testing.T) {
	now := time.Now()
	dur := time.Second * 5
	if GoTime(Time(now)) != now { t.Error("GoTime failed") }
	if !GoTime(Time(now)).Equal(now) { t.Error("GoTime Equal failed") }
	if !GoTime(nil).IsZero() { t.Error("GoTime(nil) should be zero") }
	if GoTimeDuration(TimeDuration(dur)) != dur { t.Error("GoTimeDuration failed") }
	if GoTimeDuration(nil) != 0 { t.Error("GoTimeDuration(nil) failed") }
}

func TestIsNil(t *testing.T) {
	var p *int = nil
	if !IsNil(p) { t.Error("IsNil failed for nil pointer") }
	var m map[string]int = nil
	if !IsNil(m) { t.Error("IsNil failed for nil map") }
	var s []int = nil
	if !IsNil(s) { t.Error("IsNil failed for nil slice") }
	if IsNil(123) { t.Error("IsNil failed for non-nil value") }
}