package search

import (
	"testing"
	//"github.com/afajl/assert"
	//"github.com/afajl/ctrl/remote"
	"github.com/afajl/ctrl/config"
	"net/url"
	"os"
	"strings"
	"io/ioutil"
)


func TestJsonFac(t *testing.T) {
	type utest struct {
		desc, rawurl string
		ok bool
	}

	tests := []utest{
		{"bad scheme", "badscheme://p", false},
		{"invalid path", "json:///THIS_SHOULD_NOT_EXIST", false},
		{"missing path", "json:///", false},
		{"dir", "json://testdata/", false},
		{"relative path", "json://rel/ok.json", true},
		// TODO add tilde expansion and environ 
	}

	os.Chmod("testdata/rel/noaccess.json", 0000)
	defer os.Chmod("testdata/rel/noaccess.json", 0644)
	tests = append(tests, utest{"unauthorized path", "json://rel/noaccess.json", false})

	// absfile test
	absfile := "/tmp/CTRL_JSON_OK_XXX.json"
	ioutil.WriteFile(absfile, []byte("{}"), 0644)
	defer os.Remove(absfile)
	tests = append(tests, utest{"absolute path", "json://"+absfile, true})

	// rel symlink
	if err := os.Symlink(absfile, "testdata/rel/ok-link.json"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("testdata/rel/ok-link.json")
	tests = append(tests, utest{"relative link", "json://rel/ok-link.json", true})

	// abs symlink
	abslink := "/tmp/CTRL_JSON_OK_XXX-LINK.json"
	if err := os.Symlink(absfile, abslink); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(abslink)
	tests = append(tests, utest{"absolute link", "json://"+abslink, true})


	config.StartConfig.Rootdir = "testdata"

	for _, test := range tests {
		u, err := url.Parse(test.rawurl)
		if err != nil {
			t.Fatal(err)
		}
		_, err = JsonFac(u)
        if test.ok && err != nil {
			t.Error(test.desc, "should succeed:", err)
		}
		if !test.ok && err == nil {
			t.Error(test.desc, "should fail")
		}
	}
}

func TestParsing(t *testing.T) {
	type jtest struct {
		desc, json string
		ok bool
	}
	tests := []jtest{
		{"empty ok", "{}", true},
		{"empty host ok", `{"ahost": {}}`, true},
		{"ok attr", `{"ahost": {"name": "ahost"}}`, true},
		{"empty id", `{"": {}}`, false},
		{"id not obj", `{"id": "i should be an obj"}`, false},
		{"bad type", `{"id": {"Name": 1}}`, false},
		// TODO: Requires 'manual' parsing
		//{"duplicate ids", `{"dup": {}, "dup": {}}`, false},
		//{"bad attr", `{"id": {"UNKNOWN": 1}}`, false},
	}

	for _, test := range tests {
		js, err := JsonFromReader(strings.NewReader(test.json))
		_, err = js.Id("_")

		if test.ok && err != nil {
			t.Error(test.desc, "should succeed:", err)
		}
		if !test.ok && err == nil {
			t.Error(test.desc, "should fail")
		}
	}
}

func TestId(t *testing.T) {

}
