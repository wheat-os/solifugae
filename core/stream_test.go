package core

import (
	"context"
	"testing"
)

// TestBaseStream 表驱动测试 BaseStream 的方法
func TestBaseStream(t *testing.T) {
	// 创建测试用例
	tests := []struct {
		name     string
		setup    func() *BaseStream
		testFunc func(t *testing.T, bs *BaseStream)
	}{
		{
			name: "Context test",
			setup: func() *BaseStream {
				ctx := context.WithValue(context.Background(), "key", "value")
				bs := &BaseStream{ctx: ctx}
				return bs
			},
			testFunc: func(t *testing.T, bs *BaseStream) {
				if bs.Context() == nil {
					t.Error("Expected context not to be nil")
				}
				if bs.Context().Value("key") != "value" {
					t.Error("Expected context value to be 'value'")
				}
			},
		},
		{
			name: "WithContext test",
			setup: func() *BaseStream {
				return &BaseStream{ctx: context.Background()}
			},
			testFunc: func(t *testing.T, bs *BaseStream) {
				newCtx := context.WithValue(context.Background(), "new-key", "new-value")
				bs.WithContext(newCtx)

				if bs.Context() != newCtx {
					t.Error("Expected context to be updated")
				}
			},
		},
		{
			name: "Name test",
			setup: func() *BaseStream {
				return &BaseStream{}
			},
			testFunc: func(t *testing.T, bs *BaseStream) {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected Name() to panic")
					}
				}()
				bs.Name()
			},
		},
		{
			name: "Codec test",
			setup: func() *BaseStream {
				return &BaseStream{}
			},
			testFunc: func(t *testing.T, bs *BaseStream) {
				codec, ok := bs.Codec()
				if !ok {
					t.Error("Expected Codec() to return true")
				}

				if codec == nil {
					t.Error("Expected Codec() to return a non-nil codec")
				}
			},
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := tt.setup()
			tt.testFunc(t, bs)
		})
	}
}
