package template

import (
	"fmt"

	"github.com/tsukinoko-kun/codehere/internal/settings"
)

const (
	Template_Go = uint8(iota)
)

func Create(s settings.Init) error {
	switch s.Template {
	case Template_Go:
		return create_Go(s)
	}
	return fmt.Errorf("no Template with code %q", s.Template)
}
