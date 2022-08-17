/**
 * @Author ZhangDafang
 * @create 2022-08-15 23:02
 */

package glimiter

import "context"

type Option struct {

	// max connection 
	MaxConnections int64

	//max qps
	MaxQPS         int64

	//filter , include user_id , ip etc .
	ReqFilterInfo   	   interface{}
}

type Adaper interface {

	// All interfaces add context , you can add special parameters through the content value, and make judgments through RefFilterInfo in Adapter, such as ip, user_id, etc
	// is acquire 
	Acquire(ctx context.Context) bool
	// reset limiter status
	ResetStatus(ctx context.Context) 
	// get current status 
	Status(ctx context.Context) (currentConnctions , currentQPS int64) 
	// update options 
	UpdateOption(options []*Option) bool 
}


type localAdapter = Adaper

type Limiter struct {
	localAdapter 
}

func NewLimiter(adapter Adaper, options... *Option ) *Limiter {

	if len(options) > 0{
		adapter.UpdateOption(options)
	}

	return &Limiter{
		localAdapter: adapter,
	}
}

func (l *Limiter) Acquire(ctx context.Context) bool {
	return l.Acquire(ctx)
}

func (l *Limiter) ResetStatus(ctx context.Context) {
	l.ResetStatus(ctx)
}

func (l *Limiter)Status(ctx context.Context) (currentConnctions , currentQPS int64)  {
	return l.Status(ctx)
}

func (l *Limiter) UpdateOption(ctx context.Context, options []*Option ) bool {
	return l.UpdateOption(ctx, options)
}


