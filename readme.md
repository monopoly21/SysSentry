# System Call Blocker

A lightweight eBPF-based tool to intercept and block specific system calls (e.g., openat) triggered by particular users or mount namespace IDs. Useful for security, sandboxing, or monitoring.

---

## Objective

- **Intercept** syscalls at the kernel level using eBPF (via kprobes).
- **Filter** based on:
  - User IDs (UIDs)
  - Mount namespace IDs (extension possible)
- **Block** the syscall if filters match.
- **Log** attempts with bpf_printk, viewable in trace_pipe.

---

## Repository Structure

```
.
├── README.md
├── Makefile
├── scripts/
│   └── test.sh          # Shell script for testing
└── src/
    ├── main.go          # User-space control logic (Go)
    ├── loader.go        # eBPF object loader + map functions
    └── bpf/
        ├── bpf_prog.c   # eBPF kprobe program
        ├── bpf_prog.h   # Shared map definitions
        └── vmlinux.h    # Auto-generated from bpftool
```

---

## Prerequisites

- **Linux kernel 5.8+** (for modern eBPF support)
- **clang**, **llvm** (to compile BPF programs)
- **bpftool**, **libbpf-dev** (to manage eBPF objects)
- **Go 1.18+** (for user-space logic)
- **make**

### Install on Debian/Ubuntu

```bash
sudo apt update && sudo apt install -y \
    clang llvm libbpf-dev bpftool \
    golang-go make
```

---

## Setup & Build

1. **Clone** the repository:

   ```bash
   git clone <repo-url> && cd syscall-blocker
   ```

2. **Generate `vmlinux.h` (if not present):**

   ```bash
   bpftool btf dump file /sys/kernel/btf/vmlinux format c > src/bpf/vmlinux.h
   ```

3. **Compile:**

   ```bash
   make
   ```

   If successful, you’ll get:
   - `src/bpf/bpf_prog.o`
   - `./syscall-blocker`

---

## Usage

1. **Run** the syscall blocker for your current UID:

   ```bash
   sudo ./syscall-blocker src/bpf/bpf_prog.o $(id -u)
   ```

   This attaches a **kprobe** to `__x64_sys_openat`, blocking the syscall for that UID.

2. **Test** in a new terminal:

   ```bash
   touch /tmp/blockedfile
   ```

   If blocked, you’ll see:

   ```bash
   touch: cannot touch '/tmp/blockedfile': Permission denied
   ```

3. **Logs** from `bpf_printk` are visible in:

   ```bash
   sudo cat /sys/kernel/debug/tracing/trace_pipe
   ```

   e.g.:

   ```bash
   🛑 BLOCKED: UID=1000 attempted openat
   ```

4. **Stop** the blocker by pressing `Ctrl+C` in its terminal.

---

## `scripts/test.sh`

A sample `test.sh` might:

```bash
#!/bin/bash

UID=$(id -u)

echo "[*] Starting the blocker for UID=$UID..."
sudo ./syscall-blocker src/bpf/bpf_prog.o $UID &

PID=$!

sleep 2

echo "[*] Testing file creation..."
touch /tmp/blocktest
ls -l /tmp/blocktest

echo "[*] Killing blocker..."
kill $PID
```

Feel free to customize it.

---

## Extending

- **Mount Namespace Filtering**: Add a second map keyed by namespace ID or combine key `(uid, mnt_ns)`.
- **Additional Syscalls**: Attach kprobes for `execve`, `unlink`, etc.
- **Dynamic Control**: Add CLI commands to add/remove UIDs or syscalls at runtime.
- **Container Ops**: List and kill containers (Docker CLI, Podman, etc.) for deeper sandboxing.

---

## License

MIT – or per your assignment’s requirement.

**Happy coding!** If you have any questions or want advanced features, feel free to reach out.
