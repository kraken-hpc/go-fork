package fork

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
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
	// Where to send stdout (default: os.Stdout)
	Stdout *os.File
	// Where to send stderr (default: os.Stderr)
	Stderr *os.File
	// Where to get stdin (default: os.Stdin)
	Stdin *os.File

	// contains filtered or unexported fields
	c  exec.Cmd
	fn reflect.Value
}

// NewFork createas and initializes a Fork
// A Fork object can be manipluated to control how a process is launched.
// E.g. you can set new namespaces in the SysProcAttr property...
//      or, you can set custom args with the (optional) variatic args aparameters.
//      If you set args, the first should be the program name (Args[0]), which may
//		Which may or may not match the  executable.
// If no args are specified, args is set to []string{os.Args[0]}
func NewFork(n string, fn interface{}, args ...string) (f *Function) {
	f = &Function{}
	f.c = exec.Cmd{}
	// Process and ProcessState don't actually exist at this point
	f.SysProcAttr = f.c.SysProcAttr
	// os.Executable might not be the most robust way to do this, but it is portable.
	f.c.Path, _ = os.Executable()
	f.c.Args = args
	f.c.Stderr = os.Stderr
	f.c.Stdout = os.Stdout
	f.c.Stdin = os.Stdin
	if len(args) == 0 {
		f.c.Args = []string{os.Args[0]}
	}
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
func (f *Function) Fork(args ...interface{}) (err error) {
	if err = f.validateArgs(args...); err != nil {
		return
	}
	f.c.Stderr = f.Stderr
	f.c.Stdout = f.Stdout
	f.c.Stdin = f.Stdin
	f.c.Env = os.Environ()
	f.c.Env = append(f.c.Env, nameVar+"="+f.Name)
	af, err := ioutil.TempFile("", "gofork_*")
	f.c.Env = append(f.c.Env, argsVar+"="+af.Name())
	if err != nil {
		return
	}
	enc := gob.NewEncoder(af)
	for _, iv := range args {
		enc.EncodeValue(reflect.ValueOf(iv))
	}
	af.Close()
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

// private

func (f *Function) validateArgs(args ...interface{}) (err error) {
	t := f.fn.Type()
	if len(args) != t.NumIn() {
		return fmt.Errorf("incorrect number of args for: %s", t.String())
	}
	for i := 0; i < t.NumIn(); i++ {
		if t.In(i).Kind() != reflect.TypeOf(args[i]).Kind() {
			return fmt.Errorf("argument mismatch (1) %s != %s", reflect.TypeOf(args[i]).Kind(), t.In(i).Kind())
		}
	}
	return
}
