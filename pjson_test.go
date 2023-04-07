package pjson_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/byrnedo/pjson"
)

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

func TestArray(t *testing.T) {
	bytes := []byte(`[{"type": "a", "a": "AAA"},{"type": "a", "a": "AAA1"},{"type": "b", "b": "BBB"}]`)
	items, err := pjson.New([]pjson.Variant{A{}, B{}}, pjson.WithVariantField("type")).UnmarshalArray(bytes)
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
	item, err := pjson.New([]pjson.Variant{A{}, B{}}).UnmarshalObject(bytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(item)
	t.Log(reflect.TypeOf(item))
}
func TestObjectNoTagMatch(t *testing.T) {
	bytes := []byte(`{"type": "x"}`)
	_, err := pjson.New([]pjson.Variant{A{}, B{}}).UnmarshalObject(bytes)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}

func TestArrayNotArray(t *testing.T) {
	bytes := []byte(`{"type": "a"}`)
	_, err := pjson.New([]pjson.Variant{A{}, B{}}).UnmarshalArray(bytes)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}

func TestObjectNotObject(t *testing.T) {
	bytes := []byte(`[{"type": "a"}]`)
	_, err := pjson.New([]pjson.Variant{A{}, B{}}).UnmarshalObject(bytes)
	if err == nil {
		t.Fatal("should have error")
	}

	t.Log(err)
}

func TestPjson_MarshalObject(t *testing.T) {
	b, err := pjson.New([]pjson.Variant{}).MarshalObject(A{A: "AAA"})
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(b) {
		t.Fatal("json invalid")
	}
	t.Log(string(b))
}

func TestPjson_MarshalArray(t *testing.T) {
	b, err := pjson.New([]pjson.Variant{}).MarshalArray([]pjson.Variant{A{A: "AA0"}, A{A: "AA1"}, A{A: "AA2"}, B{B: "BB0"}})
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(b) {
		t.Fatal("json invalid")
	}
	t.Log(string(b))
}

type ABFaces []pjson.Variant

func (f ABFaces) MarshalJSON() ([]byte, error) {
	return pjson.New(ABFaces{}).MarshalArray(f)
}

func (f *ABFaces) UnmarshalJSON(bytes []byte) (err error) {
	*f, err = pjson.New(ABFaces{A{}, B{}}).UnmarshalArray(bytes)
	return
}

type SuperObject struct {
	FieldA string  `json:"field_a"`
	FieldB int     `json:"field_b"`
	Slice  ABFaces `json:"slice"`
}

func TestSuperObject(t *testing.T) {

	s := SuperObject{
		FieldA: "A",
		FieldB: 1,
		Slice:  []pjson.Variant{A{}, A{}, B{}},
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

	if reflect.TypeOf(s2.Slice[0]) != reflect.TypeOf(A{}) {
		t.Fatal("wrong type")
	}

	if reflect.TypeOf(s2.Slice[2]) != reflect.TypeOf(B{}) {
		t.Fatal("wrong type")
	}

}

type Foo struct {
	pjson.Variant
}

func (f Foo) Variants() []pjson.Variant {
	return []pjson.Variant{
		A{}, B{},
	}
}

var fooPjson = pjson.New(Foo{}.Variants())

func (f *Foo) UnmarshalJSON(bytes []byte) error {
	v, err := fooPjson.UnmarshalObject(bytes)
	if err != nil {
		return err
	}
	f.Variant = v
	return nil
}

func BenchmarkMarshal(b *testing.B) {

	bytes := []byte(`{"type": "a", "a": "AAAA", "a_one": 1, "a_two": "two"}`)
	b.Run("with pjson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			f := Foo{}
			err := json.Unmarshal(bytes, &f)
			if err != nil {
				b.Fatal(err)
			}

		}
	})

	b.Run("without pjson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {

			a := A{}
			err := json.Unmarshal(bytes, &a)
			if err != nil {
				b.Fatal(err)
			}

		}
	})
}
