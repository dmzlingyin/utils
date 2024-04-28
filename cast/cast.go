package cast

import (
	"github.com/spf13/cast"
	"time"
)

func ToBool(v any, d ...bool) bool {
	val, err := cast.ToBoolE(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToTime(v any, d ...time.Time) time.Time {
	val, err := cast.ToTimeE(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToTimeInDefaultLocation(v any, location *time.Location) time.Time {
	val, _ := cast.ToTimeInDefaultLocationE(v, location)
	return val
}

func ToDuration(v any, d ...time.Duration) time.Duration {
	val, err := cast.ToDurationE(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToFloat64(v any, d ...float64) float64 {
	val, err := cast.ToFloat64E(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToFloat32(v any, d ...float32) float32 {
	val, err := cast.ToFloat32E(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToInt64(v any, d ...int64) int64 {
	val, err := cast.ToInt64E(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToInt32(v any, d ...int32) int32 {
	val, err := cast.ToInt32E(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToInt16(v any, d ...int16) int16 {
	val, err := cast.ToInt16E(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToInt8(v any, d ...int8) int8 {
	val, err := cast.ToInt8E(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToInt(v any, d ...int) int {
	val, err := cast.ToIntE(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToUint(v any, d ...uint) uint {
	val, err := cast.ToUintE(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToUint64(v any, d ...uint64) uint64 {
	val, err := cast.ToUint64E(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToUint32(v any, d ...uint32) uint32 {
	val, err := cast.ToUint32E(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToUint16(v any, d ...uint16) uint16 {
	val, err := cast.ToUint16E(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToUint8(v any, d ...uint8) uint8 {
	val, err := cast.ToUint8E(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}

func ToString(v any, d ...string) string {
	val, err := cast.ToStringE(v)
	if err != nil && len(d) > 0 {
		return d[0]
	}
	return val
}
