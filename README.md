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

## Getting started


### Prerequisites
You need to declare and interface for your tagged type. That interface must implement `Variant() string`.
Then create structs implementing `Variant`, returning their tag value.

```go
type MyFace interface {
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
```

### Instantiating

Use the `pjson.New` function.

```go
var pj := pjson.New([]SomeInterface{SomeStruct1{}, SomeStruct2{}})
```

### Variants
You should have every possible variant in the slice passed to `New`. If you forget a variant, and `pjson` comes across one when unmarshalling, you'll get an error: `failed to find variant field 'foo' in json object`

The default json field that `pjson` looks at is `type`, but you can override that with `WithVariantField` option:
```go
var pj := pjson.New([]SomeInterface{SomeStruct1{}, SomeStruct2{}}, pjson.WithVariantField("variant")) // will look for "variant" in json.
```

Technically, you don't have to provide the variants if just using `MarshalX` methods, but you must for `UnmarshalX` methods.

### Struct vs Pointer

You'll get back the same as what you pass to `New`:
```go
pjson.New([]MyFace{&A{}, &B{}}) // Unmarshal will give you pointers
pjson.New([]MyFace{A{}, B{}}) // Unmarshal will give you structs
```

### Unmarshal

You can then unmarshal either an object or a slice of object.

To Unmarshal an Object:
```go
bytes := []byte(`{"type": "b"}`)
obj, err := pj.UnmarshalObject(bytes)
// obj is a MyFace, of type 'B'
```

To Unmarshal an Array:
```go
bytes := []byte(`[{"type": "b"}, {"type": "a"}]`)
slice, err := pj.UnmarshalArray(bytes)
// slice is a []MyFace{B{},A{}}
```

### Marshal

You can then marshal either an object or a slice of object.

To Marshal an Object:
```go
bytes, err := pj.MarshalObject(B{})
// bytes has '{"type": "b" ...}"
```

To Marshal an Array (slice really, it's called Array to conform to UnmarshalArray):
```go
bytes, err := pj.MarshalArray([]MyFace{B{},A{}})
// bytes has '[{"type": "b" ...}{"type":"a", ...}]"
```

The only thing special that the `Marshal` methods do is to add the variant field so you don't have to have that as a struct field.

## Expanded example
[Run it on go playground](https://go.dev/play/p/jHqZ-TnXq-e)

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
