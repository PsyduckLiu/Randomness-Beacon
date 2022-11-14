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
		panic("[Command line arguments]Usage: input id")
	}

	// start entropy node[id]
	id, _ := strconv.Atoi(os.Args[1])
	go node.StartEntropyNode(id)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("===>[Start]Entropy Node[%d] is running at PID[%s]\n", id, pid)

	sig := <-sigCh
	fmt.Printf("===>[Finish]Entropy Node[%d] is finished by signal[%s]\n", id, sig.String())
}
