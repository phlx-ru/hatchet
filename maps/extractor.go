package maps

import (
	"encoding/json"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Extractor struct {
	value gjson.Result
}

func MakeExtractor(in any) (*Extractor, error) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	value := gjson.ParseBytes(bytes)
	return &Extractor{value: value}, nil
}

func MustMakeExtractor(in any) *Extractor {
	extractor, err := MakeExtractor(in)
	if err != nil {
		panic(err)
	}
	return extractor
}

func (e *Extractor) Get(path string) *Extractor {
	clone := *e
	clone.value = clone.value.Get(path)
	return &clone
}

// GetByCargoType get string from <path>_uid instead of json from path if <path>_type is `jxs:string`
func (e *Extractor) GetByCargoType(path string) *Extractor {
	clone := *e
	clone.value = clone.value.Get(path)
	if !clone.value.Exists() {
		if e.value.Get(path+`_type`).String() == `jxs:string` {
			clone.value = e.value.Get(path + `_uid`)
		}
	}
	return &clone
}

func (e *Extractor) Exists() bool {
	if e.value.Type == gjson.Null {
		return false
	}
	return e.value.Exists()
}

func (e *Extractor) ToString() string {
	return e.value.String()
}

func (e *Extractor) ToPointerString() *string {
	if !e.Exists() {
		return nil
	}
	return pointer.ToString(e.value.String())
}

func (e *Extractor) ToPointerStringOrNil() *string {
	if !e.Exists() {
		return nil
	}
	return pointer.ToStringOrNil(e.value.String())
}

func (e *Extractor) ToBool() bool {
	return e.value.Bool()
}

func (e *Extractor) ToPointerBool() *bool {
	if !e.Exists() {
		return nil
	}
	return pointer.ToBool(e.value.Bool())
}

func (e *Extractor) ToInt() int {
	return int(e.value.Int())
}

func (e *Extractor) ToPointerInt() *int {
	if !e.Exists() {
		return nil
	}
	return pointer.ToInt(int(e.value.Int()))
}

func (e *Extractor) ToPointerIntOrNil() *int {
	if !e.Exists() {
		return nil
	}
	return pointer.ToIntOrNil(int(e.value.Int()))
}

func (e *Extractor) ToInt32() int32 {
	return int32(e.value.Int())
}

func (e *Extractor) ToPointerInt32() *int32 {
	if !e.Exists() {
		return nil
	}
	return pointer.ToInt32(int32(e.value.Int()))
}

func (e *Extractor) ToPointerInt32OrNil() *int32 {
	if !e.Exists() {
		return nil
	}
	return pointer.ToInt32OrNil(int32(e.value.Int()))
}

func (e *Extractor) ToInt64() int64 {
	return e.value.Int()
}

func (e *Extractor) ToPointerInt64() *int64 {
	if !e.Exists() {
		return nil
	}
	return pointer.ToInt64(e.value.Int())
}

func (e *Extractor) ToPointerInt64OrNil() *int64 {
	if !e.Exists() {
		return nil
	}
	return pointer.ToInt64OrNil(e.value.Int())
}

func (e *Extractor) ToUint64() uint64 {
	return uint64(e.value.Int())
}

func (e *Extractor) ToPointerUint64() *uint64 {
	if !e.Exists() {
		return nil
	}
	return pointer.ToUint64(e.ToUint64())
}

func (e *Extractor) ToPointerUint64OrNil() *uint64 {
	if !e.Exists() {
		return nil
	}
	return pointer.ToUint64OrNil(e.ToUint64())
}

func (e *Extractor) ToFloat32() float32 {
	return float32(e.value.Float())
}

func (e *Extractor) ToPointerFloat32() *float32 {
	if !e.Exists() {
		return nil
	}
	return pointer.ToFloat32(float32(e.value.Float()))
}

func (e *Extractor) ToPointerFloat32OrNil() *float32 {
	if !e.Exists() {
		return nil
	}
	return pointer.ToFloat32OrNil(float32(e.value.Float()))
}

func (e *Extractor) ToFloat64() float64 {
	return e.value.Float()
}

func (e *Extractor) ToPointerFloat64() *float64 {
	if !e.Exists() {
		return nil
	}
	return pointer.ToFloat64(e.value.Float())
}

func (e *Extractor) ToPointerTimestamp() *timestamppb.Timestamp {
	if !e.Exists() {
		return nil
	}
	return timestamppb.New(e.value.Time())
}

func (e *Extractor) ToTimeString() string {
	return e.value.Time().Format(time.RFC3339)
}

func (e *Extractor) ToPointerTimeString() *string {
	if !e.Exists() {
		return nil
	}
	return pointer.ToString(e.ToTimeString())
}

func (e *Extractor) IsArray() bool {
	return e.value.IsArray()
}

func (e *Extractor) IsJSON() bool {
	return e.value.Type == gjson.JSON
}

func (e *Extractor) IsString() bool {
	return e.value.Type == gjson.String
}

func (e *Extractor) ToArray() []*Extractor {
	array := []*Extractor{}
	for _, value := range e.value.Array() {
		extractor := &Extractor{value: value}
		array = append(array, extractor)
	}
	return array
}
