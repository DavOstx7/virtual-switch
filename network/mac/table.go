package mac

import (
	"fmt"
	"sync"
)

type Table struct {
	entries       map[string]string // MAC Address to Port Name
	outputChanges bool
	mu            sync.Mutex
}

func NewTable(outputChanges bool) *Table {
	return &Table{
		entries:       make(map[string]string),
		outputChanges: outputChanges,
	}
}

func (mt *Table) LearnMAC(sourceMAC string, destinationPort string) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if mt.outputChanges {
		if existingPort, ok := mt.entries[sourceMAC]; ok && existingPort == destinationPort {
			return
		}
		mt.entries[sourceMAC] = destinationPort
		fmt.Println(mt._string())
	} else {
		mt.entries[sourceMAC] = destinationPort
	}
}

func (t *Table) LookupPort(destinationMAC string) (string, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	portName, exists := t.entries[destinationMAC]
	return portName, exists
}

func (t *Table) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t._string()
}

func (t *Table) _string() string {
	var result string

	border := "+----------------------+------------+"

	result += border

	result += fmt.Sprintf("\n| %-20s | %-10s |\n", "MAC Address", "Ports")
	result += "|----------------------|------------|\n"

	for mac, port := range t.entries {
		result += fmt.Sprintf("| %-20s | %-10s |\n", mac, port)
	}

	result += border

	return result
}
