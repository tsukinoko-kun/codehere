package license

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/tsukinoko-kun/codehere/internal/settings"
	"github.com/tsukinoko-kun/codehere/internal/util"
)

const (
	License_MIT = uint8(iota)
)

const (
	mit = `Copyright %s %s

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`
)

func Create(s settings.Init) error {
	switch s.License {
	case License_MIT:
		return create_MIT(s)
	}
	return fmt.Errorf("no License with code %q", s.License)
}

func create_MIT(s settings.Init) error {
	var copyrightHolder string
	year := strconv.Itoa(time.Now().Year())

	if gitDisplayname, err := util.GitDisplayname(); err == nil {
		copyrightHolder = gitDisplayname
	}

	f := huh.NewForm(huh.NewGroup(
		huh.NewInput().Value(&copyrightHolder).Title("Copyright Holder"),
		huh.NewInput().Value(&year).Title("Year").Suggestions([]string{year}),
	))

	if err := f.Run(); err != nil {
		return err
	}

	if _, err := util.WriteFile(filepath.Join(s.Loc, "LICENSE"), fmt.Sprintf(mit, year, copyrightHolder)); err != nil {
		return err
	}

	return nil
}
