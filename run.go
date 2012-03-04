package ctrl

import (
	"flag"
	"fmt"
	"github.com/afajl/ctrl/config"
	"github.com/afajl/ctrl/log"
	"github.com/afajl/ctrl/queue"
	"github.com/afajl/ctrl/remote"
	"os"
)

var (
	listCmds   = flag.Bool("l", false, "list commands")
	configfile = flag.String("c", "", "config file")
)

type Run struct {
	Config   *config.Config
	Queue    *queue.Queue
	Hosts    []*remote.Host
	Cmds     []*RoutedCmd
	Log, Out *log.WriteLogger
	loggers  *log.RunLogs
}

func NewRun() *Run {
	run := &Run{}
	run.Queue = queue.NewQueue()
	return run
}

// Wrap a cmd to create a QueuedCmd
func makeQueuedCmd(cmd Cmd, ctrl Ctrl) queue.QueuedCmd {
	return func() error {
		return cmd(ctrl)
	}
}

func getConfig() *config.Config {
	if err := config.Init(*configfile); err != nil {
		exit_usage("could not load config: ", err)
	}
	conf := config.StartConfig

	if len(conf.Hosts) == 0 {
		exit_usage("no hosts specified")
	}

	return conf
}

func Start(routes *Routes) {
	flag.Usage = usage
	flag.Parse()

	if *listCmds {
		routes.Print()
		os.Exit(0)
	}

	run := NewRun()

	conf := getConfig()

	run.loggers = log.NewRunLogs(conf.Logdir, os.Args[1:], conf.Verbose, !conf.DontLog)
	run.Log = run.loggers.GetRunLog()
	run.Out = run.loggers.GetRunOut()

	var err error
	if run.Hosts, err = remote.NewHosts(conf.Hosts); err != nil {
		exit_usage(err)
	}

	if run.Cmds, err = routes.Parse(flag.Args()); err != nil {
		exit_usage(err)
	}

	if err := run.Run(); err != nil {
		run.Fail(err)
	}
}

func (run *Run) Run() error {
	for _, cmd := range run.Cmds {
		run.QueueCmd(cmd.GetCmd(), run.Hosts...)
	}
	return run.Queue.Run()
}

func (run *Run) QueueCmd(cmd Cmd, hosts ...*remote.Host) {
	for _, host := range hosts {
		ctrl := NewCtrl(run).ForHost(*host)
		ctrl.log = run.loggers.GetHostLog(host.Name)
		ctrl.out = run.loggers.GetHostOut(host.Name)
		qcmd := makeQueuedCmd(cmd, ctrl)
		run.Queue.Add(qcmd)
	}
}

func (run *Run) Fail(a ...interface{}) {
	run.Out.Fatalln(append([]interface{}{"run stopped:"}, a...)...)
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s [OPTION]... COMMANDS...\n", os.Args[0])
	flag.PrintDefaults()
}

func exit_usage(args ...interface{}) {
	fmt.Fprintf(os.Stderr, "\nERROR: %s\n\n", fmt.Sprint(args...))
	usage()
	os.Exit(1)
}
