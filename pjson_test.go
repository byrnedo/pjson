package pjson_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/byrnedo/pjson"
)

type ABDisc struct{}

func (ab ABDisc) Field() string {
	return "type"
}
func (ab ABDisc) Variants() []pjson.Variant {
	return []pjson.Variant{A{}, B{}}
}

type A struct {
	A    string `json:"a"`
	AOne int    `json:"a_one,omitempty"`
	ATwo string `json:"a_two,omitempty"`
}

func (a A) Variant() string {
	return "a"
}

type B struct {
	B    string `json:"b"`
	BOne int    `json:"b_one,omitempty"`
	BTwo string `json:"b_two,omitempty"`
}

func (b B) Variant() string {
	return "b"
}

func TestObjectNoTagMatch(t *testing.T) {
	bytes := []byte(`{"type": "x"}`)

	f := pjson.Tagged[ABDisc]{}
	err := json.Unmarshal(bytes, &f)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}

func TestNil(t *testing.T) {
	bytes := []byte(`null`)

	f := pjson.Tagged[ABDisc]{}
	err := json.Unmarshal(bytes, &f)
	if err != nil {
		t.Fatal(err)
	}
	if f.Value != nil {
		t.Fatal("should be nil")
	}

	bytes, err = json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bytes))
}

//	func TestArrayNotArray(t *testing.T) {
//		bytes := []byte(`{"type": "a"}`)
//		_, err := pjson.New([]pjson.Variant{A{}, B{}}).UnmarshalArray(bytes)
//		if err == nil {
//			t.Fatal("should have error")
//		}
//
//		t.Log(err)
//	}
func TestObjectNotObject(t *testing.T) {
	bytes := []byte(`[{"type": "a"}]`)
	f := pjson.Tagged[ABDisc]{}
	err := json.Unmarshal(bytes, &f)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}
func TestPjson_MarshalObject(t *testing.T) {
	f := pjson.Tagged[ABDisc]{Value: A{}}
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(b) {
		t.Fatal("json invalid")
	}

	if string(b) != "{\"a\":\"\",\"type\":\"a\"}" {
		t.Fatal(string(b))
	}
}

type SuperObject struct {
	FieldA string                 `json:"field_a"`
	FieldB int                    `json:"field_b"`
	Slice  []pjson.Tagged[ABDisc] `json:"slice"`
}

func TestSuperObject(t *testing.T) {

	s := SuperObject{
		FieldA: "A",
		FieldB: 1,
		Slice:  []pjson.Tagged[ABDisc]{{Value: A{}}, {Value: A{}}, {Value: B{}}},
	}

	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
	if string(b) != `{"field_a":"A","field_b":1,"slice":[{"a":"","type":"a"},{"a":"","type":"a"},{"b":"","type":"b"}]}` {
		t.Fatal(string(b))
	}

	s2 := SuperObject{}

	err = json.Unmarshal([]byte(`{"field_a":"A","field_b":1,"slice":[{"a":"AA1","type":"a"},{"a":"AA2","type":"a"},{"b":"BB1","type":"b"}]}`), &s2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s2)
	if len(s2.Slice) != 3 {
		t.Fatal("not 3 elems")
	}

	if reflect.TypeOf(s2.Slice[0].Value) != reflect.TypeOf(&A{}) {
		t.Fatal("wrong type")
	}

	if reflect.TypeOf(s2.Slice[2].Value) != reflect.TypeOf(&B{}) {
		t.Fatal("wrong type")
	}

}

func Benchmark(b *testing.B) {

	bytes := []byte(`{"type": "a", "a": "AAAA", "a_one": 1, "a_two": "two"}`)
	b.Run("unmarshal with pjson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			f := pjson.Tagged[ABDisc]{}
			err := json.Unmarshal(bytes, &f)
			if err != nil {
				b.Fatal(err)
			}

		}
	})

	b.Run("unmarshal without pjson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			a := A{}
			err := json.Unmarshal(bytes, &a)
			if err != nil {
				b.Fatal(err)
			}

		}
	})

	b.Run("marshal with pjson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			f := pjson.Tagged[ABDisc]{Value: A{}}
			_, err := json.Marshal(f)
			if err != nil {
				b.Fatal(err)
			}

		}
	})

	b.Run("marshal without pjson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			a := A{}
			_, err := json.Marshal(a)
			if err != nil {
				b.Fatal(err)
			}

		}
	})
}

type Foo struct {
	A string `json:"a"`
}

// set it's tag value
func (a Foo) Variant() string {
	return "foo"
}

type Bar struct {
	B string `json:"b"`
}

func (b Bar) Variant() string {
	return "bar"
}

// specify the union
type FooBarUnion struct{}

func (u FooBarUnion) Field() string { return "type" }

func (u FooBarUnion) Variants() []pjson.Variant {
	return []pjson.Variant{
		Foo{}, Bar{},
	}
}

func ExampleReadme() {
	// now that we have our types we can use OneOf
	o := pjson.Tagged[FooBarUnion]{}

	bytes := []byte(`{"type": "foo", "a": "AAAA"}`)

	err := json.Unmarshal(bytes, &o)
	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.TypeOf(o.Value), o.Value)

	bytes, _ = json.Marshal(o)
	fmt.Println(string(bytes))

	// Output: *pjson_test.Foo &{AAAA}
	// {"a":"AAAA","type":"foo"}

}
