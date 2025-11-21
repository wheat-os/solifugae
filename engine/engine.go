package engine

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/wheat-os/solifugae/core"
	"github.com/wheat-os/solifugae/downloader"
	"github.com/wheat-os/solifugae/outputter"
	"github.com/wheat-os/solifugae/scheduler"
	"github.com/wheat-os/solifugae/spider"
)

type FastEngine struct {
	scheduler    scheduler.Scheduler
	downloader   downloader.Downloader
	outputter    outputter.Outputter
	spiderCenter *spider.SpiderRefCenter

	activeWorker int32 // 激活的worker数量
	workerCount  int32 // worker总数
	workerGroup  sync.WaitGroup

	cancel func()
	closed int32
}

// 启动引擎，指定worker数量
func (f *FastEngine) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	f.cancel = cancel
	// 捕获系统信号
	f.captureSignals()

	// 调用开始爬虫函数

	// 启动worker
	for i := 0; i < int(f.workerCount); i++ {
		f.workerGroup.Add(1)
		go func() {
			defer f.workerGroup.Done()
			f.worker(ctx)
		}()
	}
	f.workerGroup.Wait()
	f.Close()
	// 回收其他资源
	f.downloader.Close()
	f.outputter.Close()
	f.scheduler.Close()
}

func (f *FastEngine) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			atomic.AddInt32(&f.activeWorker, 1)
			req, err := f.scheduler.GetStream()
			if err != nil {
				continue
			}
			f.handleRequestStream(req)
			atomic.AddInt32(&f.activeWorker, -1)
		}
		// 满足 close 条件, 退出
		if !f.scheduler.HasPending() && atomic.LoadInt32(&f.activeWorker) == 0 {
			return
		}
	}
}

// handleScheduleStream 处理来自调度器的 stream 流
func (f *FastEngine) handleRequestStream(stream core.StreamData) {
	respStream, err := f.downloader.Download(stream)
	if err != nil {
		return
	}
	f.handleResponseStream(respStream)
}

// handleResponseStream 处理来自下载器的 stream 流
func (f *FastEngine) handleResponseStream(stream core.StreamData) {
	parseStreams, err := f.spiderCenter.Parse(stream)
	if err != nil {
		return
	}
	for sd := range parseStreams {
		switch sd.Type() {
		case core.StreamTypeRequest:
			f.scheduler.SetStream(sd)
		case core.StreamTypeResponse:
			f.handleResponseStream(sd)
		case core.StreamTypeItem:
			f.outputter.Output(sd)
		default:
		}
	}
}

func (f *FastEngine) Close() {
	// 确保只关闭一次
	if !atomic.CompareAndSwapInt32(&f.closed, 0, 1) {
		return
	}
	f.cancel()
}

// captureSignals 捕获系统信号
func (f *FastEngine) captureSignals() {
	// 创建信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动信号处理协程
	go func() {
		<-sigChan
		// 收到退出信号，关闭引擎
		f.Close()
	}()
}
