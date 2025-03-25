package cast

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"strings"
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

// StructToMap 将结构体转换为map，并根据tag设定map的key值
func StructToMap(in any) (map[string]any, error) {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("not supported non struct type")
	}

	result := make(map[string]any)
	err := structToMap(v, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func structToMap(v reflect.Value, result map[string]any) error {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		// 处理匿名字段（组合结构体）
		if fieldType.Anonymous {
			if err := structToMap(fieldValue, result); err != nil {
				return err
			}
			continue
		}

		tag := fieldType.Tag.Get("map")
		if tag == "" || tag == "-" {
			continue
		}
		omitempty := false
		if items := strings.Split(tag, ","); len(items) == 2 {
			tag = items[0]
			omitempty = items[1] == "omitempty"
		}

		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				continue
			}
			fieldValue = fieldValue.Elem()
		}
		if omitempty && fieldValue.IsZero() {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Struct:
			nestedResult := make(map[string]any)
			if err := structToMap(fieldValue, nestedResult); err != nil {
				return err
			}
			if omitempty && len(nestedResult) <= 0 {
				continue
			}
			result[tag] = nestedResult
		case reflect.Map:
			mapResult := make(map[string]any)
			for _, key := range fieldValue.MapKeys() {
				mapValue := fieldValue.MapIndex(key)
				if mapValue.Kind() == reflect.Ptr && mapValue.IsNil() {
					continue
				}
				if mapValue.Kind() == reflect.Ptr {
					mapValue = mapValue.Elem()
				}
				if omitempty && mapValue.IsZero() {
					continue
				}
				if mapValue.Kind() == reflect.Struct {
					nestedMapResult := make(map[string]any)
					if err := structToMap(mapValue, nestedMapResult); err != nil {
						return err
					}
					if omitempty && len(nestedMapResult) <= 0 {
						continue
					}
					mapResult[key.String()] = nestedMapResult
				} else {
					mapResult[key.String()] = mapValue.Interface()
				}
			}
			if omitempty && len(mapResult) <= 0 {
				continue
			}
			result[tag] = mapResult
		default:
			result[tag] = fieldValue.Interface()
		}
	}

	return nil
}

// MapToStruct 将map转换为结构体，并根据tag设定字段值
func MapToStruct(m map[string]any, out any) error {
	v := reflect.ValueOf(out)
	// 确保输出是一个指针
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("output is not a pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get("map")
		if tag == "" {
			tag = fieldType.Name
		}
		if mapValue, ok := m[tag]; ok {
			convertedValue, err := convertValue(mapValue, fieldValue.Type())
			if err == nil && fieldValue.CanSet() {
				fieldValue.Set(convertedValue)
			}
		}
	}

	return nil
}

// convertValue 将map中的值转换为struct字段的类型
func convertValue(value any, targetType reflect.Type) (reflect.Value, error) {
	if value == nil {
		return reflect.Zero(targetType), nil
	}

	val := reflect.ValueOf(value)
	valType := val.Type()
	if valType.AssignableTo(targetType) {
		return val, nil
	}

	// 处理基本类型的转换
	switch targetType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if valType.Kind() >= reflect.Int && valType.Kind() <= reflect.Int64 {
			return reflect.ValueOf(val.Convert(targetType).Interface()), nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if valType.Kind() >= reflect.Uint && valType.Kind() <= reflect.Uint64 {
			return reflect.ValueOf(val.Convert(targetType).Interface()), nil
		}
	case reflect.Float32, reflect.Float64:
		if valType.Kind() == reflect.Float32 || valType.Kind() == reflect.Float64 {
			return reflect.ValueOf(val.Convert(targetType).Interface()), nil
		}
	case reflect.String:
		if valType.Kind() == reflect.String {
			return val, nil
		}
	case reflect.Struct:
		if valType.Kind() == reflect.Map {
			nestedStruct := reflect.New(targetType).Elem()
			err := MapToStruct(value.(map[string]any), nestedStruct.Addr().Interface())
			if err != nil {
				return reflect.Value{}, err
			}
			return nestedStruct, nil
		}
	case reflect.Ptr:
		if valType.Kind() == reflect.Map {
			elemType := targetType.Elem()
			nestedStruct := reflect.New(elemType).Elem()
			err := MapToStruct(value.(map[string]any), nestedStruct.Addr().Interface())
			if err != nil {
				return reflect.Value{}, err
			}
			ptr := reflect.New(elemType)
			ptr.Elem().Set(nestedStruct)
			return ptr, nil
		}
	default:
		return val, errors.New("unsupported type")
	}

	return reflect.Value{}, fmt.Errorf("cannot convert %v to %v", valType, targetType)
}
