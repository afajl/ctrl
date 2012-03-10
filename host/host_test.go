package host

import (
	"testing"
)

type hosttest struct {
	s        string
	exp_host string
	exp_user string
	exp_port string
	fail     bool
}

var hosttests = []hosttest{
	{" a ", "a", "", "", false},
	{"u@a", "a", "u", "", false},
	{"u@a:1", "a", "u", "1", false},
	{"u @a", "", "", "", true},
	{"@a:", "a", "", "", false},
	{"@a :2", "", "", "", true},
	{"u:a:2", "", "", "", true},
	{"u@a@2", "", "", "", true},
	{"u@:2", "", "", "", true},
	{"", "", "", "", true},
}

func TestHost(t *testing.T) {
	for _, test := range hosttests {
		host := &Host{}
		if err := host.Set(test.s); err == nil && test.fail {
			t.Fatalf("%q result in %q", test, err)
		} else {
			continue
		}
		if host.Name != test.exp_host {
			t.Fatalf("host %q != %q", host.Name, test)
		}
		if host.User != test.exp_user {
			t.Fatalf("user %q != %q", host.User, test)
		}
		if host.Port != test.exp_port {
			t.Fatalf("port %q != %q", host.Port, test)
		}
	}
}
