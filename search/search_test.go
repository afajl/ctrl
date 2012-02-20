package search

import (
    "github.com/afajl/ctrl/remote"
    "net/url"
	"testing"
)

type testdriver struct{
    got_byname string
    got_bytag []string
    handles string
}

func (t *testdriver) ByName(s string) ([]*remote.Host, error) {
    t.got_byname = s
    return nil, nil
}

func (t *testdriver) ByTag(s ...string) ([]*remote.Host, error) {
    t.got_bytag = s
    return nil, nil
}

func (t *testdriver) Handle(url *url.URL) bool {
    if url.Scheme == t.handles {
        return true
    }
    return false
}


func TestFuncDriver(t *testing.T) {
    td := &testdriver{handles: "json"}
    Register("json", td)
}
