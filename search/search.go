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
	Id(...string) ([]*remote.Host, error)
	Tags(...string) ([]*remote.Host, error)
	String() string
}

type MultiSearcher struct {
	searcherFacs []SearcherFac
	searchers    []Searcher
	sourcesAdded bool
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

func (s *MultiSearcher) getSearcher(url *url.URL) (Searcher, error) {
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
		searcher, err := s.getSearcher(u)
		if err != nil {
			return err
		}
		s.searchers = append(s.searchers, searcher)
	}
	s.sourcesAdded = true
	return nil
}

func verifyMatches(hosts []*remote.Host) error {
	var seen = make(map[string]bool, len(hosts))
	for _, host := range hosts {
		if host.Id == "" {
			return errors.New("search: found zero length id: " + host.Id)
		}
		if _, duplicate := seen[host.Id]; duplicate {
			return errors.New("search: found duplicate match: " + host.Id)
		}
		seen[host.Id] = true
	}
	return nil
}

func groupMatches(matches [][]*remote.Host) ([]*remote.Host, error) {
	var hostmap = make(map[string][]*remote.Host)
	for _, hosts := range matches {
		for _, host := range hosts {
			groups, seen := hostmap[host.Id]
			if !seen {
				hostmap[host.Id] = []*remote.Host{host}
			} else {
				hostmap[host.Id] = append(groups, host)
			}
		}
	}
	return foldHosts(hostmap)
}

func foldHosts(hostgroups map[string][]*remote.Host) ([]*remote.Host, error) {
	var hostlist = make([]*remote.Host, 0, len(hostgroups))
	for _, hosts := range hostgroups {
		host := remote.Fold(hosts...)
		hostlist = append(hostlist, host)
	}
	return hostlist, nil
}

func (s *MultiSearcher) Id(ids ...string) (hosts []*remote.Host, err error) {
	if !s.sourcesAdded {
		return nil, errors.New("search: add source before starting search")
	}

	var allmatches = make([][]*remote.Host, 0, len(s.searchers))
	for i, searcher := range s.searchers {
		if hosts, err = searcher.Id(ids...); err != nil {
			return
		}
		if err = verifyMatches(hosts); err != nil {
			return nil, fmt.Errorf("search: %s: %s", searcher, err)
		}
		allmatches[i] = hosts
	}
	return groupMatches(allmatches)
}

//func ByName(s string) ([]*remote.Host, error) {
//
//}
