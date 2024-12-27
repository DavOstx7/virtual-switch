package stubs

/*
type Frame interface {
	Bytes() []byte
	SourceMAC() string
	DestinationMAC() string
}
*/

type FrameStub struct {
	bytes          []byte
	sourceMAC      string
	destinationMAC string
}

func (f *FrameStub) Bytes() []byte {
	return f.bytes
}

func (f *FrameStub) SourceMAC() string {
	return f.sourceMAC
}

func (f *FrameStub) DestinationMAC() string {
	return f.destinationMAC
}
