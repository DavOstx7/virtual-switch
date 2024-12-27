package net

type Frame interface {
	Bytes() []byte
	SourceMAC() string
	DestinationMAC() string
}
