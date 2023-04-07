package pjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Variant interface {
	Variant() string
}

type Discriminator interface {
	Field() string
	Variants() []Variant
}
type Tagged[T Discriminator] struct {
	d     T
	Value Variant
}

func (o Tagged[T]) MarshalJSON() ([]byte, error) {
	if o.Value == nil {
		return json.Marshal(o.Value)
	}
	variant := o.Value.Variant()

	b, err := json.Marshal(o.Value)
	if err != nil {
		return nil, err
	}

	return sjson.SetBytes(b, o.d.Field(), variant)
}
func (o *Tagged[T]) UnmarshalJSON(bytes []byte) error {

	if len(bytes) == 0 || string(bytes) == "null" {
		return nil
	}

	jRes := gjson.ParseBytes(bytes)

	variants := o.d.Variants()
	if !jRes.IsObject() {
		return fmt.Errorf("did not hold an Object")
	}

	variantRes := jRes.Get(o.d.Field())
	if !variantRes.Exists() {
		return fmt.Errorf("failed to find variant field '%s' in json object", o.d.Field())
	}
	variantValue := strings.TrimSpace(variantRes.String())
	if variantValue == "" {
		return fmt.Errorf("variant field '%s' was empty", o.d.Field())
	}

	for _, obj := range variants {
		if obj.Variant() != variantValue {
			continue
		}

		t := reflect.TypeOf(obj)
		dest := obj
		// a pointer works just fine, but if it's not we need to get one
		if t.Kind() != reflect.Pointer {
			dest = reflect.New(t).Interface().(Variant)
		}

		if err := json.Unmarshal([]byte(jRes.Raw), &dest); err != nil {
			return fmt.Errorf("failed to unmarshal variant '%s': %w", variantValue, err)
		}

		o.Value = dest
		return nil
	}

	return fmt.Errorf("no variant matched type '%s'", variantValue)
}
