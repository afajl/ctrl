package search

import (
	"github.com/afajl/assert"
	"github.com/afajl/ctrl/remote"
	"net/url"
	"testing"
)

type testSearcher struct {
	got_byname string
	got_bytag  []string
}

func (t *testSearcher) ByName(s string) ([]*remote.Host, error) {
	t.got_byname = s
	return nil, nil
}

func (t *testSearcher) ByTags(s ...string) ([]*remote.Host, error) {
	t.got_bytag = s
	return nil, nil
}

func clearState() {
	searcherFacs = make([]SearcherFac, 0, 2)
}

func TestSetSourcesOk(t *testing.T) {
	defer clearState()

	var got_scheme string

	facAny := func(u *url.URL) (Searcher, error) {
		got_scheme = u.Scheme
		return &testSearcher{}, nil
	}
	Register(facAny)
	err := SetSources([]string{"myscheme://path"})
	if err != nil {
		t.Fatal("SetSources failed", err)
	}
	assert.Equal(t, got_scheme, "myscheme")
}

func TestSetSourcesNoHandler(t *testing.T) {
	defer clearState()
	facNone := func(u *url.URL) (Searcher, error) {
		return nil, nil
	}
	Register(facNone)
	err := SetSources([]string{"myscheme://path"})
	if err == nil {
		t.Fatalf("Setsources should fail when no handlers")
	}
}

func TestSetSourcesOrdering(t *testing.T) {
	defer clearState()
	var called = 1
	facOne := func(u *url.URL) (Searcher, error) {
		called *= 10
		return &testSearcher{}, nil
	}
	facTwo := func(u *url.URL) (Searcher, error) {
		called--
		return &testSearcher{}, nil
	}
	Register(facOne)
	Register(facTwo)

	err := SetSources([]string{"myscheme://path"})
	if err != nil {
		t.Fatal("SetSources failed", err)
	}

	assert.Equal(t, called, 0, "the last added fac should be first")
}
