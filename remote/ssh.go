package remote

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"github.com/afajl/ctrl/shell"
	"io"
	"strings"
)

// Make interface of ssh.Session similar to exec.Cmd
type Cmd struct {
	*ssh.Session
	Shell      string
	Dir        string
	Format     string
	Logger     io.Writer
	FormatArgs []interface{}
}

func Command(host *Host, cmdf string, a ...interface{}) (cmd *Cmd, err error) {
	cmd = &Cmd{}
	cmd.Session, err = newSession(host)
	if err != nil {
		return
	}
	cmd.Format = cmdf
	cmd.FormatArgs = a
	return
}

func (c *Cmd) String() string {
	_, _, cmd, err := shell.ShCommandf(c.Shell, c.Format, c.FormatArgs...)
	if err != nil {
		return "error formatting cmd: " + err.Error()
	}
	return cmd
}

func (c *Cmd) SetEnvMap(env map[string]string) error {
	for k, v := range env {
		c.Setenv(k, v)
	}
	return nil
}

func (c *Cmd) formatCmd() (string, error) {
	cmdf := c.Format
	if c.Dir != "" {
		cmdf = fmt.Sprintf("cd %s && %s", c.Dir, cmdf)
	}

	shell, shellargs, cmd, err := shell.ShCommandf(c.Shell, cmdf, c.FormatArgs...)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`%s %s "%s"`, shell, strings.Join(shellargs, " "), cmd), nil
}

func (c *Cmd) Start() error {
	cmd, err := c.formatCmd()
	if err != nil {
		return err
	}
	return c.Session.Start(cmd)
}

func (c *Cmd) Run() error {
	defer c.Close()
	err := c.Start()
	if err != nil {
		return err
	}
	return c.Session.Wait()
}

func (c *Cmd) CombinedOutput() (string, error) {
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

/*func Put(host *Host, src, target string) error {*/
/*srctar := exec.Cmd("tar", "czf" "-", src)*/

/*}*/
