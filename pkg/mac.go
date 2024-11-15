package pkg

import (
	"sync"
)

type MACTable struct {
	entries map[string]*VirtualPort // MAC Address to Virtual Ports
	mu      sync.Mutex
}

func NewMACTable() *MACTable {
	return &MACTable{
		entries: make(map[string]*VirtualPort),
	}
}

func (mt *MACTable) UpdateSourceEntry(sourceMAC string, vPort *VirtualPort) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	// No need to check if it already exists, as we would want to override it anyways
	mt.entries[sourceMAC] = vPort
}

func (mt *MACTable) GetDestinationPort(destinationMAC string) *VirtualPort {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if vPort, ok := mt.entries[destinationMAC]; ok {
		return vPort
	} else {
		return nil
	}
}
