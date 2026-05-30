# eBPF Labs

A Go-based development environment for eBPF programs using the Cilium eBPF library.

This repository contains a simple "Hello World" eBPF program that hooks into the `sys_enter_execve` tracepoint (executing processes) and prints kernel trace logs back to Go.

## Prerequisites

1. **Linux Kernel**: Ensure you are running on Linux.
2. **Go**: Version 1.18 or higher (installed).
3. **Clang/LLVM**: Used to compile the eBPF C code (installed).
4. **BPF Headers**: System headers installed (e.g., `libbpf-devel` on Fedora/RHEL/CentOS, or `libbpf-dev` on Debian/Ubuntu).

## Project Structure

* [hello.c](file:///home/baxter/Projects/eBPF-Labs/hello.c): eBPF program written in C.
* [main.go](file:///home/baxter/Projects/eBPF-Labs/main.go): Go program that loads the compiled eBPF code, attaches it to the kernel tracepoint, and reads the trace log output.
* `.gitignore`: Git configuration.

## How to Build and Run

### 1. Compile C to eBPF Bytecode and Generate Go Bindings

Use `go generate` which triggers Cilium's `bpf2go` tool to compile the C code and produce Go wrappers:

```bash
go generate
```

This generates `hello_bpfel.go` and `hello_bpfel.o`.

### 2. Build the Go Binary

```bash
go build -o ebpf-hello
```

### 3. Run the Program

Since eBPF programs require root privileges to load and attach programs to the kernel:

```bash
sudo ./ebpf-hello
```

Once running, try executing commands in another terminal (e.g., running `ls`, `mkdir`, or any CLI command). You should see output from the tracepipe showing the command name and PID that triggered the execution:

```
eBPF program successfully loaded and attached!
Tracing execve syscalls... Press Ctrl+C to exit.
To see logs, run commands in another terminal.

           <...>-12345   [001] d..3  1234.567890: bpf_trace_printk: Hello World! execve called by PID 12345 (ls)
```

## How it Works

1. **C BPF Hook**: The eBPF code in [hello.c](file:///home/baxter/Projects/eBPF-Labs/hello.c) defines a function `hello_execve` associated with the tracepoint `syscalls/sys_enter_execve`. When a process runs, the kernel triggers this tracepoint.
2. **Retrieve Information**: In the BPF context, `bpf_get_current_pid_tgid()` retrieves the process ID, and `bpf_get_current_comm()` gets the executable name.
3. **Log to Trace Pipe**: It prints the log using `bpf_printk`, which writes to `/sys/kernel/tracing/trace_pipe`.
4. **Go Consumer**: [main.go](file:///home/baxter/Projects/eBPF-Labs/main.go) loads the program, registers it, opens `/sys/kernel/tracing/trace_pipe`, and prints incoming events in real-time.
