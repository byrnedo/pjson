package pjson

import (
	"reflect"
	"testing"
)

type MyFace interface {
	Type() string
}

type A struct {
	A string `json:"a"`
}

func (a A) Type() string {
	return "a"
}

type B struct {
	B string `json:"b"`
}

func (b B) Type() string {
	return "b"
}

func TestSomething(t *testing.T) {

	bytes := []byte(`[{"type": "a", "a": "AAA"},{"type": "b", "b": "BBB"}]`)
	items, err := Unmarshal[MyFace](bytes, A{}, B{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(items)
	t.Log(reflect.TypeOf(items[0]))
	t.Log(reflect.TypeOf(items[1]))
}
