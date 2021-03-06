package remote

import (
	"code.google.com/p/go.crypto/ssh"
	"github.com/afajl/ctrl/host"
	"fmt"
)

var keys = new(keyring)
var conns = make(map[string]*ssh.ClientConn)

func getConn(host *host.Host) (*ssh.ClientConn, error) {
	hostkey := host.Id
	if con, ok := conns[hostkey]; ok {
		return con, nil
	}
	if host.User == "" {
		return nil, fmt.Errorf("user not set")
	}
	for _, keyfile := range host.Keyfiles {
		// TODO add key to global keyring, ok?
		if err := keys.loadPEM(keyfile); err != nil {
			return nil, fmt.Errorf("unable to load %s: %v", keyfile, err)
		}
	}
	config := &ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthKeyring(keys),
		},
	}
	conn, err := ssh.Dial("tcp", host.ConnStr(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %s: %v", host, err)
	}
	conns[hostkey] = conn
	return conn, nil
}

func newSession(host *host.Host) (*ssh.Session, error) {
	conn, err := getConn(host)
	if err != nil {
		return nil, err
	}
	return conn.NewSession()
}
