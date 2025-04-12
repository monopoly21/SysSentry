#ifndef __BPF_PROG_H__
#define __BPF_PROG_H__

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, __u32);   // UID
    __type(value, __u8);  // 1 = block
} blocked_uids SEC(".maps");

#endif
