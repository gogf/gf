package v1

import "github.com/gogf/gf/v2/frame/g"

type HelloReq struct {
	g.Meta `path:"/hello" tags:"Hello" method:"get" summary:"api-Hello"`
}
type HelloRes struct {
	g.Meta `mime:"text/html" example:"string"`
}

// type Hello2Req struct {
//	g.Meta `path:"/hello2" tags:"Hello" method:"get" summary:"api-Hello2"`
// }
// type Hello2Res struct {
//	g.Meta `mime:"text/html" example:"string"`
// }

//type Hello3Req struct {
//	g.Meta   `path:"/hello3" tags:"Hello" method:"get" summary:"api-Hello3"`
//	PodState consts.PodPhase
//}
//type Hello3Res struct {
//}
