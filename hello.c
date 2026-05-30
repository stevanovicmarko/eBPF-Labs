//go:build ignore

#include <linux/bpf.h>
#include <linux/types.h>
#include <bpf/bpf_helpers.h>

SEC("tracepoint/syscalls/sys_enter_execve")
int hello_execve(void *ctx) {
    __u64 id = bpf_get_current_pid_tgid();
    __u32 pid = id >> 32;
    char comm[16];
    
    bpf_get_current_comm(&comm, sizeof(comm));
    bpf_printk("Hello World! execve called by PID %d (%s)\n", pid, comm);
    return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
