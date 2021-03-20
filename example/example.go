package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kraken-hpc/go-fork"
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
	if err := fork.Fork("child", 1); err != nil {
		log.Fatalf("failed to fork: %v", err)
	}
}
