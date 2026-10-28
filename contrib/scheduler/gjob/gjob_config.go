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
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
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
// Disabled tasks and tasks without a matching handler are skipped.
// A config read error, malformed job list, an unknown task type, a missing cron spec, or a handler
// type mismatch is returned so startup can fail instead of silently dropping jobs.
func NewServersFromConfig(ctx context.Context, handlers HandlerMap) ([]gapp.Server, error) {
	configs, err := loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	var (
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
			if !ok || h == nil {
				return nil, gerror.NewCodef(
					gcode.CodeInvalidConfiguration,
					"job config %s: handler must be a non-nil WorkerHandler",
					cfg.Name,
				)
			}
			workerTasks = append(workerTasks, WorkerTask{
				Name:    cfg.Name,
				Handler: h,
			})

		case taskTypeCron:
			if cfg.Spec == "" {
				return nil, gerror.NewCodef(
					gcode.CodeInvalidConfiguration,
					"job config %s: missing spec for cron task",
					cfg.Name,
				)
			}
			h, ok := handler.(CronHandler)
			if !ok || h == nil {
				return nil, gerror.NewCodef(
					gcode.CodeInvalidConfiguration,
					"job config %s: handler must be a non-nil CronHandler",
					cfg.Name,
				)
			}
			cronTasks = append(cronTasks, CronTask{
				Name:    cfg.Name,
				Spec:    cfg.Spec,
				Handler: h,
			})

		default:
			return nil, gerror.NewCodef(
				gcode.CodeInvalidConfiguration,
				"job config %s: unknown type %q",
				cfg.Name,
				cfg.Type,
			)
		}
	}

	var servers []gapp.Server
	if len(workerTasks) > 0 {
		servers = append(servers, NewWorkerServer(ctx, workerTasks...))
	}
	if len(cronTasks) > 0 {
		servers = append(servers, NewCronServer(ctx, cronTasks...))
	}
	return servers, nil
}

// loadConfig reads job configurations from the application config file.
func loadConfig(ctx context.Context) ([]*jobConfig, error) {
	cfg, err := g.Cfg().Get(ctx, configKey)
	if err != nil {
		return nil, gerror.WrapCode(gcode.CodeInvalidConfiguration, err, "failed to get job config")
	}
	if cfg.IsNil() {
		return nil, nil
	}

	value := reflect.ValueOf(cfg.Val())
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return nil, gerror.NewCode(gcode.CodeInvalidConfiguration, "scheduler.job must be a list of job objects")
	}
	// Maps conversion skips invalid elements, so validate the original collection first.
	for i := 0; i < value.Len(); i++ {
		if reflect.ValueOf(value.Index(i).Interface()).Kind() != reflect.Map {
			return nil, gerror.NewCodef(gcode.CodeInvalidConfiguration, "scheduler.job[%d] must be a job object", i)
		}
	}

	maps := cfg.Maps()
	configs := make([]*jobConfig, 0, len(maps))
	for _, m := range maps {
		var (
			name    = gconv.String(m["name"])
			jobType = gconv.String(m["type"])
		)
		if jobType == "" {
			return nil, gerror.NewCodef(gcode.CodeInvalidConfiguration, "job config %s: missing type", name)
		}
		if jobType != taskTypeWorker && jobType != taskTypeCron {
			return nil, gerror.NewCodef(
				gcode.CodeInvalidConfiguration,
				"job config %s: unknown type %q",
				name,
				jobType,
			)
		}
		configs = append(configs, &jobConfig{
			Name:   name,
			Type:   jobType,
			Enable: gconv.Bool(m["enable"]),
			Spec:   gconv.String(m["spec"]),
		})
	}
	return configs, nil
}

// String returns a human-readable description of the HandlerMap.
func (h HandlerMap) String() string {
	return fmt.Sprintf("HandlerMap(%d handlers)", len(h))
}
