package cast

import "testing"

func TestCast(t *testing.T) {
	if !ToBool(1) {
		t.Error("ToBool failed")
	}
	if ToBool(0) {
		t.Error("ToBool failed")
	}
	if ToInt(1.0) != 1 {
		t.Error("ToInt failed")
	}
	if ToInt("1") != 1 {
		t.Error("ToInt failed")
	}
	if ToInt("abc", 1) != 1 {
		t.Error("ToInt failed")
	}
	if ToInt64("1.2345", 1) != 1 {
		t.Error("ToInt64 failed")
	}
	if ToInt32(1.2345) != 1 {
		t.Error("ToInt32 failed")
	}
	if ToInt8("2.34", 2) != 2 {
		t.Error("ToInt8 failed")
	}
	if ToFloat32("1.3") != 1.3 {
		t.Error("ToFloat32 failed")
	}
	if ToFloat64("1.2345") != 1.2345 {
		t.Error("ToFloat64 failed")
	}
	if ToString(1.0) != "1" {
		t.Error("ToString failed")
	}
	if ToString("") != "" {
		t.Error("ToString failed")
	}
	if ToString(nil) != "" {
		t.Error("ToString failed")
	}
}
