package issue4242

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/issue/4242/service"

	"github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func init() {
	service.RegisterIssue4242(New())
}

type sIssue4242 struct {
}

func New() *sIssue4242 {
	return &sIssue4242{}
}

// GetDriver tests versioned import path is preserved.
func (s *sIssue4242) GetDriver(ctx context.Context) (d mysql.Driver, err error) {
	return mysql.Driver{}, nil
}

// GetRequest tests another versioned import.
func (s *sIssue4242) GetRequest(ctx context.Context) (*ghttp.Request, error) {
	g.Log().Info(ctx, "getting request")
	return nil, nil
}
