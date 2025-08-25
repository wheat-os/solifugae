package spider

import (
	"context"
	"iter"

	"github.com/wheat-os/solifugae/core"
)

// Parse 定义 spider parse 函数
type Parse func(ctx context.Context, resp core.StreamData) (iter.Seq[core.StreamData], error)

type Spider interface {
	// StartRequest start request
	StartRequest() (iter.Seq[core.StreamData], error)

	// Parse parse
	Parse(resp core.StreamData) (iter.Seq[core.StreamData], error)
}
