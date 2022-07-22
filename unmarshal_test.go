package pjson

import (
	"reflect"
	"testing"
)

type ABFace interface {
	Variant() string
}

type A struct {
	A string `json:"a"`
}

func (a A) Variant() string {
	return "a"
}

type B struct {
	B string `json:"b"`
}

func (b B) Variant() string {
	return "b"
}

func TestArray(t *testing.T) {
	bytes := []byte(`[{"type": "a", "a": "AAA"},{"type": "b", "b": "BBB"}]`)
	items, err := New([]ABFace{A{}, B{}}).UnmarshalArray(bytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(items)
	t.Log(reflect.TypeOf(items[0]))
	t.Log(reflect.TypeOf(items[1]))
}

func TestObjectHappy(t *testing.T) {
	bytes := []byte(`{"type": "a", "a": "AAA"}`)
	item, err := New([]ABFace{A{}, B{}}).UnmarshalObject(bytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(item)
	t.Log(reflect.TypeOf(item))
}
func TestObjectNoTagMatch(t *testing.T) {
	bytes := []byte(`{"type": "x"}`)
	_, err := New([]ABFace{A{}, B{}}).UnmarshalObject(bytes)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}

func TestArrayNotArray(t *testing.T) {
	bytes := []byte(`{"type": "a"}`)
	_, err := New([]ABFace{A{}, B{}}).UnmarshalArray(bytes)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}

func TestObjectNotObject(t *testing.T) {
	bytes := []byte(`[{"type": "a"}]`)
	_, err := New([]ABFace{A{}, B{}}).UnmarshalObject(bytes)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}
