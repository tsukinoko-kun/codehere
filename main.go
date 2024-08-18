package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/tsukinoko-kun/codehere/internal/license"
	"github.com/tsukinoko-kun/codehere/internal/open"
	"github.com/tsukinoko-kun/codehere/internal/settings"
	"github.com/tsukinoko-kun/codehere/internal/template"
	"github.com/tsukinoko-kun/codehere/internal/util"
)

func main() {
	initSettings := settings.Init{}
	initSettings.GitInit = true

	// default value for location
	if len(os.Args) >= 2 {
		initSettings.Loc = os.Args[1]
	} else {
		if wd, err := os.Getwd(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
			return
		} else {
			initSettings.Loc = wd
		}
	}

	// abs location
	if absProjLoc, err := filepath.Abs(initSettings.Loc); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
		return
	} else {
		initSettings.Loc = absProjLoc
		initSettings.Name = filepath.Base(absProjLoc)
	}

	// user input
	f := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Project location").
			Value(&initSettings.Loc),
		huh.NewSelect[uint8]().
			Title("Template").
			Value(&initSettings.Template).
			Options(
				huh.NewOption("Go", template.Template_Go),
			),
		huh.NewConfirm().Title("Git Init").Value(&initSettings.GitInit),
	).Title("CodeHere"))

	if err := f.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
		return
	}

	// abs location
	if absProjLoc, err := filepath.Abs(initSettings.Loc); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
		return
	} else {
		initSettings.Loc = absProjLoc
		initSettings.Name = filepath.Base(absProjLoc)
	}

	// mkdir
	if err := os.MkdirAll(initSettings.Loc, 0777); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(5)
		return
	}

	// cd
	if err := os.Chdir(initSettings.Loc); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(6)
		return
	}

	// template
	if err := template.Create(initSettings); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(7)
		return
	}

	// license
	if err := license.Create(initSettings); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(8)
		return
	}

	// git init
	if initSettings.GitInit {
		if _, err := util.WriteFile(filepath.Join(initSettings.Loc, "README.md"), fmt.Sprintf("# %s\n", initSettings.Name)); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(9)
			return
		}
		if err := exec.Command("git", "init").Run(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(10)
			return
		}
		if err := exec.Command("git", "add", ".").Run(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(11)
			return
		}
	}

	// open editor
	editor := open.Application_None
	switch initSettings.Template {
	case template.Template_Go:
		editor = open.Application_Go
	}
	if editor != open.Application_None {
		if err := open.Open(editor, initSettings.Loc); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(12)
			return
		}
	}
}
