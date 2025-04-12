#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include "bpf_prog.h"

char LICENSE[] SEC("license") = "GPL";

#define EPERM 1

SEC("kprobe/__x64_sys_openat")
int block_openat(struct pt_regs *ctx) {
    __u32 uid = bpf_get_current_uid_gid() & 0xFFFFFFFF;
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    __u32 pid = pid_tgid >> 32;

    __u8 *blocked = bpf_map_lookup_elem(&blocked_uids, &uid);
    if (blocked && *blocked == 1) {
        bpf_printk("🛑 BLOCKED: UID=%d PID=%d tried to openat\n", uid, pid);
        // Force syscall to fail with EPERM
        return -EPERM;
    }

    return 0; // allow syscall
}
