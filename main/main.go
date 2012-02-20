package main

import (
	"flag"
	"github.com/afajl/ctrl"
)

func Deploy(ct ctrl.Ctrl) error {
	ct.Out("deploying to host", ct.Host())
	ct.Local(`echo -e "%s \n %s"`, "row1", "row2")
	ct.Remote(`echo %s > /tmp/n`, "42")
	return nil
}

func Install(ct ctrl.Ctrl, name string) error {
	ct.Local("echo me")
	ct.Outf("In install for %s", name)
	return nil
}

func installFac(flag *flag.FlagSet) ctrl.Cmd {
	var name = flag.String("n", "", "install name")
	return func(ctrl ctrl.Ctrl) error {
		return Install(ctrl, *name)
	}
}

func main() {
	routes := ctrl.NewRoutes()
	routes.AddFunc("deploy", Deploy, "deploy help")
	routes.AddFlagFunc("install", installFac, "install help")

	ctrl.Start(routes)
}
