/*
Konfig kan komma från fyra håll:
    Default
    Config fil
    Flags
    Host

Den här bör vara obereonde av dem. Det den borde ha 
är vilken typ som värden förväntas vara, 

De egenskaper config bör ha (i ordning) är:
    Plocka första värden som är satt

*/

package config

import (
	"os/user"
)

var DefaultConfig *Config

type Config struct {
	Verbose bool
	Logdir  string
	DontLog bool
	Hosts   []string
	Rootdir string


	User     string
	Port     string
	Keyfiles []string

	LocalShell string
	LocalCd    string
	LocalEnv   map[string]string
}

func init() {
	c := &Config{}
	if u, err := user.Current(); err == nil {
		c.User = u.Username
	}
	c.Port = "22"
	DefaultConfig = c
}

func NewConfig() *Config {
	return &Config{
		LocalEnv: make(map[string]string),
	}
}
