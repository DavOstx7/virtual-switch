package device

import (
	"fmt"
	"sync"
)

type MACTable struct {
	entries       map[string]*VirtualPort // MAC Address to Virtual Ports
	outputChanges bool
	mu            sync.Mutex
}

func NewMACTable(outputChanges bool) *MACTable {
	return &MACTable{
		entries:       make(map[string]*VirtualPort),
		outputChanges: outputChanges,
	}
}

func (mt *MACTable) LearnMAC(sourceMAC string, vPort *VirtualPort) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if mt.outputChanges {
		if vExistingPort, ok := mt.entries[sourceMAC]; ok && vExistingPort == vPort {
			return
		}
		mt.entries[sourceMAC] = vPort
		fmt.Println(mt.unsafeString())
	} else {
		mt.entries[sourceMAC] = vPort
	}
}

func (mt *MACTable) LookupPort(destinationMAC string) *VirtualPort {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if vPort, ok := mt.entries[destinationMAC]; ok {
		return vPort
	} else {
		return nil
	}
}

func (mt *MACTable) String() string {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	return mt.unsafeString()
}

func (mt *MACTable) unsafeString() string {
	var result string

	border := "+----------------------+------------+"

	result += border

	result += fmt.Sprintf("\n| %-20s | %-10s |\n", "MAC Address", "Ports")
	result += "|----------------------|------------|\n"

	for mac, vPort := range mt.entries {
		result += fmt.Sprintf("| %-20s | %-10s |\n", mac, vPort.Name())
	}

	result += border

	return result
}
