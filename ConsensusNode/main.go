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
		panic("usage: input id")
	}

	id, _ := strconv.Atoi(os.Args[1])

	config.SetupConfig()
	node := node.NewNode(int64(id))
	node.Run()

	// entropyNode.StartEntropyNode()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("===>Randomness Beacon is running at PID[%s]\n", pid)
	fmt.Println("===============================================")
	fmt.Println("*                                             *")
	fmt.Println("*             Randomness Beacon               *")
	fmt.Println("*                                             *")
	fmt.Println("===============================================")

	sig := <-sigCh
	fmt.Printf("===>Finish by signal[%s]\n", sig.String())
}
