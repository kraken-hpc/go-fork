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
// A Fork is a wrapper around an `exec.Cmd` with some of parts the data structure exposed (that we don't need to control directly).
type Fork struct {
	// SysProcAttr holds optional, operating system-sepcific attributes.
	SysProcAttr *syscall.SysProcAttr
	// Process will hold the os.Process once the Fork has been Run().
	Process *os.Process
	// ProcessState will hold the os.ProcessState after Wait() has been called.
	ProcessState *os.ProcessState
	// Name is the string we use to identify this func
	Name string

	// contains filtered or unexported fields
	c  exec.Cmd
	fn func(interface{})
}

// NewFork createas and initializes a Fork
func NewFork(n string, fn func(interface{})) (f *Fork) {
	f = &Fork{}
	f.c = exec.Cmd{}
	f.SysProcAttr = f.c.SysProcAttr
	// Process and ProcessState don't actually exist at this point
	f.c.Path = selfExe
	// we don't check for errors here, but it would be a pretty bad thing if this failed
	f.c.Args = func() []string { s, _ := ioutil.ReadFile("/proc/self/comm"); return []string{string(s)} }()
	f.fn = fn
	f.Name = n
	return
}

// Call starts a process and prepares it to call the defined fork
func (f *Fork) Call() (err error) {
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

// Wait provides a wrapper around exec.Cmd.Wait()
func (f *Fork) Wait() (err error) {
	if err = f.c.Wait(); err != nil {
		return
	}
	f.ProcessState = f.c.ProcessState
	return
}
