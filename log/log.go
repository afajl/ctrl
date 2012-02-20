// TODO:
//  Fork the log pkg and:
//   - make out logger use time from start
//   - implement:
//      OutRed      OutGreen     Out     OutGray
//      LogErr      LogInfo      Log     LogDebug 
//
//      errors                           commands to run
//                                       command output
//
//   - make output to writes on a new line prefixed with:
//     1:32:10 WORKSTATION: ls -l /etc
//      ! /etc/passwd
//      ! /etc/hosts
//
//   - Use one go routine for run logs and one for every host
//   See https://github.com/ngmoco/timber/blob/master/timber.go
//
// RunLog
//         RunLog
//      -v EchoLog
// RunOut  RunLog
//         EchoLog
// HostLog
//         HostLog     LogWriter
//      -v EchoLog
// HostOut
//         HostLog
//         EchoLog
//
package log

import (
	"bytes"
	"fmt"
	"github.com/afajl/ctrl/path"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
    //"strings"
	"time"
)

var echoOutput io.Writer
var fileOutput io.Writer

// Logger that implements th io.Writer interface
type WriteLogger struct {
	log.Logger
	wbuf bytes.Buffer
}

func NewWriteLogger(w io.Writer, prefix string, flags int) *WriteLogger {
	l := &WriteLogger{}
	l.Logger = *log.New(w, prefix, flags)
	return l
}

func (l *WriteLogger) Write(b []byte) (n int, err error) {
    n, _ = l.wbuf.Write(b)
    defer l.wbuf.Reset()
    err = l.Output(2, l.wbuf.String())
	return
}

// Noop io.Writer for testing
type noopwriter struct{}

func (*noopwriter) Write(b []byte) (int, error) {
	return len(b), nil
}

type HostLoggers struct {
	log, out *WriteLogger
}

type RunLogs struct {
	verbose, usefs          bool
	dir, hostDir            string
	runLog, runOut, echoLog *WriteLogger
	hostLoggers             map[string]*HostLoggers
}

func NewRunLogs(dir string, args []string, verbose, usefs bool) *RunLogs {
	if echoOutput == nil {
		echoOutput = os.Stdout
	}
	l := &RunLogs{}
	l.verbose = verbose
	l.usefs = usefs
	l.hostLoggers = make(map[string]*HostLoggers)

	l.dir = l.createRundir(dir, args)
	l.hostDir = l.dir + "/hosts"

	l.echoLog = NewWriteLogger(echoOutput, "", log.Ltime)

	l.runLog, l.runOut = l.makeLogOut("CTRL ", l.dir, "run.log")
	return l
}

func (l *RunLogs) GetRunLog() *WriteLogger {
	return l.runLog
}

func (l *RunLogs) GetRunOut() *WriteLogger {
	return l.runOut
}

func (l *RunLogs) GetHostLog(hostname string) *WriteLogger {
	return l.getHostLoggers(hostname).log
}

func (l *RunLogs) GetHostOut(hostname string) *WriteLogger {
	return l.getHostLoggers(hostname).out
}

func (l *RunLogs) getHostLoggers(hostname string) *HostLoggers {
	if hostloggers, ok := l.hostLoggers[hostname]; ok {
		return hostloggers
	}
	log, out := l.makeLogOut(fmt.Sprintf("%s: ", hostname), l.hostDir, hostname)

	hostloggers := &HostLoggers{log: log, out: out}
	l.hostLoggers[hostname] = hostloggers
	return hostloggers
}

func (l *RunLogs) makeLogOut(outprefix string, path ...string) (filelog *WriteLogger, out *WriteLogger) {
	filelog = NewWriteLogger(l.getFile(path...), "", log.LstdFlags|log.Lmicroseconds)
	out = NewWriteLogger(io.MultiWriter(filelog, l.echoLog), outprefix, 0)
	if l.verbose {
		filelog = out
	}
	return
}

func (l *RunLogs) getFile(path ...string) io.Writer {
	if !l.usefs {
		if fileOutput != nil {
			return fileOutput
		}
		return &noopwriter{}
	}
	file, err := os.OpenFile(filepath.Join(path...), syscall.O_WRONLY|syscall.O_CREAT, 0660)
	if err != nil {
		panic(fmt.Sprintf("could not open logfile %s: %v", path, err))
	}
	return file
}

func (l *RunLogs) runlogPath(dir string, args []string) string {
	ts := time.Now().Format("2006-01-02-15.04.05.9999")
	ident := path.CleanPathname(args...)
	path := filepath.Join(dir, ts+"_"+ident)
	return path
}

func (l *RunLogs) createRundir(dir string, args []string) string {
	path := l.runlogPath(dir, args)
	if !l.usefs {
		return path
	}
	if err := os.MkdirAll(path+"/hosts", 0750); err != nil {
		panic(fmt.Sprintf("could not create run log dir %s: %s", path, err))
	}
	return path
}
