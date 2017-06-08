package common

import "io"

type Packet interface {
	OpCode() OpCode
	Bytes() []byte
}

type ReadablePacket interface {
	Packet

	io.Reader

	ReadShort() (uint16, error)
	ReadString() (string, error)
}

type WritablePacket interface {
	Packet

	io.Writer

	WriteShort(uint16) error
	WriteString(string) error
}

type PacketWriter interface {
	Write(Packet) error
}

type PacketHandler func(ReadablePacket, PacketWriter) error

type PacketWriterFunc func(Packet) error

func (write PacketWriterFunc) Write(p Packet) error {
	return write(p)
}
