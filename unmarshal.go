package pjson

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tidwall/gjson"
)

const (
	DefaultTagField = "type"
)

type Typer interface{ Type() string }

type Pjson[T Typer] struct {
	TagField string
	Variants []T
}

func New[T Typer](variants []T) Pjson[T] {
	return Pjson[T]{
		Variants: variants,
	}
}

func (c Pjson[T]) tagField() string {
	if c.TagField == "" {
		return DefaultTagField
	}
	return c.TagField
}

func (c Pjson[T]) UnmarshalArray(bytes []byte) (items []T, err error) {

	gj := gjson.ParseBytes(bytes)
	if !gj.IsArray() {
		return nil, fmt.Errorf("bytes did not hold an array")
	}
	results := gj.Array()
	for i, jRes := range results {

		item, err := c.unmarshalObjectGjson(jRes)
		if err != nil {
			return nil, fmt.Errorf("[%d]: %w", i, err)
		}
		items = append(items, item)

	}
	return
}

func (c Pjson[T]) unmarshalObjectGjson(jRes gjson.Result) (T, error) {

	if !jRes.IsObject() {
		jType := jRes.Type.String()
		if jType == gjson.JSON.String() {
			if jRes.IsArray() {
				jType = "Array"
			}
		}
		return c.Variants[0], fmt.Errorf("did not hold an Object, was %s", jType)
	}
	tagValue := jRes.Get(c.tagField()).String()
	if tagValue == "" {
		return c.Variants[0], fmt.Errorf("failed to find tag field `%s` in json object", c.TagField)
	}

	for _, obj := range c.Variants {
		if obj.Type() == tagValue {
			objT := reflect.TypeOf(obj)
			pv := reflect.New(objT)
			if err := json.Unmarshal([]byte(jRes.Raw), pv.Interface()); err != nil {
				return c.Variants[0], fmt.Errorf("failed to unmarshal variant '%s': %w", tagValue, err)
			}
			return pv.Elem().Interface().(T), nil
		}
	}

	return c.Variants[0], fmt.Errorf("no variant matched type '%s'", tagValue)
}

func (c Pjson[T]) UnmarshalObject(bytes []byte) (T, error) {
	gj := gjson.ParseBytes(bytes)
	return c.unmarshalObjectGjson(gj)
}
