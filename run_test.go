package ctrl

import (
	"github.com/afajl/ctrl/tests"
	"os"
	"strings"
	"testing"
)

func TestFuncRun(t *testing.T) {
	var nr_hosts int
	hosts := []string{"a", "b"}
	var seen_hosts [2]string

	cmd := func(ctrl Ctrl) error {
		seen_hosts[nr_hosts] = ctrl.Host().Name
		nr_hosts++
		return nil
	}

	routes := NewRoutes()
	routes.AddFunc("cmd", cmd, "help")

	os.Args = []string{"ctrl", "-dontlog", "-c", tests.SimpleConfig,
		"-H", strings.Join(hosts, ","), "cmd"}

	Start(routes)

	if nr_hosts != 2 {
		t.Fatal("nr hosts seen are not 2", nr_hosts)
	}
	for i := 0; i < len(hosts); i++ {
		if hosts[i] != seen_hosts[i] {
			t.Error("did not see host", hosts[i])
		}
	}
}

/*type TestCmd bool*/
/*func (c *TestCmd) Run(ctrl Ctrl) {}*/
/*type cmdmethod struct {*/
/*}*/
/*func (c *cmdmethod) F(ctrl Ctrl) error {*/
/*return nil*/
/*}*/

/*func TestMethodRun(t *testing.T) {*/
/*hosts := []string{"a", "b"}*/
/*cmd := &cmdmethod{}*/
/*routes := NewRoutes()*/
/*routes.AddFunc("cmd", cmd.F, "help")*/

/*startWithArgs([]string{"-c", emptyConfig, "-H", strings.Join(hosts, ","), "cmd"}, routes)*/
/*}*/
