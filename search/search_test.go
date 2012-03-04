package search

import (
	"github.com/afajl/assert"
	"github.com/afajl/ctrl/remote"
	"net/url"
	"testing"
)

type testSearcher struct {
	got_id  []string
	got_tag []string
}

func (t *testSearcher) Id(s ...string) ([]*remote.Host, error) {
	t.got_id = s
	return nil, nil
}

func (t *testSearcher) Tags(s ...string) ([]*remote.Host, error) {
	t.got_tag = s
	return nil, nil
}

func (t *testSearcher) String() string {
	return "testSearcher"
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
	var callorder [3]int
	var idx int

	facOne := func(u *url.URL) (Searcher, error) {
		callorder[idx] = 1
		idx++
		return &testSearcher{}, nil
	}
	facTwo := func(u *url.URL) (Searcher, error) {
		callorder[idx] = 2
		idx++
		return &testSearcher{}, nil
	}

	ms := NewMultiSearcher()
	ms.Register(facOne)
	ms.Register(facTwo)

	err := ms.AddSource("myscheme://path")
	if err != nil {
		t.Fatal("SetSources failed", err)
	}
	assert.Equal(t, callorder[0], 2, "the last added fac should be first")
	assert.Equal(t, callorder[1], 1, "the first added fac should be last")
}
