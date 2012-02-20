package ctrl

import (
	"flag"
	"fmt"
)

type FlagCmd func(*flag.FlagSet) Cmd

type RoutedCmd struct {
	name, help   string
	cmd          Cmd
	flags        *flag.FlagSet
	definesFlags bool
	parallell    *bool
}

func NewRoutedCmd(name string, cmd Cmd, help string, flags *flag.FlagSet) *RoutedCmd {
	rc := &RoutedCmd{name: name, help: help, cmd: cmd, flags: flags}
	if flags == nil {
		rc.flags = flag.NewFlagSet(name, flag.ExitOnError)
	} else {
		rc.definesFlags = true
	}
	// universal cmd flags
	// rc.parallell = rc.flags.Bool("P", false, "run in parallell")
	return rc
}

//func printUniversalFlags() {
//    fmt.Fprintf(os.Stderr, "\nall commands defines these flags:\n")
//    fmt.Fprintf(os.Stderr, "  -P=false: run in parallel\n")
//}

func (rc *RoutedCmd) Match(s string) bool {
	return rc.name == s
}

func (rc *RoutedCmd) copy() *RoutedCmd {
	rccopy := *rc
	return &rccopy
}

func (rc *RoutedCmd) Print() {
	fmt.Printf(" %s:\t%s\n", rc.name, rc.help)
	if rc.definesFlags == true {
		rc.flags.PrintDefaults()
	}
}

func (rc *RoutedCmd) Parse(args []string) []string {
	rc.flags.Parse(args)
	return rc.flags.Args()
}

func (rc *RoutedCmd) GetCmd() Cmd {
	return rc.cmd
}
