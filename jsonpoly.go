package pjson

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tidwall/gjson"
)

type Typer interface{ Type() string }

func Unmarshal[T Typer](bytes []byte, variant1 T, variants ...T) (items []T, err error) {

	discriminator := "type"
	gj := gjson.ParseBytes(bytes)
	if !gj.IsArray() {
		return nil, fmt.Errorf("bytes did not hold an array")
	}
	results := gj.Array()
	for i, jRes := range results {
		if !jRes.IsObject() {
			jType := jRes.Type.String()
			if jType == gjson.JSON.String() {
				if jRes.IsArray() {
					jType = "Array"
				}
			}
			return nil, fmt.Errorf("[%d] did not hold an Object, was %s", i, jType)
		}
		tag := jRes.Get(discriminator).String()
		if tag == "" {
			return nil, fmt.Errorf("failed to find tag field `%s` in json object", discriminator)
		}

		item, err := magic[T](tag, []byte(jRes.Raw), variant1, variants...)
		if err != nil {
			return nil, fmt.Errorf("[%d]: %w", i, err)
		}
		items = append(items, item)

	}
	return
}

func magic[T Typer](tag string, bytes []byte, variant1 T, variants ...T) (T, error) {

	for _, obj := range append([]T{variant1}, variants...) {
		if obj.Type() == tag {
			objT := reflect.TypeOf(obj)
			pv := reflect.New(objT)
			if err := json.Unmarshal(bytes, pv.Interface()); err != nil {
				return variant1, fmt.Errorf("failed to unmarshal variant '%s': %w", tag, err)
			}
			return pv.Elem().Interface().(T), nil
		}
	}

	return variant1, fmt.Errorf("no variant matched type '%s'", tag)
}
