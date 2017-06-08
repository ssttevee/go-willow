package codec

import (
	"encoding/binary"
	"io"

	"golang.ssttevee.com/willow/internal/common"
)

type receivedPacket struct {
	opcode common.OpCode
	data   []byte
	pos    int
}

func newReceivedPacket(data []byte) common.ReadablePacket {
	return &receivedPacket{
		opcode: common.OpCode(binary.LittleEndian.Uint16(data[:2])),
		data:   data[:],
		pos:    2,
	}
}

func (p *receivedPacket) OpCode() common.OpCode {
	return p.opcode
}

func (p *receivedPacket) Bytes() []byte {
	return p.data[:]
}

func (p *receivedPacket) Read(b []byte) (n int, err error) {
	n = len(p.data) - p.pos
	if len(b) < n {
		n = len(b)
	}

	for i := 0; i < n; i++ {
		b[i] = p.data[p.pos+i]
	}

	p.pos += n

	if n < len(b) {
		return n, io.EOF
	}

	return
}

func (p *receivedPacket) ReadShort() (short uint16, err error) {
	short = binary.LittleEndian.Uint16(p.data[p.pos:2])
	p.pos += 2
	return
}

func (p *receivedPacket) ReadString() (str string, err error) {
	l, err := p.ReadShort()
	if err != nil {
		return "", err
	}

	str = string(p.data[p.pos:l])
	p.pos += int(l)
	return
}
