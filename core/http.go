package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
)

var _ StreamData = (*HttpStream)(nil)

type HttpStream struct {
	*BaseStream
	HttpRequest *http.Request
}

// Name returns the name of the stream
// name is used to identify the stream, and is unique within a project
// can be used {package name}.{stream name} github.com/wheat-os/solifugae/core.baseSteam
func (h *HttpStream) Name() string {
	return "github.com/wheat-os/solifugae/core.HttpStream"
}

func (h *HttpStream) Type() StreamType {
	return StreamTypeRequest
}

// 获取流的编码器
// 返回false 表示该流不可被编码,在需要持久化的场景，该流会被忽略
// Codec returns the codec for the stream
// return false, the stream will be ignored in persistence
func (h *HttpStream) Codec() (StreamCodec, bool) {
	return httpStreamCodec{}, true
}

// NewHttpStream returns a new HttpStream
// func NewHttpStream(method, url string, body io.Reader) (*HttpStream, error) {
// 	http.NewRequest(method, url, body)
// }

type httpStreamCodec struct {
}

// serializableHttpStream 是可序列化的请求结构
type serializableHttpStream struct {
	Method         string
	URL            string
	Header         http.Header
	Body           []byte
	BaseStreamBody []byte
}

// Encode encodes the given data into the stream
func (b httpStreamCodec) Encode(sd StreamData) ([]byte, error) {
	hs, ok := sd.(*HttpStream)
	if !ok {
		return nil, fmt.Errorf("the data is not a HttpStream")
	}
	bsCodec, ok := hs.BaseStream.Codec()
	if !ok {
		return nil, fmt.Errorf("the base stream is not serializable")
	}
	bsBody, err := bsCodec.Encode(hs.BaseStream)
	if err != nil {
		return nil, fmt.Errorf("failed to encode base stream: %v", err)
	}

	// 读取请求体
	var bodyBytes []byte
	if hs.HttpRequest.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(hs.HttpRequest.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %v", err)
		}
		// 恢复原始请求体
		hs.HttpRequest.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	var sr = serializableHttpStream{
		Method:         hs.HttpRequest.Method,
		URL:            hs.HttpRequest.URL.String(),
		Header:         hs.HttpRequest.Header,
		Body:           bodyBytes,
		BaseStreamBody: bsBody,
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(&sr); err != nil {
		return nil, fmt.Errorf("failed to encode http stream: %v", err)
	}
	return buffer.Bytes(), nil
}

// Decode decodes the given data from the stream
func (b httpStreamCodec) Decode(data []byte) (StreamData, error) {
	var serializable serializableHttpStream

	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&serializable); err != nil {
		return nil, fmt.Errorf("failed to decode http stream: %v", err)
	}
	req, err := http.NewRequest(serializable.Method, serializable.Method, bytes.NewReader(serializable.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header = serializable.Header

	bs := &BaseStream{}
	bsCodec, ok := bs.Codec()
	if !ok {
		return nil, fmt.Errorf("failed to get base stream codec")
	}
	bsData, err := bsCodec.Decode(serializable.BaseStreamBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base stream body: %v", err)
	}
	bs, ok = bsData.(*BaseStream)
	if !ok {
		return nil, fmt.Errorf("failed to cast base stream data")
	}
	return &HttpStream{BaseStream: bs, HttpRequest: req}, nil
}
