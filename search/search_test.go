package search

import (
	"github.com/afajl/assert"
	"github.com/afajl/ctrl/remote"
	"net/url"
	"testing"
	"errors"
)

type testSearcher struct {
	ret []*remote.Host
	ret_err error
}

func (t *testSearcher) Id(s ...string) ([]*remote.Host, error) {
	return t.ret, t.ret_err
}

func (t *testSearcher) Tags(s ...string) ([]*remote.Host, error) {
	return t.Id(s...)
}

func (t *testSearcher) String() string {
	return "testSearcher"
}

type rhl []*remote.Host

type searcherSetup struct {
	scheme string
	ret rhl
	ret_err error
}

func buildSearcher(s searcherSetup) SearcherFac {
	ts := &testSearcher{ret: s.ret, ret_err: s.ret_err}
	return func(u *url.URL) (Searcher, error) {
		return ts, nil
	}
}

func buildMultiSearcher(setup ...searcherSetup) *MultiSearcher {
	ms := NewMultiSearcher()
	for _, s := range setup {
		ms.Register(s.scheme, buildSearcher(s))
	}
	return ms
}

func TestRegister(t *testing.T) {
	dummyFac := func(u *url.URL) (Searcher, error) {
		return nil, nil
	}
	type regtest struct {
		scheme string
		f SearcherFac
		err string
	}
	tests := []regtest{
		{"", dummyFac, "search: cannot register empty scheme"},
		{"ok", nil, "search: SearcherFac cannot be nil"},
	}
	for _, test := range tests {
		ms := NewMultiSearcher()
		assert.Panic(t, test.err, func() { ms.Register(test.scheme, test.f) })
	}

	ms := NewMultiSearcher()
	ms.Register("duplicate", dummyFac)
	assert.Panic(t, "search: duplicate scheme handlers found", 
	             func() { ms.Register("duplicate", dummyFac) })
}

func TestAddUrl (t *testing.T) {
	// empty url
	ms := NewMultiSearcher()
	if err := ms.AddUrl(""); err == nil {
		t.Fatal("empty url should return an error")
	}

	if err := ms.AddUrl("notregistered://p"); err == nil {
		t.Fatal("no handlers should fail")
	}
	ms = buildMultiSearcher(searcherSetup{"match", nil, nil})
	if err := ms.AddUrl("notregistered://p"); err == nil {
		t.Fatal("no matching handlers should fail")
	}

	if err := ms.AddUrl("match://p"); err != nil {
		t.Fatal("matching handler should be ok")
	}
}

func TestBadSearchers(t *testing.T) {
	setup := []searcherSetup{
		{"any", rhl{nil},  nil},  // nil host
		{"any", nil,  nil},       // nil hosts
		{"any", rhl{&remote.Host{Id: ""}}, nil}, // empty id
		{"any", rhl{&remote.Host{Id: "d"}, &remote.Host{Id: "d"}}, nil}, // duplicate
		{"any", rhl{},  errors.New("err")}, // error
	}

	for i, bs := range setup {
		ms := buildMultiSearcher(bs)
		ms.AddUrl("any://p")
		_, err := ms.Id("_")
		if err == nil {
			t.Errorf("test %d should fail", i)
		}
	}
}

func TestSearchUrlOrder(t *testing.T) {
	setup := []searcherSetup{
		{"ascheme", rhl{&remote.Host{Id: "x", Name: "a"}},  nil},
		{"bscheme", rhl{&remote.Host{Id: "x", Name: "b"}},  nil},
	}

	ms := buildMultiSearcher(setup...)
	ms.AddUrl("ascheme://p", "bscheme://p")

	res, err := ms.Id("_")
    if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].Name, "b", "the last added source should expand")
}


func TestSearchConcatenationHosts(t *testing.T) {
	setup := []searcherSetup{
		{"ascheme", rhl{&remote.Host{Id: "x", Name: "xa"},
					    &remote.Host{Id: "y", Name: "ya"}},  nil},
		{"bscheme", rhl{&remote.Host{Id: "y", Name: "yb"}},  nil},

		{"cscheme", rhl{&remote.Host{Id: "x", Name: "xc"},
						&remote.Host{Id: "y", Name: "yc"}},  nil},
	}

	ms := buildMultiSearcher(setup...)
	ms.AddUrl("ascheme://p", "bscheme://p")

	res, err := ms.Tags("_")
    if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(res), 2)
	for _, r := range res {
		switch r.Id {
		case "x":
			assert.Equal(t, r.Name, "xa") 
		case "y":
			assert.Equal(t, r.Name, "yb")
		default:
			t.Fatal("unkown host", r.Id)
		}
	}

	ms.AddUrl("cscheme://p")
	res, err = ms.Tags("_")
    if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(res), 2)
	for _, r := range res {
		switch r.Id {
		case "x":
			assert.Equal(t, r.Name, "xc") 
		case "y":
			assert.Equal(t, r.Name, "yc")
		default:
			t.Fatal("unkown host", r.Id)
		}
	}
}
