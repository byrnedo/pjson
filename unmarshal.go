package pjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	DefaultTagField = "type"
)

type Variant interface{ Variant() string }

type Pjson[T Variant] struct {
	TagField string
	Variants []T
}

func New[T Variant](variants []T) Pjson[T] {
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
	variantRes := jRes.Get(c.tagField())
	if !variantRes.Exists() {
		return c.Variants[0], fmt.Errorf("failed to find variant field `%s` in json object", c.TagField)
	}
	variantValue := strings.TrimSpace(variantRes.String())
	if variantValue == "" {
		return c.Variants[0], fmt.Errorf("variant field `%s` was empty", c.TagField)
	}

	for _, obj := range c.Variants {
		if obj.Variant() == variantValue {
			// TODO is there a way around using reflect?
			objT := reflect.TypeOf(obj)
			pv := reflect.New(objT)
			if err := json.Unmarshal([]byte(jRes.Raw), pv.Interface()); err != nil {
				return c.Variants[0], fmt.Errorf("failed to unmarshal variant '%s': %w", variantValue, err)
			}
			return pv.Elem().Interface().(T), nil
		}
	}

	return c.Variants[0], fmt.Errorf("no variant matched type '%s'", variantValue)
}

func (c Pjson[T]) UnmarshalObject(bytes []byte) (T, error) {
	gj := gjson.ParseBytes(bytes)
	return c.unmarshalObjectGjson(gj)
}
