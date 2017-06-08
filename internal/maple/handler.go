package maple

import (
	"errors"
	"sync"

	"golang.ssttevee.com/willow/internal/common"
)

var (
	mutex    sync.RWMutex
	handlers = make(map[common.OpCode][]func(common.ReadablePacket, common.PacketWriter) error)
)

var ErrUnhandledPacket = errors.New("Unhandled Packet")

func handlePacket(packet common.ReadablePacket, writer common.PacketWriter) error {
	mutex.RLock()
	handlers, ok := handlers[packet.OpCode()]
	mutex.RUnlock()

	if !ok {
		return ErrUnhandledPacket
	}

	for _, handler := range handlers {
		if err := handler(packet, writer); err != nil {
			return err
		}
	}

	return nil
}

func RegisterPacketHandler(opcode common.OpCode, handler func(common.ReadablePacket, common.PacketWriter) error) {
	mutex.Lock()
	handlers[opcode] = append(handlers[opcode], handler)
	mutex.Unlock()
}
