package user

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genservice/service"
)

func init() {
	service.RegisterSuperVipUser(&sSuperVipUser{
		sVipUser: &sVipUser{},
	})
}

type sSuperVipUser struct {
	*sVipUser `gen:"extend"`
}

// Get supper vip user level
func (s sSuperVipUser) GetVipLevel(ctx context.Context) (vipLevel int, err error) {
	return 1, nil
}

// Set supper vip user level
func (s *sSuperVipUser) SetVipLevel(ctx context.Context, id int, vipLevel int) (err error) {
	return nil
}
