package reflections

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"
)

const (
	FieldStringSimple  = `StringSimple`
	FieldStringPointer = `StringPointer`
	FieldIntSimple     = `IntSimple`
	FieldIntPointer    = `IntPointer`
	FieldBoolSimple    = `BoolSimple`
	FieldBoolPointer   = `BoolPointer`
	FieldTimeSimple    = `TimeSimple`
	FieldTimePointer   = `TimePointer`
)

type testStruct struct {
	StringSimple  string     `json:"string_simple,omitempty"`
	StringPointer *string    `json:"string_pointer,omitempty"`
	IntSimple     int        `json:"int_simple,omitempty"`
	IntPointer    *int       `json:"int_pointer,omitempty"`
	BoolSimple    bool       `json:"bool_simple,omitempty"`
	BoolPointer   *bool      `json:"bool_pointer,omitempty"`
	TimeSimple    time.Time  `json:"time_simple,omitempty"`
	TimePointer   *time.Time `json:"time_pointer,omitempty"`
}

func TestConvertValue(t *testing.T) {
	testCases := []struct {
		name           string
		targetType     reflect.Type
		value          any
		expectedString string
		expectedType   string
		expectedError  string
	}{
		{
			name:           "nil_to_string",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldStringSimple).Type(),
			value:          nil,
			expectedType:   "string",
			expectedString: "",
		},
		{
			name:           "nil_to_*string",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldStringPointer).Type(),
			value:          nil,
			expectedType:   "*string",
			expectedString: "<nil>",
		},
		{
			name:           "string_to_string",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldStringSimple).Type(),
			value:          `hello!`,
			expectedType:   "string",
			expectedString: `hello!`,
		},
		{
			name:           "*string_to_string",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldStringSimple).Type(),
			value:          pointer.ToString(`hello!`),
			expectedType:   "string",
			expectedString: `hello!`,
		},
		{
			name:           "string_to_*string",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldStringPointer).Type(),
			value:          `hello!`,
			expectedType:   "*string",
			expectedString: `hello!`,
		},
		{
			name:           "*string_to_*string",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldStringPointer).Type(),
			value:          pointer.ToString(`hello!`),
			expectedType:   "*string",
			expectedString: `hello!`,
		},
		{
			name:          "string_to_int_error",
			targetType:    reflect.ValueOf(testStruct{}).FieldByName(FieldIntSimple).Type(),
			value:         `hello!`,
			expectedError: "unable to convert type [string] to type [int] with value [hello!]",
		},
		{
			name:          "string_to_*int_error",
			targetType:    reflect.ValueOf(testStruct{}).FieldByName(FieldIntPointer).Type(),
			value:         `hello!`,
			expectedError: "unable to convert type [string] to type [*int] with value [hello!]",
		},
		{
			name:          "string_to_int_error",
			targetType:    reflect.ValueOf(testStruct{}).FieldByName(FieldIntSimple).Type(),
			value:         pointer.ToString(`hello!`),
			expectedError: "unable to convert type [*string] to type [int] with value [hello!]",
		},
		{
			name:          "string_to_*int_error",
			targetType:    reflect.ValueOf(testStruct{}).FieldByName(FieldIntPointer).Type(),
			value:         pointer.ToString(`hello!`),
			expectedError: "unable to convert type [*string] to type [*int] with value [hello!]",
		},
		{
			name:           "int_to_string",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldStringSimple).Type(),
			value:          1089746,
			expectedType:   "string",
			expectedString: "1089746",
		},
		{
			name:           "*int_to_string",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldStringSimple).Type(),
			value:          pointer.ToInt(1089746),
			expectedType:   "string",
			expectedString: "1089746",
		},
		{
			name:           "int_to_*string",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldStringPointer).Type(),
			value:          1089746,
			expectedType:   "*string",
			expectedString: "1089746",
		},
		{
			name:           "*int_to_*string",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldStringPointer).Type(),
			value:          pointer.ToInt(1089746),
			expectedType:   "*string",
			expectedString: "1089746",
		},
		{
			name:           "int_to_int",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldIntSimple).Type(),
			value:          1089746,
			expectedType:   "int",
			expectedString: "1089746",
		},
		{
			name:           "int_to_*int",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldIntPointer).Type(),
			value:          1089746,
			expectedType:   "*int",
			expectedString: "1089746",
		},
		{
			name:           "bool_to_bool_with_true",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldBoolSimple).Type(),
			value:          true,
			expectedType:   "bool",
			expectedString: "true",
		},
		{
			name:           "bool_to_bool_with_false",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldBoolSimple).Type(),
			value:          false,
			expectedType:   "bool",
			expectedString: "false",
		},
		{
			name:           "bool_to_*bool_with_true",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldBoolPointer).Type(),
			value:          true,
			expectedType:   "*bool",
			expectedString: "true",
		},
		{
			name:           "bool_to_*bool_with_false",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldBoolPointer).Type(),
			value:          false,
			expectedType:   "*bool",
			expectedString: "false",
		},
		{
			name:           "*bool_to_bool_with_true",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldBoolSimple).Type(),
			value:          pointer.ToBool(true),
			expectedType:   "bool",
			expectedString: "true",
		},
		{
			name:           "bool_to_bool_with_false",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldBoolSimple).Type(),
			value:          pointer.ToBool(false),
			expectedType:   "bool",
			expectedString: "false",
		},
		{
			name:           "bool_to_*bool_with_true",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldBoolPointer).Type(),
			value:          pointer.ToBool(true),
			expectedType:   "*bool",
			expectedString: "true",
		},
		{
			name:           "bool_to_*bool_with_false",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldBoolPointer).Type(),
			value:          pointer.ToBool(false),
			expectedType:   "*bool",
			expectedString: "false",
		},
		{
			name:           "nil_to_*time_with",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldTimePointer).Type(),
			value:          nil,
			expectedType:   "*time.Time",
			expectedString: "<nil>",
		},
		{
			name:           "nil_to_time_with",
			targetType:     reflect.ValueOf(testStruct{}).FieldByName(FieldTimeSimple).Type(),
			value:          nil,
			expectedType:   "time.Time",
			expectedString: "0001-01-01 00:00:00 +0000 UTC",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualResult, actualError := ConvertValue(testCase.targetType, testCase.value)
			if testCase.expectedError != "" {
				require.Error(t, actualError)
				require.Equal(t, testCase.expectedError, actualError.Error())
			} else {
				require.NoError(t, actualError)
				require.NotNil(t, actualResult)
				indirect := reflect.Indirect(*actualResult)
				var actualValue any
				if !IsEmpty(&indirect) && indirect.CanInterface() {
					actualValue = indirect.Interface()
				}
				actualString := fmt.Sprintf(`%v`, actualValue)
				require.Equal(t, testCase.expectedString, actualString)
				actualType := (*actualResult).Type().String()
				require.Equal(t, testCase.expectedType, actualType)
			}
		})
	}
}
