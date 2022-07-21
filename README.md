# pjson

Help to unmarshal tagged unions in go

```go
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

func main() {

    bytes := []byte(`[{"type": "a", "a": "AAA"},{"type": "b", "b": "BBB"}]`)
    items, err := New([]MyFace{A{}, B{}}).Unmarshal(bytes)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Println(items)

}
```
