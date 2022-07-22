package pjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	DefaultVariantField = "type"
)

type Variant interface {
	Variant() string
}

type Options struct {
	VariantField string
}

type Pjson[T Variant] struct {
	Options
	Variants []T
}

type OptionFn func(p *Options)

func New[T Variant](variants []T, options ...OptionFn) Pjson[T] {
	pj := Pjson[T]{
		Options: Options{
			VariantField: DefaultVariantField,
		},
		Variants: variants,
	}
	for _, optFn := range options {
		optFn(&pj.Options)
	}
	return pj
}

func WithVariantField(name string) OptionFn {
	return func(p *Options) {
		p.VariantField = name
	}
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
	variantRes := jRes.Get(c.VariantField)
	if !variantRes.Exists() {
		return c.Variants[0], fmt.Errorf("failed to find variant field '%s' in json object", c.VariantField)
	}
	variantValue := strings.TrimSpace(variantRes.String())
	if variantValue == "" {
		return c.Variants[0], fmt.Errorf("variant field '%s' was empty", c.VariantField)
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

func (c Pjson[T]) MarshalArray(items []T) (bytes []byte, err error) {
	var singleObjBytes []string
	for i, item := range items {
		b, err := c.MarshalObject(item)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal items[%d]: %w", i, err)
		}
		singleObjBytes = append(singleObjBytes, string(b))
	}
	return []byte("[" + strings.Join(singleObjBytes, ",") + "]"), nil
}

func (c Pjson[T]) MarshalObject(item T) (bytes []byte, err error) {

	variant := item.Variant()

	b, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	return sjson.SetBytes(b, c.VariantField, variant)
}
