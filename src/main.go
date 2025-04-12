package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/cilium/ebpf/link"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./syscall-blocker <bpf_obj_file> <uid>")
		os.Exit(1)
	}

	bpfObj := os.Args[1]
	uid, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid UID: %v", err)
	}

	// Load the compiled BPF object
	coll, err := loadBpfObjects(bpfObj)
	if err != nil {
		log.Fatalf("Error loading BPF object: %v", err)
	}
	defer coll.Close()

	// Attach to kprobe for __x64_sys_openat
	kp, err := link.Kprobe("__x64_sys_openat", coll.Programs["block_openat"], nil)
	if err != nil {
		log.Fatalf("Failed to attach kprobe: %v", err)
	}
	defer kp.Close()

	// Insert UID into the block map
	if err := addUIDToBlockMap(coll, uint32(uid)); err != nil {
		log.Fatalf("Error adding UID to map: %v", err)
	}

	log.Printf("🚫 Blocking 'openat' syscall for UID %d using kprobe", uid)

	// Wait for Ctrl+C to exit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("🛑 Received signal. Exiting...")
}
