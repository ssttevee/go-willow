package handlers

import (
	"golang.ssttevee.com/willow/internal/common"
	"golang.ssttevee.com/willow/internal/maple"
	"golang.ssttevee.com/willow/internal/maple/packet"
)

func init() {
	maple.RegisterPacketHandler(
		common.RecvOpCode_Unknown_0067,
		func(_ common.ReadablePacket, writer common.PacketWriter) error {
			p := packet.NewPacketWriter(common.SendOpCode_Unknown_0017)
			if _, err := p.Write([]byte{0, 0, 0, 0, 0}); err != nil {
				return err
			}

			return writer.Write(p)
		},
	)

	maple.RegisterPacketHandler(
		common.RecvOpCode_Unknown_0098,
		func(_ common.ReadablePacket, writer common.PacketWriter) error {
			p := packet.NewPacketWriter(common.SendOpCode_Unknown_0026)
			if _, err := p.Write([]byte{0}); err != nil {
				return err
			}

			return writer.Write(p)
		},
	)

	maple.RegisterPacketHandler(
		common.RecvOpCode_Unknown_00B1,
		func(_ common.ReadablePacket, writer common.PacketWriter) error {
			p := packet.NewPacketWriter(common.SendOpCode_Unknown_0036)
			if _, err := p.Write([]byte{0x70, 0x96, 0x0B, 0xCE, 0xBB}); err != nil {
				return err
			}

			if _, err := p.Write([]byte{0xDF, 0xD2, 0x01}); err != nil {
				return err
			}

			if err := p.WriteString("MapleLogin1"); err != nil {
				return err
			}

			return writer.Write(p)
		},
	)

	// handle client launch data
	maple.RegisterPacketHandler(
		common.RecvOpCode_Unknown_009A,
		func(_ common.ReadablePacket, writer common.PacketWriter) error {
			// do nothing
			return nil
		},
	)

	// handle pong
	maple.RegisterPacketHandler(
		common.RecvOpCode_Unknown_0093,
		func(_ common.ReadablePacket, writer common.PacketWriter) error {
			// do nothing
			return nil
		},
	)
}
