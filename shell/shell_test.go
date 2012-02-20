package shell

import (
	"bytes"
	"testing"
)

type islice []interface{}

func TestShellEscape(t *testing.T) {
	type shtest struct {
		args   islice
		expect islice
	}
	tests := []shtest{
		{islice{"a"}, islice{"a"}},

		// escaped
		{
			islice{" ", "\t", "'", `"`, `\`, `$`},
			islice{`\ `, "\\\t", `\'`, `\"`, `\\`, `\$`},
		},

		// allowed
		{islice{[]byte("foo")}, islice{[]byte("foo")}},

        // unescaped
		{islice{1, 3.14}, islice{1, 3.14}},
		{islice{int8(1), float64(1)}, islice{int8(1), float64(1)}},
	}

	for _, test := range tests {
		res := ShellEscape(test.args...)
		if len(res) != len(test.expect) {
			t.Fatalf("wrong len of res: %#v != %#v", test.args, res)
		}
		for i := 0; i < len(test.expect); i++ {
            // test return of []byte
			if v, ok := test.expect[i].([]byte); ok {
				if r, ok := res[i].([]byte); ok {
					if !bytes.Equal(v, r) {
						t.Fatalf("[]bytes %s != %s", res[i], test.expect[i])
					}
				}
			} else if res[i] != test.expect[i] {
				t.Fatalf("%#v != %#v", res[i], test.expect[i])
			}
		}
	}
}

func TestCommandf(t *testing.T) {
	type cmdtest struct {
		format string
		args   islice
		expect string
		ok     bool
	}
	tests := []cmdtest{
		{"cmd %v", islice{1}, "cmd 1", true},

		// bad nr args
		{"%d %d", islice{1}, "", false},
		{"%d", islice{1, 2}, "", false},

		// bad type
		{"%d", islice{"foo"}, "", false},
		{"%f", islice{1}, "", false},

		// escaping
		{"%s %%", islice{"s"}, "s %", true},
		{"%s %%s", islice{"s"}, "s %s", true},
	}
	for _, test := range tests {
		res, err := Commandf(test.format, test.args...)
		if err == nil && !test.ok {
			t.Fatalf("should fail %#v, got: %s", test, res)
		}
		if err != nil && test.ok {
			t.Fatalf("should not fail (%s): %#v", err, test)
		}
		if res != test.expect {
			t.Fatalf("%q != %q", res, test.expect)
		}
	}
}

func TestShCommand(t *testing.T) {
	type shtest struct {
		sh, cmd string
		args    islice
		expShell  string
		expShArgs []string
		expCmd string
	}
	shtests := []shtest{
		{"bash -c", "echo %s", islice{"yes"},
		 "bash", []string{"-c"}, `echo yes`},
	}
	for _, test := range shtests {
		shell, shArgs, cmd, err := ShCommandf(test.sh, test.cmd, test.args...)
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if shell != test.expShell {
			t.Errorf("bin %v != %v", shell, test.expShell)
		}
		if len(shArgs) != len(test.expShArgs) {
			t.Fatalf("bad nr args %v != %v", shArgs, test.expShArgs)
		}
		for i := 0; i < len(shArgs); i++ {
			if shArgs[i] != test.expShArgs[i] {
				t.Errorf("arg %d %v != %v", i, shArgs[i], test.expShArgs[i])
			}
		}
        if cmd != test.expCmd {
            t.Errorf("cmd %v != %v", cmd, test.expCmd)
        }
	}
}


func TestCommandSimple(t *testing.T) {
    type cmdtest struct {
        cmdf string
        args islice
        exp string
    }

}
