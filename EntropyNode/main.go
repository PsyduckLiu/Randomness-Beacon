package main

import (
	"entropyNode/node"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		panic("usage: input id")
	}

	id, _ := strconv.Atoi(os.Args[1])

	node.StartEntropyNode(id)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("===>Entropy Node is running at PID[%s]\n", pid)

	sig := <-sigCh
	fmt.Printf("===>Finish by signal[%s]\n", sig.String())
}
