package search

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"github.com/afajl/ctrl/config"
	"github.com/afajl/ctrl/host"
	"io"
	"bytes"
	"encoding/json"
	"sort"
)

func init() {
	Register("json", JsonFac)
}

type JsonSearcher struct {
	id string
	file io.Reader
	hosts map[string]*host.Host
	tags map[string][]*host.Host
}


func JsonFac(url *url.URL) (Searcher, error) {
	if url.Scheme != "json" {
		return nil, fmt.Errorf("json searcher got unknown scheme: %s", url)
	}
	var path string
	if url.Host != "" { // relative path
		rootdir := config.StartConfig.Rootdir
		if rootdir == "" {
			// TODO any usecase for this?
			rootdir, _ = os.Getwd()
		}
		path = filepath.Join(rootdir, url.Host, url.Path)
	} else { // absolute path
		path = url.Path
	}
    path = filepath.Clean(path)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if st, err := file.Stat(); err != nil {
		return nil, err
	} else if st.IsDir() {
		return nil, fmt.Errorf("json cannot parse dir: %s", url)
	}

	return &JsonSearcher{id: url.String(), file: file}, nil
}


func JsonFromReader(r io.Reader) (Searcher, error) {
	return &JsonSearcher{id: "<io.Reader>", file: r}, nil
}


func (js *JsonSearcher) parse() error {
	if js.hosts != nil {
		return nil
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(js.file); err != nil {
		return err
	}
	if err := json.Unmarshal(buf.Bytes(), &js.hosts); err != nil {
		return fmt.Errorf("json error: %s", err)
	}

	for id, h := range js.hosts {
		if id == "" {
			return fmt.Errorf("json contains empty id")
		}
		// sort tags
		sort.Strings(h.Tags)

		h.Id = id
	}
	return nil
}


func (js *JsonSearcher) Id(ids ...string) (hosts []*host.Host, err error) {
	if err = js.parse(); err != nil {
		return
	}
	for _, id := range ids {
		host, present := js.hosts[id]
		if !present {
			continue
		}
		hosts = append(hosts, host)
	}
	return
}

func (js *JsonSearcher) Tags(tags ...string) (hosts []*host.Host, err error) {
	if err = js.parse(); err != nil {
		return
	}
	NEXT_HOST: for _, h := range js.hosts {
		for _, tag := range tags {
			n := sort.SearchStrings(h.Tags, tag)
			if n == len(h.Tags) || h.Tags[n] != tag {
				continue NEXT_HOST
			}
		}
		hosts = append(hosts, h)
	}
	return

}

func (js *JsonSearcher) String() string {
	return js.id
}
