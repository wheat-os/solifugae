package outputter

import (
	"github.com/wheat-os/solifugae/core"
)

type Outputter interface {
	// Output 调用输出器
	Output(item core.StreamData) error
	// Close 关闭输出器
	Close() error
}
