# pjson

Help to unmarshal tagged unions in go

```go
package main

import (
	"github.com/byrnedo/pjson"
	"fmt"
)

type MyFace interface {
	// Type func is required
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

func main() {

	bytes := []byte(`[{"type": "a", "a": "AAA"},{"type": "b", "b": "BBB"}]`)
	items, err := pjson.New([]MyFace{A{}, B{}}).UnmarshalArray(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(items)

	bytes = []byte(`{"type": "a", "a": "AAA"}`)
	item, err := pjson.New([]MyFace{A{}, B{}}).UnmarshalObject(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(item)
}
```
