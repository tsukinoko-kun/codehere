package open

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Application uint8

const (
	Application_None Application = iota
	Application_Go               = iota
)

var (
	ApplicationNotFoundError = errors.New("application not found")
)

func Open(a Application, wd string) error {
	var err error

	switch a {
	case Application_Go:
		err = open_GoLand(wd)
	default:
		return errors.Join(ApplicationNotFoundError, fmt.Errorf("no application with code %q", a))
	}
	if err == nil {
		return nil
	}

	if visualEditor, ok := os.LookupEnv("VISUAL"); ok {
		if err := open(visualEditor, wd); err == nil {
			return nil
		}
	}

	if editor, ok := os.LookupEnv("EDITOR"); ok {
		if err := open(editor, wd); err == nil {
			return nil
		}
	}

	return err
}

func open(a string, wd string) error {
	p, err := exec.LookPath(a)
	if err != nil {
		return errors.Join(ApplicationNotFoundError, err)
	}
	cmd := exec.Command(p, wd)
	cmd.Dir = filepath.Dir(wd)
	return cmd.Start()
}

func open_GoLand(wd string) error {
	return open("goland", wd)
}
