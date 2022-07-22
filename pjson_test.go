package pjson_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/byrnedo/pjson"
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
	bytes := []byte(`[{"type": "a", "a": "AAA"},{"type": "a", "a": "AAA1"},{"type": "b", "b": "BBB"}]`)
	items, err := pjson.New([]ABFace{A{}, B{}}, pjson.WithVariantField("type")).UnmarshalArray(bytes)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 3 {
		t.Fatalf("got %d, wanted 3", len(items))
	}

	t.Log(items)

	if reflect.TypeOf(items[0]) != reflect.TypeOf(A{}) {
		t.Fatal("wrong type")
	}
	if reflect.TypeOf(items[1]) != reflect.TypeOf(A{}) {
		t.Fatal("wrong type")
	}
	if reflect.TypeOf(items[2]) != reflect.TypeOf(B{}) {
		t.Fatal("wrong type")
	}
}

func TestObjectHappy(t *testing.T) {
	bytes := []byte(`{"type": "b", "b": "BBB"}`)
	item, err := pjson.New([]ABFace{A{}, B{}}).UnmarshalObject(bytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(item)
	t.Log(reflect.TypeOf(item))
}
func TestObjectNoTagMatch(t *testing.T) {
	bytes := []byte(`{"type": "x"}`)
	_, err := pjson.New([]ABFace{A{}, B{}}).UnmarshalObject(bytes)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}

func TestArrayNotArray(t *testing.T) {
	bytes := []byte(`{"type": "a"}`)
	_, err := pjson.New([]ABFace{A{}, B{}}).UnmarshalArray(bytes)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}

func TestObjectNotObject(t *testing.T) {
	bytes := []byte(`[{"type": "a"}]`)
	_, err := pjson.New([]ABFace{A{}, B{}}).UnmarshalObject(bytes)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}

func TestPjson_MarshalObject(t *testing.T) {
	b, err := pjson.New([]ABFace{}).MarshalObject(A{A: "AAA"})
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(b) {
		t.Fatal("json invalid")
	}
	t.Log(string(b))
}

func TestPjson_MarshalArray(t *testing.T) {
	b, err := pjson.New([]ABFace{}).MarshalArray([]ABFace{A{A: "AA0"}, A{A: "AA1"}, A{A: "AA2"}, B{B: "BB0"}})
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(b) {
		t.Fatal("json invalid")
	}
	t.Log(string(b))
}