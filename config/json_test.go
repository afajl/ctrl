package config

import (
	"testing"
)

func TestOverloading(t *testing.T) {
	conf := NewConfig()
	conf.DontLog = true
	conf.LocalCd = "orig"
	conf.LocalShell = "dontmodify" // not in json, don't touch

	json := []byte(`{"dontlog": false, "localcd": "json"}`)

	if err := FromJsonBytes(conf, json); err != nil {
		t.Fatal(err)
	}

	if conf.DontLog != false {
		t.Errorf("dontlog %v != %v", conf.DontLog, false)
	}
	if conf.LocalCd != "json" {
		t.Errorf("localcd %v != %v", conf.LocalCd, "json")
	}
	if conf.LocalShell != "dontmodify" {
		t.Errorf("localshell %v != %v", conf.LocalShell, "dontmodify")
	}
}

func TestLocalEnv(t *testing.T) {
	conf := NewConfig()
	json := []byte(`{"localenv": {"ka": "va", "kb": "va"}}`)

	if err := FromJsonBytes(conf, json); err != nil {
		t.Fatal(err)
	}

	env := conf.LocalEnv
	if env == nil {
		t.Fatalf("LocalEnv is nil")
	}
	if va, present := env["ka"]; present {
		if va != "va" {
			t.Errorf("localenv[ka] %v != %v", va, "va")
		}
	} else {
		t.Errorf("localenv[ka] not present")
	}
}

func TestBadJson(t *testing.T) {
	// TODO this is testing go json implementation....
	type badjson struct {
		test string
		json []byte
	}
	jsons := []badjson{
		{"bad map", []byte(`{"localenv": {"ka": "va", "kb", "va"}}`)},
		{"syntax err", []byte(`{"`)},
		// SHOULD FAIL {"bad key", []byte(`{"NO_SUCH_KEY": true}`)},
	}

	for _, test := range jsons {
		conf := NewConfig()
		if err := FromJsonBytes(conf, test.json); err == nil {
			t.Errorf("bad json %q should fail", test.test)
		}
	}
}

func TestInheritence(t *testing.T) {
	/*aconf := NewConfig(nil)*/
	/*aconf.SetUser("a")*/
	/*if aconf.GetUser() != "a" {*/
	/*t.Fatalf("%v != %v", aconf.GetUser(), "a")*/
	/*}*/

	/*bconf := NewConfig(aconf)*/
	/*if bconf.User.Get() != "a" {*/
	/*t.Fatalf("%v != %v", bconf.User.Get(), "a")*/
	/*}*/
	/*bconf.User.Set("b")*/

	/*aconf := New("test")*/
	/*aconf.SetUser("a")*/

	/*bconf := New("test")*/
	/*bconf.SetParent(aconf)*/

	/*if bconf.GetUser() != "a" {*/
	/*t.Fatalf("%v != %v", bconf.GetUser(), "a")*/
	/*}*/

	/*bconf.SetUser("b")*/
	/*if bconf.GetUser() != "b" {*/
	/*t.Fatalf("%v != %v", bconf.GetUser(), "b")*/
	/*}*/
}
