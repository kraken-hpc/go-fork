package fork

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
)

// the registry of forks we know about
var forks map[string]*Function

const (
	nameVar = "GOFORK_NAME"
	argsVar = "GOFORK_ARGS"
)

func init() {
	forks = make(map[string]*Function)
}

// Register records a Fork in the internal fork map.
// Register must be called before Fork on agiven function (func init() is a common place)
func Register(f *Function) {
	if f.Name == "" {
		panic("tried to register fork with no name")
	}
	forks[f.Name] = f
}

// RegisterFunc records a function as a fork in the internal fork map.
// Register must be called before Fork on agiven function (func init() is a common place)
func RegisterFunc(n string, fn interface{}) {
	Register(NewFork(n, fn))
}

// Init should be called at the point that forks should begin to execute.
// This should likely be very early in main() or within init() (to skip main entirely)
//
// At the point where this is called, either:
//
// We are identified as a fork, we execute that function, and exit.
// or, we are not identified as afork and simply return.
func Init() {
	var name string
	if name = os.Getenv(nameVar); name == "" {
		// no func is defined
		return
	}
	os.Unsetenv(nameVar)
	// we appear to be a fork
	if f, ok := forks[name]; ok {
		v := f.fn
		t := v.Type()
		args := []reflect.Value{}
		if argsFile := os.Getenv(argsVar); argsFile != "" {
			// get our arguments
			f, err := os.Open(argsFile)
			if err != nil {
				panic("failed to open args file: " + err.Error())
			}
			dec := gob.NewDecoder(f)
			for i := 0; i < t.NumIn(); i++ {
				v := reflect.Indirect(reflect.New(t.In(i)))
				err := dec.DecodeValue(v)
				if err != nil {
					panic("failed to decode arguments from args file: " + err.Error())
				}
				args = append(args, v)
			}
			f.Close()
			os.Remove(argsFile)
			os.Unsetenv(argsVar)
		}
		if t.NumIn() != len(args) {
			panic("fork failed: incorrect number of args supplied")
		}
		if err := v.Call(args); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
	panic("no fork by name: " + name)
}

// Fork calls a registered fork
func Fork(name string, args ...interface{}) (err error) {
	f, ok := forks[name]
	if !ok {
		return fmt.Errorf("no registered function by name: %s", name)
	}
	return f.Fork(args...)
}
