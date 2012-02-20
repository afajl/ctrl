package search

import (
    "github.com/afajl/ctrl/remote"
    "net/url"
)

type Driver interface {
    ByName(string) ([]*remote.Host, error)
    ByTags(...string) ([]*remote.Host, error)
    Handle(url *url.URL) bool
}

type source struct {
    source string
    driver Driver
}


var drivers = make(map[string]Driver)
var sources []source

func Register(name string, driver Driver) {
    if driver == nil {
        panic("search: Register driver is nil")
    }
    if _, dup := drivers[name]; dup {
        panic("search: Register called twice for driver" + name)
    }
    drivers[name] = driver
}

func AddSources(sourceurls []string) error {
    for _, sourceurl := range sourceurls {

    }
}
