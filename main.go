package main

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpfel hello hello.c

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

func main() {
	// Subscribe to signals for terminating the program.
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("failed to remove memlock: %v", err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs := helloObjects{}
	if err := loadHelloObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// Attach the program to the tracepoint.
	tp, err := link.Tracepoint("syscalls", "sys_enter_execve", objs.HelloExecve, nil)
	if err != nil {
		log.Fatalf("opening tracepoint: %s", err)
	}
	defer tp.Close()

	fmt.Println("eBPF program successfully loaded and attached!")
	fmt.Println("Tracing execve syscalls... Press Ctrl+C to exit.")
	fmt.Println("To see logs, run commands in another terminal.")
	fmt.Println()

	// Start reading from trace pipe in a goroutine.
	go readTracePipe(stopper)

	// Wait for signal.
	<-stopper
	fmt.Println("\nDetaching program and exiting...")
}

func readTracePipe(stopper chan os.Signal) {
	// Open the trace pipe to read printk output.
	file, err := os.Open("/sys/kernel/tracing/trace_pipe")
	if err != nil {
		log.Printf("error opening trace_pipe: %v (are you running as root/sudo?)", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error reading trace_pipe: %v", err)
			return
		}
		fmt.Print(line)
	}
}
