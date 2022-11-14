package main

import (
	"consensusNode/config"
	"consensusNode/node"
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

	// start consensus node[id]
	id, _ := strconv.Atoi(os.Args[1])
	config.SetupConfig()
	node := node.NewNode(int64(id))
	go node.Run()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("\n===>[Start]Consensus Node[%d] is running at PID[%s]\n", id, pid)

	sig := <-sigCh
	fmt.Printf("===>[Finish]Consensus Node[%d] is finished by signal[%s]\n", id, sig.String())
}
