package balancer_test

import (
	"testing"

	"github.com/gogf/gf/contrib/balancer/v2"
	"github.com/gogf/gf/v2/net/gsel"
)

func Test_Register(t *testing.T) {
	balancer.Register("test", gsel.NewSelectorRandom())
}
