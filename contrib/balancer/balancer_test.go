package balancer_test

import (
	"testing"

	"github.com/gogf/gf/v2/contrib/balancer"
	"github.com/gogf/gf/v2/net/gsel"
)

func Test_Register(t *testing.T) {
	balancer.Register("test", gsel.NewSelectorRandom())
}
