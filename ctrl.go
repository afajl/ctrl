package ctrl

import (
	"fmt"
	"github.com/afajl/ctrl/config"
	"github.com/afajl/ctrl/log"
	"github.com/afajl/ctrl/remote"
	"github.com/afajl/ctrl/shell"
)

type Cmd func(Ctrl) error

type Ctrl struct {
	config   config.Config
	host     *remote.Host
	run      *Run
	log, out *log.WriteLogger
}

func NewCtrl(run *Run) Ctrl {
	return Ctrl{run: run, config: *config.StartConfig}
}

func (c Ctrl) ForHost(host interface{}) Ctrl {
	switch ht := host.(type) {
	case string:
		h, err := remote.ParseHost(ht)
		if err != nil {
			panic(fmt.Sprintf("could not parse host: %s", err))
		}
		c.host = h
	case *remote.Host:
		c.host = new(remote.Host)
		*c.host = *ht
	case remote.Host:
		c.host = &ht
	default:
		panic(fmt.Sprintf("unkown type %q", ht))
	}
	return c
}

func (c Ctrl) LocalEnv(env map[string]string) Ctrl {
	c.config.LocalEnv = env
	return c
}

func (c Ctrl) LocalCd(path string) Ctrl {
	c.config.LocalCd = path
	return c
}

func (c Ctrl) RemoteCd(path string) Ctrl {
	return c
}

// other
func (c Ctrl) Host() *remote.Host {
	return c.host
}

// Generic commands
type GenericCommand interface {
	CombinedOutput() (string, error)
	String() string
}

type CommandBuilder func(Ctrl, string, ...interface{}) (GenericCommand, error)

func (c Ctrl) outputTry(cmdbuild CommandBuilder, cmdf string, a ...interface{}) (string, error) {
	cmd, err := cmdbuild(c, cmdf, a...)
	if err != nil {
		return "", err
	}
	c.Log(cmd)
	return cmd.CombinedOutput()
}

func (c Ctrl) output(cmdbuild CommandBuilder, cmdf string, a ...interface{}) string {
	out, err := c.outputTry(cmdbuild, cmdf, a...)
	if err != nil {
		c.run.Fail(err)
	}
	return out
}

// Local
func (c Ctrl) LocalCommand(cmdf string, args ...interface{}) (GenericCommand, error) {
	cmd := shell.Command(cmdf, args...)
	cmd.Dir = c.config.LocalCd
	cmd.Shell = c.config.LocalShell
	cmd.Logger = c.log
	err := cmd.SetEnvMap(c.config.LocalEnv)
	return cmd, err
}

func (c Ctrl) LocalTry(cmdf string, args ...interface{}) (string, error) {
	return c.outputTry((Ctrl).LocalCommand, cmdf, args...)
}

func (c Ctrl) Local(cmdf string, args ...interface{}) string {
	return c.output((Ctrl).LocalCommand, cmdf, args...)
}

// Remote
func (c Ctrl) RemoteCommand(cmdf string, args ...interface{}) (GenericCommand, error) {
	if c.host.OnWorkstation {
		return c.localRemoteCommand(cmdf, args...)
	}
	rcmd, err := remote.Command(c.host, cmdf, args...)
	if err != nil {
		return nil, err
	}
	rcmd.Dir = c.config.LocalCd
	rcmd.Shell = c.config.LocalShell
	rcmd.Logger = c.log
	err = rcmd.SetEnvMap(c.config.LocalEnv)
	return rcmd, err
}

func (c Ctrl) localRemoteCommand(cmdf string, args ...interface{}) (GenericCommand, error) {
	cmd := shell.Command(cmdf, args...)
	cmd.Dir = c.host.RemoteCd
	cmd.Shell = c.host.RemoteShell
	cmd.Logger = c.log
	err := cmd.SetEnvMap(c.host.RemoteEnv)
	return cmd, err
}

func (c Ctrl) RemoteTry(cmdf string, args ...interface{}) (string, error) {
	return c.outputTry((Ctrl).RemoteCommand, cmdf, args...)
}

func (c Ctrl) Remote(cmdf string, args ...interface{}) string {
	return c.output((Ctrl).RemoteCommand, cmdf, args...)
}

// Log
func (c Ctrl) Outf(format string, a ...interface{}) {
	c.out.Printf(format, a...)
}

func (c Ctrl) Out(a ...interface{}) {
	c.out.Println(a...)
}

func (c Ctrl) Logf(format string, a ...interface{}) {
	c.log.Printf(format, a...)
}

func (c Ctrl) Log(a ...interface{}) {
	c.log.Println(a...)
}
