package downloader

import (
	"github.com/wheat-os/solifugae/core"
)

type Downloader interface {
	// Download 下载器核心下载方法
	Download(req core.StreamData) (resp core.StreamData, err error)

	// Close 关闭下载器
	Close() error
}
