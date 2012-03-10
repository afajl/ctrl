package search

import (
	"errors"
	"fmt"
	"github.com/afajl/ctrl/host"
	"net/url"
)

// Factory that should return a Searcher set up
// for the url or nil if it doesnt want to handle
// this url
type SearcherFac func(*url.URL) (Searcher, error)

type Searcher interface {
	Id(...string) ([]*host.Host, error)
	Tags(...string) ([]*host.Host, error)
	String() string
}

// Global searcher
var globalSearcher = NewMultiSearcher()

func Register(scheme string, searcherFac SearcherFac) {
	globalSearcher.Register(scheme, searcherFac)
}

func AddUrl(rawurls ...string) error {
	return globalSearcher.AddUrl(rawurls...)
}




type MultiSearcher struct {
	searcherFacs map[string]SearcherFac
	searchers    []Searcher
}

func NewMultiSearcher() *MultiSearcher {
	s := &MultiSearcher{
		searcherFacs: map[string]SearcherFac{},
		searchers: []Searcher{},
	}
	return s
}

// Append a new searcher factory. Later added searchers
// will be tried before previous ones
func (s *MultiSearcher) Register(scheme string, searcherFac SearcherFac) {
	if scheme == "" {
		panic("search: cannot register empty scheme")
	}
	if searcherFac == nil {
		panic("search: SearcherFac cannot be nil")
	}
	if _, present := s.searcherFacs[scheme]; present {
		panic("search: duplicate scheme handlers found")
	}
	s.searcherFacs[scheme] = searcherFac
}


func (s *MultiSearcher) AddUrl(rawurls ...string) error {
	for _, rawurl := range rawurls {
		u, err := url.Parse(rawurl)
		if err != nil || u == nil {
			return fmt.Errorf("search: invalid url %s: %s", rawurl, err)
		}

		searcherFac, present := s.searcherFacs[u.Scheme]
		if !present {
			return errors.New("search: no handler found for scheme: " + u.Scheme)
		}

		searcher, err := searcherFac(u)
		if err != nil {
			return err
		}
		s.searchers = append(s.searchers, searcher)
	}
	return nil
}


// Return hosts matching Ids. Order is undefined.
func (s *MultiSearcher) Id(ids ...string) (hosts []*host.Host, err error) {
    fn := func(searcher Searcher) ([]*host.Host, error) {
		return searcher.Id(ids...)
	}
	return searchTempl(s.searchers, fn)

}

// Return hosts that has all tags. Order is undefined.
func (s *MultiSearcher) Tags(tags ...string) (hosts []*host.Host, err error) {
    fn := func(searcher Searcher) ([]*host.Host, error) {
		return searcher.Tags(tags...)
	}
	return searchTempl(s.searchers, fn)
}


type searcherDelegate func(Searcher) ([]*host.Host, error)

func searchTempl(searchers []Searcher, fn searcherDelegate) (hosts []*host.Host, err error) {
	if len(searchers) == 0 {
		return nil, errors.New("no sources added")
	}
	var idHost = map[string]*host.Host{}

	for _, searcher := range searchers {
		if hosts, err = fn(searcher); err != nil {
			return
		}
		if err := checkHosts(hosts); err != nil {
			return nil, fmt.Errorf("%s: %s", searcher, err)
		}
		for _, h := range hosts {
			if v, present := idHost[h.Id]; present {
				idHost[h.Id] = host.Update(v, h)
			} else {
				idHost[h.Id] = h
			}
		}
	}
    var res = make([]*host.Host, 0, len(idHost))
	for _, v := range idHost {
		res = append(res, v)
	}
	return res, nil
}

func checkHosts(hosts []*host.Host) error {
	if hosts == nil {
		return errors.New("hosts nil")
	}
	var seen = map[string]int{}
	for i, h := range hosts {
		if h == nil {
			return fmt.Errorf("nil host at pos %d", i)
		}
		if h.Id == "" {
			return fmt.Errorf("empty Id for host at pos %d", i)
		}
		if v, duplicate := seen[h.Id]; duplicate {
			return fmt.Errorf("duplicate Id '%s' at pos: %d, %d", h.Id, v, i)
		}
		seen[h.Id] = i
	}
	return nil
}
