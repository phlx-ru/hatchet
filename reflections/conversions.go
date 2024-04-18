package reflections

import (
	"fmt"
	"reflect"

	"github.com/AlekSi/pointer"
)

func ConvertValue(toType reflect.Type, value any) (target *reflect.Value, err error) {
	target = pointer.To[reflect.Value](reflect.Zero(toType))
	if value == nil {
		return target, err
	}
	isTargetPointer := toType.Kind() == reflect.Pointer
	if isTargetPointer {
		target = pointer.To[reflect.Value](reflect.New(toType.Elem()).Elem())
	}
	source := reflect.ValueOf(value)
	if IsEmpty(&source) {
		return target, err
	}
	sourceType := source.Type()
	isSourcePointer := source.Kind() == reflect.Pointer
	if isSourcePointer {
		source = source.Elem()
	}
	targetType := toType
	if isTargetPointer {
		targetType = toType.Elem()
	}
	if source.CanConvert(targetType) {
		targetValue := source.Convert(targetType)
		if targetType.String() == `string` {
			targetValue = reflect.ValueOf(fmt.Sprintf("%v", source.Interface()))
		}
		if isTargetPointer {
			target.Set(targetValue)
			target = pointer.To[reflect.Value](target.Addr())
		} else {
			target = pointer.To[reflect.Value](targetValue)
		}
		return target, err
	}

	return nil, fmt.Errorf(`unable to convert type [%s] to type [%s] with value [%v]`,
		sourceType.String(), toType.String(), reflect.Indirect(source).Interface())
}
