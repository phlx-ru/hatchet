package reflections

import "reflect"

// IsEmpty check reflected value is not valid or nil, but not checks on Zero value
func IsEmpty(value *reflect.Value) bool {
	if value == nil {
		return true
	}
	if !value.IsValid() {
		return true
	}
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return true
	}
	return false
}
