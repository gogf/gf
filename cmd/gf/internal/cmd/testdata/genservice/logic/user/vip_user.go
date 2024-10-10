package user

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genservice/service"
)

func init() {
	service.RegisterVipUser(&sVipUser{})
}

type sVipUser struct {
	sUser
}

// Create creates a new vip user.
func (s *sVipUser) Create(ctx context.Context, name string, vipLevel int) (id int, err error) {
	return 0, nil
}

// Get vip user level
func (s *sVipUser) GetVipLevel() (vipLevel int, err error) {
	return 1, nil
}
