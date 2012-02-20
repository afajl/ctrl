package remote

import (
	"flag"
	"testing"
)

var (
	sshuser    = flag.String("ssh.user", "", "ssh username")
	sshpass    = flag.String("ssh.pass", "", "ssh password")
	sshprivkey = flag.String("ssh.privkey", "", "ssh privkey file")
)

func TestFuncPublickeyAuth(t *testing.T) {
	if *sshuser == "" {
		t.Log("ssh.user not defined, skipping test")
		return
	}
	host, err := NewHost("localhost:22")
	if err != nil {
		t.Fatal(err)
	}
	host.User = *sshuser
	host.Keyfiles = []string{*sshprivkey}

	cmd, err := Command(host, "echo %s", 42)
	if err != nil {
		t.Fatal(err)
	}

	res, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if string(res) != "42" {
		t.Fatalf("stdout: %s != %v", res, 42)
	}
}
