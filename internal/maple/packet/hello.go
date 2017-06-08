package packet

import (
	"bytes"
	"encoding/binary"

	"golang.ssttevee.com/willow/internal/common"
)

type helloWorldPacket []byte

func (packet helloWorldPacket) OpCode() common.OpCode {
	return common.OpCode(0xFFFF)
}

func (packet helloWorldPacket) Bytes() []byte {
	return []byte(packet)[:]
}

func NewHelloPacket(version uint16, patch string, ivRecv, ivSend []byte) common.Packet {
	var hello bytes.Buffer

	// append version number
	binary.Write(&hello, binary.LittleEndian, version)

	// append patch str
	patchBytes := []byte(patch)
	binary.Write(&hello, binary.LittleEndian, uint16(len(patchBytes)))
	hello.Write(patchBytes)

	// append IVs
	hello.Write(ivRecv)
	hello.Write(ivSend)

	// append locale
	binary.Write(&hello, binary.LittleEndian, uint16(8))

	// not sure what everything else is...

	//binary.Write(&hello, binary.LittleEndian, version)
	//binary.Write(&hello, binary.LittleEndian, version)
	//binary.Write(&hello, binary.LittleEndian, uint16(0))
	//hello.Write(ivRecv)
	//hello.Write(ivSend)
	//
	//hello.Write([]byte{
	//	8, 2, 0, 0,
	//	0, 2, 0, 0,
	//	0, 0, 0, 0,
	//	0, 1, 0,
	//})

	// get packet length
	packetLength := make([]byte, 2)
	binary.LittleEndian.PutUint16(packetLength, uint16(hello.Len()))

	// put the packet length in front of the packet
	return helloWorldPacket(append(packetLength, hello.Bytes()...))
}
