package main

import (
	"github.com/cilium/ebpf"
	"os"
)

func loadBpfObjects(objPath string) (*ebpf.Collection, error) {
	spec, err := ebpf.LoadCollectionSpec(objPath)
	if err != nil {
		return nil, err
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		return nil, err
	}
	return coll, nil
}

func addUIDToBlockMap(coll *ebpf.Collection, uid uint32) error {
	m := coll.Maps["blocked_uids"]
	if m == nil {
		return os.ErrNotExist
	}
	val := uint8(1)
	return m.Put(uid, val)
}
