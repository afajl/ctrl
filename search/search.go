package search

import (
	"errors"
	"fmt"
	"github.com/afajl/ctrl/remote"
	"net/url"
)

var globalSearcher = NewMultiSearcher()

func Register(searcherFac SearcherFac) {
	globalSearcher.Register(searcherFac)
}

func AddSource(rawurls ...string) error {
	return globalSearcher.AddSource(rawurls...)
}

type SearcherFac func(*url.URL) (Searcher, error)

type Searcher interface {
	Name(...string) ([]*remote.Host, error)
	Tags(...string) ([]*remote.Host, error)
}

type MultiSearcher struct {
	searcherFacs []SearcherFac
	searchers    []Searcher
}

func NewMultiSearcher() *MultiSearcher {
	s := &MultiSearcher{searcherFacs: make([]SearcherFac, 0, 2),
		searchers: make([]Searcher, 0, 2)}
	return s
}

func (s *MultiSearcher) Register(searcherFac SearcherFac) {
	if searcherFac == nil {
		panic("search: Register SearcherFac is nil")
	}
	s.searcherFacs = append([]SearcherFac{searcherFac}, s.searcherFacs...)
}

func parseUrls(rawurls ...string) ([]*url.URL, error) {
	urls := make([]*url.URL, len(rawurls))
	for i, rawurl := range rawurls {
		u, err := url.Parse(rawurl)
		if u == nil || err != nil {
			return nil, fmt.Errorf("search: invalid url %s: %s", rawurl, err)
		}
		urls[i] = u
	}
	return urls, nil
}

func (s *MultiSearcher) urlSearcher(url *url.URL) (Searcher, error) {
	for _, searcherFac := range s.searcherFacs {
		searcher, err := searcherFac(url)
		if err != nil {
			return nil, err
		}
		if searcher != nil {
			return searcher, nil
		}
	}
	return nil, errors.New("search: no handler for: " + url.String())
}

func (s *MultiSearcher) AddSource(rawurls ...string) error {
	urls, err := parseUrls(rawurls...)
	if err != nil {
		return err
	}
	for _, u := range urls {
		searcher, err := s.urlSearcher(u)
		if err != nil {
			return err
		}
		s.searchers = append(s.searchers, searcher)
	}
	return nil
}

/*func (s *MultiSearcher) Name(name ...string) ([]*remote.Host, error) {*/

/*}*/

//func ByName(s string) ([]*remote.Host, error) {
//
//}
