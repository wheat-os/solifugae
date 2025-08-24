package core

import (
	"context"
)

type StreamCodec interface {
	// Encode encodes the stream data into bytes
	Encode(sd StreamData) ([]byte, error)

	// Decode decodes the stream data from bytes
	Decode(data []byte) (StreamData, error)
}

type StreamType uint8

const (
	// StreamTypeUnknown is the type for an unknown stream, used for error handling
	StreamTypeUnknown StreamType = iota
	// StreamTypeRequest is the type for a request stream
	StreamTypeRequest
	// StreamTypeResponse is the type for a response stream
	StreamTypeResponse
	// StreamTypeItem is the type for an item stream
	StreamTypeItem
)

type StreamData interface {
	// Context returns the context for the stream
	Context() context.Context
	// WithContext returns a new stream with the given context
	WithContext(ctx context.Context)

	// Type returns the type of the stream
	Type() StreamType

	// Name returns the name of the stream
	// name is used to identify the stream, and is unique within a project
	// can be used {package name}.{stream name} github.com/wheat-os/solifugae/core.baseSteam
	Name() string

	// 获取流的编码器
	// 返回false 表示该流不可被编码,在需要持久化的场景，该流会被忽略
	// Codec returns the codec for the stream
	// return false, the stream will be ignored in persistence
	Codec() (StreamCodec, bool)
}

type BaseStream struct {
	ctx context.Context
}

// Context returns the context for the stream
func (b *BaseStream) Context() context.Context {
	return b.ctx
}

// WithContext returns a new stream with the given context
func (b *BaseStream) WithContext(ctx context.Context) {
	b.ctx = ctx
}

// Type
func (b *BaseStream) Type() StreamType {
	panic("not implemented") // TODO: Implement
}

// Name returns the name of the stream
// name is used to identify the stream, and is unique within a project
// can be used {package name}.{stream name} github.com/wheat-os/solifugae/core.baseSteam
func (b *BaseStream) Name() string {
	panic("not implemented") // TODO: Implement
}

// 获取流的编码器
// 返回false 表示该流不可被编码,在需要持久化的场景，该流会被忽略
// Codec returns the codec for the stream
// return false, the stream will be ignored in persistence
func (b *BaseStream) Codec() (StreamCodec, bool) {
	return baseStreamCodec{}, true
}

type baseStreamCodec struct {
}

// Encode encodes the given data into the stream
func (b baseStreamCodec) Encode(sd StreamData) ([]byte, error) {
	return []byte{}, nil
	// buf := bytes.NewBuffer(nil)
	// if err := gob.NewEncoder(buf).Encode(sd); err != nil {
	// 	return nil, err
	// }
	// return buf.Bytes(), nil
}

// Decode decodes the given data from the stream
func (b baseStreamCodec) Decode(data []byte) (StreamData, error) {
	return &BaseStream{}, nil
	// var sd = &BaseStream{}
	// if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&sd); err != nil {
	// 	return nil, err
	// }
	// return sd, nil
}
