package issue3835

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/issue/3835/service"

	"github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func init() {
	service.RegisterItest(New())
}

type sItest struct {
}

func New() *sItest {
	return &sItest{}
}

func (s *sItest) F(ctx context.Context) (d mysql.Driver, err error) {
	return mysql.Driver{}, nil
}
