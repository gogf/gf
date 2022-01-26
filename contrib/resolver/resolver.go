package resolver

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"google.golang.org/grpc/resolver"
)

const (
	RawSvcKeyInSubConnInfo = `RawService`
)

func init() {
	// It uses default builder handling the DNS for grpc service requests.
	resolver.Register(&Builder{})
}

// SetRegistry sets the default Registry implements as your own implemented interface.
func SetRegistry(registry gsvc.Registry) {
	if registry == nil {
		panic(gerror.New(`invalid Registry value "nil" given`))
	}
	gsvc.SetRegistry(registry)
}
