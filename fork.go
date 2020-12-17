// +build linux

package fork

import (
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

// A Fork struct describes a fork process.  Usually this is only used internally, but if you want a bit more control of the sub-process,
// say, to assing new namespaces to the child, this will provide it.
//
// A Fork is really a wrapper around an `exec.Cmd` with some of parts the data structure exposed (that we don't need to control directly).
type Fork struct {
	// SysProcAttr holds optional, operating system-sepcific attributes.
	SysProcAttr *syscall.SysProcAttr
	// Process will hold the os.Process once the Fork has been Run().
	Process *os.Process
	// ProcessState will hold the os.ProcessState after Wait() has been called.
	ProcessState *os.ProcessState

	// contains filtered or unexported fields
	c        exec.Cmd
	args     []interface{}
	pkgName  string
	funcName string
}

// NewFork createas and initializes a Fork
func NewFork(fn *func(...interface{}), args ...[]interface{}) (f *Fork) {
	f = &Fork{}
	f.c = exec.Cmd{}
	f.SysProcAttr = f.c.SysProcAttr
	// Process and ProcessState don't actually exist at this point
	f.c.Path = "/proc/self/exe"
	// we don't check for errors here, but it would be a pretty bad thing if this failed
	f.c.Args = func() []string { s, _ := ioutil.ReadFile("/proc/self/comm"); return []string{string(s)} }()
	f.funcName = reflect.
	return
}

func (f *Fork) Run() (err error) {
	f.c.Env = os.Environ()
	f.c.Stderr = os.Stderr
	f.c.Stdout = os.Stdout
	f.c.Stdin = os.Stdin
	if err = f.c.Start(); err != nil {
		return
	}
	f.Process = f.c.Process
	return
}

func (f *Fork) Wait() (err error) {
	if err = f.c.Wait(); err != nil {
		return
	}
	f.ProcessState = f.c.ProcessState
	return
}
