package tests

import (
	"io/ioutil"
	"os"
)

var (
	wd, _        = os.Getwd()
	TestFiles    = wd + "/tests/files"
	EmptyConfig  = TestFiles + "/emptyctrl.conf"
	SimpleConfig = TestFiles + "/simplectrl.conf"
)

func TempDir() string {
	dir, err := ioutil.TempDir("/tmp", "ctrl_tmp_")
	if err != nil {
		panic(err)
	}
	return dir
}

func ClearDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}
}
