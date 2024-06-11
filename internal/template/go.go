package template

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/tsukinoko-kun/codehere/internal/settings"
	"github.com/tsukinoko-kun/codehere/internal/util"
)

const (
	mainGoFile = `package main

import "fmt"

func main() {
	fmt.Println("hello world")
}
`

	goEditorconfigFile = `# EditorConfig is awesome: https://EditorConfig.org
root = true

[*]
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[{Makefile,go.mod,go.sum,*.go,*.templ,.gitmodules}]
charset = utf-8
indent_style = tab
tab_width = 4
indent_size = 4

[*.md]
indent_size = 4
trim_trailing_whitespace = false

eclint_indent_style = unset

[Dockerfile]
indent_style = space
indent_size = 4
`

	goGitignore = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work
go.work.sum

# env file
.env

# Editor config
**/.idea
**/.vs
**/.vscode
**/.fleet
`
)

func create_Go(s settings.Init) error {
	var packageName string

	if gh, err := util.GhUsername(); err != nil {
		packageName = s.Name
	} else {
		packageName = fmt.Sprintf("github.com/%s/%s", gh, s.Name)
	}

	f := huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("module name").Value(&packageName),
	).Title("Go"))

	if err := f.Run(); err != nil {
		return err
	}

	if _, err := util.WriteFile(filepath.Join(s.Loc, "main.go"), mainGoFile); err != nil {
		return err
	}

	goModLoc := filepath.Join(s.Loc, "go.mod")
	if exist, err := util.Exist(goModLoc); err != nil {
		return err
	} else if exist {
		if err := os.Remove(goModLoc); err != nil {
			return err
		}
	}

	if err := exec.Command("go", "mod", "init", packageName).Run(); err != nil {
		return errors.Join(fmt.Errorf("failed to go mod init %q", packageName))
	}

	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return errors.Join(errors.New("failed to go mod tidy"))
	}

	if _, err := util.WriteFile(filepath.Join(s.Loc, ".editorconfig"), goEditorconfigFile); err != nil {
		return err
	}

	if _, err := util.WriteFile(filepath.Join(s.Loc, ".gitignore"), goGitignore); err != nil {
		return err
	}

	return nil
}
