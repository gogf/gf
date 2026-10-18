package pointer

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genservice/service"
)

type sPointer struct{}

func (s *sPointer) TestPointer(ctx context.Context, opts *service.PointerOptions) error {
	return nil
}
