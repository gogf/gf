// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package zookeeper

import (
	"context"
	"path"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
)

// serviceInstanceLeaf returns a unique ZooKeeper leaf name for one service instance.
// Derived from endpoints so multiple instances of the same service name can coexist
// under the same service prefix path (#4149 / #4717).
func serviceInstanceLeaf(service gsvc.Service) string {
	return strings.NewReplacer(":", "-", ",", "_").Replace(service.GetEndpoints().String())
}

// Register registers `service` to Registry.
// Note that it returns a new Service if it changes the input Service with custom one.
func (r *Registry) Register(_ context.Context, service gsvc.Service) (gsvc.Service, error) {
	var (
		data []byte
		err  error
	)
	if err = r.ensureName(r.opts.namespace, []byte(""), 0); err != nil {
		return service, gerror.Wrapf(
			err,
			"Error Creat node which name is %s",
			r.opts.namespace,
		)
	}
	prefix := strings.Trim(strings.ReplaceAll(service.GetPrefix(), "/", "-"), "-")
	servicePrefixPath := path.Join(r.opts.namespace, prefix)
	if err = r.ensureName(servicePrefixPath, []byte(""), 0); err != nil {
		return service, gerror.Wrapf(
			err,
			"Error Creat node which name is %s",
			servicePrefixPath,
		)
	}

	if data, err = marshal(&Content{
		Key:   service.GetKey(),
		Value: service.GetValue(),
	}); err != nil {
		return service, gerror.Wrapf(
			err,
			"Error with marshal Content to Json string",
		)
	}
	// Unique leaf per instance (endpoint), not per service name, so multi-instance
	// registration does not delete the previous instance's ephemeral node.
	servicePath := path.Join(servicePrefixPath, serviceInstanceLeaf(service))
	if err = r.ensureName(servicePath, data, zk.FlagEphemeral); err != nil {
		return service, gerror.Wrapf(
			err,
			"Error Creat node which name is %s",
			servicePath,
		)
	}
	go r.reRegister(servicePath, data)
	return service, nil
}

// Deregister off-lines and removes `service` from the Registry.
func (r *Registry) Deregister(ctx context.Context, service gsvc.Service) error {
	ch := make(chan error, 1)
	prefix := strings.Trim(strings.ReplaceAll(service.GetPrefix(), "/", "-"), "-")
	servicePath := path.Join(r.opts.namespace, prefix, serviceInstanceLeaf(service))
	go func() {
		err := r.conn.Delete(servicePath, -1)
		ch <- err
	}()
	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ch:
	}
	return gerror.Wrapf(err,
		"Error with deregister service:%s",
		service.GetName(),
	)
}

// ensureName ensure node exists, if not exist, create and set data
func (r *Registry) ensureName(path string, data []byte, flags int32) error {
	exists, stat, err := r.conn.Exists(path)
	if err != nil {
		return gerror.Wrapf(err,
			"Error with check node exist which name is %s",
			path,
		)
	}
	// Ephemeral nodes: re-create after crash/restart. Only delete when the node
	// already exists (stat valid) — also avoids wiping a sibling instance once
	// each instance uses a unique leaf path.
	if exists && flags&zk.FlagEphemeral == zk.FlagEphemeral {
		err = r.conn.Delete(path, stat.Version)
		if err != nil && err != zk.ErrNoNode {
			return gerror.Wrapf(err,
				"Error with delete node which name is %s",
				path,
			)
		}
		exists = false
	}
	if !exists {
		if len(r.opts.user) > 0 && len(r.opts.password) > 0 {
			_, err = r.conn.Create(path, data, flags, zk.DigestACL(zk.PermAll, r.opts.user, r.opts.password))
		} else {
			_, err = r.conn.Create(path, data, flags, zk.WorldACL(zk.PermAll))
		}
		if err != nil {
			return gerror.Wrapf(err,
				"Error with create node which name is %s",
				path,
			)
		}
	}
	return nil
}

// reRegister re-register data node info when bad connection recovered
func (r *Registry) reRegister(path string, data []byte) {
	sessionID := r.conn.SessionID()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		cur := r.conn.SessionID()
		// sessionID changed
		if cur > 0 && sessionID != cur {
			// re-ensureName
			if err := r.ensureName(path, data, zk.FlagEphemeral); err != nil {
				return
			}
			sessionID = cur
		}
	}
}
