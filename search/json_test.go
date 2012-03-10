package search

import (
	"testing"
	"github.com/afajl/assert"
	"github.com/afajl/ctrl/host"
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

var test_json = `
{
	"a": {
		"name": "aname",
		"tags": ["aonly", "all", "ab", "ac"],
		"port": "8080",
		"user": "auser",
		"keyfiles": ["akeyx", "akeyy"],
		"remoteshell": "bash -a",
		"remotecd": "/a",
		"remoteenv": {
			"aenvx": "x",
			"aenvy": "y"
		}
	},
	"b": {
		"name": "bname",
		"tags": ["bonly", "all", "ab", "bc"],
		"port": "8080",
		"user": "buser",
		"keyfiles": ["bkeyx", "bkeyy"],
		"remoteshell": "bash -b",
		"remotecd": "/",
		"remoteenv": {
			"benvx": "x",
			"benvy": "y"
		}
	},
 	"c": {
		"name": "cname",
		"tags": ["conly", "all", "bc", "ac"]
	}
}
`

func jsonTest(t *testing.T) Searcher {
    js, err := JsonFromReader(strings.NewReader(test_json))
	if err != nil {
		t.Fatal(err)
	}
	return js
}


func idTest(t *testing.T, ids ...string) []*host.Host {
	js := jsonTest(t)
	res, err := js.Id(ids...)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func tagTest(t *testing.T, tags ...string) []*host.Host {
	js := jsonTest(t)
	res, err := js.Tags(tags...)
	if err != nil {
		t.Fatal(err)
	}
	return res
}



func TestIdHost(t *testing.T) {
	res := idTest(t, "a")

	assert.Equal(t, len(res), 1)

	a := res[0]

	assert.Equal(t, a.Id, "a")
	assert.Equal(t, a.Name, "aname")
	assert.Equal(t, a.Tags, []string{"ab", "ac", "all", "aonly"})
	assert.Equal(t, a.Port, "8080")
	assert.Equal(t, a.User, "auser")
	assert.Equal(t, a.Keyfiles, []string{"akeyx", "akeyy"})
	assert.Equal(t, a.RemoteShell, "bash -a")
	assert.Equal(t, a.RemoteCd, "/a")
	assert.Equal(t, a.RemoteEnv, map[string]string{"aenvx": "x", "aenvy": "y"})
}

func TestMultipleIds(t *testing.T) {
	res := idTest(t, "a", "b")

	assert.Equal(t, len(res), 2)

	a := res[0]
	b := res[1]

	assert.Equal(t, a.Id, "a")
	assert.Equal(t, a.Name, "aname")
	assert.Equal(t, b.Id, "b")
	assert.Equal(t, b.Name, "bname")
}

func TestMissingIds(t *testing.T) {
	res := idTest(t, "X")
	assert.Equal(t, len(res), 0)
}

func TestTagSingle(t *testing.T) {
	res := tagTest(t, "aonly")
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].Id, "a")

	res = tagTest(t, "conly")
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].Id, "c")
}

func TestTagUnion(t *testing.T) {
	res := tagTest(t, "ab", "ac")
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].Id, "a")
}

func TestTagMany(t *testing.T) {
	res := tagTest(t, "ac")
	assert.Equal(t, len(res), 2)
	for _, r := range res {
		switch r.Id {
		case "a", "c":
			//pass
		default:
			t.Fatal("tags should not return id", r.Id)
		}
	}
}
