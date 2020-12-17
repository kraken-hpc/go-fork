package main

import (
	"fmt"
	"os"

	"github.com/jlowellwofford/go-fork"
)

func init() {
	fork.RegisterFunc("child", child)
	fork.Init()
}

func child(n int) {
	fmt.Printf("child(%d) pid: %d\n", n, os.Getpid())
}

func main() {
	fmt.Printf("main() pid: %d\n", os.Getpid())
	fork.Fork("child", 1, 2)
}
