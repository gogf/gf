package user

import (
	"context"

	bbb "github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genservice/logic/base"
	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genservice/service"
)

func init() {
	service.RegisterVipUser(&sVipUser{})
}

type mybase = bbb.Base

type sVipUser struct {
	sUser  `gen:"extend"`
	mybase `gen:"extend"`
}

// Create creates a new vip user.
func (s *sVipUser) Create(ctx context.Context, name string, vipLevel int) (id int, err error) {
	return 0, nil
}

// Get vip user level
func (s *sVipUser) GetVipLevel() (vipLevel int, err error) {
	return 1, nil
}
