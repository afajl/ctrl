package remote

import (
	"fmt"
	"github.com/afajl/ctrl/config"
	"strings"
	"unicode"
)

const OnWorkstationName = "WORKSTATION"

type Host struct {
	Id            string // Unique id for the host, default Name
	Tags		[]string
	Name          string // IP or hostname to connect to
	Port          string
	User          string
	Keyfiles      []string
	OnWorkstation bool
	RemoteShell string
	RemoteCd    string
	RemoteEnv   map[string]string
}

func NewHost(s string) (*Host, error) {
	host := &Host{RemoteEnv: make(map[string]string)}

	// copy config settings
	host.Port = config.StartConfig.Port
	host.User = config.StartConfig.User
	host.Keyfiles = config.StartConfig.Keyfiles

	// TODO use better parser
	err := host.Set(s)
	return host, err
}

func NewHosts(h []string) (hosts []*Host, err error) {
	nr_hosts := len(h)
	hosts = make([]*Host, nr_hosts)
	for i := 0; i < nr_hosts; i++ {
		hosts[i], err = NewHost(h[i])
		if err != nil {
			return
		}
	}
	return
}

func combine(a *Host, b *Host) *Host {
	c := new(Host)
	*c = *a

	if b.Id != "" { c.Id = b.Id }
	c.Tags = append(c.Tags, b.Tags...)
	if b.Name != "" { c.Name = b.Name }
	if b.Port != "" { c.Port = b.Port }
	if b.User != "" { c.User = b.User }
	c.Keyfiles = append(c.Keyfiles, b.Keyfiles...)

	if b.RemoteShell != "" { c.RemoteShell = b.RemoteShell }
	if b.RemoteCd != "" { c.RemoteCd = b.RemoteCd }
	for k, v := range b.RemoteEnv {
		c.RemoteEnv[k] = v
	}
	return c
}

func Fold(hosts ...*Host) (host *Host) {
	switch n := len(hosts); n {
	case 0:
		// pass
	case 1:
		host = hosts[0]
	default:
		host = hosts[0]
		for i := 1; i < n; i++ {
			host = combine(host, hosts[i])
		}
	}
    return
}

func (h *Host) String() string {
	return h.Id
}

func (h *Host) ConnStr() string {
	if h.Port != "" {
		return h.Name + ":" + h.Port
	}
	return h.Name
}

/*
func getuser(rawhost string) (user, rest string, err error) {
	for i := 0; len(rawhost) < i; i++ {
		c := rawhost[i]
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' || c == '_' || c == '-':
			// continue
		case c == '@':
			if i == 0 {
				return "", "", errors.New("missing username")
			}
			return rawhost[0:i], rawhost[i+1:], nil
		default:
			// invalid character
			return "", "", fmt.Errorf("invalid character in username: %s", c)
		}
	}
	return "", rawhost, nil
}

func gethost(rawhost string) (host, rest string, err error) {
    if strings.Contains(rawhost, "@") {
        _, rawhost, err = getuser(rawhost)
        if err != nil {
            return "", "", err
        }
    }
    for i := 0; len(rawhost) < i; i++ {
        c := rawhost[i]
        switch {
        case 'a' <= c && c <= 'z' || '0' <= c && c <= '9': 
            // continue
        case c == '-':
            if i == 0 {
                return "", "", errors.New("hostnames cant start with a dash")
            }
        case c == ':':
            if i == 0 {
                return "", "", errors.New("missing hostname")
            }
            return rawhost[0:1], rawhost[i+1:], nil
    }
    return rawhost, "", nil
}

func getport(rawhost string) (port string, err error) {
    if strings.Contains(rawhost, ":") {
        _, rawhost, err = gethost(rawhost)
        if err != nil {
            return "", "", err
        }
    }
    for i := 0; len(rawhost) < i; i++ {

}
*/

func (h *Host) Set(s string) (err error) {
	s = strings.TrimSpace(s)

	if strings.IndexFunc(s, unicode.IsSpace) >= 0 {
		err = fmt.Errorf("host contains space")
		return
	}

	switch c := strings.Count(s, "@"); {
	case c == 0:
	case c == 1:
		userhost := strings.SplitN(s, "@", 2)
		h.User = userhost[0]
		s = userhost[1]
	case c > 1:
		return fmt.Errorf("more then one @ in host")
	}

	switch c := strings.Count(s, ":"); {
	case c == 0:
		h.Name = s
	case c == 1:
		hostport := strings.SplitN(s, ":", 2)
		h.Name = hostport[0]
		h.Port = hostport[1]
	case c > 1:
		return fmt.Errorf("more then one @ in host")
	}
	if h.Name == "" {
		return fmt.Errorf("zero length host")
	}

	if h.Name == OnWorkstationName {
		h.OnWorkstation = true
	}
	return
}
