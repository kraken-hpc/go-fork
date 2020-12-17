package fork

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"syscall"
)

// A Function struct describes a fork process.  Usually this is only used internally, but if you want a bit more control of the sub-process,
// say, to assing new namespaces to the child, this will provide it.
//
// A Fork is a wrapper around an `exec.Cmd` with some of parts the data structure exposed (that we don't need to control directly).
type Function struct {
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
	fn reflect.Value
}

// NewFork createas and initializes a Fork
func NewFork(n string, fn interface{}) (f *Function) {
	f = &Function{}
	f.c = exec.Cmd{}
	f.SysProcAttr = f.c.SysProcAttr
	// Process and ProcessState don't actually exist at this point
	f.c.Path = getSelfExe()
	f.c.Args = []string{f.c.Path}
	// we don't check for errors here, but it would be a pretty bad thing if this failed
	//f.c.Args = func() []string { s, _ := ioutil.ReadFile("/proc/self/comm"); return []string{string(s)} }()
	f.fn = reflect.ValueOf(fn)
	if f.fn.Kind() != reflect.Func {
		return nil
	}
	f.Name = n
	return
}

// Fork starts a process and prepares it to call the defined fork
func (f *Function) Fork(args interface{}) (err error) {
	f.c.Env = os.Environ()
	f.c.Stderr = os.Stderr
	f.c.Stdout = os.Stdout
	f.c.Stdin = os.Stdin
	f.c.Env = os.Environ()
	f.c.Env = append(f.c.Env, nameVar+"="+f.Name)
	af, err := ioutil.TempFile("", "gofork_*")
	f.c.Env = append(f.c.Env, argsVar+"="+af.Name())
	if err != nil {
		return
	}
	enc := gob.NewEncoder(af)
	is, ok := args.([]interface{})
	if !ok {
		is = []interface{}{args}
	}
	for _, iv := range is {
		enc.EncodeValue(reflect.ValueOf(iv))
	}
	fmt.Println()
	af.Close()
	if err = f.validateArgs(args); err != nil {
		return
	}
	if err = f.c.Start(); err != nil {
		return
	}
	f.Process = f.c.Process
	return
}

// Wait provides a wrapper around exec.Cmd.Wait()
func (f *Function) Wait() (err error) {
	if err = f.c.Wait(); err != nil {
		return
	}
	f.ProcessState = f.c.ProcessState
	return
}

func (f *Function) validateArgs(a interface{}) (err error) {
	return
}

func getSelfExe() string {
	if runtime.GOOS == "linux" {
		return "/proc/self/exe"
	}
	return os.Args[0]
}
