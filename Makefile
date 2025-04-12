BPF_CLANG_FLAGS = -O2 -g -Wall -target bpf -D__TARGET_ARCH_x86

all: syscall-blocker

src/bpf/bpf_prog.o: src/bpf/bpf_prog.c src/bpf/bpf_prog.h
	clang $(BPF_CLANG_FLAGS) -c $< -o $@

syscall-blocker: src/main.go src/loader.go src/bpf/bpf_prog.o
	go build -o syscall-blocker src/main.go src/loader.go

clean:
	rm -f src/bpf/*.o syscall-blocker
