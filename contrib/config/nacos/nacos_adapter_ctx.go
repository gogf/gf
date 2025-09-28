package nacos

import (
	"context"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	// KeyNamespace is the context key for namespace
	KeyNamespace gctx.StrKey = "namespace"
	// KeyGroup is the context key for group
	KeyGroup gctx.StrKey = "group"
	// KeyDataId is the context key for dataId
	KeyDataId gctx.StrKey = "dataId"
)

const (
	// OperationUpdate represents the update operation
	OperationUpdate = "update"
)

// NacosAdapterCtx is the context adapter for Nacos configuration
type NacosAdapterCtx struct {
	Ctx context.Context
}

// NewNacosAdapterCtx creates and returns a new NacosAdapterCtx.
// If context is provided, it will be used; otherwise, a background context is created.
func NewNacosAdapterCtx(ctx ...context.Context) *NacosAdapterCtx {
	if len(ctx) > 0 {
		return &NacosAdapterCtx{
			Ctx: ctx[0],
		}
	}
	return &NacosAdapterCtx{Ctx: context.Background()}
}

// GetNacosAdapterCtx creates a new NacosAdapterCtx with the given context
func GetNacosAdapterCtx(ctx context.Context) *NacosAdapterCtx {
	return &NacosAdapterCtx{Ctx: ctx}
}

// WithOperation sets the operation in the context
func (n *NacosAdapterCtx) WithOperation(operation ...string) *NacosAdapterCtx {
	if len(operation) > 0 {
		n.Ctx = context.WithValue(n.Ctx, gcfg.OperationWrite, operation[0])
	}
	return n
}

// WithNamespace sets the namespace in the context
func (n *NacosAdapterCtx) WithNamespace(namespace ...string) *NacosAdapterCtx {
	if len(namespace) > 0 {
		n.Ctx = context.WithValue(n.Ctx, KeyNamespace, namespace[0])
	}
	return n
}

// WithGroup sets the group in the context
func (n *NacosAdapterCtx) WithGroup(group ...string) *NacosAdapterCtx {
	if len(group) > 0 {
		n.Ctx = context.WithValue(n.Ctx, KeyGroup, group[0])
	}
	return n
}

// WithDataId sets the dataId in the context
func (n *NacosAdapterCtx) WithDataId(dataId ...string) *NacosAdapterCtx {
	if len(dataId) > 0 {
		n.Ctx = context.WithValue(n.Ctx, KeyDataId, dataId[0])
	}
	return n
}

// WithSetContent sets the content in the context
func (n *NacosAdapterCtx) WithSetContent(content ...string) *NacosAdapterCtx {
	if len(content) > 0 {
		n.Ctx = context.WithValue(n.Ctx, gcfg.KeySetContent, content[0])
	}
	return n
}

// GetNamespace retrieves the namespace from the context
func (n *NacosAdapterCtx) GetNamespace() string {
	if v := n.Ctx.Value(KeyNamespace); v != nil {
		return v.(string)
	}
	return ""
}

// GetGroup retrieves the group from the context
func (n *NacosAdapterCtx) GetGroup() string {
	if v := n.Ctx.Value(KeyGroup); v != nil {
		return v.(string)
	}
	return ""
}

// GetDataId retrieves the dataId from the context
func (n *NacosAdapterCtx) GetDataId() string {
	if v := n.Ctx.Value(KeyDataId); v != nil {
		return v.(string)
	}
	return ""
}

// GetSetContent retrieves the content from the context
func (n *NacosAdapterCtx) GetSetContent() string {
	if v := n.Ctx.Value(gcfg.KeySetContent); v != nil {
		return v.(string)
	}
	return ""
}

// GetOperation retrieves the operation from the context
func (n *NacosAdapterCtx) GetOperation() string {
	if v := n.Ctx.Value(OperationUpdate); v != nil {
		return v.(string)
	}
	return ""
}
