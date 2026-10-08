// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/agents"
)

var (
	// Agents is the `gf agents` command.
	Agents = cAgents{}
)

type cAgents struct {
	g.Meta `name:"agents" brief:"link local AI coding-agent paths to .agents and AGENTS.md" eg:"{cAgentsEg}"`
}

const (
	cAgentsEg = `
gf agents
gf agents --agent=claude
gf agents --agent=claude --action=unlink
gf agents --agent=claude --force
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cAgentsEg`: cAgentsEg,
	})
}

type cAgentsInput struct {
	g.Meta `name:"agents"`
	Agent  string `name:"agent" short:"a" brief:"single agent name, for example claude"`
	Action string `name:"action" brief:"link or unlink, default link"`
	Force  bool   `name:"force" short:"f" brief:"rebuild mismatched links and git symlink stubs" orphan:"true"`
}

type cAgentsOutput struct{}

// Index runs gf agents.
func (c cAgents) Index(_ context.Context, in cAgentsInput) (*cAgentsOutput, error) {
	err := agents.Run(agents.Request{
		Agent:  in.Agent,
		Action: in.Action,
		Force:  in.Force,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
	})
	if err != nil {
		return nil, err
	}
	return &cAgentsOutput{}, nil
}
