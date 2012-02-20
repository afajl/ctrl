package shell

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

type Result struct {
	Output string
	Error  error
}

const defaultShell = "bash -c"

var shellRx = regexp.MustCompile(`([ \t'\"\\$])`)

func ShellEscape(args ...interface{}) []interface{} {
	for i, arg := range args {
		switch arg.(type) {
		case string:
			args[i] = shellRx.ReplaceAllString(fmt.Sprintf("%s", arg), `\$1`)
		case []byte:
			args[i] = []byte(shellRx.ReplaceAllString(fmt.Sprintf("%s", arg), `\$1`))

		}
	}
	return args
}

func Commandf(format string, args ...interface{}) (string, error) {
	res := fmt.Sprintf(format, ShellEscape(args...)...)
	// TODO Fix bug in fmt,   MISSING is not starting with %!
	for _, m := range []string{"%!", "(MISSING)"} {
		if strings.Contains(res, m) {
			return "", fmt.Errorf("error formating cmd: %s", res)
		}
	}
	return res, nil
}

func ShCommandf(sh string, cmdf string, a ...interface{}) (shell string, shellargs []string, cmd string, err error) {
	if sh == "" {
		sh = defaultShell
	}
	shparts := strings.Split(sh, " ")

	shell = shparts[0]
    shellargs = shparts[1:]
	cmd, err = Commandf(cmdf, a...)
	return
}

func mapToEnv(env map[string]string) ([]string, error) {
	ret := make([]string, len(env))
	for k, v := range env {
		item, err := Commandf("%s=%q", k, v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, item)
	}
	return ret, nil
}

type Cmd struct {
	*exec.Cmd
	Shell      string
	Format     string
	Logger     io.Writer
	FormatArgs []interface{}
}

func Command(cmdf string, a ...interface{}) *Cmd {
	cmd := &Cmd{}
	cmd.Cmd = &exec.Cmd{}
	cmd.Format = cmdf
	cmd.FormatArgs = a
	return cmd
}

func (c *Cmd) String() string {
    _, _, cmd, err := ShCommandf(c.Shell, c.Format, c.FormatArgs...)
    if err != nil {
        return "error formatting cmd: " + err.Error()
    }
    return cmd
}

func (c *Cmd) SetEnvMap(env map[string]string) (err error) {
	c.Env, err = mapToEnv(env)
	return
}

func (c *Cmd) formatCmd() (shell string, args []string, err error) {
    var relShell, cmd string
	relShell, args, cmd, err = ShCommandf(c.Shell, c.Format, c.FormatArgs...)
    if err != nil {
        return
    }
	shell, err = exec.LookPath(relShell)
	if err != nil {
		shell = relShell
	}
    args = append(append([]string{shell}, args...), cmd)
    return
}


func (c *Cmd) Start() (err error) {
    c.Path, c.Args, err = c.formatCmd()
	if err != nil {
		return
	}
	return c.Cmd.Start()
}

func (c *Cmd) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *Cmd) CombinedOutput() (string, error) {
	if c.Stdout != nil {
		return "", errors.New("exec: Stdout already set")
	}
	if c.Stderr != nil {
		return "", errors.New("exec: Stderr already set")
	}
	var b bytes.Buffer
	if c.Logger != nil {
		mw := io.MultiWriter(c.Logger, &b)
		c.Stdout = mw
		c.Stderr = mw
	} else {
		c.Stdout = &b
		c.Stderr = &b
	}
	err := c.Run()
	return b.String(), err
}
