# go-fork

Go fork is a simple pkg to provide fork-like *emulation* for Go.  In reality, it's not very fork-like, but provides for a common thing `fork()` is used for, namely making a subroutine (goroute) run in a new process space.

Go, by its natively threaded nature, is not fork-safe.  This pkg does not implement a *real* fork. 

Instead, this package makes a way to reroute go to execute a goroutine other than `main()` in a new process.  This allows easy creation of some fork-like behaviors without actually using forks.

Here are a couple of particularly important diffierence between real forks and go-forks:

- Go fork does not continue execution from the fork, but rather starts the specified goroutine.
- Go fork doesn't share memory, file pointers, mutexes, or really anything with the child process.  

Go-fork does allow you to pass an `interface {}` to the child process, but it's up to the child process to interpret it. This is handled by passing a gob to the child process.
