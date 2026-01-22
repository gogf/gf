package issue4242alias

import (
	"context"

	// Anonymous import (should be skipped)
	_ "github.com/gogf/gf/v2/os/gres"

	// Versioned import without alias
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/issue/4242/service"

	// Explicit alias import
	mysqlDriver "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func init() {
	service.RegisterIssue4242Alias(New())
}

type sIssue4242Alias struct {
}

func New() *sIssue4242Alias {
	return &sIssue4242Alias{}
}

// GetDriver tests explicit alias import.
func (s *sIssue4242Alias) GetDriver(ctx context.Context) (d mysqlDriver.Driver, err error) {
	return mysqlDriver.Driver{}, nil
}

// GetRequest tests versioned import.
func (s *sIssue4242Alias) GetRequest(ctx context.Context) (*ghttp.Request, error) {
	return nil, nil
}
