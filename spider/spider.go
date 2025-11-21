package spider

import (
	"iter"

	"github.com/wheat-os/solifugae/core"
)

// ParseFunc 定义 spider parse 函数
type ParseFunc func(resp core.StreamData) (iter.Seq[core.StreamData], error)

type Spider interface {
	// StartRequest start request
	StartRequest() (iter.Seq[core.StreamData], error)
	// Parse parse
	Parse(resp core.StreamData) (iter.Seq[core.StreamData], error)
	// name 返回爬虫的唯一名称
	Name() string
}
