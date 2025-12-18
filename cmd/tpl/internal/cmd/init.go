package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tpl/internal/logic"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

// 预定义模板列表
var defaultTemplates = []struct {
	Name string
	Repo string
	Desc string
}{
	{"template-single", "github.com/gogf/template-single", "单体项目模板"},
	{"template-mono", "github.com/gogf/template-mono", "大仓项目模板"},
}

var (
	Init = gcmd.Command{
		Name:        "init",
		Brief:       "Initialize a new project",
		Description: "Download a remote template and generate a new project.",
		Examples: `
            tpl init                                                                          # 交互式选择模板
            tpl init github.com/gogf/template-single my-project                               # 单体项目模板
            tpl init github.com/gogf/template-mono my-mono                                    # 大仓项目模板
            tpl init github.com/gogf/gf/cmd/gf/v2 my-gf                                       # 嵌套 Go 模块
            tpl init github.com/gogf/gf/cmd/gf/v2 my-gf -s                                    # 交互式版本选择
            tpl init github.com/gogf/gf/cmd/gf/v2@v2.8.0 my-gf                                # 指定版本
            tpl init github.com/gogf/template-single my-project -m github.com/myorg/myproject # 自定义模块路径
            tpl init github.com/gogf/examples/httpserver/jwt my-jwt                           # 子目录 (via git)
        `,
		Arguments: []gcmd.Argument{
			{Name: "repo", Short: "r", Brief: "Remote repository URL", IsArg: true},
			{Name: "name", Short: "n", Brief: "Project name (optional)", IsArg: true},
			{Name: "select", Short: "s", Brief: "Enable interactive version selection", Orphan: true},
			{Name: "mod", Short: "m", Brief: "Go module path (e.g., github.com/xxx/xxx)"},
			{Name: "upgrade", Short: "u", Brief: "Upgrade dependencies to latest (go get -u ./...)", Orphan: true},
		},
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			repo := parser.GetArg(2).String()
			name := parser.GetArg(3).String()
			selectMode := parser.GetOpt("select") != nil
			modPath := parser.GetOpt("mod").String()
			upgradeDeps := parser.GetOpt("upgrade") != nil

			// 如果没有提供 repo，进入交互模式
			if repo == "" {
				var upgrade bool
				repo, name, modPath, upgrade, err = interactiveMode()
				if err != nil {
					return err
				}
				upgradeDeps = upgrade
			}

			g.Log().Info(ctx, "Parsed args - Repo:", repo, "Name:", name, "SelectMode:", selectMode, "ModPath:", modPath, "UpgradeDeps:", upgradeDeps)

			// Call logic with options
			return logic.Process(ctx, repo, name, &logic.ProcessOptions{
				SelectVersion: selectMode,
				ModulePath:    modPath,
				UpgradeDeps:   upgradeDeps,
			})
		},
	}
)

// interactiveMode 交互式选择模板和输入项目信息
func interactiveMode() (repo, name, modPath string, upgradeDeps bool, err error) {
	reader := bufio.NewReader(os.Stdin)

	// 1. 选择模板
	fmt.Println("\n请选择项目模板:")
	fmt.Println(strings.Repeat("-", 50))
	for i, t := range defaultTemplates {
		fmt.Printf("  [%d] %s - %s\n", i+1, t.Name, t.Desc)
	}
	fmt.Println(strings.Repeat("-", 50))

	for {
		fmt.Printf("选择模板 [1-%d]: ", len(defaultTemplates))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		idx, e := strconv.Atoi(input)
		if e != nil || idx < 1 || idx > len(defaultTemplates) {
			fmt.Printf("无效选择，请输入 1-%d 之间的数字\n", len(defaultTemplates))
			continue
		}
		repo = defaultTemplates[idx-1].Repo
		fmt.Printf("已选择: %s\n\n", repo)
		break
	}

	// 2. 输入项目名称
	for {
		fmt.Print("请输入项目名称: ")
		input, _ := reader.ReadString('\n')
		name = strings.TrimSpace(input)
		if name == "" {
			fmt.Println("项目名称不能为空")
			continue
		}
		break
	}

	// 3. 输入模块路径（可选）
	fmt.Printf("请输入 Go 模块路径 (留空则使用项目名 \"%s\"): ", name)
	input, _ := reader.ReadString('\n')
	modPath = strings.TrimSpace(input)

	// 4. 是否升级依赖到最新
	fmt.Print("是否升级依赖到最新版本 (go get -u)? [y/N]: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	upgradeDeps = input == "y" || input == "yes"

	fmt.Println()
	return repo, name, modPath, upgradeDeps, nil
}

func init() {
	Main.AddCommand(&Init)
}
