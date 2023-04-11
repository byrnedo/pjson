# pjson

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/byrnedo/pjson/blob/master/LICENSE.txt)
[![Go Reference](https://pkg.go.dev/badge/github.com/byrnedo/pjson.svg)](https://pkg.go.dev/github.com/byrnedo/pjson)
[![Go Coverage](https://github.com/byrnedo/pjson/wiki/coverage.svg)](https://raw.githack.com/wiki/byrnedo/pjson/coverage.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/byrnedo/pjson)](https://goreportcard.com/report/github.com/byrnedo/pjson)

Help to easily JSON marshal / unmarshal tagged unions in go

A tagged union / discriminating type is, for instance with the following JSON:

```json
[
  {
    "type": "a",
    "a_name": "AName",
    "a_foo": "FOO"
  },
  {
    "type": "b",
    "b_name": "BName",
    "b_goo": "GOO"
  }
]
```

The `type` field denotes which type the object is. So many object share a common discriminating field.
In some languages this is supported, but not in go.

**Pjson** gives us a helper `pjson.Tagged` type to create these pseudo tagged unions that can be automatically
serialized and deserialized to and from JSON.

## Usage

```go
package readme

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/byrnedo/pjson"
)

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
	// now that we have our types we can use Tagged
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
```

## Benchmarks

Macbook Pro M1 2022

```
Benchmark/unmarshal_with_pjson
Benchmark/unmarshal_with_pjson-10         	  867177	      1356 ns/op
Benchmark/unmarshal_without_pjson
Benchmark/unmarshal_without_pjson-10      	 1793629	       670.6 ns/op
Benchmark/marshal_with_pjson
Benchmark/marshal_with_pjson-10           	 2415705	       488.7 ns/op
Benchmark/marshal_without_pjson
Benchmark/marshal_without_pjson-10        	10956252	       109.8 ns/op
```
