package util

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

func Exist(loc string) (bool, error) {
	if _, err := os.Stat(loc); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func WriteFile(loc string, content string) (bool, error) {
	if exist, err := Exist(loc); err != nil {
		return false, err
	} else if exist {
		var overwrite bool
		f := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Value(&overwrite).
				Title(fmt.Sprintf("Overwrite existing file %s ?", loc)),
		))
		if err := f.Run(); err != nil {
			return false, err
		}
		if !overwrite {
			return false, nil
		}
	}

	if f, err := os.Create(loc); err != nil {
		return false, err
	} else {
		if _, err := f.WriteString(content); err != nil {
			return false, err
		}
	}

	return true, nil
}

func Download(loc string, url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}

	var buf strings.Builder
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return false, err
	}

	return WriteFile(loc, buf.String())
}
