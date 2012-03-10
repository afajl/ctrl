package config

import (
	"flag"
	"strings"
)

var logdir = flag.String("logdir", "", "run log directory")
var keyfile = flag.String("i", "", "add keyfile")
var verbose = flag.Bool("v", false, "verbose")
var dontlog = flag.Bool("dontlog", false, "dont write any logs")
var hosts = new(hostList)

func init() {
	flag.Var(hosts, "H", "hosts to run against")
}

type hostList []string

func (h *hostList) Set(s string) error {
	*h = strings.Split(s, ",")
	return nil
}

func (h *hostList) String() string {
	return strings.Join(*h, ", ")
}

func FromFlags(c *Config) {
	if !flag.Parsed() {
		panic("FlagConfig called before flags parsed")
	}
	if *logdir != "" {
		c.Logdir = *logdir
	}
	if *dontlog {
		c.DontLog = true
	}
	if *verbose {
		c.Verbose = true
	}
	if *hosts != nil {
		c.Hosts = *hosts
	}
	if *keyfile != "" {
		c.Keyfiles = append(c.Keyfiles, *keyfile)
	}
}
