package redis

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gmlock"
	"io"
	"strings"
)

var cacheScriptLock = gmlock.New()

type Script struct {
	src, hash string
}

func NewScript(src string) *Script {
	h := sha1.New()
	_, _ = io.WriteString(h, src)
	return &Script{
		src:  src,
		hash: hex.EncodeToString(h.Sum(nil)),
	}
}

func (s *Script) Hash() string {
	return s.hash
}

func (s *Script) Load(ctx context.Context, c gredis.IGroupScript) (string, error) {
	return c.ScriptLoad(ctx, s.src)
}

func (s *Script) Exist(ctx context.Context, c gredis.IGroupScript) (bool, error) {
	exists, err := c.ScriptExists(ctx, s.hash)
	if err != nil {
		return false, err
	}
	return exists[s.hash], nil
}

func (s *Script) Eval(ctx context.Context, c gredis.IGroupScript, keys []string, args ...interface{}) (*gvar.Var, error) {
	return c.Eval(ctx, s.src, int64(len(keys)), keys, args)
}

func (s *Script) EvalSha(ctx context.Context, c gredis.IGroupScript, keys []string, args ...interface{}) (*gvar.Var, error) {
	return c.EvalSha(ctx, s.hash, int64(len(keys)), keys, args)
}

// Run optimistically uses EVALSHA to run the script. If script does not exist
// it is retried using EVAL and cache the script
func (s *Script) Run(ctx context.Context, c gredis.IGroupScript, keys []string, args ...interface{}) (val *gvar.Var, err error) {
	if val, err = s.EvalSha(ctx, c, keys, args...); err != nil {
		if strings.Contains(err.Error(), "NOSCRIPT") {
			go s.tryCacheScript(c)
			return s.Eval(ctx, c, keys, args...)
		}
	}
	return
}

func (s *Script) tryCacheScript(c gredis.IGroupScript) {
	cacheScriptLock.TryLockFunc(s.hash, func() {
		_, _ = s.Load(context.Background(), c)
	})
}
