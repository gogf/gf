package resolver

import (
	"google.golang.org/grpc/resolver"
)

// Register registers the resolver builder to the resolver map. b.Scheme will be
// used as the scheme registered with this builder.
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe. If multiple Resolvers are
// registered with the same name, the one registered last will take effect.
func Register(builder resolver.Builder) {
	resolver.Register(builder)
}
