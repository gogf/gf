package goai

// Callback is specified by OpenAPI/Swagger standard version 3.0.
type Callback map[string]*Path

type Callbacks map[string]*CallbackRef

type CallbackRef struct {
	Ref   string
	Value *Callback
}
