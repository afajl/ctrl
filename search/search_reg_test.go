package search

import (
	"github.com/afajl/assert"
	"github.com/afajl/ctrl/remote"
	"net/url"
	"testing"
	"fmt"
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


func searcherFacNever (u *url.URL) (Searcher, error) {
	return nil, nil
}

func searcherFacAlways (u *url.URL) (Searcher, error) {
	return &testSearcher{}, nil
}


func TestGlobalSearcher(t *testing.T) {
	defer func() {
		globalSearcher = NewMultiSearcher()
	}()

	Register(searcherFacAlways)
	err := AddSource("myscheme://path")
	if err != nil {
		t.Fatal("AddSource failed", err)
	}
}

func TestSetSourcesNoHandler(t *testing.T) {
	ms := NewMultiSearcher()
	ms.Register(searcherFacNever)
	err := ms.AddSource("myscheme://path")
	if err == nil {
		t.Fatalf("AddSource should fail when no handlers wants it")
	}
}

func TestSetSourcesBad(t *testing.T) {
	ms := NewMultiSearcher()
	ms.Register(searcherFacAlways)

	type test struct {
		u string
		ok bool
	}
	for _, ut := range []test{
		{"ok://path", true},
		{"", false},
	} {
		err := ms.AddSource(ut.u)
		if ut.ok {
			assert.Equal(t, err, nil, ut.u + " should not fail")
		}
	}
}


func buildSearcherFac(id int, idx *int, res *string, scheme string) SearcherFac {
	return func(u *url.URL) (Searcher, error) {
		*idx++
		*res = fmt.Sprintf("searcherFac %d was called %d", id, *idx)
		if u.Scheme == scheme {
			*res += " and matched"
			return &testSearcher{}, nil
		}
		return nil, nil
	}
}

func TestSetSourcesOrdering(t *testing.T) {
	var resOne, resTwo string
	var idx int

	ms := NewMultiSearcher()
	ms.Register(buildSearcherFac(1, &idx, &resOne, "http"))
	ms.Register(buildSearcherFac(2, &idx, &resTwo, "http"))

	err := ms.AddSource("http://path")
	if err != nil {
		t.Fatal("SetSources failed", err)
	}
	assert.Equal(t, resOne, "", "searcherFac1 should not be called if last added matches")
	assert.Equal(t, resTwo, "searcherFac 2 was called 1 and matched")
}

func TestSetSourcesAccept(t *testing.T) {
	var resOne, resTwo string
	var idx int

	ms := NewMultiSearcher()
	ms.Register(buildSearcherFac(1, &idx, &resOne, "http"))
	ms.Register(buildSearcherFac(2, &idx, &resTwo, "NO MATCH"))

	err := ms.AddSource("http://path")
	if err != nil {
		t.Fatal("SetSources failed", err)
	}
	assert.Equal(t, resOne, "searcherFac 1 was called 2 and matched")
	assert.Equal(t, resTwo, "searcherFac 2 was called 1")
}

