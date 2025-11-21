package spider

import (
	"errors"
	"fmt"
	"iter"
	"reflect"

	"github.com/wheat-os/solifugae/core"
	"github.com/wheat-os/solifugae/xiter"
)

type SpiderRefCenter struct {
	spiders        map[string]Spider
	spiderCallFunc map[string]ParseFunc
}

// RegisterSpider 添加爬虫到管理器中心
func (s *SpiderRefCenter) RegisterSpider(sp Spider) {
	if _, ok := s.spiders[sp.Name()]; ok {
		panic(fmt.Sprintf("spider %s already exists", sp.Name()))
	}
	s.spiders[sp.Name()] = sp
	t := reflect.TypeOf(sp)
	v := reflect.ValueOf(sp)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		methodName := method.Name
		if s.isValidPaseFuncWrapperSignature(method.Type) {
			// 创建闭包，只在收集时使用一次反射
			methodValue := v.Method(i)
			funcName := fmt.Sprintf("%s.%s", sp.Name(), methodName)
			s.spiderCallFunc[funcName] = s.createPaseFuncWrapper(sp, methodValue)
		}
	}
}

// createPaseFuncWrapper 创建解析函数包装器
func (s *SpiderRefCenter) createPaseFuncWrapper(spider Spider, methodValue reflect.Value) ParseFunc {
	return func(resp core.StreamData) (iter.Seq[core.StreamData], error) {
		args := []reflect.Value{reflect.ValueOf(resp)}
		results := methodValue.Call(args)
		result := results[0].Interface().(iter.Seq[core.StreamData])
		var err error
		if results[1].IsNil() {
			err = nil
		} else {
			err = results[1].Interface().(error)
		}
		if err != nil {
			return nil, err
		}
		resultIterator := xiter.Before(result, func(stream core.StreamData) {
			if stream == nil {
				return
			}
			stream.SetMetadata(SPIDER_CONTEXT_NAME, spider.Name()) // 设置spider名称
		})
		return resultIterator, nil
	}
}

func (s *SpiderRefCenter) isValidPaseFuncWrapperSignature(methodType reflect.Type) bool {
	if methodType.NumIn() != 2 || methodType.NumOut() != 2 {
		return false
	}
	if methodType.In(1) != reflect.TypeOf((*core.StreamData)(nil)).Elem() {
		return false
	}
	if methodType.Out(0) != reflect.TypeOf((*iter.Seq[core.StreamData])(nil)).Elem() {
		return false
	}
	if methodType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return false
	}
	return true
}

func (s *SpiderRefCenter) Parse(resp core.StreamData) (iter.Seq[core.StreamData], error) {
	spiderName, ok := resp.GetMetadata(SPIDER_CONTEXT_NAME)
	if !ok {
		return nil, errors.New("spiderName not found")
	}
	callbackName := "Parse"
	callback, ok := resp.GetMetadata(SPIDER_CONTEXT_CALLBACK)
	if ok {
		callbackName = callback.(string)
	}

	funcName := fmt.Sprintf("%s.%s", spiderName, callbackName)
	parseFunc, ok := s.spiderCallFunc[funcName]
	if !ok {
		return nil, errors.New("parse func not found")
	}
	return parseFunc(resp)
}

// NewSpiderRefCenter 创建爬虫管理中心
func NewSpiderRefCenter() *SpiderRefCenter {
	return &SpiderRefCenter{
		spiders:        make(map[string]Spider, 0),
		spiderCallFunc: make(map[string]ParseFunc),
	}
}
