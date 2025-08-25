package downloader

import (
	"context"

	"github.com/wheat-os/solifugae/core"
)

type Downloader interface {
	// Download 下载器核心下载方法
	Download(req core.StreamData) (core.StreamData, error)
	// 下载器启动方法
	Start(ctx context.Context) error
}


