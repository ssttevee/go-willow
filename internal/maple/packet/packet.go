package packet

import (
	"bytes"
	"encoding/binary"

	"golang.ssttevee.com/willow/internal/common"
)

type sendingPacket struct {
	buf    bytes.Buffer
	opcode common.OpCode
}

func NewPacketWriter(opcode common.OpCode) common.WritablePacket {
	return &sendingPacket{
		opcode: opcode,
	}
}

func (p *sendingPacket) OpCode() common.OpCode {
	return p.opcode
}

func (p *sendingPacket) Bytes() []byte {
	ret := make([]byte, 2)
	binary.LittleEndian.PutUint16(ret, uint16(p.opcode))

	return append(ret, p.buf.Bytes()...)
}

func (p *sendingPacket) Write(b []byte) (int, error) {
	return p.buf.Write(b)
}

func (p *sendingPacket) WriteShort(short uint16) error {
	binary.Write(&p.buf, binary.LittleEndian, short)

	return nil
}

func (p *sendingPacket) WriteString(str string) error {
	patchBytes := []byte(str)
	binary.Write(&p.buf, binary.LittleEndian, uint16(len(patchBytes)))
	p.buf.Write(patchBytes)

	return nil
}

type RawPacket []byte

func (p RawPacket) OpCode() uint16 {
	return 0
}

func (p RawPacket) Bytes() []byte {
	return []byte(p)
}
