package fork

import (
	"fmt"
	"os"
	"reflect"
)

// the registry of forks we know about
var forks map[string]*Fork

const (
	nameVar = "GOFORK_NAME"
	selfExe = "/proc/self/exe"
)

func init() {
	forks = make(map[string]*Fork)
}

// Register records a Fork in the internal fork map.
// Register must be called before Fork on agiven function (func init() is a common place)
func Register(f *Fork) {
	if f.Name == "" {
		panic("tried to register fork with no name")
	}
	forks[f.Name] = f
}

// RegisterFunc records a function as a fork in the internal fork map.
// Register must be called before Fork on agiven function (func init() is a common place)
func RegisterFunc(n string, fn func(interface{})) {
	Register(NewFork(n, fn))
}

// ForkInit should be called a the point that forks should begin to execute.
// This should likely be very early in main() or within init() (to skip main entirely)
//
// At the point where this is called, either:
//
// We are identified as a fork, we execute that function, and exit.
// or, we are not identified as afork and simply return.
func ForkInit() {
	fmt.Printf("ForkInit(): ")
	var name string
	if os.Args[0] != selfExe {
		// we're not a self-exec
		fmt.Printf("not a fork\n")
		return
	}
	if name = os.Getenv(nameVar); name == "" {
		// no func is defined
		fmt.Printf("not a fork\n")
		return
	}
	fmt.Printf("a fork: %s\n", name)
	os.Unsetenv(nameVar)
	// we appear to be a fork
	if f, ok := forks[name]; ok {
		v := reflect.ValueOf(f.fn)
		if err := v.Call([]reflect.Value{}); err != nil {
			fmt.Printf("fork failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	panic("no fork by name: " + name)
}
