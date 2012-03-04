package search

import (
	"errors"
	"github.com/afajl/ctrl/remote"
	"net/url"
)

type SearcherFac func(*url.URL) (Searcher, error)

type Searcher interface {
	ByName(string) ([]*remote.Host, error)
	ByTags(...string) ([]*remote.Host, error)
}

var searcherFacs = make([]SearcherFac, 0, 2)
var searchers []Searcher

// Prepends the searcher factory to the list
func Register(searcherFac SearcherFac) {
	if searcherFac == nil {
		panic("search: Register SearcherFac is nil")
	}
	searcherFacs = append([]SearcherFac{searcherFac}, searcherFacs...)
}

func SetSources(strings []string) error {
	searchers = make([]Searcher, 0, len(strings))
	for _, s := range strings {
		sourceUrl, err := url.Parse(s)
		if sourceUrl == nil {
			return errors.New("search: invalid url " + s)
		}
		if err != nil {
			return err
		}
		handled := false
		for _, searcherFac := range searcherFacs {
			searcher, err := searcherFac(sourceUrl)
			if err != nil {
				return err
			}
			if searcher != nil {
				searchers = append(searchers, searcher)
				handled = true
			}
		}
		if !handled {
			return errors.New("search: could not find handler for " + s)
		}
	}
	return nil
}

//func ByName(s string) ([]*remote.Host, error) {
//
//}
