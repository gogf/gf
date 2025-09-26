package json

import (
	"testing"

	"github.com/gogf/gf/v2/errors/gerror"
)

func TestError_MarshalJSON(t *testing.T) {
	errNormal := gerror.New("test")
	errWithSign := gerror.New("test \"\"")
	_, err := Marshal(errNormal)
	if err != nil {
		t.Fail()
	}
	_, err = Marshal(errWithSign)
	if err != nil {
		t.Fail()
	}
}
