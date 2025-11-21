package spider

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wheat-os/solifugae/core"
)

type testSpiderV1 struct {
}

func (t *testSpiderV1) StartRequest() (iter.Seq[core.StreamData], error) {
	return nil, nil
}

func (t *testSpiderV1) Parse(resp core.StreamData) (iter.Seq[core.StreamData], error) {
	n := 10
	return func(yield func(core.StreamData) bool) {
		b := core.BaseStream{}
		for i := 0; i < n; i++ {
			b.SetMetadata("n", i)
			if !yield(&b) {
				break
			}
		}
	}, nil
}

func (t *testSpiderV1) ParseNil(resp core.StreamData) (iter.Seq[core.StreamData], error) {
	return func(yield func(core.StreamData) bool) {
		if !yield(nil) {
			return
		}
	}, nil
}

func (t *testSpiderV1) Name() string {
	return "testSpiderV1"
}

func TestSpiderRefCenter_AddSpider(t *testing.T) {
	spiderCenter := NewSpiderRefCenter()
	spiderCenter.RegisterSpider(&testSpiderV1{})
	resp, err := spiderCenter.spiderCallFunc["testSpiderV1.Parse"](&core.BaseStream{})
	require.NoError(t, err)
	n := 0
	for b := range resp {
		bn, ok := b.GetMetadata("n")
		require.True(t, ok)
		require.Equal(t, bn, n)
		n++
	}
	require.NoError(t, err)
	resp, err = spiderCenter.spiderCallFunc["testSpiderV1.ParseNil"](&core.BaseStream{})
	require.NoError(t, err)
	for b := range resp {
		require.Nil(t, b)
	}
}
