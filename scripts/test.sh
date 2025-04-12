#!/bin/bash

echo "[*] Testing syscall blocker..."
uid=$(id -u)

echo "[*] Running syscall-blocker for UID=$uid"
sudo ./syscall-blocker src/bpf/bpf_prog.o $uid &
pid=$!
sleep 2

echo "[*] Trying to open file (should fail if blocked)..."
touch /tmp/testblockfile 2>/dev/null

kill $pid
