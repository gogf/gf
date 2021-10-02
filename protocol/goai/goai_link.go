package goai

// Link is specified by OpenAPI/Swagger standard version 3.0.
type Link struct {
	OperationID  string                 `json:"operationId,omitempty"  yaml:"operationId,omitempty"`
	OperationRef string                 `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	Description  string                 `json:"description,omitempty"  yaml:"description,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"   yaml:"parameters,omitempty"`
	Server       *Server                `json:"server,omitempty"       yaml:"server,omitempty"`
	RequestBody  interface{}            `json:"requestBody,omitempty"  yaml:"requestBody,omitempty"`
}

type Links map[string]LinkRef

type LinkRef struct {
	Ref   string
	Value *Link
}
