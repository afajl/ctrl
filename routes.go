package ctrl

import (
	"flag"
	"fmt"
)

type Routes struct {
	cmds []*RoutedCmd
}

func NewRoutes() *Routes {
	return &Routes{}
}

func (r *Routes) addCmd(rc *RoutedCmd) {
	if _, exists := r.Lookup(rc.name); exists {
		panic("cmd name already added")
	}
	r.cmds = append(r.cmds, rc)
}

func (r *Routes) AddFunc(name string, cmd Cmd, help string) {
	r.addCmd(NewRoutedCmd(name, cmd, help, nil))
}

func (r *Routes) AddFlagFunc(name string, fcmd FlagCmd, help string) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	cmd := fcmd(fs)
	r.addCmd(NewRoutedCmd(name, cmd, help, fs))
}

func (r *Routes) MountRoutes(prefix string, cmdmap *Routes) {
	for _, rc := range cmdmap.cmds {
		rccopy := rc.copy()
		rccopy.name = prefix + "." + rc.name
		r.addCmd(rccopy)
	}
}

func (r *Routes) Visit(fn func(*RoutedCmd) bool) {
	for _, rc := range r.cmds {
		if !fn(rc) {
			return
		}
	}
}

func (r *Routes) Lookup(s string) (cmd *RoutedCmd, ok bool) {
	r.Visit(func(rc *RoutedCmd) bool {
		if rc.Match(s) {
			cmd = rc
			ok = true
			return false // stop iteration
		}
		return true
	})
	return
}

func (r *Routes) Print() {
	r.Visit(func(rc *RoutedCmd) bool {
		rc.Print()
		return true
	})
	//printUniversalFlags()
}

func (r *Routes) Parse(args []string) ([]*RoutedCmd, error) {
	rcs := make([]*RoutedCmd, 0, 2)
	return r.parseRec(args, rcs)
}

func (r *Routes) parseRec(args []string, rcs []*RoutedCmd) ([]*RoutedCmd, error) {
	if len(args) == 0 {
		return rcs, nil
	}

	name := args[0]
	/* test this, or remove
	if name[0:1] == "-" {
		return nil, fmt.Errorf("%s does not define any flags", rcs[len(rcs)-1].name)
	}
	*/
	rc, ok := r.Lookup(name)
	if !ok {
		return nil, fmt.Errorf("could not find command %s, use -l flag", name)
	}
	args = rc.Parse(args[1:])
	rcs = append(rcs, rc)
	return r.parseRec(args, rcs)
}
