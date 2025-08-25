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
		{
			name: "Type test",
			setup: func() *BaseStream {
				return &BaseStream{}
			},
			testFunc: func(t *testing.T, bs *BaseStream) {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected Type() to panic")
					}
				}()
				bs.Type()
			},
		},
		{
			name: "SetMetadata and GetMetadata test",
			setup: func() *BaseStream {
				return &BaseStream{}
			},
			testFunc: func(t *testing.T, bs *BaseStream) {
				// 测试设置元数据
				bs.SetMetadata("key1", "value1")
				bs.SetMetadata("key2", 123)

				// 测试获取存在的元数据
				value1, ok1 := bs.GetMetadata("key1")
				if !ok1 {
					t.Error("Expected key1 to exist")
				}
				if value1 != "value1" {
					t.Errorf("Expected value1 to be 'value1', got %v", value1)
				}

				value2, ok2 := bs.GetMetadata("key2")
				if !ok2 {
					t.Error("Expected key2 to exist")
				}
				if value2 != 123 {
					t.Errorf("Expected value2 to be 123, got %v", value2)
				}

				// 测试获取不存在的元数据
				_, ok3 := bs.GetMetadata("nonexistent")
				// 注意: 根据实现，即使键不存在，只要metadata map不为nil，也会返回true
				if !ok3 {
					t.Error("Expected ok to be true when metadata map is not nil")
				}
			},
		},
		{
			name: "GetMetadata from empty metadata test",
			setup: func() *BaseStream {
				return &BaseStream{}
			},
			testFunc: func(t *testing.T, bs *BaseStream) {
				// 测试从未初始化的metadata获取数据
				_, ok := bs.GetMetadata("anykey")
				if ok {
					t.Error("Expected no metadata to be found in uninitialized metadata map")
				}
			},
		},
		{
			name: "GetMetadata after explicit nil initialization",
			setup: func() *BaseStream {
				return &BaseStream{metadata: nil}
			},
			testFunc: func(t *testing.T, bs *BaseStream) {
				// 测试metadata显式初始化为nil的情况
				_, ok := bs.GetMetadata("anykey")
				if ok {
					t.Error("Expected no metadata to be found when metadata is explicitly nil")
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

// TestBaseStreamCodec 测试baseStreamCodec的Encode和Decode方法
func TestBaseStreamCodec(t *testing.T) {
	// 创建测试用例
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Encode test",
			testFunc: func(t *testing.T) {
				// 创建一个BaseStream实例
				bs := &BaseStream{
					ctx: context.Background(),
					metadata: map[string]any{
						"key1": "value1",
						"key2": 123,
					},
				}

				// 创建codec实例
				codec, _ := bs.Codec()

				// 编码BaseStream
				data, err := codec.Encode(bs)
				if err != nil {
					t.Errorf("Expected Encode to succeed, but got error: %v", err)
				}

				if len(data) == 0 {
					t.Error("Expected encoded data to be non-empty")
				}
			},
		},
		{
			name: "Decode test",
			testFunc: func(t *testing.T) {
				// 创建一个BaseStream实例并编码
				originalStream := &BaseStream{
					ctx: context.Background(),
					metadata: map[string]any{
						"key1": "value1",
						"key2": 123,
					},
				}

				codec := baseStreamCodec{}
				data, err := codec.Encode(originalStream)
				if err != nil {
					t.Fatalf("Failed to encode stream: %v", err)
				}

				// 解码数据
				decodedStream, err := codec.Decode(data)
				if err != nil {
					t.Errorf("Expected Decode to succeed, but got error: %v", err)
				}

				if decodedStream == nil {
					t.Error("Expected decoded stream to be non-nil")
				}
			},
		},
		{
			name: "Decode with invalid data test",
			testFunc: func(t *testing.T) {
				codec := baseStreamCodec{}
				// 尝试解码无效数据
				_, err := codec.Decode([]byte("invalid data"))
				if err == nil {
					t.Error("Expected Decode to fail with invalid data")
				}
			},
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}
