package scheduler

import (
	"github.com/wheat-os/solifugae/core"
)

type Scheduler interface {
	SetStream(req core.StreamData) error

	GetStream() (core.StreamData, error)

	HasPending() bool

	Close() error
}
