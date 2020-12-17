# go-fork

Go fork is a package to provide fork-like *emulation* for Go.  In reality, it's not very fork-like, but provides for a common thing `fork()` is used for, namely making a subroutine (goroutine) run in a new process spaces.

Go, by its natively threaded nature, is not fork-safe.  This pkg does not implement a *real* fork. 

Instead, this package makes a way to reroute go to execute a goroutine other than `main()` in a new process.  This allows easy creation of some fork-like behaviors without actually using forks.

Moreover, while `go-fork` children do not share any memory space with `go-fork` parents, the parent can pass arbitrary data to `go-fork` via function arguments.

Here are a couple of particularly important diffierence between real forks and go-forks:

- Go fork does not continue execution from the fork, but rather starts the specified goroutine from whatever point `fork.Init()` is called (usually early in `main()` or in `init()`).
- Go fork doesn't share memory, file pointers, mutexes, or really anything with the child process.  It passes function arguments to the child by encoding them to/decoding them from a temporary file.

# How it works

To use `go-fork` you must do two things:
1. `Register` functions to be forkable (this must be done before trying to fork).b
2. Call `fork.Init()` somewhere early in the code.  If the process is a child, execution will be taken over from this point, and the code will `os.Exit` when it's done, never returning execution to anything after `fork.Init()` is called.

It should be noted that `go-fork` is not able to detect function calling errors at build time.  Errors like incorrect argument assignments are *runtime* errors.

`go-fork` determines that an run is a fork by looking for special values in its environment, which it will immediately unset once read.

`go-fork` passes arguments by `encoding/gob` endoding/decoding them in a temporary file. This fill will be cleaned up immediately after it is read.

# Example:

```go
func init() {
	fork.RegisterFunc("child", child)
	fork.Init()
}

func child(n int) {
	fmt.Printf("child(%d) pid: %d\n", n, os.Getpid())
}

func main() {
	fmt.Printf("main() pid: %d\n", os.Getpid())
	if err := fork.Fork("child", 1); err != nil {
		log.Fatalf("failed to fork: %v", err)
	}
}
```

This will output:

```
$ go run example/example.go 
main() pid: 164120
child(1) pid: 164125
```