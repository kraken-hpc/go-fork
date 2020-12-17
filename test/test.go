package main

import (
	"fmt"
	"os"

	"github.com/jlowellwofford/go-fork"
)

func init() {
	fork.RegisterFunc("child", child)
	fork.ForkInit()
}

func child(interface{}) {
	fmt.Printf("child() pid; %d\n", os.Getpid())
}

func main() {
	fmt.Printf("main() pid: %d\n", os.Getpid())
}
