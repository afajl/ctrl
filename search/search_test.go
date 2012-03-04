package search

import (
	"github.com/afajl/assert"
	"github.com/afajl/ctrl/remote"
	"net/url"
	"testing"
)

type testSearcher struct {
	got_byname []string
	got_bytag  []string
}

func (t *testSearcher) Name(s ...string) ([]*remote.Host, error) {
	t.got_byname = s
	return nil, nil
}

func (t *testSearcher) Tags(s ...string) ([]*remote.Host, error) {
	t.got_bytag = s
	return nil, nil
}

func TestGlobalSearcher(t *testing.T) {
	defer func() {
		globalSearcher = NewMultiSearcher()
	}()

	var got_scheme string

	facAny := func(u *url.URL) (Searcher, error) {
		got_scheme = u.Scheme
		return &testSearcher{}, nil
	}
	Register(facAny)
	err := AddSource("myscheme://path")
	if err != nil {
		t.Fatal("AddSource failed", err)
	}
	assert.Equal(t, got_scheme, "myscheme")
}

func TestSetSourcesNoHandler(t *testing.T) {
	facNone := func(u *url.URL) (Searcher, error) {
		return nil, nil
	}
	ms := NewMultiSearcher()
	ms.Register(facNone)
	err := ms.AddSource("myscheme://path")
	if err == nil {
		t.Fatalf("AddSource should fail when no handlers")
	}
}

func TestSetSourcesOrdering(t *testing.T) {
	var called = 1
	facOne := func(u *url.URL) (Searcher, error) {
		called *= 10
		return &testSearcher{}, nil
	}
	facTwo := func(u *url.URL) (Searcher, error) {
		called--
		return &testSearcher{}, nil
	}
	ms := NewMultiSearcher()
	ms.Register(facOne)
	ms.Register(facTwo)

	err := ms.AddSource("myscheme://path")
	if err != nil {
		t.Fatal("SetSources failed", err)
	}

	assert.Equal(t, called, 0, "the last added fac should be first")
}
