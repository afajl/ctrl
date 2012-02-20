package log

import (
	"bytes"
	"github.com/afajl/ctrl/tests"
	"io/ioutil"
	"strings"
	"testing"
)

func TestFunc(t *testing.T) {
	logdir := tests.TempDir()
	defer tests.ClearDir(logdir)

	var out bytes.Buffer
	echoOutput = &out

	l := NewRunLogs(logdir, []string{"deploy", "-P"}, false, true)

	l.GetRunOut().Println("RUNOUT")
	l.GetRunLog().Println("RUNLOG")
	l.GetHostOut("host").Println("HOSTOUT")
	l.GetHostLog("host").Println("HOSTLOG")

	contains := func(name, s string, lines ...string) {
		for _, line := range lines {
			if !strings.Contains(s, line) {
				t.Fatalf("%s does not contain '%s': %s", name, s, line)
			}
		}
	}
	filecontains := func(path string, lines ...string) {
		s, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatalf("could not open %s: %s", s, err)
		}
		contains(path, string(s), lines...)
	}

	contains("stdout", out.String(), "RUNOUT", "HOSTOUT")
	filecontains(l.dir+"/run.log", "RUNOUT", "RUNLOG")
	filecontains(l.hostDir+"/host", "HOSTOUT", "HOSTLOG")

}
