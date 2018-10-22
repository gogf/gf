package yaml_test

import (
	. "gitee.com/johng/gf/third/gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})
