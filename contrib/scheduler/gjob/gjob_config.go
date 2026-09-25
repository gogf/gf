// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Configuration-based helper for creating job servers.

package gjob

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gapp"
	"github.com/gogf/gf/v2/util/gconv"
)

// configKey is the configuration key for job tasks.
const configKey = "scheduler.job"

// taskTypeWorker identifies worker tasks in configuration.
const taskTypeWorker = "worker"

// taskTypeCron identifies cron tasks in configuration.
const taskTypeCron = "cron"

// jobConfig holds the configuration for a single job task.
type jobConfig struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Enable bool   `json:"enable"`
	Spec   string `json:"spec"`
}

// HandlerMap maps task names to their handler functions.
// The value type depends on the task type:
//   - worker: WorkerHandler (func(ctx context.Context) func())
//   - cron: CronHandler (func(ctx context.Context) error)
type HandlerMap map[string]any

// NewServersFromConfig creates gapp.Server instances based on configuration.
// It reads job configurations from the "scheduler.job" config key and matches
// them with handlers from the provided HandlerMap.
//
// Configuration format (YAML example):
//
//	scheduler:
//	  job:
//	    - name: my-worker
//	      type: worker
//	      enable: true
//	    - name: my-cron
//	      type: cron
//	      enable: true
//	      spec: "*/2 * * * * *"
//
// The returned slice contains only the servers that have at least one task.
// Disabled tasks and tasks without matching handlers are skipped.
func NewServersFromConfig(ctx context.Context, handlers HandlerMap) []gapp.Server {
	var (
		configs     = loadConfig(ctx)
		workerTasks []WorkerTask
		cronTasks   []CronTask
	)

	for _, cfg := range configs {
		if !cfg.Enable {
			continue
		}

		handler, ok := handlers[cfg.Name]
		if !ok {
			g.Log().Warningf(ctx, "job config %s: handler not found, skipping", cfg.Name)
			continue
		}

		switch cfg.Type {
		case taskTypeWorker:
			h, ok := handler.(WorkerHandler)
			if !ok {
				g.Log().Warningf(ctx, "job config %s: handler is not WorkerHandler, skipping", cfg.Name)
				continue
			}
			workerTasks = append(workerTasks, WorkerTask{
				Name:    cfg.Name,
				Handler: h,
			})

		case taskTypeCron:
			if cfg.Spec == "" {
				g.Log().Warningf(ctx, "job config %s: missing spec for cron task, skipping", cfg.Name)
				continue
			}
			h, ok := handler.(CronHandler)
			if !ok {
				g.Log().Warningf(ctx, "job config %s: handler is not CronHandler, skipping", cfg.Name)
				continue
			}
			cronTasks = append(cronTasks, CronTask{
				Name:    cfg.Name,
				Spec:    cfg.Spec,
				Handler: h,
			})

		default:
			g.Log().Warningf(ctx, "job config %s: unknown type %q, skipping", cfg.Name, cfg.Type)
		}
	}

	var servers []gapp.Server

	if len(workerTasks) > 0 {
		servers = append(servers, NewWorkerServer(ctx, workerTasks...))
	}

	if len(cronTasks) > 0 {
		servers = append(servers, NewCronServer(ctx, cronTasks...))
	}

	return servers
}

// loadConfig reads job configurations from the application config file.
func loadConfig(ctx context.Context) []*jobConfig {
	cfg, err := g.Cfg().Get(ctx, configKey)
	if err != nil {
		g.Log().Errorf(ctx, "failed to get job config: %v", err)
		return nil
	}

	if cfg.IsNil() || cfg.IsEmpty() {
		return nil
	}

	maps := cfg.Maps()
	if len(maps) == 0 {
		return nil
	}

	var configs []*jobConfig
	for _, m := range maps {
		jobType := gconv.String(m["type"])
		if jobType == "" {
			name := gconv.String(m["name"])
			g.Log().Warningf(ctx, "job config %s: missing type, skipping", name)
			continue
		}

		if jobType != taskTypeWorker && jobType != taskTypeCron {
			name := gconv.String(m["name"])
			g.Log().Warningf(ctx, "job config %s: invalid type %q, skipping", name, jobType)
			continue
		}

		configs = append(configs, &jobConfig{
			Name:   gconv.String(m["name"]),
			Type:   jobType,
			Enable: gconv.Bool(m["enable"]),
			Spec:   gconv.String(m["spec"]),
		})
	}

	return configs
}

// String returns a human-readable description of the HandlerMap.
func (h HandlerMap) String() string {
	return fmt.Sprintf("HandlerMap(%d handlers)", len(h))
}
