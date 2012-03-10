package search

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"github.com/afajl/ctrl/config"
	"github.com/afajl/ctrl/remote"
	"io"
	"bytes"
	"encoding/json"
)

func init() {
	Register("json", JsonFac)
}

type JsonSearcher struct {
	id string
	file io.Reader
	hosts map[string]*remote.Host
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
		return err
	}

	if _, emptyid := js.hosts[""]; emptyid {
		return fmt.Errorf("json contains empty id")
	}
	return nil
}

func (js *JsonSearcher) Id(s ...string) ([]*remote.Host, error) {
	if err := js.parse(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (js *JsonSearcher) Tags(s ...string) ([]*remote.Host, error) {
	return nil, nil
}

func (js *JsonSearcher) String() string {
	return "testSearcher"
}
