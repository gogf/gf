// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/geninit"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/allyes"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

var (
	// Init .
	Init = cInit{}
)

type cInit struct {
	g.Meta `name:"init" brief:"{cInitBrief}" eg:"{cInitEg}"`
}

const (
	cInitRepoPrefix  = `github.com/gogf/`
	cInitMonoRepo    = `template-mono`
	cInitMonoRepoApp = `template-mono-app`
	cInitSingleRepo  = `template-single`
	cInitBrief       = `create and initialize an empty GoFrame project`
	cInitEg          = `
gf init my-project
gf init my-mono-repo -m
gf init my-mono-repo -a
gf init my-project -u
gf init my-project -g "github.com/myorg/myproject"
gf init -r github.com/gogf/template-single my-project
gf init -r github.com/gogf/template-single my-project -s
gf init -r github.com/gogf/examples/httpserver/jwt my-jwt
gf init -i
`
	cInitNameBrief = `
name for the project. It will create a folder with NAME in current directory.
The NAME will also be the module name for the project.
`
	// cInitGitDir the git directory
	cInitGitDir = ".git"
	// cInitGitignore the gitignore file
	cInitGitignore = ".gitignore"
)

// defaultTemplates is the list of predefined templates for interactive selection
var defaultTemplates = []struct {
	Name string
	Repo string
	Desc string
}{
	{"template-single", "github.com/gogf/template-single", "Single project template"},
	{"template-mono", "github.com/gogf/template-mono", "Mono-repo project template"},
}

func init() {
	gtag.Sets(g.MapStrStr{
		`cInitBrief`:     cInitBrief,
		`cInitEg`:        cInitEg,
		`cInitNameBrief`: cInitNameBrief,
	})
}

type cInitInput struct {
	g.Meta      `name:"init"`
	Name        string `name:"NAME" arg:"true" brief:"{cInitNameBrief}"`
	Mono        bool   `name:"mono" short:"m" brief:"initialize a mono-repo instead a single-repo" orphan:"true"`
	MonoApp     bool   `name:"monoApp" short:"a" brief:"initialize a mono-repo-app instead a single-repo" orphan:"true"`
	Update      bool   `name:"update" short:"u" brief:"update to the latest goframe version" orphan:"true"`
	Module      string `name:"module" short:"g" brief:"custom go module"`
	Repo        string `name:"repo" short:"r" brief:"remote repository URL for template download"`
	SelectVer   bool   `name:"select" short:"s" brief:"enable interactive version selection for remote template" orphan:"true"`
	Interactive bool   `name:"interactive" short:"i" brief:"enable interactive mode to select template" orphan:"true"`
}

type cInitOutput struct{}

func (c cInit) Index(ctx context.Context, in cInitInput) (out *cInitOutput, err error) {
	// Check if using remote template mode
	if in.Repo != "" || in.Interactive {
		return c.initFromRemote(ctx, in)
	}

	// If no name provided and no remote mode, enter interactive mode
	if in.Name == "" {
		return c.initInteractive(ctx, in)
	}

	// Default: use built-in template
	return c.initFromBuiltin(ctx, in)
}

// initFromRemote initializes project from remote repository
func (c cInit) initFromRemote(ctx context.Context, in cInitInput) (out *cInitOutput, err error) {
	repo := in.Repo
	name := in.Name

	// If interactive mode and no repo specified, let user select
	if in.Interactive && repo == "" {
		var modPath string
		var upgradeDeps bool
		repo, name, modPath, upgradeDeps, err = interactiveSelectTemplate()
		if err != nil {
			return nil, err
		}
		if modPath != "" {
			in.Module = modPath
		}
		if upgradeDeps {
			in.Update = true
		}
	}

	if repo == "" {
		return nil, fmt.Errorf("repository URL is required for remote template mode")
	}

	// Default name to repo basename if empty
	if name == "" {
		name = gfile.Basename(repo)
		mlog.Printf("Using repository basename as project name: %s", name)
	}

	mlog.Print("initializing from remote template...")

	opts := &geninit.ProcessOptions{
		SelectVersion: in.SelectVer,
		ModulePath:    in.Module,
		UpgradeDeps:   in.Update,
	}

	if err = geninit.Process(ctx, repo, name, opts); err != nil {
		return nil, err
	}

	mlog.Print("initialization done!")
	if name != "" && name != "." {
		mlog.Printf(`you can now run "cd %s && gf run main.go" to start your journey, enjoy!`, name)
	}
	return
}

// initFromBuiltin initializes project from built-in template
func (c cInit) initFromBuiltin(ctx context.Context, in cInitInput) (out *cInitOutput, err error) {
	var overwrote = false
	if !gfile.IsEmpty(in.Name) && !allyes.Check() {
		s := gcmd.Scanf(`the folder "%s" is not empty, files might be overwrote, continue? [y/n]: `, in.Name)
		if strings.EqualFold(s, "n") {
			return
		}
		overwrote = true
	}
	mlog.Print("initializing...")

	// Create project folder and files.
	var (
		templateRepoName string
		gitignoreFile    = in.Name + "/" + cInitGitignore
	)

	if in.Mono {
		templateRepoName = cInitMonoRepo
	} else if in.MonoApp {
		templateRepoName = cInitMonoRepoApp
	} else {
		templateRepoName = cInitSingleRepo
	}

	err = gres.Export(templateRepoName, in.Name, gres.ExportOption{
		RemovePrefix: templateRepoName,
	})
	if err != nil {
		return
	}

	// build ignoreFiles from the .gitignore file
	ignoreFiles := make([]string, 0, 10)
	ignoreFiles = append(ignoreFiles, cInitGitDir)
	// in.MonoApp is a mono-repo-app, it should ignore the .gitignore file
	if overwrote && !in.MonoApp {
		err = gfile.ReadLines(gitignoreFile, func(line string) error {
			// Add only hidden files or directories
			// If other directories are added, it may cause the entire directory to be ignored
			// such as 'main' in the .gitignore file, but the path is ' D:\main\my-project '
			if line != "" && strings.HasPrefix(line, ".") {
				ignoreFiles = append(ignoreFiles, line)
			}
			return nil
		})

		// if not found the .gitignore file will skip os.ErrNotExist error
		if err != nil && !os.IsNotExist(err) {
			return
		}
	}

	// Get template name and module name.
	if in.Module == "" {
		in.Module = gfile.Basename(gfile.RealPath(in.Name))
	}
	if in.MonoApp {
		pwd := gfile.Pwd() + string(os.PathSeparator) + in.Name
		in.Module = utils.GetImportPath(pwd)
	}

	// Replace template name to project name.
	err = gfile.ReplaceDirFunc(func(path, content string) string {
		for _, ignoreFile := range ignoreFiles {
			if strings.Contains(path, ignoreFile) {
				return content
			}
		}
		return gstr.Replace(gfile.GetContents(path), cInitRepoPrefix+templateRepoName, in.Module)
	}, in.Name, "*", true)
	if err != nil {
		return
	}

	// Update the GoFrame version.
	if in.Update {
		mlog.Print("update goframe...")
		// go get -u github.com/gogf/gf/v2@latest
		updateCommand := `go get -u github.com/gogf/gf/v2@latest`
		if in.Name != "." {
			updateCommand = fmt.Sprintf(`cd %s && %s`, in.Name, updateCommand)
		}
		if err = gproc.ShellRun(ctx, updateCommand); err != nil {
			mlog.Fatal(err)
		}
		// go mod tidy
		gomModTidyCommand := `go mod tidy`
		if in.Name != "." {
			gomModTidyCommand = fmt.Sprintf(`cd %s && %s`, in.Name, gomModTidyCommand)
		}
		if err = gproc.ShellRun(ctx, gomModTidyCommand); err != nil {
			mlog.Fatal(err)
		}
	}

	mlog.Print("initialization done! ")
	if !in.Mono {
		enjoyCommand := `gf run main.go`
		if in.Name != "." {
			enjoyCommand = fmt.Sprintf(`cd %s && %s`, in.Name, enjoyCommand)
		}
		mlog.Printf(`you can now run "%s" to start your journey, enjoy!`, enjoyCommand)
	}
	return
}

// initInteractive enters interactive mode when no arguments provided
func (c cInit) initInteractive(ctx context.Context, in cInitInput) (out *cInitOutput, err error) {
	reader := bufio.NewReader(os.Stdin)

	// Ask user which mode to use
	fmt.Println("\nPlease select initialization mode:")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("  [1] Built-in template (default)")
	fmt.Println("  [2] Remote template")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Print("Select mode [1-2] (default: 1): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		mlog.Fatalf("failed to read input: %v", err)
		return
	}
	input = strings.TrimSpace(input)

	if input == "2" {
		in.Interactive = true
		return c.initFromRemote(ctx, in)
	}

	// Built-in template mode
	fmt.Println("\nPlease select project type:")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("  [1] Single project (default)")
	fmt.Println("  [2] Mono-repo project")
	fmt.Println("  [3] Mono-repo app")
	fmt.Println(strings.Repeat("-", 50))

	fmt.Print("Select type [1-3] (default: 1): ")
	input, err = reader.ReadString('\n')
	if err != nil {
		mlog.Fatalf("failed to read input: %v", err)
		return
	}
	input = strings.TrimSpace(input)

	switch input {
	case "2":
		in.Mono = true
	case "3":
		in.MonoApp = true
	}

	// Get project name
	for {
		fmt.Print("Enter project name: ")
		input, err = reader.ReadString('\n')
		if err != nil {
			mlog.Fatalf("failed to read input: %v", err)
			return
		}
		in.Name = strings.TrimSpace(input)
		if in.Name != "" {
			break
		}
		fmt.Println("Project name cannot be empty")
	}

	// Get module path (optional)
	fmt.Printf("Enter Go module path (leave empty to use \"%s\"): ", in.Name)
	input, err = reader.ReadString('\n')
	if err != nil {
		mlog.Fatalf("failed to read input: %v", err)
		return
	}
	in.Module = strings.TrimSpace(input)

	// Ask about update
	fmt.Print("Update to latest GoFrame version? [y/N]: ")
	input, err = reader.ReadString('\n')
	if err != nil {
		mlog.Fatalf("failed to read input: %v", err)
		return
	}
	input = strings.TrimSpace(strings.ToLower(input))
	in.Update = input == "y" || input == "yes"

	fmt.Println()
	return c.initFromBuiltin(ctx, in)
}

// interactiveSelectTemplate prompts user to select a template interactively
func interactiveSelectTemplate() (repo, name, modPath string, upgradeDeps bool, err error) {
	reader := bufio.NewReader(os.Stdin)

	// 1. Select template
	fmt.Println("\nPlease select a project template:")
	fmt.Println(strings.Repeat("-", 50))
	for i, t := range defaultTemplates {
		fmt.Printf("  [%d] %s - %s\n", i+1, t.Name, t.Desc)
	}
	fmt.Printf("  [%d] Custom repository URL\n", len(defaultTemplates)+1)
	fmt.Println(strings.Repeat("-", 50))

	for {
		fmt.Printf("Select template [1-%d]: ", len(defaultTemplates)+1)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", "", "", false, fmt.Errorf("failed to read template selection: %w", err)
		}
		input = strings.TrimSpace(input)

		idx, e := strconv.Atoi(input)
		if e != nil || idx < 1 || idx > len(defaultTemplates)+1 {
			fmt.Printf("Invalid selection, please enter a number between 1-%d\n", len(defaultTemplates)+1)
			continue
		}

		if idx <= len(defaultTemplates) {
			repo = defaultTemplates[idx-1].Repo
			fmt.Printf("Selected: %s\n\n", repo)
		} else {
			// Custom URL
			fmt.Print("Enter repository URL: ")
			input, err = reader.ReadString('\n')
			if err != nil {
				return "", "", "", false, fmt.Errorf("failed to read repository URL: %w", err)
			}
			repo = strings.TrimSpace(input)
			if repo == "" {
				fmt.Println("Repository URL cannot be empty")
				continue
			}
		}
		break
	}

	// 2. Enter project name
	for {
		fmt.Print("Enter project name: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", "", "", false, fmt.Errorf("failed to read project name: %w", err)
		}
		name = strings.TrimSpace(input)
		if name == "" {
			fmt.Println("Project name cannot be empty")
			continue
		}
		break
	}

	// 3. Enter module path (optional)
	fmt.Printf("Enter Go module path (leave empty to use \"%s\"): ", name)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", false, fmt.Errorf("failed to read module path: %w", err)
	}
	modPath = strings.TrimSpace(input)

	// 4. Ask about upgrade
	fmt.Print("Upgrade dependencies to latest (go get -u)? [y/N]: ")
	input, err = reader.ReadString('\n')
	if err != nil {
		return "", "", "", false, fmt.Errorf("failed to read upgrade confirmation: %w", err)
	}
	input = strings.TrimSpace(strings.ToLower(input))
	upgradeDeps = input == "y" || input == "yes"

	fmt.Println()
	return repo, name, modPath, upgradeDeps, nil
}
