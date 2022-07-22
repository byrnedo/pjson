# pjson

Help to easily JSON marshal / unmarshal tagged unions in go

tldr: allows you to do this

```go
package foo

import (
	"github.com/byrnedo/pjson"
	"encoding/json"
)

type MyFaces []MyInterface

func (f MyFaces) MarshalJSON() ([]byte, error) {
	return pjson.New(MyFaces{}).MarshalArray(f)
}

func (f *MyFaces) UnmarshalJSON(bytes []byte) (err error) {
	*f, err = pjson.New(MyFaces{VariantA{}, VariantB{}}).UnmarshalArray(bytes)
	return
}
```

Expanded ([run it on goplay](https://go.dev/play/p/jHqZ-TnXq-e)):

```go
package main

import (
	"github.com/byrnedo/pjson"
	"encoding/json"
	"fmt"
)

type MyFace interface {
	// Variant func is required
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

func directUsage() {

	bytes := []byte(`[{"type": "a", "a": "AAA"},{"type": "b", "b": "BBB"}]`)
	items, err := pjson.New([]MyFace{A{}, B{}}).UnmarshalArray(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(items)

	bytes = []byte(`{"tag": "a", "a": "AAA"}`)
	item, err := pjson.New([]MyFace{A{}, B{}}, pjson.WithVariantField("tag")).UnmarshalObject(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(item)
}

type MyFaces []MyFace

func (f MyFaces) MarshalJSON() ([]byte, error) {
	return pjson.New(MyFaces{}).MarshalArray(f)
}

func (f *MyFaces) UnmarshalJSON(bytes []byte) (err error) {
	*f, err = pjson.New(MyFaces{A{}, B{}}).UnmarshalArray(bytes)
	return
}

func asField() {
	type customStruct struct {
		Field1 string  `json:"field_1"`
		Field2 int     `json:"field_2"`
		Slice  MyFaces `json:"slice"`
	}

	c := customStruct{
		Field1: "field1",
		Field2: 1,
		Slice:  MyFaces{A{A: "A1"}, B{B: "B1"}, A{A: "A2"}},
	}
	b, _ := json.Marshal(c)
	fmt.Println(string(b))

	c2 := customStruct{}

	json.Unmarshal(b, &c2)
	fmt.Println(c2)
}

func main() {
	directUsage()
	asField()
}
```
