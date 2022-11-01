// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"os"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Manager is a process manager maintaining multiple processes.
type Manager struct {
	processes *gmap.IntAnyMap // Process id to Process object mapping.
}

// NewManager creates and returns a new process manager.
func NewManager() *Manager {
	return &Manager{
		processes: gmap.NewIntAnyMap(true),
	}
}

// NewProcess creates and returns a Process object.
func (m *Manager) NewProcess(path string, args []string, environment []string) *Process {
	p := NewProcess(path, args, environment)
	p.Manager = m
	return p
}

// GetProcess retrieves and returns a Process object.
// It returns nil if it does not find the process with given `pid`.
func (m *Manager) GetProcess(pid int) *Process {
	if v := m.processes.Get(pid); v != nil {
		return v.(*Process)
	}
	return nil
}

// AddProcess adds a process to current manager.
// It does nothing if the process with given `pid` does not exist.
func (m *Manager) AddProcess(pid int) {
	if m.processes.Get(pid) == nil {
		if process, err := os.FindProcess(pid); err == nil {
			p := m.NewProcess("", nil, nil)
			p.Process = process
			m.processes.Set(pid, p)
		}
	}
}

// RemoveProcess removes a process from current manager.
func (m *Manager) RemoveProcess(pid int) {
	m.processes.Remove(pid)
}

// Processes retrieves and returns all processes in current manager.
func (m *Manager) Processes() []*Process {
	processes := make([]*Process, 0)
	m.processes.RLockFunc(func(m map[int]interface{}) {
		for _, v := range m {
			processes = append(processes, v.(*Process))
		}
	})
	return processes
}

// Pids retrieves and returns all process id array in current manager.
func (m *Manager) Pids() []int {
	return m.processes.Keys()
}

// WaitAll waits until all process exit.
func (m *Manager) WaitAll() {
	processes := m.Processes()
	if len(processes) > 0 {
		for _, p := range processes {
			_ = p.Wait()
		}
	}
}

// KillAll kills all processes in current manager.
func (m *Manager) KillAll() error {
	for _, p := range m.Processes() {
		if err := p.Kill(); err != nil {
			return err
		}
	}
	return nil
}

// SignalAll sends a signal `sig` to all processes in current manager.
func (m *Manager) SignalAll(sig os.Signal) error {
	for _, p := range m.Processes() {
		if err := p.Signal(sig); err != nil {
			err = gerror.Wrapf(err, `send signal to process failed for pid "%d" with signal "%s"`, p.Process.Pid, sig)
			return err
		}
	}
	return nil
}

// Send sends data bytes to all processes in current manager.
func (m *Manager) Send(data []byte) {
	for _, p := range m.Processes() {
		_ = p.Send(data)
	}
}

// SendTo sneds data bytes to specified processe in current manager.
func (m *Manager) SendTo(pid int, data []byte) error {
	return Send(pid, data)
}

// Clear removes all processes in current manager.
func (m *Manager) Clear() {
	m.processes.Clear()
}

// Size returns the size of processes in current manager.
func (m *Manager) Size() int {
	return m.processes.Size()
}
