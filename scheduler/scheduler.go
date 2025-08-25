package scheduler

import (
	"github.com/wheat-os/solifugae/core"
)

type Scheduler interface {
	SetStream(req *core.StreamData) error

	GetStream(req *core.StreamData) (core.StreamData, error)
}
