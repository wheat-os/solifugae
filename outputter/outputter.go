package outputter

import (
	"github.com/wheat-os/solifugae/core"
)

type Outputter interface {
	Save(item core.StreamData) error
}
